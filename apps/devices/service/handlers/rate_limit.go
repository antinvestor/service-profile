package handlers

import (
	"net/http"

	"github.com/pitabwire/frame/security"

	"github.com/antinvestor/service-profile/apps/devices/service/caching"
)

const prefixRateTURN = "rate:turn:"

// RateLimitTURN wraps an HTTP handler with per-caller TURN credential rate limiting.
// The caller identity is extracted from the JWT claims (subject).
func RateLimitTURN(
	next http.HandlerFunc,
	cacheSvc *caching.DeviceCacheService,
	limitPerMinute int64,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cacheSvc != nil && limitPerMinute > 0 {
			callerID := "anonymous"
			claims := security.ClaimsFromContext(r.Context())
			if claims != nil {
				if sub, err := claims.GetSubject(); err == nil && sub != "" {
					callerID = sub
				}
			}

			allowed, _ := cacheSvc.CheckRateLimit(r.Context(), prefixRateTURN, callerID, limitPerMinute)
			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "60")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":"TURN credential rate limit exceeded"}`))
				return
			}
		}

		next(w, r)
	}
}
