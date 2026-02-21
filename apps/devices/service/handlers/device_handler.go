package handlers

import (
	"context"
	"errors"
	"slices"

	"buf.build/gen/go/antinvestor/device/connectrpc/go/device/v1/devicev1connect"
	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/devices/service/business"
	"github.com/antinvestor/service-profile/internal/errorutil"
)

type DevicesServer struct {
	devicev1connect.UnimplementedDeviceServiceHandler

	deviceBusiness   business.DeviceBusiness
	presenceBusiness business.PresenceBusiness
	keyBusiness      business.KeysBusiness
	notifyBusiness   business.NotifyBusiness
}

func NewDeviceServer(_ context.Context, deviceBusiness business.DeviceBusiness,
	presenceBusiness business.PresenceBusiness, keyBusiness business.KeysBusiness,
	notifyBusiness business.NotifyBusiness) *DevicesServer {
	return &DevicesServer{
		deviceBusiness:   deviceBusiness,
		presenceBusiness: presenceBusiness,
		keyBusiness:      keyBusiness,
		notifyBusiness:   notifyBusiness,
	}
}

// GetById retrieves one or more devices by ID using a batch-optimized path.
// nolint: revive,staticcheck,nolintlint // This is an api implementation
func (ds *DevicesServer) GetById(
	ctx context.Context,
	req *connect.Request[devicev1.GetByIdRequest],
) (*connect.Response[devicev1.GetByIdResponse], error) {
	ids := req.Msg.GetId()
	if slices.Contains(ids, "") {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("device ID cannot be empty"))
	}

	devicesList, err := ds.deviceBusiness.GetDevicesByIDs(ctx, ids)
	if err != nil {
		return nil, errorutil.CleanErr(err)
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
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&devicev1.GetBySessionIdResponse{
		Data: device,
	}), nil
}

func (ds *DevicesServer) Search(
	ctx context.Context,
	req *connect.Request[devicev1.SearchRequest],
	stream *connect.ServerStream[devicev1.SearchResponse],
) error {
	// Always process the search, even for empty queries
	response, err := ds.deviceBusiness.SearchDevices(ctx, req.Msg)
	if err != nil {
		return errorutil.CleanErr(err)
	}

	for {
		res, ok := response.ReadResult(ctx)
		if !ok {
			return nil
		}

		if res.IsError() {
			return errorutil.CleanErr(res.Error())
		}

		sErr := stream.Send(&devicev1.SearchResponse{
			Data: res.Item(),
		})
		if sErr != nil {
			return errorutil.CleanErr(sErr)
		}
	}
}

func (ds *DevicesServer) Create(
	ctx context.Context,
	req *connect.Request[devicev1.CreateRequest],
) (*connect.Response[devicev1.CreateResponse], error) {
	// Generate a device ID for tracking
	deviceID := util.IDString()

	msg := req.Msg

	// Add device name to properties if provided
	var properties data.JSONMap = msg.GetProperties().AsMap()

	if msg.GetName() != "" {
		properties["device_name"] = msg.GetName()
	}

	properties["ip"] = GetClientIP(ctx)

	// Extract Cloudflare geo headers so the queue handler can skip external GeoIP lookups.
	if cfGeo := ExtractCloudflareGeo(ctx, req.Header()); len(cfGeo) > 0 {
		properties["cf_geo"] = cfGeo
	}

	sessionID := properties.GetString("session_id")

	// Log device activity to trigger device analysis and creation
	_, err := ds.deviceBusiness.LogDeviceActivity(ctx, deviceID, sessionID, properties)
	if err != nil {
		return nil, errorutil.CleanErr(err)
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

func (ds *DevicesServer) Update(
	ctx context.Context,
	req *connect.Request[devicev1.UpdateRequest],
) (*connect.Response[devicev1.UpdateResponse], error) {
	msg := req.Msg
	device, err := ds.deviceBusiness.SaveDevice(
		ctx, msg.GetId(), msg.GetName(), msg.GetProperties().AsMap(),
	)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&devicev1.UpdateResponse{
		Data: device,
	}), nil
}

func (ds *DevicesServer) Link(
	ctx context.Context,
	req *connect.Request[devicev1.LinkRequest],
) (*connect.Response[devicev1.LinkResponse], error) {
	msg := req.Msg
	device, err := ds.deviceBusiness.LinkDeviceToProfile(
		ctx, msg.GetId(), msg.GetProfileId(), msg.GetProperties().AsMap(),
	)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&devicev1.LinkResponse{
		Data: device,
	}), nil
}

func (ds *DevicesServer) Remove(
	ctx context.Context,
	req *connect.Request[devicev1.RemoveRequest],
) (*connect.Response[devicev1.RemoveResponse], error) {
	msg := req.Msg

	dev, err := ds.deviceBusiness.RemoveDevice(ctx, msg.GetId())
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&devicev1.RemoveResponse{
		Data: dev,
	}), nil
}

func (ds *DevicesServer) Log(
	ctx context.Context,
	req *connect.Request[devicev1.LogRequest],
) (*connect.Response[devicev1.LogResponse], error) {
	msg := req.Msg

	payload := msg.GetExtras().AsMap()

	payload["ip"] = GetClientIP(ctx)

	// Extract Cloudflare geo headers so the queue handler can skip external GeoIP lookups.
	if cfGeo := ExtractCloudflareGeo(ctx, req.Header()); len(cfGeo) > 0 {
		payload["cf_geo"] = cfGeo
	}

	deviceLog, err := ds.deviceBusiness.LogDeviceActivity(ctx, msg.GetDeviceId(), msg.GetSessionId(), payload)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return connect.NewResponse(&devicev1.LogResponse{
		Data: deviceLog,
	}), nil
}

func (ds *DevicesServer) ListLogs(
	ctx context.Context,
	req *connect.Request[devicev1.ListLogsRequest],
	stream *connect.ServerStream[devicev1.ListLogsResponse],
) error {
	response, err := ds.deviceBusiness.GetDeviceLogs(ctx, req.Msg.GetDeviceId())
	if err != nil {
		return errorutil.CleanErr(err)
	}

	for {
		res, ok := response.ReadResult(ctx)
		if !ok {
			return nil
		}

		if res.IsError() {
			return errorutil.CleanErr(res.Error())
		}

		sErr := stream.Send(&devicev1.ListLogsResponse{
			Data: res.Item(),
		})
		if sErr != nil {
			return errorutil.CleanErr(sErr)
		}
	}
}
