package models

import (
	"time"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"github.com/pitabwire/frame/data"
	"google.golang.org/protobuf/encoding/protojson"
)

// Device represents a core device identity.
type Device struct {
	data.BaseModel
	ProfileID string `gorm:"index;size:40" json:"profile_id"`
	Name      string `gorm:"size:255"      json:"name"`
	OS        string `gorm:"size:255"      json:"os"`
}

func (d *Device) ToAPI(session *DeviceSession) *devicev1.DeviceObject {
	ownerProperties := data.JSONMap{"owner": d.ProfileID}

	obj := &devicev1.DeviceObject{
		Id:         d.GetID(),
		Name:       d.Name,
		Os:         d.OS,
		Properties: ownerProperties.ToProtoStruct(),
	}

	if session != nil {
		obj.SessionId = session.GetID()
		obj.UserAgent = session.UserAgent
		obj.Ip = session.IP

		if len(session.Locale) > 0 {
			var locale devicev1.Locale
			if err := protojson.Unmarshal(session.Locale, &locale); err == nil {
				obj.Locale = &locale
			}
		}
		obj.Location = session.Location.ToProtoStruct()
		obj.LastSeen = session.LastSeen.String()
	}

	return obj
}

// DeviceSession represents a single session of a device.
type DeviceSession struct {
	data.BaseModel
	DeviceID  string       `gorm:"index"      json:"device_id"`
	UserAgent string       `gorm:"size:512"   json:"user_agent"`
	IP        string       `gorm:"size:45"    json:"ip"`
	Locale    []byte       `gorm:"type:bytea" json:"locale"`
	Location  data.JSONMap `                  json:"location"`
	LastSeen  time.Time    `                  json:"last_seen"`
}

// DeviceKey holds encryption keys for a device.
type DeviceKey struct {
	data.BaseModel
	DeviceID  string           `gorm:"index"      json:"device_id"`
	KeyType   devicev1.KeyType `                  json:"key_type"`
	Key       []byte           `gorm:"type:bytea" json:"key"`
	ExpiresAt *time.Time       `                  json:"expires_at,omitempty"`
	Extra     data.JSONMap     `                  json:"extra"`
}

func (k *DeviceKey) ToAPI() *devicev1.KeyObject {
	expiryTime := ""
	isActive := true
	if k.ExpiresAt != nil && !k.ExpiresAt.IsZero() {
		expiryTime = k.ExpiresAt.UTC().Format(time.RFC3339)
		isActive = k.ExpiresAt.After(time.Now())
	}
	return &devicev1.KeyObject{
		Id:        k.GetID(),
		DeviceId:  k.DeviceID,
		KeyType:   k.KeyType,
		Key:       k.Key,
		CreatedAt: k.CreatedAt.UTC().Format(time.RFC3339),
		ExpiresAt: expiryTime,
		IsActive:  isActive,
		Extra:     k.Extra.ToProtoStruct(),
	}
}

// DeviceLog records activities for a device.
type DeviceLog struct {
	data.BaseModel
	DeviceID        string       `gorm:"index"`
	DeviceSessionID string       `gorm:"index"`
	Data            data.JSONMap `gorm:"type:jsonb"`
}

func (dl *DeviceLog) ToAPI() *devicev1.DeviceLog {
	return &devicev1.DeviceLog{
		Id:        dl.GetID(),
		DeviceId:  dl.DeviceID,
		SessionId: dl.DeviceSessionID,
		Extra:     dl.Data.ToProtoStruct(),
	}
}

// DevicePresence records availablity of a device.
type DevicePresence struct {
	data.BaseModel
	DeviceID      string                  `gorm:"uniqueIndex"`
	ProfileID     string                  `gorm:"index"`
	Status        devicev1.PresenceStatus `gorm:"type:int"`
	StatusMessage string
	ExpiryTime    *time.Time
	Data          data.JSONMap
}

func (p *DevicePresence) ToAPI() *devicev1.PresenceObject {
	presence := &devicev1.PresenceObject{
		DeviceId:      p.DeviceID,
		ProfileId:     p.ProfileID,
		Status:        p.Status,
		StatusMessage: p.StatusMessage,
		Extras:        p.Data.ToProtoStruct(),
	}

	presence.LastActive = p.CreatedAt.UTC().Format(time.RFC3339)

	return presence
}
