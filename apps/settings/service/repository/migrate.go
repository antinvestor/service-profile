package repository

import (
	"context"

	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
)

func Migrate(ctx context.Context, svc *frame.Service, migrationPath string) error {
	return svc.MigrateDatastore(ctx, migrationPath,
		models.SettingRef{}, models.SettingVal{}, models.SettingAudit{})
}
