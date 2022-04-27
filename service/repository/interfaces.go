package repository

import (
	"context"
	papi "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/service-profile/service/models"
)

type ProfileRepository interface {
	GetByID(ctx context.Context, id string) (*models.Profile, error)
	Save(ctx context.Context, profile *models.Profile) error
	Delete(ctx context.Context, id string) error

	GetTypeByID(ctx context.Context, profileTypeId string) (*models.ProfileType, error)
	GetTypeByUID(ctx context.Context, profileType papi.ProfileType) (*models.ProfileType, error)
}

type ContactRepository interface {
	GetByID(ctx context.Context, id string) (*models.Contact, error)
	GetByProfileID(ctx context.Context, profileId string) ([]*models.Contact, error)
	GetByDetail(ctx context.Context, detail string) (*models.Contact, error)
	Save(ctx context.Context, contact *models.Contact) (*models.Contact, error)
	Delete(ctx context.Context, id string) error

	ContactType(ctx context.Context, contactType papi.ContactType) (*models.ContactType, error)
	ContactTypeByID(ctx context.Context, contactTypeID string) (*models.ContactType, error)
	CommunicationLevel(ctx context.Context, communicationLevel papi.CommunicationLevel) (*models.CommunicationLevel, error)
	CommunicationLevelByID(ctx context.Context, communicationLevelID string) (*models.CommunicationLevel, error)

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
