package models

import (
	"encoding/json"

	"github.com/pitabwire/frame/v2/data"
)

const (
	SessionStatusActive = "active"
	SessionStatusReady  = "ready"
	SessionStatusEnded  = "ended"
)

// IntakeContext is a versioned product context definition (prompt + fields).
type IntakeContext struct {
	data.BaseModel

	ContextKey string `gorm:"type:varchar(128);not null;index:idx_chat_ctx_key,priority:1"`
	Version    int    `gorm:"not null;index:idx_chat_ctx_key,priority:2"`
	Purpose    string `gorm:"type:text"`
	// DefinitionJSON is the full ContextDef snapshot.
	DefinitionJSON data.JSONMap `gorm:"type:jsonb"`
	Active         bool         `gorm:"not null;default:true;index"`
}

func (IntakeContext) TableName() string { return "chat_contexts" }

// Session is one subject's intake conversation under a context.
type Session struct {
	data.BaseModel

	SubjectID      string       `gorm:"type:varchar(64);not null;index:idx_chat_sess_subject,priority:1"`
	ContextKey     string       `gorm:"type:varchar(128);not null"`
	ContextVersion int          `gorm:"not null;default:1"`
	ConfigSnapshot data.JSONMap `gorm:"type:jsonb"`
	Fields         data.JSONMap `gorm:"type:jsonb"`
	Runtime        data.JSONMap `gorm:"type:jsonb"`
	Documents      data.JSONMap `gorm:"type:jsonb"` // {items:[{name,kind,text}]}
	// Channel is the omnichannel binding snapshot (Notification delivery target).
	// Empty / web means replies are returned on the RPC only.
	Channel      data.JSONMap `gorm:"type:jsonb"`
	ChannelName  string       `gorm:"type:varchar(32);index:idx_chat_sess_channel,priority:1"` // denormalized for lookup
	ContactID    string       `gorm:"type:varchar(64);index:idx_chat_sess_channel,priority:2"`
	Status       string       `gorm:"type:varchar(32);not null;default:'active';index"`
	Ready        bool         `gorm:"not null;default:false;index"`
	MessageCount int          `gorm:"not null;default:0"`
}

func (Session) TableName() string { return "chat_sessions" }

// Message is one transcript row.
type Message struct {
	data.BaseModel

	SessionID string `gorm:"type:varchar(50);not null;index:idx_chat_msg_sess,priority:1"`
	Seq       int    `gorm:"not null;index:idx_chat_msg_sess,priority:2"`
	Role      string `gorm:"type:varchar(16);not null"`
	Content   string `gorm:"type:text;not null"`
}

func (Message) TableName() string { return "chat_messages" }

// DefinitionBytes marshals a map for storage convenience.
func DefinitionBytes(m map[string]any) data.JSONMap {
	if m == nil {
		return data.JSONMap{}
	}
	return data.JSONMap(m)
}

// MapFromJSONMap copies JSONMap to map[string]any.
func MapFromJSONMap(m data.JSONMap) map[string]any {
	if m == nil {
		return map[string]any{}
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// JSONMapFromStruct marshals any value into JSONMap via JSON.
func JSONMapFromStruct(v any) (data.JSONMap, error) {
	if v == nil {
		return data.JSONMap{}, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m data.JSONMap
	if uerr := json.Unmarshal(b, &m); uerr != nil {
		return nil, uerr
	}
	if m == nil {
		return data.JSONMap{}, nil
	}
	return m, nil
}

// UnmarshalJSONMap decodes JSONMap into dest.
func UnmarshalJSONMap(m data.JSONMap, dest any) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}
