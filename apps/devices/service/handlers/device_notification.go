package handlers

import (
	"context"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/security/authorizer"

	"github.com/antinvestor/service-profile/internal/errorutil"
)

func (ds *DevicesServer) RegisterKey(
	ctx context.Context,
	req *connect.Request[devicev1.RegisterKeyRequest],
) (*connect.Response[devicev1.RegisterKeyResponse], error) {
	if err := ds.authz.CanManageDevices(ctx); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	response, err := ds.notifyBusiness.RegisterKey(ctx, req.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&devicev1.RegisterKeyResponse{
		Data: response,
	}), nil
}

func (ds *DevicesServer) DeRegisterKey(
	ctx context.Context,
	req *connect.Request[devicev1.DeRegisterKeyRequest],
) (*connect.Response[devicev1.DeRegisterKeyResponse], error) {
	if err := ds.authz.CanManageDevices(ctx); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	err := ds.notifyBusiness.DeRegisterKey(ctx, req.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&devicev1.DeRegisterKeyResponse{
		Success: true,
		Message: "ok",
	}), nil
}

func (ds *DevicesServer) Notify(
	ctx context.Context,
	req *connect.Request[devicev1.NotifyRequest],
) (*connect.Response[devicev1.NotifyResponse], error) {
	if err := ds.authz.CanManageDevices(ctx); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	response, err := ds.notifyBusiness.Notify(ctx, req.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&devicev1.NotifyResponse{
		Results: response,
	}), nil
}
