package repository

import (
	"context"

	"github.com/antinvestor/service-profile/apps/devices/service/models"

	"github.com/pitabwire/frame"
)

func Migrate(ctx context.Context, svc *frame.Service, migrationPath string) error {
	return svc.MigrateDatastore(ctx, migrationPath, &models.Device{}, &models.DeviceLog{})
}
