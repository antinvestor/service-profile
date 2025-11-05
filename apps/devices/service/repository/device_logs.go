package repository

import (
	"context"

	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
)

type deviceLogRepository struct {
	datastore.BaseRepository[*models.DeviceLog]
}

func NewDeviceLogRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) DeviceLogRepository {
	return &deviceLogRepository{
		BaseRepository: datastore.NewBaseRepository[*models.DeviceLog](
			ctx, dbPool, workMan, func() *models.DeviceLog { return &models.DeviceLog{} },
		),
	}
}

func (dlr *deviceLogRepository) GetByDeviceID(
	ctx context.Context,
	deviceID string,
) (workerpool.JobResultPipe[[]*models.DeviceLog], error) {
	query := data.NewSearchQuery(data.WithSearchFiltersAndByValue(map[string]any{
		"device_id": deviceID,
	}))

	return dlr.Search(ctx, query)
}
