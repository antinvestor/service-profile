package repository

import (
	"context"
	profilev1 "github.com/antinvestor/apis/profile/v1"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/pitabwire/frame"
	"gorm.io/gorm/clause"
)

type relationshipRepository struct {
	service *frame.Service
}

func (ar *relationshipRepository) GetByID(ctx context.Context, id string) (*models.Relationship, error) {
	relationship := &models.Relationship{}
	err := ar.service.DB(ctx, true).Preload(clause.Associations).First(relationship, "id = ?", id).Error
	return relationship, err
}

func (ar *relationshipRepository) List(ctx context.Context, parent, parentId string, childrenIds []string, lastRelationshipId string, count int) ([]*models.Relationship, error) {
	var relationshipList []*models.Relationship

	if count == 0 {
		count = 100
	}

	database := ar.service.DB(ctx, true).Preload(clause.Associations).
		Limit(count).Where(
		"(( parent_object = ? AND parent_object_id = ? ) OR ( child_object = ? AND child_object_id = ?)) ",
		parent, parentId, parent, parentId)

	if lastRelationshipId != "" {
		database = database.Where("id > ?", lastRelationshipId)
	}

	if len(childrenIds) > 0 {
		database = database.Where("child_object_id IN ?", childrenIds)
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

func (ar *relationshipRepository) RelationshipTypeByID(ctx context.Context, profileTypeId string) (*models.RelationshipType, error) {
	relationshipType := &models.RelationshipType{}
	err := ar.service.DB(ctx, true).First(relationshipType, "id = ?", profileTypeId).Error
	return relationshipType, err
}

func (ar *relationshipRepository) RelationshipType(ctx context.Context, profileType profilev1.RelationshipType) (*models.RelationshipType, error) {

	relationshipTypeUId := models.RelationshipTypeIDMap[profileType]
	relationshipTypeM := &models.RelationshipType{}
	err := ar.service.DB(ctx, true).First(relationshipTypeM, "uid = ?", relationshipTypeUId).Error
	return relationshipTypeM, err
}

func NewRelationshipRepository(service *frame.Service) RelationshipRepository {
	repository := relationshipRepository{
		service: service,
	}
	return &repository
}
