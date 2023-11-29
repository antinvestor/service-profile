package business

import (
	"context"
	profilev1 "github.com/antinvestor/apis/profile"
	"github.com/antinvestor/service-profile/service"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/pitabwire/frame"
)

type RelationshipBusiness interface {
	ListRelationships(ctx context.Context, request *profilev1.ProfileListRelationshipRequest) ([]*models.Relationship, error)
	CreateRelationship(ctx context.Context, request *profilev1.ProfileAddRelationshipRequest) (*profilev1.RelationshipObject, error)
	DeleteRelationship(ctx context.Context, request *profilev1.ProfileDeleteRelationshipRequest) (*profilev1.RelationshipObject, error)

	ToAPI(ctx context.Context, sourceParent, sourceParentID string, relationship *models.Relationship) (*profilev1.RelationshipObject, error)
}

func NewRelationshipBusiness(ctx context.Context, service *frame.Service, profileEncryptionKey []byte) RelationshipBusiness {
	relationshipRepo := repository.NewRelationshipRepository(service)
	profileBiz := NewProfileBusiness(ctx, service)

	return &relationshipBusiness{
		service:              service,
		profileBusiness:      profileBiz,
		profileEncryptionKey: profileEncryptionKey,
		relationshipRepo:     relationshipRepo,
	}
}

type relationshipBusiness struct {
	service              *frame.Service
	profileBusiness      ProfileBusiness
	profileEncryptionKey []byte
	relationshipRepo     repository.RelationshipRepository
}

func (aB *relationshipBusiness) ToAPI(ctx context.Context, sourceParent, sourceParentID string, relationship *models.Relationship) (*profilev1.RelationshipObject, error) {

	if relationship == nil {
		return nil, nil
	}

	parentId := relationship.ParentObjectID
	if sourceParent != relationship.ParentObject && sourceParentID != relationship.ParentObjectID {

		if relationship.ChildObject != "Profile" {
			//TODO: only support relationships between profiles alone
			return nil, nil
		}

		parentId = relationship.ChildObjectID

	}

	relationshipObj := &profilev1.RelationshipObject{
		ID:         relationship.GetID(),
		Type:       profilev1.RelationshipType(relationship.RelationshipType.UID),
		Properties: frame.DBPropertiesToMap(relationship.Properties),
	}

	if relationship.ChildObject == "Profile" {
		profileObj, err := aB.profileBusiness.GetByID(ctx, aB.profileEncryptionKey, parentId)
		if err != nil {
			return nil, err
		}

		relationshipObj.Child = &profilev1.RelationshipObject_Profile{Profile: profileObj}
	}

	return relationshipObj, nil

}

func (aB *relationshipBusiness) ListRelationships(ctx context.Context, request *profilev1.ProfileListRelationshipRequest) ([]*models.Relationship, error) {

	if request.GetParent() == "Profile" {
		profileObj, err := aB.profileBusiness.GetByID(ctx, aB.profileEncryptionKey, request.GetParentID())
		if err != nil {
			return nil, err
		}

		if profileObj == nil {
			return nil, service.ErrorProfileDoesNotExist
		}
	}

	return aB.relationshipRepo.List(ctx, request.GetParent(), request.GetParentID(), request.GetRelatedChildrenID(), request.GetLastRelationshipID(), int(request.GetCount()))
}

func (aB *relationshipBusiness) CreateRelationship(ctx context.Context, request *profilev1.ProfileAddRelationshipRequest) (*profilev1.RelationshipObject, error) {

	var profileObj *profilev1.ProfileObject
	logger := aB.service.L().WithField("request", request)

	relationships, err := aB.relationshipRepo.List(ctx, request.GetParent(), request.GetParentID(), []string{request.GetChildID()}, "", 2)
	if err != nil {
		logger.WithError(err).Warn("get existing relationship error")

		if !frame.DBErrorIsRecordNotFound(err) {
			return nil, err
		}
	}

	if len(relationships) > 0 {

		return aB.ToAPI(ctx, request.GetParent(), request.GetParentID(), relationships[0])
	}

	if request.GetParent() == "Profile" {
		profileObj, err = aB.profileBusiness.GetByID(ctx, aB.profileEncryptionKey, request.GetParentID())
		if err != nil {
			return nil, err
		}

		if profileObj == nil {
			return nil, service.ErrorProfileDoesNotExist
		}
	}

	if request.GetChild() == "Profile" {
		profileObj, err = aB.profileBusiness.GetByID(ctx, aB.profileEncryptionKey, request.GetChildID())
		if err != nil {
			return nil, err
		}

		if profileObj == nil {
			return nil, service.ErrorProfileDoesNotExist
		}
	}

	relationshipType, err := aB.relationshipRepo.RelationshipType(ctx, request.GetType())
	if err != nil {
		return nil, err
	}

	a := models.Relationship{
		BaseModel: frame.BaseModel{
			ID: request.GetID(),
		},
		ParentObject:       request.GetParent(),
		ParentObjectID:     request.GetParentID(),
		RelationshipTypeID: relationshipType.GetID(),
		RelationshipType:   relationshipType,
		ChildObject:        request.GetChild(),
		ChildObjectID:      request.GetChildID(),
		Properties:         nil,
	}

	err = aB.relationshipRepo.Save(ctx, &a)
	if err != nil {
		return nil, err
	}

	return aB.ToAPI(ctx, request.GetParent(), request.GetParentID(), &a)

}
func (aB *relationshipBusiness) DeleteRelationship(ctx context.Context, request *profilev1.ProfileDeleteRelationshipRequest) (*profilev1.RelationshipObject, error) {

	relationship, err := aB.relationshipRepo.GetByID(ctx, request.GetID())
	if err != nil || relationship == nil {
		return nil, err
	}

	if request.GetParentID() == "" || request.GetParentID() == relationship.ParentObjectID {

		err = aB.relationshipRepo.Delete(ctx, request.GetID())
		if err != nil {
			return nil, err
		}

		return aB.ToAPI(ctx, relationship.ParentObject, relationship.ParentObjectID, relationship)
	}

	return nil, nil
}
