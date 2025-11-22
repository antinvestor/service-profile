package business

import (
	"context"
	"errors"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/util"

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
	ListRelationships(
		ctx context.Context,
		request *profilev1.ListRelationshipRequest,
	) ([]*models.Relationship, *connect.Error)
	CreateRelationship(
		ctx context.Context,
		request *profilev1.AddRelationshipRequest,
	) (*profilev1.RelationshipObject, *connect.Error)
	DeleteRelationship(
		ctx context.Context,
		request *profilev1.DeleteRelationshipRequest,
	) (*profilev1.RelationshipObject, *connect.Error)

	ToAPI(
		ctx context.Context,
		relationship *models.Relationship,
		invertRelationship bool,
	) (*profilev1.RelationshipObject, *connect.Error)
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
) ([]*models.Relationship, *connect.Error) {
	if request.GetPeerName() == ProfilePeerName {
		profileObj, err := rb.profileBusiness.GetByID(ctx, request.GetPeerId())
		if err != nil {
			return nil, err
		}

		if profileObj == nil {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("profile does not exist"))
		}
	}

	relationships, err := rb.relationshipRepo.List(
		ctx,
		request.GetPeerName(),
		request.GetPeerId(),
		request.GetInvertRelation(),
		request.GetRelatedChildrenId(),
		request.GetLastRelationshipId(),
		int(request.GetCount()),
	)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}
	return relationships, nil
}

func (rb *relationshipBusiness) CreateRelationship(
	ctx context.Context,
	request *profilev1.AddRelationshipRequest,
) (*profilev1.RelationshipObject, *connect.Error) {
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
			return nil, data.ErrorConvertToAPI(err)
		}
	}

	if len(relationships) > 0 {
		relationship := relationships[0]

		return relationship.ToAPI(), nil
	}

	relationshipType, err := rb.relationshipRepo.RelationshipType(ctx, request.GetType())
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
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
		return nil, data.ErrorConvertToAPI(err)
	}

	logger.Debug("successfully add relationship relationship")

	return relationship.ToAPI(), nil
}

func (rb *relationshipBusiness) DeleteRelationship(
	ctx context.Context,
	request *profilev1.DeleteRelationshipRequest,
) (*profilev1.RelationshipObject, *connect.Error) {
	relationship, err := rb.relationshipRepo.GetByID(ctx, request.GetId())
	if err != nil {
		if data.ErrorIsNoRows(err) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("relationship not found"))
		}
		return nil, data.ErrorConvertToAPI(err)
	}

	relationshipObject, apiErr := rb.ToAPI(ctx, relationship, false)
	if apiErr != nil {
		return nil, apiErr
	}

	deleteErr := rb.relationshipRepo.Delete(ctx, request.GetId())
	if deleteErr != nil {
		return nil, data.ErrorConvertToAPI(deleteErr)
	}

	return relationshipObject, nil
}

// Define sentinel errors.
var (
	ErrNilRelationship = connect.NewError(connect.CodeInvalidArgument, errors.New("relationship is nil"))
)

func (rb *relationshipBusiness) ToAPI(
	ctx context.Context,
	relationship *models.Relationship,
	invertRelationship bool,
) (*profilev1.RelationshipObject, *connect.Error) {
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
	}

	return relationshipObj, nil
}
