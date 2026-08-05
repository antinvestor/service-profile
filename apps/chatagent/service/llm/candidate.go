package llm

import (
	"fmt"
	"strings"
	"time"
)

// Provider identifies a known inference backend.
type Provider string

const (
	ProviderOpenAI Provider = "openai"
	ProviderGoogle Provider = "google"
	ProviderCustom Provider = "custom"
)

// Default bases and completions paths for OpenAI-compatible surfaces.
const (
	DefaultOpenAIBaseURL = "https://api.openai.com"
	DefaultGoogleBaseURL = "https://generativelanguage.googleapis.com/v1beta/openai"

	PathOpenAICompletions = "/v1/chat/completions"
	PathGoogleCompletions = "/chat/completions"
)

// Candidate is one (provider, model, key) slot in the failover chain.
type Candidate struct {
	Provider        Provider
	BaseURL         string
	CompletionsPath string
	Model           string
	APIKey          string
	// Label is a short human label for logs (e.g. "primary", "secondary").
	Label string
}

// KeyFingerprint returns a non-secret identifier for logs (last 4 runes).
func KeyFingerprint(apiKey string) string {
	r := []rune(strings.TrimSpace(apiKey))
	if len(r) == 0 {
		return "(empty)"
	}
	if len(r) <= 4 {
		return "…" + string(r)
	}
	return "…" + string(r[len(r)-4:])
}

// String returns a safe summary for logs (no full key).
func (c Candidate) String() string {
	return fmt.Sprintf("%s model=%s key=%s", c.Provider, c.Model, KeyFingerprint(c.APIKey))
}

// ParseProvider normalizes a provider string.
func ParseProvider(s string) (Provider, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", string(ProviderCustom):
		return ProviderCustom, nil
	case string(ProviderOpenAI):
		return ProviderOpenAI, nil
	case string(ProviderGoogle):
		return ProviderGoogle, nil
	default:
		return "", fmt.Errorf("llm: unknown provider %q (want openai, google, or custom)", s)
	}
}

// ResolveBaseURL returns base URL and completions path for a provider slot.
// Explicit baseURL overrides the provider default when non-empty.
func ResolveBaseURL(provider Provider, baseURL string) (base string, path string, err error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	switch provider {
	case ProviderOpenAI:
		if baseURL == "" {
			baseURL = DefaultOpenAIBaseURL
		}
		return baseURL, PathOpenAICompletions, nil
	case ProviderGoogle:
		if baseURL == "" {
			baseURL = DefaultGoogleBaseURL
		}
		return baseURL, PathGoogleCompletions, nil
	case ProviderCustom:
		if baseURL == "" {
			return "", "", fmt.Errorf("llm: custom provider requires base URL")
		}
		return baseURL, PathOpenAICompletions, nil
	default:
		return "", "", fmt.Errorf("llm: unknown provider %q", provider)
	}
}

// ParseKeys splits a multi-key string on commas and whitespace, trims,
// drops empties, and de-duplicates while preserving first-seen order.
func ParseKeys(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	// Normalize commas to spaces then fields-split.
	normalized := strings.ReplaceAll(raw, ",", " ")
	parts := strings.Fields(normalized)
	out := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out
}

// MergeKeys prefers multi-key list; falls back to single key.
func MergeKeys(multi, single string) []string {
	keys := ParseKeys(multi)
	if len(keys) > 0 {
		return keys
	}
	return ParseKeys(single)
}

// SlotConfig describes one provider slot (primary or secondary).
type SlotConfig struct {
	Provider Provider
	BaseURL  string
	Model    string
	Keys     []string
	Label    string
}

// BuildCandidates expands slots into ordered candidates (keys-first within each slot).
func BuildCandidates(slots ...SlotConfig) ([]Candidate, error) {
	var out []Candidate
	for _, slot := range slots {
		if len(slot.Keys) == 0 {
			continue
		}
		model := strings.TrimSpace(slot.Model)
		if model == "" {
			return nil, fmt.Errorf("llm: %s slot has keys but empty model", slot.Label)
		}
		base, path, err := ResolveBaseURL(slot.Provider, slot.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("llm: %s slot: %w", slot.Label, err)
		}
		for _, key := range slot.Keys {
			out = append(out, Candidate{
				Provider:        slot.Provider,
				BaseURL:         base,
				CompletionsPath: path,
				Model:           model,
				APIKey:          key,
				Label:           slot.Label,
			})
		}
	}
	return out, nil
}

// DefaultCooldown is used when failover cooldown is unset or invalid.
const DefaultCooldown = 2 * time.Minute
