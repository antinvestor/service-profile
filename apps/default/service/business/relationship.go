package business

import (
	"context"
	"errors"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/default/service"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

// Constants for pagination and limits.
const (
	// MaxRelationshipsToCheck is the maximum number of relationships to check when creating a new relationship.
	MaxRelationshipsToCheck = 2
	DefaultListLimit        = 20
	MaxListLimit            = 100
	ProfilePeerName         = "Profile"
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
	profileBiz ProfileBusiness,
	relationshipRepo repository.RelationshipRepository,
) RelationshipBusiness {
	return &relationshipBusiness{
		profileBusiness:  profileBiz,
		relationshipRepo: relationshipRepo,
	}
}

type relationshipBusiness struct {
	profileBusiness  ProfileBusiness
	relationshipRepo repository.RelationshipRepository
}

func (rb *relationshipBusiness) ListRelationships(
	ctx context.Context,
	request *profilev1.ListRelationshipRequest,
) ([]*models.Relationship, error) {
	if request.GetPeerName() == ProfilePeerName {
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
	logger := util.Log(ctx).WithField("request", request)

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

		if !data.ErrorIsNoRows(err) {
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

	requestProperties := data.JSONMap{}

	relationship := models.Relationship{
		ParentObject:       request.GetParent(),
		ParentObjectID:     request.GetParentId(),
		RelationshipTypeID: relationshipType.GetID(),
		RelationshipType:   relationshipType,
		ChildObject:        request.GetChild(),
		ChildObjectID:      request.GetChildId(),
		Properties:         requestProperties.FromProtoStruct(request.GetProperties()),
	}
	relationship.GenID(ctx)
	if relationship.ValidXID(request.GetId()) {
		relationship.ID = request.GetId()
	}

	err = rb.relationshipRepo.Create(ctx, &relationship)
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
		if data.ErrorIsNoRows(err) {
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

// Define sentinel errors.
var (
	ErrNilRelationship = errors.New("relationship is nil")
)

func (rb *relationshipBusiness) ToAPI(
	ctx context.Context,
	relationship *models.Relationship,
	invertRelationship bool,
) (*profilev1.RelationshipObject, error) {
	if relationship == nil {
		return nil, ErrNilRelationship
	}

	relationshipObj := relationship.ToAPI()

	peerProfileID := ""

	if !invertRelationship {
		if relationship.ChildObject == ProfilePeerName {
			peerProfileID = relationship.ChildObjectID
		}
	} else {
		if relationship.ParentObject == ProfilePeerName {
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
