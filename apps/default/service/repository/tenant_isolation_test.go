package repository_test

import (
	"testing"

	"github.com/pitabwire/frame/v2/datastore"
	"github.com/pitabwire/frame/v2/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

// TestRosterTenantIsolation verifies Postgres RLS actually hides
// tenant-scoped rows across tenants. Rosters are tenant-private (the
// repository does NOT use security.SkipTenancyChecksOnClaims), unlike
// profiles/contacts which are deliberately globally readable.
func (rts *RepositoryTestSuite) TestRosterTenantIsolation() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)

		dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)
		workMan := svc.WorkManager()
		contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
		rosterRepo := repository.NewRosterRepository(ctx, dbPool, workMan)

		ctxA := rts.WithAuthClaims(ctx, "tenant-a", "partition-a", util.IDString())
		ctxB := rts.WithAuthClaims(ctx, "tenant-b", "partition-b", util.IDString())

		profileID := util.IDString()

		contact := &models.Contact{
			LookUpToken:     []byte("isolation-token-" + util.IDString()),
			EncryptedDetail: []byte("encrypted-detail"),
			EncryptionKeyID: "test-key-id",
			ContactType:     "EMAIL",
			ProfileID:       profileID,
		}
		contact.GenID(ctxA)
		require.NoError(t, contactRepo.Create(ctxA, contact))

		roster := &models.Roster{
			ProfileID: profileID,
			ContactID: contact.GetID(),
			Name:      "isolation-roster",
		}
		roster.GenID(ctxA)
		require.NoError(t, rosterRepo.Create(ctxA, roster))

		// Owning tenant can read its roster entry back.
		got, err := rosterRepo.GetByID(ctxA, roster.GetID())
		require.NoError(t, err)
		require.Equal(t, roster.GetID(), got.GetID())

		// Another tenant must not be able to read it.
		_, err = rosterRepo.GetByID(ctxB, roster.GetID())
		require.Error(t, err, "tenant B must not read tenant A's roster entry")

		entries, err := rosterRepo.GetByContactIDsAndProfileID(ctxB, []string{contact.GetID()}, profileID)
		require.NoError(t, err)
		require.Empty(t, entries, "tenant B must not list tenant A's roster entries")
	})
}

// TestVerificationTenantIsolation verifies contact verifications are
// tenant-private even though the contacts they belong to are globally
// readable by design.
func (rts *RepositoryTestSuite) TestVerificationTenantIsolation() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)
		_, _, verificationRepo := rts.getRepositories(ctx, svc)

		ctxA := rts.WithAuthClaims(ctx, "tenant-a", "partition-a", util.IDString())
		ctxB := rts.WithAuthClaims(ctx, "tenant-b", "partition-b", util.IDString())

		verification := &models.Verification{
			ProfileID: util.IDString(),
			ContactID: util.IDString(),
			Code:      "123456",
		}
		verification.GenID(ctxA)
		require.NoError(t, verificationRepo.Create(ctxA, verification))

		// Owning tenant can read the verification back.
		got, err := verificationRepo.GetByID(ctxA, verification.GetID())
		require.NoError(t, err)
		require.Equal(t, verification.GetID(), got.GetID())

		// Another tenant must not see it.
		_, err = verificationRepo.GetByID(ctxB, verification.GetID())
		require.Error(t, err, "tenant B must not read tenant A's verification")
	})
}
