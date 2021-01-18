package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	papi "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service"
	"github.com/pitabwire/frame"
	"gorm.io/gorm"

	"math/rand"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/ttacon/libphonenumber"

)

type PropertyMap map[string]interface{}

func (p PropertyMap) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

func (p *PropertyMap) Scan(src interface{}) error {

	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		*p = map[string]interface{}{}
	}

	return nil
}

var profileTypeIDMap = map[papi.ProfileType]uint{
	papi.ProfileType_PERSON:      0,
	papi.ProfileType_INSTITUTION: 1,
	papi.ProfileType_BOT:         2,
}

type ProfileType struct {
	ProfileTypeID string `gorm:"type:varchar(50);primary_key"`
	UID           uint   `sql:"unique"`
	Name          string
	Description   string
	frame.BaseModel
}

func (pt *ProfileType) From(db *gorm.DB, profileType papi.ProfileType) {
	pt.UID = profileTypeIDMap[profileType]
	db.First(pt)
}

func (pt *ProfileType) ToEnum() papi.ProfileType {
	for key, val := range profileTypeIDMap {
		if val == pt.UID {
			return key
		}
	}
	return papi.ProfileType_PERSON
}

type Profile struct {

	Properties PropertyMap `sql:"type:jsonb;"`

	ProfileTypeUID uint
	ProfileType    ProfileType
	frame.BaseModel
}

func (p *Profile) GetByID(db *gorm.DB) error {
	modelID := strings.TrimSpace(p.ID)
	if err := db.Where("id = ?", modelID).First(p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return service.ErrorProfileDoesNotExist
		}
		return err
	}
	return nil
}

func (p *Profile) UpdateProperties(db *gorm.DB, params map[string]interface{}) error {

	for key, value := range params {
		if value != "" && value != p.Properties[key] {
			p.Properties[key] = value
		}
	}

	return db.Model(p).Update("Properties", p.Properties).Error
}

func (p *Profile) Create(db *gorm.DB, profileType papi.ProfileType,
	properties map[string]interface{}, ) error {

	pt := ProfileType{}
	pt.From(db, profileType)

	p.ProfileTypeUID = pt.UID
	p.ProfileType = pt
	p.Properties = properties

	return db.Save(p).Error
}

func (p *Profile) ToObject(db *gorm.DB) (*papi.ProfileObject, error) {
	profileObject := papi.ProfileObject{}
	profileObject.ID = p.ID
	profileObject.Type = p.ProfileType.ToEnum()
	profileObject.Properties = map[string]string{}

	for key, val := range p.Properties {
		profileObject.Properties[key] = fmt.Sprintf("%v", val)
	}

	var contactObjects []*papi.ContactObject
	contacts, err := GetContactsByProfile(db, p)
	if err != nil {
		return nil, err
	}

	for _, c := range contacts {
		contactObjects = append(contactObjects, c.ToObject())
	}
	profileObject.Contacts = contactObjects

	var addressObjects []*papi.AddressObject
	addresses, err2 := GetProfileAddresses(db, p)
	if err2 != nil {
		return nil, err2
	}

	for _, a := range addresses {
		addressObjects = append(addressObjects, a.ToObject(db))
	}
	profileObject.Addresses = addressObjects

	return &profileObject, nil

}

var contactTypeUIDMap = map[papi.ContactType]uint{
	papi.ContactType_EMAIL: 0,
	papi.ContactType_PHONE: 1,
}

var communicationLevelUIDMap = map[papi.CommunicationLevel]uint{
	papi.CommunicationLevel_ALL:           0,
	papi.CommunicationLevel_SYSTEM_ALERTS: 1,
	papi.CommunicationLevel_NO_CONTACT:    2,
}

type ContactType struct {
	UID uint `sql:"unique"`

	Name        string
	Description string

	frame.BaseModel
}

func (ct *ContactType) From(db *gorm.DB, contactType papi.ContactType) {
	ct.UID = contactTypeUIDMap[contactType]
	db.First(ct)
}

