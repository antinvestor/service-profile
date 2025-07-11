package models

import (
	"math"
	"time"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame"
)

// Profile type and relationship type constants.
const (
	// Profile type IDs.
	ProfileTypePersonID      uint = 1
	ProfileTypeBotID         uint = 2
	ProfileTypeInstitutionID uint = 3

	// Relationship type IDs.
	RelationshipTypeMemberID      uint = 1
	RelationshipTypeAffiliatedID  uint = 2
	RelationshipTypeBlackListedID uint = 3
)

// ProfileTypeIDMap maps profile types to their respective IDs.
var ProfileTypeIDMap = map[profilev1.ProfileType]uint{
	profilev1.ProfileType_PERSON:      ProfileTypePersonID,
	profilev1.ProfileType_BOT:         ProfileTypeBotID,
	profilev1.ProfileType_INSTITUTION: ProfileTypeInstitutionID,
}

// RelationshipTypeIDMap maps relationship types to their respective IDs.
var RelationshipTypeIDMap = map[profilev1.RelationshipType]uint{
	profilev1.RelationshipType_MEMBER:       RelationshipTypeMemberID,
	profilev1.RelationshipType_AFFILIATED:   RelationshipTypeAffiliatedID,
	profilev1.RelationshipType_BLACK_LISTED: RelationshipTypeBlackListedID,
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
	Properties frame.JSONMap

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

	Properties frame.JSONMap
}

type Roster struct {
	frame.BaseModel

	ProfileID string `gorm:"type:varchar(50);uniqueIndex:roster_composite_index;index:profile_id"`

	ContactID string `gorm:"type:varchar(50);uniqueIndex:roster_composite_index;"`
	Contact   *Contact

	Properties frame.JSONMap
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
	ISO2         string `              sql:"unique"`
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

	Properties frame.JSONMap
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

	Properties frame.JSONMap
}

func (r *Relationship) ToAPI() *profilev1.RelationshipObject {
	// Safe conversion from uint to int32
	var relationshipTypeValue int32
	if r.RelationshipType.UID <= uint(math.MaxInt32) {
		relationshipTypeValue = int32(r.RelationshipType.UID)
	} else {
		relationshipTypeValue = math.MaxInt32
	}

	relationshipObj := &profilev1.RelationshipObject{
		Id:         r.GetID(),
		Type:       profilev1.RelationshipType(relationshipTypeValue),
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
