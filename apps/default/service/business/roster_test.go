package business_test

import (
	"context"
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

func (rts *RosterTestSuite) getRosterBusiness(ctx context.Context, svc *frame.Service) business.RosterBusiness {
	evtsMan := svc.EventsManager()
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	cfg := svc.Config().(*config.ProfileConfig)

	contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
	verificationRepo := repository.NewVerificationRepository(ctx, dbPool, workMan)

	contactBusiness := business.NewContactBusiness(ctx, cfg, evtsMan, contactRepo, verificationRepo)

	rosterRepo := repository.NewRosterRepository(ctx, dbPool, workMan)
	return business.NewRosterBusiness(ctx, contactBusiness, rosterRepo)
}

func (rts *RosterTestSuite) createRoster(
	ctx context.Context,
	rb business.RosterBusiness,
	profileID string,
	contacts map[string]data.JSONMap,
) (map[string]*profilev1.RosterObject, error) {
	var requestData []*profilev1.AddContactRequest
	for detail, extra := range contacts {
		requestData = append(requestData, &profilev1.AddContactRequest{
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

	rts.WithTestDependancies(t, func(t *testing.T, _ *definition.DependencyOption) {
		contact := &models.Contact{
			Detail:             "+256757546244",
			ProfileID:          "ownersId123",
			ContactType:        profilev1.ContactType_MSISDN.String(),
			CommunicationLevel: profilev1.CommunicationLevel_ALL.String(),
		}
		roster := &models.Roster{
			Contact:    contact,
			ProfileID:  "profile123",
			Properties: data.JSONMap{"key1": "value1"},
		}

		result := roster.ToAPI()
		require.Equal(t, "ownersId123", result.GetProfileId(), "Profile ID should match")
	})
}

func (rts *RosterTestSuite) TestRosterBusiness_GetByID() {
	t := rts.T()

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := rts.CreateService(t, dep)
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
		svc, ctx := rts.CreateService(t, dep)
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
		svc, ctx := rts.CreateService(t, dep)
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
		svc, ctx := rts.CreateService(t, dep)
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
					Query:      "searchrostercontact@test.com",
					Properties: []string{"name", "age"},
				},
				profileID:   profileID,
				wantError:   require.NoError,
				resultCount: 1,
			},
			{
				name: "Valid search request partial contact match",
				request: &profilev1.SearchRosterRequest{
					Query:      "+25675571829",
					Properties: []string{"name", "age"},
				},
				profileID:   profileID,
				wantError:   require.NoError,
				resultCount: 2,
			},
			{
				name: "Valid search request partial properties match",
				request: &profilev1.SearchRosterRequest{
					Query:      "Osogo",
					Properties: []string{"name", "age"},
				},
				profileID:   profileID,
				wantError:   require.NoError,
				resultCount: 2,
			},
			{
				name: "Valid cross field partial match",
				request: &profilev1.SearchRosterRequest{
					Query:      "search",
					Properties: []string{"name", "age"},
				},
				profileID:   profileID,
				wantError:   require.NoError,
				resultCount: 3,
			},
			{
				name: "Request with non-existent profile ID",
				request: &profilev1.SearchRosterRequest{
					Query:      "non existent profile",
					Properties: []string{"name", "age"},
				},
				profileID:   "nonExistentProfileId",
				wantError:   require.NoError,
				resultCount: 0,
			},
			{
				name: "Empty search request",
				request: &profilev1.SearchRosterRequest{
					Query:      "",
					Properties: []string{"name", "age"},
				},
				profileID:   profileID,
				wantError:   require.NoError,
				resultCount: 4,
			},
			{
				name: "Empty search request with wrong profile",
				request: &profilev1.SearchRosterRequest{
					Query:      "",
					Properties: []string{"name", "age"},
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
				for {
					result, ok := jobResult.ReadResult(ctx)
					if result == nil || result.IsError() || !ok {
						break
					}
					rosterList = append(rosterList, result.Item()...)
				}

				require.Len(t, rosterList, tt.resultCount, "Roster count mismatch")
			})
		}
	})
}
