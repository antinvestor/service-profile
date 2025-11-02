package handlers

import (
	"context"

	"connectrpc.com/connect"
	devicev1 "github.com/antinvestor/apis/go/device/v1"
)

func (ds *DevicesServer) RegisterKey(ctx context.Context, req *connect.Request[devicev1.RegisterKeyRequest]) (*connect.Response[devicev1.RegisterKeyResponse], error) {

}

func (ds *DevicesServer) DeRegisterKey(ctx context.Context, req *connect.Request[devicev1.DeRegisterKeyRequest]) (*connect.Response[devicev1.DeRegisterKeyResponse], error) {

}

func (ds *DevicesServer) Notify(ctx context.Context, req *connect.Request[devicev1.NotifyRequest]) (*connect.Response[devicev1.NotifyResponse], error) {

}
