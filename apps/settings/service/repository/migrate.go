package repository

import (
	"context"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
	"github.com/pitabwire/frame"
)

func Migrate(ctx context.Context, svc *frame.Service, migrationPath string) error {
	return svc.MigrateDatastore(ctx, migrationPath,
		models.SettingRef{}, models.SettingVal{}, models.SettingAudit{})
}