func (ct *ContactType) FromDetail(db *gorm.DB, detail string) error {

	if govalidator.IsEmail(detail) {

		ct.From(db, papi.ContactType_EMAIL)
		return nil
	} else {

		possibleNumber, err := libphonenumber.Parse(detail, "")

		if err == nil && libphonenumber.IsValidNumber(possibleNumber) {
			ct.From(db, papi.ContactType_PHONE)
			return nil
		}
	}

	return service.ErrorContactDetailsNotValid
}

func (ct *ContactType) ToEnum() papi.ContactType {
	for key, val := range contactTypeUIDMap {
		if val == ct.UID {
			return key
		}
	}
	return papi.ContactType_EMAIL
}

type CommunicationLevel struct {
	UID uint `sql:"unique"`

	Name        string
	Description string

	frame.BaseModel
}

func (cl *CommunicationLevel) From(db *gorm.DB, communicationLevel papi.CommunicationLevel) {
	cl.UID = communicationLevelUIDMap[communicationLevel]
	db.First(cl)
}

func (cl *CommunicationLevel) ToEnum() papi.CommunicationLevel {
	for key, val := range communicationLevelUIDMap {
		if val == cl.UID {
			return key
		}
	}
	return papi.CommunicationLevel_ALL
}

type Contact struct {

	Detail string `gorm:"type:varchar(100);unique_index"`

	ContactTypeUID uint
	ContactType    ContactType

	Language string

	CommunicationLevelUID uint
	CommunicationLevel    CommunicationLevel

	ProfileID string
	Profile   Profile

	frame.BaseModel
}

func (contact *Contact) GetByDetail(db *gorm.DB) error {

	detail := strings.TrimSpace(contact.Detail)
	detail = strings.ToLower(detail)
	if err := db.Last(contact, "detail = ?", detail).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return service.ErrorContactDoesNotExist
		}
		return err
	}
	return nil
}

func GetContactsByProfile(db *gorm.DB, p *Profile) ([]Contact, error) {

	var profileContacts []Contact

	err := db.Where("profile_id = ?", p.ID).Find(&profileContacts).Error

	if err != nil {
		return nil, err
	}

	return profileContacts, nil
}

func (contact *Contact) Create(db *gorm.DB, profileID string, contactDetail string) error {

	detail := strings.TrimSpace(contactDetail)
	detail = strings.ToLower(detail)
	ct := ContactType{}
	err := ct.FromDetail(db, detail)
	if err != nil{
		return err
	}
	contact.ContactType = ct
	contact.ContactTypeUID = ct.UID

	cl := CommunicationLevel{}
	cl.From(db, papi.CommunicationLevel_ALL)
	contact.CommunicationLevel = cl
	contact.CommunicationLevelUID = cl.UID

	contact.ProfileID = profileID
	contact.Detail = contactDetail
	err = db.Save(contact).Error
	if err != nil {
		return err
	}

	return nil

}

func (contact *Contact) ToObject() *papi.ContactObject {

	contactObject := papi.ContactObject{}
	contactObject.ID = contact.ID
	contactObject.Detail = contact.Detail
	contactObject.Type = contact.ContactType.ToEnum()
	contactObject.CommunicationLevel = contact.CommunicationLevel.ToEnum()

	return &contactObject
}

type Verification struct {
	VerificationID string `gorm:"type:varchar(50);primary_key"`
	ContactID      string `gorm:"type:varchar(50);"`
	Contact        Contact

	ProductID string `gorm:"type:varchar(50);"`

	Pin      string `gorm:"type:varchar(10)"`
	LinkHash string `gorm:"type:varchar(100)"`

	ExpiresAt  *time.Time
	VerifiedAt *time.Time
	frame.BaseModel
}

func (v *Verification) Create(db *gorm.DB, productId string, contact Contact, expiryTimeInSec int) error {

	v.ProductID = productId

	v.Contact = contact
	v.ContactID = contact.ID

	v.Pin = GeneratePin(config.ConfigLengthOfVerificationPin)
	v.LinkHash = GeneratePin(config.ConfigLengthOfVerificationLinkHash)

	if expiryTimeInSec > 0 {
		expiryTime := time.Now().Add(time.Duration(expiryTimeInSec))
		v.ExpiresAt = &expiryTime
	}

	return db.Save(v).Error
}

// GeneratePin returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GeneratePin(n int) string {

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

