package business_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame/tests/testdef"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/apps/default/tests"
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

func (pts *ProfileTestSuite) Test_profileBusiness_GetByContact() {
	t := pts.T()

	pts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := pts.CreateService(t, dep)
		pb := business.NewProfileBusiness(ctx, svc)

		// Create a profile first
		createReq := &profilev1.CreateRequest{
			Type:    profilev1.ProfileType_PERSON,
			Contact: "getbycontact@testing.com",
			Properties: map[string]string{
				"name": "Get By Contact Test",
			},
		}

		profile, err := pb.CreateProfile(ctx, createReq)
		if err != nil {
			t.Errorf("CreateProfile() error = %v", err)
			return
		}

		// Test getting by contact detail
		got, err := pb.GetByContact(ctx, "getbycontact@testing.com")
		if err != nil {
			t.Errorf("GetByContact() error = %v", err)
			return
		}

		if got.GetId() != profile.GetId() {
			t.Errorf("GetByContact() got profile ID = %v, want %v", got.GetId(), profile.GetId())
		}

		// Test getting by contact ID
		if len(profile.GetContacts()) > 0 {
			contactID := profile.GetContacts()[0].GetId()
			got2, getErr := pb.GetByContact(ctx, contactID)
			if getErr != nil {
				t.Errorf("GetByContact() error = %v", getErr)
				return
			}

			if got2.GetId() != profile.GetId() {
				t.Errorf("GetByContact() got profile ID = %v, want %v", got2.GetId(), profile.GetId())
			}
		}
	})
}

func (pts *ProfileTestSuite) Test_profileBusiness_UpdateProfile() {
	t := pts.T()

	pts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := pts.CreateService(t, dep)
		pb := business.NewProfileBusiness(ctx, svc)

		// Create a profile first
		createReq := &profilev1.CreateRequest{
			Type:    profilev1.ProfileType_PERSON,
			Contact: "update@testing.com",
			Properties: map[string]string{
				"name": "Original Name",
			},
		}

		profile, err := pb.CreateProfile(ctx, createReq)
		if err != nil {
			t.Errorf("CreateProfile() error = %v", err)
			return
		}

		// Update the profile
		updateReq := &profilev1.UpdateRequest{
			Id: profile.GetId(),
			Properties: map[string]string{
				"name": "Updated Name",
				"age":  "30",
			},
		}

		updated, err := pb.UpdateProfile(ctx, updateReq)
		if err != nil {
			t.Errorf("UpdateProfile() error = %v", err)
			return
		}

		if updated.GetProperties()["name"] != "Updated Name" {
			t.Errorf(
				"UpdateProfile() name not updated, got = %v, want = %v",
				updated.GetProperties()["name"],
				"Updated Name",
			)
		}

		if updated.GetProperties()["age"] != "30" {
			t.Errorf("UpdateProfile() age not added, got = %v, want = %v", updated.GetProperties()["age"], "30")
		}
	})
}

func (pts *ProfileTestSuite) Test_profileBusiness_MergeProfile() {
	t := pts.T()

	pts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := pts.CreateService(t, dep)
		pb := business.NewProfileBusiness(ctx, svc)

		// Create target profile
		targetReq := &profilev1.CreateRequest{
			Type:    profilev1.ProfileType_PERSON,
			Contact: "target@testing.com",
			Properties: map[string]string{
				"name": "Target Profile",
			},
		}

		target, err := pb.CreateProfile(ctx, targetReq)
		if err != nil {
			t.Errorf("CreateProfile() target error = %v", err)
			return
		}

		// Create profile to merge
		mergeReq := &profilev1.CreateRequest{
			Type:    profilev1.ProfileType_PERSON,
			Contact: "merge@testing.com",
			Properties: map[string]string{
				"age":     "25",
				"country": "Kenya",
			},
		}

		merging, err := pb.CreateProfile(ctx, mergeReq)
		if err != nil {
			t.Errorf("CreateProfile() merging error = %v", err)
			return
		}

		// Merge profiles
		mergeRequest := &profilev1.MergeRequest{
			Id:      target.GetId(),
			Mergeid: merging.GetId(),
		}

		merged, err := pb.MergeProfile(ctx, mergeRequest)
		if err != nil {
			t.Errorf("MergeProfile() error = %v", err)
			return
		}

		// Check merged properties
		if merged.GetProperties()["name"] != "Target Profile" {
			t.Errorf("MergeProfile() original property lost, got = %v", merged.GetProperties()["name"])
		}

		if merged.GetProperties()["age"] != "25" {
			t.Errorf("MergeProfile() merged property not added, got = %v", merged.GetProperties()["age"])
		}

		if merged.GetProperties()["country"] != "Kenya" {
			t.Errorf("MergeProfile() merged property not added, got = %v", merged.GetProperties()["country"])
		}
	})
}

