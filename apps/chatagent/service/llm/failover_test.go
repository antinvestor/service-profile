package llm_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/chatagent/service/llm"
)

func chatOK(content string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{"content": content}},
			},
		})
	}
}

func chatStatus(code int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(code)
		_, _ = io.WriteString(w, body)
	}
}

// multiKeyServer routes by Authorization Bearer token so one server can
// represent multiple logical keys.
func multiKeyServer(t *testing.T, handlers map[string]http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		const prefix = "Bearer "
		key := ""
		if len(auth) > len(prefix) {
			key = auth[len(prefix):]
		}
		h, ok := handlers[key]
		if !ok {
			http.Error(w, "unknown key "+key, http.StatusInternalServerError)
			return
		}
		h(w, r)
	}))
}

func candidatesFor(baseURL string, keys ...string) []llm.Candidate {
	out := make([]llm.Candidate, 0, len(keys))
	for _, k := range keys {
		out = append(out, llm.Candidate{
			Provider:        llm.ProviderCustom,
			BaseURL:         baseURL,
			CompletionsPath: llm.PathOpenAICompletions,
			Model:           "test-model",
			APIKey:          k,
			Label:           "primary",
		})
	}
	return out
}

func TestFailover_StickyUsesPrimaryWhileHealthy(t *testing.T) {
	t.Parallel()
	var hits1, hits2 atomic.Int32
	srv := multiKeyServer(t, map[string]http.HandlerFunc{
		"key1": func(w http.ResponseWriter, r *http.Request) {
			hits1.Add(1)
			chatOK(`{"fields":{},"reply":"ok"}`)(w, r)
		},
		"key2": func(w http.ResponseWriter, r *http.Request) {
			hits2.Add(1)
			chatOK(`{"fields":{},"reply":"ok2"}`)(w, r)
		},
	})
	t.Cleanup(srv.Close)

	fc, err := llm.NewFailover(candidatesFor(srv.URL, "key1", "key2"), srv.Client(), llm.WithCooldown(time.Minute))
	require.NoError(t, err)

	ctx := context.Background()
	for range 3 {
		out, cerr := fc.Complete(ctx, "hello")
		require.NoError(t, cerr)
		require.Contains(t, out, "ok")
	}
	require.Equal(t, int32(3), hits1.Load(), "primary key must handle all healthy calls")
	require.Equal(t, int32(0), hits2.Load(), "must not rotate to second key")
}

func TestFailover_SameRequestThenStickyOnBackup(t *testing.T) {
	t.Parallel()
	var hits1, hits2 atomic.Int32
	srv := multiKeyServer(t, map[string]http.HandlerFunc{
		"key1": func(w http.ResponseWriter, r *http.Request) {
			hits1.Add(1)
			chatStatus(http.StatusTooManyRequests, "rate limited")(w, r)
		},
		"key2": func(w http.ResponseWriter, r *http.Request) {
			hits2.Add(1)
			chatOK(`backup-ok`)(w, r)
		},
	})
	t.Cleanup(srv.Close)

	fc, err := llm.NewFailover(candidatesFor(srv.URL, "key1", "key2"), srv.Client(), llm.WithCooldown(time.Hour))
	require.NoError(t, err)

	ctx := context.Background()
	out, err := fc.Complete(ctx, "hello")
	require.NoError(t, err)
	require.Equal(t, "backup-ok", out)
	require.Equal(t, int32(1), hits1.Load())
	require.Equal(t, int32(1), hits2.Load())

	// Next requests stick to key2 while key1 is degraded — no rotation back to key1.
	out, err = fc.Complete(ctx, "again")
	require.NoError(t, err)
	require.Equal(t, "backup-ok", out)
	require.Equal(t, int32(1), hits1.Load(), "degraded primary must not be hit again during cooldown")
	require.Equal(t, int32(2), hits2.Load())
}

