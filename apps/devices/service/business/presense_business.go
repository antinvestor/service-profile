package business

import (
	"context"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
	"github.com/pitabwire/frame/queue"
	"github.com/pitabwire/frame/workerpool"
)

type PresenceBusiness interface {
	UpdatePresence(ctx context.Context, req *devicev1.UpdatePresenceRequest) error
}

type presenceBusiness struct {
	cfg *config.DevicesConfig

	qMan    queue.Manager
	workMan workerpool.Manager

	deviceRepo repository.DeviceRepository
}

// NewPresenceBusiness creates a new instance of PresenceBusiness.
func NewPresenceBusiness(_ context.Context, cfg *config.DevicesConfig, qMan queue.Manager,
	workMan workerpool.Manager, deviceRepo repository.DeviceRepository) PresenceBusiness {
	return &presenceBusiness{
		cfg:        cfg,
		qMan:       qMan,
		workMan:    workMan,
		deviceRepo: deviceRepo,
	}
}

func (p presenceBusiness) UpdatePresence(ctx context.Context, req *devicev1.UpdatePresenceRequest) error {
	return nil
}
