package business_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/apps/default/tests"
)

type ContactTestSuite struct {
	tests.BaseTestSuite
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
			newPin := util.RandomString(tt.args.n)
			require.Lenf(t, newPin, tt.want, "GeneratePin(%v)", tt.args.n)
		})
	}
}

func (cts *ContactTestSuite) Test_contactBusiness_CreateContact() {
	t := cts.T()

	type args struct {
		detail string
		extra  map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Contact
		wantErr require.ErrorAssertionFunc
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
			wantErr: require.NoError, // No error expected
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
			wantErr: require.NoError, // No error expected
		},
		{
			name: "Create contact with invalid detail",
			args: args{
				detail: "invalid-detail", // Invalid data, e.g., malformed MSISDN or email
				extra:  map[string]string{"type": "unknown"},
			},
			want:    nil,           // Expect no valid contact to be created
			wantErr: require.Error, // Error is expected
		},
		{
			name: "Create contact with empty detail",
			args: args{
				detail: "", // Empty detail
				extra:  map[string]string{"type": "email"},
			},
			want:    nil,           // Contact should not be created with empty details
			wantErr: require.Error, // Error is expected
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
			wantErr: require.NoError, // No error expected if type can be inferred or defaults
		},
	}

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				svc, ctx := cts.CreateService(t, dep)

				cb := business.NewContactBusiness(ctx, svc)
				got, err := cb.CreateContact(ctx, tt.args.detail, tt.args.extra)
				tt.wantErr(t, err, fmt.Sprintf("CreateContact(ctx, %v, %v)", tt.args.detail, tt.args.extra))

				if tt.want != nil && got != nil {
					require.Equalf(
						t,
						tt.want.Detail,
						got.Detail,
						"CreateContact(ctx, %v, %v)",
						tt.args.detail,
						tt.args.extra,
					)
					require.Equalf(
						t,
						tt.want.Properties,
						got.Properties,
						"CreateContact(ctx, %v, %v)",
						tt.args.detail,
						tt.args.extra,
					)
				}
			})
		}
	})
}

func (cts *ContactTestSuite) createContacts(
	ctx context.Context,
	cb business.ContactBusiness,
	contacDetails ...string,
) (map[string]*models.Contact, error) {
	result := map[string]*models.Contact{}
	for _, detail := range contacDetails {
		contact, err := cb.CreateContact(ctx, detail, nil)
		if err != nil {
			return nil, err
		}
		result[contact.Detail] = contact
	}
	return result, nil
}

func (cts *ContactTestSuite) Test_contactBusiness_GetByDetail() {
	t := cts.T()

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := cts.CreateService(t, dep)

		cb := business.NewContactBusiness(ctx, svc)
		existingContacts, err := cts.createContacts(ctx, cb, "+256757546215", "+256757532244", "bwire@gmail.com")
		require.NoError(t, err)

		type args struct {
			detail string
		}
		tests := []struct {
			name    string
			args    args
			want    *models.Contact
			wantErr require.ErrorAssertionFunc
		}{
			{
				name: "Get contact by valid MSISDN",
				args: args{
					detail: "+256757546215",
				},
				want:    existingContacts["+256757546215"],
				wantErr: require.NoError,
			},
			{
				name: "Get contact by valid Email",
				args: args{
					detail: "bwire@gmail.com",
				},
				want:    existingContacts["bwire@gmail.com"],
				wantErr: require.NoError,
			},
			{
				name: "Get contact by invalid detail",
				args: args{
					detail: "invalid-detail",
				},
				want:    nil,
				wantErr: require.Error,
			},
			{
				name: "Get contact by empty detail",
				args: args{
					detail: "",
				},
				want:    nil,
				wantErr: require.Error,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var got *models.Contact
				got, err = cb.GetByDetail(ctx, tt.args.detail)
				tt.wantErr(t, err, fmt.Sprintf("GetByDetail(%v, %v)", ctx, tt.args.detail))
				require.Equalf(t, tt.want, got, "GetByDetail(%v, %v)", ctx, tt.args.detail)
			})
		}
	})
}

