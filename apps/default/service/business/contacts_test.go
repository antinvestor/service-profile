package business_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/apps/default/tests"
)

type ContactTestSuite struct {
	tests.ProfileBaseTestSuite
}

func TestContactSuite(t *testing.T) {
	suite.Run(t, new(ContactTestSuite))
}

func (cts *ContactTestSuite) getContactBusiness(
	ctx context.Context,
	svc *frame.Service,
) (business.ContactBusiness, repository.VerificationRepository) {
	evtsMan := svc.EventsManager()
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	cfg := svc.Config().(*config.ProfileConfig)

	contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
	verificationRepo := repository.NewVerificationRepository(ctx, dbPool, workMan)

	return business.NewContactBusiness(ctx, cfg, evtsMan, contactRepo, verificationRepo), verificationRepo
}

func (cts *ContactTestSuite) TestGeneratePin() {
	t := cts.T()

	type args struct {
		n int
	}
	testCases := []struct {
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
	for _, tt := range testCases {
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
		extra  data.JSONMap
	}
	testCases := []struct {
		name    string
		args    args
		want    *models.Contact
		wantErr require.ErrorAssertionFunc
	}{

		{
			name: "Create contact with valid MSISDN",
			args: args{
				detail: "+256757546244", // Valid MSISDN
				extra:  data.JSONMap{"type": "msisdn"},
			},
			want: &models.Contact{ // Expected result
				Detail:     "+256757546244",
				Properties: data.JSONMap{"type": "msisdn"},
			},
			wantErr: require.NoError, // No error expected
		},
		{
			name: "Create contact with valid Email",
			args: args{
				detail: "test@example.com", // Valid Email
				extra:  data.JSONMap{"type": "email"},
			},
			want: &models.Contact{ // Expected result
				Detail:     "test@example.com",
				Properties: data.JSONMap{"type": "email"},
			},
			wantErr: require.NoError, // No error expected
		},
		{
			name: "Create contact with invalid detail",
			args: args{
				detail: "invalid-detail", // Invalid data, e.g., malformed MSISDN or email
				extra:  data.JSONMap{"type": "unknown"},
			},
			want:    nil,           // Expect no valid contact to be created
			wantErr: require.Error, // Error is expected
		},
		{
			name: "Create contact with empty detail",
			args: args{
				detail: "", // Empty detail
				extra:  data.JSONMap{"type": "email"},
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
				Detail:     "test2@example.com",
				Properties: data.JSONMap{},
			},
			wantErr: require.NoError, // No error expected if type can be inferred or defaults
		},
	}

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				ctx, svc := cts.CreateService(t, dep)

				cb, _ := cts.getContactBusiness(ctx, svc)
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

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := cts.CreateService(t, dep)

		cb, _ := cts.getContactBusiness(ctx, svc)
		existingContacts, err := cts.createContacts(ctx, cb, "+256757546215", "+256757532244", "bwire@gmail.com")
		require.NoError(t, err)

		type args struct {
			detail string
		}
		testCases := []struct {
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

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				var got *models.Contact
				got, err = cb.GetByDetail(ctx, tt.args.detail)
				tt.wantErr(t, err, fmt.Sprintf("GetByDetail(ctx, %v)", tt.args.detail))
				if err != nil {
					return
				}
				require.Equalf(t, tt.want.GetID(), got.GetID(), "GetByDetail(ctx, %v) ID check", tt.args.detail)
				require.Equalf(t, tt.want.Version, got.Version, "GetByDetail(ctx, %v) Version check", tt.args.detail)
				require.Equalf(t, tt.want.Detail, got.Detail, "GetByDetail(ctx, %v) Detail check", tt.args.detail)
				require.Equalf(
					t,
					tt.want.ContactType,
					got.ContactType,
					"GetByDetail(ctx, %v) ContactType check",
					tt.args.detail,
				)
				require.Equalf(t, tt.want.Language, got.Language, "GetByDetail(ctx, %v) Language check", tt.args.detail)
				require.Equalf(
					t,
					tt.want.ProfileID,
					got.ProfileID,
					"GetByDetail(ctx, %v) ProfileID check",
					tt.args.detail,
				)
				require.Equalf(
					t,
					tt.want.Properties,
					got.Properties,
					"GetByDetail(ctx, %v) Properties check",
					tt.args.detail,
				)
			})
		}
	})
}

