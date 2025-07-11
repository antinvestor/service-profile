package business

import (
	"context"
	"errors"
	"strings"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame"
	"github.com/rs/xid"

	"github.com/antinvestor/service-profile/apps/default/service"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/internal/dbutil"
)

type ProfileBusiness interface {
	GetByID(ctx context.Context, profileID string) (*profilev1.ProfileObject, error)
	GetByContact(ctx context.Context, detail string) (*profilev1.ProfileObject, error)

	SearchProfile(ctx context.Context, request *profilev1.SearchRequest) (frame.JobResultPipe[[]*models.Profile], error)

	CreateProfile(ctx context.Context, request *profilev1.CreateRequest) (*profilev1.ProfileObject, error)

	UpdateProfile(ctx context.Context, request *profilev1.UpdateRequest) (*profilev1.ProfileObject, error)

	MergeProfile(ctx context.Context, request *profilev1.MergeRequest) (*profilev1.ProfileObject, error)

	AddAddress(ctx context.Context, address *profilev1.AddAddressRequest) (*profilev1.ProfileObject, error)

	AddContact(ctx context.Context, contact *profilev1.AddContactRequest) (*profilev1.ProfileObject, error)

	RemoveContact(ctx context.Context, contact *profilev1.RemoveContactRequest) (*profilev1.ProfileObject, error)
	ToAPI(ctx context.Context, profile *models.Profile) (*profilev1.ProfileObject, error)
}

func NewProfileBusiness(ctx context.Context, service *frame.Service) ProfileBusiness {
	return &profileBusiness{
		service:         service,
		contactBusiness: NewContactBusiness(ctx, service),
		addressBusiness: NewAddressBusiness(ctx, service),
		profileRepo:     repository.NewProfileRepository(service),
	}
}

type profileBusiness struct {
	service *frame.Service

	contactBusiness ContactBusiness
	addressBusiness AddressBusiness

	profileRepo repository.ProfileRepository
}

func (pb *profileBusiness) ToAPI(ctx context.Context,
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
		ctObj, ctErr := pb.contactBusiness.ToAPI(ctx, c, true)
		if ctErr != nil {
			return nil, ctErr
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

	return pb.ToAPI(ctx, profile)
}

func (pb *profileBusiness) SearchProfile(ctx context.Context,
	request *profilev1.SearchRequest) (frame.JobResultPipe[[]*models.Profile], error) {
	ctx = frame.SkipTenancyChecksOnClaims(ctx)

	profileID := ""
	claims := frame.ClaimsFromContext(ctx)
	if claims != nil {
		profileID, _ = claims.GetSubject()
	}

	query, err := dbutil.NewSearchQuery(
		ctx,
		profileID,
		request.GetQuery(),
		request.GetProperties(),
		request.GetStartDate(),
		request.GetEndDate(),
		int(request.GetCount()),
		int(request.GetPage()),
	)
	if err != nil {
		return nil, err
	}

	return pb.profileRepo.Search(ctx, query)
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

	return pb.ToAPI(ctx, target)
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

	return pb.ToAPI(ctx, profile)
}

func (pb *profileBusiness) CreateProfile(
	ctx context.Context,
	request *profilev1.CreateRequest) (*profilev1.ProfileObject, error) {
	ctx = frame.SkipTenancyChecksOnClaims(ctx)

	contactDetail := strings.TrimSpace(request.GetContact())

	if contactDetail == "" {
		return nil, service.ErrContactDetailsNotValid
	}

	p := models.Profile{}
	p.Properties = frame.DBPropertiesFromMap(request.GetProperties())

	contact, err := pb.contactBusiness.GetByDetail(ctx, contactDetail)
	if err == nil {
		return pb.GetByID(ctx, contact.ProfileID)
	}

	if !errors.Is(err, service.ErrContactDoesNotExist) {
		return nil, err
	}

	var pt *models.ProfileType
	pt, err = pb.profileRepo.GetTypeByUID(ctx, request.GetType())
	if err != nil {
		return nil, err
	}

	p.ProfileType = *pt
	p.ProfileTypeID = pt.ID

	err = pb.profileRepo.Save(ctx, &p)
	if err != nil {
		return nil, err
	}

	contact, err = pb.contactBusiness.CreateContact(ctx, contactDetail, map[string]string{})
	if err != nil {
		return nil, err
	}

	contact, err = pb.contactBusiness.UpdateContact(ctx, contact.GetID(), p.GetID(), map[string]string{})
	if err != nil {
		return nil, err
	}

	return pb.GetByID(ctx, contact.ProfileID)
}

// func (pb *profileBusiness) UpdateProperties(db *gorm.DB, params map[string]any) error {
//
//	storedPropertiesMap := make(map[string]any)
//	attributeMap, err := p.PropertiesToSearchOn.MarshalJSON()
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
//	err = p.PropertiesToSearchOn.UnmarshalJSON(stringProperties)
//	if err != nil {
//		return err
//	}
//
//	return db.Model(p).Update("PropertiesToSearchOn", p.PropertiesToSearchOn).Error
// }

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

func (pb *profileBusiness) RemoveContact(
	ctx context.Context,
	request *profilev1.RemoveContactRequest) (*profilev1.ProfileObject, error) {
	return pb.GetByID(ctx, request.GetId())
}
