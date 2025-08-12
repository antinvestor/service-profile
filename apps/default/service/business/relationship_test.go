package business_test

import (
	"context"
	"testing"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/tests"
)

type RelationshipTestSuite struct {
	tests.BaseTestSuite
}

func TestRelationshipSuite(t *testing.T) {
	suite.Run(t, new(RelationshipTestSuite))
}

func (rts *RelationshipTestSuite) TestNewRelationshipBusiness() {
	t := rts.T()

	testcases := []struct {
		name string
	}{
		{
			name: "New relationship business test",
		},
	}

	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				svc, ctx := rts.CreateService(t, dep)

				profileBusiness := business.NewProfileBusiness(ctx, svc)
				if got := business.NewRelationshipBusiness(ctx, svc, profileBusiness); got == nil {
					t.Errorf("NewRelationshipBusiness() = %v, is nil", got)
				}
			})
		}
	})
}

func (rts *RelationshipTestSuite) Test_relationshipBusiness_CreateRelationship() {
	t := rts.T()
	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := rts.CreateService(t, dep)
		profileBusiness := business.NewProfileBusiness(ctx, svc)

		testProfiles, err := rts.CreateTestProfiles(
			ctx,
			svc,
			[]string{"new.relationship.1@ant.com", "new.relationship.2@ant.com"},
		)
		if err != nil {
			t.Errorf(" CreateProfile failed with %+v", err)
			return
		}

		type args struct {
			request *profilev1.AddRelationshipRequest
		}
		tests := []struct {
			name    string
			args    args
			want    *profilev1.RelationshipObject
			wantErr bool
		}{
			{
				name: "Create a relationship object",
				args: args{
					request: &profilev1.AddRelationshipRequest{
						Parent:     "Profile",
						ParentId:   testProfiles[0].GetId(),
						Child:      "Profile",
						ChildId:    testProfiles[1].GetId(),
						Type:       profilev1.RelationshipType_MEMBER,
						Properties: nil,
					},
				},
				want: &profilev1.RelationshipObject{
					Id:         "",
					Type:       0,
					Properties: nil,
					ChildEntry: &profilev1.EntryItem{ObjectName: "Profile", ObjectId: testProfiles[0].GetId()},
				},
				wantErr: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				aB := business.NewRelationshipBusiness(ctx, svc, profileBusiness)
				got, err1 := aB.CreateRelationship(ctx, tt.args.request)
				if (err1 != nil) != tt.wantErr {
					t.Errorf("CreateRelationship() error = %v, wantErr %v", err1, tt.wantErr)
					return
				}

				if tt.wantErr {
					return
				}

				gotProfile := got.GetChildEntry()

				wantProfile := got.GetChildEntry()

				if gotProfile.GetObjectId() != wantProfile.GetObjectId() {
					t.Errorf("CreateRelationship() got = %v, want %v", gotProfile, wantProfile)
				}
			})
		}
	})
}

func (rts *RelationshipTestSuite) Test_relationshipBusiness_DeleteRelationship() {
	t := rts.T()
	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := rts.CreateService(t, dep)

		profileBusiness := business.NewProfileBusiness(ctx, svc)
		aB := business.NewRelationshipBusiness(ctx, svc, profileBusiness)

		testProfiles, err := rts.CreateTestProfiles(
			ctx,
			svc,
			[]string{"delete.relationship.1@ant.com", "delete.relationship.2@ant.com"},
		)
		if err != nil {
			t.Errorf(" Delete profile failed with %+v", err)
			return
		}

		existingRelation, err := aB.CreateRelationship(ctx, &profilev1.AddRelationshipRequest{
			Parent:     "Profile",
			ParentId:   testProfiles[0].GetId(),
			Child:      "Profile",
			ChildId:    testProfiles[1].GetId(),
			Type:       profilev1.RelationshipType_MEMBER,
			Properties: nil,
		})
		if err != nil {
			t.Errorf("DeleteRelationship() error = %v", err)
			return
		}

		type args struct {
			request *profilev1.DeleteRelationshipRequest
		}
		testcases := []struct {
			name    string
			args    args
			want    *profilev1.RelationshipObject
			wantErr bool
		}{
			{
				name: "Delete existing relation",
				args: args{
					request: &profilev1.DeleteRelationshipRequest{
						Id: existingRelation.GetId(),
					},
				},
				want:    nil,
				wantErr: false,
			},
			{
				name: "Delete deleted relation",
				args: args{
					request: &profilev1.DeleteRelationshipRequest{
						Id: existingRelation.GetId(),
					},
				},
				want:    nil,
				wantErr: true,
			},
		}
		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				_, deleteErr := aB.DeleteRelationship(ctx, tt.args.request)
				if (deleteErr != nil) != tt.wantErr {
					t.Errorf("DeleteRelationship() error = %v, wantErr %v", deleteErr, tt.wantErr)
					return
				}
			})
		}
	})
}