func (cts *ContactTestSuite) Test_contactBusiness_GetByID() {
	t := cts.T()

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := cts.CreateService(t, dep)

		cb := business.NewContactBusiness(ctx, svc)
		existingContacts, err := cts.createContacts(ctx, cb, "+256757592215", "+254757532244", "bwireid@gmail.com")
		require.NoError(t, err)

		type args struct {
			contactID string
		}
		testCases := []struct {
			name    string
			args    args
			want    *models.Contact
			wantErr require.ErrorAssertionFunc
		}{
			{
				name: "Get contact with valid ID",
				args: args{
					contactID: existingContacts["+256757592215"].GetID(),
				},
				want:    existingContacts["+256757592215"],
				wantErr: require.NoError,
			},
			{
				name: "Get contact with invalid ID",
				args: args{
					contactID: "invalid-contact-id",
				},
				want:    nil,
				wantErr: require.Error,
			},
			{
				name: "Get contact with empty ID",
				args: args{
					contactID: "",
				},
				want:    nil,
				wantErr: require.Error,
			},
		}

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				got, err0 := cb.GetByID(ctx, tt.args.contactID)
				tt.wantErr(t, err0, fmt.Sprintf("GetByID(ctx, %v)", tt.args.contactID))
				if tt.want != nil {
					require.Equalf(t, tt.want, got, "GetByID(ctx, %v)", tt.args.contactID)
				}
			})
		}
	})
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
		wantErr require.ErrorAssertionFunc
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
			wantErr: require.NoError,
		},
		{
			name: "Get contacts by invalid profile ID",
			args: args{
				profileID: "invalid-profile-id",
			},
			want:    []*models.Contact{},
			wantErr: require.NoError,
		},
		{
			name: "Get contacts by empty profile ID",
			args: args{
				profileID: "",
			},
			want:    []*models.Contact{},
			wantErr: require.Error,
		},
	}

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				svc, ctx := cts.CreateService(t, dep)

				cb := business.NewContactBusiness(ctx, svc)
				existingContacts, err := cts.createContacts(ctx, cb, tt.args.contactDetails...)
				require.NoError(t, err)

				for _, contact := range existingContacts {
					_, err = cb.UpdateContact(ctx, contact.GetID(), tt.args.profileID, map[string]string{})
					require.NoError(t, err)
				}

				var got []*models.Contact
				got, err = cb.GetByProfile(ctx, tt.args.profileID)
				tt.wantErr(t, err, fmt.Sprintf("GetByProfile(%v, %v)", ctx, tt.args.profileID))

				require.Lenf(t, got, len(tt.want), "GetByProfile(%v, %v)", ctx, tt.args.profileID)
			})
		}
	})
}

func (cts *ContactTestSuite) Test_contactBusiness_UpdateContact() {
	t := cts.T()

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := cts.CreateService(t, dep)
		cb := business.NewContactBusiness(ctx, svc)

		// Create a contact first
		contact, err := cb.CreateContact(ctx, "update@testing.com", map[string]string{
			"name": "Original Name",
		})
		require.NoError(t, err)
		require.NotNil(t, contact)

		// Update the contact
		updated, err := cb.UpdateContact(ctx, contact.GetID(), util.IDString(), map[string]string{
			"name": "Updated Name",
			"age":  "30",
		})
		require.NoError(t, err)
		require.NotNil(t, updated)
		require.Equal(t, "Updated Name", updated.Properties["name"])
		require.Equal(t, "30", updated.Properties["age"])

		// Test updating non-existent contact
		_, err = cb.UpdateContact(ctx, util.IDString(), util.IDString(), map[string]string{})
		require.Error(t, err)
	})
}

func (cts *ContactTestSuite) Test_contactBusiness_RemoveContact() {
	t := cts.T()

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := cts.CreateService(t, dep)
		cb := business.NewContactBusiness(ctx, svc)

		// Create a contact with profile ID
		contact, err := cb.CreateContact(ctx, "remove@testing.com", map[string]string{})
		require.NoError(t, err)
		require.NotNil(t, contact)

		profileID := util.IDString()
		// Update contact to link to profile
		updated, err := cb.UpdateContact(ctx, contact.GetID(), profileID, map[string]string{})
		require.NoError(t, err)
		require.Equal(t, profileID, updated.ProfileID)

		// Remove the contact
		removed, err := cb.RemoveContact(ctx, contact.GetID(), profileID)
		require.NoError(t, err)
		require.NotNil(t, removed)
		require.Empty(t, removed.ProfileID)

		// Test removing with wrong profile ID
		_, err = cb.RemoveContact(ctx, contact.GetID(), util.IDString())
		require.Error(t, err)

		// Test removing non-existent contact
		_, err = cb.RemoveContact(ctx, util.IDString(), profileID)
		require.Error(t, err)
	})
}

func (cts *ContactTestSuite) Test_contactBusiness_VerifyContact() {
	t := cts.T()

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := cts.CreateService(t, dep)
		cb := business.NewContactBusiness(ctx, svc)

		// Create a contact first
		contact, err := cb.CreateContact(ctx, "verify@testing.com", map[string]string{})
		require.NoError(t, err)
		require.NotNil(t, contact)

		// Link to profile
		profileID := util.IDString()
		updated, err := cb.UpdateContact(ctx, contact.GetID(), profileID, map[string]string{})
		require.NoError(t, err)

		// Verify contact
		verification, err := cb.VerifyContact(ctx, updated, "", "123456", 0)
		require.NoError(t, err)
		require.NotNil(t, verification)
		require.Equal(t, "123456", verification.Code)
		require.Equal(t, profileID, verification.ProfileID)
		require.Equal(t, updated.GetID(), verification.ContactID)

		// Test with custom verification ID
		customVerificationID := util.IDString()
		verification2, err := cb.VerifyContact(ctx, updated, customVerificationID, "654321", 0)
		require.NoError(t, err)
		require.NotNil(t, verification2)
		require.Equal(t, customVerificationID, verification2.GetID())

		// Test with nil contact
		_, err = cb.VerifyContact(ctx, nil, "", "123456", 0)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no contact specified")
	})
}

