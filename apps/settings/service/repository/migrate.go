package repository

import (
	"context"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
)

func Migrate(ctx context.Context, dbManager datastore.Manager, dbPool pool.Pool, migrationPath string) error {
	return dbManager.Migrate(ctx, dbPool, migrationPath,
		models.SettingRef{}, models.SettingVal{}, models.SettingAudit{},
	)
}
