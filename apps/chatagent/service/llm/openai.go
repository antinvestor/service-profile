package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Completer posts OpenAI-compatible chat completions for a single endpoint/key.
type Completer struct {
	BaseURL         string
	CompletionsPath string
	APIKey          string
	Model           string
	HTTPClient      *http.Client
}

// New creates a Completer with the OpenAI-style path (/v1/chat/completions).
// httpClient may be nil (uses a 90s timeout client).
func New(baseURL, apiKey, model string, httpClient *http.Client) *Completer {
	return NewWithPath(baseURL, PathOpenAICompletions, apiKey, model, httpClient)
}

// NewWithPath creates a Completer with an explicit completions path
// (e.g. Google OpenAI-compat uses /chat/completions).
func NewWithPath(baseURL, completionsPath, apiKey, model string, httpClient *http.Client) *Completer {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 90 * time.Second}
	}
	if strings.TrimSpace(completionsPath) == "" {
		completionsPath = PathOpenAICompletions
	}
	return &Completer{
		BaseURL:         strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		CompletionsPath: normalizePath(completionsPath),
		APIKey:          apiKey,
		Model:           model,
		HTTPClient:      httpClient,
	}
}

// NewFromCandidate builds a Completer for one failover candidate.
func NewFromCandidate(c Candidate, httpClient *http.Client) *Completer {
	return NewWithPath(c.BaseURL, c.CompletionsPath, c.APIKey, c.Model, httpClient)
}

func normalizePath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return PathOpenAICompletions
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return p
}

// CompletionsURL returns the full POST URL.
func (c *Completer) CompletionsURL() string {
	if c == nil {
		return ""
	}
	return c.BaseURL + c.CompletionsPath
}

// Complete implements engine.Completer.
func (c *Completer) Complete(ctx context.Context, prompt string) (string, error) {
	if c == nil || c.BaseURL == "" {
		return "", fmt.Errorf("llm: no base URL")
	}
	body := map[string]any{
		"model": c.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"stream":          false,
		"max_tokens":      4096,
		"response_format": map[string]string{"type": "json_object"},
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("llm: marshal: %w", err)
	}
	url := c.CompletionsURL()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("llm: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm: do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", &HTTPStatusError{StatusCode: resp.StatusCode, Body: string(b)}
	}
	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if derr := json.NewDecoder(resp.Body).Decode(&out); derr != nil {
		return "", fmt.Errorf("llm: decode: %w", derr)
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("llm: empty choices")
	}
	return out.Choices[0].Message.Content, nil
}
