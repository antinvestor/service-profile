package business

import (
	"context"
	"encoding/json"
	"errors"
	"slices"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/v2/data"
	"github.com/pitabwire/frame/v2/queue"
	"github.com/pitabwire/frame/v2/workerpool"
	"golang.org/x/sync/singleflight"

	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/caching"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
	"github.com/antinvestor/service-profile/internal/outbox"
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
	device, err := b.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	// Key material is checked against the schema for its type before it is
	// stored, and the extra document the service keeps is the normalised one
	// (GFOS K4).
	normalisedExtra, err := ValidateKeyMaterial(keyType, key, extra)
	if err != nil {
		return nil, err
	}

	deviceKey := &models.DeviceKey{
		DeviceID: deviceID,
		KeyType:  keyType,
		Key:      key,
		Extra:    normalisedExtra,
	}

	// The key row and the device.key.added fact are written together, so a
	// registered key always has a fact and a fact always has a key (GFOS K5).
	err = b.deviceKeyRepo.CreateWithFact(ctx, deviceKey, func(stored *models.DeviceKey) *outbox.Event {
		return DeviceKeyAddedFact(ctx, stored, device.ProfileID)
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateKeyMaterial) {
			return b.resolveDuplicateKey(ctx, deviceID, keyType, key)
		}
		return nil, err
	}

	// Invalidate keys cache after adding a key.
	if b.cache != nil {
		b.cache.InvalidateDeviceKeys(ctx, deviceID)
	}

	return deviceKey.ToAPI(), nil
}

// resolveDuplicateKey decides what a (key_type, key) collision means. The same
// device re-registering the same material is the client retrying, so the
// existing key is returned and no second fact is published. A different device
// claiming material that is already registered is refused: one public key
// belongs to one device.
func (b *keysBusiness) resolveDuplicateKey(
	ctx context.Context,
	deviceID string,
	keyType devicev1.KeyType,
	key []byte,
) (*devicev1.KeyObject, error) {
	existing, err := b.deviceKeyRepo.GetByTypeAndKey(ctx, keyType, key)
	if err != nil {
		return nil, err
	}
	if existing.DeviceID == deviceID {
		return existing.ToAPI(), nil
	}
	return nil, connect.NewError(
		connect.CodeAlreadyExists,
		errors.New("key material is already registered to another device"),
	)
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
func (b *keysBusiness) getDeviceKeysWithCache(ctx context.Context, deviceID string) ([]*models.DeviceKey, error) {
	if keys, hit := b.tryKeysCache(ctx, deviceID); hit {
		return keys, nil
	}

	// Use singleflight to collapse concurrent DB fetches for the same device.
	val, err, _ := b.sfKeys.Do("keys:"+deviceID, func() (any, error) {
		return b.fetchKeysAndCache(ctx, deviceID)
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

// tryKeysCache attempts to read device keys from cache.
func (b *keysBusiness) tryKeysCache(ctx context.Context, deviceID string) ([]*models.DeviceKey, bool) {
	if b.cache == nil {
		return nil, false
	}
	cached, found := b.cache.GetDeviceKeys(ctx, deviceID)
	if !found {
		return nil, false
	}
	var entry cachedKeyEntry
	if err := json.Unmarshal(cached, &entry); err != nil {
		return nil, false
	}
	return entry.Keys, true
}

// fetchKeysAndCache loads device keys from DB, populates cache, and returns them.
func (b *keysBusiness) fetchKeysAndCache(ctx context.Context, deviceID string) ([]*models.DeviceKey, error) {
	// Double-check cache after acquiring singleflight.
	if keys, hit := b.tryKeysCache(ctx, deviceID); hit {
		return keys, nil
	}

	keys, err := b.deviceKeyRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	if b.cache != nil {
		entry := cachedKeyEntry{DeviceID: deviceID, Keys: keys}
		encoded, encErr := json.Marshal(entry)
		if encErr == nil {
			b.cache.SetDeviceKeys(ctx, deviceID, encoded)
		}
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
				removedKey, err := b.removeKey(ctx, keyID)
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

// removeKey withdraws one key and stages the device.key.removed fact in the
// same transaction. The owning profile is read first so the fact can name it
// without a join the consumer would otherwise have to do.
func (b *keysBusiness) removeKey(ctx context.Context, keyID string) (*models.DeviceKey, error) {
	existing, err := b.deviceKeyRepo.GetByID(ctx, keyID)
	if err != nil {
		return nil, err
	}

	var profileID string
	if device, devErr := b.deviceRepo.GetByID(ctx, existing.DeviceID); devErr == nil {
		profileID = device.ProfileID
	}

	return b.deviceKeyRepo.RemoveByIDWithFact(ctx, keyID, func(removed *models.DeviceKey) *outbox.Event {
		return DeviceKeyRemovedFact(ctx, removed, profileID)
	})
}
