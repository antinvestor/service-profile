package business_test

import (
	"context"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service/business"
	"reflect"
	"testing"
)

func Test_profileBusiness_CreateProfile(t *testing.T) {

	ctx, srv := getTestService()
	encryptionKey := getEncryptionKey()

	tests := []struct {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := business.NewProfileBusiness(ctx, srv, func() []byte { return encryptionKey })
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
}

func Test_profileBusiness_GetByID(t *testing.T) {

	ctx, srv := getTestService()
	encryptionKey := getEncryptionKey()

	var profileAvailable []string
	pbc := business.NewProfileBusiness(ctx, srv, func() []byte { return encryptionKey })

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

	type args struct {
		ctx                  context.Context
		createProfileRequest *profilev1.CreateRequest
		profileID            string
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

			pb := business.NewProfileBusiness(ctx, srv, func() []byte { return encryptionKey })

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
}
