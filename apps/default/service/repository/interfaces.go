package repository

import (
	"context"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/workerpool"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)

type ProfileRepository interface {
	datastore.BaseRepository[*models.Profile]
	Search(ctx context.Context, query *data.SearchQuery) (workerpool.JobResultPipe[[]*models.Profile], error)

	GetTypeByID(ctx context.Context, profileTypeID string) (*models.ProfileType, error)
	GetTypeByUID(ctx context.Context, profileType profilev1.ProfileType) (*models.ProfileType, error)
}

type ContactRepository interface {
	datastore.BaseRepository[*models.Contact]
	GetByProfileID(ctx context.Context, profileID string) ([]*models.Contact, error)
	GetByDetail(ctx context.Context, detail string) (*models.Contact, error)
	DelinkFromProfile(ctx context.Context, id, profileID string) (*models.Contact, error)
}

type VerificationRepository interface {
	datastore.BaseRepository[*models.Verification]

	GetAttempts(ctx context.Context, verificationID string) ([]*models.VerificationAttempt, error)
	SaveAttempt(ctx context.Context, verificationAttempt *models.VerificationAttempt) error
}

type RosterRepository interface {
	datastore.BaseRepository[*models.Roster]
	GetByContactAndProfileID(ctx context.Context, profileID, contactID string) (*models.Roster, error)
	Search(ctx context.Context, query *data.SearchQuery) (workerpool.JobResultPipe[[]*models.Roster], error)
}

type AddressRepository interface {
	datastore.BaseRepository[*models.Address]
	GetByNameAdminUnitAndCountry(
		ctx context.Context,
		name string,
		adminUnit string,
		countryID string,
	) (*models.Address, error)

	GetByProfileID(ctx context.Context, profileID string) ([]*models.ProfileAddress, error)
	SaveLink(ctx context.Context, profileAddress *models.ProfileAddress) error
	DeleteLink(ctx context.Context, id string) error

	CountryGetByISO3(ctx context.Context, countryISO3 string) (*models.Country, error)
	CountryGetByAny(ctx context.Context, c string) (*models.Country, error)
	CountryGetByName(ctx context.Context, name string) (*models.Country, error)
}

type RelationshipRepository interface {
	datastore.BaseRepository[*models.Relationship]
	List(ctx context.Context,
		peerName string, peerID string,
		inverseRelation bool, relatedChildrenIDs []string,
		lastRelationshipID string, count int,
	) ([]*models.Relationship, error)

	RelationshipType(ctx context.Context, relationshipType profilev1.RelationshipType) (*models.RelationshipType, error)
	RelationshipTypeByID(ctx context.Context, relationshipTypeID string) (*models.RelationshipType, error)
}