func TestFailover_CooldownRestoresPrimary(t *testing.T) {
	t.Parallel()
	var hits1, hits2 atomic.Int32
	var key1Mode atomic.Int32 // 0=429, 1=200
	srv := multiKeyServer(t, map[string]http.HandlerFunc{
		"key1": func(w http.ResponseWriter, r *http.Request) {
			hits1.Add(1)
			if key1Mode.Load() == 0 {
				chatStatus(http.StatusTooManyRequests, "rl")(w, r)
				return
			}
			chatOK("primary-ok")(w, r)
		},
		"key2": func(w http.ResponseWriter, r *http.Request) {
			hits2.Add(1)
			chatOK("backup-ok")(w, r)
		},
	})
	t.Cleanup(srv.Close)

	now := time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)
	var mu sync.Mutex
	fc, err := llm.NewFailover(
		candidatesFor(srv.URL, "key1", "key2"),
		srv.Client(),
		llm.WithCooldown(2*time.Minute),
		llm.WithClock(func() time.Time {
			mu.Lock()
			defer mu.Unlock()
			return now
		}),
	)
	require.NoError(t, err)

	ctx := context.Background()
	out, err := fc.Complete(ctx, "1")
	require.NoError(t, err)
	require.Equal(t, "backup-ok", out)

	// Still in cooldown.
	out, err = fc.Complete(ctx, "2")
	require.NoError(t, err)
	require.Equal(t, "backup-ok", out)
	require.Equal(t, int32(1), hits1.Load())

	// Advance past cooldown; primary healthy again.
	mu.Lock()
	now = now.Add(3 * time.Minute)
	mu.Unlock()
	key1Mode.Store(1)

	out, err = fc.Complete(ctx, "3")
	require.NoError(t, err)
	require.Equal(t, "primary-ok", out)
	require.Equal(t, int32(2), hits1.Load(), "primary preferred after cooldown")
	require.Equal(t, int32(2), hits2.Load())
}

func TestFailover_Hard400DoesNotTryNextKey(t *testing.T) {
	t.Parallel()
	var hits1, hits2 atomic.Int32
	srv := multiKeyServer(t, map[string]http.HandlerFunc{
		"key1": func(w http.ResponseWriter, r *http.Request) {
			hits1.Add(1)
			chatStatus(http.StatusBadRequest, `{"error":"bad json schema"}`)(w, r)
		},
		"key2": func(w http.ResponseWriter, r *http.Request) {
			hits2.Add(1)
			chatOK("should-not-run")(w, r)
		},
	})
	t.Cleanup(srv.Close)

	fc, err := llm.NewFailover(candidatesFor(srv.URL, "key1", "key2"), srv.Client(), llm.WithCooldown(time.Minute))
	require.NoError(t, err)

	_, err = fc.Complete(context.Background(), "hello")
	require.Error(t, err)
	require.Equal(t, int32(1), hits1.Load())
	require.Equal(t, int32(0), hits2.Load(), "hard errors must not fan out to other keys")
}

func TestFailover_CallerCancelIsHard(t *testing.T) {
	t.Parallel()
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		<-r.Context().Done()
	}))
	t.Cleanup(srv.Close)

	fc, err := llm.NewFailover(candidatesFor(srv.URL, "key1", "key2"), srv.Client(), llm.WithCooldown(time.Minute))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = fc.Complete(ctx, "hello")
	require.Error(t, err)
	// May or may not have started HTTP; must not succeed on key2 with canceled ctx.
	require.ErrorIs(t, err, context.Canceled)
}

func TestFailover_AllKeysFail(t *testing.T) {
	t.Parallel()
	srv := multiKeyServer(t, map[string]http.HandlerFunc{
		"key1": chatStatus(http.StatusServiceUnavailable, "down"),
		"key2": chatStatus(http.StatusTooManyRequests, "rl"),
	})
	t.Cleanup(srv.Close)

	fc, err := llm.NewFailover(candidatesFor(srv.URL, "key1", "key2"), srv.Client(), llm.WithCooldown(time.Minute))
	require.NoError(t, err)

	_, err = fc.Complete(context.Background(), "hello")
	require.Error(t, err)
	require.Contains(t, err.Error(), "all candidates failed")
	require.Contains(t, err.Error(), "key1"[len("key1")-4:]) // fingerprint ends with key1
}

func TestFailover_Auth401DegradesToNext(t *testing.T) {
	t.Parallel()
	srv := multiKeyServer(t, map[string]http.HandlerFunc{
		"bad":  chatStatus(http.StatusUnauthorized, "invalid"),
		"good": chatOK("ok"),
	})
	t.Cleanup(srv.Close)

	fc, err := llm.NewFailover(candidatesFor(srv.URL, "bad", "good"), srv.Client(), llm.WithCooldown(time.Hour))
	require.NoError(t, err)

	out, err := fc.Complete(context.Background(), "x")
	require.NoError(t, err)
	require.Equal(t, "ok", out)
}

