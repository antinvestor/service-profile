package business_test

import (
	profilev1 "github.com/antinvestor/apis/profile/v1"
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
