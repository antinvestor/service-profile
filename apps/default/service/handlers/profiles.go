package handlers

import (
	"context"
	"math"
	"net"
	"strings"
	"time"

	"connectrpc.com/connect"
	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/apis/go/profile/v1/profilev1connect"
	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/datastore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// Constants for pagination and batch sizes.
const (
	// MaxBatchSize defines the maximum number of items to process in a single batch.
	MaxBatchSize = 50
)

type ProfileServer struct {
	Service              *frame.Service
	NotificationCli      *notificationv1.NotificationClient
	profileBusiness      business.ProfileBusiness
	rosterBusiness       business.RosterBusiness
	relationshipBusiness business.RelationshipBusiness

	profilev1connect.UnimplementedProfileServiceHandler
}

// NewProfileServer creates a new ProfileServer with the profile business already initialized.
func NewProfileServer(
	ctx context.Context,
	svc *frame.Service,
	notificationCli *notificationv1.NotificationClient,
) *ProfileServer {

	evtsMan := svc.EventsManager(ctx)
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	cfg := svc.Config().(*config.ProfileConfig)

	contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
	verificationRepo := repository.NewVerificationRepository(ctx, dbPool, workMan)

	contactBusiness := business.NewContactBusiness(ctx, cfg, evtsMan, contactRepo, verificationRepo)

	addressRepo := repository.NewAddressRepository(ctx, dbPool, workMan)
	addressBusiness := business.NewAddressBusiness(ctx, addressRepo)

	profileRepo := repository.NewProfileRepository(ctx, dbPool, workMan)
	profileBusiness := business.NewProfileBusiness(ctx, evtsMan, contactBusiness, addressBusiness, profileRepo)

	rosterRepo := repository.NewRosterRepository(ctx, dbPool, workMan)
	rosterBusiness := business.NewRosterBusiness(ctx, contactBusiness, rosterRepo)

	relationshipRepo := repository.NewRelationshipRepository(ctx, dbPool, workMan)
	relationshipBusiness := business.NewRelationshipBusiness(ctx, profileBusiness, relationshipRepo)

	return &ProfileServer{
		Service:              svc,
		NotificationCli:      notificationCli,
		profileBusiness:      profileBusiness,
		rosterBusiness:       rosterBusiness,
		relationshipBusiness: relationshipBusiness,
	}
}

func (ps *ProfileServer) toAPIError(err error) error {
	grpcError, ok := status.FromError(err)

	if ok {
		return grpcError.Err()
	}

	if data.ErrorIsNoRows(err) {
		return status.Error(codes.NotFound, err.Error())
	}

	return grpcError.Err()
}

//nolint:revive,staticcheck // server implementation
func (ps *ProfileServer) GetById(ctx context.Context,
	request *connect.Request[profilev1.GetByIdRequest]) (
	*connect.Response[profilev1.GetByIdResponse], error) {
	profileObj, err := ps.profileBusiness.GetByID(ctx, request.Msg.GetId())
	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return connect.NewResponse(&profilev1.GetByIdResponse{Data: profileObj}), nil
}

func (ps *ProfileServer) GetByContact(ctx context.Context,
	request *connect.Request[profilev1.GetByContactRequest]) (
	*connect.Response[profilev1.GetByContactResponse], error) {
	profileObj, err := ps.profileBusiness.GetByContact(ctx, request.Msg.GetContact())

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return connect.NewResponse(&profilev1.GetByContactResponse{Data: profileObj}), nil
}

func (ps *ProfileServer) Search(ctx context.Context,
	request *connect.Request[profilev1.SearchRequest],
	stream *connect.ServerStream[profilev1.SearchResponse],
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobResult, err := ps.profileBusiness.SearchProfile(ctx, request.Msg)
	if err != nil {
		return ps.toAPIError(err)
	}

	for {
		result, ok := jobResult.ReadResult(ctx)

		if !ok {
			return nil
		}

		if result.IsError() {
			return ps.toAPIError(result.Error())
		}

		for _, profile := range result.Item() {
			profileObject, err1 := ps.profileBusiness.ToAPI(ctx, profile)
			if err1 != nil {
				return err1
			}
			err1 = stream.Send(&profilev1.SearchResponse{Data: []*profilev1.ProfileObject{profileObject}})
			if err1 != nil {
				return err1
			}
		}
	}
}

func (ps *ProfileServer) Merge(ctx context.Context, request *connect.Request[profilev1.MergeRequest]) (*connect.Response[profilev1.MergeResponse], error) {
	profileObj, err := ps.profileBusiness.MergeProfile(ctx, request.Msg)
	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return connect.NewResponse(&profilev1.MergeResponse{Data: profileObj}), nil
}

func (ps *ProfileServer) Create(ctx context.Context, request *connect.Request[profilev1.CreateRequest]) (*connect.Response[profilev1.CreateResponse], error) {
	profileObj, err := ps.profileBusiness.CreateProfile(ctx, request.Msg)

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return connect.NewResponse(&profilev1.CreateResponse{Data: profileObj}), nil
}

func (ps *ProfileServer) Update(ctx context.Context, request *connect.Request[profilev1.UpdateRequest]) (*connect.Response[profilev1.UpdateResponse], error) {
	profileObj, err := ps.profileBusiness.UpdateProfile(ctx, request.Msg)

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return connect.NewResponse(&profilev1.UpdateResponse{Data: profileObj}), nil
}

// AddAddress Adds a new address based on the request.
func (ps *ProfileServer) AddAddress(ctx context.Context,
	request *connect.Request[profilev1.AddAddressRequest]) (*connect.Response[profilev1.AddAddressResponse], error) {
	profileObj, err := ps.profileBusiness.AddAddress(ctx, request.Msg)
	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return connect.NewResponse(&profilev1.AddAddressResponse{Data: profileObj}), nil
}

func (ps *ProfileServer) AddContact(ctx context.Context, request *connect.Request[profilev1.AddContactRequest]) (*connect.Response[profilev1.AddContactResponse], error) {
	profileObj, verificationID, err := ps.profileBusiness.AddContact(ctx, request.Msg)

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return connect.NewResponse(&profilev1.AddContactResponse{Data: profileObj, VerificationId: verificationID}), nil
}

func (ps *ProfileServer) CreateContact(
	ctx context.Context,
	request *connect.Request[profilev1.CreateContactRequest]) (*connect.Response[profilev1.CreateContactResponse], error) {
	contactObj, err := ps.profileBusiness.CreateContact(ctx, request.Msg)

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return connect.NewResponse(&profilev1.CreateContactResponse{Data: contactObj}), nil
}

func (ps *ProfileServer) CreateContactVerification(
	ctx context.Context,
	request *connect.Request[profilev1.CreateContactVerificationRequest]) (
	*connect.Response[profilev1.CreateContactVerificationResponse], error) {
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
		return nil, ps.toAPIError(err)
	}

	return connect.NewResponse(&profilev1.CreateContactVerificationResponse{
		Id:      verificationID,
		Success: true,
	}), nil
}

