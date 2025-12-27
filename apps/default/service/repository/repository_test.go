package repository_test

import (
	"context"
	"testing"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/apps/default/tests"
)

type RepositoryTestSuite struct {
	tests.ProfileBaseTestSuite
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (rts *RepositoryTestSuite) getRepositories(ctx context.Context, svc *frame.Service) (
	repository.ContactRepository,
	repository.ProfileRepository,
	repository.VerificationRepository,
) {
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)
	workMan := svc.WorkManager()

	contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
	profileRepo := repository.NewProfileRepository(ctx, dbPool, workMan)
	verificationRepo := repository.NewVerificationRepository(ctx, dbPool, workMan)

	return contactRepo, profileRepo, verificationRepo
}

func (rts *RepositoryTestSuite) TestContactRepository_Create() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)
		contactRepo, _, _ := rts.getRepositories(ctx, svc)

		contact := &models.Contact{
			LookUpToken:     []byte("test-lookup-token"),
			EncryptedDetail: []byte("encrypted-detail"),
			EncryptionKeyID: "test-key-id",
			ContactType:     "EMAIL",
		}
		contact.GenID(ctx)

		err := contactRepo.Create(ctx, contact)
		require.NoError(t, err)
		require.NotEmpty(t, contact.GetID())

		// Retrieve the contact
		retrieved, err := contactRepo.GetByID(ctx, contact.GetID())
		require.NoError(t, err)
		require.Equal(t, contact.GetID(), retrieved.GetID())
		require.Equal(t, contact.ContactType, retrieved.ContactType)
	})
}

func (rts *RepositoryTestSuite) TestContactRepository_GetByLookupToken() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)
		contactRepo, _, _ := rts.getRepositories(ctx, svc)

		lookupToken := []byte("unique-lookup-token-" + util.IDString())
		contact := &models.Contact{
			LookUpToken:     lookupToken,
			EncryptedDetail: []byte("encrypted-detail"),
			EncryptionKeyID: "test-key-id",
			ContactType:     "MSISDN",
		}
		contact.GenID(ctx)

		err := contactRepo.Create(ctx, contact)
		require.NoError(t, err)

		// Get by lookup token
		contacts, err := contactRepo.GetByLookupToken(ctx, lookupToken)
		require.NoError(t, err)
		require.Len(t, contacts, 1)
		require.Equal(t, contact.GetID(), contacts[0].GetID())
	})
}

func (rts *RepositoryTestSuite) TestContactRepository_GetByProfileID() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)
		contactRepo, _, _ := rts.getRepositories(ctx, svc)

		profileID := util.IDString()

		// Create multiple contacts for the same profile
		for i := 0; i < 3; i++ {
			contact := &models.Contact{
				LookUpToken:     []byte("lookup-token-" + util.IDString()),
				EncryptedDetail: []byte("encrypted-detail"),
				EncryptionKeyID: "test-key-id",
				ContactType:     "EMAIL",
				ProfileID:       profileID,
			}
			contact.GenID(ctx)
			err := contactRepo.Create(ctx, contact)
			require.NoError(t, err)
		}

		// Get by profile ID
		contacts, err := contactRepo.GetByProfileID(ctx, profileID)
		require.NoError(t, err)
		require.Len(t, contacts, 3)
	})
}

func (rts *RepositoryTestSuite) TestProfileRepository_Create() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)
		_, profileRepo, _ := rts.getRepositories(ctx, svc)

		// Create a profile type first
		profileType := &models.ProfileType{
			UID:         1,
			Name:        "person",
			Description: "A person profile",
		}
		profileType.GenID(ctx)
		dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)
		err := dbPool.DB(ctx, false).Create(profileType).Error
		require.NoError(t, err)

		profile := &models.Profile{
			Properties:    map[string]any{"name": "Test User"},
			ProfileTypeID: profileType.GetID(),
		}
		profile.GenID(ctx)

		err = profileRepo.Create(ctx, profile)
		require.NoError(t, err)
		require.NotEmpty(t, profile.GetID())

		// Retrieve the profile
		retrieved, err := profileRepo.GetByID(ctx, profile.GetID())
		require.NoError(t, err)
		require.Equal(t, profile.GetID(), retrieved.GetID())
	})
}

