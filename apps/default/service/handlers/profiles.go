package handlers

import (
	"context"
	"math"
	"net"
	"strings"
	"time"

	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/antinvestor/service-profile/apps/default/service/business"
)

// Constants for pagination and batch sizes.
const (
	// MaxBatchSize defines the maximum number of items to process in a single batch.
	MaxBatchSize = 50
)

type ProfileServer struct {
	Service         *frame.Service
	NotificationCli *notificationv1.NotificationClient
	ProfileBusiness business.ProfileBusiness

	profilev1.UnimplementedProfileServiceServer
}

// NewProfileServer creates a new ProfileServer with the profile business already initialized.
func NewProfileServer(
	ctx context.Context,
	service *frame.Service,
	notificationCli *notificationv1.NotificationClient,
) *ProfileServer {
	return &ProfileServer{
		Service:         service,
		NotificationCli: notificationCli,
		ProfileBusiness: business.NewProfileBusiness(ctx, service),
	}
}

func (ps *ProfileServer) toAPIError(err error) error {
	grpcError, ok := status.FromError(err)

	if ok {
		return grpcError.Err()
	}

	if frame.ErrorIsNoRows(err) {
		return status.Error(codes.NotFound, err.Error())
	}

	return grpcError.Err()
}

func (ps *ProfileServer) GetByID(ctx context.Context,
	request *profilev1.GetByIdRequest) (*profilev1.GetByIdResponse, error) {
	profileObj, err := ps.ProfileBusiness.GetByID(ctx, request.GetId())
	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return &profilev1.GetByIdResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) GetByContact(ctx context.Context,
	request *profilev1.GetByContactRequest) (*profilev1.GetByContactResponse, error) {
	profileObj, err := ps.ProfileBusiness.GetByContact(ctx, request.GetContact())

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return &profilev1.GetByContactResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) Search(request *profilev1.SearchRequest, stream profilev1.ProfileService_SearchServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	jobResult, err := ps.ProfileBusiness.SearchProfile(ctx, request)
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
			profileObject, err1 := ps.ProfileBusiness.ToAPI(ctx, profile)
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

func (ps *ProfileServer) Merge(ctx context.Context, request *profilev1.MergeRequest) (
	*profilev1.MergeResponse, error) {
	profileObj, err := ps.ProfileBusiness.MergeProfile(ctx, request)
	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return &profilev1.MergeResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) Create(ctx context.Context, request *profilev1.CreateRequest) (
	*profilev1.CreateResponse, error) {
	profileObj, err := ps.ProfileBusiness.CreateProfile(ctx, request)

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return &profilev1.CreateResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) Update(ctx context.Context, request *profilev1.UpdateRequest) (
	*profilev1.UpdateResponse, error) {
	profileObj, err := ps.ProfileBusiness.UpdateProfile(ctx, request)

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return &profilev1.UpdateResponse{Data: profileObj}, nil
}

// AddAddress Adds a new address based on the request.
func (ps *ProfileServer) AddAddress(ctx context.Context,
	request *profilev1.AddAddressRequest) (*profilev1.AddAddressResponse, error) {
	profileObj, err := ps.ProfileBusiness.AddAddress(ctx, request)
	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return &profilev1.AddAddressResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) AddContact(ctx context.Context, request *profilev1.AddContactRequest,
) (*profilev1.AddContactResponse, error) {
	profileObj, err := ps.ProfileBusiness.AddContact(ctx, request)

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return &profilev1.AddContactResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) CreateContactVerification(
	ctx context.Context,
	request *profilev1.CreateContactVerificationRequest,
) (*profilev1.CreateContactVerificationResponse, error) {
	expiryDuration, err := time.ParseDuration(request.GetDurationToExpire())
	if err != nil {
		expiryDuration = 0
	}

	verificationID, err := ps.ProfileBusiness.VerifyContact(
		ctx,
		request.GetContactId(),
		request.GetId(),
		request.GetCode(),
		expiryDuration,
	)
	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return &profilev1.CreateContactVerificationResponse{
		Id:      verificationID,
		Success: true,
	}, nil
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
	request *profilev1.CheckVerificationRequest,
) (*profilev1.CheckVerificationResponse, error) {
	verificationAttempts, verified, err := ps.ProfileBusiness.CheckVerification(
		ctx,
		request.GetId(),
		request.GetCode(),
		getClientIP(ctx),
	)

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return &profilev1.CheckVerificationResponse{
		Id:            request.GetId(),
		CheckAttempts: int32(verificationAttempts),
		Success:       verified,
	}, nil
}

func (ps *ProfileServer) RemoveContact(
	ctx context.Context,
	request *profilev1.RemoveContactRequest,
) (*profilev1.RemoveContactResponse, error) {
	profileObj, err := ps.ProfileBusiness.RemoveContact(ctx, request)

	if err != nil {
		return nil, ps.toAPIError(err)
	}
	return &profilev1.RemoveContactResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) SearchRoster(
	request *profilev1.SearchRosterRequest,
	stream grpc.ServerStreamingServer[profilev1.SearchRosterResponse],
) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	rosterBusiness := business.NewRosterBusiness(ctx, ps.Service)
	jobResult, err := rosterBusiness.Search(ctx, request)
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
			rosterObject, err1 := rosterBusiness.ToAPI(ctx, roster)
			if err1 != nil {
				return ps.toAPIError(err1)
			}

			rosterList = append(rosterList, rosterObject)
		}

		err = stream.Send(&profilev1.SearchRosterResponse{Data: rosterList})
		if err != nil {
			return ps.toAPIError(err)
		}
	}
}

