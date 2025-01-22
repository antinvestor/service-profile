package repository

import (
	"context"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/pitabwire/frame"
	"time"
)

const defaultBatchSize = 30

type SearchQuery struct {
	ProfileID            string
	Query                string
	PropertiesToSearchOn []string
	StartAt              *time.Time
	EndAt                *time.Time

	Pagination *Paginator
}

func NewSearchQuery(ctx context.Context, query string, props []string, startAt, endAt string, resultPage, resultCount int) (*SearchQuery, error) {

	if resultCount == 0 {
		resultCount = defaultBatchSize
	}

	profileID := ""
	claims := frame.ClaimsFromContext(ctx)
	if claims != nil {
		profileID, _ = claims.GetSubject()
	}

	sq := &SearchQuery{
		ProfileID:            profileID,
		Query:                query,
		PropertiesToSearchOn: props,
		Pagination: &Paginator{
			offset:    resultPage * resultCount,
			limit:     resultCount,
			batchSize: defaultBatchSize,
		},
	}

	if startAt != "" {

		parsedTime, err := time.Parse(time.DateTime, startAt)
		if err != nil {
			return nil, err
		}
		sq.StartAt = &parsedTime
	}

	if endAt != "" {

		parsedTime, err := time.Parse(time.DateTime, endAt)
		if err != nil {
			return nil, err
		}
		sq.EndAt = &parsedTime
	}

	return sq, nil
}

type Paginator struct {
	offset int
	limit  int

	batchSize int
}

func (sq *Paginator) canLoad() bool {
	return sq.offset < sq.limit
}

func (sq *Paginator) stop(loadedCount int) bool {
	sq.offset += loadedCount
	if sq.offset+sq.batchSize > sq.limit {
		sq.batchSize = sq.limit - sq.offset
	}

	return loadedCount < sq.batchSize
}

type ProfileRepository interface {
	GetByID(ctx context.Context, id string) (*models.Profile, error)
	Search(ctx context.Context, query *SearchQuery) (frame.JobResultPipe, error)
	Save(ctx context.Context, profile *models.Profile) error
	Delete(ctx context.Context, id string) error

	GetTypeByID(ctx context.Context, profileTypeId string) (*models.ProfileType, error)
	GetTypeByUID(ctx context.Context, profileType profilev1.ProfileType) (*models.ProfileType, error)
}

type ContactRepository interface {
	GetByID(ctx context.Context, id string) (*models.Contact, error)
	GetByProfileID(ctx context.Context, profileId string) ([]*models.Contact, error)
	GetByDetail(ctx context.Context, detail string) (*models.Contact, error)
	Save(ctx context.Context, contact *models.Contact) (*models.Contact, error)
	Delete(ctx context.Context, id string) error
	DelinkFromProfile(ctx context.Context, id, profileID string) (*models.Contact, error)

	GetVerificationByContactID(ctx context.Context, contactID string) (*models.Verification, error)
	VerificationSave(ctx context.Context, verification *models.Verification) error
	VerificationAttemptSave(ctx context.Context, attempt *models.VerificationAttempt) error
}

type RosterRepository interface {
	GetByID(ctx context.Context, id string) (*models.Roster, error)
	GetByContactAndProfileID(ctx context.Context, profileID, contactID string) (*models.Roster, error)
	Search(ctx context.Context, query *SearchQuery) (frame.JobResultPipe, error)
	Save(ctx context.Context, contact *models.Roster) (*models.Roster, error)
	Delete(ctx context.Context, id string) error
}

type AddressRepository interface {
	GetByID(ctx context.Context, id string) (*models.Address, error)
	GetByNameAdminUnitAndCountry(ctx context.Context, name string, adminUnit string, countryID string) (*models.Address, error)
	Save(ctx context.Context, address *models.Address) error
	Delete(ctx context.Context, id string) error

	GetByProfileID(ctx context.Context, profileId string) ([]*models.ProfileAddress, error)
	SaveLink(ctx context.Context, profileAddress *models.ProfileAddress) error
	DeleteLink(ctx context.Context, id string) error

	CountryGetByISO3(ctx context.Context, countryISO3 string) (*models.Country, error)
	CountryGetByAny(ctx context.Context, c string) (*models.Country, error)
	CountryGetByName(ctx context.Context, name string) (*models.Country, error)
}

type RelationshipRepository interface {
	GetByID(ctx context.Context, id string) (*models.Relationship, error)
	Save(ctx context.Context, relationship *models.Relationship) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, peerName string, peerId string, inverseRelation bool, relatedChildrenIds []string, lastRelationshipId string, count int) ([]*models.Relationship, error)

	RelationshipType(ctx context.Context, relationshipType profilev1.RelationshipType) (*models.RelationshipType, error)
	RelationshipTypeByID(ctx context.Context, relationshipTypeID string) (*models.RelationshipType, error)
}

type DeviceRepository interface {
	Save(ctx context.Context, device *models.Device) error
	GetByID(ctx context.Context, id string) (*models.Device, error)
	GetByLinkID(ctx context.Context, linkId string) (*models.Device, error)
	List(ctx context.Context, profileId string) ([]*models.Device, error)
	ListByEmbedding(ctx context.Context, embedding []float32) ([]*models.Device, error)
}

type DeviceLogRepository interface {
	Save(ctx context.Context, deviceLog *models.DeviceLog) error
	GetByID(ctx context.Context, id string) (*models.DeviceLog, error)
	GetByLinkID(ctx context.Context, linkID string) (*models.DeviceLog, error)
	ListByDeviceID(ctx context.Context, deviceLogID string) ([]*models.DeviceLog, error)
}
