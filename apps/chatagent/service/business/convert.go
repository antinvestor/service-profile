package business

import (
	"encoding/json"
	"strings"

	"github.com/pitabwire/frame/v2/data"
	"google.golang.org/protobuf/types/known/structpb"

	chatagentv1 "github.com/antinvestor/service-profile/gen/go/chatagent/v1"

	"github.com/antinvestor/service-profile/apps/chatagent/service/engine"
	"github.com/antinvestor/service-profile/apps/chatagent/service/models"
)

func protoToContextDef(p *chatagentv1.ContextDefinition) engine.ContextDef {
	if p == nil {
		return engine.ContextDef{}
	}
	def := engine.ContextDef{
		Key:          strings.TrimSpace(p.GetContextKey()),
		Purpose:      strings.TrimSpace(p.GetPurpose()),
		SystemPrompt: strings.TrimSpace(p.GetSystemPrompt()),
		ExtractRules: strings.TrimSpace(p.GetExtractRules()),
	}
	if rp := p.GetReplyPolicy(); rp != nil {
		def.ReplyPolicy = engine.ReplyPolicy{
			MaxSentences:      int(rp.GetMaxSentences()),
			AskOneMissingOnly: rp.GetAskOneMissingOnly() || rp.GetMaxSentences() == 0,
			CompleteMessage:   rp.GetCompleteMessage(),
		}
		// Default ask-one when policy message present without explicit false —
		// proto3 bool defaults false; we treat missing policy as ask-one.
		if p.GetReplyPolicy() != nil && !rp.GetAskOneMissingOnly() && rp.GetMaxSentences() == 0 && rp.GetCompleteMessage() == "" {
			def.ReplyPolicy.AskOneMissingOnly = true
		}
	} else {
		def.ReplyPolicy = engine.ReplyPolicy{MaxSentences: 3, AskOneMissingOnly: true}
	}
	for _, f := range p.GetFields() {
		def.Fields = append(def.Fields, engine.FieldDef{
			Name:          strings.TrimSpace(f.GetName()),
			Type:          protoFieldType(f.GetType()),
			Required:      f.GetRequired(),
			Priority:      int(f.GetPriority()),
			Description:   f.GetDescription(),
			Enum:          f.GetEnumValues(),
			MinLength:     int(f.GetMinLength()),
			Ask:           f.GetAsk(),
			Why:           f.GetWhy(),
			EvidenceHints: f.GetEvidenceHints(),
		})
	}
	return def
}

func contextDefToProto(def engine.ContextDef) *chatagentv1.ContextDefinition {
	p := &chatagentv1.ContextDefinition{
		ContextKey:   def.Key,
		Purpose:      def.Purpose,
		SystemPrompt: def.SystemPrompt,
		ExtractRules: def.ExtractRules,
		ReplyPolicy: &chatagentv1.ReplyPolicy{
			MaxSentences:      int32(def.ReplyPolicy.MaxSentences),
			AskOneMissingOnly: def.ReplyPolicy.AskOneMissingOnly,
			CompleteMessage:   def.ReplyPolicy.CompleteMessage,
		},
	}
	for _, f := range def.Fields {
		p.Fields = append(p.Fields, &chatagentv1.FieldDef{
			Name:          f.Name,
			Type:          fieldTypeToProto(f.Type),
			Required:      f.Required,
			Priority:      int32(f.Priority),
			Description:   f.Description,
			EnumValues:    f.Enum,
			MinLength:     int32(f.MinLength),
			Ask:           f.Ask,
			Why:           f.Why,
			EvidenceHints: f.EvidenceHints,
		})
	}
	return p
}

func protoFieldType(t chatagentv1.FieldType) engine.FieldType {
	switch t {
	case chatagentv1.FieldType_FIELD_TYPE_NUMBER:
		return engine.FieldNumber
	case chatagentv1.FieldType_FIELD_TYPE_STRING_LIST:
		return engine.FieldStringList
	case chatagentv1.FieldType_FIELD_TYPE_BOOL:
		return engine.FieldBool
	case chatagentv1.FieldType_FIELD_TYPE_OBJECT:
		return engine.FieldObject
	default:
		return engine.FieldString
	}
}

func fieldTypeToProto(t engine.FieldType) chatagentv1.FieldType {
	switch t {
	case engine.FieldNumber:
		return chatagentv1.FieldType_FIELD_TYPE_NUMBER
	case engine.FieldStringList:
		return chatagentv1.FieldType_FIELD_TYPE_STRING_LIST
	case engine.FieldBool:
		return chatagentv1.FieldType_FIELD_TYPE_BOOL
	case engine.FieldObject:
		return chatagentv1.FieldType_FIELD_TYPE_OBJECT
	default:
		return chatagentv1.FieldType_FIELD_TYPE_STRING
	}
}

