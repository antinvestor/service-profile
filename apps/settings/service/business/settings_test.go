package business_test

import (
	"context"
	"fmt"
	"testing"

	settingsv1 "buf.build/gen/go/antinvestor/settingz/protocolbuffers/go/settings/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	"github.com/antinvestor/service-profile/apps/settings/service/business"
	"github.com/antinvestor/service-profile/apps/settings/service/models"
	"github.com/antinvestor/service-profile/apps/settings/service/repository"
	"github.com/antinvestor/service-profile/apps/settings/tests"
)

type SettingsTestSuite struct {
	tests.SettingsBaseTestSuite
}

func TestSettings(t *testing.T) {
	suite.Run(t, new(SettingsTestSuite))
}

func (ts *SettingsTestSuite) getSettingBusiness(
	ctx context.Context,
	svc *frame.Service,
) (business.SettingsBusiness, repository.ReferenceRepository, repository.SettingValRepository) {
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	refRepo := repository.NewReferenceRepository(ctx, dbPool, workMan)
	valRepo := repository.NewSettingValRepository(ctx, dbPool, workMan)

	return business.NewSettingsBusiness(refRepo, valRepo), refRepo, valRepo
}

// TestNewSettingsBusiness tests the creation of a new settings business.
func (ts *SettingsTestSuite) TestNewSettingsBusiness() {
	testcases := []struct {
		name       string
		want       business.SettingsBusiness
		nilService bool
		wantErr    require.ErrorAssertionFunc
	}{

		{
			name:    "NewSettingsBusiness",
			wantErr: require.NoError,
		},

		{
			name:       "NewSettingsBusinessWithNils",
			nilService: true,
			wantErr:    require.Error},
	}

	ts.WithTestDependancies(ts.T(), func(t *testing.T, depOpt *definition.DependencyOption) {
		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				var svc *frame.Service
				var ctx context.Context

				if !tt.nilService {
					svc, ctx = ts.CreateService(t, depOpt)
				}

				got, _, _ := ts.getSettingBusiness(ctx, svc)
				require.NotNil(t, got)
			})
		}
	})
}

