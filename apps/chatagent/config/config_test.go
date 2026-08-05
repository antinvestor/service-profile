package config_test

import (
	"os"
	"testing"

	"github.com/pitabwire/frame/v2/config"
	"github.com/stretchr/testify/require"

	aconfig "github.com/antinvestor/service-profile/apps/chatagent/config"
)

func TestChatAgentConfig_LoadsInferenceEnv(t *testing.T) {
	t.Setenv("INFERENCE_PROVIDER", "google")
	t.Setenv("INFERENCE_MODEL", "gemini-2.0-flash")
	t.Setenv("INFERENCE_API_KEYS", "key-a,key-b")
	t.Setenv("INFERENCE_SECONDARY_PROVIDER", "openai")
	t.Setenv("INFERENCE_SECONDARY_MODEL", "gpt-4o-mini")
	t.Setenv("INFERENCE_FAILOVER_COOLDOWN", "3m")
	t.Setenv("NOTIFICATION_SERVICE_URI", "https://api.stawi.org/notification")

	// Ensure unrelated required env defaults do not break parse
	_ = os.Getenv("DATABASE_URL")

	cfg, err := config.FromEnv[aconfig.ChatAgentConfig]()
	require.NoError(t, err)
	require.Equal(t, "google", cfg.InferenceProvider)
	require.Equal(t, "gemini-2.0-flash", cfg.InferenceModel)
	require.Equal(t, "key-a,key-b", cfg.InferenceAPIKeys)
	require.Equal(t, "openai", cfg.InferenceSecondaryProvider)
	require.Equal(t, "gpt-4o-mini", cfg.InferenceSecondaryModel)
	require.Equal(t, "3m", cfg.InferenceFailoverCooldown)
	require.Equal(t, "https://api.stawi.org/notification", cfg.NotificationSvcURI)
}
