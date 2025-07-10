package repository

import (
	"context"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"gorm.io/gorm/clause"

	"github.com/pitabwire/frame"
)

type relationshipRepository struct {
	service *frame.Service
}

func (ar *relationshipRepository) GetByID(ctx context.Context, id string) (*models.Relationship, error) {
	relationship := &models.Relationship{}
	err := ar.service.DB(ctx, true).Preload(clause.Associations).First(relationship, "id = ?", id).Error
	return relationship, err
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

	database := ar.service.DB(ctx, true).Preload(clause.Associations)

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

func (ar *relationshipRepository) Save(ctx context.Context, relationship *models.Relationship) error {
	return ar.service.DB(ctx, false).Save(relationship).Error
}

func (ar *relationshipRepository) Delete(ctx context.Context, id string) error {
	relationship, err := ar.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return ar.service.DB(ctx, false).Delete(relationship).Error
}

func (ar *relationshipRepository) RelationshipTypeByID(
	ctx context.Context,
	profileTypeID string,
) (*models.RelationshipType, error) {
	relationshipType := &models.RelationshipType{}
	err := ar.service.DB(ctx, true).First(relationshipType, "id = ?", profileTypeID).Error
	return relationshipType, err
}

func (ar *relationshipRepository) RelationshipType(
	ctx context.Context,
	profileType profilev1.RelationshipType,
) (*models.RelationshipType, error) {
	relationshipTypeUID := models.RelationshipTypeIDMap[profileType]
	relationshipTypeM := &models.RelationshipType{}
	err := ar.service.DB(ctx, true).First(relationshipTypeM, "uid = ?", relationshipTypeUID).Error
	return relationshipTypeM, err
}

func NewRelationshipRepository(service *frame.Service) RelationshipRepository {
	repository := relationshipRepository{
		service: service,
	}
	return &repository
}
