package repository

import (
	"context"

	"github.com/pitabwire/frame/datastore"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
)

func Migrate(ctx context.Context, dbManager datastore.Manager, migrationPath string) error {
	dbPool := dbManager.GetPool(ctx, datastore.DefaultMigrationPoolName)

	// Models MUST be registered as pointers: tenancy.Tenanted is only
	// satisfied by *data.BaseModel receivers, so value-registered models
	// are silently skipped during RLS policy installation.
	return dbManager.Migrate(ctx, dbPool, migrationPath,
		&models.SettingRef{}, &models.SettingVal{}, &models.SettingAudit{},
	)
}
