package handlers

import (
	"strings"
	"testing"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"github.com/stretchr/testify/require"
)

func TestPickProfileEmail_PrefersVerified(t *testing.T) {
	t.Parallel()
	contacts := []*profilev1.ContactObject{
		{Type: profilev1.ContactType_EMAIL, Detail: "unverified@example.com", Verified: false},
		{Type: profilev1.ContactType_EMAIL, Detail: "verified@example.com", Verified: true},
		{Type: profilev1.ContactType_MSISDN, Detail: "+254700", Verified: true},
	}
	require.Equal(t, "verified@example.com", pickProfileEmail(contacts))
}

func TestPickProfileEmail_FallbackUnverified(t *testing.T) {
	t.Parallel()
	contacts := []*profilev1.ContactObject{
		{Type: profilev1.ContactType_EMAIL, Detail: " only@example.com ", Verified: false},
	}
	require.Equal(t, "only@example.com", pickProfileEmail(contacts))
}

func TestGravatarURLForEmail(t *testing.T) {
	t.Parallel()
	require.Empty(t, gravatarURLForEmail("", 80))
	url := gravatarURLForEmail("  Test@Example.COM  ", 80)
	require.True(t, strings.HasPrefix(url, "https://www.gravatar.com/avatar/"))
	require.Contains(t, url, "s=80")
	require.Contains(t, url, "d=identicon")
	// SHA-256 of "test@example.com"
	require.Contains(t, url, "973dfe463ec85785f5f95af5ba3906eedb2d931c24e69824a89ea65dba4e813b")
}
