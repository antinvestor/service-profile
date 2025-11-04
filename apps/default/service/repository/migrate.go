package repository

import (
	"context"

	"github.com/pitabwire/frame/datastore"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)

func Migrate(ctx context.Context, dbManager datastore.Manager, migrationPath string) error {
	dbPool := dbManager.GetPool(ctx, datastore.DefaultMigrationPoolName)

	return dbManager.Migrate(ctx, dbPool, migrationPath,
		&models.ProfileType{}, &models.Profile{}, &models.Contact{}, &models.Country{},
		&models.Address{}, &models.ProfileAddress{}, &models.Verification{}, &models.VerificationAttempt{},
		&models.RelationshipType{}, &models.Relationship{}, &models.Roster{},
	)
}
