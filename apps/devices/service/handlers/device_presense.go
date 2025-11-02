package handlers

import (
	"context"

	"connectrpc.com/connect"
	devicev1 "github.com/antinvestor/apis/go/device/v1"
)

func (ds *DevicesServer) UpdatePresence(ctx context.Context, req *connect.Request[devicev1.UpdatePresenceRequest]) (*connect.Response[devicev1.UpdatePresenceResponse], error) {
	presence, err := ds.presenceBusiness.UpdatePresence(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&devicev1.UpdatePresenceResponse{
		Data: presence,
	}), nil
}
