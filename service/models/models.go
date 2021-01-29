package models

import (
	"encoding/json"
	papi "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service"
	"github.com/go-errors/errors"
	"github.com/pitabwire/frame"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"math/rand"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/ttacon/libphonenumber"
)

var profileTypeIDMap = map[papi.ProfileType]uint{
	papi.ProfileType_PERSON:      0,
	papi.ProfileType_INSTITUTION: 1,
	papi.ProfileType_BOT:         2,
}

type ProfileType struct {
	frame.BaseModel
	UID         uint `sql:"unique"`
	Name        string
	Description string
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

	frame.BaseModel
	Properties datatypes.JSON

	ProfileTypeID string `gorm:"type:varchar(50)"`
	ProfileType ProfileType
}

func (p *Profile) GetByID(db *gorm.DB) error {
	modelID := strings.TrimSpace(p.ID)
	if err := db.Preload(clause.Associations).Where("id = ?", modelID).First(p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return service.ErrorProfileDoesNotExist
		}
		return errors.Wrap(err, 1)
	}
	return nil
}

func (p *Profile) UpdateProperties(db *gorm.DB, params map[string]interface{}) error {

	storedPropertiesMap := make(map[string]interface{})
	attributeMap, err := p.Properties.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, 1)
	}

	err = json.Unmarshal(attributeMap, &storedPropertiesMap)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	for key, value := range params {
		if value != nil && value != "" && value != storedPropertiesMap[key] {
			storedPropertiesMap[key] = value
		}
	}

	stringProperties, err := json.Marshal(storedPropertiesMap)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	err = p.Properties.UnmarshalJSON(stringProperties)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	return db.Model(p).Update("Properties", p.Properties).Error
}

func (p *Profile) Create(db *gorm.DB, profileType papi.ProfileType, properties map[string]interface{}, ) error {

	pt := ProfileType{}
	pt.From(db, profileType)

	p.ProfileType = pt
	p.ProfileTypeID = pt.ID

	stringProperties, err := json.Marshal(properties)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	err = p.Properties.UnmarshalJSON(stringProperties)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	return db.Save(p).Error
}

func (p *Profile) ToObject(db *gorm.DB) (*papi.ProfileObject, error) {
	profileObject := papi.ProfileObject{}
	profileObject.ID = p.ID
	profileObject.Type = p.ProfileType.ToEnum()
	profileObject.Properties = map[string]string{}

	attributeMap, err := p.Properties.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, 1)
	}

	err = json.Unmarshal(attributeMap, &profileObject.Properties)
	if err != nil {
		return nil, errors.Wrap(err, 1)
	}

	var contactObjects []*papi.ContactObject
	contacts, err := GetContactsByProfile(db, p)
	if err != nil {
		return nil, errors.Wrap(err, 1)
	}

	for _, c := range contacts {
		contactObjects = append(contactObjects, c.ToObject())
	}
	profileObject.Contacts = contactObjects

	var addressObjects []*papi.AddressObject
	addresses, err2 := GetProfileAddresses(db, p)
	if err2 != nil {
		return nil, errors.Wrap(err2, 1)
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

	frame.BaseModel
	UID uint `sql:"unique"`

	Name        string
	Description string

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
	frame.BaseModel
	UID uint `sql:"unique"`

	Name        string
	Description string

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
	frame.BaseModel
	Detail string `gorm:"type:varchar(100);unique"`

	ContactTypeID string `gorm:"type:varchar(50)"`
	ContactType        ContactType

	CommunicationLevelID string `gorm:"type:varchar(50)"`
	CommunicationLevel CommunicationLevel

	Language           string

	ProfileID string `gorm:"type:varchar(50)"`
	Profile            Profile

}

func (contact *Contact) GetByDetail(db *gorm.DB) error {

	detail := strings.TrimSpace(contact.Detail)
	detail = strings.ToLower(detail)
	if err := db.Preload(clause.Associations).Last(contact, "detail = ?", detail).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return service.ErrorContactDoesNotExist
		}
		return errors.Wrap(err, 1)
	}
	return nil
}

func GetContactsByProfile(db *gorm.DB, p *Profile) ([]Contact, error) {

	var profileContacts []Contact

	err := db.Preload(clause.Associations).Where("profile_id = ?", p.ID).Find(&profileContacts).Error

	if err != nil {
		return nil, errors.Wrap(err, 1)
	}

	return profileContacts, nil
}