func (cts *ContactTestSuite) Test_contactBusiness_GetByID() {
	t := cts.T()

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := cts.CreateService(t, dep)

		cb, _ := cts.getContactBusiness(ctx, svc)
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
					require.Equalf(t, tt.want.GetID(), got.GetID(), "GetByID(ctx, %v) ID check", tt.args.contactID)
					require.Equalf(t, tt.want.Detail, got.Detail, "GetByID(ctx, %v) Detail check", tt.args.contactID)
					require.Equalf(
						t,
						tt.want.ContactType,
						got.ContactType,
						"GetByID(ctx, %v) ContactType check",
						tt.args.contactID,
					)
					require.Equalf(
						t,
						tt.want.ProfileID,
						got.ProfileID,
						"GetByID(ctx, %v) ProfileID check",
						tt.args.contactID,
					)
					require.Equalf(
						t,
						tt.want.Properties,
						got.Properties,
						"GetByID(ctx, %v) Properties check",
						tt.args.contactID,
					)
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
	testCases := []struct {
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

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				ctx, svc := cts.CreateService(t, dep)

				cb, _ := cts.getContactBusiness(ctx, svc)
				existingContacts, err := cts.createContacts(ctx, cb, tt.args.contactDetails...)
				require.NoError(t, err)

				for _, contact := range existingContacts {
					_, err = cb.UpdateContact(ctx, contact.GetID(), tt.args.profileID, nil)
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

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := cts.CreateService(t, dep)
		cb, _ := cts.getContactBusiness(ctx, svc)

		// Create a contact first
		contact, err := cb.CreateContact(ctx, "update@testing.com", data.JSONMap{
			"name": "Original Name",
		})
		require.NoError(t, err)
		require.NotNil(t, contact)

		// Update the contact
		updated, err := cb.UpdateContact(ctx, contact.GetID(), util.IDString(), data.JSONMap{
			"name": "Updated Name",
			"age":  "30",
		})
		require.NoError(t, err)
		require.NotNil(t, updated)

		require.Equal(t, "Updated Name", updated.Properties["name"])
		require.Equal(t, "30", updated.Properties["age"])

		// Test updating non-existent contact
		_, err = cb.UpdateContact(ctx, util.IDString(), util.IDString(), data.JSONMap{})
		require.Error(t, err)
	})
}

