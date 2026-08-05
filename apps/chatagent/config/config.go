package config

import "github.com/pitabwire/frame/v2/config"

// ChatAgentConfig holds env configuration for the chatagent app.
type ChatAgentConfig struct {
	config.ConfigurationDefault

	SecurelyRunService bool `default:"true" envconfig:"SECURELY_RUN_SERVICE"`

	// Inference — OpenAI-compatible primary (optional; evidence-only when no candidates).
	//
	// Sticky multi-key failover: INFERENCE_API_KEYS (ordered) are tried for the
	// primary provider/model until a key degrades, then the next key. Optional
	// secondary provider is used only after all primary keys are degraded.
	// Keys are never rotated every request for load balancing.
	InferenceProvider string `envconfig:"INFERENCE_PROVIDER"` // openai | google | custom
	InferenceBaseURL  string `envconfig:"INFERENCE_BASE_URL"`
	InferenceAPIKey   string `envconfig:"INFERENCE_API_KEY"`  // legacy single key
	InferenceAPIKeys  string `envconfig:"INFERENCE_API_KEYS"` // comma/whitespace list, preferred
	InferenceModel    string `envconfig:"INFERENCE_MODEL" default:"meta/llama-3.1-8b-instruct"`

	// Optional secondary provider (used after primary key pool is degraded).
	InferenceSecondaryProvider string `envconfig:"INFERENCE_SECONDARY_PROVIDER"`
	InferenceSecondaryBaseURL  string `envconfig:"INFERENCE_SECONDARY_BASE_URL"`
	InferenceSecondaryAPIKey   string `envconfig:"INFERENCE_SECONDARY_API_KEY"`
	InferenceSecondaryAPIKeys  string `envconfig:"INFERENCE_SECONDARY_API_KEYS"`
	InferenceSecondaryModel    string `envconfig:"INFERENCE_SECONDARY_MODEL"`

	// How long a failed key stays skipped before the primary is preferred again.
	InferenceFailoverCooldown string `envconfig:"INFERENCE_FAILOVER_COOLDOWN" default:"2m"`

	// Notification service — omnichannel outbound for non-web session channels.
	// Empty URI disables delivery (web RPC replies still work).
	NotificationSvcURI                       string `envconfig:"NOTIFICATION_SERVICE_URI"`
	NotificationServiceWorkloadAPITargetPath string `envconfig:"NOTIFICATION_SERVICE_WORKLOAD_API_TARGET_PATH" default:"/ns/notifications/sa/service-notification"`
}
