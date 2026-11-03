package business

import (
	"strings"
	"testing"

	"github.com/pitabwire/frame/v3/data"
	"github.com/stretchr/testify/require"
)

func TestNormalizeHandle(t *testing.T) {
	valid := []struct{ in, want string }{
		{"peter.bwire", "peter.bwire"},
		{" @Peter.Bwire ", "peter.bwire"},
		{"user.a1b2c3d4", "user.a1b2c3d4"},
		{"d4ac8fnpf2b6ruj0", "d4ac8fnpf2b6ruj0"},
	}
	for _, tc := range valid {
		got, err := NormalizeHandle(tc.in)
		require.NoError(t, err, tc.in)
		require.Equal(t, tc.want, got, tc.in)
	}
	invalid := []string{
		"", "@", "ab", "peter..bwire", ".peter", "peter.", "peter bwire", "pe%", "a_b_c", strings.Repeat("a", 41),
	}
	for _, in := range invalid {
		_, err := NormalizeHandle(in)
		require.ErrorIs(t, err, ErrInvalidHandle, in)
	}
}

func TestSlugHandle(t *testing.T) {
	for in, want := range map[string]string{
		"Peter Bwire":                        "peter.bwire",
		"  José   Ñúñez-O'Brien ":            "jose.nunez.o.brien",
		"stawi-jobs-api":                     "stawi.jobs.api",
		"Ünïcödé":                            "unicode",
		"李小龙":                                "",
		"a very long display name indeed ok": "a.very.long.display.name",
		"---":                                "",
	} {
		got := slugHandle(in)
		require.Equal(t, want, got, in)
		if got != "" {
			_, err := NormalizeHandle(got)
			require.NoError(t, err, "slug %q of %q must be a valid handle", got, in)
		}
	}
}

func TestHandleBase(t *testing.T) {
	require.Equal(t, "peter.bwire", handleBase(data.JSONMap{"name": "Peter Bwire", "au_name": "pb"}))
	require.Equal(t, "peter.bwire", handleBase(data.JSONMap{"name": "李", "au_name": "Peter Bwire"}))
	require.Equal(t, "user", handleBase(data.JSONMap{"name": "Al"}))
	require.Equal(t, "user", handleBase(data.JSONMap{"name": 42}))
	require.Equal(t, "user", handleBase(nil))
}

func TestHandleCandidates(t *testing.T) {
	id := "d4ac8fnpf2b6ruj0abcd"
	require.Equal(t,
		[]string{"peter.bwire", "peter.bwire.ruj0abcd", id},
		handleCandidates("peter.bwire", id))
	require.Equal(t,
		[]string{"user.ruj0abcd", id},
		handleCandidates(handleFallback, id),
		"a profile without a name never gets a bare user")
	for _, h := range handleCandidates("peter.bwire", id) {
		_, err := NormalizeHandle(h)
		require.NoError(t, err, h)
	}
}