func getClientIP(ctx context.Context) string {
	// First, try to get the IP from the X-Forwarded-For header.
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
			// X-Forwarded-For can be a comma-separated list of IPs.
			// The first one is the original client.
			ips := strings.Split(xff[0], ",")
			if len(ips) > 0 {
				return strings.TrimSpace(ips[0])
			}
		}
		if xrip := md.Get("x-real-ip"); len(xrip) > 0 {
			return xrip[0]
		}
	}

	// If not available, fall back to the peer's address.
	p, ok := peer.FromContext(ctx)
	if ok {
		tcpAddr, tcok := p.Addr.(*net.TCPAddr)
		if tcok {
			return tcpAddr.IP.String()
		}
		return p.Addr.String()
	}

	return "unknown"
}

func (ps *ProfileServer) CheckVerification(
	ctx context.Context,
	request *connect.Request[profilev1.CheckVerificationRequest]) (*connect.Response[profilev1.CheckVerificationResponse], error) {
	verificationAttempts, verified, err := ps.profileBusiness.CheckVerification(
		ctx,
		request.Msg.GetId(),
		request.Msg.GetCode(),
		getClientIP(ctx),
	)

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return connect.NewResponse(
		&profilev1.CheckVerificationResponse{
			Id:            request.Msg.GetId(),
			CheckAttempts: int32(verificationAttempts), // #nosec G115 - verificationAttempts is bounded by business logic
			Success:       verified,
		}), nil
}

