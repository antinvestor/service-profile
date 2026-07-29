// Package engine is a product-agnostic conversational data-collection tool.
//
// Products supply only a ContextDef (required fields + guidance). The engine
// evaluates all evidence already available (seed fields, documents such as CV
// text, prior conversation) before asking the user for anything new.
package engine

import (
	"context"
	"encoding/json"
)

// Completer is the minimal LLM surface. Nil disables AI and uses empty extract + guided reply.
type Completer interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

// FieldType classifies values in Fields maps.
type FieldType string

const (
	FieldString     FieldType = "string"
	FieldNumber     FieldType = "number"
	FieldStringList FieldType = "string_list"
	FieldBool       FieldType = "bool"
	FieldObject     FieldType = "object"
)

// FieldDef describes one piece of information to collect.
type FieldDef struct {
	Name          string    `json:"name"`
	Type          FieldType `json:"type"`
	Required      bool      `json:"required"`
	Priority      int       `json:"priority"`
	Description   string    `json:"description,omitempty"`
	Enum          []string  `json:"enum,omitempty"`
	MinLength     int       `json:"min_length,omitempty"`
	Ask           string    `json:"ask,omitempty"`
	Why           string    `json:"why,omitempty"`
	EvidenceHints []string  `json:"evidence_hints,omitempty"`
}

// ReplyPolicy controls assistant wording.
type ReplyPolicy struct {
	MaxSentences      int    `json:"max_sentences,omitempty"`
	AskOneMissingOnly bool   `json:"ask_one_missing_only"`
	CompleteMessage   string `json:"complete_message,omitempty"`
}

// ContextDef is the only product-specific configuration.
type ContextDef struct {
	Key          string      `json:"context_key"`
	Purpose      string      `json:"purpose"`
	SystemPrompt string      `json:"system_prompt,omitempty"`
	Fields       []FieldDef  `json:"fields"`
	ReplyPolicy  ReplyPolicy `json:"reply_policy"`
	ExtractRules string      `json:"extract_rules,omitempty"`
}

// Document is prior material to evaluate before asking (CV, notes, form dump).
type Document struct {
	Name string `json:"name,omitempty"`
	Kind string `json:"kind,omitempty"`
	Text string `json:"text"`
}

// Message is one transcript turn.
type Message struct {
	Role    string `json:"role"` // user | assistant
	Content string `json:"content"`
}

// Fields is a generic structured map (product-agnostic).
type Fields map[string]any

// Evidence is everything already known about the subject.
type Evidence struct {
	// SeedFields are answers already in the system (draft, profile props).
	SeedFields Fields `json:"seed_fields,omitempty"`
	// Documents are long-form evidence (CV text, uploads).
	Documents []Document `json:"documents,omitempty"`
	// Messages are prior conversation (excluding the latest user message if separate).
	Messages []Message `json:"messages,omitempty"`
	// Structured are non-free-text answers from this turn (widgets).
	Structured Fields `json:"structured,omitempty"`
}

// TurnInput is one collection step.
type TurnInput struct {
	// Message is free-text from the user; empty is OK when documents/structured only.
	Message string
	// Evidence is re-evaluated every turn (full known state).
	Evidence Evidence
	// PriorFields is the session field map before this turn.
	PriorFields Fields
}

// FieldStatus is per-field readiness.
type FieldStatus struct {
	OK     bool   `json:"ok"`
	Value  string `json:"value,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// TurnResult is the outcome of one evaluation/turn.
type TurnResult struct {
	Fields      Fields                 `json:"fields"`
	Missing     []string               `json:"missing"`
	Ready       bool                   `json:"ready"`
	FieldStatus map[string]FieldStatus `json:"field_status"`
	Reply       string                 `json:"reply"`
	// Source: evidence | llm | heuristic | llm+heuristic | evidence+llm
	Source   string    `json:"source"`
	Messages []Message `json:"messages,omitempty"`
}

// CloneFields deep-copies a Fields map via JSON (safe for nested values).
func CloneFields(in Fields) Fields {
	if in == nil {
		return Fields{}
	}
	b, err := json.Marshal(in)
	if err != nil {
		out := make(Fields, len(in))
		for k, v := range in {
			out[k] = v
		}
		return out
	}
	out := Fields{}
	_ = json.Unmarshal(b, &out)
	if out == nil {
		return Fields{}
	}
	return out
}
