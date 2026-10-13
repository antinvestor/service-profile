package handlers

import (
	"context"
	"math"
	"time"

	"buf.build/gen/go/antinvestor/notification/connectrpc/go/v2/notification/v1/notificationv1connect"
	"buf.build/gen/go/antinvestor/profile/connectrpc/go/v2/profile/v1/profilev1connect"
	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect/v2"
	auditlib "github.com/antinvestor/common/audit"
	"github.com/pitabwire/frame/v3"
	"github.com/pitabwire/frame/v3/data"
	"github.com/pitabwire/frame/v3/datastore"
	"github.com/pitabwire/frame/v3/security"
	"github.com/pitabwire/frame/v3/security/authorizer"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/accounts"
	"github.com/antinvestor/service-profile/apps/default/service/authz"
	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/pkg/errorutil"
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
	checker              *authorizer.FunctionChecker
	profileBusiness      business.ProfileBusiness
	accountBusiness      business.AccountBusiness
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
	checker *authorizer.FunctionChecker,
	deriver *accounts.Deriver,
) *ProfileServer {
	evtsMan := svc.EventsManager()
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	cfg, _ := svc.Config().(*config.ProfileConfig)

	contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
	verificationRepo := repository.NewVerificationRepository(ctx, dbPool, workMan)

	contactBusiness := business.NewContactBusiness(
		ctx,
		cfg,
		dek,
		evtsMan,
		contactRepo,
		verificationRepo,
	)

	addressRepo := repository.NewAddressRepository(ctx, dbPool, workMan)
	addressBusiness := business.NewAddressBusiness(ctx, addressRepo)

	profileRepo := repository.NewProfileRepository(ctx, dbPool, workMan)
	propertyEntryRepo := repository.NewPropertyEntryRepository(ctx, dbPool, workMan)
	accountBusiness := business.NewAccountBusiness(
		deriver,
		repository.NewProfileAccountRepository(ctx, dbPool, workMan),
		profileRepo,
		cfg.AccountBackfillBatchSize,
		// Service principals are those holding account_resolve, which is
		// bound to ROLE_SERVICE alone.
		func(ctx context.Context) bool {
			return checker.Check(ctx, authz.PermissionAccountResolve) == nil
		},
	)
	profileBusiness := business.NewProfileBusiness(
		ctx,
		cfg,
		dek,
		evtsMan,
		contactBusiness,
		addressBusiness,
		profileRepo,
		propertyEntryRepo,
		accountBusiness,
	)

	rosterRepo := repository.NewRosterRepository(ctx, dbPool, workMan)
	rosterBusiness := business.NewRosterBusiness(ctx, cfg, dek, contactBusiness, rosterRepo)

	relationshipRepo := repository.NewRelationshipRepository(ctx, dbPool, workMan)
	relationshipBusiness := business.NewRelationshipBusiness(ctx, profileBusiness, relationshipRepo)

	return &ProfileServer{
		Service:              svc,
		DEK:                  dek,
		NotificationCli:      notificationCli,
		checker:              checker,
		profileBusiness:      profileBusiness,
		accountBusiness:      accountBusiness,
		contactBusiness:      contactBusiness,
		rosterBusiness:       rosterBusiness,
		relationshipBusiness: relationshipBusiness,
	}
}

