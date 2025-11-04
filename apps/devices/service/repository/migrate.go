package repository

import (
	"context"
	"errors"

	"github.com/pitabwire/frame/datastore"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
)

func Migrate(ctx context.Context, dbManager datastore.Manager, migrationPath string) error {

	pool := dbManager.GetPool(ctx, datastore.DefaultMigrationPoolName)
	if pool == nil {
		return errors.New("datastore pool is not initialized")
	}

	return dbManager.Migrate(ctx, pool, migrationPath,
		&models.Device{},
		&models.DeviceSession{},
		&models.DeviceKey{},
		&models.DeviceLog{},
	)
}
