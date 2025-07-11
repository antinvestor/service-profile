package business

import (
	"context"
	"errors"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/default/service"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

// Constants for pagination and limits.
const (
	// MaxRelationshipsToCheck is the maximum number of relationships to check when creating a new relationship.
	MaxRelationshipsToCheck = 2
)

type RelationshipBusiness interface {
	ListRelationships(ctx context.Context, request *profilev1.ListRelationshipRequest) ([]*models.Relationship, error)
	CreateRelationship(
		ctx context.Context,
		request *profilev1.AddRelationshipRequest,
	) (*profilev1.RelationshipObject, error)
	DeleteRelationship(
		ctx context.Context,
		request *profilev1.DeleteRelationshipRequest,
	) (*profilev1.RelationshipObject, error)

	ToAPI(
		ctx context.Context,
		relationship *models.Relationship,
		invertRelationship bool,
	) (*profilev1.RelationshipObject, error)
}

func NewRelationshipBusiness(
	_ context.Context,
	service *frame.Service,
	profileBiz ProfileBusiness,
) RelationshipBusiness {
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

func (rb *relationshipBusiness) ListRelationships(
	ctx context.Context,
	request *profilev1.ListRelationshipRequest,
) ([]*models.Relationship, error) {
	if request.GetPeerName() == "Profile" {
		profileObj, err := rb.profileBusiness.GetByID(ctx, request.GetPeerId())
		if err != nil {
			return nil, err
		}

		if profileObj == nil {
			return nil, service.ErrProfileDoesNotExist
		}
	}

	return rb.relationshipRepo.List(
		ctx,
		request.GetPeerName(),
		request.GetPeerId(),
		request.GetInvertRelation(),
		request.GetRelatedChildrenId(),
		request.GetLastRelationshipId(),
		int(request.GetCount()),
	)
}

func (rb *relationshipBusiness) CreateRelationship(
	ctx context.Context,
	request *profilev1.AddRelationshipRequest,
) (*profilev1.RelationshipObject, error) {
	logger := rb.service.Log(ctx).WithField("request", request)

	relationships, err := rb.relationshipRepo.List(
		ctx,
		request.GetParent(),
		request.GetParentId(),
		false,
		[]string{request.GetChildId()},
		"",
		MaxRelationshipsToCheck,
	)
	if err != nil {
		logger.WithError(err).Warn("get existing relationship error")

		if !frame.ErrorIsNoRows(err) {
			return nil, err
		}
	}

	if len(relationships) > 0 {
		relationship := relationships[0]

		return relationship.ToAPI(), nil
	}

	relationshipType, err := rb.relationshipRepo.RelationshipType(ctx, request.GetType())
	if err != nil {
		return nil, err
	}

	relationship := models.Relationship{
		ParentObject:       request.GetParent(),
		ParentObjectID:     request.GetParentId(),
		RelationshipTypeID: relationshipType.GetID(),
		RelationshipType:   relationshipType,
		ChildObject:        request.GetChild(),
		ChildObjectID:      request.GetChildId(),
		Properties:         frame.DBPropertiesFromMap(request.GetProperties()),
	}
	relationship.GenID(ctx)
	if relationship.ValidXID(request.GetId()) {
		relationship.ID = request.GetId()
	}

	err = rb.relationshipRepo.Save(ctx, &relationship)
	if err != nil {
		return nil, err
	}

	logger.Debug("successfully add relationship relationship")

	return relationship.ToAPI(), nil
}

func (rb *relationshipBusiness) DeleteRelationship(
	ctx context.Context,
	request *profilev1.DeleteRelationshipRequest,
) (*profilev1.RelationshipObject, error) {
	relationship, err := rb.relationshipRepo.GetByID(ctx, request.GetId())
	if err != nil {
		if frame.ErrorIsNoRows(err) {
			return nil, errors.New("relationship not found")
		}
		return nil, err
	}

	relationshipObject, err := rb.ToAPI(ctx, relationship, false)
	if err != nil {
		return nil, err
	}

	err = rb.relationshipRepo.Delete(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	return relationshipObject, nil
}

func (rb *relationshipBusiness) ToAPI(
	ctx context.Context,
	relationship *models.Relationship,
	invertRelationship bool,
) (*profilev1.RelationshipObject, error) {
	if relationship == nil {
		return nil, nil
	}

	relationshipObj := relationship.ToAPI()

	peerProfileID := ""

	if !invertRelationship {
		if relationship.ChildObject == "Profile" {
			peerProfileID = relationship.ChildObjectID
		}
	} else {
		if relationship.ParentObject == "Profile" {
			peerProfileID = relationship.ParentObjectID
		}
	}

	if peerProfileID != "" {
		profileObj, err := rb.profileBusiness.GetByID(ctx, peerProfileID)
		if err == nil {
			relationshipObj.PeerProfile = profileObj
		}

		relationshipObj.PeerProfile = profileObj
	}

	return relationshipObj, nil
}
