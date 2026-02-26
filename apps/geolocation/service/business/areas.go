package business

import (
	"context"
	"errors"
	"fmt"
	"maps"

	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/events"
	"github.com/pitabwire/util"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
)

// AreaChangedEventName is the internal frame event name for area changes.
const AreaChangedEventName = "area.changed"

// StateActive corresponds to common.v1.STATE_ACTIVE = 2.
const StateActive int32 = 2

// StateDeleted corresponds to common.v1.STATE_DELETED = 4.
const StateDeleted int32 = 4

type areaBusiness struct {
	eventsMan events.Manager
	areaRepo  repository.AreaRepository
	stateRepo repository.GeofenceStateRepository
}

// NewAreaBusiness creates a new AreaBusiness.
func NewAreaBusiness(
	eventsMan events.Manager,
	areaRepo repository.AreaRepository,
	stateRepo repository.GeofenceStateRepository,
) AreaBusiness {
	return &areaBusiness{
		eventsMan: eventsMan,
		areaRepo:  areaRepo,
		stateRepo: stateRepo,
	}
}

func (b *areaBusiness) CreateArea(ctx context.Context, req *models.CreateAreaRequest) (*models.AreaAPI, error) {
	log := util.Log(ctx)

	if req == nil || req.Data == nil {
		return nil, errors.New("create area request data is nil")
	}

	apiData := req.Data

	// Validate.
	if err := models.ValidateAreaName(apiData.Name); err != nil {
		return nil, fmt.Errorf("invalid area name: %w", err)
	}
	if err := models.ValidateGeoJSON(apiData.GeometryJSON); err != nil {
		return nil, fmt.Errorf("invalid geometry: %w", err)
	}
	if apiData.OwnerID == "" {
		return nil, errors.New("owner_id is required")
	}

	area := &models.Area{
		OwnerID:      apiData.OwnerID,
		Name:         apiData.Name,
		Description:  apiData.Description,
		AreaType:     apiData.AreaType,
		GeometryJSON: apiData.GeometryJSON,
		State:        StateActive,
		Extras:       models.StructToJSONMap(apiData.Extras),
	}
	area.GenID(ctx)

	// Persist the area row and set geometry in a single transaction.
	// This prevents orphaned rows without spatial data on partial failure.
	db := b.areaRepo.Pool().DB(ctx, false)
	txErr := db.Transaction(func(tx *gorm.DB) error {
		if createErr := tx.Create(area).Error; createErr != nil {
			return fmt.Errorf("create area: %w", createErr)
		}
		if geomErr := b.areaRepo.UpdateGeometryTx(tx, area.GetID(), apiData.GeometryJSON); geomErr != nil {
			return fmt.Errorf("set area geometry: %w", geomErr)
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	// Re-read to get computed metrics (area_m2, perimeter_m).
	persisted, err := b.areaRepo.GetByID(ctx, area.GetID())
	if err != nil {
		return nil, fmt.Errorf("read back created area: %w", err)
	}

	b.emitAreaChanged(ctx, persisted.GetID(), persisted.OwnerID, "created")

	log.Info("area created", "area_id", persisted.GetID(), "name", persisted.Name)
	return persisted.ToAPI(), nil
}

func (b *areaBusiness) UpdateArea(ctx context.Context, req *models.UpdateAreaRequest) (*models.AreaAPI, error) {
	log := util.Log(ctx)

	if req == nil || req.ID == "" {
		return nil, errors.New("update area request requires an ID")
	}

	area, err := b.areaRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("area not found: %w", err)
	}

	// Apply field updates.
	if req.Name != "" {
		if vErr := models.ValidateAreaName(req.Name); vErr != nil {
			return nil, fmt.Errorf("invalid area name: %w", vErr)
		}
		area.Name = req.Name
	}
	if req.Description != "" {
		area.Description = req.Description
	}
	if req.AreaType != nil {
		area.AreaType = *req.AreaType
	}
	if req.Extras != nil {
		existing := models.StructToJSONMap(area.ToAPI().Extras)
		maps.Copy(existing, models.StructToJSONMap(req.Extras))
		area.Extras = existing
	}

	// Save non-spatial updates.
	if _, err = b.areaRepo.Update(ctx, area); err != nil {
		return nil, fmt.Errorf("update area: %w", err)
	}

	// Update geometry if provided.
	if req.Geometry != "" {
		if vErr := models.ValidateGeoJSON(req.Geometry); vErr != nil {
			return nil, fmt.Errorf("invalid geometry: %w", vErr)
		}
		if gErr := b.areaRepo.UpdateGeometry(ctx, area.GetID(), req.Geometry); gErr != nil {
			return nil, fmt.Errorf("update area geometry: %w", gErr)
		}
	}

	// Re-read to get updated computed metrics.
	persisted, err := b.areaRepo.GetByID(ctx, area.GetID())
	if err != nil {
		return nil, fmt.Errorf("read back updated area: %w", err)
	}

	b.emitAreaChanged(ctx, persisted.GetID(), persisted.OwnerID, "updated")

	log.Info("area updated", "area_id", persisted.GetID())
	return persisted.ToAPI(), nil
}

func (b *areaBusiness) DeleteArea(ctx context.Context, areaID string) error {
	log := util.Log(ctx)

	area, err := b.areaRepo.GetByID(ctx, areaID)
	if err != nil {
		return fmt.Errorf("area not found: %w", err)
	}

	// Soft delete: set state to DELETED.
	area.State = StateDeleted
	if _, err = b.areaRepo.Update(ctx, area); err != nil {
		return fmt.Errorf("soft delete area: %w", err)
	}

	// Clean up geofence states for the deleted area to prevent stale "inside" entries.
	if delErr := b.stateRepo.DeleteByArea(ctx, areaID); delErr != nil {
		log.WithError(delErr).Error("failed to clean up geofence states for deleted area",
			"area_id", areaID,
		)
		// Non-fatal: area is already marked deleted, spatial queries will skip it.
	}

	b.emitAreaChanged(ctx, area.GetID(), area.OwnerID, "deleted")

	log.Info("area deleted", "area_id", area.GetID())
	return nil
}

func (b *areaBusiness) GetArea(ctx context.Context, areaID string) (*models.AreaAPI, error) {
	area, err := b.areaRepo.GetByID(ctx, areaID)
	if err != nil {
		return nil, fmt.Errorf("get area: %w", err)
	}
	return area.ToAPI(), nil
}

func (b *areaBusiness) SearchAreas(
	ctx context.Context,
	query string,
	ownerID string,
	limit int,
) ([]*models.AreaAPI, error) {
	if limit <= 0 {
		limit = defaultSearchLimit
	}

	var areas []*models.Area
	var err error

	switch {
	case ownerID != "":
		areas, err = b.areaRepo.SearchByOwner(ctx, ownerID, limit)
	case query != "":
		areas, err = b.areaRepo.SearchByQuery(ctx, query, limit)
	default:
		return nil, errors.New("either query or owner_id is required for area search")
	}

	if err != nil {
		return nil, fmt.Errorf("search areas: %w", err)
	}

	result := make([]*models.AreaAPI, 0, len(areas))
	for _, a := range areas {
		result = append(result, a.ToAPI())
	}
	return result, nil
}

func (b *areaBusiness) emitAreaChanged(ctx context.Context, areaID, ownerID, action string) {
	event := &models.AreaChangedEvent{
		AreaID:  areaID,
		OwnerID: ownerID,
		Action:  action,
	}

	if err := b.eventsMan.Emit(ctx, AreaChangedEventName, event); err != nil {
		util.Log(ctx).WithError(err).Error("failed to emit area changed event",
			"area_id", areaID,
			"action", action,
		)
	}
}

// Ensure data package is used (for JSONMap operations).
var _ = data.JSONMap{}
