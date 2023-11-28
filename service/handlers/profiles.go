package handlers

import (
	"context"
	napi "github.com/antinvestor/apis/notification"
	papi "github.com/antinvestor/apis/profile"
	"github.com/antinvestor/service-profile/service/business"
	"github.com/pitabwire/frame"

	"strings"
)

type ProfileServer struct {
	EncryptionKey []byte

	Service         *frame.Service
	NotificationCli *napi.NotificationClient

	papi.ProfileServiceServer
}

func (ps *ProfileServer) GetByID(ctx context.Context,
	request *papi.ProfileIDRequest) (*papi.ProfileObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	profileID := strings.TrimSpace(request.GetID())

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	return profileBusiness.GetByID(ctx, ps.EncryptionKey, profileID)
}

func (ps *ProfileServer) Search(request *papi.ProfileSearchRequest, stream papi.ProfileService_SearchServer) error {

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	return profileBusiness.SearchProfile(ctx, ps.EncryptionKey, request, stream)

}

func (ps *ProfileServer) Merge(ctx context.Context, request *papi.ProfileMergeRequest) (
	*papi.ProfileObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	return profileBusiness.MergeProfile(ctx, ps.EncryptionKey, request)

}

func (ps *ProfileServer) Create(ctx context.Context, request *papi.ProfileCreateRequest) (
	*papi.ProfileObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	return profileBusiness.CreateProfile(ctx, ps.EncryptionKey, request)

}

func (ps *ProfileServer) Update(ctx context.Context, request *papi.ProfileUpdateRequest) (
	*papi.ProfileObject, error) {

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	return profileBusiness.UpdateProfile(ctx, ps.EncryptionKey, request)

}

// AddAddress Adds a new address based on the request.
func (ps *ProfileServer) AddAddress(ctx context.Context,
	request *papi.ProfileAddAddressRequest) (*papi.ProfileObject, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	return profileBusiness.AddAddress(ctx, ps.EncryptionKey, request)
}

func (ps *ProfileServer) GetByContact(ctx context.Context,
	request *papi.ProfileContactRequest) (*papi.ProfileObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	return profileBusiness.GetByContact(ctx, ps.EncryptionKey, request.GetContact())
}

func (ps *ProfileServer) AddContact(ctx context.Context, request *papi.ProfileAddContactRequest,
) (*papi.ProfileObject, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	return profileBusiness.AddContact(ctx, ps.EncryptionKey, request)
}

func (ps *ProfileServer) AddRelationship(ctx context.Context, request *papi.ProfileAddRelationshipRequest) (*papi.RelationshipObject, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	relationshipBusiness := business.NewRelationshipBusiness(ctx, ps.Service, ps.EncryptionKey)
	return relationshipBusiness.CreateRelationship(ctx, request)
}

func (ps *ProfileServer) DeleteRelationship(ctx context.Context, request *papi.ProfileDeleteRelationshipRequest) (*papi.RelationshipObject, error) {

	if err := request.Validate(); err != nil {
		return nil, err
	}

	relationshipBusiness := business.NewRelationshipBusiness(ctx, ps.Service, ps.EncryptionKey)
	return relationshipBusiness.DeleteRelationship(ctx, request)
}

func (ps *ProfileServer) ListRelationships(request *papi.ProfileListRelationshipRequest, server papi.ProfileService_ListRelationshipsServer) error {

	if err := request.Validate(); err != nil {
		return err
	}

	ctx := server.Context()

	relationshipBusiness := business.NewRelationshipBusiness(ctx, ps.Service, ps.EncryptionKey)
	relationships, err := relationshipBusiness.ListRelationships(ctx, request)
	if err != nil {
		return err
	}

	for _, relationship := range relationships {

		relationshipObject, err1 := relationshipBusiness.ToAPI(ctx, request.GetParent(), request.GetParentID(), relationship)
		if err1 != nil {
			return err
		}

		err = server.Send(relationshipObject)
		if err != nil {
			return err
		}
	}

	return nil

}
