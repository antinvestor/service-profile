package repository

import (
	"context"

	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/pitabwire/frame"
)

func Migrate(ctx context.Context, svc *frame.Service, migrationPath string) error {
	return svc.MigrateDatastore(ctx, migrationPath,
		&models.ProfileType{}, &models.Profile{}, &models.Contact{}, &models.Country{},
		&models.Address{}, &models.ProfileAddress{}, &models.Verification{}, &models.VerificationAttempt{},
		&models.RelationshipType{}, &models.Relationship{}, &models.Roster{})
}
