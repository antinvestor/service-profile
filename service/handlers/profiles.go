package handlers

import (
	"context"
	napi "github.com/antinvestor/service-notification-api"
	papi "github.com/antinvestor/service-profile-api"
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
	request *papi.ProfileIDRequest, ) (*papi.ProfileObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	profileID := strings.TrimSpace(request.GetID())

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	return profileBusiness.GetByID(ctx, ps.EncryptionKey, profileID)
}

func (ps *ProfileServer) Search(request *papi.ProfileSearchRequest, stream papi.ProfileService_SearchServer, ) error {

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	return profileBusiness.SearchProfile(ctx, ps.EncryptionKey, request, stream)

}

func (ps *ProfileServer) Merge(ctx context.Context, request *papi.ProfileMergeRequest, ) (
	*papi.ProfileObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	return profileBusiness.MergeProfile(ctx, ps.EncryptionKey, request)

}

func (ps *ProfileServer) Create(ctx context.Context, request *papi.ProfileCreateRequest, ) (
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