func (ps *ProfileServer) RemoveContact(
	ctx context.Context,
	request *connect.Request[profilev1.RemoveContactRequest]) (*connect.Response[profilev1.RemoveContactResponse], error) {
	profileObj, err := ps.profileBusiness.RemoveContact(ctx, request.Msg)

	if err != nil {
		return nil, ps.toAPIError(err)
	}
	return connect.NewResponse(&profilev1.RemoveContactResponse{Data: profileObj}), nil
}

func (ps *ProfileServer) SearchRoster(ctx context.Context,
	request *connect.Request[profilev1.SearchRosterRequest], stream *connect.ServerStream[profilev1.SearchRosterResponse],
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobResult, err := ps.rosterBusiness.Search(ctx, request.Msg)
	if err != nil {
		return ps.toAPIError(err)
	}

	for {
		result, ok := jobResult.ReadResult(ctx)

		if !ok {
			return nil
		}

		if result.IsError() {
			return ps.toAPIError(result.Error())
		}

		// Preallocate slice to optimize memory allocation.
		rosterList := make([]*profilev1.RosterObject, 0, len(result.Item()))
		for _, roster := range result.Item() {
			rosterList = append(rosterList, roster.ToAPI())
		}

		err = stream.Send(&profilev1.SearchRosterResponse{Data: rosterList})
		if err != nil {
			return ps.toAPIError(err)
		}
	}
}

func (ps *ProfileServer) AddRoster(
	ctx context.Context,
	request *connect.Request[profilev1.AddRosterRequest]) (*connect.Response[profilev1.AddRosterResponse], error) {

	roster, err := ps.rosterBusiness.CreateRoster(ctx, request.Msg)

	if err != nil {
		return nil, ps.toAPIError(err)
	}
	return connect.NewResponse(
		&profilev1.AddRosterResponse{
			Data: roster,
		}), nil
}

func (ps *ProfileServer) RemoveRoster(
	ctx context.Context,
	request *connect.Request[profilev1.RemoveRosterRequest]) (*connect.Response[profilev1.RemoveRosterResponse], error) {

	roster, err := ps.rosterBusiness.RemoveRoster(ctx, request.Msg.GetId())

	if err != nil {
		return nil, ps.toAPIError(err)
	}
	return connect.NewResponse(
		&profilev1.RemoveRosterResponse{
			Roster: roster,
		}), nil
}

func (ps *ProfileServer) AddRelationship(ctx context.Context,
	request *connect.Request[profilev1.AddRelationshipRequest]) (*connect.Response[profilev1.AddRelationshipResponse], error) {

	relationshipObj, err := ps.relationshipBusiness.CreateRelationship(ctx, request.Msg)

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return connect.NewResponse(
		&profilev1.AddRelationshipResponse{Data: relationshipObj}), nil
}

func (ps *ProfileServer) DeleteRelationship(ctx context.Context,
	request *connect.Request[profilev1.DeleteRelationshipRequest]) (*connect.Response[profilev1.DeleteRelationshipResponse], error) {

	relationshipObj, err := ps.relationshipBusiness.DeleteRelationship(ctx, request.Msg)

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return connect.NewResponse(
		&profilev1.DeleteRelationshipResponse{Data: relationshipObj}), nil
}

func (ps *ProfileServer) ListRelationships(ctx context.Context,
	request *connect.Request[profilev1.ListRelationshipRequest], stream *connect.ServerStream[profilev1.ListRelationshipResponse],
) error {

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
			return ps.toAPIError(err)
		}

		var responseList []*profilev1.RelationshipObject

		for _, relationship := range relationships {
			responseList = append(responseList, relationship.ToAPI())
		}

		err = stream.Send(&profilev1.ListRelationshipResponse{Data: responseList})
		if err != nil {
			return ps.toAPIError(err)
		}

		totalSent += len(relationships)

		if totalSent >= requiredCount || len(relationships) < remainingCount {
			break
		}
	}

	return nil
}
