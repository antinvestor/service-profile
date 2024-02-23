package handlers

import (
	"context"
	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service/business"
	"github.com/pitabwire/frame"
)

type ProfileServer struct {
	EncryptionKeyFunc func() []byte

	Service         *frame.Service
	NotificationCli *notificationv1.NotificationClient

	profilev1.UnimplementedProfileServiceServer
}

func (ps *ProfileServer) GetById(ctx context.Context,
	request *profilev1.GetByIdRequest) (*profilev1.GetByIdResponse, error) {

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	profileObj, err := profileBusiness.GetByID(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	return &profilev1.GetByIdResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) GetByContact(ctx context.Context,
	request *profilev1.GetByContactRequest) (*profilev1.GetByContactResponse, error) {

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	profileObj, err := profileBusiness.GetByContact(ctx, request.GetContact())

	if err != nil {
		return nil, err
	}

	return &profilev1.GetByContactResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) Search(request *profilev1.SearchRequest, stream profilev1.ProfileService_SearchServer) error {

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	return profileBusiness.SearchProfile(ctx, request, stream)

}

func (ps *ProfileServer) Merge(ctx context.Context, request *profilev1.MergeRequest) (
	*profilev1.MergeResponse, error) {

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	profileObj, err := profileBusiness.MergeProfile(ctx, request)
	if err != nil {
		return nil, err
	}

	return &profilev1.MergeResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) Create(ctx context.Context, request *profilev1.CreateRequest) (
	*profilev1.CreateResponse, error) {

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	profileObj, err := profileBusiness.CreateProfile(ctx, request)

	if err != nil {
		return nil, err
	}

	return &profilev1.CreateResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) Update(ctx context.Context, request *profilev1.UpdateRequest) (
	*profilev1.UpdateResponse, error) {

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	profileObj, err := profileBusiness.UpdateProfile(ctx, request)

	if err != nil {
		return nil, err
	}

	return &profilev1.UpdateResponse{Data: profileObj}, nil
}

// AddAddress Adds a new address based on the request.
func (ps *ProfileServer) AddAddress(ctx context.Context,
	request *profilev1.AddAddressRequest) (*profilev1.AddAddressResponse, error) {

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	profileObj, err := profileBusiness.AddAddress(ctx, request)
	if err != nil {
		return nil, err
	}

	return &profilev1.AddAddressResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) AddContact(ctx context.Context, request *profilev1.AddContactRequest,
) (*profilev1.AddContactResponse, error) {

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	profileObj, err := profileBusiness.AddContact(ctx, request)

	if err != nil {
		return nil, err
	}

	return &profilev1.AddContactResponse{Data: profileObj}, nil
}

func (ps *ProfileServer) AddRelationship(ctx context.Context,
	request *profilev1.AddRelationshipRequest) (*profilev1.AddRelationshipResponse, error) {

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	relationshipBusiness := business.NewRelationshipBusiness(ctx, ps.Service, profileBusiness)
	relationshipObj, err := relationshipBusiness.CreateRelationship(ctx, request)

	if err != nil {
		return nil, err
	}

	return &profilev1.AddRelationshipResponse{Data: relationshipObj}, nil
}

func (ps *ProfileServer) DeleteRelationship(ctx context.Context,
	request *profilev1.DeleteRelationshipRequest) (*profilev1.DeleteRelationshipResponse, error) {

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	relationshipBusiness := business.NewRelationshipBusiness(ctx, ps.Service, profileBusiness)
	relationshipObj, err := relationshipBusiness.DeleteRelationship(ctx, request)

	if err != nil {
		return nil, err
	}

	return &profilev1.DeleteRelationshipResponse{Data: relationshipObj}, nil
}

func (ps *ProfileServer) ListRelationships(request *profilev1.ListRelationshipRequest, server profilev1.ProfileService_ListRelationshipServer) error {

	ctx := server.Context()

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	relationshipBusiness := business.NewRelationshipBusiness(ctx, ps.Service, profileBusiness)
	relationships, err := relationshipBusiness.ListRelationships(ctx, request)
	if err != nil {
		return err
	}

	var responseList []*profilev1.RelationshipObject

	for _, relationship := range relationships {

		relationshipObject, err1 := relationshipBusiness.ToAPI(ctx, request.GetParent(), request.GetParentId(), relationship)
		if err1 != nil {
			return err
		}

		responseList = append(responseList, relationshipObject)
	}

	err = server.Send(&profilev1.ListRelationshipResponse{Data: responseList})
	if err != nil {
		return err
	}

	return nil

}
