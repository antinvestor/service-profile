package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test started when the test binary is started. Only calls main.
func TestSystem(t *testing.T) {
	// assert equality
	require.Equal(t, 123, 123, "Just ensuring we are testing something")
}
