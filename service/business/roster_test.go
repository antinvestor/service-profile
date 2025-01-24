package business_test

import (
	"context"
	"errors"
	"github.com/antinvestor/service-profile/service/business"
	"github.com/pitabwire/frame"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type RosterTestSuite struct {
	BaseTestSuite
}

func TestRosterSuite(t *testing.T) {
	suite.Run(t, new(RosterTestSuite))
}

func (rts *RosterTestSuite) createRoster(ctx context.Context, rb business.RosterBusiness, profileID string, contacts map[string]map[string]string) (map[string]*profilev1.RosterObject, error) {
	result := map[string]*profilev1.RosterObject{}
	var requestData []*profilev1.AddContactRequest
	for detail, extra := range contacts {

		requestData = append(requestData, &profilev1.AddContactRequest{
			Contact: detail,
			Extras:  extra,
		})
	}

	claims := frame.ClaimsFromMap(map[string]string{
		"sub":          profileID,
		"tenant_id":    "tenantx",
		"partition_id": "party",
	})

	ctx = claims.ClaimsToContext(ctx)

	rosterList, err := rb.CreateRoster(ctx, &profilev1.AddRosterRequest{Data: requestData})
	if err != nil {
		return nil, err
	}

	for _, roster := range rosterList {
		result[roster.GetContact().GetDetail()] = roster
	}

	return result, nil
}

func (rts *RosterTestSuite) TestRosterBusiness_ToApi() {

	t := rts.T()
	rb := business.NewRosterBusiness(rts.ctx, rts.service)

	contact := &models.Contact{
		Detail:             "+256757546244",
		ContactType:        profilev1.ContactType_MSISDN.String(),
		CommunicationLevel: profilev1.CommunicationLevel_ALL.String(),
	}
	roster := &models.Roster{
		Contact:    contact,
		ProfileID:  "profile123",
		Properties: map[string]interface{}{"key1": "value1"},
	}

	result, err := rb.ToApi(rts.ctx, roster)

	require.NoError(t, err, "ToApi should succeed")
	require.Equal(t, "profile123", result.ProfileId, "Profile ID should match")

}

func (rts *RosterTestSuite) TestRosterBusiness_GetByID() {

	t := rts.T()
	rb := business.NewRosterBusiness(rts.ctx, rts.service)

	rosterMap, err := rts.createRoster(rts.ctx, rb, "profile123", map[string]map[string]string{
		"roster@test.com": {"key1": "value1"},
	})
	assert.NoError(t, err)

	expectedRoster := rosterMap["roster@test.com"]
	rosterID := expectedRoster.GetId()

	result, err := rb.GetByID(rts.ctx, rosterID)

	assert.NoError(t, err)
	assert.Equal(t, expectedRoster.GetProfileId(), result.ProfileID, "Profile ID should match")
	assert.Equal(t, expectedRoster.GetContact().GetId(), result.ContactID, "Contact ID should match")
}

func (rts *RosterTestSuite) TestRosterBusiness_RemoveRoster() {
	t := rts.T()
	rb := business.NewRosterBusiness(rts.ctx, rts.service)

	rosterMap, err := rts.createRoster(rts.ctx, rb, "profRemov123", map[string]map[string]string{
		"rosterremove@test.com": {"key1": "value1"},
	})
	assert.NoError(t, err)

	roster := rosterMap["rosterremove@test.com"]
	rosterID := roster.GetId()

	result, err := rb.RemoveRoster(rts.ctx, rosterID)

	assert.NoError(t, err)
	assert.Equal(t, roster.GetProfileId(), result.GetProfileId())

	_, err = rb.GetByID(rts.ctx, rosterID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "Error should be 'gorm.ErrRecordNotFound'")
}

func (rts *RosterTestSuite) TestRosterBusiness_RemoveRoster_NotFound() {
	t := rts.T()
	rb := business.NewRosterBusiness(rts.ctx, rts.service)

	rosterID := "nonexistent"

	result, err := rb.RemoveRoster(rts.ctx, rosterID)

	require.Error(t, err, "A not found error should be returned")
	require.Nil(t, result, "Result should be nil")
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound), "Error should be 'gorm.ErrRecordNotFound'")
}

func (rts *RosterTestSuite) TestRosterBusiness_Search() {

	t := rts.T()
	rb := business.NewRosterBusiness(rts.ctx, rts.service)

	profileID := "searchProfileC1"

	_, err := rts.createRoster(rts.ctx, rb, profileID, map[string]map[string]string{
		"rostersearch@test.com":        {"name": "John Osogo", "age": "33"},
		"searchrostercontact@test.com": {"name": "Thomas Balindhe", "age": "36"},
		"+256755718293":                {"name": "Mary Osogo", "age": "21"},
		"+256755718291":                {"name": "Julius Search Best", "age": "51"},
	})
	require.NoError(t, err)

	tests := []struct {
		name        string
		request     *profilev1.SearchRosterRequest
		wantError   require.ErrorAssertionFunc
		profileId   string
		resultCount int
	}{
		{
			name: "Valid search request with matching data",
			request: &profilev1.SearchRosterRequest{
				Query:      "searchrostercontact@test.com",
				Properties: []string{"name", "age"},
			},
			profileId:   profileID,
			wantError:   require.NoError,
			resultCount: 1,
		},
		{
			name: "Valid search request partial contact match",
			request: &profilev1.SearchRosterRequest{
				Query:      "+25675571829",
				Properties: []string{"name", "age"},
			},
			profileId:   profileID,
			wantError:   require.NoError,
			resultCount: 2,
		},
		{
			name: "Valid search request partial properties match",
			request: &profilev1.SearchRosterRequest{
				Query:      "Osogo",
				Properties: []string{"name", "age"},
			},
			profileId:   profileID,
			wantError:   require.NoError,
			resultCount: 2,
		},
		{
			name: "Valid cross field partial match",
			request: &profilev1.SearchRosterRequest{
				Query:      "search",
				Properties: []string{"name", "age"},
			},
			profileId:   profileID,
			wantError:   require.NoError,
			resultCount: 3,
		},
		{
			name: "Request with non-existent profile ID",
			request: &profilev1.SearchRosterRequest{
				Query:      "non existent profile",
				Properties: []string{"name", "age"},
			},
			profileId:   "nonExistentProfileId",
			wantError:   require.NoError,
			resultCount: 0,
		},
		{
			name: "Empty search request",
			request: &profilev1.SearchRosterRequest{
				Query:      "",
				Properties: []string{"name", "age"},
			},
			profileId:   profileID,
			wantError:   require.NoError,
			resultCount: 4,
		},
		{
			name: "Empty search request with wrong profile",
			request: &profilev1.SearchRosterRequest{
				Query:      "",
				Properties: []string{"name", "age"},
			},
			profileId:   "funnyProfileId",
			wantError:   require.NoError,
			resultCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			claims := frame.ClaimsFromMap(map[string]string{
				"sub":          tt.profileId,
				"tenant_id":    "tenantx",
				"partition_id": "party",
			})

			ctx := claims.ClaimsToContext(rts.ctx)

			jobResult, err0 := rb.Search(ctx, tt.request)
			require.NoError(t, err0)

			var rosterList []*models.Roster
			for {
				result, ok := jobResult.ReadResult(rts.ctx)
				if result == nil || result.IsError() || !ok {
					break
				}
				rosterList = append(rosterList, result.Item()...)

			}

			require.Equal(t, tt.resultCount, len(rosterList), "Roster count mismatch")
		})
	}
}
