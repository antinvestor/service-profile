package service

import (
	"bitbucket.org/antinvestor/service-profile/profile"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/xid"

	"database/sql/driver"
	"encoding/json"
	"errors"

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

// AntMigration Our simple table holding all the migration data
type AntMigration struct {
	AntMigrationID string `gorm:"type:varchar(50);primary_key"`
	Name           string `gorm:"type:varchar(50);unique_index"`
	Patch          string `gorm:"type:text"`
	AppliedAt      *time.Time
	CreatedAt      time.Time
	ModifiedAt     time.Time
	Version        uint32 `gorm:"DEFAULT 0"`
}

// BeforeCreate Ensures we update a migrations time stamps
func (model *AntMigration) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("AntMigrationID", xid.New().String())
	scope.SetColumn("CreatedAt", time.Now())
	return scope.SetColumn("ModifiedAt", time.Now())
}

// BeforeUpdate Updates time stamp every time we update status of a migration
func (model *AntMigration) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("Version", model.Version+1)
	return scope.SetColumn("ModifiedAt", time.Now())
}

var profileTypeIDMap = map[profile.ProfileType]uint{
	profile.ProfileType_PERSON:      0,
	profile.ProfileType_INSTITUTION: 1,
	profile.ProfileType_BOT:         2,
}

type ProfileType struct {
	ProfileTypeID string `gorm:"type:varchar(50);primary_key"`
	UID           uint   `sql:"unique"`
	Name          string
	Description   string
	CreatedAt     time.Time
	ModifiedAt    time.Time
	DeletedAt     *time.Time `sql:"index"`
	Version       uint32     `gorm:"DEFAULT 0"`
}

func (pt *ProfileType) From(db *gorm.DB, profileType profile.ProfileType) {
	pt.UID = profileTypeIDMap[profileType]
	db.First(pt)
}

func (pt *ProfileType) ToEnum() profile.ProfileType {
	for key, val := range profileTypeIDMap {
		if val == pt.UID {
			return key
		}
	}
	return profile.ProfileType_PERSON
}

type Profile struct {
	ProfileID string `gorm:"type:varchar(50);primary_key"`

	Properties PropertyMap `sql:"type:jsonb;"`

	ProfileTypeUID uint
	ProfileType    ProfileType

	CreatedAt  time.Time
	ModifiedAt time.Time
	DeletedAt  *time.Time `sql:"index"`
	Version    uint32     `gorm:"DEFAULT 0"`
}

// BeforeCreate Ensures we update a migrations time stamps
func (p *Profile) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ProfileID", xid.New().String())
	scope.SetColumn("CreatedAt", time.Now())
	return scope.SetColumn("ModifiedAt", time.Now())
}

// BeforeUpdate Updates time stamp every time we update status of a migration
func (p *Profile) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("Version", p.Version+1)
	return scope.SetColumn("ModifiedAt", time.Now())
}

