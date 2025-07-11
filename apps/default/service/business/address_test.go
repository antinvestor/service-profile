package business_test

import (
	"reflect"
	"testing"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame/tests/testdef"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/internal/tests"
)

type AddressTestSuite struct {
	tests.BaseTestSuite
}

func TestAddressSuite(t *testing.T) {
	suite.Run(t, new(AddressTestSuite))
}

func (ats *AddressTestSuite) TestNewAddressBusiness() {
	t := ats.T()

	tests := []struct {
		name string
		want business.AddressBusiness
	}{
		{
			name: "New Address Business test",
		},
	}

	ats.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				svc, ctx := ats.CreateService(t, dep)

				if got := business.NewAddressBusiness(ctx, svc); got == nil {
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

	tests := []struct {
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

	ats.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				svc, ctx := ats.CreateService(t, dep)

				aB := business.NewAddressBusiness(ctx, svc)
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

	ats.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := ats.CreateService(t, dep)

		testProfiles, err := ats.CreateTestProfiles(ctx, svc, []string{"testing@ant.com"})
		if err != nil {
			t.Errorf(" CreateProfile failed with %+v", err)
			return
		}

		profile := testProfiles[0]

		addBuss := business.NewAddressBusiness(ctx, svc)

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

func (ats *AddressTestSuite) Test_addressBusiness_LinkAddressToProfile() {
	t := ats.T()

	type args struct {
		profileID string
		name      string
		address   *profilev1.AddressObject
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	ats.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				svc, ctx := ats.CreateService(t, dep)

				aB := business.NewAddressBusiness(ctx, svc)
				if err := aB.LinkAddressToProfile(ctx, tt.args.profileID, tt.args.name, tt.args.address); (err != nil) != tt.wantErr {
					t.Errorf("LinkAddressToProfile() error = %v, wantErr %+v", err, tt.wantErr)
				}
			})
		}
	})
}

func (ats *AddressTestSuite) Test_addressBusiness_ToAPI() {
	t := ats.T()

	type args struct {
		address *models.Address
	}
	tests := []struct {
		name string
		args args
		want *profilev1.AddressObject
	}{
		// TODO: Add test cases.
	}

	ats.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				svc, ctx := ats.CreateService(t, dep)

				aB := business.NewAddressBusiness(ctx, svc)
				if got := aB.ToAPI(tt.args.address); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ToAPI() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
