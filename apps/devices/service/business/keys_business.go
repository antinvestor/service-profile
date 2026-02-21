package business

import (
	"context"
	"encoding/json"
	"errors"
	"slices"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/queue"
	"github.com/pitabwire/frame/workerpool"
	"golang.org/x/sync/singleflight"

	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/caching"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
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
		keys ...devicev1.KeyType,
	) (<-chan workerpool.JobResult[[]*devicev1.KeyObject], error)
	RemoveKeys(ctx context.Context, id ...string) (<-chan workerpool.JobResult[[]*devicev1.KeyObject], error)
}

type keysBusiness struct {
	cfg *config.DevicesConfig

	qMan    queue.Manager
	workMan workerpool.Manager

	deviceRepo    repository.DeviceRepository
	deviceKeyRepo repository.DeviceKeyRepository

	cache  *caching.DeviceCacheService
	sfKeys singleflight.Group
}

// NewKeysBusiness creates a new instance of KeysBusiness.
func NewKeysBusiness(_ context.Context, cfg *config.DevicesConfig,
	qMan queue.Manager, workMan workerpool.Manager, deviceRepo repository.DeviceRepository,
	deviceKeyRepo repository.DeviceKeyRepository, cacheSvc *caching.DeviceCacheService) KeysBusiness {
	return &keysBusiness{
		cfg:           cfg,
		qMan:          qMan,
		workMan:       workMan,
		deviceRepo:    deviceRepo,
		deviceKeyRepo: deviceKeyRepo,
		cache:         cacheSvc,
	}
}

func (b *keysBusiness) AddKey(
	ctx context.Context,
	deviceID string,
	keyType devicev1.KeyType,
	key []byte,
	extra data.JSONMap,
) (*devicev1.KeyObject, error) {
	// Validate that the device exists before adding a key.
	_, err := b.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	deviceKey := &models.DeviceKey{
		DeviceID: deviceID,
		KeyType:  keyType,
		Key:      key,
		Extra:    extra,
	}

	err = b.deviceKeyRepo.Create(ctx, deviceKey)
	if err != nil {
		return nil, err
	}

	// Invalidate keys cache after adding a key.
	if b.cache != nil {
		b.cache.InvalidateDeviceKeys(ctx, deviceID)
	}

	return deviceKey.ToAPI(), nil
}

// cachedKeyEntry is a serializable container for device keys stored in cache.
type cachedKeyEntry struct {
	DeviceID string              `json:"device_id"`
	Keys     []*models.DeviceKey `json:"keys"`
}

func (b *keysBusiness) GetKeys(
	ctx context.Context,
	deviceID string,
	keyType ...devicev1.KeyType,
) (<-chan workerpool.JobResult[[]*devicev1.KeyObject], error) {
	resultPipe := workerpool.NewJob[[]*devicev1.KeyObject](
		func(ctx context.Context, result workerpool.JobResultPipe[[]*devicev1.KeyObject]) error {
			keys, err := b.getDeviceKeysWithCache(ctx, deviceID)
			if err != nil {
				return err
			}

			apiKeys := make([]*devicev1.KeyObject, 0, len(keys))
			for _, key := range keys {
				if len(keyType) == 0 || slices.Contains(keyType, key.KeyType) {
					apiKeys = append(apiKeys, key.ToAPI())
				}
			}

			return result.WriteResult(ctx, apiKeys)
		},
	)

	if err := workerpool.SubmitJob(ctx, b.workMan, resultPipe); err != nil {
		return nil, err
	}

	return resultPipe.ResultChan(), nil
}

// getDeviceKeysWithCache retrieves device keys using cache-aside pattern with singleflight
// to collapse concurrent requests for the same device's keys.
//
//nolint:gocognit // Complexity from singleflight double-check pattern is intentional.
func (b *keysBusiness) getDeviceKeysWithCache(ctx context.Context, deviceID string) ([]*models.DeviceKey, error) {
	// Try cache first.
	if b.cache != nil {
		if cached, found := b.cache.GetDeviceKeys(ctx, deviceID); found {
			var entry cachedKeyEntry
			if err := json.Unmarshal(cached, &entry); err == nil {
				return entry.Keys, nil
			}
		}
	}

	// Use singleflight to collapse concurrent DB fetches for the same device.
	val, err, _ := b.sfKeys.Do("keys:"+deviceID, func() (any, error) {
		// Double-check cache after acquiring singleflight.
		if b.cache != nil {
			if cached, found := b.cache.GetDeviceKeys(ctx, deviceID); found {
				var entry cachedKeyEntry
				if err := json.Unmarshal(cached, &entry); err == nil {
					return entry.Keys, nil
				}
			}
		}

		keys, err := b.deviceKeyRepo.GetByDeviceID(ctx, deviceID)
		if err != nil {
			return nil, err
		}

		// Populate cache.
		if b.cache != nil {
			entry := cachedKeyEntry{DeviceID: deviceID, Keys: keys}
			encoded, encErr := json.Marshal(entry)
			if encErr == nil {
				b.cache.SetDeviceKeys(ctx, deviceID, encoded)
			}
		}

		return keys, nil
	})
	if err != nil {
		return nil, err
	}

	keys, ok := val.([]*models.DeviceKey)
	if !ok {
		return nil, errors.New("unexpected type in singleflight result")
	}
	return keys, nil
}

func (b *keysBusiness) RemoveKeys(
	ctx context.Context,
	id ...string,
) (<-chan workerpool.JobResult[[]*devicev1.KeyObject], error) {
	resultPipe := workerpool.NewJob[[]*devicev1.KeyObject](
		func(ctx context.Context, result workerpool.JobResultPipe[[]*devicev1.KeyObject]) error {
			var removedKeys []*devicev1.KeyObject

			for _, keyID := range id {
				removedKey, err := b.deviceKeyRepo.RemoveByID(ctx, keyID)
				if err != nil {
					return err
				}
				if removedKey != nil {
					removedKeys = append(removedKeys, removedKey.ToAPI())

					// Invalidate keys cache for the affected device.
					if b.cache != nil {
						b.cache.InvalidateDeviceKeys(ctx, removedKey.DeviceID)
					}
				}
			}

			return result.WriteResult(ctx, removedKeys)
		},
	)

	if err := workerpool.SubmitJob(ctx, b.workMan, resultPipe); err != nil {
		return nil, err
	}

	return resultPipe.ResultChan(), nil
}
