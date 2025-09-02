package business

import (
	"context"
	"errors"
	"strings"
	"time"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/framedata"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/default/service"
	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

type ProfileBusiness interface {
	GetByID(ctx context.Context, profileID string) (*profilev1.ProfileObject, error)
	GetByContact(ctx context.Context, detail string) (*profilev1.ProfileObject, error)

	SearchProfile(ctx context.Context, request *profilev1.SearchRequest) (frame.JobResultPipe[[]*models.Profile], error)

	CreateProfile(ctx context.Context, request *profilev1.CreateRequest) (*profilev1.ProfileObject, error)

	UpdateProfile(ctx context.Context, request *profilev1.UpdateRequest) (*profilev1.ProfileObject, error)

	MergeProfile(ctx context.Context, request *profilev1.MergeRequest) (*profilev1.ProfileObject, error)

	AddAddress(ctx context.Context, address *profilev1.AddAddressRequest) (*profilev1.ProfileObject, error)

	AddContact(ctx context.Context, contact *profilev1.AddContactRequest) (*profilev1.ProfileObject, string, error)

	CreateContact(ctx context.Context, contact *profilev1.CreateContactRequest) (*profilev1.ContactObject, error)

	RemoveContact(ctx context.Context, contact *profilev1.RemoveContactRequest) (*profilev1.ProfileObject, error)

	GetContactByID(ctx context.Context, contactID string) (*profilev1.ContactObject, error)
	VerifyContact(
		ctx context.Context,
		contactID string,
		verificationID, code string,
		expiryDuration time.Duration,
	) (string, error)
	CheckVerification(ctx context.Context, verificationID string, code, ipAddress string) (int, bool, error)
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
	profileObject.Properties = p.Properties.ToProtoStruct()

	var contactObjects []*profilev1.ContactObject
	contactList, err := pb.contactBusiness.GetByProfile(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	for _, c := range contactList {
		contactObjects = append(contactObjects, c.ToAPI(true))
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

	_, err := ContactTypeFromDetail(ctx, contactData)
	if err == nil {
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
	profileID := ""
	claims := frame.ClaimsFromContext(ctx)
	if claims != nil {
		profileID, _ = claims.GetSubject()
	}

	searchProperties := frame.JSONMap{
		"profile_id": profileID,
		"start_date": request.GetStartDate(),
		"end_date":   request.GetEndDate(),
	}

	for _, p := range request.GetProperties() {
		searchProperties[p] = request.GetQuery()
	}

	query := framedata.NewSearchQuery(
		request.GetQuery(), searchProperties,
		int(request.GetPage()),
		int(request.GetCount()),
	)

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

	requestProperties := frame.JSONMap{}
	requestProperties = requestProperties.FromProtoStruct(request.GetProperties())
	profile.Properties = profile.Properties.Update(requestProperties)

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

	p.Properties = request.GetProperties().AsMap()

	var contact *models.Contact
	_, err := ContactTypeFromDetail(ctx, contactDetail)
	if err == nil {
		contact, err = pb.contactBusiness.GetByDetail(ctx, contactDetail)
		if !errors.Is(err, service.ErrContactDoesNotExist) {
			return nil, err
		}
	} else {
		contact, err = pb.contactBusiness.GetByID(ctx, contactDetail)
		if err != nil {
			return nil, err
		}
	}

	if contact != nil && contact.ProfileID != "" {
		return pb.GetByID(ctx, contact.ProfileID)
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

	if contact == nil {
		contact, err = pb.contactBusiness.CreateContact(ctx, contactDetail, frame.JSONMap{})
		if err != nil {
			return nil, err
		}
	}

	contact, err = pb.contactBusiness.UpdateContact(ctx, contact.GetID(), p.GetID(), nil)
	if err != nil {
		return nil, err
	}

	return pb.GetByID(ctx, contact.ProfileID)
}

// func (pb *profileBusiness) UpdateProperties(db *gorm.DB, params frame.JSONMap) error {
//
//	storedPropertiesMap := make(frame.JSONMap)
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
	request *profilev1.AddContactRequest) (*profilev1.ProfileObject, string, error) {
	claims := frame.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, "", service.ErrContactProfileNotValid
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return nil, "", service.ErrContactProfileNotValid
	}
	if request.GetId() != subject {
		return nil, "", service.ErrContactProfileNotValid
	}

	profile, err := pb.GetByID(ctx, subject)
	if err != nil {
		return nil, "", err
	}

	for _, contact := range profile.GetContacts() {
		if contact.GetDetail() == request.GetContact() {
			return profile, "", nil
		}
	}

	resp, err := pb.CreateContact(ctx, &profilev1.CreateContactRequest{
		Contact: request.GetContact(),
		Extras:  request.GetExtras(),
	})

	if err != nil {
		return nil, "", err
	}

	verificationID, err := pb.VerifyContact(ctx, resp.GetId(), "", "", 0)
	if err != nil {
		return nil, "", err
	}

	return profile, verificationID, nil
}

func (pb *profileBusiness) CreateContact(
	ctx context.Context,
	request *profilev1.CreateContactRequest) (*profilev1.ContactObject, error) {
	requestProperties := frame.JSONMap{}
	contact, err := pb.contactBusiness.CreateContact(
		ctx,
		request.GetContact(),
		requestProperties.FromProtoStruct(request.GetExtras()),
	)
	if err != nil {
		return nil, err
	}

	return contact.ToAPI(true), nil
}

func (pb *profileBusiness) RemoveContact(
	ctx context.Context,
	request *profilev1.RemoveContactRequest) (*profilev1.ProfileObject, error) {
	return pb.GetByID(ctx, request.GetId())
}

func (pb *profileBusiness) GetContactByID(ctx context.Context, contactID string) (*profilev1.ContactObject, error) {
	contact, err := pb.contactBusiness.GetByID(ctx, contactID)
	if err != nil {
		return nil, err
	}

	return contact.ToAPI(true), nil
}

func (pb *profileBusiness) VerifyContact(
	ctx context.Context,
	contactID, verificationID, code string,
	expiryDuration time.Duration,
) (string, error) {
	contact, err := pb.contactBusiness.GetByID(ctx, contactID)
	if err != nil {
		return "", err
	}

	verification, err := pb.contactBusiness.VerifyContact(ctx, contact, verificationID, code, expiryDuration)
	if err != nil {
		return "", err
	}

	return verification.GetID(), nil
}

func (pb *profileBusiness) CheckVerification(
	ctx context.Context,
	verificationID string,
	code string,
	ipAddress string,
) (int, bool, error) {
	logger := pb.service.Log(ctx).WithField("verificationID", verificationID)

	verification, err := pb.contactBusiness.GetVerification(ctx, verificationID)
	if err != nil {
		return 0, false, err
	}

	attempts, err := pb.contactBusiness.GetVerificationAttempts(ctx, verificationID)
	if err != nil {
		return 0, false, err
	}

	verificationAttempts := len(attempts) + 1

	deviceID := ""
	claim := frame.ClaimsFromContext(ctx)
	if claim != nil {
		deviceID = claim.DeviceID
	}

	newAttempt := &models.VerificationAttempt{
		VerificationID: verificationID,
		Data:           code,
		State:          "Fail",
		DeviceID:       deviceID,
		IPAddress:      ipAddress,
		RequestID:      util.GetRequestID(ctx),
	}

	newAttempt.GenID(ctx)

	codeMatches := false
	if verification.Code == code {
		newAttempt.State = "Success"
		codeMatches = true
	}

	err = pb.service.Emit(ctx, events.VerificationAttemptEventHandlerName, newAttempt)
	if err != nil {
		logger.WithError(err).Error("could not emit verification attempt event")
	}

	verificationExpired := verification.ExpiresAt.Before(time.Now())

	return verificationAttempts, codeMatches && !verificationExpired, nil
}
