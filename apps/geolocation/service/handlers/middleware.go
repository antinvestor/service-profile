package handlers

import (
	"context"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiterConfig holds configuration for the per-IP rate limiter.
type RateLimiterConfig struct {
	// RequestsPerWindow is the maximum number of requests allowed per window.
	RequestsPerWindow int
	// WindowDuration is the duration of the sliding window.
	WindowDuration time.Duration
	// CleanupInterval is how often stale entries are purged.
	CleanupInterval time.Duration
	// TrustedProxies is a list of CIDR ranges or IPs whose X-Forwarded-For header is trusted.
	// If empty, X-Forwarded-For is never used (only RemoteAddr is used for IP extraction).
	TrustedProxies []string
}

// rateLimiter implements a simple per-IP token bucket rate limiter.
// Uses an in-memory map with periodic cleanup of stale entries.
type rateLimiter struct {
	mu             sync.Mutex
	buckets        map[string]*bucket
	cfg            RateLimiterConfig
	trustedNets    []*net.IPNet
	stopCleanup    context.CancelFunc
	cleanupStopped chan struct{}
}

type bucket struct {
	tokens    int
	lastReset time.Time
}

// Rate limiter defaults.
const (
	defaultRequestsPerWindow         = 600
	defaultCleanupIntervalMultiplier = 5

	// CIDR mask sizes for single-IP trusted proxy entries.
	ipv4MaskBits = 32
	ipv6MaskBits = 128
)

// RateLimitResult holds the middleware and stop function returned by NewRateLimitMiddleware.
type RateLimitResult struct {
	Middleware func(http.Handler) http.Handler
	Stop       func()
}

// NewRateLimitMiddleware creates an HTTP middleware that enforces per-IP rate limiting.
// If cfg is nil or has zero values, sensible defaults are used.
// Call Stop during graceful shutdown to stop the cleanup goroutine.
func NewRateLimitMiddleware(cfg *RateLimiterConfig) *RateLimitResult {
	c := RateLimiterConfig{
		RequestsPerWindow: defaultRequestsPerWindow,
		WindowDuration:    time.Minute,
		CleanupInterval:   defaultCleanupIntervalMultiplier * time.Minute,
	}
	if cfg != nil {
		if cfg.RequestsPerWindow > 0 {
			c.RequestsPerWindow = cfg.RequestsPerWindow
		}
		if cfg.WindowDuration > 0 {
			c.WindowDuration = cfg.WindowDuration
		}
		if cfg.CleanupInterval > 0 {
			c.CleanupInterval = cfg.CleanupInterval
		}
		c.TrustedProxies = cfg.TrustedProxies
	}

	ctx, cancel := context.WithCancel(context.Background())
	rl := &rateLimiter{
		buckets:        make(map[string]*bucket),
		cfg:            c,
		trustedNets:    parseTrustedProxies(c.TrustedProxies),
		stopCleanup:    cancel,
		cleanupStopped: make(chan struct{}),
	}

	// Start background cleanup goroutine.
	go rl.cleanupLoop(ctx)

	return &RateLimitResult{
		Middleware: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ip := rl.extractClientIP(r)
				if !rl.allow(ip) {
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("Retry-After", "60")
					w.WriteHeader(http.StatusTooManyRequests)
					_, _ = w.Write([]byte(`{"error":"rate limit exceeded"}`))
					return
				}
				next.ServeHTTP(w, r)
			})
		},
		Stop: func() {
			rl.stopCleanup()
			<-rl.cleanupStopped
		},
	}
}

// allow checks if the IP is allowed to make a request.
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.buckets[ip]
	if !exists || now.Sub(b.lastReset) >= rl.cfg.WindowDuration {
		rl.buckets[ip] = &bucket{
			tokens:    rl.cfg.RequestsPerWindow - 1,
			lastReset: now,
		}
		return true
	}

	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

// cleanupLoop periodically removes stale buckets. Stops when ctx is cancelled.
func (rl *rateLimiter) cleanupLoop(ctx context.Context) {
	defer close(rl.cleanupStopped)

	ticker := time.NewTicker(rl.cfg.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rl.mu.Lock()
			cutoff := time.Now().Add(-2 * rl.cfg.WindowDuration)
			for ip, b := range rl.buckets {
				if b.lastReset.Before(cutoff) {
					delete(rl.buckets, ip)
				}
			}
			rl.mu.Unlock()
		}
	}
}

// extractClientIP extracts the client IP from the request.
// Only trusts X-Forwarded-For if the direct connection comes from a trusted proxy.
// Otherwise, uses RemoteAddr to prevent IP spoofing via header injection.
func (rl *rateLimiter) extractClientIP(r *http.Request) string {
	remoteIP := extractIPFromAddr(r.RemoteAddr)

	// Only trust X-Forwarded-For when the direct peer is a configured trusted proxy.
	if len(rl.trustedNets) > 0 {
		parsed := net.ParseIP(remoteIP)
		if parsed != nil && rl.isTrustedProxy(parsed) {
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				// Take the rightmost untrusted IP (the one added by the last trusted proxy).
				return extractRightmostUntrustedIP(xff, rl)
			}
		}
	}

	return remoteIP
}

// extractIPFromAddr extracts the IP portion from an addr:port string.
func extractIPFromAddr(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}

// isTrustedProxy checks if the given IP is within any of the trusted proxy networks.
func (rl *rateLimiter) isTrustedProxy(ip net.IP) bool {
	for _, n := range rl.trustedNets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

// extractRightmostUntrustedIP walks the X-Forwarded-For chain from right to left
// and returns the first IP that is NOT a trusted proxy. This is the client IP
// as seen by the outermost trusted proxy.
func extractRightmostUntrustedIP(xff string, rl *rateLimiter) string {
	parts := strings.Split(xff, ",")
	for i := len(parts) - 1; i >= 0; i-- {
		ip := strings.TrimSpace(parts[i])
		parsed := net.ParseIP(ip)
		if parsed == nil {
			continue
		}
		if !rl.isTrustedProxy(parsed) {
			return ip
		}
	}
	// All IPs in the chain are trusted — use the leftmost (original client).
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return ""
}

// parseTrustedProxies converts a list of CIDR or IP strings into net.IPNet objects.
// Single IPs are treated as /32 (IPv4) or /128 (IPv6).
func parseTrustedProxies(proxies []string) []*net.IPNet {
	var nets []*net.IPNet
	for _, p := range proxies {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		// Try as CIDR first.
		_, ipNet, err := net.ParseCIDR(p)
		if err == nil {
			nets = append(nets, ipNet)
			continue
		}
		// Try as single IP.
		ip := net.ParseIP(p)
		if ip == nil {
			continue
		}
		mask := net.CIDRMask(ipv6MaskBits, ipv6MaskBits)
		if ip.To4() != nil {
			mask = net.CIDRMask(ipv4MaskBits, ipv4MaskBits)
		}
		nets = append(nets, &net.IPNet{IP: ip, Mask: mask})
	}
	return nets
}
