package business

import (
	"context"
	"errors"
	"fmt"
	profileV1 "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/service-profile/service"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/pitabwire/frame"
	"strings"
)

type ProfileBusiness interface {
	GetByID(ctx context.Context, encryptionKey []byte, profileID string) (*profileV1.ProfileObject, error)
	GetByContact(ctx context.Context, encryptionKey []byte, detail string) (*profileV1.ProfileObject, error)

	SearchProfile(ctx context.Context, encryptionKey []byte, request *profileV1.ProfileSearchRequest, stream profileV1.ProfileService_SearchServer) error
	CreateProfile(ctx context.Context, encryptionKey []byte, request *profileV1.ProfileCreateRequest) (*profileV1.ProfileObject, error)
	UpdateProfile(ctx context.Context, encryptionKey []byte, request *profileV1.ProfileUpdateRequest) (*profileV1.ProfileObject, error)
	MergeProfile(ctx context.Context, encryptionKey []byte, request *profileV1.ProfileMergeRequest) (*profileV1.ProfileObject, error)

	AddAddress(ctx context.Context, encryptionKey []byte, address *profileV1.ProfileAddAddressRequest) (*profileV1.ProfileObject, error)
	AddContact(ctx context.Context, encryptionKey []byte, contact *profileV1.ProfileAddContactRequest) (*profileV1.ProfileObject, error)
}

func NewProfileBusiness(ctx context.Context, service *frame.Service) ProfileBusiness {
	profileRepo := repository.NewProfileRepository(service)
	contactBusiness := NewContactBusiness(ctx, service)
	addressBusiness := NewAddressBusiness(ctx, service)
	return &profileBusiness{
		service:         service,
		contactBusiness: contactBusiness,
		addressBusiness: addressBusiness,
		profileRepo:     profileRepo,
	}
}

type profileBusiness struct {
	service         *frame.Service
	contactBusiness ContactBusiness
	addressBusiness AddressBusiness

	profileRepo repository.ProfileRepository
}

