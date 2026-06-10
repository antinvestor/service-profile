package business_test

import (
	"testing"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
	"github.com/antinvestor/service-profile/apps/settings/service/repository"
)

// TestSettingValTenantIsolation verifies Postgres RLS hides one
// tenant's setting values from another tenant. The suite drops
// application queries to an unprivileged role, so the tenancy
// policies installed during migration are actually enforced.
func (ts *SettingsTestSuite) TestSettingValTenantIsolation() {
	ts.WithTestDependancies(ts.T(), func(t *testing.T, depOpt *definition.DependencyOption) {
		ctx, svc := ts.CreateService(t, depOpt)

		workMan := svc.WorkManager()
		dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)
		refRepo := repository.NewReferenceRepository(ctx, dbPool, workMan)
		valRepo := repository.NewSettingValRepository(ctx, dbPool, workMan)

		ctxA := ts.WithAuthClaims(ctx, "tenant-a", "partition-a", util.IDString())
		ctxB := ts.WithAuthClaims(ctx, "tenant-b", "partition-b", util.IDString())

		ref := &models.SettingRef{
			Name:     "isolation-setting",
			Object:   "profile",
			ObjectID: util.IDString(),
			Language: "en",
			Module:   "test",
		}
		ref.GenID(ctxA)
		require.NoError(t, refRepo.Create(ctxA, ref))

		val := &models.SettingVal{
			Ref:    ref.GetID(),
			Detail: "tenant-a-secret-value",
		}
		val.GenID(ctxA)
		require.NoError(t, valRepo.Create(ctxA, val))

		// Owning tenant reads its value back.
		got, err := valRepo.GetByID(ctxA, val.GetID())
		require.NoError(t, err)
		require.Equal(t, "tenant-a-secret-value", got.Detail)

		// Another tenant must not be able to read it.
		_, err = valRepo.GetByID(ctxB, val.GetID())
		require.Error(t, err, "tenant B must not read tenant A's setting value")

		vals, err := valRepo.GetByRef(ctxB, ref.GetID())
		require.NoError(t, err)
		require.Empty(t, vals, "tenant B must not list tenant A's setting values")

		// The reference itself is tenant-scoped too.
		_, err = refRepo.GetByID(ctxB, ref.GetID())
		require.Error(t, err, "tenant B must not read tenant A's setting reference")
	})
}
