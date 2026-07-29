package config

import "github.com/pitabwire/frame/v2/config"

// ChatAgentConfig holds env configuration for the chatagent app.
type ChatAgentConfig struct {
	config.ConfigurationDefault

	SecurelyRunService bool `default:"true" envconfig:"SECURELY_RUN_SERVICE"`

	// OpenAI-compatible inference (optional — evidence-only mode when empty).
	InferenceBaseURL string `envconfig:"INFERENCE_BASE_URL"`
	InferenceAPIKey  string `envconfig:"INFERENCE_API_KEY"`
	InferenceModel   string `envconfig:"INFERENCE_MODEL" default:"meta/llama-3.1-8b-instruct"`
}
