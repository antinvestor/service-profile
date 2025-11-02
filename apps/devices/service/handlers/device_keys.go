package handlers

import (
	"context"

	"connectrpc.com/connect"
	devicev1 "github.com/antinvestor/apis/go/device/v1"
)

func (ds *DevicesServer) AddKey(ctx context.Context, req *connect.Request[devicev1.AddKeyRequest]) (*connect.Response[devicev1.AddKeyResponse], error) {

	msg := req.Msg
	deviceKey, err := ds.keyBusiness.AddKey(ctx, msg.GetDeviceId(), msg.GetKeyType(), msg.GetData(), msg.GetExtras().AsMap())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&devicev1.AddKeyResponse{
		Data: deviceKey,
	}), nil
}

func (ds *DevicesServer) RemoveKey(
	ctx context.Context, req *connect.Request[devicev1.RemoveKeyRequest]) (*connect.Response[devicev1.RemoveKeyResponse], error) {

	var keyIDList []string
	response, err := ds.keyBusiness.RemoveKeys(ctx, req.Msg.GetId()...)
	if err != nil {
		return nil, err
	}
	for res := range response {
		if res.IsError() {
			return nil, res.Error()
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
	ctx context.Context, req *connect.Request[devicev1.SearchKeyRequest]) (*connect.Response[devicev1.SearchKeyResponse], error) {

	msg := req.Msg
	response, err := ds.keyBusiness.GetKeys(ctx, msg.GetDeviceId(), msg.GetKeyTypes())
	if err != nil {
		return nil, err
	}

	var keyObjList []*devicev1.KeyObject
	for res := range response {
		if res.IsError() {
			return nil, res.Error()
		}

		for _, key := range res.Item() {
			keyObjList = append(keyObjList, key)
		}
	}

	return connect.NewResponse(&devicev1.SearchKeyResponse{
		Data: keyObjList,
	}), nil

}
