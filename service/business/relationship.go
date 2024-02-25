package business

import (
	"context"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/pitabwire/frame"
)

type RelationshipBusiness interface {
	ListRelationships(ctx context.Context, request *profilev1.ListRelationshipRequest) ([]*models.Relationship, error)
	CreateRelationship(ctx context.Context, request *profilev1.AddRelationshipRequest) (*profilev1.RelationshipObject, error)
	DeleteRelationship(ctx context.Context, request *profilev1.DeleteRelationshipRequest) (*profilev1.RelationshipObject, error)

	ToAPI(ctx context.Context, sourceParent, sourceParentID string, relationship *models.Relationship) (*profilev1.RelationshipObject, error)
}

func NewRelationshipBusiness(_ context.Context, service *frame.Service, profileBiz ProfileBusiness) RelationshipBusiness {
	relationshipRepo := repository.NewRelationshipRepository(service)

	return &relationshipBusiness{
		service:          service,
		profileBusiness:  profileBiz,
		relationshipRepo: relationshipRepo,
	}
}

type relationshipBusiness struct {
	service          *frame.Service
	profileBusiness  ProfileBusiness
	relationshipRepo repository.RelationshipRepository
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
		Id:         relationship.GetID(),
		Type:       profilev1.RelationshipType(relationship.RelationshipType.UID),
		Properties: frame.DBPropertiesToMap(relationship.Properties),
	}

	if relationship.ChildObject == "Profile" {
		profileObj, err := aB.profileBusiness.GetByID(ctx, parentId)
		if err != nil {
			return nil, err
		}

		relationshipObj.Child = &profilev1.RelationshipObject_Profile{Profile: profileObj}
	}

	return relationshipObj, nil

}

func (aB *relationshipBusiness) ListRelationships(ctx context.Context, request *profilev1.ListRelationshipRequest) ([]*models.Relationship, error) {

	if request.GetParent() == "Profile" {
		profileObj, err := aB.profileBusiness.GetByID(ctx, request.GetParentId())
		if err != nil {
			return nil, err
		}

		if profileObj == nil {
			return nil, service.ErrorProfileDoesNotExist
		}
	}

	return aB.relationshipRepo.List(ctx, request.GetParent(), request.GetParentId(), request.GetRelatedChildrenId(), request.GetLastRelationshipId(), int(request.GetCount()))
}

func (aB *relationshipBusiness) CreateRelationship(ctx context.Context, request *profilev1.AddRelationshipRequest) (*profilev1.RelationshipObject, error) {

	var profileObj *profilev1.ProfileObject
	logger := aB.service.L().WithField("request", request)

	relationships, err := aB.relationshipRepo.List(ctx, request.GetParent(), request.GetParentId(), []string{request.GetChildId()}, "", 2)
	if err != nil {
		logger.WithError(err).Warn("get existing relationship error")

		if !frame.DBErrorIsRecordNotFound(err) {
			return nil, err
		}
	}

	if len(relationships) > 0 {

		return aB.ToAPI(ctx, request.GetParent(), request.GetParentId(), relationships[0])
	}

	if request.GetParent() == "Profile" {
		profileObj, err = aB.profileBusiness.GetByID(ctx, request.GetParentId())
		if err != nil {
			return nil, err
		}

		if profileObj == nil {
			return nil, service.ErrorProfileDoesNotExist
		}
	}

	if request.GetChild() == "Profile" {
		profileObj, err = aB.profileBusiness.GetByID(ctx, request.GetChildId())
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
		ParentObject:       request.GetParent(),
		ParentObjectID:     request.GetParentId(),
		RelationshipTypeID: relationshipType.GetID(),
		RelationshipType:   relationshipType,
		ChildObject:        request.GetChild(),
		ChildObjectID:      request.GetChildId(),
		Properties:         nil,
	}
	a.GenID(ctx)
	if a.ValidXID(request.GetId()) {
		a.ID = request.GetId()
	}

	err = aB.relationshipRepo.Save(ctx, &a)
	if err != nil {
		return nil, err
	}

	return aB.ToAPI(ctx, request.GetParent(), request.GetParentId(), &a)

}
func (aB *relationshipBusiness) DeleteRelationship(ctx context.Context, request *profilev1.DeleteRelationshipRequest) (*profilev1.RelationshipObject, error) {

	relationship, err := aB.relationshipRepo.GetByID(ctx, request.GetId())
	if err != nil || relationship == nil {
		return nil, err
	}

	if request.GetParentId() == "" || request.GetParentId() == relationship.ParentObjectID {

		err = aB.relationshipRepo.Delete(ctx, request.GetId())
		if err != nil {
			return nil, err
		}

		return aB.ToAPI(ctx, relationship.ParentObject, relationship.ParentObjectID, relationship)
	}

	return nil, nil
}
