package business

import (
	"context"
	"crypto/hmac"
	"crypto/sha1" //nolint:gosec // SHA-1 required by TURN REST API spec (RFC draft-uberti-behave-turn-rest-00).
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pitabwire/frame/client"

	"github.com/antinvestor/service-profile/apps/devices/config"
)

const (
	turnTTLMin = 60
	turnTTLMax = 86400

	// maxCloudflareErrorBodyLen caps the size of Cloudflare error response bodies included in error messages.
	maxCloudflareErrorBodyLen = 256
)

// validTURNURLPrefixes are the allowed schemes for ICE server URLs per the WebRTC spec.
//
//nolint:gochecknoglobals // Package-level lookup table for URL validation.
var validTURNURLPrefixes = []string{"turn:", "turns:", "stun:"}

// ICEServer represents a single ICE server entry compatible with RTCPeerConnection.
type ICEServer struct {
	URLs       []string `json:"urls"`
	Username   string   `json:"username,omitempty"`
	Credential string   `json:"credential,omitempty"`
}

// TURNCredentialsResponse is the response returned to clients for WebRTC peer connection setup.
type TURNCredentialsResponse struct {
	ICEServers []ICEServer `json:"iceServers"`
}

// TURNBusiness handles TURN credential generation.
type TURNBusiness interface {
	GetTurnCredentials(ctx context.Context) (*TURNCredentialsResponse, error)
}

// TURNProvider is a pluggable credential generator.
type TURNProvider interface {
	GenerateCredentials(ctx context.Context, ttl int) (*TURNCredentialsResponse, error)
}

type turnBusiness struct {
	cfg      *config.DevicesConfig
	provider TURNProvider
}

// NewTURNBusiness creates a TURNBusiness using the provider configured via TURN_PROVIDER.
// Returns (nil, nil) when TURN is intentionally disabled (empty provider with no config).
func NewTURNBusiness(cfg *config.DevicesConfig, clientMgr client.Manager) (TURNBusiness, error) {
	// Validate and clamp TTL.
	if cfg.TURNTTL < turnTTLMin || cfg.TURNTTL > turnTTLMax {
		return nil, fmt.Errorf("TURN_TTL must be between %d and %d, got %d", turnTTLMin, turnTTLMax, cfg.TURNTTL)
	}

	var provider TURNProvider

	switch strings.ToLower(strings.TrimSpace(cfg.TURNProvider)) {
	case "cloudflare":
		tokenID := strings.TrimSpace(cfg.CloudflareTURNTokenID)
		apiToken := strings.TrimSpace(cfg.CloudflareTURNAPIToken)
		if tokenID == "" || apiToken == "" {
			return nil, errors.New(
				"cloudflare TURN provider requires CLOUDFLARE_TURN_TOKEN_ID and CLOUDFLARE_TURN_API_TOKEN",
			)
		}
		provider = newCloudflareTURNProvider(tokenID, apiToken, clientMgr, cfg.TURNTTL)
	case "static", "":
		sharedSecret := strings.TrimSpace(cfg.TURNSharedSecret)
		if sharedSecret == "" || cfg.TURNServerURLs == "" {
			return nil, errors.New("static TURN provider requires TURN_SHARED_SECRET and TURN_SERVER_URLS")
		}
		urls, err := parseTURNURLs(cfg.TURNServerURLs)
		if err != nil {
			return nil, err
		}
		hmacAlg := strings.ToLower(strings.TrimSpace(cfg.TURNHMACAlgorithm))
		var newHash func() hash.Hash
		switch hmacAlg {
		case "sha1", "":
			newHash = sha1.New
		case "sha256":
			newHash = sha256.New
		default:
			return nil, fmt.Errorf("unsupported TURN_HMAC_ALGORITHM: %q (supported: sha1, sha256)", hmacAlg)
		}
		provider = &staticTURNProvider{
			sharedSecret: sharedSecret,
			serverURLs:   urls,
			newHash:      newHash,
		}
	default:
		return nil, fmt.Errorf("unknown TURN provider: %q (supported: cloudflare, static)", cfg.TURNProvider)
	}

	return &turnBusiness{cfg: cfg, provider: provider}, nil
}

func (t *turnBusiness) GetTurnCredentials(ctx context.Context) (*TURNCredentialsResponse, error) {
	return t.provider.GenerateCredentials(ctx, t.cfg.TURNTTL)
}