func (cts *ContactTestSuite) Test_contactBusiness_RemoveContact() {
	t := cts.T()

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := cts.CreateService(t, dep)
		cb, _ := cts.getContactBusiness(ctx, svc)

		// Create a contact with profile ID
		contact, err := cb.CreateContact(ctx, "remove@testing.com", data.JSONMap{})
		require.NoError(t, err)
		require.NotNil(t, contact)

		profileID := util.IDString()
		// Update contact to link to profile
		updated, err := cb.UpdateContact(ctx, contact.GetID(), profileID, data.JSONMap{})
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

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := cts.CreateService(t, dep)
		cb, _ := cts.getContactBusiness(ctx, svc)

		// Create a contact first
		contact, err := cb.CreateContact(ctx, "verify@testing.com", data.JSONMap{})
		require.NoError(t, err)
		require.NotNil(t, contact)

		// Link to profile
		profileID := util.IDString()
		updated, err := cb.UpdateContact(ctx, contact.GetID(), profileID, data.JSONMap{})
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

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := cts.CreateService(t, dep)

		cb, _ := cts.getContactBusiness(ctx, svc)

		// Create a contact and verification using the business layer
		contact, err := cb.CreateContact(ctx, "verify@example.com", data.JSONMap{})
		require.NoError(t, err)

		profileID := util.IDString()
		updated, err := cb.UpdateContact(ctx, contact.GetID(), profileID, data.JSONMap{})
		require.NoError(t, err)

		verification, err := cb.VerifyContact(ctx, updated, "", "123456", 0)
		require.NoError(t, err)

		result, cErr := tests.WaitForConditionWithResult(ctx, func() (*models.Verification, error) {
			v, vErr := cb.GetVerification(ctx, verification.GetID())
			if vErr != nil {
				return nil, vErr
			}
			return v, nil
		}, 5*time.Second, 100*time.Millisecond)
		require.NoError(t, cErr)
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

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := cts.CreateService(t, dep)

		cb, verificationRepo := cts.getContactBusiness(ctx, svc)

		// Create a verification first
		verification := &models.Verification{
			ProfileID: util.IDString(),
			ContactID: util.IDString(),
			Code:      "123456",
		}
		verification.GenID(ctx)

		err := verificationRepo.Create(ctx, verification)
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

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _ := cts.CreateService(t, dep)

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

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := cts.CreateService(t, dep)

		cb, _ := cts.getContactBusiness(ctx, svc)
		profileID := util.IDString()

		// Create multiple contacts for the same profile using valid email formats only
		contact1, err := cb.CreateContact(ctx, "test1@example.com", data.JSONMap{"type": "email"})
		require.NoError(t, err)

		// Update contact with profile ID
		_, err = cb.UpdateContact(ctx, contact1.GetID(), profileID, data.JSONMap{})
		require.NoError(t, err)

		contact2, err := cb.CreateContact(ctx, "test2@example.com", data.JSONMap{"type": "email"})
		require.NoError(t, err)

		// Update contact with profile ID
		_, err = cb.UpdateContact(ctx, contact2.GetID(), profileID, data.JSONMap{})
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

	cts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := cts.CreateService(t, dep)
		cb, _ := cts.getContactBusiness(ctx, svc)

		// Test with empty detail
		contact, err := cb.CreateContact(ctx, "", data.JSONMap{})
		require.Error(t, err)
		require.Nil(t, contact)

		// Test with nil extra map
		contact, err = cb.CreateContact(ctx, "test@example.com", nil)
		require.NoError(t, err)
		require.NotNil(t, contact)
		require.Equal(t, "test@example.com", contact.Detail)

		// Test with empty extra map
		contact, err = cb.CreateContact(ctx, "test2@example.com", data.JSONMap{})
		require.NoError(t, err)
		require.NotNil(t, contact)
		require.Equal(t, "test2@example.com", contact.Detail)
	})
}