func (p *Profile) GetByID(db *gorm.DB) error {
	modelID := strings.TrimSpace(p.ProfileID)
	if db.Where("ProfileID = ?", modelID).First(p).RecordNotFound() {
		return profile.ErrorProfileDoesNotExist
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

func (p *Profile) Create(db *gorm.DB, profileType profile.ProfileType,
	contactDetail string, properties map[string]interface{}, ) error {

	if contactDetail == "" {
		return profile.ErrorEmptyValueSupplied
	}

	contact := Contact{Detail: contactDetail}

	err := contact.GetByDetail(db)

	if err != nil {

		if err != profile.ErrorContactDoesNotExist {
			return err
		}

		pt := ProfileType{}
		pt.From(db, profileType)

		p.ProfileTypeUID = pt.UID
		p.ProfileType = pt
		p.Properties = properties

		err := db.Save(p).Error
		if err != nil {
			return err
		}

		return contact.Create(db, p.ProfileID, contactDetail)

	}

	p.ProfileID = contact.ProfileID

	err2 := db.First(p).Error
	if err2 == nil {
		return err2
	}

	return p.UpdateProperties(db, properties)

}

func (p *Profile) ToObject(db *gorm.DB) (*profile.ProfileObject, error) {
	profileObject := profile.ProfileObject{}
	profileObject.ID = p.ProfileID
	profileObject.Type = p.ProfileType.ToEnum()
	profileObject.Properties = map[string]string{}

	for key, val := range p.Properties {
		profileObject.Properties[key] = fmt.Sprintf("%v", val)
	}

	var contactObjects []*profile.ContactObject
	contacts, err := GetContactsByProfile(db, p)
	if err != nil {
		return nil, err
	}

	for _, c := range contacts {
		contactObjects = append(contactObjects, c.ToObject())
	}
	profileObject.Contacts = contactObjects

	var addressObjects []*profile.AddressObject
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

var contactTypeUIDMap = map[profile.ContactType]uint{
	profile.ContactType_EMAIL: 0,
	profile.ContactType_PHONE: 1,
}

var communicationLevelUIDMap = map[profile.CommunicationLevel]uint{
	profile.CommunicationLevel_ALL:           0,
	profile.CommunicationLevel_SYSTEM_ALERTS: 1,
	profile.CommunicationLevel_NO_CONTACT:    2,
}

type ContactType struct {
	ContactTypeID string `gorm:"type:varchar(50);primary_key"`

	UID uint `sql:"unique"`

	Name        string
	Description string

	CreatedAt  time.Time
	ModifiedAt time.Time
	DeletedAt  *time.Time `sql:"index"`
	Version    uint32     `gorm:"DEFAULT 0"`
}

func (ct *ContactType) From(db *gorm.DB, contactType profile.ContactType) {
	ct.UID = contactTypeUIDMap[contactType]
	db.First(ct)
}

func (ct *ContactType) FromDetail(db *gorm.DB, detail string) error {

	if govalidator.IsEmail(detail) {

		ct.From(db, profile.ContactType_EMAIL)
		return nil
	} else {

		possibleNumber, err := libphonenumber.Parse(detail, "")

		if err == nil && libphonenumber.IsValidNumber(possibleNumber) {
			ct.From(db, profile.ContactType_PHONE)
			return nil
		}
	}

	return profile.ErrorContactDetailsNotValid
}

func (ct *ContactType) ToEnum() profile.ContactType {
	for key, val := range contactTypeUIDMap {
		if val == ct.UID {
			return key
		}
	}
	return profile.ContactType_EMAIL
}

type CommunicationLevel struct {
	CommunicationLevelID string `gorm:"type:varchar(50);primary_key"`

	UID uint `sql:"unique"`

	Name        string
	Description string

	CreatedAt  time.Time
	ModifiedAt time.Time
	DeletedAt  *time.Time `sql:"index"`
	Version    uint32     `gorm:"DEFAULT 0"`
}

func (cl *CommunicationLevel) From(db *gorm.DB, communicationLevel profile.CommunicationLevel) {
	cl.UID = communicationLevelUIDMap[communicationLevel]
	db.First(cl)
}

func (cl *CommunicationLevel) ToEnum() profile.CommunicationLevel {
	for key, val := range communicationLevelUIDMap {
		if val == cl.UID {
			return key
		}
	}
	return profile.CommunicationLevel_ALL
}

type Contact struct {
	ContactID string `gorm:"type:varchar(50);primary_key"`

	Detail string `gorm:"type:varchar(100);unique_index"`

	ContactTypeUID uint
	ContactType    ContactType

	CommunicationLevelUID uint
	CommunicationLevel    CommunicationLevel

	ProfileID string
	Profile   Profile

	CreatedAt  time.Time
	ModifiedAt time.Time
	DeletedAt  *time.Time `sql:"index"`
	Version    uint32     `gorm:"DEFAULT 0"`
}

// BeforeCreate Ensures we update a migrations time stamps
func (model *Contact) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ContactID", xid.New().String())
	scope.SetColumn("CreatedAt", time.Now())
	return scope.SetColumn("ModifiedAt", time.Now())
}

// BeforeUpdate Updates time stamp every time we update status of a migration
func (model *Contact) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("Version", model.Version+1)
	return scope.SetColumn("ModifiedAt", time.Now())
}

func (c *Contact) GetByDetail(db *gorm.DB) error {

	detail := strings.TrimSpace(c.Detail)
	detail = strings.ToLower(detail)
	if db.Find(c, "detail = ?", detail).RecordNotFound() {
		return profile.ErrorContactDoesNotExist
	}
	return nil
}

func GetContactsByProfile(db *gorm.DB, p *Profile) ([]Contact, error) {

	var profileContacts []Contact

	err := db.Where("profile_id = ?", p.ProfileID).Find(&profileContacts).Error

	if err != nil {
		return nil, err
	}

	return profileContacts, nil
}

func (c *Contact) Create(db *gorm.DB, profileID string, contactDetail string) error {

	detail := strings.TrimSpace(c.Detail)
	detail = strings.ToLower(detail)
	ct := ContactType{}
	ct.FromDetail(db, detail)
	c.ContactType = ct
	c.ContactTypeUID = ct.UID

	cl := CommunicationLevel{}
	cl.From(db, profile.CommunicationLevel_ALL)
	c.CommunicationLevel = cl
	c.CommunicationLevelUID = cl.UID

	c.ProfileID = profileID
	c.Detail = contactDetail
	return db.Save(c).Error
}

func (c *Contact) ToObject() *profile.ContactObject {

	contactObject := profile.ContactObject{}
	contactObject.ID = c.ContactID
	contactObject.Detail = c.Detail
	contactObject.Type = c.ContactType.ToEnum()
	contactObject.CommunicationLevel = c.CommunicationLevel.ToEnum()

	return &contactObject
}

type Country struct {
	CountryID string `gorm:"type:varchar(50);primary_key"`
	Name      string
	City      string
	ISO2      string
	ISO3      string

	CreatedAt  time.Time
	ModifiedAt time.Time
	DeletedAt  *time.Time `sql:"index"`
	Version    uint32     `gorm:"DEFAULT 0"`
}

func (country *Country) GetByID(db *gorm.DB, countryID string) error {
	return db.Where("CountryID = ?", countryID).First(country).Error
}

func (country *Country) GetByISO3(db *gorm.DB, iso3Code string) error {

	if iso3Code == "" {
		return profile.ErrorCountryDoesNotExist
	}

	return db.Where("ISO3 = ?", iso3Code).First(country).Error
}

func (country *Country) From(db *gorm.DB, name string) error {
	return db.Where("name = ?", name).First(country).Error
}

func (country *Country) Create(db *gorm.DB, name string, iso2 string, iso3 string) error {

	err := country.GetByISO3(db, iso3)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			country.Name = name
			country.ISO2 = iso2
			country.ISO3 = iso3
			return db.Save(country).Error
		}
	}

	return nil
}

