package repository

import (
	"context"

	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
)

func Migrate(ctx context.Context, svc *frame.Service, migrationPath string) error {
	err := svc.MigrateDatastore(ctx, migrationPath,
		&models.Device{},
		&models.DeviceSession{},
		&models.DeviceKey{},
		&models.DeviceLog{},
	)
	if err != nil {
		return err
	}

	return nil
}
