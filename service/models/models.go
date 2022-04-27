package models

import (
	"encoding/json"
	papi "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/service-profile/service"
	"github.com/pitabwire/frame"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

var ProfileTypeIDMap = map[papi.ProfileType]uint{
	papi.ProfileType_PERSON:      0,
	papi.ProfileType_INSTITUTION: 1,
	papi.ProfileType_BOT:         2,
}

func ProfileTypeIDToEnum(profileTypeID uint) papi.ProfileType {
	for key, val := range ProfileTypeIDMap {
		if val == profileTypeID {
			return key
		}
	}
	return papi.ProfileType_PERSON
}

type ProfileType struct {
	frame.BaseModel
	UID         uint `sql:"unique"`
	Name        string
	Description string
}

type Profile struct {
	frame.BaseModel
	Properties datatypes.JSONMap

	ProfileTypeID string `gorm:"type:varchar(50)"`
	ProfileType   *ProfileType
}

func (p *Profile) GetByID(db *gorm.DB) error {
	modelID := strings.TrimSpace(p.ID)
	if err := db.Preload(clause.Associations).Where("id = ?", modelID).First(p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return service.ErrorProfileDoesNotExist
		}
		return err
	}
	return nil
}

func (p *Profile) Create(db *gorm.DB, profileType papi.ProfileType, properties map[string]interface{}) error {

	stringProperties, err := json.Marshal(properties)
	if err != nil {
		return err
	}

	err = p.Properties.UnmarshalJSON(stringProperties)
	if err != nil {
		return err
	}

	return db.Save(p).Error
}

var ContactTypeUIDMap = map[papi.ContactType]uint{
	papi.ContactType_EMAIL: 0,
	papi.ContactType_PHONE: 1,
}

var CommunicationLevelUIDMap = map[papi.CommunicationLevel]uint{
	papi.CommunicationLevel_ALL:           0,
	papi.CommunicationLevel_SYSTEM_ALERTS: 1,
	papi.CommunicationLevel_NO_CONTACT:    2,
}

type ContactType struct {
	frame.BaseModel
	UID uint `sql:"unique"`

	Name        string
	Description string
}

func ContactTypeIDToEnum(contactTypeID uint) papi.ContactType {
	for key, val := range ContactTypeUIDMap {
		if val == contactTypeID {
			return key
		}
	}
	return papi.ContactType_EMAIL
}

type CommunicationLevel struct {
	frame.BaseModel
	UID uint `sql:"unique"`

	Name        string
	Description string
}

func (cl *CommunicationLevel) From(db *gorm.DB, communicationLevel papi.CommunicationLevel) {
	cl.UID = CommunicationLevelUIDMap[communicationLevel]
	db.First(cl)
}

func CommunicationLevelIDToEnum(communicationLevelID uint) papi.CommunicationLevel {
	for key, val := range CommunicationLevelUIDMap {
		if val == communicationLevelID {
			return key
		}
	}
	return papi.CommunicationLevel_ALL
}

type Contact struct {
	frame.BaseModel
	Detail []byte `gorm:"type:bytea;unique"`
	Nonce  []byte `gorm:"type:bytea"`
	Tokens string `gorm:"type:tsvector"`

	ContactTypeID string `gorm:"type:varchar(50)"`
	ContactType   *ContactType

	CommunicationLevelID string `gorm:"type:varchar(50)"`
	CommunicationLevel   *CommunicationLevel

	Language string

	ProfileID string `gorm:"type:varchar(50)"`
	Profile   Profile
}

func GetContactsByProfile(db *gorm.DB, p *Profile) ([]Contact, error) {

	var profileContacts []Contact

	err := db.Preload(clause.Associations).Where("profile_id = ?", p.ID).Find(&profileContacts).Error

	if err != nil {
		return nil, err
	}

	return profileContacts, nil
}

type Verification struct {
	frame.BaseModel
	ContactID string `gorm:"type:varchar(50)"`
	Contact   Contact

	Pin      string `gorm:"type:varchar(10)"`
	LinkHash string `gorm:"type:varchar(100)"`

	ExpiresAt  *time.Time
	VerifiedAt *time.Time
}

type VerificationAttempt struct {
	frame.BaseModel
	VerificationID string `gorm:"type:varchar(50)"`
	Verification   Verification

	ContactID string `gorm:"type:varchar(50)"`
	Contact   Contact

	State string `gorm:"type:varchar(10)"`
}

type Country struct {
	frame.BaseModel
	ISO3         string `gorm:"unique"`
	ISO2         string `sql:"unique"`
	Name         string
	Code         int
	LatitudeAvg  float64
	LongitudeAvg float64
	City         string
}

type Address struct {
	frame.BaseModel
	Name      string
	AdminUnit string

	ParentID string `gorm:"type:varchar(50)"`
	Parent   *Address

	CountryID string `gorm:"type:varchar(50)"`
	Country   *Country

	Properties datatypes.JSONMap
}

type ProfileAddress struct {
	frame.BaseModel
	Name string

	AddressID string `gorm:"type:varchar(50)"`
	Address   *Address

	ProfileID string `gorm:"type:varchar(50)"`
	Profile   *Profile
}
