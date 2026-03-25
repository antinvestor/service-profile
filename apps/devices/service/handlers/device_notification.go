package handlers

import (
	"context"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"connectrpc.com/connect"

	"github.com/antinvestor/service-profile/internal/errorutil"
)

func (ds *DevicesServer) RegisterKey(
	ctx context.Context,
	req *connect.Request[devicev1.RegisterKeyRequest],
) (*connect.Response[devicev1.RegisterKeyResponse], error) {
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
	response, err := ds.notifyBusiness.Notify(ctx, req.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&devicev1.NotifyResponse{
		Results: response,
	}), nil
}
