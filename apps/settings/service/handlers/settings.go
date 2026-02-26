package handlers

import (
	"context"
	"errors"

	commonv1 "buf.build/gen/go/antinvestor/common/protocolbuffers/go/common/v1"
	"buf.build/gen/go/antinvestor/settingz/connectrpc/go/settings/v1/settingsv1connect"
	settingsv1 "buf.build/gen/go/antinvestor/settingz/protocolbuffers/go/settings/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/security/authorizer"

	"github.com/antinvestor/service-profile/apps/settings/service/authz"
	"github.com/antinvestor/service-profile/apps/settings/service/business"
	"github.com/antinvestor/service-profile/apps/settings/service/repository"
	"github.com/antinvestor/service-profile/internal/errorutil"
)

// toConnectError converts authorisation errors into appropriate connect errors.
func toConnectError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, authorizer.ErrInvalidSubject) || errors.Is(err, authorizer.ErrInvalidObject) {
		return connect.NewError(connect.CodeUnauthenticated, err)
	}

	var permErr *authorizer.PermissionDeniedError
	if errors.As(err, &permErr) {
		return connect.NewError(connect.CodePermissionDenied, err)
	}

	return connect.NewError(connect.CodeInternal, err)
}

type SettingsServer struct {
	authz           authz.Middleware
	settingBusiness business.SettingsBusiness
	settingsv1connect.UnimplementedSettingsServiceHandler
}

func NewSettingsServer(ctx context.Context, svc *frame.Service, authzMiddleware authz.Middleware) *SettingsServer {
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	refRepo := repository.NewReferenceRepository(ctx, dbPool, workMan)
	valRepo := repository.NewSettingValRepository(ctx, dbPool, workMan)

	return &SettingsServer{
		authz:           authzMiddleware,
		settingBusiness: business.NewSettingsBusiness(refRepo, valRepo),
	}
}

// Get a single setting and its stored value.
func (s *SettingsServer) Get(
	ctx context.Context,
	req *connect.Request[settingsv1.GetRequest],
) (*connect.Response[settingsv1.GetResponse], error) {
	if err := s.authz.CanViewSettings(ctx); err != nil {
		return nil, toConnectError(err)
	}

	resp, err := s.settingBusiness.Get(ctx, req.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return connect.NewResponse(resp), nil
}

// Set save the setting value appropriately.
func (s *SettingsServer) Set(
	ctx context.Context,
	req *connect.Request[settingsv1.SetRequest],
) (*connect.Response[settingsv1.SetResponse], error) {
	if err := s.authz.CanManageSettings(ctx); err != nil {
		return nil, toConnectError(err)
	}

	resp, err := s.settingBusiness.Set(ctx, req.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return connect.NewResponse(resp), nil
}

// List Pulls all setting values that match some criteria in the name & any other setting properties.
func (s *SettingsServer) List(
	ctx context.Context,
	req *connect.Request[settingsv1.ListRequest],
	stream *connect.ServerStream[settingsv1.ListResponse],
) error {
	if err := s.authz.CanViewSettings(ctx); err != nil {
		return toConnectError(err)
	}

	response, err := s.settingBusiness.List(ctx, req.Msg)

	if err != nil {
		return errorutil.CleanErr(err)
	}
	return stream.Send(&settingsv1.ListResponse{Data: response})
}

// Search Streams setting values that match some criteria in the name & any other setting properties.
func (s *SettingsServer) Search(
	ctx context.Context,
	req *connect.Request[commonv1.SearchRequest],
	stream *connect.ServerStream[settingsv1.SearchResponse],
) error {
	if err := s.authz.CanViewSettings(ctx); err != nil {
		return toConnectError(err)
	}

	resp, err := s.settingBusiness.Search(ctx, req.Msg)

	if err != nil {
		return errorutil.CleanErr(err)
	}

	for {
		result, ok := resp.ReadResult(ctx)

		if !ok {
			return nil
		}

		if result.IsError() {
			return errorutil.CleanErr(result.Error())
		}

		sErr := stream.Send(&settingsv1.SearchResponse{Data: result.Item()})
		if sErr != nil {
			return errorutil.CleanErr(sErr)
		}
	}
}
