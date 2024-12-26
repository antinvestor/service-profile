package business

import (
	"context"
	"errors"
	"fmt"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/pitabwire/frame"
	"github.com/rs/xid"
	"strings"
)

type ProfileBusiness interface {
	GetByID(ctx context.Context, profileID string) (*profilev1.ProfileObject, error)
	GetByContact(ctx context.Context, detail string) (*profilev1.ProfileObject, error)

	SearchProfile(ctx context.Context, request *profilev1.SearchRequest, stream profilev1.ProfileService_SearchServer) error

	CreateProfile(ctx context.Context, request *profilev1.CreateRequest) (*profilev1.ProfileObject, error)

	UpdateProfile(ctx context.Context, request *profilev1.UpdateRequest) (*profilev1.ProfileObject, error)

	MergeProfile(ctx context.Context, request *profilev1.MergeRequest) (*profilev1.ProfileObject, error)

	AddAddress(ctx context.Context, address *profilev1.AddAddressRequest) (*profilev1.ProfileObject, error)

	AddContact(ctx context.Context, contact *profilev1.AddContactRequest) (*profilev1.ProfileObject, error)

	EncryptionKeyFunc() []byte
}

func NewProfileBusiness(ctx context.Context, service *frame.Service, encryptionKeyFunc func() []byte) ProfileBusiness {
	return &profileBusiness{
		service:         service,
		encryptionKey:   encryptionKeyFunc(),
		contactBusiness: NewContactBusiness(ctx, service),
		addressBusiness: NewAddressBusiness(ctx, service),
		profileRepo:     repository.NewProfileRepository(service),
	}
}

type profileBusiness struct {
	service *frame.Service

	encryptionKey []byte

	contactBusiness ContactBusiness
	addressBusiness AddressBusiness

	profileRepo repository.ProfileRepository
}

func (pb *profileBusiness) EncryptionKeyFunc() []byte {
	return pb.encryptionKey
}

func (pb *profileBusiness) ProfileToAPI(ctx context.Context,
	p *models.Profile) (*profilev1.ProfileObject, error) {
	profileObject := profilev1.ProfileObject{}
	profileObject.Id = p.ID

	profileObject.Type = models.ProfileTypeIDToEnum(p.ProfileType.UID)
	profileObject.Properties = frame.DBPropertiesToMap(p.Properties)

	var contactObjects []*profilev1.ContactObject
	contactList, err := pb.contactBusiness.GetByProfile(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	for _, c := range contactList {
		ctObj, err := pb.contactBusiness.ToAPI(ctx, c, pb.EncryptionKeyFunc())
		if err != nil {
			return nil, err
		}

		contactObjects = append(contactObjects, ctObj)
	}
	profileObject.Contacts = contactObjects

	var addressObjects []*profilev1.AddressObject
	addressList, err := pb.addressBusiness.GetByProfile(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	for _, a := range addressList {
		address := pb.addressBusiness.ToAPI(a.Address)
		addressObjects = append(addressObjects, address)
	}
	profileObject.Addresses = addressObjects

	return &profileObject, nil
}

func (pb *profileBusiness) GetByContact(
	ctx context.Context,
	contactData string) (*profilev1.ProfileObject, error) {

	ctx = frame.SkipTenancyChecksOnClaims(ctx)

	var contact *models.Contact

	_, err := xid.FromString(contactData)
	if err != nil {

		contact, err = pb.contactBusiness.GetByDetail(ctx, contactData)
		if err != nil {
			return nil, err
		}
	} else {
		contact, err = pb.contactBusiness.GetByID(ctx, contactData)
		if err != nil {
			return nil, err
		}
	}

	return pb.GetByID(ctx, contact.ProfileID)
}

func (pb *profileBusiness) GetByID(
	ctx context.Context,
	profileID string) (*profilev1.ProfileObject, error) {

	ctx = frame.SkipTenancyChecksOnClaims(ctx)

	profile, err := pb.profileRepo.GetByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return pb.ProfileToAPI(ctx, profile)

}

func (pb *profileBusiness) SearchProfile(ctx context.Context,
	request *profilev1.SearchRequest, stream profilev1.ProfileService_SearchServer) error {

	ctx = frame.SkipTenancyChecksOnClaims(ctx)

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
		profileObject, err := pb.ProfileToAPI(ctx, profile)
		if err != nil {
			return err
		}
		err = stream.Send(&profilev1.SearchResponse{Data: []*profilev1.ProfileObject{profileObject}})
		if err != nil {
			return err
		}
	}

	return nil
}

func (pb *profileBusiness) MergeProfile(ctx context.Context,
	request *profilev1.MergeRequest) (*profilev1.ProfileObject, error) {

	ctx = frame.SkipTenancyChecksOnClaims(ctx)

	target, err := pb.profileRepo.GetByID(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	merging, err := pb.profileRepo.GetByID(ctx, request.GetMergeid())
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

	return pb.ProfileToAPI(ctx, target)
}

func (pb *profileBusiness) UpdateProfile(
	ctx context.Context,
	request *profilev1.UpdateRequest) (*profilev1.ProfileObject, error) {

	profile, err := pb.profileRepo.GetByID(ctx, request.GetId())
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

	return pb.ProfileToAPI(ctx, profile)
}

func (pb *profileBusiness) CreateProfile(
	ctx context.Context,
	request *profilev1.CreateRequest) (*profilev1.ProfileObject, error) {

	ctx = frame.SkipTenancyChecksOnClaims(ctx)

	contactDetail := strings.TrimSpace(request.GetContact())

	if contactDetail == "" {
		return nil, service.ErrorContactDetailsNotValid
	}

	p := models.Profile{}
	p.Properties = frame.DBPropertiesFromMap(request.GetProperties())

	contact, err := pb.contactBusiness.GetByDetail(ctx, contactDetail)
	if err != nil {
		if !errors.Is(err, service.ErrorContactDoesNotExist) {
			return nil, err
		}

		pt, err := pb.profileRepo.GetTypeByUID(ctx, request.GetType())
		if err != nil {
			return nil, err
		}

		p.ProfileType = *pt
		p.ProfileTypeID = pt.ID

		err = pb.profileRepo.Save(ctx, &p)
		if err != nil {
			return nil, err
		}

		err = pb.contactBusiness.CreateContact(ctx, pb.EncryptionKeyFunc(), p.GetID(), contactDetail)
		if err != nil {
			return nil, err
		}

		return pb.GetByID(ctx, p.GetID())
	}

	return pb.GetByID(ctx, contact.ProfileID)

}

//func (pb *profileBusiness) UpdateProperties(db *gorm.DB, params map[string]any) error {
//
//	storedPropertiesMap := make(map[string]any)
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

func (pb *profileBusiness) AddAddress(
	ctx context.Context,
	request *profilev1.AddAddressRequest) (*profilev1.ProfileObject, error) {

	address, err := pb.addressBusiness.CreateAddress(ctx, request.GetAddress())
	if err != nil {
		return nil, err
	}

	err = pb.addressBusiness.LinkAddressToProfile(ctx,
		request.GetId(), request.GetAddress().GetExtra(), address)
	if err != nil {
		return nil, err
	}

	return pb.GetByID(ctx, request.GetId())

}

func (pb *profileBusiness) AddContact(
	ctx context.Context,
	request *profilev1.AddContactRequest) (*profilev1.ProfileObject, error) {

	return pb.GetByID(ctx, request.GetId())

}
