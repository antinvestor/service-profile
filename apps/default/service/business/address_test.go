package business_test

import (
	"context"
	"encoding/base64"
	"testing"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/apps/default/tests"
)

type AddressTestSuite struct {
	tests.ProfileBaseTestSuite
}

func TestAddressSuite(t *testing.T) {
	suite.Run(t, new(AddressTestSuite))
}

// Helper function to create consistent test DEK.
func createAddressTestDEK(cfg *config.ProfileConfig) *config.DEK {
	// Decode base64 keys
	key, _ := base64.StdEncoding.DecodeString(cfg.DEKActiveAES256GCMKey)
	lookupKey, _ := base64.StdEncoding.DecodeString(cfg.DEKLookupTokenHMACSHA256Key)

	return &config.DEK{
		KeyID:     cfg.DEKActiveKeyID,
		Key:       key,
		OldKeyID:  "old-key-id",
		OldKey:    []byte("1234567890123456"), // 16 bytes for old key
		LookUpKey: lookupKey,
	}
}

func (ats *AddressTestSuite) getProfileBusiness(
	ctx context.Context,
	svc *frame.Service,
) (business.ProfileBusiness, repository.AddressRepository) {
	evtsMan := svc.EventsManager()
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	cfg := svc.Config().(*config.ProfileConfig)

	contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
	verificationRepo := repository.NewVerificationRepository(ctx, dbPool, workMan)

	contactBusiness := business.NewContactBusiness(
		ctx,
		cfg,
		createAddressTestDEK(cfg),
		evtsMan,
		contactRepo,
		verificationRepo,
	)

	addressRepo := repository.NewAddressRepository(ctx, dbPool, workMan)
	addressBusiness := business.NewAddressBusiness(ctx, addressRepo)

	profileRepo := repository.NewProfileRepository(ctx, dbPool, workMan)
	propertyEntryRepo := repository.NewPropertyEntryRepository(ctx, dbPool, workMan)
	return business.NewProfileBusiness(
		ctx,
		cfg,
		createAddressTestDEK(cfg),
		evtsMan,
		contactBusiness,
		addressBusiness,
		profileRepo,
		propertyEntryRepo,
	), addressRepo
}

func (ats *AddressTestSuite) TestNewAddressBusiness() {
	testcases := []struct {
		name string
		want business.AddressBusiness
	}{
		{
			name: "New Address Business test",
		},
	}

	ats.WithTestDependancies(ats.T(), func(t *testing.T, dep *definition.DependencyOption) {
		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				ctx, svc := ats.CreateService(t, dep)

				dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)
				addressRepo := repository.NewAddressRepository(ctx, dbPool, svc.WorkManager())

				if got := business.NewAddressBusiness(ctx, addressRepo); got == nil {
					t.Errorf("NewAddressBusiness() = %v, want non nil address business", got)
				}
			})
		}
	})
}

func (ats *AddressTestSuite) Test_addressBusiness_CreateAddress() {
	t := ats.T()

	adObj := &profilev1.AddressObject{
		Name:    "test address",
		Area:    "Town",
		Country: "KEN",
	}

	testCases := []struct {
		name    string
		request *profilev1.AddressObject
		want    *profilev1.AddressObject
		wantErr require.ErrorAssertionFunc
	}{
		{
			name:    "Create Address test",
			request: adObj,
			want:    nil,
			wantErr: require.NoError,
		},
	}

	ats.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				ctx, svc := ats.CreateService(t, dep)

				dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)
				addressRepo := repository.NewAddressRepository(ctx, dbPool, svc.WorkManager())

				aB := business.NewAddressBusiness(ctx, addressRepo)
				got, err := aB.CreateAddress(ctx, tt.request)
				tt.wantErr(t, err)

				_ = got // Use got to avoid unused variable warning
			})
		}
	})
}

func (ats *AddressTestSuite) Test_addressBusiness_GetByProfile() {
	t := ats.T()

	ats.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := ats.CreateService(t, dep)

		profileBusiness, addressRepo := ats.getProfileBusiness(ctx, svc)

		testProfiles, err := ats.CreateTestProfiles(ctx, profileBusiness, []string{"testing@ant.com"})
		if err != nil {
			t.Errorf(" CreateProfile failed with %+v", err)
			return
		}

		profile := testProfiles[0]

		addBuss := business.NewAddressBusiness(ctx, addressRepo)

		adObj := &profilev1.AddressObject{
			Name:    "Linked address",
			Area:    "Town",
			Country: "KEN",
		}

		add, err := addBuss.CreateAddress(ctx, adObj)
		if err != nil {
			t.Errorf(" CreateAddress failed with %+v", err)
			return
		}

		err = addBuss.LinkAddressToProfile(ctx, profile.GetId(), "Test Link", add)
		if err != nil {
			t.Errorf(" LinkAddressToProfile failed with %+v", err)
			return
		}

		addresses, err := addBuss.GetByProfile(ctx, profile.GetId())
		if err != nil {
			t.Errorf(" GetByProfile failed with %+v", err)
			return
		}

		if len(addresses) == 0 {
			t.Errorf(" GetByProfile failed with %+v", err)
		}
	})
}

// Test_addressBusiness_GetByProfile_WithTenantContext verifies that addresses
// linked to a profile under one tenant are visible from a different tenant context.
// Addresses are cross-tenant identity data.
func (ats *AddressTestSuite) Test_addressBusiness_GetByProfile_WithTenantContext() {
	t := ats.T()

	ats.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := ats.CreateService(t, dep)

		// Create profile under tenant A
		tenantA := util.IDString()
		partitionA := util.IDString()
		ctxA := ats.WithAuthClaims(ctx, tenantA, partitionA, util.IDString())

		profileBusiness, addressRepo := ats.getProfileBusiness(ctxA, svc)

		testProfiles, err := ats.CreateTestProfiles(ctxA, profileBusiness, []string{"addr.tenant@testing.com"})
		require.NoError(t, err, "CreateTestProfiles should succeed under tenant A")

		profile := testProfiles[0]

		// Add an address under tenant A context
		addBuss := business.NewAddressBusiness(ctxA, addressRepo)

		adObj := &profilev1.AddressObject{
			Name:    "Cross-tenant address",
			Area:    "Town",
			Country: "KEN",
		}

		add, err := addBuss.CreateAddress(ctxA, adObj)
		require.NoError(t, err, "CreateAddress should succeed under tenant A")

		err = addBuss.LinkAddressToProfile(ctxA, profile.GetId(), "Test Link", add)
		require.NoError(t, err, "LinkAddressToProfile should succeed under tenant A")

		// Switch to tenant B context and read addresses
		tenantB := util.IDString()
		partitionB := util.IDString()
		ctxB := ats.WithAuthClaims(ctx, tenantB, partitionB, util.IDString())

		// Re-create business with tenant B context repos
		_, addressRepoB := ats.getProfileBusiness(ctxB, svc)
		addBussB := business.NewAddressBusiness(ctxB, addressRepoB)

		addresses, err := addBussB.GetByProfile(ctxB, profile.GetId())
		require.NoError(t, err, "GetByProfile should work cross-tenant")
		require.NotEmpty(t, addresses, "addresses should be visible from tenant B")
		require.Equal(t, add.GetId(), addresses[0].AddressID, "address ID should match")
	})
}
