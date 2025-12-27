package business

import (
	"context"
	"errors"
	"strings"
	"time"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/data"
	frevents "github.com/pitabwire/frame/events"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/workerpool"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

// ErrContactNotFound is returned when a contact lookup finds no matching contact.
var ErrContactNotFound = errors.New("contact not found")

type ProfileBusiness interface {
	GetByID(ctx context.Context, profileID string) (*profilev1.ProfileObject, error)
	GetByContact(ctx context.Context, detail string) (*profilev1.ProfileObject, error)

	SearchProfile(
		ctx context.Context,
		request *profilev1.SearchRequest,
	) (workerpool.JobResultPipe[[]*models.Profile], error)

	CreateProfile(ctx context.Context, request *profilev1.CreateRequest) (*profilev1.ProfileObject, error)

	UpdateProfile(ctx context.Context, request *profilev1.UpdateRequest) (*profilev1.ProfileObject, error)

	MergeProfile(ctx context.Context, request *profilev1.MergeRequest) (*profilev1.ProfileObject, error)

	AddAddress(ctx context.Context, address *profilev1.AddAddressRequest) (*profilev1.ProfileObject, error)

	AddContact(
		ctx context.Context,
		contact *profilev1.AddContactRequest,
	) (*profilev1.ProfileObject, string, error)

	RemoveContact(
		ctx context.Context,
		contact *profilev1.RemoveContactRequest,
	) (*profilev1.ProfileObject, error)

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

func NewProfileBusiness(_ context.Context, cfg *config.ProfileConfig, dek *config.DEK,
	eventsMan frevents.Manager,
	contactBusiness ContactBusiness, addressBusiness AddressBusiness,
	profileRepo repository.ProfileRepository) ProfileBusiness {
	return &profileBusiness{
		cfg:             cfg,
		dek:             dek,
		contactBusiness: contactBusiness,
		addressBusiness: addressBusiness,
		profileRepo:     profileRepo,
		eventsMan:       eventsMan,
	}
}

type profileBusiness struct {
	cfg             *config.ProfileConfig
	dek             *config.DEK
	contactBusiness ContactBusiness
	addressBusiness AddressBusiness

	profileRepo repository.ProfileRepository

	eventsMan frevents.Manager
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
		contactObj, toAPIErr := c.ToAPI(pb.dek, true)
		if toAPIErr != nil {
			return nil, toAPIErr
		}
		contactObjects = append(contactObjects, contactObj)
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
	var err error
	var contact *models.Contact

	_, contactTypeErr := ContactTypeFromDetail(ctx, contactData)
	if contactTypeErr == nil {
		contactList, detailErr := pb.contactBusiness.GetByDetail(ctx, contactData)
		if detailErr != nil {
			return nil, detailErr
		}

		if len(contactList) == 0 {
			return nil, errors.New("contact not found")
		}
		contact = contactList[0]
	} else {
		contact, err = pb.contactBusiness.GetByID(ctx, contactData)
		if err != nil {
			return nil, err
		}
	}

	if contact.ProfileID == "" {
		profileObject := profilev1.ProfileObject{}

		profileObject.Type = models.ProfileTypeIDToEnum(models.ProfileTypePersonID)
		props := data.JSONMap{}
		profileObject.Properties = props.ToProtoStruct()
		contactObj, toAPIErr := contact.ToAPI(pb.dek, true)
		if toAPIErr != nil {
			return nil, toAPIErr
		}
		profileObject.Contacts = []*profilev1.ContactObject{contactObj}
		profileObject.Addresses = []*profilev1.AddressObject{}

		return &profileObject, nil
	}

	return pb.GetByID(ctx, contact.ProfileID)
}

func (pb *profileBusiness) GetByID(
	ctx context.Context,
	profileID string) (*profilev1.ProfileObject, error) {
	profile, err := pb.profileRepo.GetByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return pb.ToAPI(ctx, profile)
}

func (pb *profileBusiness) SearchProfile(ctx context.Context,
	request *profilev1.SearchRequest) (workerpool.JobResultPipe[[]*models.Profile], error) {
	profileID := ""
	claims := security.ClaimsFromContext(ctx)
	if claims != nil {
		profileID, _ = claims.GetSubject()
	}

	query := data.NewSearchQuery(
		data.WithSearchLimit(int(request.GetCount())),
		data.WithSearchOffset(int(request.GetPage())),
		data.WithSearchFiltersAndByValue(map[string]any{"profile_id": profileID}),
		data.WithSearchFiltersOrByValue(map[string]any{
			"searchable @@ websearch_to_tsquery( 'english', ?) ": request.GetQuery(),
		}))

	result, err := pb.profileRepo.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (pb *profileBusiness) MergeProfile(ctx context.Context,
	request *profilev1.MergeRequest) (*profilev1.ProfileObject, error) {
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

	_, err = pb.profileRepo.Update(ctx, target, "properties")
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

	requestProperties := data.JSONMap{}
	requestProperties = requestProperties.FromProtoStruct(request.GetProperties())
	profile.Properties = profile.Properties.Update(requestProperties)

	_, err = pb.profileRepo.Update(ctx, profile, "properties")
	if err != nil {
		return nil, err
	}

	return pb.ToAPI(ctx, profile)
}

// lookupContactByDetail attempts to find a contact by detail or ID.
// Returns the contact if found, nil if not found, or an error if lookup fails.
func (pb *profileBusiness) lookupContactByDetail(
	ctx context.Context,
	contactDetail string,
) (*models.Contact, error) {
	_, contactTypeErr := ContactTypeFromDetail(ctx, contactDetail)
	if contactTypeErr != nil {
		// Not a valid contact format, try lookup by ID
		return pb.contactBusiness.GetByID(ctx, contactDetail)
	}

	// Valid contact format, lookup by detail
	contactList, detailErr := pb.contactBusiness.GetByDetail(ctx, contactDetail)
	if detailErr != nil && !frame.ErrorIsNotFound(detailErr) {
		return nil, detailErr
	}

	if len(contactList) != 0 {
		return contactList[0], nil
	}

	return nil, ErrContactNotFound
}

func (pb *profileBusiness) CreateProfile(
	ctx context.Context,
	request *profilev1.CreateRequest) (*profilev1.ProfileObject, error) {
	contactDetail := strings.TrimSpace(request.GetContact())

	if contactDetail == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("contact details are invalid"))
	}

	p := models.Profile{}

	p.Properties = request.GetProperties().AsMap()

	contact, lookupErr := pb.lookupContactByDetail(ctx, contactDetail)
	if lookupErr != nil && !errors.Is(lookupErr, ErrContactNotFound) {
		return nil, lookupErr
	}

	if contact != nil && contact.ProfileID != "" {
		return pb.GetByID(ctx, contact.ProfileID)
	}

	var pt *models.ProfileType
	pt, repoErr := pb.profileRepo.GetTypeByUID(ctx, request.GetType())
	if repoErr != nil {
		return nil, data.ErrorConvertToAPI(repoErr)
	}

	p.ProfileType = *pt
	p.ProfileTypeID = pt.ID

	createErr := pb.profileRepo.Create(ctx, &p)
	if createErr != nil {
		return nil, data.ErrorConvertToAPI(createErr)
	}

	var err error
	if contact == nil {
		contact, err = pb.contactBusiness.CreateContact(ctx, contactDetail, data.JSONMap{})
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

// func (pb *profileBusiness) UpdateProperties(db *gorm.DB, params data.JSONMap) error {
//
//	storedPropertiesMap := make(data.JSONMap)
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
	claims := security.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, "", connect.NewError(connect.CodeInvalidArgument, errors.New("contact profile is invalid"))
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return nil, "", connect.NewError(connect.CodeInvalidArgument, errors.New("contact profile is invalid"))
	}
	if request.GetId() != subject {
		return nil, "", connect.NewError(connect.CodeInvalidArgument, errors.New("contact profile is invalid"))
	}

	profile, profileErr := pb.GetByID(ctx, subject)
	if profileErr != nil {
		return nil, "", profileErr
	}

	for _, contact := range profile.GetContacts() {
		normalizedContact := request.GetContact()
		if contact.GetDetail() == normalizedContact {
			return profile, "", nil
		}
	}

	extrasMap := data.JSONMap{}

	resp, contactErr := pb.contactBusiness.CreateContact(
		ctx,
		request.GetContact(),
		extrasMap.FromProtoStruct(request.GetExtras()),
	)

	if contactErr != nil {
		return nil, "", contactErr
	}

	verificationID, verifyErr := pb.VerifyContact(ctx, resp.GetID(), "", "", 0)
	if verifyErr != nil {
		return nil, "", verifyErr
	}

	return profile, verificationID, nil
}

func (pb *profileBusiness) RemoveContact(
	ctx context.Context,
	request *profilev1.RemoveContactRequest) (*profilev1.ProfileObject, error) {
	return pb.GetByID(ctx, request.GetId())
}

func (pb *profileBusiness) GetContactByID(
	ctx context.Context,
	contactID string,
) (*profilev1.ContactObject, error) {
	contact, err := pb.contactBusiness.GetByID(ctx, contactID)
	if err != nil {
		return nil, err
	}

	contactObj, err := contact.ToAPI(pb.dek, true)
	if err != nil {
		return nil, err
	}
	return contactObj, nil
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
	logger := util.Log(ctx).WithField("verificationID", verificationID)

	verification, verifyErr := pb.contactBusiness.GetVerification(ctx, verificationID)
	if verifyErr != nil {
		return 0, false, verifyErr
	}

	attempts, attemptsErr := pb.contactBusiness.GetVerificationAttempts(ctx, verificationID)
	if attemptsErr != nil {
		return 0, false, attemptsErr
	}

	verificationAttempts := len(attempts) + 1

	deviceID := ""
	claim := security.ClaimsFromContext(ctx)
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

	eventErr := pb.eventsMan.Emit(ctx, events.VerificationAttemptEventHandlerName, newAttempt)
	if eventErr != nil {
		logger.WithError(eventErr).Error("could not emit verification attempt event")
	}

	verificationExpired := verification.ExpiresAt.Before(time.Now())

	return verificationAttempts, codeMatches && !verificationExpired, nil
}
