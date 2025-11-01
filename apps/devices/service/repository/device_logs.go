package repository

import (
	"context"
	"time"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/workerpool"
)

type deviceLogRepository struct {
	service *frame.Service
}

func NewDeviceLogRepository(service *frame.Service) DeviceLogRepository {
	return &deviceLogRepository{service: service}
}

func (dlr *deviceLogRepository) Save(ctx context.Context, log *models.DeviceLog) error {
	return dlr.service.DB(ctx, false).Create(log).Error
}

func (dlr *deviceLogRepository) GetByID(ctx context.Context, id string) (*models.DeviceLog, error) {
	var log models.DeviceLog
	if err := dlr.service.DB(ctx, true).First(&log, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (dlr *deviceLogRepository) GetByDeviceID(
	ctx context.Context,
	query *data.SearchQuery,
) (workerpool.JobResultPipe[[]*models.DeviceLog], error) {
	return data.StableSearch[models.DeviceLog](ctx, dlr.service, query, func(
		ctx context.Context,
		query *data.SearchQuery,
	) ([]*models.DeviceLog, error) {
		var deviceLogList []*models.DeviceLog

		paginator := query.Pagination

		db := dlr.service.DB(ctx, true).
			Limit(paginator.Limit).Offset(paginator.Offset)

		if query.Fields != nil {
			startAt, sok := query.Fields["start_date"]
			stopAt, stok := query.Fields["end_date"]
			if sok && startAt != nil && stok && stopAt != nil {
				startDate, ok1 := startAt.(*time.Time)
				endDate, ok2 := stopAt.(*time.Time)
				if ok1 && ok2 {
					db = db.Where(
						"created_at BETWEEN ? AND ? ",
						startDate.Format("2020-01-31T00:00:00Z"),
						endDate.Format("2020-01-31T00:00:00Z"),
					)
				}
			}

			deviceID, pok := query.Fields["device_id"]
			if pok {
				db = db.Where("device_id = ?", deviceID)
			}
		}

		if query.Query != "" {
			db = db.Where(" search_data @@ plainto_tsquery(?) ", query.Query)
		}

		err := db.Find(&deviceLogList).Error
		if err != nil {
			return nil, err
		}

		return deviceLogList, nil
	})
}
