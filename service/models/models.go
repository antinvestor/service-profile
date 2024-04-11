package models

import (
	"encoding/json"
	"errors"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service"
	"github.com/pitabwire/frame"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

var ProfileTypeIDMap = map[profilev1.ProfileType]uint{
	profilev1.ProfileType_PERSON:      0,
	profilev1.ProfileType_INSTITUTION: 1,
	profilev1.ProfileType_BOT:         2,
}

var RelationshipTypeIDMap = map[profilev1.RelationshipType]uint{
	profilev1.RelationshipType_MEMBER:       0,
	profilev1.RelationshipType_AFFILIATED:   1,
	profilev1.RelationshipType_BLACK_LISTED: 1,
}

func ProfileTypeIDToEnum(profileTypeID uint) profilev1.ProfileType {
	for key, val := range ProfileTypeIDMap {
		if val == profileTypeID {
			return key
		}
	}
	return profilev1.ProfileType_PERSON
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
	modelID := strings.TrimSpace(p.GetID())
	if err := db.Preload(clause.Associations).Where("id = ?", modelID).First(p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrorProfileDoesNotExist
		}
		return err
	}
	return nil
}

func (p *Profile) Create(db *gorm.DB, profileType profilev1.ProfileType, properties map[string]interface{}) error {

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

var ContactTypeUIDMap = map[profilev1.ContactType]uint{
	profilev1.ContactType_EMAIL: 0,
	profilev1.ContactType_PHONE: 1,
}

var CommunicationLevelUIDMap = map[profilev1.CommunicationLevel]uint{
	profilev1.CommunicationLevel_ALL:           0,
	profilev1.CommunicationLevel_SYSTEM_ALERTS: 1,
	profilev1.CommunicationLevel_NO_CONTACT:    2,
}

type ContactType struct {
	frame.BaseModel
	UID uint `sql:"unique"`

	Name        string
	Description string
}

func ContactTypeIDToEnum(contactTypeID uint) profilev1.ContactType {
	for key, val := range ContactTypeUIDMap {
		if val == contactTypeID {
			return key
		}
	}
	return profilev1.ContactType_EMAIL
}

type CommunicationLevel struct {
	frame.BaseModel
	UID uint `sql:"unique"`

	Name        string
	Description string
}

func (cl *CommunicationLevel) From(db *gorm.DB, communicationLevel profilev1.CommunicationLevel) {
	cl.UID = CommunicationLevelUIDMap[communicationLevel]
	db.First(cl)
}

func CommunicationLevelIDToEnum(communicationLevelID uint) profilev1.CommunicationLevel {
	for key, val := range CommunicationLevelUIDMap {
		if val == communicationLevelID {
			return key
		}
	}
	return profilev1.CommunicationLevel_ALL
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

	err := db.Preload(clause.Associations).Where("profile_id = ?", p.GetID()).Find(&profileContacts).Error

	if err != nil {
		return nil, err
	}

	return profileContacts, nil
}

type Verification struct {
	frame.BaseModel
	ProfileID string `gorm:"type:varchar(50)"`
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

type RelationshipType struct {
	frame.BaseModel
	UID         uint `sql:"unique"`
	Name        string
	Description string
}

func RelationshipTypeIDToEnum(relationshipTypeID uint) profilev1.RelationshipType {
	for key, val := range RelationshipTypeIDMap {
		if val == relationshipTypeID {
			return key
		}
	}
	return profilev1.RelationshipType_MEMBER
}

type Relationship struct {
	frame.BaseModel

	ParentObject   string `gorm:"type:varchar(50)"`
	ParentObjectID string `gorm:"type:varchar(50)"`

	ChildObject   string `gorm:"type:varchar(50)"`
	ChildObjectID string `gorm:"type:varchar(50)"`

	RelationshipTypeID string `gorm:"type:varchar(50)"`
	RelationshipType   *RelationshipType

	Properties datatypes.JSONMap
}

func (r *Relationship) ToAPI() *profilev1.RelationshipObject {

	relationshipObj := &profilev1.RelationshipObject{
		Id:         r.GetID(),
		Type:       profilev1.RelationshipType(r.RelationshipType.UID),
		Properties: frame.DBPropertiesToMap(r.Properties),
		ChildEntry: &profilev1.EntryItem{
			ObjectName: r.ChildObject,
			ObjectId:   r.ChildObjectID,
		},
		ParentEntry: &profilev1.EntryItem{
			ObjectName: r.ParentObject,
			ObjectId:   r.ParentObjectID,
		},
	}

	return relationshipObj

}
