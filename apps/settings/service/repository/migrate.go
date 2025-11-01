package repository

import (
	"context"
	"errors"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"
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
		models.SettingRef{}, models.SettingVal{}, models.SettingAudit{},
	)
}
