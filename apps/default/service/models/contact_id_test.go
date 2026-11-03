package models_test

import (
	"testing"

	"github.com/pitabwire/util"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)

func TestContactIDFromLookupToken(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")

	first := models.ContactIDFromLookupToken(util.ComputeLookupToken(key, "+256757546244"))
	again := models.ContactIDFromLookupToken(util.ComputeLookupToken(key, "+256757546244"))
	other := models.ContactIDFromLookupToken(util.ComputeLookupToken(key, "+256757546245"))
	otherKey := models.ContactIDFromLookupToken(
		util.ComputeLookupToken([]byte("another-lookup-key-another-key!!"), "+256757546244"))

	require.Equal(t, first, again, "the same detail must give the same id")
	require.NotEqual(t, first, other, "a different detail must give a different id")
	require.NotEqual(t, first, otherKey, "the id depends on the lookup key, not the detail alone")

	require.Len(t, first, len(util.IDString()), "the id is as long as a frame xid")
	_, err := xid.FromString(first)
	require.NoError(t, err, "the id is a valid xid")

	require.Empty(t, models.ContactIDFromLookupToken([]byte("short")))
	require.Empty(t, models.ContactIDFromLookupToken(nil))
}
