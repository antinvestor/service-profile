package business_test

import (
	"context"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service/business"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/pitabwire/frame"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type AddressTestSuite struct {
	BaseTestSuite
}

func (ats *AddressTestSuite) SetupSuite() {
	ats.BaseTestSuite.SetupSuite()

}

func TestAddressSuite(t *testing.T) {
	suite.Run(t, new(AddressTestSuite))
}

func (ats *AddressTestSuite) TestNewAddressBusiness() {
	t := ats.T()

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
				ctx:     ats.ctx,
				service: ats.service,
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

func (ats *AddressTestSuite) Test_addressBusiness_CreateAddress() {
	t := ats.T()

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
				service: ats.service,
			},
			args: args{
				ctx:     ats.ctx,
				request: adObj,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aB := business.NewAddressBusiness(ats.ctx, tt.fields.service)
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

func (ats *AddressTestSuite) Test_addressBusiness_GetByProfile() {

	t := ats.T()

	ctx := ats.ctx

	testProfiles, err := ats.createTestProfiles([]string{"testing@ant.com"})
	if err != nil {
		t.Errorf(" CreateProfile failed with %+v", err)
		return
	}

	profile := testProfiles[0]

	addBuss := business.NewAddressBusiness(ctx, ats.service)

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

func (ats *AddressTestSuite) Test_addressBusiness_LinkAddressToProfile() {
	t := ats.T()
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

func (ats *AddressTestSuite) Test_addressBusiness_ToAPI() {
	t := ats.T()
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
