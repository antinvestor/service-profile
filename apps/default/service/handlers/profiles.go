package handlers

import (
	"context"
	"math"
	"time"

	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	"buf.build/gen/go/antinvestor/profile/connectrpc/go/profile/v1/profilev1connect"
	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/security/authorizer"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/authz"
	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/internal/errorutil"
)

// Constants for pagination and batch sizes.
const (
	// MaxBatchSize defines the maximum number of items to process in a single batch.
	MaxBatchSize = 50
)

type ProfileServer struct {
	Service              *frame.Service
	DEK                  *config.DEK
	NotificationCli      notificationv1connect.NotificationServiceClient
	authz                authz.Middleware
	profileBusiness      business.ProfileBusiness
	contactBusiness      business.ContactBusiness
	rosterBusiness       business.RosterBusiness
	relationshipBusiness business.RelationshipBusiness

	profilev1connect.UnimplementedProfileServiceHandler
}

// NewProfileServer creates a new ProfileServer with the profile business already initialized.
func NewProfileServer(
	ctx context.Context,
	svc *frame.Service,
	dek *config.DEK,
	notificationCli notificationv1connect.NotificationServiceClient,
	authzMiddleware authz.Middleware,
) *ProfileServer {
	evtsMan := svc.EventsManager()
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	cfg, _ := svc.Config().(*config.ProfileConfig)

	contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
	verificationRepo := repository.NewVerificationRepository(ctx, dbPool, workMan)

	contactBusiness := business.NewContactBusiness(ctx, cfg, dek, evtsMan, contactRepo, verificationRepo)

	addressRepo := repository.NewAddressRepository(ctx, dbPool, workMan)
	addressBusiness := business.NewAddressBusiness(ctx, addressRepo)

	profileRepo := repository.NewProfileRepository(ctx, dbPool, workMan)
	profileBusiness := business.NewProfileBusiness(
		ctx,
		cfg,
		dek,
		evtsMan,
		contactBusiness,
		addressBusiness,
		profileRepo,
	)

	rosterRepo := repository.NewRosterRepository(ctx, dbPool, workMan)
	rosterBusiness := business.NewRosterBusiness(ctx, cfg, dek, contactBusiness, rosterRepo)

	relationshipRepo := repository.NewRelationshipRepository(ctx, dbPool, workMan)
	relationshipBusiness := business.NewRelationshipBusiness(ctx, profileBusiness, relationshipRepo)

	return &ProfileServer{
		Service:              svc,
		DEK:                  dek,
		NotificationCli:      notificationCli,
		authz:                authzMiddleware,
		profileBusiness:      profileBusiness,
		contactBusiness:      contactBusiness,
		rosterBusiness:       rosterBusiness,
		relationshipBusiness: relationshipBusiness,
	}
}

//nolint:revive,staticcheck // server implementation
func (ps *ProfileServer) GetById(ctx context.Context,
	request *connect.Request[profilev1.GetByIdRequest]) (
	*connect.Response[profilev1.GetByIdResponse], error) {
	if err := ps.authz.CanViewProfileSelf(ctx, request.Msg.GetId()); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	profileObj, err := ps.profileBusiness.GetByID(ctx, request.Msg.GetId())
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&profilev1.GetByIdResponse{Data: profileObj}), nil
}

func (ps *ProfileServer) GetByContact(ctx context.Context,
	request *connect.Request[profilev1.GetByContactRequest]) (
	*connect.Response[profilev1.GetByContactResponse], error) {
	if err := ps.authz.CanViewProfile(ctx); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	profileObj, err := ps.profileBusiness.GetByContact(ctx, request.Msg.GetContact())

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&profilev1.GetByContactResponse{Data: profileObj}), nil
}

func (ps *ProfileServer) Search(ctx context.Context,
	request *connect.Request[profilev1.SearchRequest],
	stream *connect.ServerStream[profilev1.SearchResponse],
) error {
	if err := ps.authz.CanViewProfile(ctx); err != nil {
		return authorizer.ToConnectError(err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobResult, err := ps.profileBusiness.SearchProfile(ctx, request.Msg)
	if err != nil {
		return errorutil.CleanErr(err)
	}

	for {
		result, ok := jobResult.ReadResult(ctx)

		if !ok {
			return nil
		}

		if result.IsError() {
			return errorutil.CleanErr(result.Error())
		}

		for _, profile := range result.Item() {
			profileObject, err1 := ps.profileBusiness.ToAPI(ctx, profile)
			if err1 != nil {
				return errorutil.CleanErr(err1)
			}
			sErr := stream.Send(&profilev1.SearchResponse{Data: []*profilev1.ProfileObject{profileObject}})
			if sErr != nil {
				return errorutil.CleanErr(sErr)
			}
		}
	}
}

func (ps *ProfileServer) Merge(
	ctx context.Context,
	request *connect.Request[profilev1.MergeRequest],
) (*connect.Response[profilev1.MergeResponse], error) {
	if err := ps.authz.CanMergeProfiles(ctx); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	profileObj, err := ps.profileBusiness.MergeProfile(ctx, request.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&profilev1.MergeResponse{Data: profileObj}), nil
}

func (ps *ProfileServer) Create(
	ctx context.Context,
	request *connect.Request[profilev1.CreateRequest],
) (*connect.Response[profilev1.CreateResponse], error) {
	if err := ps.authz.CanCreateProfile(ctx); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	profileObj, err := ps.profileBusiness.CreateProfile(ctx, request.Msg)

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&profilev1.CreateResponse{Data: profileObj}), nil
}

func (ps *ProfileServer) Update(
	ctx context.Context,
	request *connect.Request[profilev1.UpdateRequest],
) (*connect.Response[profilev1.UpdateResponse], error) {
	if err := ps.authz.CanUpdateProfileSelf(ctx, request.Msg.GetId()); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	profileObj, err := ps.profileBusiness.UpdateProfile(ctx, request.Msg)

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&profilev1.UpdateResponse{Data: profileObj}), nil
}

