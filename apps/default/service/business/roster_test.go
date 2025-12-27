package business_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/frame/security"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/apps/default/tests"
)

type RosterTestSuite struct {
	tests.ProfileBaseTestSuite
}

func TestRosterSuite(t *testing.T) {
	suite.Run(t, new(RosterTestSuite))
}

// Helper function to create consistent test DEK
func createRosterTestDEK(cfg *config.ProfileConfig) *config.DEK {
	// Decode base64 keys
	key, err := base64.StdEncoding.DecodeString(cfg.DEKActiveAES256GCMKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode DEKActiveAES256GCMKey: %v", err))
	}
	lookupKey, err := base64.StdEncoding.DecodeString(cfg.DEKLookupTokenHMACSHA256Key)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode DEKLookupTokenHMACSHA256Key: %v", err))
	}
	
	return &config.DEK{
		KeyID:     cfg.DEKActiveKeyID,
		Key:       key,
		OldKeyID:  "old-key-id",
		OldKey:    []byte("1234567890123456"), // 16 bytes for old key
		LookUpKey: lookupKey,
	}
}

// Helper function to create consistent test DEK for contacts
func createRosterContactTestDEK(cfg *config.ProfileConfig) *config.DEK {
	// Decode base64 keys
	key, err := base64.StdEncoding.DecodeString(cfg.DEKActiveAES256GCMKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode DEKActiveAES256GCMKey: %v", err))
	}
	lookupKey, err := base64.StdEncoding.DecodeString(cfg.DEKLookupTokenHMACSHA256Key)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode DEKLookupTokenHMACSHA256Key: %v", err))
	}
	
	return &config.DEK{
		KeyID:     cfg.DEKActiveKeyID,
		Key:       key,
		OldKeyID:  "old-key-id",
		OldKey:    []byte("1234567890123456"), // 16 bytes for old key
		LookUpKey: lookupKey,
	}
}

func (rts *RosterTestSuite) getRosterBusiness(ctx context.Context, svc *frame.Service) business.RosterBusiness {
	evtsMan := svc.EventsManager()
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	cfg := svc.Config().(*config.ProfileConfig)

	contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
	verificationRepo := repository.NewVerificationRepository(ctx, dbPool, workMan)

	contactBusiness := business.NewContactBusiness(ctx, cfg, createRosterContactTestDEK(cfg), evtsMan, contactRepo, verificationRepo)

	rosterRepo := repository.NewRosterRepository(ctx, dbPool, workMan)
	return business.NewRosterBusiness(ctx, cfg, createRosterContactTestDEK(cfg), contactBusiness, rosterRepo)
}

func (rts *RosterTestSuite) getContactBusiness(
	ctx context.Context,
	svc *frame.Service,
) (business.ContactBusiness, repository.VerificationRepository) {
	evtsMan := svc.EventsManager()
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	cfg := svc.Config().(*config.ProfileConfig)

	contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
	verificationRepo := repository.NewVerificationRepository(ctx, dbPool, workMan)

	return business.NewContactBusiness(ctx, cfg, createRosterContactTestDEK(cfg), evtsMan, contactRepo, verificationRepo), verificationRepo
}

func (rts *RosterTestSuite) createRoster(
	ctx context.Context,
	rb business.RosterBusiness,
	profileID string,
	contacts map[string]data.JSONMap,
) (map[string]*profilev1.RosterObject, error) {
	var requestData []*profilev1.RawContact
	for detail, extra := range contacts {
		requestData = append(requestData, &profilev1.RawContact{
			Contact: detail,
			Extras:  extra.ToProtoStruct(),
		})
	}

	claims := security.ClaimsFromMap(map[string]string{
		"sub":          profileID,
		"tenant_id":    "tenantx",
		"partition_id": "party",
	})

	ctx = claims.ClaimsToContext(ctx)

	rosterList, err := rb.CreateRoster(ctx, &profilev1.AddRosterRequest{Data: requestData})
	if err != nil {
		return nil, err
	}

	result := map[string]*profilev1.RosterObject{}
	for _, roster := range rosterList {
		result[roster.GetContact().GetDetail()] = roster
	}

	return result, nil
}

func (rts *RosterTestSuite) TestRosterBusiness_ToApi() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)

		// Create a real contact using the business layer
		cb, _ := rts.getContactBusiness(ctx, svc)
		contact, err := cb.CreateContact(ctx, "+256757546244", data.JSONMap{"type": "msisdn"})
		require.NoError(t, err)
		require.NotNil(t, contact)

		roster := &models.Roster{
			Contact:    contact,
			ProfileID:  "profile123",
			Properties: data.JSONMap{"key1": "value1"},
		}

		// Use the same DEK that was used to create the contact
		dek := createRosterTestDEK(svc.Config().(*config.ProfileConfig))
		result, err := roster.ToAPI(dek)
		require.NoError(t, err)
		require.Equal(t, "profile123", result.GetProfileId(), "Profile ID should match")
	})
}

