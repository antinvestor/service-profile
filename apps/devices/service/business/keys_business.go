package business

import (
	"context"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/queue"
	"github.com/pitabwire/frame/workerpool"
)

type KeysBusiness interface {
	AddKey(
		ctx context.Context,
		deviceID string,
		_ devicev1.KeyType,
		key []byte,
		extra data.JSONMap,
	) (*devicev1.KeyObject, error)
	GetKeys(
		ctx context.Context,
		deviceID string,
		_ []devicev1.KeyType,
	) (<-chan workerpool.JobResult[[]*devicev1.KeyObject], error)
	RemoveKeys(ctx context.Context, id ...string) (<-chan workerpool.JobResult[[]*devicev1.KeyObject], error)
}

type keysBusiness struct {
	cfg *config.DevicesConfig

	qMan    queue.Manager
	workMan workerpool.Manager

	deviceRepo    repository.DeviceRepository
	deviceKeyRepo repository.DeviceKeyRepository
}

// NewKeysBusiness creates a new instance of DeviceBusiness.
func NewKeysBusiness(_ context.Context, cfg *config.DevicesConfig,
	qMan queue.Manager, workMan workerpool.Manager, deviceRepo repository.DeviceRepository,
	deviceKeyRepo repository.DeviceKeyRepository) KeysBusiness {
	return &keysBusiness{
		cfg:           cfg,
		qMan:          qMan,
		workMan:       workMan,
		deviceRepo:    deviceRepo,
		deviceKeyRepo: deviceKeyRepo,
	}
}

func (b *keysBusiness) AddKey(
	ctx context.Context,
	deviceID string,
	_ devicev1.KeyType,
	key []byte,
	extra data.JSONMap,
) (*devicev1.KeyObject, error) {
	// Validate that the device exists before adding a key
	_, err := b.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	deviceKey := &models.DeviceKey{
		DeviceID: deviceID,
		Key:      key,
		Extra:    extra,
	}

	err = b.deviceKeyRepo.Create(ctx, deviceKey)
	if err != nil {
		return nil, err
	}

	return deviceKey.ToAPI(), nil
}

func (b *keysBusiness) GetKeys(
	ctx context.Context,
	deviceID string,
	_ []devicev1.KeyType,
) (<-chan workerpool.JobResult[[]*devicev1.KeyObject], error) {
	out := make(chan workerpool.JobResult[[]*devicev1.KeyObject])

	go func() {
		defer close(out)

		keys, err := b.deviceKeyRepo.GetByDeviceID(ctx, deviceID)
		if err != nil {
			out <- workerpool.ErrorResult[[]*devicev1.KeyObject](err)
			return
		}

		apiKeys := make([]*devicev1.KeyObject, len(keys))
		for i, key := range keys {
			apiKeys[i] = key.ToAPI()
		}

		out <- workerpool.Result[[]*devicev1.KeyObject](apiKeys)
	}()

	return out, nil
}

func (b *keysBusiness) RemoveKeys(
	ctx context.Context,
	id ...string,
) (<-chan workerpool.JobResult[[]*devicev1.KeyObject], error) {
	out := make(chan workerpool.JobResult[[]*devicev1.KeyObject])

	go func() {
		defer close(out)

		var removedKeys []*devicev1.KeyObject

		for _, keyID := range id {
			removedKey, err := b.deviceKeyRepo.RemoveByID(ctx, keyID)
			if err != nil {
				out <- workerpool.ErrorResult[[]*devicev1.KeyObject](err)
				return
			}
			if removedKey != nil {
				removedKeys = append(removedKeys, removedKey.ToAPI())
			}
		}

		out <- workerpool.Result[[]*devicev1.KeyObject](removedKeys)
	}()

	return out, nil
}
