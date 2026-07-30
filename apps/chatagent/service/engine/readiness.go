package engine

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Assess evaluates required fields against current values. Ready is true only
// when every required field passes its type/min-length/enum checks.
func Assess(def ContextDef, fields Fields) (map[string]FieldStatus, []string, bool) {
	status := make(map[string]FieldStatus, len(def.Fields))
	var required []FieldDef
	for _, f := range def.Fields {
		st := assessOne(f, fields)
		status[f.Name] = st
		if f.Required {
			required = append(required, f)
		}
	}
	sort.SliceStable(required, func(i, j int) bool {
		if required[i].Priority == required[j].Priority {
			return required[i].Name < required[j].Name
		}
		return required[i].Priority < required[j].Priority
	})
	var missing []string
	for _, f := range required {
		if !status[f.Name].OK {
			missing = append(missing, f.Name)
		}
	}
	return status, missing, len(missing) == 0
}

func assessOne(def FieldDef, fields Fields) FieldStatus {
	v, ok := fields[def.Name]
	if !ok || v == nil {
		return FieldStatus{OK: false, Reason: missingReason(def)}
	}
	switch def.Type {
	case FieldString, "":
		s := strings.TrimSpace(asString(v))
		if s == "" {
			return FieldStatus{OK: false, Reason: missingReason(def)}
		}
		if def.MinLength > 0 && len([]rune(s)) < def.MinLength {
			return FieldStatus{OK: false, Reason: fmt.Sprintf("need at least %d characters", def.MinLength)}
		}
		if len(def.Enum) > 0 && !enumContains(def.Enum, s) {
			return FieldStatus{OK: false, Reason: "value not in allowed set"}
		}
		return FieldStatus{OK: true, Value: truncate(s, 200)}
	case FieldNumber:
		n, okNum := asFloat(v)
		if !okNum {
			return FieldStatus{OK: false, Reason: missingReason(def)}
		}
		return FieldStatus{OK: true, Value: fmt.Sprintf("%g", n)}
	case FieldStringList:
		list := asStringList(v)
		if len(list) == 0 {
			return FieldStatus{OK: false, Reason: missingReason(def)}
		}
		if len(def.Enum) > 0 {
			filtered := make([]string, 0, len(list))
			for _, item := range list {
				if enumContains(def.Enum, item) {
					filtered = append(filtered, item)
				}
			}
			if len(filtered) == 0 {
				return FieldStatus{OK: false, Reason: "no allowed list values"}
			}
			list = filtered
		}
		return FieldStatus{OK: true, Value: strings.Join(list, ", ")}
	case FieldBool:
		b, okBool := asBool(v)
		if !okBool {
			return FieldStatus{OK: false, Reason: missingReason(def)}
		}
		return FieldStatus{OK: true, Value: fmt.Sprintf("%v", b)}
	case FieldObject:
		if m, okMap := v.(map[string]any); okMap && len(m) > 0 {
			return FieldStatus{OK: true, Value: "(object)"}
		}
		return FieldStatus{OK: false, Reason: missingReason(def)}
	default:
		s := strings.TrimSpace(asString(v))
		if s == "" {
			return FieldStatus{OK: false, Reason: missingReason(def)}
		}
		return FieldStatus{OK: true, Value: truncate(s, 200)}
	}
}

func missingReason(def FieldDef) string {
	if def.Ask != "" {
		return def.Ask
	}
	if def.Description != "" {
		return def.Description
	}
	return "required"
}

func enumContains(enum []string, v string) bool {
	v = strings.TrimSpace(strings.ToLower(v))
	for _, e := range enum {
		if strings.EqualFold(strings.TrimSpace(e), v) {
			return true
		}
	}
	return false
}

func asString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case json.Number:
		return t.String()
	case fmt.Stringer:
		return t.String()
	case float64:
		return fmt.Sprintf("%g", t)
	case float32:
		return fmt.Sprintf("%g", t)
	case int:
		return fmt.Sprintf("%d", t)
	case int64:
		return fmt.Sprintf("%d", t)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprint(v)
		}
		s := strings.TrimSpace(string(b))
		return strings.Trim(s, `"`)
	}
}

func asFloat(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case json.Number:
		f, err := t.Float64()
		return f, err == nil
	case string:
		var f float64
		_, err := fmt.Sscanf(strings.TrimSpace(t), "%f", &f)
		return f, err == nil
	default:
		return 0, false
	}
}

func asBool(v any) (bool, bool) {
	switch t := v.(type) {
	case bool:
		return t, true
	case string:
		switch strings.ToLower(strings.TrimSpace(t)) {
		case "true", "yes", "1":
			return true, true
		case "false", "no", "0":
			return false, true
		}
	}
	return false, false
}

func asStringList(v any) []string {
	switch t := v.(type) {
	case []string:
		return nonEmptyStrings(t)
	case []any:
		out := make([]string, 0, len(t))
		for _, item := range t {
			if s := strings.TrimSpace(asString(item)); s != "" {
				out = append(out, s)
			}
		}
		return out
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return nil
		}
		parts := strings.Split(s, ",")
		return nonEmptyStrings(parts)
	default:
		return nil
	}
}

func nonEmptyStrings(in []string) []string {
	var out []string
	for _, s := range in {
		if t := strings.TrimSpace(s); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func truncate(s string, maxLen int) string {
	r := []rune(s)
	if len(r) <= maxLen {
		return s
	}
	return string(r[:maxLen]) + "…"
}
