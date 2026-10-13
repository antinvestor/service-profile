package handlers

import (
	"context"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"

	"github.com/antinvestor/service-profile/pkg/errorutil"
)

func (ds *DevicesServer) RegisterKey(
	ctx context.Context,
	req *devicev1.RegisterKeyRequest,
) (*devicev1.RegisterKeyResponse, error) {
	response, err := ds.notifyBusiness.RegisterKey(ctx, req)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &devicev1.RegisterKeyResponse{
		Data: response,
	}, nil
}

func (ds *DevicesServer) DeRegisterKey(
	ctx context.Context,
	req *devicev1.DeRegisterKeyRequest,
) (*devicev1.DeRegisterKeyResponse, error) {
	err := ds.notifyBusiness.DeRegisterKey(ctx, req)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &devicev1.DeRegisterKeyResponse{
		Success: true,
		Message: "ok",
	}, nil
}

func (ds *DevicesServer) Notify(
	ctx context.Context,
	req *devicev1.NotifyRequest,
) (*devicev1.NotifyResponse, error) {
	response, err := ds.notifyBusiness.Notify(ctx, req)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &devicev1.NotifyResponse{
		Results: response,
	}, nil
}