// Temporarily commented out due to nil pointer dereference in business layer.
func (cts *ContactTestSuite) Test_contactBusiness_GetVerification() {
	t := cts.T()

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := cts.CreateService(t, dep)

		cb := business.NewContactBusiness(ctx, svc)

		// Create a contact and verification using the business layer
		contact, err := cb.CreateContact(ctx, "verify@example.com", map[string]string{})
		require.NoError(t, err)

		profileID := util.IDString()
		updated, err := cb.UpdateContact(ctx, contact.GetID(), profileID, map[string]string{})
		require.NoError(t, err)

		verification, err := cb.VerifyContact(ctx, updated, "", "123456", 0)
		require.NoError(t, err)

		result, err := tests.WaitForConditionWithResult(ctx, func() (*models.Verification, error) {
			return cb.GetVerification(ctx, verification.GetID())
		}, 5*time.Second, 100*time.Millisecond)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, verification.GetID(), result.GetID())

		// Test with non-existent verification
		result, err = cb.GetVerification(ctx, util.IDString())
		require.Error(t, err) // Should return error for non-existent verification
		require.Empty(t, result.GetID())
	})
}

// Temporarily commented out due to nil pointer dereference in verification repository.
func (cts *ContactTestSuite) Test_contactBusiness_GetVerificationAttempts() {
	t := cts.T()

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := cts.CreateService(t, dep)

		cb := business.NewContactBusiness(ctx, svc)

		// Create a verification first
		verification := &models.Verification{
			ProfileID: util.IDString(),
			ContactID: util.IDString(),
			Code:      "123456",
		}
		verification.GenID(ctx)

		verificationRepo := repository.NewVerificationRepository(svc)
		err := verificationRepo.Save(ctx, verification)
		require.NoError(t, err)

		// Note: Using verification repository to save attempts since VerificationAttemptRepository may not exist
		// This is a simplified test approach

		// Test getting verification attempts
		attempts, err := cb.GetVerificationAttempts(ctx, verification.GetID())
		require.NoError(t, err)
		require.NotNil(t, attempts) // May be empty if no attempts saved

		// Test with non-existent verification
		attempts, err = cb.GetVerificationAttempts(ctx, util.IDString())
		require.NoError(t, err)
		require.Empty(t, attempts)
	})
}

func (cts *ContactTestSuite) Test_contactBusiness_ToAPI() {
	t := cts.T()

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		_, ctx := cts.CreateService(t, dep)

		// Create a contact
		contact := &models.Contact{
			Detail:             "test@example.com",
			ContactType:        profilev1.ContactType_name[int32(profilev1.ContactType_EMAIL)],
			CommunicationLevel: "primary",
			ProfileID:          util.IDString(),
		}
		contact.GenID(ctx)

		// Test ToAPI conversion
		apiContact := contact.ToAPI(false)
		require.NotNil(t, apiContact)
		require.Equal(t, contact.Detail, apiContact.GetDetail())
		require.Equal(t, contact.GetID(), apiContact.GetId())

		// Test with partial flag
		apiContact = contact.ToAPI(true)
		require.NotNil(t, apiContact)
	})
}

func (cts *ContactTestSuite) Test_contactBusiness_GetByProfile_Extended() {
	t := cts.T()

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := cts.CreateService(t, dep)

		cb := business.NewContactBusiness(ctx, svc)
		profileID := util.IDString()

		// Create multiple contacts for the same profile using valid email formats only
		contact1, err := cb.CreateContact(ctx, "test1@example.com", map[string]string{"type": "email"})
		require.NoError(t, err)

		// Update contact with profile ID
		_, err = cb.UpdateContact(ctx, contact1.GetID(), profileID, map[string]string{})
		require.NoError(t, err)

		contact2, err := cb.CreateContact(ctx, "test2@example.com", map[string]string{"type": "email"})
		require.NoError(t, err)

		// Update contact with profile ID
		_, err = cb.UpdateContact(ctx, contact2.GetID(), profileID, map[string]string{})
		require.NoError(t, err)

		// Test getting contacts by profile
		contacts, err := cb.GetByProfile(ctx, profileID)
		require.NoError(t, err)
		require.Len(t, contacts, 2)

		// Test with non-existent profile
		contacts, err = cb.GetByProfile(ctx, util.IDString())
		require.NoError(t, err)
		require.Empty(t, contacts)
	})
}

func (cts *ContactTestSuite) Test_contactBusiness_CreateContact_EdgeCases() {
	t := cts.T()

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := cts.CreateService(t, dep)
		cb := business.NewContactBusiness(ctx, svc)

		// Test with empty detail
		contact, err := cb.CreateContact(ctx, "", map[string]string{})
		require.Error(t, err)
		require.Nil(t, contact)

		// Test with nil extra map
		contact, err = cb.CreateContact(ctx, "test@example.com", nil)
		require.NoError(t, err)
		require.NotNil(t, contact)
		require.Equal(t, "test@example.com", contact.Detail)

		// Test with empty extra map
		contact, err = cb.CreateContact(ctx, "test2@example.com", map[string]string{})
		require.NoError(t, err)
		require.NotNil(t, contact)
		require.Equal(t, "test2@example.com", contact.Detail)
	})
}
