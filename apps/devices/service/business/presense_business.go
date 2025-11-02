package business

import (
	"context"
	"errors"
	"slices"
	"time"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/queue"
	"github.com/pitabwire/frame/workerpool"
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
}

// NewPresenceBusiness creates a new instance of PresenceBusiness.
func NewPresenceBusiness(_ context.Context, cfg *config.DevicesConfig, qMan queue.Manager,
	workMan workerpool.Manager, deviceRepo repository.DeviceRepository,
	presenceRepo repository.DevicePresenceRepository) PresenceBusiness {
	return &presenceBusiness{
		cfg:          cfg,
		qMan:         qMan,
		workMan:      workMan,
		deviceRepo:   deviceRepo,
		presenceRepo: presenceRepo,
	}
}

func (p presenceBusiness) UpdatePresence(ctx context.Context, req *devicev1.UpdatePresenceRequest) (*devicev1.PresenceObject, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	device, err := p.deviceRepo.GetByID(ctx, req.GetDeviceId())
	if err != nil {
		return nil, err
	}

	extras := data.JSONMap{}
	extras = extras.FromProtoStruct(req.GetExtras())

	var expirtyTime *time.Time
	if slices.Contains([]devicev1.PresenceStatus{}, req.GetStatus()) {
		t := time.Now().Add(1 * time.Hour)
		expirtyTime = &t
	}

	presence := &models.DevicePresence{
		DeviceID:      device.GetID(),
		ProfileID:     device.ProfileID,
		Status:        req.GetStatus(),
		StatusMessage: req.GetStatusMessage(),
		ExpiryTime:    expirtyTime,
		Data:          extras,
	}

	err = p.presenceRepo.Create(ctx, presence)
	if err != nil {
		return nil, err
	}

	return presence.ToAPI(), nil
}
