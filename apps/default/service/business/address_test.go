package business_test

import (
	"context"
	"testing"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests/definition"
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

	contactBusiness := business.NewContactBusiness(ctx, cfg, evtsMan, contactRepo, verificationRepo)

	addressRepo := repository.NewAddressRepository(ctx, dbPool, workMan)
	addressBusiness := business.NewAddressBusiness(ctx, addressRepo)

	profileRepo := repository.NewProfileRepository(ctx, dbPool, workMan)
	return business.NewProfileBusiness(ctx, evtsMan, contactBusiness, addressBusiness, profileRepo), addressRepo
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
				svc, ctx := ats.CreateService(t, dep)

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
				svc, ctx := ats.CreateService(t, dep)

				dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)
				addressRepo := repository.NewAddressRepository(ctx, dbPool, svc.WorkManager())

				aB := business.NewAddressBusiness(ctx, addressRepo)
				got, err := aB.CreateAddress(ctx, tt.request)
				tt.wantErr(t, err)

				if got == nil || got.GetId() == "" || got.GetName() != adObj.GetName() ||
					got.GetArea() != adObj.GetArea() ||
					got.GetCountry() != "Kenya" {
					t.Errorf("CreateAddress() got = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func (ats *AddressTestSuite) Test_addressBusiness_GetByProfile() {
	t := ats.T()

	ats.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := ats.CreateService(t, dep)

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