//nolint:revive,staticcheck // server implementation
func (ps *ProfileServer) GetById(ctx context.Context,
	request *profilev1.GetByIdRequest) (*profilev1.GetByIdResponse, error) {
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != request.GetId() {
		if err := ps.checker.Check(ctx, "profile_view"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	profileObj, err := ps.profileBusiness.GetByID(ctx, request.GetId())
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &profilev1.GetByIdResponse{Data: profileObj}, nil
}

// ResolveAccounts maps chain account addresses to the profiles that own them.
// account_resolve is bound to ROLE_SERVICE only, so owners and admins are
// refused; it is checked here (and excluded from the interceptor map) so the
// rule is enforced on every path into the handler.
func (ps *ProfileServer) ResolveAccounts(ctx context.Context,
	request *profilev1.ResolveAccountsRequest) (*profilev1.ResolveAccountsResponse, error) {
	if err := ps.checker.Check(ctx, authz.PermissionAccountResolve); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	found, err := ps.accountBusiness.Resolve(ctx, request.GetAddresses())
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	owners := make([]*profilev1.AccountOwner, 0, len(found))
	for _, a := range found {
		owners = append(owners, &profilev1.AccountOwner{
			Address:   accounts.HexAddress(a.Address),
			ProfileId: a.ProfileID,
		})
	}
	return &profilev1.ResolveAccountsResponse{Data: owners}, nil
}

func (ps *ProfileServer) GetByContact(ctx context.Context,
	request *profilev1.GetByContactRequest) (*profilev1.GetByContactResponse, error) {
	profileObj, err := ps.profileBusiness.GetByContact(ctx, request.GetContact())

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &profilev1.GetByContactResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) Search(ctx context.Context,
	request *profilev1.SearchRequest,
	stream profilev1connect.ProfileServiceSearchServerStream,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobResult, err := ps.profileBusiness.SearchProfile(ctx, request)
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
			sErr := stream.Send(
				&profilev1.SearchResponse{Data: []*profilev1.ProfileObject{profileObject}},
			)
			if sErr != nil {
				return errorutil.CleanErr(sErr)
			}
		}
	}
}

