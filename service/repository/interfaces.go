package repository

import (
	"context"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service/models"
)

type ProfileRepository interface {
	GetByID(ctx context.Context, id string) (*models.Profile, error)
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

	ContactType(ctx context.Context, contactType profilev1.ContactType) (*models.ContactType, error)
	ContactTypeByID(ctx context.Context, contactTypeID string) (*models.ContactType, error)
	CommunicationLevel(ctx context.Context, communicationLevel profilev1.CommunicationLevel) (*models.CommunicationLevel, error)
	CommunicationLevelByID(ctx context.Context, communicationLevelID string) (*models.CommunicationLevel, error)

	GetVerificationByContactID(ctx context.Context, contactID string) (*models.Verification, error)
	VerificationSave(ctx context.Context, verification *models.Verification) error
	VerificationAttemptSave(ctx context.Context, attempt *models.VerificationAttempt) error
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
	List(ctx context.Context, parent string, parentId string, relatedChildrenIds []string, lastRelationshipId string, count int) ([]*models.Relationship, error)

	RelationshipType(ctx context.Context, relationshipType profilev1.RelationshipType) (*models.RelationshipType, error)
	RelationshipTypeByID(ctx context.Context, relationshipTypeID string) (*models.RelationshipType, error)
}