// AddAddress Adds a new address based on the request.
func (ps *ProfileServer) AddAddress(ctx context.Context,
	request *connect.Request[profilev1.AddAddressRequest]) (*connect.Response[profilev1.AddAddressResponse], error) {
	if err := ps.authz.CanManageContactsSelf(ctx, request.Msg.GetId()); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	profileObj, err := ps.profileBusiness.AddAddress(ctx, request.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&profilev1.AddAddressResponse{Data: profileObj}), nil
}

func (ps *ProfileServer) AddContact(
	ctx context.Context,
	request *connect.Request[profilev1.AddContactRequest],
) (*connect.Response[profilev1.AddContactResponse], error) {
	if err := ps.authz.CanManageContactsSelf(ctx, request.Msg.GetId()); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	profileObj, verificationID, err := ps.profileBusiness.AddContact(ctx, request.Msg)

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&profilev1.AddContactResponse{Data: profileObj, VerificationId: verificationID}), nil
}

func (ps *ProfileServer) CreateContact(
	ctx context.Context,
	request *connect.Request[profilev1.CreateContactRequest],
) (*connect.Response[profilev1.CreateContactResponse], error) {
	if err := ps.authz.CanManageContacts(ctx); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	createReq := request.Msg

	contactList, err := ps.contactBusiness.GetByDetail(ctx, createReq.GetContact())

	if err != nil {
		if !frame.ErrorIsNotFound(err) {
			return nil, errorutil.CleanErr(err)
		}
	}

	if len(contactList) > 0 {
		contact, decryptErr := contactList[0].ToAPI(ps.DEK, true)
		if decryptErr != nil {
			return nil, errorutil.CleanErr(decryptErr)
		}
		return connect.NewResponse(&profilev1.CreateContactResponse{Data: contact}), nil
	}

	requestProperties := data.JSONMap{}
	contact, err := ps.contactBusiness.CreateContact(
		ctx,
		createReq.GetContact(),
		requestProperties.FromProtoStruct(createReq.GetExtras()),
	)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	contactObj, err := contact.ToAPI(ps.DEK, true)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&profilev1.CreateContactResponse{Data: contactObj}), nil
}

func (ps *ProfileServer) CreateContactVerification(
	ctx context.Context,
	request *connect.Request[profilev1.CreateContactVerificationRequest]) (
	*connect.Response[profilev1.CreateContactVerificationResponse], error) {
	if err := ps.authz.CanManageContacts(ctx); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	expiryDuration, err := time.ParseDuration(request.Msg.GetDurationToExpire())
	if err != nil {
		expiryDuration = 0
	}

	verificationID, err := ps.profileBusiness.VerifyContact(
		ctx,
		request.Msg.GetContactId(),
		request.Msg.GetId(),
		request.Msg.GetCode(),
		expiryDuration,
	)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&profilev1.CreateContactVerificationResponse{
		Id:      verificationID,
		Success: true,
	}), nil
}

func (ps *ProfileServer) CheckVerification(
	ctx context.Context,
	request *connect.Request[profilev1.CheckVerificationRequest],
) (*connect.Response[profilev1.CheckVerificationResponse], error) {
	if err := ps.authz.CanManageContacts(ctx); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	verificationAttempts, verified, err := ps.profileBusiness.CheckVerification(
		ctx,
		request.Msg.GetId(),
		request.Msg.GetCode(),
		request.Peer().Addr,
	)

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	// #nosec G115: verificationAttempts is bounded by business logic constraints
	return connect.NewResponse(
		&profilev1.CheckVerificationResponse{
			Id: request.Msg.GetId(),
			CheckAttempts: int32(
				verificationAttempts,
			),
			Success: verified,
		}), nil
}

func (ps *ProfileServer) RemoveContact(
	ctx context.Context,
	request *connect.Request[profilev1.RemoveContactRequest],
) (*connect.Response[profilev1.RemoveContactResponse], error) {
	if err := ps.authz.CanManageContactsSelf(ctx, request.Msg.GetId()); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	profileObj, err := ps.profileBusiness.RemoveContact(ctx, request.Msg)

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return connect.NewResponse(&profilev1.RemoveContactResponse{Data: profileObj}), nil
}

