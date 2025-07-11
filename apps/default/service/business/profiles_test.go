package business_test

import (
	"reflect"
	"testing"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame/tests/testdef"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/internal/tests"
)

type ProfileTestSuite struct {
	tests.BaseTestSuite
}

func TestProfileSuite(t *testing.T) {
	suite.Run(t, new(ProfileTestSuite))
}

func (pts *ProfileTestSuite) Test_profileBusiness_CreateProfile() {
	t := pts.T()

	testcases := []struct {
		name    string
		request *profilev1.CreateRequest
		wantErr bool
	}{
		{
			name: "Happy path create a profile",
			request: &profilev1.CreateRequest{
				Type:    profilev1.ProfileType_PERSON,
				Contact: "profile.create@testing.com",
				Properties: map[string]string{
					"name": "Profile Tester",
				},
			},
			wantErr: false,
		},
	}

	pts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				svc, ctx := pts.CreateService(t, dep)
				pb := business.NewProfileBusiness(ctx, svc)
				got, err := pb.CreateProfile(ctx, tt.request)
				if (err != nil) != tt.wantErr {
					t.Errorf("CreateProfile() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if len(got.GetContacts()) != 1 {
					t.Errorf("CreateProfile() does not have a contact ")
				}

				if !reflect.DeepEqual(got.GetProperties(), tt.request.GetProperties()) {
					t.Errorf("CreateProfile() got = %v, want %v", got.GetProperties(), tt.request.GetProperties())
				}
			})
		}
	})
}

func (pts *ProfileTestSuite) Test_profileBusiness_GetByID() {
	t := pts.T()

	pts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := pts.CreateService(t, dep)

		var profileAvailable []string
		pbc := business.NewProfileBusiness(ctx, svc)

		for _, val := range []*profilev1.CreateRequest{
			{
				Type:    profilev1.ProfileType_PERSON,
				Contact: "profile.get.one@testing.com",
				Properties: map[string]string{
					"name": "Profile Tester Get",
				},
			},
			{
				Type:    profilev1.ProfileType_PERSON,
				Contact: "profile.create@testing.com",
				Properties: map[string]string{
					"name": "Profile Tester Get 2",
				},
			},
		} {
			got, err := pbc.CreateProfile(ctx, val)
			if err != nil {
				t.Errorf("CreateProfile() error = %v", err)
				return
			}

			profileAvailable = append(profileAvailable, got.GetId())
		}

		tests := []struct {
			name      string
			profileID string
			wantErr   bool
		}{
			{
				name:      "Happy case 1",
				profileID: profileAvailable[0],
				wantErr:   false,
			},
			{
				name:      "Happy case 2",
				profileID: profileAvailable[1],
				wantErr:   false,
			},
			{
				name:      "Not existing case",
				profileID: "clt0p9viopfc73fdoagg",
				wantErr:   true,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				pb := business.NewProfileBusiness(ctx, svc)

				p, err := pb.GetByID(ctx, tt.profileID)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if err == nil && p == nil {
					t.Error("GetByID() a nil profile should not exist")
				}
			})
		}
	})
}