type VerificationAttempt struct {

	VerificationID string `gorm:"type:varchar(50);"`
	Verification   Verification

	ContactID string `gorm:"type:varchar(50);"`
	Contact   Contact

	State string `gorm:"type:varchar(10)"`
	frame.BaseModel
}

type Country struct {
	ISO3 string `gorm:"type:varchar(50);primary_key"`
	Name string
	City string
	ISO2 string `sql:"unique"`

	frame.BaseModel
}

func (country *Country) GetByID(db *gorm.DB, countryID string) error {
	return db.Where("ISO3 = ?", countryID).First(country).Error
}

func (country *Country) GetByAny(db *gorm.DB, c string) error {

	if c == "" {
		return service.ErrorCountryDoesNotExist
	}

	upperC := strings.ToUpper(c)

	return db.Where("ISO3 = ? OR ISO2 = ? OR Name = ?", upperC, upperC, upperC).First(country).Error
}

func (country *Country) From(db *gorm.DB, name string) error {
	return db.Where("name = ?", name).First(country).Error
}

type Address struct {
	Area     string
	Street   string
	House    string
	PostCode string

	Latitude  float64
	Longitude float64

	CountryID string
	Country   Country

	frame.BaseModel
}

func (address *Address) GetByID(db *gorm.DB, addressID string) error {
	return db.Where("ID = ? ", addressID).First(address).Error
}

func (address *Address) GetByAll(db *gorm.DB, countryID string, area, street, house, postcode string,
	latitude, longitude float64, ) error {

	return db.Where("country_id = ? AND area = ? AND street = ? AND house = ? AND postCode = ? AND latitude = ? AND longitude = ?",
		countryID, area, street, house, postcode, latitude, longitude).First(address).Error

}

func (address *Address) Create(db *gorm.DB, countryID string, area, street, house, postcode string,
	latitude, longitude float64, ) error {
	address.Area = area
	address.Street = street
	address.House = house
	address.PostCode = postcode
	address.Latitude = latitude
	address.Longitude = longitude
	address.CountryID = countryID
	return db.Save(address).Error
}

func (address *Address) CreateFull(db *gorm.DB, country, area, street,
	house, postcode string, latitude, longitude float64, ) error {

	countryRecord := Country{}
	if err := countryRecord.GetByAny(db, country); err != nil {
		return err
	}

	addressRecord := Address{}
	if err := addressRecord.GetByAll(db, countryRecord.ISO3, area, street, house,
		postcode, latitude, longitude); err != nil {

		if err != gorm.ErrRecordNotFound {
			return err
		}

		if err2 := addressRecord.Create(db, countryRecord.ISO3, area, street,
			house, postcode, latitude, longitude); err2 != nil {
			return err2
		}
	}

	return nil
}

type ProfileAddress struct {

	Name      string
	AddressID string
	Address   Address

	ProfileID string
	Profile   Profile

	frame.BaseModel
}

func (profileAddress *ProfileAddress) Create(db *gorm.DB, profileID string, addressID string, name string) error {
	profileAddress.ProfileID = profileID
	profileAddress.AddressID = addressID
	profileAddress.Name = name
	return db.Save(profileAddress).Error
}

func (profileAddress *ProfileAddress) ToObject(db *gorm.DB) *papi.AddressObject {

	obj := &papi.AddressObject{}

	if err := profileAddress.Address.GetByID(db, profileAddress.AddressID); err != nil {

	}

	obj.Name = profileAddress.Address.Area
	obj.Area = profileAddress.Address.Area
	obj.Street = profileAddress.Address.Street
	obj.House = profileAddress.Address.House
	obj.Postcode = profileAddress.Address.PostCode
	obj.Latitude = profileAddress.Address.Latitude
	obj.Longitude = profileAddress.Address.Longitude

	if err := profileAddress.Address.Country.GetByID(db, profileAddress.Address.CountryID); err != nil {

	}

	obj.Country = profileAddress.Address.Country.Name
	obj.City = profileAddress.Address.Country.City
	return obj

}

func GetProfileAddresses(db *gorm.DB, p *Profile) ([]ProfileAddress, error) {
	var addresses []ProfileAddress
	if err := db.Where("profile_id = ?", p.ID).Find(&addresses).Error; err != nil {
		return nil, err
	}

	return addresses, nil
}