func (ps *ProfileServer) SearchRoster(
	ctx context.Context,
	request *connect.Request[profilev1.SearchRosterRequest],
	stream *connect.ServerStream[profilev1.SearchRosterResponse],
) error {
	if err := ps.authz.CanManageRosterSelf(ctx, request.Msg.GetProfileId()); err != nil {
		return authorizer.ToConnectError(err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobResult, err := ps.rosterBusiness.Search(ctx, request.Msg)
	if err != nil {
		return errorutil.CleanErr(err)
	}

	for {
		result, ok := jobResult.ReadResult(ctx)

		if !ok {
			return nil
		}

		if result.IsError() {
			return errorutil.CleanErr(result.Error())
		}

		// Preallocate slice to optimize memory allocation.
		rosterList := make([]*profilev1.RosterObject, 0, len(result.Item()))
		for _, roster := range result.Item() {
			rosterObj, rosterErr := roster.ToAPI(ps.DEK)
			if rosterErr != nil {
				return errorutil.CleanErr(rosterErr)
			}

			rosterList = append(rosterList, rosterObj)
		}

		sErr := stream.Send(&profilev1.SearchRosterResponse{Data: rosterList})
		if sErr != nil {
			return errorutil.CleanErr(sErr)
		}
	}
}

func (ps *ProfileServer) AddRoster(
	ctx context.Context,
	request *connect.Request[profilev1.AddRosterRequest]) (*connect.Response[profilev1.AddRosterResponse], error) {
	if err := ps.authz.CanManageRoster(ctx); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	roster, err := ps.rosterBusiness.CreateRoster(ctx, request.Msg)

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return connect.NewResponse(
		&profilev1.AddRosterResponse{
			Data: roster,
		}), nil
}

func (ps *ProfileServer) RemoveRoster(
	ctx context.Context,
	request *connect.Request[profilev1.RemoveRosterRequest],
) (*connect.Response[profilev1.RemoveRosterResponse], error) {
	if err := ps.authz.CanManageRoster(ctx); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	roster, err := ps.rosterBusiness.RemoveRoster(ctx, request.Msg.GetId())

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return connect.NewResponse(
		&profilev1.RemoveRosterResponse{
			Roster: roster,
		}), nil
}

func (ps *ProfileServer) AddRelationship(
	ctx context.Context,
	request *connect.Request[profilev1.AddRelationshipRequest],
) (*connect.Response[profilev1.AddRelationshipResponse], error) {
	if err := ps.authz.CanManageRelationshipsSelf(ctx, request.Msg.GetId()); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	relationshipObj, err := ps.relationshipBusiness.CreateRelationship(ctx, request.Msg)

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(
		&profilev1.AddRelationshipResponse{Data: relationshipObj}), nil
}

func (ps *ProfileServer) DeleteRelationship(
	ctx context.Context,
	request *connect.Request[profilev1.DeleteRelationshipRequest],
) (*connect.Response[profilev1.DeleteRelationshipResponse], error) {
	if err := ps.authz.CanManageRelationships(ctx); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	relationshipObj, err := ps.relationshipBusiness.DeleteRelationship(ctx, request.Msg)

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(
		&profilev1.DeleteRelationshipResponse{Data: relationshipObj}), nil
}

func (ps *ProfileServer) ListRelationships(
	ctx context.Context,
	request *connect.Request[profilev1.ListRelationshipRequest],
	stream *connect.ServerStream[profilev1.ListRelationshipResponse],
) error {
	if err := ps.authz.CanManageRelationshipsSelf(ctx, request.Msg.GetPeerId()); err != nil {
		return authorizer.ToConnectError(err)
	}

	totalSent := 0
	requiredCount := int(request.Msg.GetCount())
	if requiredCount == 0 {
		requiredCount = 1000000
	}
	for {
		remainingCount := requiredCount - totalSent
		if remainingCount > MaxBatchSize {
			remainingCount = MaxBatchSize
		}
		// Apply count limits with bounds check
		if remainingCount > math.MaxInt32 {
			request.Msg.Count = math.MaxInt32
		} else {
			// Safe conversion with bounds check
			request.Msg.Count = int32(remainingCount) // #nosec G115 -- bounds checked above
		}

		relationships, err := ps.relationshipBusiness.ListRelationships(ctx, request.Msg)
		if err != nil {
			return errorutil.CleanErr(err)
		}

		var responseList []*profilev1.RelationshipObject

		for _, relationship := range relationships {
			responseList = append(responseList, relationship.ToAPI())
		}

		sErr := stream.Send(&profilev1.ListRelationshipResponse{Data: responseList})
		if sErr != nil {
			return errorutil.CleanErr(sErr)
		}

		totalSent += len(relationships)

		if totalSent >= requiredCount || len(relationships) < remainingCount {
			break
		}
	}

	return nil
}