type Address struct {
	AddressID string `gorm:"type:varchar(50);primary_key"`

	Area     string
	Street   string
	House    string
	PostCode string

	Latitude  float64
	Longitude float64

	CountryID string
	Country   Country

	CreatedAt  time.Time
	ModifiedAt time.Time
	DeletedAt  *time.Time `sql:"index"`
	Version    uint32     `gorm:"DEFAULT 0"`
}

func (address *Address) GetByID(db *gorm.DB, addressID string) error {
	return db.Where("AddressID = ? ").First(address).Error
}

func (address *Address) GetByAll(db *gorm.DB, countryID string, area, street, house, postcode string,
	latitude, longitude float64, ) error {

	return db.Where("CountryID = ? AND Area = ? AND Street = ? AND House = ? AND PostCode = ? AND Latitude = ? AND Longitude = ?",
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

func (address *Address) CreateFull(db *gorm.DB, country, town, location, area, street,
	house, postcode string, latitude, longitude float64, ) error {

	countryRecord := Country{}
	if err := countryRecord.GetByISO3(db, country); err != nil {
		return err
	}

	addressRecord := Address{}
	if err := addressRecord.GetByAll(db, countryRecord.CountryID, area, street, house,
		postcode, latitude, longitude, ); err != nil {

		if err != gorm.ErrRecordNotFound {
			return err
		}

		if err2 := addressRecord.Create(db, countryRecord.CountryID, area, street,
			house, postcode, latitude, longitude, ); err2 != nil {
			return err2
		}
	}

	return nil
}

type ProfileAddress struct {
	ProfileAddressID string `gorm:"type:varchar(50);primary_key"`

	Name      string
	AddressID string
	Address   Address

	ProfileID string
	Profile   Profile

	CreatedAt  time.Time
	ModifiedAt time.Time
	DeletedAt  *time.Time `sql:"index"`
	Version    uint32     `gorm:"DEFAULT 0"`
}

func (profileAddress *ProfileAddress) Create(db *gorm.DB, profileID string, addressID string, name string) error {

	profileAddress.ProfileID = profileID
	profileAddress.AddressID = addressID
	profileAddress.Name = name
	return db.Save(profileAddress).Error
}

func (profileAddress *ProfileAddress) ToObject(db *gorm.DB) *profile.AddressObject {

	obj := &profile.AddressObject{}

	if err := profileAddress.Address.GetByID(db, profileAddress.AddressID); err != nil {

		if err == gorm.ErrRecordNotFound {
			err = profile.ErrorAddressDoesNotExist
		}
	}

	obj.Name = profileAddress.Address.Area
	obj.Area = profileAddress.Address.Area
	obj.Street = profileAddress.Address.Street
	obj.House = profileAddress.Address.House
	obj.Postcode = profileAddress.Address.PostCode
	obj.Latitude = profileAddress.Address.Latitude
	obj.Longitude = profileAddress.Address.Longitude

	if err := profileAddress.Address.Country.GetByID(db, profileAddress.Address.CountryID); err != nil {

		if err == gorm.ErrRecordNotFound {
			err = profile.ErrorCountryDoesNotExist
		}
	}

	obj.Country = profileAddress.Address.Country.Name
	obj.City = profileAddress.Address.Country.City
	return obj

}

func GetProfileAddresses(db *gorm.DB, p *Profile) ([]ProfileAddress, error) {
	var addresses []ProfileAddress
	if err := db.Where("ProfileID = ?", p.ProfileID).Find(&addresses).Error; err != nil {
		return nil, err
	}

	return addresses, nil
}
