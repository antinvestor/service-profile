package llm

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pitabwire/util"
)

// clock is injectable for tests.
type clock func() time.Time

// FailoverCompleter selects the highest-priority healthy candidate and sticks
// to it until it degrades. It never rotates healthy keys for load balancing.
// After cooldown, higher-priority keys are preferred again.
type FailoverCompleter struct {
	candidates []Candidate
	httpClient *http.Client
	cooldown   time.Duration
	now        clock

	mu            sync.Mutex
	degradedUntil []time.Time
}

// FailoverOption configures FailoverCompleter.
type FailoverOption func(*FailoverCompleter)

// WithCooldown sets how long a key stays degraded after a degradable failure.
func WithCooldown(d time.Duration) FailoverOption {
	return func(f *FailoverCompleter) {
		if d > 0 {
			f.cooldown = d
		}
	}
}

// WithClock injects a time source (tests).
func WithClock(fn clock) FailoverOption {
	return func(f *FailoverCompleter) {
		if fn != nil {
			f.now = fn
		}
	}
}

// NewFailover builds a sticky failover completer. candidates must be non-empty
// and ordered highest priority first.
func NewFailover(candidates []Candidate, httpClient *http.Client, opts ...FailoverOption) (*FailoverCompleter, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("llm: no candidates")
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 90 * time.Second}
	}
	f := &FailoverCompleter{
		candidates:    append([]Candidate(nil), candidates...),
		httpClient:    httpClient,
		cooldown:      DefaultCooldown,
		now:           time.Now,
		degradedUntil: make([]time.Time, len(candidates)),
	}
	for _, opt := range opts {
		opt(f)
	}
	if f.cooldown <= 0 {
		f.cooldown = DefaultCooldown
	}
	return f, nil
}

// Complete implements engine.Completer.
func (f *FailoverCompleter) Complete(ctx context.Context, prompt string) (string, error) {
	if f == nil || len(f.candidates) == 0 {
		return "", fmt.Errorf("llm: no candidates")
	}

	log := util.Log(ctx)
	order := f.attemptOrder()
	var attempts []*AttemptError

	for _, idx := range order {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		c := f.candidates[idx]
		client := NewFromCandidate(c, f.httpClient)
		content, err := client.Complete(ctx, prompt)
		if err == nil {
			log.Debug("chatagent llm: completion ok",
				"provider", string(c.Provider),
				"model", c.Model,
				"key", KeyFingerprint(c.APIKey),
				"label", c.Label,
				"candidate_index", idx,
			)
			return content, nil
		}

		class := Classify(ctx, err)
		ae := &AttemptError{
			Index:      idx,
			Candidate:  c,
			StatusCode: statusCodeFrom(err),
			Class:      class,
			Err:        err,
		}
		attempts = append(attempts, ae)

		if class == ClassHard {
			log.WithError(err).Warn("chatagent llm: hard error; not trying other keys",
				"provider", string(c.Provider),
				"model", c.Model,
				"key", KeyFingerprint(c.APIKey),
				"label", c.Label,
				"status", ae.StatusCode,
				"candidate_index", idx,
			)
			// If this is the only attempt, return it directly; else wrap.
			if len(attempts) == 1 {
				return "", ae
			}
			return "", &MultiError{Attempts: attempts}
		}

		until := f.markDegraded(idx)
		log.WithError(err).Warn("chatagent llm: candidate degraded; trying next",
			"provider", string(c.Provider),
			"model", c.Model,
			"key", KeyFingerprint(c.APIKey),
			"label", c.Label,
			"status", ae.StatusCode,
			"candidate_index", idx,
			"degraded_until", until.UTC().Format(time.RFC3339),
			"cooldown", f.cooldown.String(),
		)
	}

	if len(attempts) == 0 {
		return "", fmt.Errorf("llm: no candidates attempted")
	}
	log.Error("chatagent llm: all candidates failed", "attempt_count", len(attempts))
	return "", &MultiError{Attempts: attempts}
}

// attemptOrder returns candidate indices to try.
// Prefer healthy (not in cooldown) keys in priority order.
// If none are healthy, probe all in priority order (last-resort availability).
func (f *FailoverCompleter) attemptOrder() []int {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := f.now()
	n := len(f.candidates)
	healthy := make([]int, 0, n)
	for i := range n {
		until := f.degradedUntil[i]
		if until.IsZero() || !now.Before(until) {
			// Cooldown expired: clear so next classify starts clean.
			if !until.IsZero() && !now.Before(until) {
				f.degradedUntil[i] = time.Time{}
			}
			healthy = append(healthy, i)
		}
	}
	if len(healthy) > 0 {
		return healthy
	}
	// All cooling down: still probe every key once in priority order.
	all := make([]int, n)
	for i := range all {
		all[i] = i
	}
	return all
}

func (f *FailoverCompleter) markDegraded(idx int) time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	until := f.now().Add(f.cooldown)
	f.degradedUntil[idx] = until
	return until
}

// DegradedUntil returns the degrade deadline for candidate i (tests / diagnostics).
func (f *FailoverCompleter) DegradedUntil(idx int) time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	if idx < 0 || idx >= len(f.degradedUntil) {
		return time.Time{}
	}
	return f.degradedUntil[idx]
}

// CandidateCount returns how many candidates are configured.
func (f *FailoverCompleter) CandidateCount() int {
	if f == nil {
		return 0
	}
	return len(f.candidates)
}
