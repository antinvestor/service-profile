package handlers

import (
	"context"
	"slices"
	"time"

	"buf.build/gen/go/antinvestor/device/connectrpc/go/v2/device/v1/devicev1connect"
	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"connectrpc.com/connect/v2"
	"github.com/pitabwire/frame/v3/data"
	"github.com/pitabwire/frame/v3/security"
	"github.com/pitabwire/frame/v3/security/authorizer"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/devices/service/authz"
	"github.com/antinvestor/service-profile/apps/devices/service/business"
	"github.com/antinvestor/service-profile/apps/devices/service/caching"
	"github.com/antinvestor/service-profile/pkg/errorutil"
)

const prefixRateTURN = "rate:turn:"

type DevicesServer struct {
	devicev1connect.UnimplementedDeviceServiceHandler

	checker          *authorizer.FunctionChecker
	deviceBusiness   business.DeviceBusiness
	presenceBusiness business.PresenceBusiness
	keyBusiness      business.KeysBusiness
	notifyBusiness   business.NotifyBusiness
	turnBusiness     business.TURNBusiness

	cacheSvc               *caching.DeviceCacheService
	turnTTL                int32
	rateLimitTURNPerMinute int64
}

func NewDeviceServer(_ context.Context, checker *authorizer.FunctionChecker,
	deviceBusiness business.DeviceBusiness,
	presenceBusiness business.PresenceBusiness, keyBusiness business.KeysBusiness,
	notifyBusiness business.NotifyBusiness, turnBusiness business.TURNBusiness,
	cacheSvc *caching.DeviceCacheService, turnTTL int32, rateLimitTURNPerMinute int64,
) *DevicesServer {
	return &DevicesServer{
		checker:                checker,
		deviceBusiness:         deviceBusiness,
		presenceBusiness:       presenceBusiness,
		keyBusiness:            keyBusiness,
		notifyBusiness:         notifyBusiness,
		turnBusiness:           turnBusiness,
		cacheSvc:               cacheSvc,
		turnTTL:                turnTTL,
		rateLimitTURNPerMinute: rateLimitTURNPerMinute,
	}
}

// GetById retrieves one or more devices by ID using a batch-optimized path.
// nolint: revive,staticcheck,nolintlint // This is an api implementation
func (ds *DevicesServer) GetById(
	ctx context.Context,
	req *devicev1.GetByIdRequest,
) (*devicev1.GetByIdResponse, error) {
	ids := req.GetId()
	if slices.Contains(ids, "") {
		return nil, connect.NewError(connect.CodeInvalidArgument, "device ID cannot be empty")
	}

	devicesList, err := ds.deviceBusiness.GetDevicesByIDs(ctx, ids)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &devicev1.GetByIdResponse{
		Data: devicesList,
	}, nil
}

// GetBySessionId retrieves a device by session ID
// nolint: revive,staticcheck,nolintlint // This is an api implementation
func (ds *DevicesServer) GetBySessionId(
	ctx context.Context,
	req *devicev1.GetBySessionIdRequest) (*devicev1.GetBySessionIdResponse, error) {
	device, err := ds.deviceBusiness.GetDeviceBySessionID(ctx, req.GetId())
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &devicev1.GetBySessionIdResponse{
		Data: device,
	}, nil
}

