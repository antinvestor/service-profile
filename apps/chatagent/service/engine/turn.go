package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pitabwire/util"
)

// Agent runs evidence-first data collection turns.
type Agent struct {
	LLM Completer
}

// NewAgent returns an agent. llm may be nil (guided replies without extract).
func NewAgent(llm Completer) *Agent {
	return &Agent{LLM: llm}
}

// Turn evaluates all evidence, extracts/merges fields, assesses readiness, and
// produces a reply that only asks for what is still missing.
func (a *Agent) Turn(ctx context.Context, def ContextDef, in TurnInput) (TurnResult, error) {
	log := util.Log(ctx)

	fields := CloneFields(in.PriorFields)
	fields = MergeFields(fields, in.Evidence.SeedFields)
	fields = ApplyDocuments(def, fields, in.Evidence.Documents)
	fields = ApplyStructured(fields, in.Evidence.Structured)
	fields = Sanitize(def, fields)

	source := "evidence"
	_, missingBefore, readyBefore := Assess(def, fields)

	var llmReply string
	msg := strings.TrimSpace(in.Message)

	// Run LLM extract when we have any signal (message, docs, or incomplete seed).
	// When inference is configured, failures are hard errors — never fall through
	// to canned/template guided replies that look like a successful assistant turn.
	needExtract := a.LLM != nil && (msg != "" || len(in.Evidence.Documents) > 0 || !readyBefore || len(in.Evidence.Messages) > 0)
	if needExtract {
		prompt := BuildExtractPrompt(def, fields, missingBefore, in.Evidence, msg)
		raw, err := a.LLM.Complete(ctx, prompt)
		if err != nil {
			log.WithError(err).Error("chatagent: LLM extract failed")
			return TurnResult{}, fmt.Errorf("chatagent: llm extract failed: %w", err)
		}
		extracted, reply, perr := parseExtractResponse(raw)
		if perr != nil {
			log.WithError(perr).Error("chatagent: LLM JSON parse failed")
			return TurnResult{}, fmt.Errorf("chatagent: llm response unusable: %w", perr)
		}
		fields = MergeFields(fields, extracted)
		fields = Sanitize(def, fields)
		llmReply = strings.TrimSpace(reply)
		if source == "evidence" {
			source = "evidence+llm"
		} else {
			source = "llm"
		}
	}

	// Re-apply structured/documents so model cannot wipe trusted inputs.
	fields = ApplyDocuments(def, fields, in.Evidence.Documents)
	fields = ApplyStructured(fields, in.Evidence.Structured)
	fields = Sanitize(def, fields)

	status, missing, ready := Assess(def, fields)
	reply := ComposeReply(def, fields, missing, ready, llmReply)

	// Build transcript append for caller convenience.
	var msgs []Message
	if msg != "" {
		msgs = append(msgs, Message{Role: "user", Content: truncateRunes(msg, 12_000)})
	} else if len(in.Evidence.Documents) > 0 {
		// Synthetic user turn for document-only evaluation.
		names := make([]string, 0, len(in.Evidence.Documents))
		for _, d := range in.Evidence.Documents {
			names = append(names, firstNonEmpty(d.Name, d.Kind, "document"))
		}
		msgs = append(msgs, Message{
			Role:    "user",
			Content: "I've provided: " + strings.Join(names, ", "),
		})
	}
	if reply != "" {
		msgs = append(msgs, Message{Role: "assistant", Content: truncateRunes(reply, 4_000)})
	}

	return TurnResult{
		Fields:      fields,
		Missing:     missing,
		Ready:       ready,
		FieldStatus: status,
		Reply:       reply,
		Source:      source,
		Messages:    msgs,
	}, nil
}

// Evaluate runs a turn with no new user message — used at CreateSession when
// seed fields / documents should populate readiness before the first chat turn.
func (a *Agent) Evaluate(ctx context.Context, def ContextDef, prior Fields, evidence Evidence) (TurnResult, error) {
	return a.Turn(ctx, def, TurnInput{
		Message:     "",
		Evidence:    evidence,
		PriorFields: prior,
	})
}

