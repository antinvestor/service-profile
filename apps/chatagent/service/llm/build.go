package llm

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Config is the inference configuration used to build a Completer.
// Mirrors chatagent env settings without importing the config package
// (keeps llm unit-testable in isolation).
type Config struct {
	Provider string
	BaseURL  string
	Model    string
	APIKey   string // legacy single key
	APIKeys  string // multi-key preferred

	SecondaryProvider string
	SecondaryBaseURL  string
	SecondaryModel    string
	SecondaryAPIKey   string
	SecondaryAPIKeys  string

	// FailoverCooldown is a Go duration string (e.g. "2m"). Empty → DefaultCooldown.
	FailoverCooldown string
}

// BuildCompleter constructs a sticky failover completer from config.
// When inference is not configured, returns a nil completer and a nil error
// (evidence-only mode). Invalid config returns a non-nil error.
func BuildCompleter(cfg Config, httpClient *http.Client) (*FailoverCompleter, error) {
	candidates, err := CandidatesFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		// Evidence-only mode: nil completer is intentional, not an error.
		return nil, nil //nolint:nilnil // disabled inference is not a failure
	}
	cooldown := DefaultCooldown
	if s := strings.TrimSpace(cfg.FailoverCooldown); s != "" {
		d, perr := time.ParseDuration(s)
		if perr != nil {
			return nil, fmt.Errorf("llm: INFERENCE_FAILOVER_COOLDOWN: %w", perr)
		}
		if d > 0 {
			cooldown = d
		}
	}
	return NewFailover(candidates, httpClient, WithCooldown(cooldown))
}

// CandidatesFromConfig expands env-style config into an ordered candidate list.
// Empty result means inference is disabled (not an error).
func CandidatesFromConfig(cfg Config) ([]Candidate, error) {
	primaryKeys := MergeKeys(cfg.APIKeys, cfg.APIKey)
	secondaryKeys := MergeKeys(cfg.SecondaryAPIKeys, cfg.SecondaryAPIKey)

	primaryProvider, err := resolvePrimaryProvider(cfg.Provider, cfg.BaseURL, primaryKeys)
	if err != nil {
		return nil, err
	}

	var slots []SlotConfig

	if primaryProvider != "" && len(primaryKeys) > 0 {
		model := strings.TrimSpace(cfg.Model)
		if model == "" {
			return nil, fmt.Errorf("llm: primary model is required when API keys are set")
		}
		// custom with no base URL already fails in BuildCandidates via ResolveBaseURL
		slots = append(slots, SlotConfig{
			Provider: primaryProvider,
			BaseURL:  cfg.BaseURL,
			Model:    model,
			Keys:     primaryKeys,
			Label:    "primary",
		})
	} else if strings.TrimSpace(cfg.BaseURL) != "" && len(primaryKeys) == 0 {
		// Base URL without keys: single anonymous candidate (some gateways allow this).
		model := strings.TrimSpace(cfg.Model)
		if model == "" {
			return nil, fmt.Errorf("llm: model is required when INFERENCE_BASE_URL is set")
		}
		prov := primaryProvider
		if prov == "" {
			prov = ProviderCustom
		}
		slots = append(slots, SlotConfig{
			Provider: prov,
			BaseURL:  cfg.BaseURL,
			Model:    model,
			Keys:     []string{""}, // empty key — Completer omits Authorization
			Label:    "primary",
		})
	}

	secProvRaw := strings.TrimSpace(cfg.SecondaryProvider)
	secBase := strings.TrimSpace(cfg.SecondaryBaseURL)
	// Secondary is optional: provider/base without keys is ignored (allows
	// Cloud Run to declare secondary env before SM keys are seeded).
	// Keys without provider/base remain an error.
	if len(secondaryKeys) > 0 {
		secModel := strings.TrimSpace(cfg.SecondaryModel)
		if secModel == "" {
			// Fall back to primary model when secondary model omitted.
			secModel = strings.TrimSpace(cfg.Model)
		}
		if secModel == "" {
			return nil, fmt.Errorf("llm: secondary model is required when secondary API keys are set")
		}
		secProv, perr := ParseProvider(secProvRaw)
		if perr != nil {
			return nil, fmt.Errorf("llm: secondary: %w", perr)
		}
		if secProvRaw == "" && secBase != "" {
			secProv = ProviderCustom
		}
		if secProv == ProviderCustom && secBase == "" {
			return nil, fmt.Errorf("llm: secondary custom provider requires base URL")
		}
		if secProvRaw == "" && secBase == "" {
			return nil, fmt.Errorf("llm: secondary API keys require INFERENCE_SECONDARY_PROVIDER or INFERENCE_SECONDARY_BASE_URL")
		}
		slots = append(slots, SlotConfig{
			Provider: secProv,
			BaseURL:  cfg.SecondaryBaseURL,
			Model:    secModel,
			Keys:     secondaryKeys,
			Label:    "secondary",
		})
	}
	// Secondary provider/base without keys is ignored (primary-only until keys are seeded).

	return BuildCandidates(slots...)
}

// resolvePrimaryProvider picks the primary provider enum.
// Returns ("", nil) when nothing is configured (disabled inference).
func resolvePrimaryProvider(provider, baseURL string, keys []string) (Provider, error) {
	provider = strings.TrimSpace(provider)
	baseURL = strings.TrimSpace(baseURL)

	if provider == "" && baseURL == "" && len(keys) == 0 {
		return "", nil
	}

	if provider != "" {
		p, err := ParseProvider(provider)
		if err != nil {
			return "", err
		}
		if p == ProviderCustom && baseURL == "" && len(keys) > 0 {
			// custom + keys without base is invalid unless they only meant openai/google.
			return "", fmt.Errorf("llm: custom provider requires INFERENCE_BASE_URL")
		}
		return p, nil
	}

	// No provider string: infer.
	if baseURL != "" {
		return ProviderCustom, nil
	}
	// Keys without provider or base — not enough to call an API.
	if len(keys) > 0 {
		return "", fmt.Errorf("llm: API keys set but INFERENCE_PROVIDER or INFERENCE_BASE_URL is required")
	}
	return "", nil
}