func (ps *ProfileServer) Merge(
	ctx context.Context,
	request *profilev1.MergeRequest,
) (*profilev1.MergeResponse, error) {
	profileObj, err := ps.profileBusiness.MergeProfile(ctx, request)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &profilev1.MergeResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) Create(
	ctx context.Context,
	request *profilev1.CreateRequest,
) (*profilev1.CreateResponse, error) {
	profileObj, err := ps.profileBusiness.CreateProfile(ctx, request)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	auditlib.WithResource(ctx, auditlib.ResourceProfile, profileObj.GetId())
	auditlib.WithDetail(ctx, "profile_type", request.GetType().String())
	auditlib.WithDetail(ctx, "contact", request.GetContact())

	return &profilev1.CreateResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) Update(
	ctx context.Context,
	request *profilev1.UpdateRequest,
) (*profilev1.UpdateResponse, error) {
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != request.GetId() {
		if err := ps.checker.Check(ctx, "profile_update"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	profileObj, err := ps.profileBusiness.UpdateProfile(ctx, request)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	auditlib.WithResource(ctx, auditlib.ResourceProfile, request.GetId())

	return &profilev1.UpdateResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) GetByIDAndPartition(
	ctx context.Context,
	request *profilev1.GetByIDAndPartitionRequest,
) (*profilev1.GetByIDAndPartitionResponse, error) {
	profileObj, err := ps.profileBusiness.GetByIDAndPartition(
		ctx,
		request.GetId(),
		request.GetPartitionId(),
	)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	auditlib.WithResource(ctx, auditlib.ResourceProfile, request.GetId())

	return &profilev1.GetByIDAndPartitionResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) PropertyHistory(
	ctx context.Context,
	request *profilev1.PropertyHistoryRequest,
) (*profilev1.PropertyHistoryResponse, error) {
	claims := security.ClaimsFromContext(ctx)
	callerTenantID := ""
	if claims != nil {
		callerTenantID = claims.GetTenantID()
	}

	entries, err := ps.profileBusiness.GetPropertyHistory(
		ctx,
		request.GetId(),
		request.GetKey(),
		callerTenantID,
	)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	var protoEntries []*profilev1.PropertyEntryObject
	for _, e := range entries {
		protoEntries = append(protoEntries, &profilev1.PropertyEntryObject{
			Key:       e.Key,
			Value:     e.Value,
			TenantId:  e.TenantID,
			CreatedBy: e.CreatedBy,
			CreatedAt: timestamppb.New(e.CreatedAt),
			Scoped:    e.Scoped,
		})
	}

	return &profilev1.PropertyHistoryResponse{Entries: protoEntries}, nil
}

// AddAddress Adds a new address based on the request.
func (ps *ProfileServer) AddAddress(
	ctx context.Context,
	request *profilev1.AddAddressRequest,
) (*profilev1.AddAddressResponse, error) {
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != request.GetId() {
		if err := ps.checker.Check(ctx, "contact_manage"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	profileObj, err := ps.profileBusiness.AddAddress(ctx, request)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	auditlib.WithResource(ctx, auditlib.ResourceProfile, request.GetId())
	auditlib.WithRelation(ctx, auditlib.Relation{
		ParentType: auditlib.ResourceProfile,
		ParentID:   request.GetId(),
		ChildType:  auditlib.ResourceAddress,
		ChildID:    "", // address ID not returned in this response
		Action:     auditlib.RelationAdded,
	})

	return &profilev1.AddAddressResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) AddContact(
	ctx context.Context,
	request *profilev1.AddContactRequest,
) (*profilev1.AddContactResponse, error) {
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != request.GetId() {
		if err := ps.checker.Check(ctx, "contact_manage"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	profileObj, verificationID, err := ps.profileBusiness.AddContact(ctx, request)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	auditlib.WithResource(ctx, auditlib.ResourceProfile, request.GetId())
	auditlib.WithRelation(ctx, auditlib.Relation{
		ParentType: auditlib.ResourceProfile,
		ParentID:   request.GetId(),
		ChildType:  auditlib.ResourceContact,
		ChildID:    request.GetContact(),
		Action:     auditlib.RelationAdded,
	})

	return &profilev1.AddContactResponse{Data: profileObj, VerificationId: verificationID}, nil
}

func (ps *ProfileServer) CreateContact(
	ctx context.Context,
	request *profilev1.CreateContactRequest,
) (*profilev1.CreateContactResponse, error) {
	if err := ps.checker.Check(ctx, "contact_manage"); err != nil {
		return nil, authorizer.ToConnectError(err)
	}
	createReq := request

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
		return &profilev1.CreateContactResponse{Data: contact}, nil
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

	return &profilev1.CreateContactResponse{Data: contactObj}, nil
}

// GetContacts resolves one or many contact ids (standalone or attached).
func (ps *ProfileServer) GetContacts(
	ctx context.Context,
	request *profilev1.GetContactsRequest,
) (*profilev1.GetContactsResponse, error) {
	if err := ps.checker.Check(ctx, "contact_manage"); err != nil {
		return nil, authorizer.ToConnectError(err)
	}
	found, missing, err := ps.profileBusiness.GetContacts(ctx, request.GetIds())
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return &profilev1.GetContactsResponse{
		Data:       found,
		MissingIds: missing,
	}, nil
}

func (ps *ProfileServer) CreateContactVerification(
	ctx context.Context,
	request *profilev1.CreateContactVerificationRequest) (*profilev1.CreateContactVerificationResponse, error) {
	expiryDuration, err := time.ParseDuration(request.GetDurationToExpire())
	if err != nil {
		expiryDuration = 0
	}

	verificationID, err := ps.profileBusiness.VerifyContact(
		ctx,
		request.GetContactId(),
		request.GetId(),
		request.GetCode(),
		expiryDuration,
	)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &profilev1.CreateContactVerificationResponse{
		Id:      verificationID,
		Success: true,
	}, nil
}

func (ps *ProfileServer) CheckVerification(
	ctx context.Context,
	request *profilev1.CheckVerificationRequest,
) (*profilev1.CheckVerificationResponse, error) {
	verificationAttempts, verified, err := ps.profileBusiness.CheckVerification(
		ctx,
		request.GetId(),
		request.GetCode(),
		peerAddr(ctx),
	)

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	// #nosec G115: verificationAttempts is bounded by business logic constraints
	return &profilev1.CheckVerificationResponse{
		Id: request.GetId(),
		CheckAttempts: int32(
			verificationAttempts,
		),
		Success: verified,
	}, nil
}

func (ps *ProfileServer) RemoveContact(
	ctx context.Context,
	request *profilev1.RemoveContactRequest,
) (*profilev1.RemoveContactResponse, error) {
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != request.GetId() {
		if err := ps.checker.Check(ctx, "contact_manage"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	profileObj, err := ps.profileBusiness.RemoveContact(ctx, request)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	auditlib.WithResource(ctx, auditlib.ResourceContact, request.GetId())
	auditlib.WithAction(ctx, auditlib.ActionRemove)

	return &profilev1.RemoveContactResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) SearchRoster(
	ctx context.Context,
	request *profilev1.SearchRosterRequest,
	stream profilev1connect.ProfileServiceSearchRosterServerStream,
) error {
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != request.GetProfileId() {
		if err := ps.checker.Check(ctx, "roster_manage"); err != nil {
			return authorizer.ToConnectError(err)
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobResult, err := ps.rosterBusiness.Search(ctx, request)
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
	request *profilev1.AddRosterRequest,
) (*profilev1.AddRosterResponse, error) {
	roster, err := ps.rosterBusiness.CreateRoster(ctx, request)

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return &profilev1.AddRosterResponse{
		Data: roster,
	}, nil
}

func (ps *ProfileServer) RemoveRoster(
	ctx context.Context,
	request *profilev1.RemoveRosterRequest,
) (*profilev1.RemoveRosterResponse, error) {
	roster, err := ps.rosterBusiness.RemoveRoster(ctx, request.GetId())

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return &profilev1.RemoveRosterResponse{
		Roster: roster,
	}, nil
}

func (ps *ProfileServer) AddRelationship(
	ctx context.Context,
	request *profilev1.AddRelationshipRequest,
) (*profilev1.AddRelationshipResponse, error) {
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != request.GetId() {
		if err := ps.checker.Check(ctx, "relationship_manage"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	relationshipObj, err := ps.relationshipBusiness.CreateRelationship(ctx, request)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	auditlib.WithResource(ctx, auditlib.ResourceProfile, request.GetId())
	auditlib.WithRelation(ctx, auditlib.Relation{
		ParentType: auditlib.ResourceProfile,
		ParentID:   request.GetId(),
		ChildType:  auditlib.ResourceRelationship,
		ChildID:    relationshipObj.GetId(),
		Action:     auditlib.RelationAdded,
	})
	auditlib.WithDetail(ctx, "relationship_type", request.GetType().String())

	return &profilev1.AddRelationshipResponse{Data: relationshipObj}, nil
}

func (ps *ProfileServer) DeleteRelationship(
	ctx context.Context,
	request *profilev1.DeleteRelationshipRequest,
) (*profilev1.DeleteRelationshipResponse, error) {
	relationshipObj, err := ps.relationshipBusiness.DeleteRelationship(ctx, request)

	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &profilev1.DeleteRelationshipResponse{Data: relationshipObj}, nil
}

func (ps *ProfileServer) ListRelationship(
	ctx context.Context,
	request *profilev1.ListRelationshipRequest,
	stream profilev1connect.ProfileServiceListRelationshipServerStream,
) error {
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != request.GetPeerId() {
		if err := ps.checker.Check(ctx, "relationship_manage"); err != nil {
			return authorizer.ToConnectError(err)
		}
	}

	totalSent := 0
	requiredCount := int(request.GetCount())
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
			request.Count = math.MaxInt32
		} else {
			// Safe conversion with bounds check
			request.Count = int32(remainingCount) // #nosec G115 -- bounds checked above
		}

		relationships, err := ps.relationshipBusiness.ListRelationships(ctx, request)
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

// peerAddr returns the remote address of the current RPC, or "" when the
// context carries no server call info.
func peerAddr(ctx context.Context) string {
	if info, ok := connect.CallInfoForServerContext(ctx); ok {
		return info.PeerAddr
	}
	return ""
}