func (rts *RepositoryTestSuite) TestVerificationRepository_Create() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)
		_, _, verificationRepo := rts.getRepositories(ctx, svc)

		verification := &models.Verification{
			ProfileID: util.IDString(),
			ContactID: util.IDString(),
			Code:      "123456",
		}
		verification.GenID(ctx)

		err := verificationRepo.Create(ctx, verification)
		require.NoError(t, err)
		require.NotEmpty(t, verification.GetID())

		// Retrieve the verification
		retrieved, err := verificationRepo.GetByID(ctx, verification.GetID())
		require.NoError(t, err)
		require.Equal(t, verification.GetID(), retrieved.GetID())
		require.Equal(t, verification.Code, retrieved.Code)
	})
}

func (rts *RepositoryTestSuite) TestVerificationRepository_SaveAttempt() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)
		_, _, verificationRepo := rts.getRepositories(ctx, svc)

		// Create a verification first
		verification := &models.Verification{
			ProfileID: util.IDString(),
			ContactID: util.IDString(),
			Code:      "654321",
		}
		verification.GenID(ctx)
		err := verificationRepo.Create(ctx, verification)
		require.NoError(t, err)

		// Create an attempt
		attempt := &models.VerificationAttempt{
			VerificationID: verification.GetID(),
			Data:           "test-attempt-data",
			State:          "pending",
			DeviceID:       "device-123",
			IPAddress:      "192.168.1.1",
		}
		attempt.GenID(ctx)

		err = verificationRepo.SaveAttempt(ctx, attempt)
		require.NoError(t, err)

		// Get attempts
		attempts, err := verificationRepo.GetAttempts(ctx, verification.GetID())
		require.NoError(t, err)
		require.Len(t, attempts, 1)
		require.Equal(t, attempt.Data, attempts[0].Data)
	})
}

func (rts *RepositoryTestSuite) TestMigrate() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)

		// Migrate should not fail on an already migrated database
		err := repository.Migrate(ctx, svc.DatastoreManager(), "../../migrations")
		require.NoError(t, err)
	})
}

func (rts *RepositoryTestSuite) TestContactRepository_DelinkFromProfile() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)
		contactRepo, _, _ := rts.getRepositories(ctx, svc)

		profileID := util.IDString()
		contact := &models.Contact{
			LookUpToken:     []byte("delink-lookup-token-" + util.IDString()),
			EncryptedDetail: []byte("encrypted-detail"),
			EncryptionKeyID: "test-key-id",
			ContactType:     "EMAIL",
			ProfileID:       profileID,
		}
		contact.GenID(ctx)

		err := contactRepo.Create(ctx, contact)
		require.NoError(t, err)

		// Delink from profile
		delinked, err := contactRepo.DelinkFromProfile(ctx, contact.GetID(), profileID)
		require.NoError(t, err)
		require.Empty(t, delinked.ProfileID)
	})
}

func (rts *RepositoryTestSuite) TestContactRepository_DelinkFromProfile_WrongProfile() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)
		contactRepo, _, _ := rts.getRepositories(ctx, svc)

		profileID := util.IDString()
		contact := &models.Contact{
			LookUpToken:     []byte("delink-wrong-token-" + util.IDString()),
			EncryptedDetail: []byte("encrypted-detail"),
			EncryptionKeyID: "test-key-id",
			ContactType:     "EMAIL",
			ProfileID:       profileID,
		}
		contact.GenID(ctx)

		err := contactRepo.Create(ctx, contact)
		require.NoError(t, err)

		// Delink with wrong profile - should fail
		_, err = contactRepo.DelinkFromProfile(ctx, contact.GetID(), "wrong-profile")
		require.Error(t, err)
	})
}
