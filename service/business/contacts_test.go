package business_test

import (
	"context"
	"fmt"
	"github.com/antinvestor/service-profile/service/business"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/pitabwire/frame"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ContactTestSuite struct {
	BaseTestSuite
}

func TestContactSuite(t *testing.T) {
	suite.Run(t, new(ContactTestSuite))
}

func (cts *ContactTestSuite) TestGeneratePin() {

	t := cts.T()

	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Generate 4-digit PIN",
			args: args{n: 4},
			want: 4, // Placeholder PIN for example purposes, replace with actual expected result based on logic
		},
		{
			name: "Generate 6-digit PIN",
			args: args{n: 6},
			want: 6, // Placeholder PIN for example purposes, replace with actual expected result based on logic
		},
		{
			name: "Generate PIN with zero length",
			args: args{n: 0},
			want: 0, // Expecting no output for zero-length PIN
		},
		{
			name: "Generate PIN with negative length",
			args: args{n: -1},
			want: 0, // Depending on function, it may return empty or handle error
		},
		{
			name: "Generate large PIN",
			args: args{n: 12},
			want: 12, // Placeholder PIN for example purposes
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newPin := business.GeneratePin(tt.args.n)
			assert.Equalf(t, tt.want, len(newPin), "GeneratePin(%v)", tt.args.n)
		})
	}
}

func (cts *ContactTestSuite) Test_contactBusiness_CreateContact() {

	t := cts.T()
	ctx := cts.ctx

	srv := cts.service

	type args struct {
		detail string
		extra  map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Contact
		wantErr assert.ErrorAssertionFunc
	}{

		{
			name: "Create contact with valid MSISDN",
			args: args{
				detail: "+256757546244", // Valid MSISDN
				extra:  map[string]string{"type": "msisdn"},
			},
			want: &models.Contact{ // Expected result
				Detail:     "+256757546244",
				Properties: frame.DBPropertiesFromMap(map[string]string{"type": "msisdn"}),
			},
			wantErr: assert.NoError, // No error expected
		},
		{
			name: "Create contact with valid Email",
			args: args{
				detail: "test@example.com", // Valid Email
				extra:  map[string]string{"type": "email"},
			},
			want: &models.Contact{ // Expected result
				Detail:     "test@example.com",
				Properties: frame.DBPropertiesFromMap(map[string]string{"type": "email"}),
			},
			wantErr: assert.NoError, // No error expected
		},
		{
			name: "Create contact with invalid detail",
			args: args{
				detail: "invalid-detail", // Invalid data, e.g., malformed MSISDN or email
				extra:  map[string]string{"type": "unknown"},
			},
			want:    nil,          // Expect no valid contact to be created
			wantErr: assert.Error, // Error is expected
		},
		{
			name: "Create contact with empty detail",
			args: args{
				detail: "", // Empty detail
				extra:  map[string]string{"type": "email"},
			},
			want:    nil,          // Contact should not be created with empty details
			wantErr: assert.Error, // Error is expected
		},
		{
			name: "Create contact with missing extra data",
			args: args{
				detail: "test2@example.com", // Valid Email but missing extra information
				extra:  nil,                 // Properties data missing
			},
			want: &models.Contact{
				Detail: "test2@example.com",
			},
			wantErr: assert.NoError, // No error expected if type can be inferred or defaults
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cb := business.NewContactBusiness(ctx, srv)
			got, err := cb.CreateContact(ctx, tt.args.detail, tt.args.extra)
			if tt.wantErr(t, err, fmt.Sprintf("CreateContact(ctx, %v, %v)", tt.args.detail, tt.args.extra)) {
				return
			}
			assert.Equalf(t, tt.want.Detail, got.Detail, "CreateContact(ctx, %v, %v)", tt.args.detail, tt.args.extra)
			assert.Equalf(t, tt.want.Properties, got.Properties, "CreateContact(ctx, %v, %v)", tt.args.detail, tt.args.extra)
		})
	}
}

func (cts *ContactTestSuite) createContacts(ctx context.Context, cb business.ContactBusiness, contacDetails ...string) (map[string]*models.Contact, error) {
	result := map[string]*models.Contact{}
	for _, detail := range contacDetails {

		contact, err := cb.CreateContact(ctx, detail, map[string]string{})
		if err != nil {
			return nil, err
		}
		result[contact.Detail] = contact
	}
	return result, nil
}