func (ds *DevicesServer) Search(
	ctx context.Context,
	req *devicev1.SearchRequest,
	stream devicev1connect.DeviceServiceSearchServerStream,
) error {
	// Always process the search, even for empty queries
	response, err := ds.deviceBusiness.SearchDevices(ctx, req)
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
	req *devicev1.CreateRequest,
) (*devicev1.CreateResponse, error) {
	// Generate a device ID for tracking
	deviceID := util.IDString()

	msg := req

	// Add device name to properties if provided
	var properties data.JSONMap = msg.GetProperties().AsMap()

	if msg.GetName() != "" {
		properties["device_name"] = msg.GetName()
	}

	properties["ip"] = GetClientIP(ctx)

	// Extract Cloudflare geo headers so the queue handler can skip external GeoIP lookups.
	if cfGeo := ExtractCloudflareGeo(ctx); len(cfGeo) > 0 {
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
	return &devicev1.CreateResponse{
		Data: &devicev1.DeviceObject{
			Id:   deviceID,
			Name: msg.GetName(),
		},
	}, nil
}

func (ds *DevicesServer) Update(
	ctx context.Context,
	req *devicev1.UpdateRequest,
) (*devicev1.UpdateResponse, error) {
	msg := req
	device, err := ds.deviceBusiness.SaveDevice(
		ctx, msg.GetId(), msg.GetName(), msg.GetProperties().AsMap(),
	)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &devicev1.UpdateResponse{
		Data: device,
	}, nil
}

func (ds *DevicesServer) Link(
	ctx context.Context,
	req *devicev1.LinkRequest,
) (*devicev1.LinkResponse, error) {
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != req.GetProfileId() {
		if err := ds.checker.Check(ctx, authz.PermissionDevicesManage); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	msg := req
	device, err := ds.deviceBusiness.LinkDeviceToProfile(
		ctx, msg.GetId(), msg.GetProfileId(), msg.GetProperties().AsMap(),
	)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &devicev1.LinkResponse{
		Data: device,
	}, nil
}

func (ds *DevicesServer) Remove(
	ctx context.Context,
	req *devicev1.RemoveRequest,
) (*devicev1.RemoveResponse, error) {
	msg := req

	dev, err := ds.deviceBusiness.RemoveDevice(ctx, msg.GetId())
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &devicev1.RemoveResponse{
		Data: dev,
	}, nil
}

func (ds *DevicesServer) Log(
	ctx context.Context,
	req *devicev1.LogRequest,
) (*devicev1.LogResponse, error) {
	msg := req

	payload := msg.GetExtras().AsMap()

	payload["ip"] = GetClientIP(ctx)

	// Extract Cloudflare geo headers so the queue handler can skip external GeoIP lookups.
	if cfGeo := ExtractCloudflareGeo(ctx); len(cfGeo) > 0 {
		payload["cf_geo"] = cfGeo
	}

	deviceLog, err := ds.deviceBusiness.LogDeviceActivity(ctx, msg.GetDeviceId(), msg.GetSessionId(), payload)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	return &devicev1.LogResponse{
		Data: deviceLog,
	}, nil
}

func (ds *DevicesServer) ListLogs(
	ctx context.Context,
	req *devicev1.ListLogsRequest,
	stream devicev1connect.DeviceServiceListLogsServerStream,
) error {
	response, err := ds.deviceBusiness.GetDeviceLogs(ctx, req.GetDeviceId())
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

func (ds *DevicesServer) GetTurnCredentials(
	ctx context.Context,
	_ *devicev1.GetTurnCredentialsRequest,
) (*devicev1.GetTurnCredentialsResponse, error) {
	if ds.turnBusiness == nil {
		return nil, connect.NewError(connect.CodeUnavailable, "TURN credentials provider is not configured")
	}

	// Per-caller rate limiting.
	if ds.cacheSvc != nil && ds.rateLimitTURNPerMinute > 0 {
		callerID := "anonymous"
		claims := security.ClaimsFromContext(ctx)
		if claims != nil {
			if sub, err := claims.GetSubject(); err == nil && sub != "" {
				callerID = sub
			}
		}

		allowed, _ := ds.cacheSvc.CheckRateLimit(ctx, prefixRateTURN, callerID, ds.rateLimitTURNPerMinute)
		if !allowed {
			return nil, connect.NewError(
				connect.CodeResourceExhausted,
				"TURN credential rate limit exceeded",
			)
		}
	}

	credentials, err := ds.turnBusiness.GetTurnCredentials(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, "failed to generate TURN credentials")
	}

	expiresAt := time.Now().Unix() + int64(ds.turnTTL)

	var servers []*devicev1.TurnServer
	for _, ice := range credentials.ICEServers {
		for _, url := range ice.URLs {
			servers = append(servers, &devicev1.TurnServer{
				Url:        url,
				Username:   ice.Username,
				Credential: ice.Credential,
				ExpiresAt:  expiresAt,
			})
		}
	}

	return &devicev1.GetTurnCredentialsResponse{
		Servers:    servers,
		TtlSeconds: ds.turnTTL,
	}, nil
}
