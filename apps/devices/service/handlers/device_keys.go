package handlers

import (
	"context"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"

	"github.com/antinvestor/service-profile/pkg/errorutil"
)

func (ds *DevicesServer) AddKey(
	ctx context.Context,
	req *devicev1.AddKeyRequest,
) (*devicev1.AddKeyResponse, error) {
	msg := req
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
	return &devicev1.AddKeyResponse{
		Data: deviceKey,
	}, nil
}

func (ds *DevicesServer) RemoveKey(
	ctx context.Context,
	req *devicev1.RemoveKeyRequest,
) (*devicev1.RemoveKeyResponse, error) {
	var keyIDList []string
	response, err := ds.keyBusiness.RemoveKeys(ctx, req.GetId()...)
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
	return &devicev1.RemoveKeyResponse{
		Id: keyIDList,
	}, nil
}

func (ds *DevicesServer) SearchKey(
	ctx context.Context,
	req *devicev1.SearchKeyRequest,
) (*devicev1.SearchKeyResponse, error) {
	msg := req
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

	return &devicev1.SearchKeyResponse{
		Data: keyObjList,
	}, nil
}