// Test_settingsBusiness_Set tests the Set method of the settings business.
func (ts *SettingsTestSuite) Test_settingsBusiness_Set() {
	ts.WithTestDependancies(ts.T(), func(t *testing.T, depOpt *definition.DependencyOption) {
		svc, ctx := ts.CreateService(t, depOpt)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		type args struct {
			req *settingsv1.SetRequest
		}
		testcases := []struct {
			name     string
			args     args
			want     *settingsv1.SetResponse
			response string
		}{
			{
				name: "Set successfully",
				args: args{
					req: &settingsv1.SetRequest{
						Key: &settingsv1.Setting{
							Name:     "set.test.success",
							Object:   "",
							ObjectId: "",
							Lang:     "",
							Module:   "Listings",
						}, Value: "happy path"},
				},
				response: "happy path",
			},

			{
				name: "Set fuzzy",
				args: args{
					req: &settingsv1.SetRequest{
						Key: &settingsv1.Setting{
							Name:     "set",
							Object:   "",
							ObjectId: "",
							Lang:     "",
							Module:   "",
						}, Value: "fuzzy"},
				},
				response: "fuzzy",
			},
			{
				name: "Set Fully",
				args: args{
					req: &settingsv1.SetRequest{
						Key: &settingsv1.Setting{
							Name:     "set.test.2.success",
							Object:   "testi",
							ObjectId: "idto",
							Lang:     "en",
							Module:   "Mone",
						},
						Value: "val",
					},
				},
				response: "val",
			},
		}
		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				nb, _, _ := ts.getSettingBusiness(ctx, svc)

				got, err := nb.Set(ctx, tt.args.req)
				if err != nil {
					t.Errorf("Set() error = %v", err)
					return
				}
				if got.GetData() != nil && got.GetData().GetValue() != tt.response {
					t.Errorf("Set() response got = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

// createTestSettings creates test setting references and values in the database.
func (ts *SettingsTestSuite) createTestSettings(
	ctx context.Context,
	t *testing.T,
	rRepo repository.ReferenceRepository,
	vRepo repository.SettingValRepository,
) {
	for _, ref := range []models.SettingRef{
		{
			Name:     "get.test.success",
			Object:   "",
			ObjectID: "",
			Language: "",
			Module:   "Listings",
		},
		{
			Name:     "get.test.with.id",
			Object:   "tester",
			ObjectID: "tid",
			Language: "",
			Module:   "Listings",
		}, {
			Name:     "get.test.with.lang",
			Object:   "tester",
			ObjectID: "tid",
			Language: "en",
			Module:   "Listings",
		},
	} {
		err := rRepo.Create(ctx, &ref)
		if err != nil {
			t.Errorf("Could not save setting ref for listing, %v", err)
			return
		}

		err = vRepo.Create(ctx, &models.SettingVal{
			Ref:     ref.GetID(),
			Detail:  fmt.Sprintf("Random value for : %s", ref.Name),
			Version: 0,
		})
		if err != nil {
			t.Errorf("Could not save setting val for listing, %v", err)
			return
		}
	}
}

// runGetTest runs a single Get test case.
func (ts *SettingsTestSuite) runGetTest(ctx context.Context, t *testing.T, svc *frame.Service, tt struct {
	name     string
	req      *settingsv1.GetRequest
	response string
}) {
	nb, _, _ := ts.getSettingBusiness(ctx, svc)

	got, err := nb.Get(ctx, tt.req)
	if err != nil {
		t.Errorf("Get() error = %v", err)
		return
	}

	if got.GetData() != nil && !proto.Equal(got.GetData().GetKey(), tt.req.GetKey()) {
		t.Errorf("Get() got invalid key = %v, want %v", got.GetData().GetKey(), tt.req.GetKey())
		return
	}

	if got.GetData() != nil && got.GetData().GetValue() != tt.response {
		t.Errorf("Get() got invalid value = %v, want %v", got.GetData().GetValue(), tt.response)
		return
	}
}

// Test_settingsBusiness_Get tests the Get method of the settings business.
func (ts *SettingsTestSuite) Test_settingsBusiness_Get() {
	testcases := []struct {
		name     string
		req      *settingsv1.GetRequest
		response string
	}{
		{
			name: "Get successfully",
			req: &settingsv1.GetRequest{Key: &settingsv1.Setting{
				Name:     "get.test.success",
				Object:   "",
				ObjectId: "",
				Lang:     "",
				Module:   "Listings",
			}},
			response: "Random value for : get.test.success",
		},
		{
			name: "Get fuzzy",
			req: &settingsv1.GetRequest{Key: &settingsv1.Setting{
				Name:     "get",
				Object:   "",
				ObjectId: "",
				Lang:     "",
				Module:   "",
			}},
			response: "",
		},
		{
			name: "Get Less Module",
			req: &settingsv1.GetRequest{Key: &settingsv1.Setting{
				Name:     "get.test.success",
				Object:   "",
				ObjectId: "",
				Lang:     "",
				Module:   "",
			}},
			response: "",
		},
		{
			name: "Get Random key",
			req: &settingsv1.GetRequest{Key: &settingsv1.Setting{
				Name:     "get.missing.key",
				Object:   "",
				ObjectId: "",
				Lang:     "",
				Module:   "",
			}},
			response: "",
		},
		{
			name: "Get with object id",
			req: &settingsv1.GetRequest{Key: &settingsv1.Setting{
				Name:     "get.test.with.id",
				Object:   "tester",
				ObjectId: "tid",
				Lang:     "",
				Module:   "Listings",
			}},
			response: "Random value for : get.test.with.id",
		},
		{
			name: "Get with language",
			req: &settingsv1.GetRequest{Key: &settingsv1.Setting{
				Name:     "get.test.with.lang",
				Object:   "tester",
				ObjectId: "tid",
				Lang:     "en",
				Module:   "Listings",
			}},
			response: "Random value for : get.test.with.lang",
		},
	}

	ts.WithTestDependancies(ts.T(), func(t *testing.T, depOpt *definition.DependencyOption) {
		svc, ctx := ts.CreateService(t, depOpt)

		_, rRepo, vRepo := ts.getSettingBusiness(ctx, svc)

		// Setup test data
		ts.createTestSettings(ctx, t, rRepo, vRepo)

		// Run test cases
		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				ts.runGetTest(ctx, t, svc, tt)
			})
		}
	})
}

// createListTestSettings creates test setting references and values for list testing.
func (ts *SettingsTestSuite) createListTestSettings(
	ctx context.Context,
	t *testing.T,
	rRepo repository.ReferenceRepository,
	vRepo repository.SettingValRepository,
) {
	for _, ref := range []models.SettingRef{
		{
			Name:     "testing.listing.ok",
			Object:   "",
			ObjectID: "",
			Language: "",
			Module:   "Listings",
		}, {
			Name:     "listing.test.success",
			Object:   "",
			ObjectID: "",
			Language: "",
			Module:   "Listings",
		},
		{
			Name:     "listing.test.with.id",
			Object:   "tester",
			ObjectID: "tid",
			Language: "",
			Module:   "Listings",
		},
	} {
		err := rRepo.Create(ctx, &ref)
		if err != nil {
			t.Errorf("Could not save setting ref for listing, %v", err)
			return
		}

		err = vRepo.Create(ctx, &models.SettingVal{
			Ref:     ref.GetID(),
			Detail:  fmt.Sprintf("Random value for : %s", ref.Name),
			Version: 0,
		})
		if err != nil {
			t.Errorf("Could not save setting val for listing, %v", err)
			return
		}
	}
}

// runListTest runs a single List test case.
func (ts *SettingsTestSuite) runListTest(ctx context.Context, t *testing.T, svc *frame.Service, tt struct {
	name          string
	req           *settingsv1.ListRequest
	wantSendCalls int
	wantDataLen   int
}) {
	nb, _, _ := ts.getSettingBusiness(ctx, svc)

	responses, err := nb.List(ctx, tt.req)
	if err != nil {
		t.Errorf("List() error = %v", err)
		return
	}

	if tt.wantDataLen > 0 {
		if len(responses) != tt.wantDataLen {
			t.Errorf("Data length is not as expected, Got %d, Want %d", len(responses), tt.wantDataLen)
		}
	}
}

func (ts *SettingsTestSuite) Test_settingsBusiness_List() {
	testcases := []struct {
		name          string
		req           *settingsv1.ListRequest
		wantSendCalls int
		wantDataLen   int
	}{
		{
			name: "Query Successfully",
			req: &settingsv1.ListRequest{Key: &settingsv1.Setting{
				Name:     "listing.test.success",
				Object:   "",
				ObjectId: "",
				Lang:     "",
				Module:   "Listings",
			}},
			wantSendCalls: 1,
			wantDataLen:   1,
		},
		{
			name: "Query Less Module",
			req: &settingsv1.ListRequest{Key: &settingsv1.Setting{
				Name:     "listing.test.success",
				Object:   "",
				ObjectId: "",
				Lang:     "",
				Module:   "",
			}},
			wantSendCalls: 1,
			wantDataLen:   1,
		},

		{
			name: "Query Fuzzy",
			req: &settingsv1.ListRequest{Key: &settingsv1.Setting{
				Name:     "listing",
				Object:   "",
				ObjectId: "",
				Lang:     "",
				Module:   "",
			}},
			wantSendCalls: 1,
			wantDataLen:   3,
		},
		{
			name: "Query Empty",
			req: &settingsv1.ListRequest{Key: &settingsv1.Setting{
				Name:     "listing.test.empty",
				Object:   "",
				ObjectId: "",
				Lang:     "",
				Module:   "",
			}},
			wantSendCalls: 1,
			wantDataLen:   0,
		},
	}

	ts.WithTestDependancies(ts.T(), func(t *testing.T, depOpt *definition.DependencyOption) {
		svc, ctx := ts.CreateService(t, depOpt)

		_, rRepo, vRepo := ts.getSettingBusiness(ctx, svc)

		// Setup test data
		ts.createListTestSettings(ctx, t, rRepo, vRepo)

		// Run test cases
		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				ts.runListTest(ctx, t, svc, tt)
			})
		}
	})
}