func contextDefToJSONMap(def engine.ContextDef) (data.JSONMap, error) {
	return models.JSONMapFromStruct(def)
}

func contextDefFromJSONMap(m data.JSONMap) (engine.ContextDef, error) {
	var def engine.ContextDef
	if err := models.UnmarshalJSONMap(m, &def); err != nil {
		return def, err
	}
	return def, nil
}

func structToFields(s *structpb.Struct) engine.Fields {
	if s == nil {
		return engine.Fields{}
	}
	return engine.Fields(s.AsMap())
}

func fieldsToStruct(f engine.Fields) *structpb.Struct {
	if f == nil {
		return &structpb.Struct{Fields: map[string]*structpb.Value{}}
	}
	// Sanitize for structpb (json numbers etc.).
	b, err := json.Marshal(f)
	if err != nil {
		st, _ := structpb.NewStruct(map[string]any{})
		return st
	}
	var m map[string]any
	if uerr := json.Unmarshal(b, &m); uerr != nil {
		st, _ := structpb.NewStruct(map[string]any{})
		return st
	}
	st, err := structpb.NewStruct(m)
	if err != nil {
		st, _ = structpb.NewStruct(map[string]any{})
	}
	return st
}

func protoDocs(docs []*chatagentv1.DocumentEvidence) []engine.Document {
	out := make([]engine.Document, 0, len(docs))
	for _, d := range docs {
		if d == nil {
			continue
		}
		text := strings.TrimSpace(d.GetText())
		if text == "" {
			continue
		}
		out = append(out, engine.Document{
			Name: d.GetName(),
			Kind: d.GetKind(),
			Text: text,
		})
	}
	return out
}

func protoMessages(msgs []*chatagentv1.ChatMessage) []engine.Message {
	out := make([]engine.Message, 0, len(msgs))
	for _, m := range msgs {
		if m == nil {
			continue
		}
		role := strings.ToLower(strings.TrimSpace(m.GetRole()))
		content := strings.TrimSpace(m.GetContent())
		if content == "" || (role != "user" && role != "assistant") {
			continue
		}
		out = append(out, engine.Message{Role: role, Content: content})
	}
	return out
}

func engineMessagesToProto(msgs []engine.Message) []*chatagentv1.ChatMessage {
	out := make([]*chatagentv1.ChatMessage, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, &chatagentv1.ChatMessage{Role: m.Role, Content: m.Content})
	}
	return out
}

func fieldStatusToProto(st map[string]engine.FieldStatus) map[string]*chatagentv1.FieldStatus {
	out := make(map[string]*chatagentv1.FieldStatus, len(st))
	for k, v := range st {
		out[k] = &chatagentv1.FieldStatus{Ok: v.OK, Value: v.Value, Reason: v.Reason}
	}
	return out
}

func sessionStatusProto(s string, ready bool) chatagentv1.SessionStatus {
	switch s {
	case models.SessionStatusEnded:
		return chatagentv1.SessionStatus_SESSION_STATUS_ENDED
	case models.SessionStatusReady:
		return chatagentv1.SessionStatus_SESSION_STATUS_READY
	default:
		if ready {
			return chatagentv1.SessionStatus_SESSION_STATUS_READY
		}
		return chatagentv1.SessionStatus_SESSION_STATUS_ACTIVE
	}
}

func docsToJSONMap(docs []engine.Document) data.JSONMap {
	items := make([]any, 0, len(docs))
	for _, d := range docs {
		items = append(items, map[string]any{
			"name": d.Name,
			"kind": d.Kind,
			"text": d.Text,
		})
	}
	m, _ := models.JSONMapFromStruct(map[string]any{"items": items})
	return m
}

func docsFromJSONMap(m data.JSONMap) []engine.Document {
	if m == nil {
		return nil
	}
	raw, ok := m["items"]
	if !ok {
		return nil
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	var items []engine.Document
	if uerr := json.Unmarshal(b, &items); uerr != nil {
		return nil
	}
	return items
}

func fieldsFromJSONMap(m data.JSONMap) engine.Fields {
	return engine.Fields(models.MapFromJSONMap(m))
}

func fieldsToJSONMap(f engine.Fields) data.JSONMap {
	m, err := models.JSONMapFromStruct(f)
	if err != nil {
		return data.JSONMap{}
	}
	return m
}

func stringField(m data.JSONMap, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return strings.TrimSpace(s)
}

func boolField(m data.JSONMap, key string) bool {
	v, ok := m[key]
	if !ok || v == nil {
		return false
	}
	b, _ := v.(bool)
	return b
}