func TestFailover_SecondaryAfterPrimaryPool(t *testing.T) {
	t.Parallel()
	var pathHits sync.Map
	var googlePath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathHits.Store(r.URL.Path, true)
		auth := r.Header.Get("Authorization")
		switch auth {
		case "Bearer pk1", "Bearer pk2":
			chatStatus(http.StatusTooManyRequests, "rl")(w, r)
		case "Bearer gk1":
			// Google path should be /chat/completions not /v1/...
			googlePath = r.URL.Path
			chatOK("gemini-ok")(w, r)
		default:
			http.Error(w, "bad auth", http.StatusInternalServerError)
		}
	}))
	t.Cleanup(srv.Close)

	cands := []llm.Candidate{
		{Provider: llm.ProviderOpenAI, BaseURL: srv.URL, CompletionsPath: llm.PathOpenAICompletions, Model: "gpt", APIKey: "pk1", Label: "primary"},
		{Provider: llm.ProviderOpenAI, BaseURL: srv.URL, CompletionsPath: llm.PathOpenAICompletions, Model: "gpt", APIKey: "pk2", Label: "primary"},
		{Provider: llm.ProviderGoogle, BaseURL: srv.URL, CompletionsPath: llm.PathGoogleCompletions, Model: "gemini", APIKey: "gk1", Label: "secondary"},
	}
	fc, err := llm.NewFailover(cands, srv.Client(), llm.WithCooldown(time.Hour))
	require.NoError(t, err)

	out, err := fc.Complete(context.Background(), "x")
	require.NoError(t, err)
	require.Equal(t, "gemini-ok", out)
	require.Equal(t, llm.PathGoogleCompletions, googlePath)

	_, okV1 := pathHits.Load(llm.PathOpenAICompletions)
	_, okG := pathHits.Load(llm.PathGoogleCompletions)
	require.True(t, okV1)
	require.True(t, okG)
}

func TestFailover_ConcurrentSafe(t *testing.T) {
	t.Parallel()
	srv := multiKeyServer(t, map[string]http.HandlerFunc{
		"key1": chatOK("ok"),
	})
	t.Cleanup(srv.Close)

	fc, err := llm.NewFailover(candidatesFor(srv.URL, "key1"), srv.Client(), llm.WithCooldown(time.Minute))
	require.NoError(t, err)

	var (
		wg   sync.WaitGroup
		errs atomic.Int32
	)
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, cerr := fc.Complete(context.Background(), "x"); cerr != nil {
				errs.Add(1)
			}
		}()
	}
	wg.Wait()
	require.Equal(t, int32(0), errs.Load())
}

func TestClassify_StatusCodes(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	require.Equal(t, llm.ClassDegradable, llm.Classify(ctx, &llm.HTTPStatusError{StatusCode: 429}))
	require.Equal(t, llm.ClassDegradable, llm.Classify(ctx, &llm.HTTPStatusError{StatusCode: 503}))
	require.Equal(t, llm.ClassDegradable, llm.Classify(ctx, &llm.HTTPStatusError{StatusCode: 401}))
	require.Equal(t, llm.ClassHard, llm.Classify(ctx, &llm.HTTPStatusError{StatusCode: 400}))
	require.Equal(t, llm.ClassHard, llm.Classify(ctx, &llm.HTTPStatusError{StatusCode: 422}))
}

func TestCompleter_UsesConfiguredPath(t *testing.T) {
	t.Parallel()
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		chatOK("x")(w, r)
	}))
	t.Cleanup(srv.Close)

	c := llm.NewWithPath(srv.URL, llm.PathGoogleCompletions, "k", "m", srv.Client())
	_, err := c.Complete(context.Background(), "p")
	require.NoError(t, err)
	require.Equal(t, llm.PathGoogleCompletions, gotPath)
}

func TestBuildCompleter_NilWhenDisabled(t *testing.T) {
	t.Parallel()
	fc, err := llm.BuildCompleter(llm.Config{}, nil)
	require.NoError(t, err)
	require.Nil(t, fc)
}

func TestBuildCompleter_MultiKey(t *testing.T) {
	t.Parallel()
	fc, err := llm.BuildCompleter(llm.Config{
		Provider:         "openai",
		Model:            "gpt-4o-mini",
		APIKeys:          "a,b,c",
		FailoverCooldown: "30s",
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, fc)
	require.Equal(t, 3, fc.CandidateCount())
}