func (pb *profileBusiness) ProfileToApi(ctx context.Context, p *models.Profile, key []byte) (*profileV1.ProfileObject, error) {
	profileObject := profileV1.ProfileObject{}
	profileObject.ID = p.ID

	profileType, err := pb.profileRepo.GetTypeByID(ctx, p.ProfileTypeID)
	if err != nil {
		return nil, err
	}
	profileObject.Type = models.ProfileTypeIDToEnum(profileType.UID)
	profileObject.Properties = frame.DBPropertiesToMap(p.Properties)

	var contactObjects []*profileV1.ContactObject
	contactList, err := pb.contactBusiness.GetByProfile(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	for _, c := range contactList {

		ctObj, err := pb.contactBusiness.ToApi(ctx, c, key)
		if err != nil {
			return nil, err
		}

		contactObjects = append(contactObjects, ctObj)
	}
	profileObject.Contacts = contactObjects

	var addressObjects []*profileV1.AddressObject
	addressList, err := pb.addressBusiness.GetByProfile(ctx, p.ID)
	for _, a := range addressList {
		address := pb.addressBusiness.ToApi(a.Address)
		addressObjects = append(addressObjects, address)
	}
	profileObject.Addresses = addressObjects

	return &profileObject, nil

}

func (pb *profileBusiness) GetByContact(ctx context.Context, encryptionKey []byte, detail string) (*profileV1.ProfileObject, error) {

	contact, err := pb.contactBusiness.GetByDetail(ctx, detail)
	if err != nil {
		return nil, err
	}

	return pb.GetByID(ctx, encryptionKey, contact.ProfileID)
}

func (pb *profileBusiness) GetByID(ctx context.Context, encryptionKey []byte, profileID string) (*profileV1.ProfileObject, error) {

	profile, err := pb.profileRepo.GetByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return pb.ProfileToApi(ctx, profile, encryptionKey)

}

func (pb *profileBusiness) SearchProfile(ctx context.Context, encryptionKey []byte, request *profileV1.ProfileSearchRequest, stream profileV1.ProfileService_SearchServer) error {

	var profileList []*models.Profile
	//// creating WHERE clause to query by properties JSONB
	scope := pb.service.DB(ctx, true)
	for _, property := range request.GetProperties() {
		column := fmt.Sprintf("properties->>'%s'", property)
		scope = scope.Or(column+" LIKE ?", "%"+request.GetQuery()+"%")
	}

	err := scope.Find(&profileList).Error
	if err != nil {
		return err
	}

	for _, profile := range profileList {

		profileObject, err := pb.ProfileToApi(ctx, profile, encryptionKey)
		err = stream.Send(profileObject)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pb *profileBusiness) MergeProfile(ctx context.Context, encryptionKey []byte, request *profileV1.ProfileMergeRequest) (*profileV1.ProfileObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	target, err := pb.profileRepo.GetByID(ctx, request.GetID())
	if err != nil {
		return nil, err
	}

	merging, err := pb.profileRepo.GetByID(ctx, request.GetMergeID())
	if err != nil {
		return nil, err
	}

	for key, value := range merging.Properties {

		if value == nil || target.Properties[key] == value {
			continue
		}
		target.Properties[key] = value
	}

	err = pb.profileRepo.Save(ctx, target)
	if err != nil {
		return nil, err
	}

	err = pb.profileRepo.Delete(ctx, merging.GetID())
	if err != nil {
		return nil, err
	}

	return pb.ProfileToApi(ctx, target, encryptionKey)
}

func (pb *profileBusiness) UpdateProfile(ctx context.Context, encryptionKey []byte, request *profileV1.ProfileUpdateRequest) (*profileV1.ProfileObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	profile, err := pb.profileRepo.GetByID(ctx, request.GetID())
	if err != nil {
		return nil, err
	}

	properties := frame.DBPropertiesFromMap(request.GetProperties())
	for key, value := range properties {
		if value != profile.Properties[key] {
			profile.Properties[key] = value
		}
	}

	err = pb.profileRepo.Save(ctx, profile)
	if err != nil {
		return nil, err
	}

	return pb.ProfileToApi(ctx, profile, encryptionKey)
}

func (pb *profileBusiness) CreateProfile(ctx context.Context, encryptionKey []byte, request *profileV1.ProfileCreateRequest) (*profileV1.ProfileObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	contactDetail := strings.TrimSpace(request.GetContact())

	if contactDetail == "" {
		return nil, service.ErrorContactDetailsNotValid
	}

	p := models.Profile{}
	p.Properties = frame.DBPropertiesFromMap(request.GetProperties())

	contact, err := pb.contactBusiness.GetByDetail(ctx, contactDetail)
	if err != nil {
		if !errors.Is(service.ErrorContactDoesNotExist, err) {
			return nil, err
		}

		pt, err := pb.profileRepo.GetTypeByUID(ctx, request.GetType())
		if err != nil {
			return nil, err
		}

		p.ProfileType = pt
		p.ProfileTypeID = pt.ID

		err = pb.profileRepo.Save(ctx, &p)
		if err != nil {
			return nil, err
		}

		err = pb.contactBusiness.CreateContact(ctx, encryptionKey, p.GetID(), contactDetail)
		if err != nil {
			return nil, err
		}

		return pb.GetByID(ctx, encryptionKey, p.GetID())
	}

	return pb.GetByID(ctx, encryptionKey, contact.ProfileID)

}

//func (pb *profileBusiness) UpdateProperties(db *gorm.DB, params map[string]interface{}) error {
//
//	storedPropertiesMap := make(map[string]interface{})
//	attributeMap, err := p.Properties.MarshalJSON()
//	if err != nil {
//		return err
//	}
//
//	err = json.Unmarshal(attributeMap, &storedPropertiesMap)
//	if err != nil {
//		return err
//	}
//
//	for key, value := range params {
//		if value != nil && value != "" && value != storedPropertiesMap[key] {
//			storedPropertiesMap[key] = value
//		}
//	}
//
//	stringProperties, err := json.Marshal(storedPropertiesMap)
//	if err != nil {
//		return err
//	}
//
//	err = p.Properties.UnmarshalJSON(stringProperties)
//	if err != nil {
//		return err
//	}
//
//	return db.Model(p).Update("Properties", p.Properties).Error
//}

func (pb *profileBusiness) AddAddress(ctx context.Context, encryptionKey []byte, request *profileV1.ProfileAddAddressRequest) (*profileV1.ProfileObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	address, err := pb.addressBusiness.CreateAddress(ctx, request.GetAddress())
	if err != nil {
		return nil, err
	}

	err = pb.addressBusiness.LinkAddressToProfile(ctx,
		request.GetID(), request.GetAddress().GetExtra(), address)

	return pb.GetByID(ctx, encryptionKey, request.GetID())

}

func (pb *profileBusiness) AddContact(ctx context.Context, encryptionKey []byte, request *profileV1.ProfileAddContactRequest) (*profileV1.ProfileObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return pb.GetByID(ctx, encryptionKey, request.GetID())

}
