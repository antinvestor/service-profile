package models

import (
	"fmt"
	"math"
	"time"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/util"
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
//
//nolint:gochecknoglobals // This is a mapping table that needs to be global
var ProfileTypeIDMap = map[profilev1.ProfileType]uint{
	profilev1.ProfileType_PERSON:      ProfileTypePersonID,
	profilev1.ProfileType_BOT:         ProfileTypeBotID,
	profilev1.ProfileType_INSTITUTION: ProfileTypeInstitutionID,
}

// RelationshipTypeIDMap maps relationship types to their respective IDs.
//
//nolint:gochecknoglobals // This is a mapping table that needs to be global
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
	data.BaseModel
	UID         uint `sql:"unique"`
	Name        string
	Description string
}

type Profile struct {
	data.BaseModel
	Properties data.JSONMap

	ProfileTypeID string `gorm:"type:varchar(50);index:profile_id"`
	ProfileType   ProfileType
}

type Contact struct {
	data.BaseModel

	LookUpToken     []byte `gorm:"type:bytea;uniqueIndex"`
	EncryptedDetail []byte `gorm:"type:bytea"`
	EncryptionKeyID string `gorm:"type:varchar(255)"`

	ContactType        string `gorm:"type:varchar(50)"`
	CommunicationLevel string `gorm:"type:varchar(50)"`

	Language string

	ProfileID string `gorm:"type:varchar(50);index:profile_id"`

	Properties data.JSONMap

	VerificationID string `gorm:"type:varchar(50)"`
}

func (c *Contact) DecryptDetail(decryptionKeyID string, decryptionKeyData []byte) (string, error) {

	if c.EncryptionKeyID != decryptionKeyID {
		return "", fmt.Errorf("decryption key does not match contact key id")
	}

	detailBytes, err := util.DecryptValue(decryptionKeyData, c.EncryptedDetail)
	if err != nil {
		return "", err
	}
	return string(detailBytes), nil
}

func (c *Contact) ToAPI(dek *config.DEK, partial bool) (*profilev1.ContactObject, error) {

	contactDetail, err := c.DecryptDetail(dek.KeyID, dek.Key)
	if err != nil {
		return nil, err
	}

	contactObject := profilev1.ContactObject{
		Id:     c.ID,
		Detail: contactDetail,
	}

	contactTypeID, ok := profilev1.ContactType_value[c.ContactType]
	if !ok {
		contactTypeID = int32(profilev1.ContactType_EMAIL)
	}
	contactObject.Type = profilev1.ContactType(contactTypeID)

	communicationLevel, ok := profilev1.CommunicationLevel_value[c.CommunicationLevel]
	if !ok {
		communicationLevel = int32(profilev1.CommunicationLevel_ALL)
	}
	contactObject.CommunicationLevel = profilev1.CommunicationLevel(communicationLevel)

	contactObject.Verified = false
	if !partial {
		contactObject.Verified = c.VerificationID != ""
	}

	return &contactObject, nil
}

type Roster struct {
	data.BaseModel

	ProfileID string `gorm:"type:varchar(50);uniqueIndex:roster_composite_index;index:profile_id"`

	ContactID string `gorm:"type:varchar(50);uniqueIndex:roster_composite_index;"`
	Contact   *Contact

	Properties data.JSONMap
}

func (r *Roster) ToAPI(dek *config.DEK) (*profilev1.RosterObject, error) {

	contactObj, err := r.Contact.ToAPI(dek, true)
	if err != nil {
		return nil, err
	}

	return &profilev1.RosterObject{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		Contact:   contactObj,
		Extra:     r.Properties.ToProtoStruct(),
	}, nil
}

type Verification struct {
	data.BaseModel
	ProfileID  string    `gorm:"type:varchar(50);index:profile_id" json:"profile_id"`
	ContactID  string    `gorm:"type:varchar(50);index:contact_id" json:"contact_id"`
	Code       string    `gorm:"type:varchar(255)"                 json:"code"`
	ExpiresAt  time.Time `                                         json:"expires_at"`
	VerifiedAt time.Time `                                         json:"verified_at"`
}

type VerificationAttempt struct {
	data.BaseModel
	VerificationID string `gorm:"type:varchar(50);index:verification_id"`

	Data  string `gorm:"type:varchar(250)"`
	State string `gorm:"type:varchar(10)"`

	DeviceID  string `gorm:"type:varchar(50)"`
	IPAddress string `gorm:"type:varchar(50)"`
	RequestID string `gorm:"type:varchar(50)"`
}

type Country struct {
	data.BaseModel
	ISO3         string `gorm:"unique"`
	ISO2         string `              sql:"unique"`
	Name         string
	Code         int
	LatitudeAvg  float64
	LongitudeAvg float64
	City         string
}

type Address struct {
	data.BaseModel
	Name      string
	AdminUnit string

	ParentID string `gorm:"type:varchar(50);index:parent_id"`

	CountryID string `gorm:"type:varchar(50)"`
	Country   *Country

	Properties data.JSONMap
}

type ProfileAddress struct {
	data.BaseModel
	Name string

	AddressID string `gorm:"type:varchar(50);index:address_id"`
	Address   *Address

	ProfileID string `gorm:"type:varchar(50);index:profile_id"`
	Profile   *Profile
}

type RelationshipType struct {
	data.BaseModel
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
	data.BaseModel

	ParentObject   string `gorm:"type:varchar(50)"`
	ParentObjectID string `gorm:"type:varchar(50);index:parent_obj_id"`

	ChildObject   string `gorm:"type:varchar(50)"`
	ChildObjectID string `gorm:"type:varchar(50);index:child_obj_id"`

	RelationshipTypeID string `gorm:"type:varchar(50);index:relationship_type_id"`
	RelationshipType   *RelationshipType

	Properties data.JSONMap
}

func (r *Relationship) ToAPI() *profilev1.RelationshipObject {
	// Safe conversion from uint to int32
	var relationshipTypeValue int32
	if r.RelationshipType.UID <= uint(math.MaxInt32) {
		relationshipTypeValue = int32(r.RelationshipType.UID) // #nosec G115 -- bounds checked above
	} else {
		relationshipTypeValue = math.MaxInt32
	}

	relationshipObj := &profilev1.RelationshipObject{
		Id:         r.GetID(),
		Type:       profilev1.RelationshipType(relationshipTypeValue),
		Properties: r.Properties.ToProtoStruct(),
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
