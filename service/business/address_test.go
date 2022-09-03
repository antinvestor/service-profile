package business

import (
	"context"
	"crypto/sha256"
	"fmt"
	profilev1 "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/pitabwire/frame"
	"golang.org/x/crypto/pbkdf2"
	"reflect"
	"testing"
)

func testService(ctx context.Context) *frame.Service {
	dbURL := frame.GetEnv("TEST_DATABASE_URL", "postgres://ant:secret@localhost:5434/service_profile?sslmode=disable")
	mainDB := frame.DatastoreCon(ctx, dbURL, false)

	configProfile := config.Profile{
		QueueVerification:     fmt.Sprintf("mem://%s", "QueueVerificationName"),
		QueueVerificationName: "QueueVerificationName",
	}

	verificationQueuePublisher := frame.RegisterPublisher(configProfile.QueueVerificationName, configProfile.QueueVerification)

	service := frame.NewService("profile tests", mainDB,
		verificationQueuePublisher, frame.Config(&configProfile), frame.NoopDriver())
	_ = service.Run(ctx, "")
	return service
}

func getEncryptionKey() []byte {
	return pbkdf2.Key([]byte("ualgJEcb4GNXLn3jYV9TUGtgYrdTMg"), []byte("VufLmnycUCgz"), 4096, 32, sha256.New)
}

func TestNewAddressBusiness(t *testing.T) {
	ctx := context.Background()
	srv := testService(ctx)
	type args struct {
		ctx     context.Context
		service *frame.Service
	}

	tests := []struct {
		name string
		args args
		want AddressBusiness
	}{
		{
			name: "New Address Business test",
			args: args{
				ctx:     ctx,
				service: srv,
			},
			want: &addressBusiness{
				service:     srv,
				addressRepo: repository.NewAddressRepository(srv),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAddressBusiness(tt.args.ctx, tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAddressBusiness() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addressBusiness_CreateAddress(t *testing.T) {
	ctx := context.Background()
	srv := testService(ctx)
	addRepo := repository.NewAddressRepository(srv)

	adObj := &profilev1.AddressObject{
		Name:    "test address",
		Area:    "Town",
		Country: "KEN",
	}

	type fields struct {
		service     *frame.Service
		addressRepo repository.AddressRepository
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
				service:     srv,
				addressRepo: addRepo,
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
			aB := &addressBusiness{
				service:     tt.fields.service,
				addressRepo: tt.fields.addressRepo,
			}
			got, err := aB.CreateAddress(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAddress() error = %v, wantErr %+v", err, tt.wantErr)
				return
			}
			if got == nil || got.ID == "" || got.Name != adObj.Name || got.Area != adObj.Area || got.Country != "Kenya" {
				t.Errorf("CreateAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addressBusiness_GetByProfile(t *testing.T) {

	ctx := context.Background()
	srv := testService(ctx)
	encryptionKey := getEncryptionKey()

	profBuss := NewProfileBusiness(ctx, srv)

	prof := &profilev1.ProfileCreateRequest{
		Contact: "testing@ant.com",
	}
	profile, err := profBuss.CreateProfile(ctx, encryptionKey, prof)
	if err != nil {
		t.Errorf(" CreateProfile failed with %+v", err)
		return
	}

	addBuss := NewAddressBusiness(ctx, srv)

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

	err = addBuss.LinkAddressToProfile(ctx, profile.GetID(), "Test Link", add)
	if err != nil {
		t.Errorf(" LinkAddressToProfile failed with %+v", err)
		return
	}

	addresses, err := addBuss.GetByProfile(ctx, profile.GetID())
	if err != nil {
		t.Errorf(" GetByProfile failed with %+v", err)
		return
	}

	if len(addresses) == 0 {
		t.Errorf(" GetByProfile failed with %+v", err)
	}
}

func Test_addressBusiness_LinkAddressToProfile(t *testing.T) {
	type fields struct {
		service     *frame.Service
		addressRepo repository.AddressRepository
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
			aB := &addressBusiness{
				service:     tt.fields.service,
				addressRepo: tt.fields.addressRepo,
			}
			if err := aB.LinkAddressToProfile(tt.args.ctx, tt.args.profileID, tt.args.name, tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("LinkAddressToProfile() error = %v, wantErr %+v", err, tt.wantErr)
			}
		})
	}
}

func Test_addressBusiness_ToAPI(t *testing.T) {
	type fields struct {
		service     *frame.Service
		addressRepo repository.AddressRepository
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
			aB := &addressBusiness{
				service:     tt.fields.service,
				addressRepo: tt.fields.addressRepo,
			}
			if got := aB.ToAPI(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}
