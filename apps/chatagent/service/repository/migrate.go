package repository

import (
	"context"

	"github.com/pitabwire/frame/v2/datastore"

	"github.com/antinvestor/service-profile/apps/chatagent/service/models"
)

func Migrate(ctx context.Context, dbManager datastore.Manager, migrationPath string) error {
	dbPool := dbManager.GetPool(ctx, datastore.DefaultMigrationPoolName)
	return dbManager.Migrate(ctx, dbPool, migrationPath,
		&models.IntakeContext{},
		&models.Session{},
		&models.Message{},
	)
}
