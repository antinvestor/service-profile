package repository

import (
	"context"
	"errors"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
)

func Migrate(ctx context.Context, svc *frame.Service, migrationPath string) error {
	manager := svc.DatastoreManager()
	if manager == nil {
		return errors.New("datastore manager is not initialized")
	}

	pool := manager.GetPool(ctx, datastore.DefaultPoolName)
	if pool == nil {
		return errors.New("datastore pool is not initialized")
	}

	return manager.Migrate(ctx, pool, migrationPath,
		&models.Device{},
		&models.DeviceSession{},
		&models.DeviceKey{},
		&models.DeviceLog{},
	)
}
