package events_test

import (
	"testing"

	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/settings/service/events"
	"github.com/antinvestor/service-profile/apps/settings/service/models"
	"github.com/antinvestor/service-profile/apps/settings/service/repository"
	"github.com/antinvestor/service-profile/apps/settings/tests"
)

type AuditorTestSuite struct {
	tests.SettingsBaseTestSuite
}

func TestAuditor(t *testing.T) {
	suite.Run(t, new(AuditorTestSuite))
}

func (ats *AuditorTestSuite) TestSettingAuditor_Execute() {
	t := ats.T()

	type args struct {
		payload interface{}
	}
	testCases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful save",

			args: args{
				payload: &models.SettingAudit{
					BaseModel: data.BaseModel{
						ID:          "testingSaveId",
						TenantID:    "tenantData",
						PartitionID: "partitionData",
						AccessID:    "testingAccessData",
					},
					Ref:     "linkedHistory",
					Detail:  "epochTesting",
					Version: 5,
				},
			},
			wantErr: false,
		},
	}

	ats.WithTestDependancies(t, func(t *testing.T, depOpt *definition.DependencyOption) {
		ctx, svc := ats.CreateService(t, depOpt)

		workMan := svc.WorkManager()
		dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

		nRepo := repository.NewSettingAuditRepository(ctx, dbPool, workMan)

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				e := events.NewSettingsAuditor(nRepo)
				if err := e.Execute(ctx, tt.args.payload); (err != nil) != tt.wantErr {
					t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				}

				audits, err := nRepo.GetByRef(ctx, "linkedHistory")
				if err != nil {
					t.Errorf("Search() error = %v could not retrieve saved audits", err)
					return
				}

				if audits == nil {
					t.Errorf("Matching audits could not be found")
					return
				}
			})
		}
	})
}
