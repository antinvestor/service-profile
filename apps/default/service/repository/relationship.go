package repository

import (
	"context"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"
	"gorm.io/gorm/clause"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)

type relationshipRepository struct {
	datastore.BaseRepository[*models.Relationship]
}

func NewRelationshipRepository(
	ctx context.Context,
	dbPool pool.Pool,
	workMan workerpool.Manager,
) RelationshipRepository {
	repository := relationshipRepository{
		BaseRepository: datastore.NewBaseRepository[*models.Relationship](
			ctx, dbPool, workMan, func() *models.Relationship { return &models.Relationship{} },
		),
	}
	return &repository
}

func (ar *relationshipRepository) List(
	ctx context.Context,
	peerName, peerID string,
	inverseRelation bool,
	childrenIDs []string,
	lastRelationshipID string,
	count int,
) ([]*models.Relationship, error) {
	var relationshipList []*models.Relationship

	database := ar.Pool().DB(ctx, true).Preload(clause.Associations)

	if count > 0 {
		database = database.Limit(count)
	}

	if inverseRelation {
		database = database.Where(
			" child_object = ? AND child_object_id = ? ",
			peerName, peerID)
	} else {
		database = database.Where(
			" parent_object = ? AND parent_object_id = ? ",
			peerName, peerID)
	}
	if lastRelationshipID != "" {
		database = database.Where("id > ?", lastRelationshipID)
	}

	if len(childrenIDs) > 0 {
		database = database.Where("child_object_id IN ?", childrenIDs)
	}

	err := database.Find(&relationshipList).Error
	return relationshipList, err
}

func (ar *relationshipRepository) RelationshipTypeByID(
	ctx context.Context,
	profileTypeID string,
) (*models.RelationshipType, error) {
	relationshipType := &models.RelationshipType{}
	err := ar.Pool().DB(ctx, true).First(relationshipType, "id = ?", profileTypeID).Error
	return relationshipType, err
}

func (ar *relationshipRepository) RelationshipType(
	ctx context.Context,
	profileType profilev1.RelationshipType,
) (*models.RelationshipType, error) {
	relationshipTypeUID := models.RelationshipTypeIDMap[profileType]
	relationshipTypeM := &models.RelationshipType{}
	err := ar.Pool().DB(ctx, true).First(relationshipTypeM, "uid = ?", relationshipTypeUID).Error
	return relationshipTypeM, err
}
