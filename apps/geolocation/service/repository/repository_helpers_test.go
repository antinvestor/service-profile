package repository //nolint:testpackage // tests access unexported repository helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepositoryHelpers(t *testing.T) {
	t.Parallel()

	require.Equal(t, `100\%\\\_done`, escapeLikeWildcards(`100%\_done`))
	require.Nil(t, applyTimeRange(nil, nil, nil))
}
