package business_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service/business"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/pitabwire/frame"
	"golang.org/x/crypto/pbkdf2"
	"reflect"
	"testing"
)

func getTestService() (context.Context, *frame.Service) {
	dbURL := frame.GetEnv("TEST_DATABASE_URL",
		"postgres://ant:secret@localhost:5434/service_profile?sslmode=disable")
	mainDB := frame.DatastoreCon(dbURL, false)

	configProfile := config.ProfileConfig{
		QueueVerification:     fmt.Sprintf("mem://%s", "QueueVerificationName"),
		QueueVerificationName: "QueueVerificationName",
	}

	verificationQueuePublisher := frame.RegisterPublisher(
		configProfile.QueueVerificationName, configProfile.QueueVerification)

	ctx, service := frame.NewService("profile tests", mainDB,
		verificationQueuePublisher, frame.Config(&configProfile), frame.NoopDriver())
	_ = service.Run(ctx, "")
	return ctx, service
}

func getEncryptionKey() []byte {
	return pbkdf2.Key([]byte("ualgJEcb4GNXLn3jYV9TUGtgYrdTMg"), []byte("VufLmnycUCgz"), 4096, 32, sha256.New)
}

func TestNewAddressBusiness(t *testing.T) {
	ctx, srv := getTestService()
	type args struct {
		ctx     context.Context
		service *frame.Service
	}

	tests := []struct {
		name string
		args args
		want business.AddressBusiness
	}{
		{
			name: "New Address Business test",
			args: args{
				ctx:     ctx,
				service: srv,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := business.NewAddressBusiness(tt.args.ctx, tt.args.service); got == nil {
				t.Errorf("NewAddressBusiness() = %v, want non nil address business", got)
			}
		})
	}
}

func Test_addressBusiness_CreateAddress(t *testing.T) {
	ctx, srv := getTestService()

	adObj := &profilev1.AddressObject{
		Name:    "test address",
		Area:    "Town",
		Country: "KEN",
	}

	type fields struct {
		service *frame.Service
	}
	type args struct {
		ctx     context.Context
		request *profilev1.AddressObject
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *profilev1.AddressObject
		wantErr bool
	}{
		{
			name: "Create Address test",
			fields: fields{
				service: srv,
			},
			args: args{
				ctx:     ctx,
				request: adObj,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aB := business.NewAddressBusiness(ctx, tt.fields.service)
			got, err := aB.CreateAddress(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAddress() error = %v, wantErr %+v", err, tt.wantErr)
				return
			}
			if got == nil || got.GetId() == "" || got.Name != adObj.Name || got.Area != adObj.Area || got.Country != "Kenya" {
				t.Errorf("CreateAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addressBusiness_GetByProfile(t *testing.T) {

	ctx, srv := getTestService()
	encryptionKey := getEncryptionKey()

	testProfiles, err := createTestProfiles(ctx, srv, encryptionKey, []string{"testing@ant.com"})
	if err != nil {
		t.Errorf(" CreateProfile failed with %+v", err)
		return
	}

	profile := testProfiles[0]

	addBuss := business.NewAddressBusiness(ctx, srv)

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
}

func Test_addressBusiness_LinkAddressToProfile(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		service *frame.Service
	}
	type args struct {
		ctx       context.Context
		profileID string
		name      string
		address   *profilev1.AddressObject
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aB := business.NewAddressBusiness(ctx, tt.fields.service)
			if err := aB.LinkAddressToProfile(tt.args.ctx, tt.args.profileID, tt.args.name, tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("LinkAddressToProfile() error = %v, wantErr %+v", err, tt.wantErr)
			}
		})
	}
}

func Test_addressBusiness_ToAPI(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		service *frame.Service
	}
	type args struct {
		address *models.Address
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *profilev1.AddressObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aB := business.NewAddressBusiness(ctx, tt.fields.service)
			if got := aB.ToAPI(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}