func (contact *Contact) Create(db *gorm.DB, profileID string, contactDetail string) error {

	detail := strings.TrimSpace(contactDetail)
	detail = strings.ToLower(detail)
	ct := ContactType{}
	err := ct.FromDetail(db, detail)
	if err != nil {
		return errors.Wrap(err, 1)
	}
	contact.ContactTypeID = ct.ID
	contact.ContactType = ct

	cl := CommunicationLevel{}
	cl.From(db, papi.CommunicationLevel_ALL)
	contact.CommunicationLevelID = cl.ID
	contact.CommunicationLevel = cl

	profile := Profile{}
	profile.ID = profileID
	err = profile.GetByID(db)
	if err != nil {
		return errors.Wrap(err, 1)
	}
	contact.Profile = profile
	contact.ProfileID = profile.ID
	contact.Detail = contactDetail
	err = db.Save(contact).Error
	if err != nil {
		return errors.Wrap(err, 1)
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

	frame.BaseModel
	ContactID string `gorm:"type:varchar(50)"`
	Contact Contact

	ProductID string `gorm:"type:varchar(50);"`

	Pin      string `gorm:"type:varchar(10)"`
	LinkHash string `gorm:"type:varchar(100)"`

	ExpiresAt  *time.Time
	VerifiedAt *time.Time
}

func (v *Verification) Create(db *gorm.DB, productId string, contact Contact, expiryTimeInSec int) error {

	v.ProductID = productId

	v.ContactID = contact.ID
	v.Contact = contact

	v.Pin = GeneratePin(config.LengthOfVerificationPin)
	v.LinkHash = GeneratePin(config.LengthOfVerificationLinkHash)

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

	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

type VerificationAttempt struct {

	frame.BaseModel
	VerificationID string `gorm:"type:varchar(50)"`
	Verification Verification

	ContactID string `gorm:"type:varchar(50)"`
	Contact      Contact

	State string `gorm:"type:varchar(10)"`
}

type Country struct {
	frame.BaseModel
	ISO3 string `gorm:"unique"`
	ISO2 string `sql:"unique"`
	Name string
	City string
}

func (country *Country) GetByISO3(db *gorm.DB, countryISO3 string) error {
	return db.Where("ISO3 = ?", countryISO3).First(country).Error
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
	frame.BaseModel
	Area     string
	Street   string
	House    string
	PostCode string

	Latitude  float64
	Longitude float64

	CountryID string `gorm:"type:varchar(50)"`
	Country Country
}

func (address *Address) GetByID(db *gorm.DB, addressID string) error {
	return db.Preload(clause.Associations).Where("ID = ? ", addressID).First(address).Error
}

func (address *Address) GetByAll(db *gorm.DB, countryID string, area, street, house, postcode string,
	latitude, longitude float64, ) error {

	return db.Preload(clause.Associations).Where("country_id = ? AND area = ? AND street = ? AND house = ? AND postCode = ? AND latitude = ? AND longitude = ?",
		countryID, area, street, house, postcode, latitude, longitude).First(address).Error

}

func (address *Address) Create(db *gorm.DB, countryISO3 string, area, street, house, postcode string,
	latitude, longitude float64, ) error {

	country := Country{ISO3: countryISO3}
	err := country.GetByISO3(db, countryISO3)
	if err != nil{
		return errors.Wrap(err, 1)
	}
	address.Area = area
	address.Street = street
	address.House = house
	address.PostCode = postcode
	address.Latitude = latitude
	address.Longitude = longitude
	address.Country = country
	return db.Save(address).Error
}

func (address *Address) CreateFull(db *gorm.DB, country, area, street,
	house, postcode string, latitude, longitude float64, ) error {

	countryRecord := Country{}
	if err := countryRecord.GetByAny(db, country); err != nil {
		return errors.Wrap(err, 1)
	}

	addressRecord := Address{}
	if err := addressRecord.GetByAll(db, countryRecord.ISO3, area, street, house,
		postcode, latitude, longitude); err != nil {

		if err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, 1)
		}

		if err2 := addressRecord.Create(db, countryRecord.ISO3, area, street,
			house, postcode, latitude, longitude); err2 != nil {
			return errors.Wrap(err2, 1)
		}
	}

	return nil
}

type ProfileAddress struct {

	frame.BaseModel
	Name    string

	AddressID string `gorm:"type:varchar(50)"`
	Address Address

	ProfileID string `gorm:"type:varchar(50)"`
	Profile Profile

}

func (profileAddress *ProfileAddress) Create(db *gorm.DB, profile Profile, address Address, name string) error {
	profileAddress.Profile = profile
	profileAddress.Address = address
	profileAddress.Name = name
	return db.Save(profileAddress).Error
}

func (profileAddress *ProfileAddress) ToObject(db *gorm.DB) *papi.AddressObject {

	obj := &papi.AddressObject{}

	obj.Name = profileAddress.Address.Area
	obj.Area = profileAddress.Address.Area
	obj.Street = profileAddress.Address.Street
	obj.House = profileAddress.Address.House
	obj.Postcode = profileAddress.Address.PostCode
	obj.Latitude = profileAddress.Address.Latitude
	obj.Longitude = profileAddress.Address.Longitude

	obj.Country = profileAddress.Address.Country.Name
	obj.City = profileAddress.Address.Country.City
	return obj

}

func GetProfileAddresses(db *gorm.DB, p *Profile) ([]ProfileAddress, error) {
	var addresses []ProfileAddress
	if err := db.Preload(clause.Associations).Where("profile_id = ?", p.ID).Find(&addresses).Error; err != nil {
		return nil, errors.Wrap(err, 1)
	}

	return addresses, nil
}
