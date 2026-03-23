package handlers //nolint:testpackage // tests access unexported handler helpers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCleanErrAndValidationHelpers(t *testing.T) {
	t.Parallel()

	server := &GeolocationServer{}
	ctx := context.Background()

	require.True(t, isValidationError(errors.New("invalid route name")))
	require.True(t, isValidationError(errors.New("owner_id is required")))
	require.False(t, isValidationError(nil))
	require.False(t, isValidationError(errors.New("boom")))

	require.Equal(t, connect.CodeNotFound, connect.CodeOf(server.cleanErr(ctx, gorm.ErrRecordNotFound)))
	require.Equal(t, connect.CodeNotFound, connect.CodeOf(server.cleanErr(ctx, errors.New("route not found"))))
	require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(server.cleanErr(ctx, errors.New("invalid area name"))))
	require.Equal(t, connect.CodeInternal, connect.CodeOf(server.cleanErr(ctx, errors.New("boom"))))
}

func TestRateLimiterHelpers(t *testing.T) {
	t.Parallel()

	rl := &rateLimiter{
		buckets: map[string]*bucket{},
		cfg: RateLimiterConfig{
			RequestsPerWindow: 2,
			WindowDuration:    time.Second,
			CleanupInterval:   time.Millisecond,
		},
		trustedNets:    parseTrustedProxies([]string{"203.0.113.0/24", "2001:db8::1", "bad-entry"}),
		cleanupStopped: make(chan struct{}),
	}

	require.Len(t, rl.trustedNets, 2)
	require.True(t, rl.isTrustedProxy(parseTrustedProxies([]string{"203.0.113.10"})[0].IP))
	require.Equal(t, "203.0.113.4", extractIPFromAddr("203.0.113.4:8080"))
	require.Equal(t, "invalid-addr", extractIPFromAddr("invalid-addr"))

	require.True(t, rl.allow("198.51.100.10"))
	require.True(t, rl.allow("198.51.100.10"))
	require.False(t, rl.allow("198.51.100.10"))

	rl.buckets["stale"] = &bucket{tokens: 1, lastReset: time.Now().Add(-3 * time.Second)}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		rl.cleanupLoop(ctx)
	}()
	time.Sleep(5 * time.Millisecond)
	cancel()
	<-done
	_, exists := rl.buckets["stale"]
	require.False(t, exists)

	require.Equal(t, "198.51.100.1", extractRightmostUntrustedIP(
		"198.51.100.1, 203.0.113.5, 203.0.113.6",
		rl,
	))
	require.Equal(t, "203.0.113.5", extractRightmostUntrustedIP("203.0.113.5, 203.0.113.6", rl))
	require.Equal(t, "bad-ip", extractRightmostUntrustedIP("bad-ip", rl))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "203.0.113.9:1234"
	req.Header.Set("X-Forwarded-For", "198.51.100.7, 203.0.113.8")
	require.Equal(t, "198.51.100.7", rl.extractClientIP(req))

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "198.51.100.9:1234"
	req.Header.Set("X-Forwarded-For", "192.0.2.10")
	require.Equal(t, "198.51.100.9", rl.extractClientIP(req))
}
