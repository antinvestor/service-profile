package llm

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// Class is how a completion error should affect failover.
type Class int

const (
	// ClassHard means do not try other keys (bad request, caller cancel, etc.).
	ClassHard Class = iota
	// ClassDegradable means mark the candidate degraded and try the next key.
	ClassDegradable
)

// AttemptError is a failure for one candidate (safe for logs).
type AttemptError struct {
	Index      int
	Candidate  Candidate
	StatusCode int
	Class      Class
	Err        error
}

func (e *AttemptError) Error() string {
	if e == nil {
		return "llm: attempt error"
	}
	status := e.StatusCode
	if status == 0 {
		return fmt.Sprintf("llm: %s [%s] key=%s: %v",
			e.Candidate.Label, e.Candidate.Provider, KeyFingerprint(e.Candidate.APIKey), e.Err)
	}
	return fmt.Sprintf("llm: %s [%s] key=%s status=%d: %v",
		e.Candidate.Label, e.Candidate.Provider, KeyFingerprint(e.Candidate.APIKey), status, e.Err)
}

func (e *AttemptError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// MultiError aggregates all failed attempts when no candidate succeeds.
type MultiError struct {
	Attempts []*AttemptError
}

func (e *MultiError) Error() string {
	if e == nil || len(e.Attempts) == 0 {
		return "llm: all candidates failed"
	}
	parts := make([]string, 0, len(e.Attempts))
	for _, a := range e.Attempts {
		parts = append(parts, a.Error())
	}
	return "llm: all candidates failed: " + strings.Join(parts, "; ")
}

func (e *MultiError) Unwrap() []error {
	if e == nil {
		return nil
	}
	out := make([]error, 0, len(e.Attempts))
	for _, a := range e.Attempts {
		out = append(out, a)
	}
	return out
}

// HTTPStatusError is returned by Completer on non-200 responses.
type HTTPStatusError struct {
	StatusCode int
	Body       string
}

func (e *HTTPStatusError) Error() string {
	if e == nil {
		return "llm: http status error"
	}
	return fmt.Sprintf("llm: status %d: %s", e.StatusCode, e.Body)
}

// Classify maps an error to hard vs degradable for failover.
// callerCtx is the request context; when it is canceled/deadline-exceeded and
// err wraps that, classification is hard so we do not burn other keys.
func Classify(callerCtx context.Context, err error) Class {
	if err == nil {
		return ClassHard
	}

	// Caller abandoned the request — do not fan out.
	if callerCtx != nil && callerCtx.Err() != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, callerCtx.Err()) {
			return ClassHard
		}
	}
	if errors.Is(err, context.Canceled) {
		// Cancel without caller ctx cancel usually means transport closed; still hard if pure cancel.
		if callerCtx == nil || callerCtx.Err() != nil {
			return ClassHard
		}
	}

	var httpErr *HTTPStatusError
	if errors.As(err, &httpErr) {
		return classifyStatus(httpErr.StatusCode)
	}

	// Network / timeout style failures.
	var netErr net.Error
	if errors.As(err, &netErr) {
		return ClassDegradable
	}
	if errors.Is(err, context.DeadlineExceeded) {
		// Independent of caller: client-side HTTP timeout while provider hung.
		return ClassDegradable
	}

	msg := strings.ToLower(err.Error())
	// Empty choices / decode after transport success — treat as flaky provider.
	if strings.Contains(msg, "empty choices") ||
		strings.Contains(msg, "decode") ||
		strings.Contains(msg, "connection") ||
		strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "tls") ||
		strings.Contains(msg, "reset by peer") ||
		strings.Contains(msg, "broken pipe") {
		return ClassDegradable
	}

	// Misconfiguration at call site (no base URL) is hard.
	if strings.Contains(msg, "no base url") || strings.Contains(msg, "no candidates") {
		return ClassHard
	}

	// Default: degrade and try next (prefer availability for unknown transport issues).
	return ClassDegradable
}

func classifyStatus(code int) Class {
	switch code {
	case http.StatusTooManyRequests, // 429
		http.StatusUnauthorized, // 401 — try next key / secondary provider
		http.StatusForbidden,    // 403
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return ClassDegradable
	case http.StatusBadRequest:
		// Invalid API key often returns 400 with a body message — still try
		// the next key or secondary provider rather than hard-stopping the pool.
		return ClassDegradable
	case http.StatusNotFound,
		http.StatusRequestEntityTooLarge,
		http.StatusUnprocessableEntity:
		return ClassHard
	default:
		if code >= 500 && code <= 599 {
			return ClassDegradable
		}
		if code >= 400 && code <= 499 {
			return ClassHard
		}
		return ClassDegradable
	}
}

// statusCodeFrom extracts HTTP status if present.
func statusCodeFrom(err error) int {
	var httpErr *HTTPStatusError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode
	}
	return 0
}