func (cts *ContactTestSuite) Test_contactBusiness_GetByDetail() {

	t := cts.T()

	cb := business.NewContactBusiness(cts.ctx, cts.service)
	existingContacts, err := cts.createContacts(cts.ctx, cb, "+256757546215", "+256757532244", "bwire@gmail.com")
	assert.NoError(t, err)

	type args struct {
		detail string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Contact
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Get contact by valid MSISDN",
			args: args{
				detail: "+256757546215",
			},
			want:    existingContacts["+256757546215"],
			wantErr: assert.NoError,
		},
		{
			name: "Get contact by valid Email",
			args: args{
				detail: "bwire@gmail.com",
			},
			want:    existingContacts["bwire@gmail.com"],
			wantErr: assert.NoError,
		},
		{
			name: "Get contact by invalid detail",
			args: args{
				detail: "invalid-detail",
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "Get contact by empty detail",
			args: args{
				detail: "",
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var got *models.Contact
			got, err = cb.GetByDetail(cts.ctx, tt.args.detail)
			if tt.wantErr(t, err, fmt.Sprintf("GetByDetail(%v, %v)", cts.ctx, tt.args.detail)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetByDetail(%v, %v)", cts.ctx, tt.args.detail)
		})
	}
}

func (cts *ContactTestSuite) Test_contactBusiness_GetByID() {

	t := cts.T()

	cb := business.NewContactBusiness(cts.ctx, cts.service)
	existingContacts, err := cts.createContacts(cts.ctx, cb, "+256757592215", "+254757532244", "bwireid@gmail.com")
	assert.NoError(t, err)

	type args struct {
		contactID string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Contact
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Get contact with valid ID",
			args: args{
				contactID: existingContacts["+256757592215"].GetID(),
			},
			want:    existingContacts["+256757592215"],
			wantErr: assert.NoError,
		},
		{
			name: "Get contact with invalid ID",
			args: args{
				contactID: "invalid-contact-id",
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "Get contact with empty ID",
			args: args{
				contactID: "",
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := cb.GetByID(cts.ctx, tt.args.contactID)
			if tt.wantErr(t, err, fmt.Sprintf("GetByID(ctx, %v)", tt.args.contactID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetByID(ctx, %v)", tt.args.contactID)
		})
	}
}

func (cts *ContactTestSuite) Test_contactBusiness_GetByProfile() {

	t := cts.T()

	type args struct {
		profileID      string
		contactDetails []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*models.Contact
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Get contacts by valid profile ID",
			args: args{
				profileID:      "valid-profile-id",
				contactDetails: []string{"+256757546200", "testGet@example.com"},
			},
			want: []*models.Contact{
				{Detail: "+256757546200"},
				{Detail: "testGet@example.com"},
			},
		},
		{
			name: "Get contacts by invalid profile ID",
			args: args{
				profileID: "invalid-profile-id",
			},
			want: []*models.Contact{},
		},
		{
			name: "Get contacts by empty profile ID",
			args: args{
				profileID: "",
			},
			want: []*models.Contact{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cb := business.NewContactBusiness(cts.ctx, cts.service)
			existingContacts, err := cts.createContacts(cts.ctx, cb, tt.args.contactDetails...)
			assert.NoError(t, err)

			for _, contact := range existingContacts {
				_, err = cb.UpdateContact(cts.ctx, contact.GetID(), tt.args.profileID, map[string]string{})
				assert.NoError(t, err)
			}

			var got []*models.Contact
			got, err = cb.GetByProfile(cts.ctx, tt.args.profileID)

			assert.Equalf(t, len(tt.want), len(got), "GetByProfile(%v, %v)", cts.ctx, tt.args.profileID)
		})
	}
}

//func (cts *ContactTestSuite) Test_contactBusiness_RemoveContact() {
//
//	t := cts.T()
//	cb := business.NewContactBusiness(cts.ctx, cts.service)
//	existingContacts, err := cts.createContacts(cts.ctx, cb, "+256757592215", "+254957532244", "bwireid@gmail.com")
//	assert.NoError(t, err)
//
//	type args struct {
//		ctx       context.Context
//		contactID string
//		profileID string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *models.Contact
//		wantErr assert.ErrorAssertionFunc
//	}{
//		{
//			name: "Remove existing contact by valid IDs",
//			args: args{
//				ctx:       ctx,
//				contactID: "valid-contact-id",
//				profileID: "valid-profile-id",
//			},
//			want: &models.Contact{
//				ID: "valid-contact-id",
//			},
//			wantErr: assert.NoError,
//		},
//		{
//			name: "Remove contact with invalid contact ID",
//			args: args{
//				ctx:       ctx,
//				contactID: "invalid-contact-id",
//				profileID: "valid-profile-id",
//			},
//			want:    nil,
//			wantErr: assert.Error,
//		},
//		{
//			name: "Remove contact with invalid profile ID",
//			args: args{
//				ctx:       ctx,
//				contactID: "valid-contact-id",
//				profileID: "invalid-profile-id",
//			},
//			want:    nil,
//			wantErr: assert.Error,
//		},
//		{
//			name: "Remove contact with empty IDs",
//			args: args{
//				ctx:       ctx,
//				contactID: "",
//				profileID: "",
//			},
//			want:    nil,
//			wantErr: assert.Error,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			cb := business.NewContactBusiness(ctx, srv)
//			got, err := cb.RemoveContact(tt.args.ctx, tt.args.contactID, tt.args.profileID)
//			if !tt.wantErr(t, err, fmt.Sprintf("RemoveContact(%v, %v, %v)", tt.args.ctx, tt.args.contactID, tt.args.profileID)) {
//				return
//			}
//			assert.Equalf(t, tt.want, got, "RemoveContact(%v, %v, %v)", tt.args.ctx, tt.args.contactID, tt.args.profileID)
//		})
//	}
//}
