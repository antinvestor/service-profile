package repository

import (
	"context"
	"time"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/data"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
)

type deviceRepository struct {
	service *frame.Service
}

func NewDeviceRepository(service *frame.Service) DeviceRepository {
	return &deviceRepository{service: service}
}

func (dr *deviceRepository) Save(ctx context.Context, device *models.Device) error {
	return dr.service.DB(ctx, false).Save(device).Error
}

func (dr *deviceRepository) GetByID(ctx context.Context, id string) (*models.Device, error) {
	var device models.Device
	if err := dr.service.DB(ctx, true).First(&device, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (dr *deviceRepository) Search(ctx context.Context,
	query *data.SearchQuery) (workerpool.JobResultPipe[[]*models.Device], error) {
	return data.StableSearch[models.Device](ctx, dr.service, query, func(
		ctx context.Context,
		query *data.SearchQuery,
	) ([]*models.Device, error) {
		var deviceList []*models.Device

		paginator := query.Pagination
		db := dr.service.DB(ctx, true).
			Limit(paginator.Limit).Offset(paginator.Offset)

		db = dr.applyFieldFilters(db, query.Fields)
		db = dr.applyTextSearch(db, query.Query)

		err := db.Find(&deviceList).Error
		if err != nil {
			return nil, err
		}

		return deviceList, nil
	})
}

func (dr *deviceRepository) applyFieldFilters(db *gorm.DB, fields map[string]interface{}) *gorm.DB {
	if fields == nil {
		return db
	}

	db = dr.applyDateRangeFilter(db, fields)
	db = dr.applyProfileIDFilter(db, fields)

	return db
}

func (dr *deviceRepository) applyDateRangeFilter(db *gorm.DB, fields map[string]interface{}) *gorm.DB {
	startAt, sok := fields["start_date"]
	stopAt, stok := fields["end_date"]

	if !sok || startAt == nil || !stok || stopAt == nil {
		return db
	}

	startDate, ok1 := startAt.(*time.Time)
	endDate, ok2 := stopAt.(*time.Time)

	if ok1 && ok2 {
		return db.Where(
			"created_at BETWEEN ? AND ? ",
			startDate.Format("2020-01-31T00:00:00Z"),
			endDate.Format("2020-01-31T00:00:00Z"),
		)
	}

	return db
}

func (dr *deviceRepository) applyProfileIDFilter(db *gorm.DB, fields map[string]interface{}) *gorm.DB {
	profileID, pok := fields["profile_id"]
	if !pok {
		return db
	}

	if profileID != "" {
		return db.Where("profile_id = ?", profileID)
	}

	// If profile_id is explicitly empty, return no results
	// This ensures proper test isolation when no claims are present
	return db.Where("1 = 0")
}

func (dr *deviceRepository) applyTextSearch(db *gorm.DB, searchQuery string) *gorm.DB {
	if searchQuery == "" {
		return db
	}

	likeQuery := "%" + searchQuery + "%"
	return db.Where(" name ilike ? OR os ilike ?", likeQuery, likeQuery)
}

func (dr *deviceRepository) RemoveByID(ctx context.Context, id string) (*models.Device, error) {
	var device models.Device
	if err := dr.service.DB(ctx, true).First(&device, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := dr.service.DB(ctx, false).Delete(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}
