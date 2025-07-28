package events_test

import (
	"testing"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/tests/testdef"
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
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful save",

			args: args{
				payload: &models.SettingAudit{
					BaseModel: frame.BaseModel{
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

	ats.WithTestDependancies(t, func(t *testing.T, depOpt *testdef.DependancyOption) {
		svc, ctx := ats.CreateService(t, depOpt)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				e := &events.SettingsAuditor{
					Service: svc,
				}
				if err := e.Execute(ctx, tt.args.payload); (err != nil) != tt.wantErr {
					t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				}

				nRepo := repository.NewSettingAuditRepository(ctx, svc)
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
