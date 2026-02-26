package repository

import (
	"context"

	"github.com/pitabwire/frame/datastore"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

// Migrate runs database migrations for the geolocation service.
// GORM auto-migrate handles the non-spatial columns and indexes.
// PostGIS spatial columns (geom, bbox) and triggers are managed via SQL migration files
// because GORM does not understand PostGIS geometry types natively.
func Migrate(ctx context.Context, dbManager datastore.Manager, migrationPath string) error {
	dbPool := dbManager.GetPool(ctx, datastore.DefaultMigrationPoolName)

	return dbManager.Migrate(ctx, dbPool, migrationPath,
		&models.LocationPoint{},
		&models.Area{},
		&models.GeoEvent{},
		&models.GeofenceState{},
		&models.LatestPosition{},
		&models.Route{},
		&models.RouteAssignment{},
		&models.RouteDeviationState{},
	)
}
