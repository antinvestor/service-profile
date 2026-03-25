package handlers

import (
	"context"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"connectrpc.com/connect"

	"github.com/antinvestor/service-profile/internal/errorutil"
)

func (ds *DevicesServer) AddKey(
	ctx context.Context,
	req *connect.Request[devicev1.AddKeyRequest],
) (*connect.Response[devicev1.AddKeyResponse], error) {
	msg := req.Msg
	deviceKey, err := ds.keyBusiness.AddKey(
		ctx,
		msg.GetDeviceId(),
		msg.GetKeyType(),
		msg.GetData(),
		msg.GetExtras().AsMap(),
	)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return connect.NewResponse(&devicev1.AddKeyResponse{
		Data: deviceKey,
	}), nil
}

func (ds *DevicesServer) RemoveKey(
	ctx context.Context,
	req *connect.Request[devicev1.RemoveKeyRequest],
) (*connect.Response[devicev1.RemoveKeyResponse], error) {
	var keyIDList []string
	response, err := ds.keyBusiness.RemoveKeys(ctx, req.Msg.GetId()...)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	for res := range response {
		if res.IsError() {
			return nil, errorutil.CleanErr(res.Error())
		}

		for _, key := range res.Item() {
			keyIDList = append(keyIDList, key.GetId())
		}
	}
	return connect.NewResponse(&devicev1.RemoveKeyResponse{
		Id: keyIDList,
	}), nil
}

func (ds *DevicesServer) SearchKey(
	ctx context.Context,
	req *connect.Request[devicev1.SearchKeyRequest],
) (*connect.Response[devicev1.SearchKeyResponse], error) {
	msg := req.Msg
	response, err := ds.keyBusiness.GetKeys(ctx, msg.GetDeviceId(), msg.GetKeyTypes()...)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	var keyObjList []*devicev1.KeyObject
	for res := range response {
		if res.IsError() {
			return nil, errorutil.CleanErr(res.Error())
		}

		keyObjList = append(keyObjList, res.Item()...)
	}

	return connect.NewResponse(&devicev1.SearchKeyResponse{
		Data: keyObjList,
	}), nil
}
