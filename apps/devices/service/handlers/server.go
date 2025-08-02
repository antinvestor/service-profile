package handlers

import (
	"context"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/antinvestor/service-profile/apps/devices/service/business"
)

type DevicesServer struct {
	devicev1.UnimplementedDeviceServiceServer
	Service *frame.Service

	Biz business.DeviceBusiness
}

func NewDeviceServer(ctx context.Context, svc *frame.Service) *DevicesServer {
	return &DevicesServer{
		Service: svc,
		Biz:     business.NewDeviceBusiness(ctx, svc),
	}
}

func (ds *DevicesServer) GetByID(ctx context.Context, req *devicev1.GetByIdRequest) (*devicev1.GetByIdResponse, error) {
	if len(req.GetId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "device ID is required")
	}

	var devicesList []*devicev1.DeviceObject
	var lastError error

	for _, idStr := range req.GetId() {
		if idStr == "" {
			return nil, status.Error(codes.InvalidArgument, "device ID cannot be empty")
		}

		device, err := ds.Biz.GetDeviceByID(ctx, idStr)
		if err != nil {
			lastError = err
			continue
		}
		devicesList = append(devicesList, device)
	}

	// If no devices found and we had errors, return the last error
	if len(devicesList) == 0 && lastError != nil {
		return nil, status.Error(codes.NotFound, "device not found")
	}

	return &devicev1.GetByIdResponse{
		Data: devicesList,
	}, nil
}

func (ds *DevicesServer) GetBySessionID(
	ctx context.Context,
	req *devicev1.GetBySessionIdRequest,
) (*devicev1.GetBySessionIdResponse, error) {
	device, err := ds.Biz.GetDeviceBySessionID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &devicev1.GetBySessionIdResponse{
		Data: device,
	}, nil
}

func (ds *DevicesServer) Search(req *devicev1.SearchRequest, stream devicev1.DeviceService_SearchServer) error {
	ctx := stream.Context()

	// Always process the search, even for empty queries
	response, err := ds.Biz.SearchDevices(ctx, req)
	if err != nil {
		return err
	}

	for res := range response {
		if res.IsError() {
			return res.Error()
		}

		err = stream.Send(&devicev1.SearchResponse{
			Data: res.Item(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (ds *DevicesServer) Create(ctx context.Context, req *devicev1.CreateRequest) (*devicev1.CreateResponse, error) {
	// Generate a device ID for tracking
	deviceID := util.IDString()

	// Add device name to properties if provided
	properties := req.GetProperties()
	if properties == nil {
		properties = make(map[string]string)
	}
	if req.GetName() != "" {
		properties["device_name"] = req.GetName()
	}

	// Log device activity to trigger device analysis and creation
	_, err := ds.Biz.LogDeviceActivity(ctx, deviceID, properties["session_id"], properties)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to log device activity: %v", err)
	}

	// The device will be created asynchronously by the queue handler
	// For now, return a response indicating the device creation is in progress
	return &devicev1.CreateResponse{
		Data: &devicev1.DeviceObject{
			Id:   deviceID,
			Name: req.GetName(),
		},
	}, nil
}

func (ds *DevicesServer) Update(ctx context.Context, req *devicev1.UpdateRequest) (*devicev1.UpdateResponse, error) {
	device, err := ds.Biz.SaveDevice(
		ctx, req.GetId(), req.GetName(), req.GetProperties(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update device: %v", err)
	}

	return &devicev1.UpdateResponse{
		Data: device,
	}, nil
}

func (ds *DevicesServer) Remove(ctx context.Context, req *devicev1.RemoveRequest) (*devicev1.RemoveResponse, error) {
	dev, err := ds.Biz.GetDeviceByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	err = ds.Biz.RemoveDevice(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &devicev1.RemoveResponse{
		Data: dev,
	}, nil
}

func (ds *DevicesServer) Log(ctx context.Context, req *devicev1.LogRequest) (*devicev1.LogResponse, error) {
	data := req.GetExtras()

	data["ip"] = GetClientIP(ctx)
	deviceLog, err := ds.Biz.LogDeviceActivity(ctx, req.GetDeviceId(), req.GetSessionId(), data)
	if err != nil {
		return nil, err
	}

	return &devicev1.LogResponse{
		Data: deviceLog,
	}, nil
}

func (ds *DevicesServer) ListLogs(req *devicev1.ListLogsRequest, stream devicev1.DeviceService_ListLogsServer) error {
	ctx := stream.Context()

	response, err := ds.Biz.GetDeviceLogs(ctx, req.GetDeviceId())
	if err != nil {
		return err
	}

	for res := range response {
		if res.IsError() {
			return res.Error()
		}

		err = stream.Send(&devicev1.ListLogsResponse{
			Data: res.Item(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (ds *DevicesServer) AddKey(ctx context.Context, req *devicev1.AddKeyRequest) (*devicev1.AddKeyResponse, error) {
	deviceKey, err := ds.Biz.AddKey(ctx, req.GetDeviceId(), req.GetKeyType(), req.GetData(), req.GetExtras())
	if err != nil {
		return nil, err
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
	response, err := ds.Biz.RemoveKeys(ctx, req.GetId()...)
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
	return &devicev1.RemoveKeyResponse{
		Id: keyIDList,
	}, nil
}

func (ds *DevicesServer) SearchKeys(
	req *devicev1.SearchKeyRequest,
	stream devicev1.DeviceService_SearchKeyServer,
) error {
	ctx := stream.Context()

	response, err := ds.Biz.GetKeys(ctx, req.GetDeviceId(), req.GetKeyType())
	if err != nil {
		return err
	}

	for res := range response {
		if res.IsError() {
			return res.Error()
		}
		resp := &devicev1.SearchKeyResponse{
			Data: res.Item(),
		}

		err = stream.Send(resp)
		if err != nil {
			return err
		}
	}
	return nil
}
