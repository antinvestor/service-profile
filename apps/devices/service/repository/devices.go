package repository

import (
	"context"
	"errors"

	"github.com/pitabwire/frame/v2/datastore"
	"github.com/pitabwire/frame/v2/datastore/pool"
	"github.com/pitabwire/frame/v2/workerpool"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/internal/outbox"
)

// ErrDeviceLinkConflict is returned when a device link loses a race with
// another writer, so neither the link nor its fact is written.
var ErrDeviceLinkConflict = errors.New("device link conflicted with a concurrent update")

type deviceRepository struct {
	datastore.BaseRepository[*models.Device]
}

func NewDeviceRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) DeviceRepository {
	return &deviceRepository{
		BaseRepository: datastore.NewBaseRepository[*models.Device](
			ctx, dbPool, workMan, func() *models.Device { return &models.Device{} },
		),
	}
}

func (dr *deviceRepository) RemoveByID(ctx context.Context, id string) (*models.Device, error) {
	device, err := dr.BaseRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	err = dr.Pool().DB(ctx, false).Delete(device).Error
	if err != nil {
		return nil, err
	}
	return device, nil
}

// LinkProfileWithFact binds a device to a profile and stages the accompanying
// fact in the same transaction, so a linked device and the fact that says so
// can never disagree (GFOS K5).
func (dr *deviceRepository) LinkProfileWithFact(
	ctx context.Context,
	device *models.Device,
	fact *outbox.Event,
) error {
	return dr.Pool().DB(ctx, false).Transaction(func(tx *gorm.DB) error {
		linked := tx.Model(device).
			Where("id = ? AND version = ?", device.GetID(), device.GetVersion()).
			Select("profile_id").
			Updates(device)
		if linked.Error != nil {
			return linked.Error
		}
		if linked.RowsAffected != 1 {
			// Another writer linked the device first; the fact would not be
			// true, so nothing is written.
			return ErrDeviceLinkConflict
		}
		return outbox.Enqueue(tx, fact)
	})
}