func (pts *ProfileTestSuite) Test_profileBusiness_GetContactByID() {
	t := pts.T()

	pts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := pts.CreateService(t, dep)
		pb := business.NewProfileBusiness(ctx, svc)

		// Create a profile first
		createReq := &profilev1.CreateRequest{
			Type:    profilev1.ProfileType_PERSON,
			Contact: "getcontact@testing.com",
			Properties: map[string]string{
				"name": "Get Contact Test",
			},
		}

		profile, err := pb.CreateProfile(ctx, createReq)
		if err != nil {
			t.Errorf("CreateProfile() error = %v", err)
			return
		}

		if len(profile.GetContacts()) == 0 {
			t.Errorf("CreateProfile() should have created a contact")
			return
		}

		contactID := profile.GetContacts()[0].GetId()
		contact, err := pb.GetContactByID(ctx, contactID)
		if err != nil {
			t.Errorf("GetContactByID() error = %v", err)
			return
		}

		if contact.GetId() != contactID {
			t.Errorf("GetContactByID() got ID = %v, want = %v", contact.GetId(), contactID)
		}

		if contact.GetDetail() != "getcontact@testing.com" {
			t.Errorf("GetContactByID() got detail = %v, want = %v", contact.GetDetail(), "getcontact@testing.com")
		}
	})
}

func (pts *ProfileTestSuite) Test_profileBusiness_VerifyContact() {
	t := pts.T()

	pts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := pts.CreateService(t, dep)
		pb := business.NewProfileBusiness(ctx, svc)

		// Create a profile first
		createReq := &profilev1.CreateRequest{
			Type:    profilev1.ProfileType_PERSON,
			Contact: "verify@testing.com",
			Properties: map[string]string{
				"name": "Verify Contact Test",
			},
		}

		profile, err := pb.CreateProfile(ctx, createReq)
		if err != nil {
			t.Errorf("CreateProfile() error = %v", err)
			return
		}

		if len(profile.GetContacts()) == 0 {
			t.Errorf("CreateProfile() should have created a contact")
			return
		}

		contactID := profile.GetContacts()[0].GetId()
		verificationID, err := pb.VerifyContact(ctx, contactID, "", "123456", 0)
		if err != nil {
			t.Errorf("VerifyContact() error = %v", err)
			return
		}

		if verificationID == "" {
			t.Errorf("VerifyContact() should return verification ID")
		}
	})
}

func (pts *ProfileTestSuite) Test_profileBusiness_CheckVerification_Success() {
	t := pts.T()

	pts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := pts.CreateService(t, dep)
		pb := business.NewProfileBusiness(ctx, svc)

		verificationRepo := repository.NewVerificationRepository(svc)

		// Create a profile and verify contact first
		createReq := &profilev1.CreateRequest{
			Type:    profilev1.ProfileType_PERSON,
			Contact: "test@example.com",
			Properties: map[string]string{
				"name": "Check Verify Test",
			},
		}

		profile, err := pb.CreateProfile(ctx, createReq)
		if err != nil {
			t.Errorf("CreateProfile() error = %v", err)
			return
		}

		if len(profile.GetContacts()) == 0 {
			t.Errorf("CreateProfile() should have created a contact")
			return
		}

		contactID := profile.GetContacts()[0].GetId()
		verificationID, err := pb.VerifyContact(ctx, contactID, "", "123456", 2*time.Second)
		if err != nil {
			t.Errorf("VerifyContact() error = %v", err)
			return
		}

		result, err := tests.WaitForConditionWithResult(ctx, func() (*models.Verification, error) {
			return verificationRepo.GetByID(ctx, verificationID)
		}, 5*time.Second, 100*time.Millisecond)
		if err != nil {
			t.Errorf("VerificationRepo.GetByID() error = %v", err)
			return
		}
		require.NotNil(t, result)
		require.Equal(t, verificationID, result.GetID())

		// Check verification with correct code
		attempts, verified, err := pb.CheckVerification(ctx, verificationID, "123456", "192.168.1.1")
		if err != nil {
			t.Errorf("CheckVerification() error = %v", err)
			return
		}

		if attempts != 1 {
			t.Errorf("CheckVerification() attempts = %v, want = 1", attempts)
		}

		if !verified {
			t.Errorf("CheckVerification() verified = %v, want = true", verified)
		}
	})
}

