package config

import "github.com/pitabwire/frame/v2/config"

// ChatAgentConfig holds env configuration for the chatagent app.
// Tags use Frame's caarlos0/env (`env` / `envDefault`), not legacy envconfig.
type ChatAgentConfig struct {
	config.ConfigurationDefault

	SecurelyRunService bool `env:"SECURELY_RUN_SERVICE" envDefault:"true"`

	// Inference — OpenAI-compatible primary (optional; evidence-only when no candidates).
	//
	// Sticky multi-key failover: INFERENCE_API_KEYS (ordered) are tried for the
	// primary provider/model until a key degrades, then the next key. Optional
	// secondary provider is used only after all primary keys are degraded.
	// Keys are never rotated every request for load balancing.
	InferenceProvider string `env:"INFERENCE_PROVIDER"` // openai | google | custom
	InferenceBaseURL  string `env:"INFERENCE_BASE_URL"`
	InferenceAPIKey   string `env:"INFERENCE_API_KEY"`  // legacy single key
	InferenceAPIKeys  string `env:"INFERENCE_API_KEYS"` // comma/whitespace list, preferred
	InferenceModel    string `env:"INFERENCE_MODEL" envDefault:"meta/llama-3.1-8b-instruct"`

	// Optional secondary provider (used after primary key pool is degraded).
	InferenceSecondaryProvider string `env:"INFERENCE_SECONDARY_PROVIDER"`
	InferenceSecondaryBaseURL  string `env:"INFERENCE_SECONDARY_BASE_URL"`
	InferenceSecondaryAPIKey   string `env:"INFERENCE_SECONDARY_API_KEY"`
	InferenceSecondaryAPIKeys  string `env:"INFERENCE_SECONDARY_API_KEYS"`
	InferenceSecondaryModel    string `env:"INFERENCE_SECONDARY_MODEL"`

	// How long a failed key stays skipped before the primary is preferred again.
	InferenceFailoverCooldown string `env:"INFERENCE_FAILOVER_COOLDOWN" envDefault:"2m"`

	// Notification service — omnichannel outbound for non-web session channels.
	// Empty URI disables delivery (web RPC replies still work).
	NotificationSvcURI                       string `env:"NOTIFICATION_SERVICE_URI"`
	NotificationServiceWorkloadAPITargetPath string `env:"NOTIFICATION_SERVICE_WORKLOAD_API_TARGET_PATH" envDefault:"/ns/notifications/sa/service-notification"`
}