// staticTURNProvider generates time-limited credentials using HMAC over a shared secret.
// This is the standard mechanism used by coturn (via static-auth-secret) and pion/turn.
// See: https://datatracker.ietf.org/doc/html/draft-uberti-behave-turn-rest-00
type staticTURNProvider struct {
	sharedSecret string
	serverURLs   []string
	newHash      func() hash.Hash
}

func (s *staticTURNProvider) GenerateCredentials(_ context.Context, ttl int) (*TURNCredentialsResponse, error) {
	// Username is "expiry_timestamp" — coturn and pion both parse the unix timestamp from the username.
	expiry := time.Now().Unix() + int64(ttl)
	username := strconv.FormatInt(expiry, 10)

	mac := hmac.New(s.newHash, []byte(s.sharedSecret))
	mac.Write([]byte(username))
	credential := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return &TURNCredentialsResponse{
		ICEServers: []ICEServer{
			{
				URLs:       s.serverURLs,
				Username:   username,
				Credential: credential,
			},
		},
	}, nil
}

// cloudflareTURNProvider generates credentials via Cloudflare's TURN API.
// Responses are cached for (TTL - safety margin) to avoid hitting the API on every request.
type cloudflareTURNProvider struct {
	endpointURL string
	authHeader  string // pre-computed "Bearer <token>"
	client      client.Manager
	ttl         int

	mu          sync.Mutex
	cached      *TURNCredentialsResponse
	cachedUntil time.Time
}

func newCloudflareTURNProvider(tokenID, apiToken string, cli client.Manager, ttl int) *cloudflareTURNProvider {
	return &cloudflareTURNProvider{
		endpointURL: fmt.Sprintf(
			"https://rtc.live.cloudflare.com/v1/turn/keys/%s/credentials/generate-ice-servers",
			tokenID,
		),
		authHeader: "Bearer " + apiToken,
		client:     cli,
		ttl:        ttl,
	}
}

func (c *cloudflareTURNProvider) GenerateCredentials(ctx context.Context, ttl int) (*TURNCredentialsResponse, error) {
	// Return cached credentials if still valid.
	c.mu.Lock()
	if c.cached != nil && time.Now().Before(c.cachedUntil) {
		resp := c.cached
		c.mu.Unlock()
		return resp, nil
	}
	c.mu.Unlock()

	// Fetch fresh credentials from Cloudflare.
	headers := http.Header{}
	headers.Set("Authorization", c.authHeader)

	payload := map[string]int{"ttl": ttl}

	resp, err := c.client.Invoke(ctx, http.MethodPost, c.endpointURL, payload, headers)
	if err != nil {
		return nil, fmt.Errorf("calling cloudflare TURN API: %w", err)
	}
	defer resp.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := resp.ToContent(ctx)
		// Cap the body length to avoid leaking sensitive data in logs.
		bodyStr := string(body)
		if len(bodyStr) > maxCloudflareErrorBodyLen {
			bodyStr = bodyStr[:maxCloudflareErrorBodyLen] + "..."
		}
		return nil, fmt.Errorf("cloudflare TURN API returned status %d: %s", resp.StatusCode, bodyStr)
	}

	var turnResp TURNCredentialsResponse
	if decErr := resp.Decode(ctx, &turnResp); decErr != nil {
		return nil, fmt.Errorf("decoding cloudflare TURN response: %w", decErr)
	}

	// Cache the response. Use 90% of TTL as the cache window to ensure
	// credentials returned to clients always have reasonable validity remaining.
	const cacheWindowFraction = 0.9
	c.mu.Lock()
	c.cached = &turnResp
	c.cachedUntil = time.Now().Add(time.Duration(float64(ttl)*cacheWindowFraction) * time.Second)
	c.mu.Unlock()

	return &turnResp, nil
}

// parseTURNURLs splits a comma-separated list of TURN/STUN URLs, trims whitespace,
// and validates that each URL starts with a valid ICE scheme (turn:, turns:, stun:).
func parseTURNURLs(raw string) ([]string, error) {
	parts := strings.Split(raw, ",")
	urls := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if !hasValidTURNScheme(p) {
			return nil, fmt.Errorf("invalid turn server URL %q: must use turn/turns/stun scheme", p)
		}
		urls = append(urls, p)
	}
	if len(urls) == 0 {
		return nil, errors.New("TURN_SERVER_URLS must contain at least one valid URL")
	}
	return urls, nil
}

func hasValidTURNScheme(url string) bool {
	lower := strings.ToLower(url)
	for _, prefix := range validTURNURLPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}
