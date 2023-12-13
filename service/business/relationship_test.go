package business_test

import (
	"context"
	profilev1 "github.com/antinvestor/apis/profile/v1"
	"github.com/antinvestor/service-profile/service/business"
	"github.com/pitabwire/frame"
	"testing"
)

func createTestProfiles(ctx context.Context, srv *frame.Service, encryptionKey []byte, contacts []string) ([]*profilev1.ProfileObject, error) {

	profBuss := business.NewProfileBusiness(ctx, srv, func() []byte {
		return encryptionKey
	})

	var profileSlice []*profilev1.ProfileObject

	for _, contact := range contacts {

		prof := &profilev1.CreateRequest{
			Contact: contact,
		}
		profile, err := profBuss.CreateProfile(ctx, prof)
		if err != nil {
			return nil, err
		}

		profileSlice = append(profileSlice, profile)
	}

	return profileSlice, nil
}

func TestNewRelationshipBusiness(t *testing.T) {

	ctx, srv := getTestService()

	type args struct {
		ctx                  context.Context
		service              *frame.Service
		profileEncryptionKey []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "New relationship business test",
			args: args{
				ctx:                  ctx,
				service:              srv,
				profileEncryptionKey: getEncryptionKey(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profileBusiness := business.NewProfileBusiness(ctx, srv, func() []byte {
				return tt.args.profileEncryptionKey
			})
			if got := business.NewRelationshipBusiness(tt.args.ctx, tt.args.service, profileBusiness); got == nil {
				t.Errorf("NewRelationshipBusiness() = %v, is nil", got)
			}
		})
	}
}

func Test_relationshipBusiness_CreateRelationship(t *testing.T) {

	ctx, srv := getTestService()
	profileEncryptionKey := getEncryptionKey()

	profileBusiness := business.NewProfileBusiness(ctx, srv, func() []byte {
		return profileEncryptionKey
	})

	testProfiles, err := createTestProfiles(ctx, srv, profileEncryptionKey, []string{"new.relationship.1@ant.com", "new.relationship.2@ant.com"})
	if err != nil {
		t.Errorf(" CreateProfile failed with %+v", err)
		return
	}

	type args struct {
		ctx     context.Context
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
				ctx: ctx,
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
				Child:      &profilev1.RelationshipObject_Profile{Profile: testProfiles[1]},
			},
			wantErr: false,
		},
		{
			name: "Create a fake relationship object",
			args: args{
				ctx: ctx,
				request: &profilev1.AddRelationshipRequest{
					Parent:     "Profile",
					ParentId:   testProfiles[0].GetId(),
					Child:      "Profile",
					ChildId:    "bjt4h376abi8cg3kgr80",
					Type:       profilev1.RelationshipType_MEMBER,
					Properties: nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid data relationship object",
			args: args{
				ctx: ctx,
				request: &profilev1.AddRelationshipRequest{
					Parent:     "Jokes",
					ParentId:   testProfiles[0].GetId(),
					Child:      "Profile",
					ChildId:    "",
					Type:       profilev1.RelationshipType_MEMBER,
					Properties: nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aB := business.NewRelationshipBusiness(ctx, srv, profileBusiness)
			got, err1 := aB.CreateRelationship(tt.args.ctx, tt.args.request)
			if (err1 != nil) != tt.wantErr {
				t.Errorf("CreateRelationship() error = %v, wantErr %v", err1, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			gotProfile, ok := got.GetChild().(*profilev1.RelationshipObject_Profile)
			if !ok {
				t.Errorf("CreateRelationship() child is not a profile : %v ", gotProfile)
				return
			}

			wantProfile, ok := got.GetChild().(*profilev1.RelationshipObject_Profile)
			if !ok {
				t.Errorf("CreateRelationship() child is not a profile : %v", wantProfile)
				return
			}
			if gotProfile.Profile.GetId() != wantProfile.Profile.GetId() {
				t.Errorf("CreateRelationship() got = %v, want %v", gotProfile, wantProfile)
			}
		})
	}
}

func Test_relationshipBusiness_DeleteRelationship(t *testing.T) {

	ctx, srv := getTestService()
	profileEncryptionKey := getEncryptionKey()

	profileBusiness := business.NewProfileBusiness(ctx, srv, func() []byte {
		return profileEncryptionKey
	})
	aB := business.NewRelationshipBusiness(ctx, srv, profileBusiness)

	testProfiles, err := createTestProfiles(ctx, srv, profileEncryptionKey, []string{"delete.relationship.1@ant.com", "delete.relationship.2@ant.com"})
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
		t.Errorf("CreateRelationship() error = %v", err)
		return
	}

	type args struct {
		ctx     context.Context
		request *profilev1.DeleteRelationshipRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *profilev1.RelationshipObject
		wantErr bool
	}{
		{
			name: "Delete existing relation",
			args: args{
				ctx: ctx,
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
				ctx: ctx,
				request: &profilev1.DeleteRelationshipRequest{
					Id: existingRelation.GetId(),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := aB.DeleteRelationship(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteRelationship() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func Test_relationshipBusiness_ListRelationships(t *testing.T) {

	ctx, srv := getTestService()
	profileEncryptionKey := getEncryptionKey()

	profileBusiness := business.NewProfileBusiness(ctx, srv, func() []byte {
		return profileEncryptionKey
	})

	aB := business.NewRelationshipBusiness(ctx, srv, profileBusiness)

	testProfiles, err := createTestProfiles(ctx, srv, profileEncryptionKey, []string{"list.relationship.1@ant.com", "list.relationship.2@ant.com", "list.relationship.3@ant.com", "list.relationship.4@ant.com"})
	if err != nil {
		t.Errorf(" Delete profile failed with %+v", err)
		return
	}

	for i := 0; i < 3; i++ {

		_, err = aB.CreateRelationship(ctx, &profilev1.AddRelationshipRequest{
			Parent:     "Profile",
			ParentId:   testProfiles[0].GetId(),
			Child:      "Profile",
			ChildId:    testProfiles[i+1].GetId(),
			Type:       profilev1.RelationshipType_MEMBER,
			Properties: nil,
		})

		if err != nil {
			t.Errorf(" Create relationship failed with %+v", err)
			return
		}
	}

	type args struct {
		ctx     context.Context
		request *profilev1.ListRelationshipRequest
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "None existent relationships",
			args: args{
				ctx: ctx,
				request: &profilev1.ListRelationshipRequest{
					Parent:   "Profile",
					ParentId: "bjt4h376abi8cg3kgr80",
				},
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "Existent relationships",
			args: args{
				ctx: ctx,
				request: &profilev1.ListRelationshipRequest{
					Parent:   "Profile",
					ParentId: testProfiles[0].GetId(),
				},
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name: "Limited existent relationships",
			args: args{
				ctx: ctx,
				request: &profilev1.ListRelationshipRequest{
					Parent:   "Profile",
					ParentId: testProfiles[0].GetId(),
					Count:    2,
				},
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "Specific existent relationships",
			args: args{
				ctx: ctx,
				request: &profilev1.ListRelationshipRequest{
					Parent:            "Profile",
					ParentId:          testProfiles[0].GetId(),
					RelatedChildrenId: []string{testProfiles[3].GetId()},
				},
			},
			wantCount: 1,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			aB := business.NewRelationshipBusiness(ctx, srv, profileBusiness)
			got, err := aB.ListRelationships(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListRelationships() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("ListRelationships() got = %v, want %v", len(got), tt.wantCount)
			}
		})
	}
}
