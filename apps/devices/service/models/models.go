package models

import (
	"time"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/pitabwire/frame"
	"google.golang.org/protobuf/encoding/protojson"
)

// Device represents a core device identity.
type Device struct {
	frame.BaseModel
	ProfileID string `gorm:"index;size:40"`
	Name      string `gorm:"size:255"`
	OS        string `gorm:"size:255"`
}

func (d *Device) ToAPI(session *DeviceSession) *devicev1.DeviceObject {
	obj := &devicev1.DeviceObject{
		Id:   d.GetID(),
		Name: d.Name,
		Os:   d.OS,
	}

	if session != nil {
		obj.SessionId = session.GetID()
		obj.UserAgent = session.UserAgent
		obj.Ip = session.IP

		var locale devicev1.Locale
		_ = protojson.Unmarshal(session.Locale, &locale)

		obj.Locale = &locale
		obj.Location = frame.DBPropertiesToMap(session.Location)
		obj.LastSeen = session.LastSeen.String()
	}

	return obj
}

// DeviceSession represents a single session of a device.
type DeviceSession struct {
	frame.BaseModel
	DeviceID  string `gorm:"index"`
	UserAgent string `gorm:"size:512"`
	IP        string `gorm:"size:45"`
	Locale    []byte `gorm:"type:bytea"`
	Location  frame.JSONMap
	LastSeen  time.Time
}

// DeviceKey holds encryption keys for a device.
type DeviceKey struct {
	frame.BaseModel
	DeviceID string `gorm:"index"`
	Key      []byte `gorm:"type:bytea"`
	Extra    frame.JSONMap
}

func (k *DeviceKey) ToAPI() *devicev1.KeyObject {
	return &devicev1.KeyObject{
		Id:       k.GetID(),
		DeviceId: k.DeviceID,
		Key:      k.Key,
		Extra:    frame.DBPropertiesToMap(k.Extra),
	}
}

// DeviceLog records activities for a device.
type DeviceLog struct {
	frame.BaseModel
	DeviceID        string `gorm:"index"`
	DeviceSessionID string `gorm:"index"`
	Data            frame.JSONMap
}

func (dl *DeviceLog) ToAPI() *devicev1.DeviceLog {
	extra := frame.DBPropertiesToMap(dl.Data)

	return &devicev1.DeviceLog{
		Id:        dl.GetID(),
		DeviceId:  dl.DeviceID,
		SessionId: dl.DeviceSessionID,
		Extra:     extra,
	}
}