func (rts *RosterTestSuite) TestRosterBusiness_GetByID() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)
		rb := rts.getRosterBusiness(ctx, svc)

		rosterMap, err := rts.createRoster(ctx, rb, "profile123", map[string]data.JSONMap{
			"roster@test.com": {"key1": "value1"},
		})
		require.NoError(t, err)

		expectedRoster := rosterMap["roster@test.com"]
		rosterID := expectedRoster.GetId()

		result, err := rb.GetByID(ctx, rosterID)

		require.NoError(t, err)
		require.Empty(
			t,
			expectedRoster.GetProfileId(),
			"Profile ID should be empty as the contact belongs to another not in the system",
		)
		require.Equal(t, "profile123", result.ProfileID, "Profile ID should match")
		require.Equal(t, expectedRoster.GetContact().GetId(), result.ContactID, "Contact ID should match")
	})
}

func (rts *RosterTestSuite) TestRosterBusiness_RemoveRoster() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)
		rb := rts.getRosterBusiness(ctx, svc)

		rosterMap, err := rts.createRoster(ctx, rb, "profRemov123", map[string]data.JSONMap{
			"rosterremove@test.com": {"key1": "value1"},
		})
		require.NoError(t, err)

		roster := rosterMap["rosterremove@test.com"]
		rosterID := roster.GetId()

		result, err := rb.RemoveRoster(ctx, rosterID)

		require.NoError(t, err)
		require.Equal(t, roster.GetProfileId(), result.GetProfileId())

		_, err = rb.GetByID(ctx, rosterID)

		require.Error(t, err)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound, "Error should be 'gorm.ErrRecordNotFound'")
	})
}

func (rts *RosterTestSuite) TestRosterBusiness_RemoveRoster_NotFound() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)
		rb := rts.getRosterBusiness(ctx, svc)

		rosterID := "nonexistent"

		result, err := rb.RemoveRoster(ctx, rosterID)

		require.Error(t, err, "A not found error should be returned")
		require.Nil(t, result, "Result should be nil")
		require.ErrorIs(t, err, gorm.ErrRecordNotFound, "Error should be 'gorm.ErrRecordNotFound'")
	})
}

func (rts *RosterTestSuite) TestRosterBusiness_Search() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)
		rb := rts.getRosterBusiness(ctx, svc)

		profileID := "searchProfileC1"

		_, err := rts.createRoster(ctx, rb, profileID, map[string]data.JSONMap{
			"rostersearch@test.com":        {"name": "John Osogo", "age": "33"},
			"searchrostercontact@test.com": {"name": "Thomas Balindhe", "age": "36"},
			"+256755718293":                {"name": "Mary Osogo", "age": "21"},
			"+256755718291":                {"name": "Julius Search Best", "age": "51"},
		})
		require.NoError(t, err)

		testCases := []struct {
			name        string
			request     *profilev1.SearchRosterRequest
			wantError   require.ErrorAssertionFunc
			profileID   string
			resultCount int
		}{
			{
				name: "Valid search request with matching data",
				request: &profilev1.SearchRosterRequest{
					Query: "searchrostercontact@test.com",
				},
				profileID:   profileID,
				wantError:   require.NoError,
				resultCount: 2,
			},
			{
				name: "Valid search request partial contact match",
				request: &profilev1.SearchRosterRequest{
					Query: "+25675571829",
				},
				profileID:   profileID,
				wantError:   require.NoError,
				resultCount: 2,
			},
			{
				name: "Valid search request partial properties match",
				request: &profilev1.SearchRosterRequest{
					Query: "Osogo",
				},
				profileID:   profileID,
				wantError:   require.NoError,
				resultCount: 2,
			},
			{
				name: "Valid cross field partial match",
				request: &profilev1.SearchRosterRequest{
					Query: "search",
					// Properties: []string{"name", "age"},
				},
				profileID:   profileID,
				wantError:   require.NoError,
				resultCount: 3,
			},
			{
				name: "Request with non-existent profile ID",
				request: &profilev1.SearchRosterRequest{
					Query: "non existent profile",
					// Properties: []string{"name", "age"},
				},
				profileID:   "nonExistentProfileId",
				wantError:   require.NoError,
				resultCount: 0,
			},
			{
				name: "Empty search request",
				request: &profilev1.SearchRosterRequest{
					Query: "",
					// Properties: []string{"name", "age"},
				},
				profileID:   profileID,
				wantError:   require.NoError,
				resultCount: 4,
			},
			{
				name: "Empty search request with wrong profile",
				request: &profilev1.SearchRosterRequest{
					Query: "",
					// Properties: []string{"name", "age"},
				},
				profileID:   "funnyProfileId",
				wantError:   require.NoError,
				resultCount: 0,
			},
		}

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				claims := security.ClaimsFromMap(map[string]string{
					"sub":          tt.profileID,
					"tenant_id":    "tenantx",
					"partition_id": "party",
				})

				ctxWithClaims := claims.ClaimsToContext(ctx)

				jobResult, err0 := rb.Search(ctxWithClaims, tt.request)
				require.NoError(t, err0)

				var rosterList []*models.Roster
				for result := range jobResult.ResultChan() {
					if result == nil || result.IsError() {
						break
					}
					rosterList = append(rosterList, result.Item()...)
				}

				require.Len(t, rosterList, tt.resultCount, "Roster count mismatch")
			})
		}
	})
}
