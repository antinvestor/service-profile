package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

type areaRepository struct {
	datastore.BaseRepository[*models.Area]
}

// NewAreaRepository creates a new repository for geographic areas.
func NewAreaRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) AreaRepository {
	return &areaRepository{
		BaseRepository: datastore.NewBaseRepository[*models.Area](
			ctx, dbPool, workMan, func() *models.Area { return &models.Area{} },
		),
	}
}

// GetActiveByBoundingBox returns active areas whose bounding box intersects with a point.
// This is the fast pre-filter step for geofence detection.
// The geom and bbox columns are managed by PostGIS; we use raw SQL for spatial queries.
func (r *areaRepository) GetActiveByBoundingBox(ctx context.Context, lat, lon float64) ([]*models.Area, error) {
	var areas []*models.Area

	db := r.Pool().DB(ctx, true)
	// ST_Intersects with the bounding box column for speed, filtering only active non-deleted areas.
	// State 2 corresponds to common.v1.STATE_ACTIVE.
	result := db.Where(
		"deleted_at IS NULL AND state = 2 AND bbox IS NOT NULL AND ST_Intersects(bbox, ST_SetSRID(ST_Point(?, ?), 4326))",
		lon,
		lat, // PostGIS uses (lon, lat) = (x, y) ordering
	).Find(&areas)

	if result.Error != nil {
		return nil, fmt.Errorf("get active areas by bbox at (%f, %f): %w", lat, lon, result.Error)
	}
	return areas, nil
}

// ContainsPoint checks if a specific area's actual geometry contains the given point.
// This is the precise containment test done after bbox pre-filtering.
func (r *areaRepository) ContainsPoint(ctx context.Context, areaID string, lat, lon float64) (bool, error) {
	var contains bool

	db := r.Pool().DB(ctx, true)
	result := db.Raw(
		"SELECT ST_Contains(geom, ST_SetSRID(ST_Point(?, ?), 4326)) FROM areas WHERE id = ? AND deleted_at IS NULL",
		lon, lat, areaID,
	).Scan(&contains)

	if result.Error != nil {
		return false, fmt.Errorf("contains point check for area %s: %w", areaID, result.Error)
	}
	return contains, nil
}

// UpdateGeometry sets the PostGIS geom column from GeoJSON.
// This triggers the compute_area_metrics trigger that updates bbox, area_m2, perimeter_m.
// It also updates the geometry_json column for API reads.
func (r *areaRepository) UpdateGeometry(ctx context.Context, areaID string, geoJSON string) error {
	db := r.Pool().DB(ctx, false)
	return executeUpdateGeometry(db, areaID, geoJSON)
}

// UpdateGeometryTx sets the PostGIS geom column within an existing transaction.
func (r *areaRepository) UpdateGeometryTx(tx *gorm.DB, areaID string, geoJSON string) error {
	return executeUpdateGeometry(tx, areaID, geoJSON)
}

// executeUpdateGeometry is the shared implementation for geometry updates.
func executeUpdateGeometry(db *gorm.DB, areaID string, geoJSON string) error {
	result := db.Exec(
		`UPDATE areas
		 SET geom = ST_SetSRID(ST_GeomFromGeoJSON(?), 4326),
		     geometry_json = ?,
		     modified_at = NOW()
		 WHERE id = ? AND deleted_at IS NULL`,
		geoJSON, geoJSON, areaID,
	)

	if result.Error != nil {
		return fmt.Errorf("update geometry for area %s: %w", areaID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("area %s not found or deleted", areaID)
	}
	return nil
}

// GetNearbyAreas finds areas within radiusMeters of the given point.
// Uses PostGIS ST_DWithin with geography cast for accurate meter-based distance.
func (r *areaRepository) GetNearbyAreas(
	ctx context.Context,
	lat, lon, radiusMeters float64,
	limit int,
) ([]*AreaWithDistance, error) {
	db := r.Pool().DB(ctx, true)

	type rawResult struct {
		models.Area
		DistanceMeters float64 `gorm:"column:distance_meters"`
	}

	var results []rawResult

	query := db.Raw(
		`SELECT a.*,
		        ST_Distance(a.geom::geography, ST_SetSRID(ST_Point(?, ?), 4326)::geography) AS distance_meters
		 FROM areas a
		 WHERE a.deleted_at IS NULL
		   AND a.state = 2
		   AND a.geom IS NOT NULL
		   AND ST_DWithin(a.geom::geography, ST_SetSRID(ST_Point(?, ?), 4326)::geography, ?)
		 ORDER BY distance_meters ASC
		 LIMIT ?`,
		lon, lat, lon, lat, radiusMeters, limit,
	)

	if err := query.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("get nearby areas at (%f, %f) radius %f: %w", lat, lon, radiusMeters, err)
	}

	out := make([]*AreaWithDistance, len(results))
	for i := range results {
		area := results[i].Area
		out[i] = &AreaWithDistance{
			Area:           &area,
			DistanceMeters: results[i].DistanceMeters,
		}
	}
	return out, nil
}

// SearchByOwner returns areas owned by the given owner, with a limit.
func (r *areaRepository) SearchByOwner(ctx context.Context, ownerID string, limit int) ([]*models.Area, error) {
	var areas []*models.Area
	db := r.Pool().DB(ctx, true)
	query := db.Where("owner_id = ? AND deleted_at IS NULL", ownerID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	result := query.Find(&areas)
	if result.Error != nil {
		return nil, fmt.Errorf("search areas by owner %s: %w", ownerID, result.Error)
	}
	return areas, nil
}

// SearchByQuery performs text search on area name and description.
// SQL wildcards in user input are escaped to prevent wildcard injection.
func (r *areaRepository) SearchByQuery(ctx context.Context, query string, limit int) ([]*models.Area, error) {
	var areas []*models.Area
	db := r.Pool().DB(ctx, true)
	likeQuery := "%" + escapeLikeWildcards(query) + "%"
	result := db.Where(
		"deleted_at IS NULL AND (name ILIKE ? OR description ILIKE ?)",
		likeQuery, likeQuery,
	).Order("created_at DESC").Limit(limit).Find(&areas)
	if result.Error != nil {
		return nil, fmt.Errorf("search areas by query %q: %w", query, result.Error)
	}
	return areas, nil
}

// escapeLikeWildcards escapes SQL LIKE/ILIKE special characters in user input.
// This prevents users from using % or _ as wildcards to enumerate data.
func escapeLikeWildcards(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}