func parseExtractResponse(raw string) (Fields, string, error) {
	raw = extractJSONObject(raw)
	if raw == "" {
		return nil, "", fmt.Errorf("empty llm response")
	}
	var parsed struct {
		Fields Fields `json:"fields"`
		Reply  string `json:"reply"`
	}
	dec := json.NewDecoder(strings.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&parsed); err != nil {
		return nil, "", fmt.Errorf("decode llm json: %w", err)
	}
	if parsed.Fields == nil {
		parsed.Fields = Fields{}
	}
	return parsed.Fields, parsed.Reply, nil
}

func extractJSONObject(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	// Strip markdown fences if present.
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```json")
		s = strings.TrimPrefix(s, "```JSON")
		s = strings.TrimPrefix(s, "```")
		if i := strings.LastIndex(s, "```"); i >= 0 {
			s = s[:i]
		}
		s = strings.TrimSpace(s)
	}
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}

// ComposeReply prefers the model reply. It never invents a canned “Got it…”
// intake script when the model did not produce wording — callers must fail the
// turn if extract failed. Guided follow-up is only used when inference is
// disabled (llmReply empty and no LLM was expected).
func ComposeReply(def ContextDef, fields Fields, missing []string, ready bool, llmReply string) string {
	llmReply = strings.TrimSpace(llmReply)
	if ready {
		if llmReply != "" && !looksLikeAskingForMore(llmReply) {
			return llmReply
		}
		if def.ReplyPolicy.CompleteMessage != "" {
			return def.ReplyPolicy.CompleteMessage
		}
		if llmReply != "" {
			return llmReply
		}
		return "Thanks — I have everything I need for now."
	}
	// Prefer authentic model text whenever present and not a false “ready” claim.
	if llmReply != "" {
		if looksLikeFalseReady(llmReply) {
			// Strip false completion claims; keep any steering if present.
			if looksLikeAskingForMore(llmReply) {
				return llmReply
			}
			// Fall through only when model falsely claimed done with no next ask.
		} else {
			return llmReply
		}
	}
	// No model wording: structural guided ask (evidence-only / no-inference mode).
	return guidedFollowUp(def, fields, missing)
}

func guidedFollowUp(def ContextDef, fields Fields, missing []string) string {
	if len(missing) == 0 {
		return "Thanks — I have everything I need for now."
	}
	fd := fieldByName(def, missing[0])
	ask := firstNonEmpty(fd.Ask, fd.Description, "Could you share "+missing[0]+"?")
	// Brief acknowledgment of what we already have.
	var known []string
	for _, f := range def.Fields {
		if !f.Required {
			continue
		}
		st := assessOne(f, fields)
		if st.OK && st.Value != "" && len(known) < 2 {
			known = append(known, f.Name)
		}
	}
	if len(known) > 0 {
		return fmt.Sprintf("Got it so far (%s). %s", strings.Join(known, ", "), ask)
	}
	return ask
}

func looksLikeAskingForMore(s string) bool {
	low := strings.ToLower(s)
	return strings.Contains(low, "?") ||
		strings.Contains(low, "which ") ||
		strings.Contains(low, "what ") ||
		strings.Contains(low, "where ") ||
		strings.Contains(low, "could you") ||
		strings.Contains(low, "can you") ||
		strings.Contains(low, "please ")
}

func looksLikeFalseReady(s string) bool {
	low := strings.ToLower(s)
	return strings.Contains(low, "all set") ||
		strings.Contains(low, "everything i need") ||
		strings.Contains(low, "i have everything") ||
		strings.Contains(low, "we're ready") ||
		strings.Contains(low, "you are all set")
}

func replyTargetsMissing(reply, missingKey string, def ContextDef) bool {
	low := strings.ToLower(reply)
	fd := fieldByName(def, missingKey)
	// Soft match on field name tokens and ask text.
	for _, token := range strings.Fields(strings.ReplaceAll(missingKey, "_", " ")) {
		if len(token) >= 3 && strings.Contains(low, strings.ToLower(token)) {
			return true
		}
	}
	if fd.Ask != "" {
		for _, token := range strings.Fields(fd.Ask) {
			t := strings.ToLower(strings.Trim(token, "?,.!"))
			if len(t) >= 4 && strings.Contains(low, t) {
				return true
			}
		}
	}
	return false
}