func (rts *RelationshipTestSuite) Test_relationshipBusiness_ListRelationships() {
	t := rts.T()
	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := rts.CreateService(t, dep)
		profileBusiness := business.NewProfileBusiness(ctx, svc)

		relationshipBusiness := business.NewRelationshipBusiness(ctx, svc, profileBusiness)

		testProfiles, err := rts.CreateTestProfiles(
			ctx,
			svc,
			[]string{
				"list.relationship.1@ant.com",
				"list.relationship.2@ant.com",
				"list.relationship.3@ant.com",
				"list.relationship.4@ant.com",
			},
		)
		if err != nil {
			t.Errorf(" List profile failed with %+v", err)
			return
		}

		for i := range 3 {
			_, err = relationshipBusiness.CreateRelationship(ctx, &profilev1.AddRelationshipRequest{
				Parent:     "Profile",
				ParentId:   testProfiles[0].GetId(),
				Child:      "Profile",
				ChildId:    testProfiles[i+1].GetId(),
				Type:       profilev1.RelationshipType_MEMBER,
				Properties: nil,
			})

			if err != nil {
				t.Errorf(" List relationship failed with %+v", err)
				return
			}
		}

		type args struct {
			ctx     context.Context
			request *profilev1.ListRelationshipRequest
		}
		testcases := []struct {
			name      string
			args      args
			wantCount int
			wantErr   require.ErrorAssertionFunc
		}{
			{
				name: "None existent relationships",
				args: args{
					ctx: ctx,
					request: &profilev1.ListRelationshipRequest{
						PeerName: "Profile",
						PeerId:   "bjt4h376abi8cg3kgr80",
					},
				},
				wantCount: 0,
				wantErr:   require.Error,
			},
			{
				name: "Existent relationships",
				args: args{
					ctx: ctx,
					request: &profilev1.ListRelationshipRequest{
						PeerName: "Profile",
						PeerId:   testProfiles[0].GetId(),
					},
				},
				wantCount: 3,
				wantErr:   require.NoError,
			},
			{
				name: "Limited existent relationships",
				args: args{
					ctx: ctx,
					request: &profilev1.ListRelationshipRequest{
						PeerName: "Profile",
						PeerId:   testProfiles[0].GetId(),
						Count:    2,
					},
				},
				wantCount: 2,
				wantErr:   require.NoError,
			},
			{
				name: "Specific existent relationships",
				args: args{
					ctx: ctx,
					request: &profilev1.ListRelationshipRequest{
						PeerName:          "Profile",
						PeerId:            testProfiles[0].GetId(),
						RelatedChildrenId: []string{testProfiles[3].GetId()},
					},
				},
				wantCount: 1,
				wantErr:   require.NoError,
			},
		}
		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				got, err0 := relationshipBusiness.ListRelationships(ctx, tt.args.request)
				tt.wantErr(t, err0)
				if len(got) != tt.wantCount {
					t.Errorf("ListRelationships() got = %v, want %v", len(got), tt.wantCount)
				}
			})
		}
	})
}
