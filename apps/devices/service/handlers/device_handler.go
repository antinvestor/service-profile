package handlers

import (
	"context"

	"connectrpc.com/connect"
	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/antinvestor/apis/go/device/v1/devicev1connect"
	aconfig "github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/business"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DevicesServer struct {
	devicev1connect.UnimplementedDeviceServiceHandler

	deviceBusiness   business.DeviceBusiness
	presenceBusiness business.PresenceBusiness
	keyBusiness      business.KeysBusiness
	notifyBusiness   business.NotifyBusiness
}

func NewDeviceServer(ctx context.Context, svc *frame.Service) *DevicesServer {

	qMan := svc.QueueManager(ctx)
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	cfg := svc.Config().(*aconfig.DevicesConfig)

	deviceLogRepo := repository.NewDeviceLogRepository(ctx, dbPool, workMan)
	deviceSessionRepo := repository.NewDeviceSessionRepository(ctx, dbPool, workMan)
	deviceRepo := repository.NewDeviceRepository(ctx, dbPool, workMan)
	deviceKeyRepo := repository.NewDeviceKeyRepository(ctx, dbPool, workMan)
	devicePresenceRepo := repository.NewDevicePresenceRepository(ctx, dbPool, workMan)

	deviceBusiness := business.NewDeviceBusiness(ctx, cfg, qMan, workMan, deviceRepo, deviceLogRepo, deviceSessionRepo)
	keyBusiness := business.NewKeysBusiness(ctx, cfg, qMan, workMan, deviceRepo, deviceKeyRepo)
	presenceBusiness := business.NewPresenceBusiness(ctx, cfg, qMan, workMan, deviceRepo, devicePresenceRepo)
	notifyBusiness := business.NewNotifyBusiness(ctx, cfg, qMan, workMan, keyBusiness, deviceRepo)

	return &DevicesServer{
		deviceBusiness:   deviceBusiness,
		presenceBusiness: presenceBusiness,
		keyBusiness:      keyBusiness,
		notifyBusiness:   notifyBusiness,
	}
}

// GetById retrieves a device by ID
// nolint: revive,staticcheck,nolintlint // This is an api implementation
func (ds *DevicesServer) GetById(ctx context.Context, req *connect.Request[devicev1.GetByIdRequest]) (*connect.Response[devicev1.GetByIdResponse], error) {

	var devicesList []*devicev1.DeviceObject
	var lastError error

	for _, idStr := range req.Msg.GetId() {
		if idStr == "" {
			return nil, status.Error(codes.InvalidArgument, "device ID cannot be empty")
		}

		device, err := ds.deviceBusiness.GetDeviceByID(ctx, idStr)
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

	return connect.NewResponse(&devicev1.GetByIdResponse{
		Data: devicesList,
	}), nil
}

// GetBySessionId retrieves a device by session ID
// nolint: revive,staticcheck,nolintlint // This is an api implementation
func (ds *DevicesServer) GetBySessionId(
	ctx context.Context,
	req *connect.Request[devicev1.GetBySessionIdRequest]) (*connect.Response[devicev1.GetBySessionIdResponse], error) {
	device, err := ds.deviceBusiness.GetDeviceBySessionID(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&devicev1.GetBySessionIdResponse{
		Data: device,
	}), nil
}

func (ds *DevicesServer) Search(ctx context.Context, req *connect.Request[devicev1.SearchRequest], stream *connect.ServerStream[devicev1.SearchResponse]) error {

	// Always process the search, even for empty queries
	response, err := ds.deviceBusiness.SearchDevices(ctx, req.Msg)
	if err != nil {
		return err
	}

	for {
		res, ok := response.ReadResult(ctx)
		if !ok {
			return nil
		}

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
}

func (ds *DevicesServer) Create(ctx context.Context, req *connect.Request[devicev1.CreateRequest]) (*connect.Response[devicev1.CreateResponse], error) {
	// Generate a device ID for tracking
	deviceID := util.IDString()

	msg := req.Msg

	// Add device name to properties if provided
	var properties data.JSONMap = msg.GetProperties().AsMap()

	if msg.GetName() != "" {
		properties["device_name"] = msg.GetName()
	}

	sessionID := properties.GetString("session_id")

	// Log device activity to trigger device analysis and creation
	_, err := ds.deviceBusiness.LogDeviceActivity(ctx, deviceID, sessionID, properties)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to log device activity: %v", err)
	}

	// The device will be created asynchronously by the queue handler
	// For now, return a response indicating the device creation is in progress
	return connect.NewResponse(&devicev1.CreateResponse{
		Data: &devicev1.DeviceObject{
			Id:   deviceID,
			Name: msg.GetName(),
		},
	}), nil
}

func (ds *DevicesServer) Update(ctx context.Context, req *connect.Request[devicev1.UpdateRequest]) (*connect.Response[devicev1.UpdateResponse], error) {

	msg := req.Msg
	device, err := ds.deviceBusiness.SaveDevice(
		ctx, msg.GetId(), msg.GetName(), msg.GetProperties().AsMap(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update device: %v", err)
	}

	return connect.NewResponse(&devicev1.UpdateResponse{
		Data: device,
	}), nil
}

func (ds *DevicesServer) Link(ctx context.Context, req *connect.Request[devicev1.LinkRequest]) (*connect.Response[devicev1.LinkResponse], error) {

	msg := req.Msg
	device, err := ds.deviceBusiness.LinkDeviceToProfile(
		ctx, msg.GetId(), msg.GetProfileId(), msg.GetProperties().AsMap(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to link session: %v", err)
	}

	return connect.NewResponse(&devicev1.LinkResponse{
		Data: device,
	}), nil
}

func (ds *DevicesServer) Remove(ctx context.Context, req *connect.Request[devicev1.RemoveRequest]) (*connect.Response[devicev1.RemoveResponse], error) {

	msg := req.Msg

	dev, err := ds.deviceBusiness.GetDeviceByID(ctx, msg.GetId())
	if err != nil {
		return nil, err
	}

	err = ds.deviceBusiness.RemoveDevice(ctx, msg.GetId())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&devicev1.RemoveResponse{
		Data: dev,
	}), nil
}

func (ds *DevicesServer) Log(ctx context.Context, req *connect.Request[devicev1.LogRequest]) (*connect.Response[devicev1.LogResponse], error) {

	msg := req.Msg

	payload := msg.GetExtras().AsMap()

	payload["ip"] = GetClientIP(ctx)
	deviceLog, err := ds.deviceBusiness.LogDeviceActivity(ctx, msg.GetDeviceId(), msg.GetSessionId(), payload)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&devicev1.LogResponse{
		Data: deviceLog,
	}), nil
}

func (ds *DevicesServer) ListLogs(ctx context.Context, req *connect.Request[devicev1.ListLogsRequest], stream *connect.ServerStream[devicev1.ListLogsResponse]) error {

	response, err := ds.deviceBusiness.GetDeviceLogs(ctx, req.Msg.GetDeviceId())
	if err != nil {
		return err
	}

	for {
		res, ok := response.ReadResult(ctx)
		if !ok {
			return nil
		}

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
}