func (pts *ProfileTestSuite) Test_profileBusiness_CheckVerification_WrongCode() {
	t := pts.T()

	pts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := pts.CreateService(t, dep)
		pb := business.NewProfileBusiness(ctx, svc)
		verificationRepo := repository.NewVerificationRepository(svc)

		// Setup: Create profile and verification
		verificationID := pts.setupVerificationForTest(ctx, t, pb, verificationRepo, "test2@example.com")

		// First attempt with correct code
		attempts, verified, err := pb.CheckVerification(ctx, verificationID, "123456", "192.168.1.1")
		if err != nil {
			t.Errorf("CheckVerification() error = %v", err)
			return
		}

		if attempts != 1 {
			t.Errorf("CheckVerification() attempts = %v, want = 1", attempts)
		}

		if !verified {
			t.Errorf("CheckVerification() verified = %v, want = true", verified)
		}

		// Wait for verification attempts to be processed
		pts.waitForVerificationAttempts(ctx, verificationRepo, verificationID)

		// Check verification with wrong code
		attempts, verified, err = pb.CheckVerification(ctx, verificationID, "wrong", "192.168.1.1")
		if err != nil {
			t.Errorf("CheckVerification() error = %v", err)
			return
		}

		if attempts != 2 {
			t.Errorf("CheckVerification() attempts = %v, want = 2", attempts)
		}

		if verified {
			t.Errorf("CheckVerification() verified = %v, want = false", verified)
		}
	})
}

func (pts *ProfileTestSuite) setupVerificationForTest(
	ctx context.Context,
	t *testing.T,
	pb business.ProfileBusiness,
	verificationRepo repository.VerificationRepository,
	email string,
) string {
	// Create a profile and verify contact first
	createReq := &profilev1.CreateRequest{
		Type:    profilev1.ProfileType_PERSON,
		Contact: email,
		Properties: map[string]string{
			"name": "Check Verify Test",
		},
	}

	profile, err := pb.CreateProfile(ctx, createReq)
	if err != nil {
		t.Errorf("CreateProfile() error = %v", err)
		return ""
	}

	if len(profile.GetContacts()) == 0 {
		t.Errorf("CreateProfile() should have created a contact")
		return ""
	}

	contactID := profile.GetContacts()[0].GetId()
	verificationID, err := pb.VerifyContact(ctx, contactID, "", "123456", 2*time.Second)
	if err != nil {
		t.Errorf("VerifyContact() error = %v", err)
		return ""
	}

	result, err := tests.WaitForConditionWithResult(ctx, func() (*models.Verification, error) {
		return verificationRepo.GetByID(ctx, verificationID)
	}, 5*time.Second, 100*time.Millisecond)
	if err != nil {
		t.Errorf("VerificationRepo.GetByID() error = %v", err)
		return ""
	}
	require.NotNil(t, result)
	require.Equal(t, verificationID, result.GetID())

	return verificationID
}

func (pts *ProfileTestSuite) waitForVerificationAttempts(
	ctx context.Context,
	verificationRepo repository.VerificationRepository,
	verificationID string,
) {
	tests.WaitForConditionWithResult(ctx, func() (*[]*models.VerificationAttempt, error) {
		attemptList, attErr := verificationRepo.GetAttempts(ctx, verificationID)
		if attErr != nil {
			return nil, attErr
		}
		if len(attemptList) == 0 {
			return nil, errors.New("no attempts found yet")
		}
		return &attemptList, nil
	}, 5*time.Second, 100*time.Millisecond)
}

func (pts *ProfileTestSuite) Test_profileBusiness_CreateProfile_EdgeCases() {
	t := pts.T()

	pts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := pts.CreateService(t, dep)
		pb := business.NewProfileBusiness(ctx, svc)

		t.Run("empty contact", func(t *testing.T) {
			createReq := &profilev1.CreateRequest{
				Type:    profilev1.ProfileType_PERSON,
				Contact: "",
				Properties: map[string]string{
					"name": "Invalid Contact Test",
				},
			}

			_, err := pb.CreateProfile(ctx, createReq)
			if err == nil {
				t.Errorf("CreateProfile() should return error for empty contact")
			}
		})

		t.Run("invalid contact format", func(t *testing.T) {
			createReq2 := &profilev1.CreateRequest{
				Type:    profilev1.ProfileType_PERSON,
				Contact: "invalid-contact",
				Properties: map[string]string{
					"name": "Invalid Contact Test 2",
				},
			}

			_, err := pb.CreateProfile(ctx, createReq2)
			if err == nil {
				t.Errorf("CreateProfile() should return error for invalid contact format")
			}
		})
	})
}
