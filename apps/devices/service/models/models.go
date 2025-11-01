package models

import (
	"time"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/pitabwire/frame/data"
	"google.golang.org/protobuf/encoding/protojson"
)

// Device represents a core device identity.
type Device struct {
	data.BaseModel
	ProfileID string `gorm:"index;size:40"`
	Name      string `gorm:"size:255"`
	OS        string `gorm:"size:255"`
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
	DeviceID  string `gorm:"index"`
	UserAgent string `gorm:"size:512"`
	IP        string `gorm:"size:45"`
	Locale    []byte `gorm:"type:bytea"`
	Location  data.JSONMap
	LastSeen  time.Time
}

// DeviceKey holds encryption keys for a device.
type DeviceKey struct {
	data.BaseModel
	DeviceID string `gorm:"index"`
	Key      []byte `gorm:"type:bytea"`
	Extra    data.JSONMap
}

func (k *DeviceKey) ToAPI() *devicev1.KeyObject {
	return &devicev1.KeyObject{
		Id:       k.GetID(),
		DeviceId: k.DeviceID,
		Key:      k.Key,
		Extra:    k.Extra.ToProtoStruct(),
	}
}

// DeviceLog records activities for a device.
type DeviceLog struct {
	data.BaseModel
	DeviceID        string `gorm:"index"`
	DeviceSessionID string `gorm:"index"`
	Data            data.JSONMap
}

func (dl *DeviceLog) ToAPI() *devicev1.DeviceLog {
	return &devicev1.DeviceLog{
		Id:        dl.GetID(),
		DeviceId:  dl.DeviceID,
		SessionId: dl.DeviceSessionID,
		Extra:     dl.Data.ToProtoStruct(),
	}
}
