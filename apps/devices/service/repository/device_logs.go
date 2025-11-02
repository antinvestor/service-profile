package repository

import (
	"context"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"
)

type deviceLogRepository struct {
	datastore.BaseRepository[*models.DeviceLog]
	service *frame.Service
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

	query := data.NewSearchQuery("", data.WithSearchFiltersAndByValue(map[string]any{
		"device_id": deviceID,
	}))

	return dlr.Search(ctx, query)
}
