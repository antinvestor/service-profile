package repository

import (
	"context"
	"errors"
	"time"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pitabwire/frame/v2/datastore"
	"github.com/pitabwire/frame/v2/datastore/pool"
	"github.com/pitabwire/frame/v2/workerpool"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/internal/outbox"
)

// ErrDuplicateKeyMaterial is returned when a key with the same (key_type, key)
// is already registered. The database enforces this with a unique index; this
// error is what that violation looks like to the business layer.
var ErrDuplicateKeyMaterial = errors.New("key material is already registered")

// pgUniqueViolation is the SQLSTATE Postgres raises for a unique index
// violation.
const pgUniqueViolation = "23505"

type deviceKeyRepository struct {
	datastore.BaseRepository[*models.DeviceKey]
}

func NewDeviceKeyRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) DeviceKeyRepository {
	return &deviceKeyRepository{
		BaseRepository: datastore.NewBaseRepository[*models.DeviceKey](
			ctx, dbPool, workMan, func() *models.DeviceKey { return &models.DeviceKey{} },
		),
	}
}

func (r *deviceKeyRepository) GetByDeviceID(ctx context.Context, deviceID string) ([]*models.DeviceKey, error) {
	var keys []*models.DeviceKey
	err := r.Pool().DB(ctx, true).
		Where("device_id = ?", deviceID).
		Where("expires_at IS NULL OR expires_at > ?", time.Now()).
		Find(&keys).Error
	return keys, err
}

// GetByTypeAndKey finds the live registration of a piece of key material.
func (r *deviceKeyRepository) GetByTypeAndKey(
	ctx context.Context,
	keyType devicev1.KeyType,
	key []byte,
) (*models.DeviceKey, error) {
	var found models.DeviceKey
	err := r.Pool().DB(ctx, true).
		Where("key_type = ?", keyType).
		Where("key = ?", key).
		First(&found).Error
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// CreateWithFact registers a key and stages the fact describing it in the same
// transaction (GFOS K5). The fact is built from the persisted row so it
// carries the generated id and timestamps.
//
// A (key_type, key) collision surfaces as ErrDuplicateKeyMaterial: the
// uniqueness GFOS K4 requires is enforced by the database, not by a read
// before the write, so two concurrent registrations of the same public key
// cannot both succeed.
func (r *deviceKeyRepository) CreateWithFact(
	ctx context.Context,
	key *models.DeviceKey,
	fact func(*models.DeviceKey) *outbox.Event,
) error {
	key.GenID(ctx)

	err := r.Pool().DB(ctx, false).Transaction(func(tx *gorm.DB) error {
		if createErr := tx.Create(key).Error; createErr != nil {
			return createErr
		}
		return outbox.Enqueue(tx, fact(key))
	})
	if isUniqueViolation(err) {
		return ErrDuplicateKeyMaterial
	}
	return err
}

func (r *deviceKeyRepository) RemoveByID(ctx context.Context, id string) (*models.DeviceKey, error) {
	var key models.DeviceKey
	if err := r.Pool().DB(ctx, true).First(&key, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.Pool().DB(ctx, false).Delete(&key).Error; err != nil {
		return nil, err
	}
	return &key, nil
}

// RemoveByIDWithFact withdraws a key and stages the fact describing the
// withdrawal in the same transaction.
func (r *deviceKeyRepository) RemoveByIDWithFact(
	ctx context.Context,
	id string,
	fact func(*models.DeviceKey) *outbox.Event,
) (*models.DeviceKey, error) {
	var key models.DeviceKey

	err := r.Pool().DB(ctx, false).Transaction(func(tx *gorm.DB) error {
		if findErr := tx.First(&key, "id = ?", id).Error; findErr != nil {
			return findErr
		}
		if delErr := tx.Delete(&key).Error; delErr != nil {
			return delErr
		}
		return outbox.Enqueue(tx, fact(&key))
	})
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgUniqueViolation
	}
	return false
}
