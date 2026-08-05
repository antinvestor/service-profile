package llm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/chatagent/service/llm"
)

func TestParseKeys_CommaAndWhitespace(t *testing.T) {
	t.Parallel()
	keys := llm.ParseKeys("  k1, k2  k3,,k2 ,k4 ")
	require.Equal(t, []string{"k1", "k2", "k3", "k4"}, keys)
}

func TestParseKeys_Empty(t *testing.T) {
	t.Parallel()
	require.Empty(t, llm.ParseKeys(""))
	require.Empty(t, llm.ParseKeys("  ,  "))
}

func TestMergeKeys_PrefersMulti(t *testing.T) {
	t.Parallel()
	require.Equal(t, []string{"a", "b"}, llm.MergeKeys("a,b", "legacy"))
	require.Equal(t, []string{"legacy"}, llm.MergeKeys("", "legacy"))
}

func TestKeyFingerprint(t *testing.T) {
	t.Parallel()
	require.Equal(t, "(empty)", llm.KeyFingerprint(""))
	require.Equal(t, "…abcd", llm.KeyFingerprint("xxxxabcd"))
	require.Equal(t, "…ab", llm.KeyFingerprint("ab"))
}

func TestResolveBaseURL_Providers(t *testing.T) {
	t.Parallel()

	base, path, err := llm.ResolveBaseURL(llm.ProviderOpenAI, "")
	require.NoError(t, err)
	require.Equal(t, llm.DefaultOpenAIBaseURL, base)
	require.Equal(t, llm.PathOpenAICompletions, path)

	base, path, err = llm.ResolveBaseURL(llm.ProviderGoogle, "")
	require.NoError(t, err)
	require.Equal(t, llm.DefaultGoogleBaseURL, base)
	require.Equal(t, llm.PathGoogleCompletions, path)
	require.NotContains(t, path, "/v1/")

	base, path, err = llm.ResolveBaseURL(llm.ProviderCustom, "https://gw.example/v1/")
	require.NoError(t, err)
	require.Equal(t, "https://gw.example/v1", base)
	require.Equal(t, llm.PathOpenAICompletions, path)

	_, _, err = llm.ResolveBaseURL(llm.ProviderCustom, "")
	require.Error(t, err)
}

func TestBuildCandidates_KeysFirstThenSecondary(t *testing.T) {
	t.Parallel()
	cands, err := llm.BuildCandidates(
		llm.SlotConfig{
			Provider: llm.ProviderOpenAI,
			Model:    "gpt-4o-mini",
			Keys:     []string{"pk1", "pk2"},
			Label:    "primary",
		},
		llm.SlotConfig{
			Provider: llm.ProviderGoogle,
			Model:    "gemini-2.0-flash",
			Keys:     []string{"gk1"},
			Label:    "secondary",
		},
	)
	require.NoError(t, err)
	require.Len(t, cands, 3)
	require.Equal(t, "pk1", cands[0].APIKey)
	require.Equal(t, "pk2", cands[1].APIKey)
	require.Equal(t, "gk1", cands[2].APIKey)
	require.Equal(t, llm.ProviderGoogle, cands[2].Provider)
	require.Equal(t, llm.PathGoogleCompletions, cands[2].CompletionsPath)
	require.Equal(t, llm.PathOpenAICompletions, cands[0].CompletionsPath)
}

func TestCandidatesFromConfig_LegacySingleKey(t *testing.T) {
	t.Parallel()
	cands, err := llm.CandidatesFromConfig(llm.Config{
		BaseURL: "https://llm.example",
		APIKey:  "sk-legacy",
		Model:   "meta/llama",
	})
	require.NoError(t, err)
	require.Len(t, cands, 1)
	require.Equal(t, "sk-legacy", cands[0].APIKey)
	require.Equal(t, llm.ProviderCustom, cands[0].Provider)
}

func TestCandidatesFromConfig_OpenAIMultiKeyPlusSecondary(t *testing.T) {
	t.Parallel()
	cands, err := llm.CandidatesFromConfig(llm.Config{
		Provider:          "openai",
		Model:             "gpt-4o-mini",
		APIKeys:           "k1,k2",
		SecondaryProvider: "google",
		SecondaryModel:    "gemini-2.0-flash",
		SecondaryAPIKeys:  "g1",
	})
	require.NoError(t, err)
	require.Len(t, cands, 3)
	require.Equal(t, llm.DefaultOpenAIBaseURL, cands[0].BaseURL)
	require.Equal(t, llm.DefaultGoogleBaseURL, cands[2].BaseURL)
	require.Equal(t, "g1", cands[2].APIKey)
}

func TestCandidatesFromConfig_Disabled(t *testing.T) {
	t.Parallel()
	cands, err := llm.CandidatesFromConfig(llm.Config{})
	require.NoError(t, err)
	require.Empty(t, cands)
}

func TestCandidatesFromConfig_KeysWithoutProvider(t *testing.T) {
	t.Parallel()
	_, err := llm.CandidatesFromConfig(llm.Config{
		APIKeys: "k1",
		Model:   "x",
	})
	require.Error(t, err)
}

func TestCandidatesFromConfig_SecondaryWithoutKeysIsSkipped(t *testing.T) {
	t.Parallel()
	// Allows Cloud Run to declare secondary provider env before keys are seeded.
	cands, err := llm.CandidatesFromConfig(llm.Config{
		Provider:          "openai",
		APIKeys:           "k1",
		Model:             "gpt",
		SecondaryProvider: "google",
		SecondaryModel:    "gemini-2.0-flash",
	})
	require.NoError(t, err)
	require.Len(t, cands, 1)
	require.Equal(t, "k1", cands[0].APIKey)
}
