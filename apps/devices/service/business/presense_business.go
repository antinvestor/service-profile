package business

import (
	"context"
	"errors"
	"slices"
	"time"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/queue"
	"github.com/pitabwire/frame/workerpool"
	"github.com/pitabwire/util"
	"go.opentelemetry.io/otel/attribute"

	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/caching"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

// presenceExpiryStatuses lists statuses that should auto-expire.
// ONLINE and AWAY expire after the default TTL; OFFLINE and DO_NOT_DISTURB persist until changed.
//
//nolint:gochecknoglobals // package-level lookup table for presence expiry logic
var presenceExpiryStatuses = []devicev1.PresenceStatus{
	devicev1.PresenceStatus_ONLINE,
	devicev1.PresenceStatus_AWAY,
}

const (
	defaultPresenceExpiry   = 1 * time.Hour
	offlinePresenceCacheTTL = 5 * time.Minute
)

type PresenceBusiness interface {
	UpdatePresence(ctx context.Context, req *devicev1.UpdatePresenceRequest) (*devicev1.PresenceObject, error)
}

type presenceBusiness struct {
	cfg *config.DevicesConfig

	qMan    queue.Manager
	workMan workerpool.Manager

	deviceRepo   repository.DeviceRepository
	presenceRepo repository.DevicePresenceRepository
	cache        *caching.DeviceCacheService
}

// NewPresenceBusiness creates a new instance of PresenceBusiness.
func NewPresenceBusiness(_ context.Context, cfg *config.DevicesConfig, qMan queue.Manager,
	workMan workerpool.Manager, deviceRepo repository.DeviceRepository,
	presenceRepo repository.DevicePresenceRepository,
	cacheSvc *caching.DeviceCacheService) PresenceBusiness {
	return &presenceBusiness{
		cfg:          cfg,
		qMan:         qMan,
		workMan:      workMan,
		deviceRepo:   deviceRepo,
		presenceRepo: presenceRepo,
		cache:        cacheSvc,
	}
}

func (p *presenceBusiness) UpdatePresence(
	ctx context.Context,
	req *devicev1.UpdatePresenceRequest,
) (*devicev1.PresenceObject, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	deviceID := req.GetDeviceId()
	if deviceID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("device ID is required"))
	}

	ctx, span := caching.StartSpan(ctx, "UpdatePresence",
		attribute.String("device_id", deviceID),
		attribute.String("status", req.GetStatus().String()))
	defer caching.EndSpan(ctx, span, nil)

	// Rate limit presence updates per device.
	if p.cache != nil && p.cfg.RateLimitPresencePerMinute > 0 {
		allowed, count := p.cache.CheckPresenceRateLimit(ctx, deviceID, p.cfg.RateLimitPresencePerMinute)
		if !allowed {
			caching.RecordRateLimited(ctx, "update_presence")
			util.Log(ctx).WithField("device_id", deviceID).WithField("count", count).
				Warn("presence update rate limited")
			return nil, connect.NewError(connect.CodeResourceExhausted,
				errors.New("presence update rate limit exceeded"))
		}
	}

	device, err := p.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	extras := data.JSONMap{}
	extras = extras.FromProtoStruct(req.GetExtras())

	// Compute expiry: statuses that imply active use get a TTL so they
	// auto-expire if the device disappears without sending OFFLINE.
	var expiryTime *time.Time
	cacheTTL := caching.TTLPresence
	if slices.Contains(presenceExpiryStatuses, req.GetStatus()) {
		t := time.Now().Add(defaultPresenceExpiry)
		expiryTime = &t
		cacheTTL = defaultPresenceExpiry
	}

	// For OFFLINE, use a short cache TTL so queries stop returning stale state.
	if req.GetStatus() == devicev1.PresenceStatus_OFFLINE {
		cacheTTL = offlinePresenceCacheTTL
	}

	now := time.Now().UTC()

	presence := &models.DevicePresence{
		DeviceID:      device.GetID(),
		ProfileID:     device.ProfileID,
		Status:        req.GetStatus(),
		StatusMessage: req.GetStatusMessage(),
		ExpiryTime:    expiryTime,
		Data:          extras,
	}

	// Write to DB for history.
	if err = p.presenceRepo.Create(ctx, presence); err != nil {
		return nil, err
	}

	// Write to cache for fast reads of current state.
	if p.cache != nil {
		expiryStr := ""
		if expiryTime != nil {
			expiryStr = expiryTime.UTC().Format(time.RFC3339)
		}
		entry := &caching.PresenceEntry{
			DeviceID:      device.GetID(),
			ProfileID:     device.ProfileID,
			Status:        int32(req.GetStatus()),
			StatusMessage: req.GetStatusMessage(),
			ExpiryTime:    expiryStr,
			LastActive:    now.Format(time.RFC3339),
		}
		p.cache.SetPresence(ctx, device.GetID(), entry, cacheTTL)
	}

	return presence.ToAPI(), nil
}
