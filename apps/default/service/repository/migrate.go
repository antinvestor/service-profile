package repository

import (
	"context"

	"github.com/pitabwire/frame/v2/datastore"

	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/internal/outbox"
)

func Migrate(ctx context.Context, dbManager datastore.Manager, migrationPath string) error {
	dbPool := dbManager.GetPool(ctx, datastore.DefaultMigrationPoolName)

	return dbManager.Migrate(ctx, dbPool, migrationPath,
		&models.ProfileType{}, &models.Profile{}, &models.PropertyEntry{}, &models.Contact{}, &models.Country{},
		&models.Address{}, &models.ProfileAddress{}, &models.Verification{}, &models.VerificationAttempt{},
		&models.RelationshipType{}, &models.Relationship{}, &models.Roster{},
		&models.ProfileAccount{},
		&outbox.Event{},
	)
}
