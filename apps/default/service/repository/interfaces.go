package repository

import (
	"context"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/framedata"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)

type ProfileRepository interface {
	GetByID(ctx context.Context, id string) (*models.Profile, error)
	Search(ctx context.Context, query *framedata.SearchQuery) (frame.JobResultPipe[[]*models.Profile], error)
	Save(ctx context.Context, profile *models.Profile) error
	Delete(ctx context.Context, id string) error

	GetTypeByID(ctx context.Context, profileTypeID string) (*models.ProfileType, error)
	GetTypeByUID(ctx context.Context, profileType profilev1.ProfileType) (*models.ProfileType, error)
}

type ContactRepository interface {
	GetByID(ctx context.Context, id string) (*models.Contact, error)
	GetByProfileID(ctx context.Context, profileID string) ([]*models.Contact, error)
	GetByDetail(ctx context.Context, detail string) (*models.Contact, error)
	Save(ctx context.Context, contact *models.Contact) (*models.Contact, error)
	Delete(ctx context.Context, id string) error
	DelinkFromProfile(ctx context.Context, id, profileID string) (*models.Contact, error)
}

type VerificationRepository interface {
	GetByID(ctx context.Context, verificationID string) (*models.Verification, error)
	Save(ctx context.Context, verification *models.Verification) error

	GetAttempts(ctx context.Context, verificationID string) ([]*models.VerificationAttempt, error)
	SaveAttempt(ctx context.Context, verificationAttempt *models.VerificationAttempt) error
}

type RosterRepository interface {
	GetByID(ctx context.Context, id string) (*models.Roster, error)
	GetByContactAndProfileID(ctx context.Context, profileID, contactID string) (*models.Roster, error)
	Search(ctx context.Context, query *framedata.SearchQuery) (frame.JobResultPipe[[]*models.Roster], error)
	Save(ctx context.Context, contact *models.Roster) (*models.Roster, error)
	Delete(ctx context.Context, id string) error
}

type AddressRepository interface {
	GetByID(ctx context.Context, id string) (*models.Address, error)
	GetByNameAdminUnitAndCountry(
		ctx context.Context,
		name string,
		adminUnit string,
		countryID string,
	) (*models.Address, error)
	Save(ctx context.Context, address *models.Address) error
	Delete(ctx context.Context, id string) error

	GetByProfileID(ctx context.Context, profileID string) ([]*models.ProfileAddress, error)
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
	List(ctx context.Context,
		peerName string, peerID string,
		inverseRelation bool, relatedChildrenIDs []string,
		lastRelationshipID string, count int,
	) ([]*models.Relationship, error)

	RelationshipType(ctx context.Context, relationshipType profilev1.RelationshipType) (*models.RelationshipType, error)
	RelationshipTypeByID(ctx context.Context, relationshipTypeID string) (*models.RelationshipType, error)
}
