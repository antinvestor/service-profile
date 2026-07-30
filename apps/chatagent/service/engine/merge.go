package engine

import (
	"strings"
)

// MergeFields overlays non-empty values from overlay onto base without wiping
// existing evidence. Longer string values win for free-text fields.
func MergeFields(base, overlay Fields) Fields {
	out := CloneFields(base)
	if overlay == nil {
		return out
	}
	for k, v := range overlay {
		if isEmptyValue(v) {
			continue
		}
		// Prefer longer string for document-like fields.
		if existing, ok := out[k]; ok {
			es := strings.TrimSpace(asString(existing))
			ns := strings.TrimSpace(asString(v))
			if es != "" && ns != "" && len(ns) < len(es) {
				// Keep existing longer text unless types differ meaningfully.
				if _, isList := v.([]any); !isList {
					if _, isList2 := v.([]string); !isList2 {
						continue
					}
				}
			}
		}
		out[k] = v
	}
	return out
}

// ApplyStructured copies structured inputs onto fields using field names as keys.
func ApplyStructured(fields Fields, structured Fields) Fields {
	return MergeFields(fields, structured)
}

// ApplyDocuments folds document text into fields when a field name matches
// document name/kind, or into a field that lists "document" evidence hints
// and is still empty (first such required field by priority).
func ApplyDocuments(def ContextDef, fields Fields, docs []Document) Fields {
	out := CloneFields(fields)
	for _, doc := range docs {
		text := strings.TrimSpace(doc.Text)
		if text == "" {
			continue
		}
		// Direct name match: document name == field name.
		if doc.Name != "" {
			if existing, ok := out[doc.Name]; !ok || isEmptyValue(existing) || len(asString(existing)) < len(text) {
				out[doc.Name] = text
			}
		}
		// Kind match.
		if doc.Kind != "" && doc.Kind != doc.Name {
			if existing, ok := out[doc.Kind]; !ok || isEmptyValue(existing) || len(asString(existing)) < len(text) {
				out[doc.Kind] = text
			}
		}
	}
	// If a required string field has evidence_hints containing "document" and is empty,
	// attach the longest document text once.
	longest := longestDocText(docs)
	if longest == "" {
		return out
	}
	for _, f := range orderedFields(def) {
		if !f.Required {
			continue
		}
		if f.Type != FieldString && f.Type != "" {
			continue
		}
		if !hintsContain(f.EvidenceHints, "document") {
			continue
		}
		if existing, ok := out[f.Name]; ok && !isEmptyValue(existing) {
			// Prefer longer.
			if len(asString(existing)) >= len(longest) {
				continue
			}
		}
		out[f.Name] = longest
		break
	}
	return out
}

func longestDocText(docs []Document) string {
	var best string
	for _, d := range docs {
		t := strings.TrimSpace(d.Text)
		if len(t) > len(best) {
			best = t
		}
	}
	return best
}

func orderedFields(def ContextDef) []FieldDef {
	out := append([]FieldDef(nil), def.Fields...)
	// sort by priority ascending
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].Priority < out[i].Priority ||
				(out[j].Priority == out[i].Priority && out[j].Name < out[i].Name) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}

func hintsContain(hints []string, want string) bool {
	for _, h := range hints {
		if strings.EqualFold(strings.TrimSpace(h), want) {
			return true
		}
	}
	return false
}

func isEmptyValue(v any) bool {
	if v == nil {
		return true
	}
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t) == ""
	case []any:
		return len(t) == 0
	case []string:
		return len(t) == 0
	case map[string]any:
		return len(t) == 0
	case Fields:
		return len(t) == 0
	default:
		return strings.TrimSpace(asString(v)) == ""
	}
}

// Sanitize drops empty values and trims strings; enforces enums when set.
func Sanitize(def ContextDef, fields Fields) Fields {
	out := Fields{}
	byName := make(map[string]FieldDef, len(def.Fields))
	for _, f := range def.Fields {
		byName[f.Name] = f
	}
	for k, v := range fields {
		if isEmptyValue(v) {
			continue
		}
		fd, hasDef := byName[k]
		if !hasDef {
			// Keep unknown keys from product seed (opaque bag) if non-empty.
			out[k] = v
			continue
		}
		switch fd.Type {
		case FieldString, "":
			s := strings.TrimSpace(asString(v))
			if s == "" {
				continue
			}
			if len(fd.Enum) > 0 {
				matched := false
				for _, e := range fd.Enum {
					if strings.EqualFold(e, s) {
						out[k] = e
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			} else {
				out[k] = s
			}
		case FieldStringList:
			list := asStringList(v)
			if len(list) == 0 {
				continue
			}
			if len(fd.Enum) > 0 {
				var filtered []string
				for _, item := range list {
					for _, e := range fd.Enum {
						if strings.EqualFold(e, item) {
							filtered = append(filtered, e)
							break
						}
					}
				}
				if len(filtered) == 0 {
					continue
				}
				out[k] = filtered
			} else {
				out[k] = list
			}
		case FieldNumber:
			if n, ok := asFloat(v); ok {
				out[k] = n
			}
		case FieldBool:
			if b, ok := asBool(v); ok {
				out[k] = b
			}
		default:
			out[k] = v
		}
	}
	return out
}
