package engine

import (
	"encoding/json"
	"fmt"
	"strings"
)

// BuildExtractPrompt constructs a generic, context-driven extraction prompt.
// Product differences live entirely in ContextDef (purpose, fields, extract_rules).
func BuildExtractPrompt(def ContextDef, prior Fields, missing []string, evidence Evidence, latest string) string {
	var b strings.Builder
	purpose := strings.TrimSpace(def.Purpose)
	if purpose == "" {
		purpose = "Collect the required structured information for this intake."
	}
	sys := strings.TrimSpace(def.SystemPrompt)
	if sys != "" {
		b.WriteString(sys)
		b.WriteString("\n\n")
	}

	fmt.Fprintf(&b, `You are a careful data-collection agent.
## Purpose
%s

## Rules (critical)
1. Extract structured fields ONLY from evidence already provided: existing fields, documents, conversation, and the latest message.
2. NEVER invent values. If unsure, leave the field empty.
3. Prefer values already present in documents and seed fields over re-asking.
4. Do not ask for information that is already clearly present in the evidence.
5. reply: brief (1–%d sentences). Acknowledge what you already know. If anything required is still missing, ask ONLY for the single highest-priority missing field. If nothing is missing, confirm completion without inventing extra questions.
6. Merge with existing_fields — do not wipe prior values.

`, purpose, maxSentences(def))

	if rules := strings.TrimSpace(def.ExtractRules); rules != "" {
		b.WriteString("## Product extract rules\n")
		b.WriteString(rules)
		b.WriteString("\n\n")
	}

	b.WriteString("## Field schema (JSON keys to fill)\n")
	schemaJSON, _ := json.MarshalIndent(schemaForPrompt(def), "", "  ")
	b.Write(schemaJSON)
	b.WriteString("\n\n")

	b.WriteString("## Collection status\n")
	if len(missing) == 0 {
		b.WriteString("All required fields are present. Confirm briefly.\n\n")
	} else {
		fmt.Fprintf(&b, "Next focus: %s\n", missing[0])
		for _, name := range missing {
			fd := fieldByName(def, name)
			if fd.Name == "" {
				continue
			}
			fmt.Fprintf(&b, "- %s: %s", name, firstNonEmpty(fd.Why, fd.Description, fd.Ask))
			if fd.Ask != "" {
				fmt.Fprintf(&b, " Ask: %s", fd.Ask)
			}
			b.WriteByte('\n')
		}
		b.WriteByte('\n')
	}

	priorJSON, _ := json.Marshal(prior)
	missingJSON, _ := json.Marshal(missing)
	fmt.Fprintf(&b, "## existing_fields\n%s\n\n## missing_before_turn\n%s\n\n", string(priorJSON), string(missingJSON))

	b.WriteString("## documents\n")
	if len(evidence.Documents) == 0 {
		b.WriteString("(none)\n\n")
	} else {
		for _, d := range evidence.Documents {
			fmt.Fprintf(&b, "### %s (%s)\n%s\n\n",
				firstNonEmpty(d.Name, "document"),
				firstNonEmpty(d.Kind, "text"),
				truncateRunes(d.Text, 12_000))
		}
	}

	b.WriteString("## conversation\n")
	hist := formatHistory(evidence.Messages, 24, 14_000)
	if hist == "" {
		b.WriteString("(none)\n\n")
	} else {
		b.WriteString(hist)
		b.WriteString("\n")
	}

	b.WriteString("## full_user_corpus\n")
	b.WriteString(truncateRunes(BuildUserCorpus(evidence, latest), 16_000))
	b.WriteString("\n\n## latest_user_message\n")
	b.WriteString(truncateRunes(latest, 8_000))
	b.WriteString("\n\n## Output\nReturn ONLY a single JSON object (no markdown fences):\n")
	b.WriteString(`{"fields":{...},"reply":"..."}`)
	b.WriteString("\nfields keys must match the schema. Use null/omit for unknown values.\n")
	return b.String()
}

func schemaForPrompt(def ContextDef) []map[string]any {
	out := make([]map[string]any, 0, len(def.Fields))
	for _, f := range def.Fields {
		item := map[string]any{
			"name":        f.Name,
			"type":        string(f.Type),
			"required":    f.Required,
			"priority":    f.Priority,
			"description": f.Description,
		}
		if len(f.Enum) > 0 {
			item["enum"] = f.Enum
		}
		if f.MinLength > 0 {
			item["min_length"] = f.MinLength
		}
		out = append(out, item)
	}
	return out
}

func fieldByName(def ContextDef, name string) FieldDef {
	for _, f := range def.Fields {
		if f.Name == name {
			return f
		}
	}
	return FieldDef{}
}

func maxSentences(def ContextDef) int {
	if def.ReplyPolicy.MaxSentences > 0 {
		return def.ReplyPolicy.MaxSentences
	}
	return 3
}

func formatHistory(msgs []Message, maxTurns int, maxRunes int) string {
	if len(msgs) == 0 {
		return ""
	}
	start := 0
	if maxTurns > 0 && len(msgs) > maxTurns {
		start = len(msgs) - maxTurns
	}
	var b strings.Builder
	for _, m := range msgs[start:] {
		role := strings.ToLower(strings.TrimSpace(m.Role))
		if role != "user" && role != "assistant" {
			continue
		}
		if b.Len() > maxRunes {
			break
		}
		fmt.Fprintf(&b, "%s: %s\n", role, truncateRunes(m.Content, 3000))
	}
	return b.String()
}

// BuildUserCorpus concatenates user-authored evidence for extractors.
func BuildUserCorpus(ev Evidence, latest string) string {
	var b strings.Builder
	for _, d := range ev.Documents {
		t := strings.TrimSpace(d.Text)
		if t == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(truncateRunes(t, 12_000))
	}
	for _, m := range ev.Messages {
		if !strings.EqualFold(strings.TrimSpace(m.Role), "user") {
			continue
		}
		c := strings.TrimSpace(m.Content)
		if c == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(truncateRunes(c, 8_000))
		if b.Len() > 40_000 {
			break
		}
	}
	if latest = strings.TrimSpace(latest); latest != "" {
		if b.Len() > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(truncateRunes(latest, 12_000))
	}
	// Structured as key=value lines.
	if len(ev.Structured) > 0 {
		if b.Len() > 0 {
			b.WriteString("\n\n")
		}
		raw, _ := json.Marshal(ev.Structured)
		b.WriteString("structured: ")
		b.Write(raw)
	}
	return b.String()
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func truncateRunes(s string, max int) string {
	r := []rune(s)
	if max <= 0 || len(r) <= max {
		return s
	}
	return string(r[:max])
}
