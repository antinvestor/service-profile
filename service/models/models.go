package models

import (
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pgvector/pgvector-go"
	"github.com/pitabwire/frame"
	"gorm.io/datatypes"
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

	ProfileTypeID string `gorm:"type:varchar(50);index:profile_id"`
	ProfileType   ProfileType
}

type Contact struct {
	frame.BaseModel
	Detail string `gorm:"type:varchar(50);unique"`

	ContactType        string `gorm:"type:varchar(50)"`
	CommunicationLevel string `gorm:"type:varchar(50)"`

	Language string

	ProfileID string `gorm:"type:varchar(50);index:profile_id"`

	Properties datatypes.JSONMap
}

type Roster struct {
	frame.BaseModel

	ProfileID string `gorm:"type:varchar(50);uniqueIndex:roster_composite_index;index:profile_id"`

	ContactID string `gorm:"type:varchar(50);uniqueIndex:roster_composite_index;"`
	Contact   *Contact

	Properties datatypes.JSONMap
}

type Verification struct {
	frame.BaseModel
	ProfileID string `gorm:"type:varchar(50);index:profile_id"`
	ContactID string `gorm:"type:varchar(50);index:contact_id"`
	Contact   Contact

	Pin      string `gorm:"type:varchar(10)"`
	LinkHash string `gorm:"type:varchar(100)"`

	ExpiresAt  *time.Time
	VerifiedAt *time.Time
}

type VerificationAttempt struct {
	frame.BaseModel
	VerificationID string `gorm:"type:varchar(50);index:verification_id"`
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

	ParentID string `gorm:"type:varchar(50);index:parent_id"`

	CountryID string `gorm:"type:varchar(50)"`
	Country   *Country

	Properties datatypes.JSONMap
}

type ProfileAddress struct {
	frame.BaseModel
	Name string

	AddressID string `gorm:"type:varchar(50);index:address_id"`
	Address   *Address

	ProfileID string `gorm:"type:varchar(50);index:profile_id"`
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
	ParentObjectID string `gorm:"type:varchar(50);index:parent_obj_id"`

	ChildObject   string `gorm:"type:varchar(50)"`
	ChildObjectID string `gorm:"type:varchar(50);index:child_obj_id"`

	RelationshipTypeID string `gorm:"type:varchar(50);index:relationship_type_id"`
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

type Device struct {
	frame.BaseModel
	ProfileID string            `json:"profile_id" gorm:"type:varchar(50);index:profile_id"`
	LinkID    string            `json:"link_id" gorm:"type:varchar(50);index:link_id"`
	Name      string            `json:"name" gorm:"type:varchar(50)"`
	Browser   string            `json:"browser" gorm:"type:varchar(50)"`
	OS        string            `json:"os" gorm:"type:varchar(50)"`
	IP        string            `json:"ip" gorm:"type:varchar(50)"`
	Locale    datatypes.JSONMap `json:"locale"`
	Location  datatypes.JSONMap `json:"location"`
	LastSeen  time.Time         `json:"last_seen"`
	Embedding *pgvector.Vector  `json:"-" gorm:"type:vector(512)"`
}

type DeviceLog struct {
	frame.BaseModel
	DeviceID  string            `json:"device_id" gorm:"type:varchar(50)"`
	LinkID    string            `json:"link_id" gorm:"type:varchar(255)"`
	Data      datatypes.JSONMap `json:"data"`
	Embedding *pgvector.Vector  `json:"-" gorm:"type:vector(512)"`
}