func (cts *ContactTestSuite) TestContactTypeFromDetail() {
	t := cts.T()

	type args struct {
		detail string
	}
	testCases := []struct {
		name    string
		args    args
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		// Valid Email Tests
		{
			name:    "Valid simple email",
			args:    args{detail: "test@example.com"},
			want:    profilev1.ContactType_EMAIL.String(),
			wantErr: require.NoError,
		},
		{
			name:    "Valid email with subdomain",
			args:    args{detail: "user@mail.example.com"},
			want:    profilev1.ContactType_EMAIL.String(),
			wantErr: require.NoError,
		},
		{
			name:    "Valid email with numbers",
			args:    args{detail: "user123@example123.com"},
			want:    profilev1.ContactType_EMAIL.String(),
			wantErr: require.NoError,
		},
		{
			name:    "Valid email with special characters",
			args:    args{detail: "user.name+tag@example.com"},
			want:    profilev1.ContactType_EMAIL.String(),
			wantErr: require.NoError,
		},
		{
			name:    "Valid email with hyphen in domain",
			args:    args{detail: "user@my-domain.com"},
			want:    profilev1.ContactType_EMAIL.String(),
			wantErr: require.NoError,
		},
		{
			name:    "Valid email with underscore",
			args:    args{detail: "user_name@example.com"},
			want:    profilev1.ContactType_EMAIL.String(),
			wantErr: require.NoError,
		},
		{
			name:    "Valid email with long TLD",
			args:    args{detail: "user@example.museum"},
			want:    profilev1.ContactType_EMAIL.String(),
			wantErr: require.NoError,
		},

		// Valid Phone Number Tests (MSISDN)
		{
			name:    "Valid international phone number with country code",
			args:    args{detail: "+256757546244"},
			want:    profilev1.ContactType_MSISDN.String(),
			wantErr: require.NoError,
		},
		{
			name:    "Random phone number",
			args:    args{detail: "+12345678900"},
			want:    profilev1.ContactType_MSISDN.String(),
			wantErr: require.NoError,
		},
		{
			name:    "Valid US phone number",
			args:    args{detail: "+12025551234"},
			want:    profilev1.ContactType_MSISDN.String(),
			wantErr: require.NoError,
		},
		{
			name:    "Valid UK phone number",
			args:    args{detail: "+442071234567"},
			want:    profilev1.ContactType_MSISDN.String(),
			wantErr: require.NoError,
		},
		{
			name:    "Valid Kenya phone number",
			args:    args{detail: "+254701234567"},
			want:    profilev1.ContactType_MSISDN.String(),
			wantErr: require.NoError,
		},
		{
			name:    "Valid phone number with spaces (should be normalized)",
			args:    args{detail: "+1 202 555 1234"},
			want:    profilev1.ContactType_MSISDN.String(),
			wantErr: require.NoError,
		},

		// Invalid Email Tests
		{
			name:    "Invalid email - missing @",
			args:    args{detail: "testexample.com"},
			want:    "",
			wantErr: require.Error,
		},
		{
			name:    "Invalid email - missing domain",
			args:    args{detail: "test@"},
			want:    "",
			wantErr: require.Error,
		},
		{
			name:    "Invalid email - missing local part",
			args:    args{detail: "@example.com"},
			want:    "",
			wantErr: require.Error,
		},
		{
			name:    "Invalid email - double @",
			args:    args{detail: "test@@example.com"},
			want:    "",
			wantErr: require.Error,
		},
		{
			name:    "Invalid email - spaces",
			args:    args{detail: "test @example.com"},
			want:    "",
			wantErr: require.Error,
		},
		{
			name:    "Invalid email - invalid characters",
			args:    args{detail: "test<>@example.com"},
			want:    "",
			wantErr: require.Error,
		},

		// Invalid Phone Number Tests
		{
			name:    "Invalid phone - too short",
			args:    args{detail: "+1234"},
			want:    "",
			wantErr: require.Error,
		},
		{
			name:    "Invalid phone - no country code",
			args:    args{detail: "1234567890"},
			want:    "",
			wantErr: require.Error,
		},
		{
			name:    "Invalid phone - letters",
			args:    args{detail: "+1abc2345678"},
			want:    "",
			wantErr: require.Error,
		},
		{
			name:    "Invalid phone - too many digits",
			args:    args{detail: "+123456789012345678901"},
			want:    "",
			wantErr: require.Error,
		},

		// Edge Cases
		{
			name:    "Empty string",
			args:    args{detail: ""},
			want:    "",
			wantErr: require.Error,
		},
		{
			name:    "Whitespace only",
			args:    args{detail: "   "},
			want:    "",
			wantErr: require.Error,
		},
		{
			name:    "Random text",
			args:    args{detail: "random-text-123"},
			want:    "",
			wantErr: require.Error,
		},
		{
			name:    "URL-like string",
			args:    args{detail: "http://example.com"},
			want:    "",
			wantErr: require.Error,
		},
		{
			name:    "Number without plus sign",
			args:    args{detail: "256757546244"},
			want:    "",
			wantErr: require.Error,
		},
		{
			name:    "Email-like but invalid TLD",
			args:    args{detail: "test@example"},
			want:    profilev1.ContactType_EMAIL.String(),
			wantErr: require.NoError,
		},
		{
			name:    "Phone with invalid country code",
			args:    args{detail: "+999999999999"},
			want:    "",
			wantErr: require.Error,
		},

		// Boundary Cases
		{
			name: "Very long email",
			args: args{
				detail: "verylongemailaddressthatmightcauseproblemswithvalidation@verylongdomainnamethatmightalsocauseissues.com",
			},
			want:    profilev1.ContactType_EMAIL.String(),
			wantErr: require.NoError,
		},
		{
			name:    "Minimum valid email",
			args:    args{detail: "a@b.co"},
			want:    profilev1.ContactType_EMAIL.String(),
			wantErr: require.NoError,
		},
		{
			name:    "Email with all allowed special chars",
			args:    args{detail: "test.email+tag@example-domain.co.uk"},
			want:    profilev1.ContactType_EMAIL.String(),
			wantErr: require.NoError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := business.ContactTypeFromDetail(ctx, tt.args.detail)
			tt.wantErr(t, err, fmt.Sprintf("ContactTypeFromDetail(ctx, %v)", tt.args.detail))
			require.Equalf(t, tt.want, got, "ContactTypeFromDetail(ctx, %v)", tt.args.detail)
		})
	}
}
