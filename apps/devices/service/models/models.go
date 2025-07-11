package models

import (
	"time"

	"github.com/pitabwire/frame"
)

type Device struct {
	frame.BaseModel
	ProfileID string        `json:"profile_id" gorm:"type:varchar(50);index:profile_id"`
	LinkID    string        `json:"link_id"    gorm:"type:varchar(50);index:link_id"`
	Name      string        `json:"name"       gorm:"type:varchar(50)"`
	Browser   string        `json:"browser"    gorm:"type:varchar(50)"`
	OS        string        `json:"os"         gorm:"type:varchar(50)"`
	IP        string        `json:"ip"         gorm:"type:varchar(50)"`
	Locale    frame.JSONMap `json:"locale"`
	Location  frame.JSONMap `json:"location"`
	LastSeen  time.Time     `json:"last_seen"`
}

type DeviceLog struct {
	frame.BaseModel
	DeviceID string        `json:"device_id" gorm:"type:varchar(50)"`
	LinkID   string        `json:"link_id"   gorm:"type:varchar(255)"`
	Data     frame.JSONMap `json:"data"`
}
