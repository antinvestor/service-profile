package models

import (
	settingsV1 "github.com/antinvestor/apis/go/settings/v1"
	"github.com/pitabwire/frame/data"
)

// SettingRef Table holds the templete details.
type SettingRef struct {
	data.BaseModel

	Name     string `gorm:"type:varchar(255)"`
	Object   string `gorm:"type:varchar(255)"`
	ObjectID string `gorm:"type:varchar(255)"`
	Language string `gorm:"type:varchar(255)"`
	Module   string `gorm:"type:varchar(255)"`
}

func (model *SettingRef) ToAPI() *settingsV1.Setting {
	setting := settingsV1.Setting{
		Name:     model.Name,
		Object:   model.Object,
		ObjectId: model.ObjectID,
		Lang:     model.Language,
		Module:   model.Module,
	}
	return &setting
}

type SettingVal struct {
	data.BaseModel
	Ref     string `gorm:"type:varchar(50);unique_index"`
	Detail  string `gorm:"type:text"`
	Version int
}

func (model *SettingVal) ToAPI(sRef *SettingRef) *settingsV1.SettingObject {
	response := settingsV1.SettingObject{
		Id:      model.ID,
		Key:     sRef.ToAPI(),
		Value:   model.Detail,
		Updated: model.ModifiedAt.String(),
	}
	return &response
}

// SettingAudit table holds a history of all the setting values overtime.
type SettingAudit struct {
	data.BaseModel

	Ref     string `gorm:"type:varchar(50);unique_index"`
	Detail  string `gorm:"type:text"`
	Version int
}
