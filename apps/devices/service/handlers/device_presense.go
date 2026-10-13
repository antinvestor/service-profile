package handlers

import (
	"context"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"

	"github.com/antinvestor/service-profile/pkg/errorutil"
)

func (ds *DevicesServer) UpdatePresence(
	ctx context.Context,
	req *devicev1.UpdatePresenceRequest,
) (*devicev1.UpdatePresenceResponse, error) {
	presence, err := ds.presenceBusiness.UpdatePresence(ctx, req)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &devicev1.UpdatePresenceResponse{
		Data: presence,
	}, nil
}