func (ps *ProfileServer) AddRoster(
	ctx context.Context,
	request *profilev1.AddRosterRequest,
) (*profilev1.AddRosterResponse, error) {
	rosterBusiness := business.NewRosterBusiness(ctx, ps.Service)
	roster, err := rosterBusiness.CreateRoster(ctx, request)

	if err != nil {
		return nil, ps.toAPIError(err)
	}
	return &profilev1.AddRosterResponse{
		Data: roster,
	}, nil
}

func (ps *ProfileServer) RemoveRoster(
	ctx context.Context,
	request *profilev1.RemoveRosterRequest,
) (*profilev1.RemoveRosterResponse, error) {
	rosterBusiness := business.NewRosterBusiness(ctx, ps.Service)
	roster, err := rosterBusiness.RemoveRoster(ctx, request.GetId())

	if err != nil {
		return nil, ps.toAPIError(err)
	}
	return &profilev1.RemoveRosterResponse{
		Roster: roster,
	}, nil
}

func (ps *ProfileServer) AddRelationship(ctx context.Context,
	request *profilev1.AddRelationshipRequest) (*profilev1.AddRelationshipResponse, error) {
	relationshipBusiness := business.NewRelationshipBusiness(ctx, ps.Service, ps.ProfileBusiness)
	relationshipObj, err := relationshipBusiness.CreateRelationship(ctx, request)

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return &profilev1.AddRelationshipResponse{Data: relationshipObj}, nil
}

func (ps *ProfileServer) DeleteRelationship(ctx context.Context,
	request *profilev1.DeleteRelationshipRequest) (*profilev1.DeleteRelationshipResponse, error) {
	relationshipBusiness := business.NewRelationshipBusiness(ctx, ps.Service, ps.ProfileBusiness)
	relationshipObj, err := relationshipBusiness.DeleteRelationship(ctx, request)

	if err != nil {
		return nil, ps.toAPIError(err)
	}

	return &profilev1.DeleteRelationshipResponse{Data: relationshipObj}, nil
}

func (ps *ProfileServer) ListRelationships(
	request *profilev1.ListRelationshipRequest,
	server profilev1.ProfileService_ListRelationshipServer,
) error {
	ctx := server.Context()

	relationshipBusiness := business.NewRelationshipBusiness(ctx, ps.Service, ps.ProfileBusiness)

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

		relationships, err := relationshipBusiness.ListRelationships(ctx, request)
		if err != nil {
			return ps.toAPIError(err)
		}

		var responseList []*profilev1.RelationshipObject

		for _, relationship := range relationships {
			responseList = append(responseList, relationship.ToAPI())
		}

		err = server.Send(&profilev1.ListRelationshipResponse{Data: responseList})
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
