package handlers

import (
	"context"

	commonv1 "buf.build/gen/go/antinvestor/common/protocolbuffers/go/common/v1"
	"buf.build/gen/go/antinvestor/settingz/connectrpc/go/settings/v1/settingsv1connect"
	settingsv1 "buf.build/gen/go/antinvestor/settingz/protocolbuffers/go/settings/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"

	"github.com/antinvestor/service-profile/apps/settings/service/business"
	"github.com/antinvestor/service-profile/apps/settings/service/repository"
)

type SettingsServer struct {
	settingBusiness business.SettingsBusiness
	settingsv1connect.UnimplementedSettingsServiceHandler
}

func NewSettingsServer(ctx context.Context, svc *frame.Service) *SettingsServer {
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	refRepo := repository.NewReferenceRepository(ctx, dbPool, workMan)
	valRepo := repository.NewSettingValRepository(ctx, dbPool, workMan)

	return &SettingsServer{
		settingBusiness: business.NewSettingsBusiness(refRepo, valRepo),
	}
}

// Get a single setting and its stored value.
func (s *SettingsServer) Get(
	ctx context.Context,
	req *connect.Request[settingsv1.GetRequest],
) (*connect.Response[settingsv1.GetResponse], error) {
	resp, err := s.settingBusiness.Get(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// Set save the setting value appropriately.
func (s *SettingsServer) Set(
	ctx context.Context,
	req *connect.Request[settingsv1.SetRequest],
) (*connect.Response[settingsv1.SetResponse], error) {
	resp, err := s.settingBusiness.Set(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// List Pulls all setting values that match some criteria in the name & any other setting properties.
func (s *SettingsServer) List(
	ctx context.Context,
	req *connect.Request[settingsv1.ListRequest],
	stream *connect.ServerStream[settingsv1.ListResponse],
) error {
	response, err := s.settingBusiness.List(ctx, req.Msg)

	if err != nil {
		return err
	}
	return stream.Send(&settingsv1.ListResponse{Data: response})
}

// Search Streams setting values that match some criteria in the name & any other setting properties.
func (s *SettingsServer) Search(
	ctx context.Context,
	req *connect.Request[commonv1.SearchRequest],
	stream *connect.ServerStream[settingsv1.SearchResponse],
) error {
	resp, err := s.settingBusiness.Search(ctx, req.Msg)

	if err != nil {
		return err
	}

	for {
		result, ok := resp.ReadResult(ctx)

		if !ok {
			return nil
		}

		if result.IsError() {
			return result.Error()
		}

		sErr := stream.Send(&settingsv1.SearchResponse{Data: result.Item()})
		if sErr != nil {
			return sErr
		}
	}
}
