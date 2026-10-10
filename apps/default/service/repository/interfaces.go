package repository

import (
	"context"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"github.com/pitabwire/frame/v2/data"
	"github.com/pitabwire/frame/v2/datastore"
	"github.com/pitabwire/frame/v2/workerpool"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/internal/outbox"
)

type ProfileRepository interface {
	datastore.BaseRepository[*models.Profile]
	Search(
		ctx context.Context,
		query *data.SearchQuery,
	) (workerpool.JobResultPipe[[]*models.Profile], error)

	GetTypeByID(ctx context.Context, profileTypeID string) (*models.ProfileType, error)
	GetTypeByUID(
		ctx context.Context,
		profileType profilev1.ProfileType,
	) (*models.ProfileType, error)
	CreateWithFact(
		ctx context.Context,
		profile *models.Profile,
		fact func(*models.Profile) *outbox.Event,
	) error
	CreateWithAccount(
		ctx context.Context,
		profile *models.Profile,
		account *models.ProfileAccount,
		facts func(*models.Profile) []*outbox.Event,
	) error
	// Merge folds merging into target in one transaction: inTx runs first
	// on that transaction (moving accounts) and returns the facts to stage;
	// then target's properties are saved and merging is deleted.
	Merge(
		ctx context.Context,
		target, merging *models.Profile,
		inTx func(tx *gorm.DB) ([]*outbox.Event, error),
	) error
}

type ContactRepository interface {
	datastore.BaseRepository[*models.Contact]
	GetByProfileID(ctx context.Context, profileID string) ([]*models.Contact, error)
	GetByLookupToken(ctx context.Context, lookUpToken ...[]byte) ([]*models.Contact, error)
	DelinkFromProfile(ctx context.Context, id, profileID string) (*models.Contact, error)

	// GetByIDFromPrimary reads from the primary connection for
	// read-your-writes (e.g. linking a just-created contact to a profile).
	GetByIDFromPrimary(ctx context.Context, id string) (*models.Contact, error)
	// GetByIDs loads contacts by id in one query (standalone or attached).
	GetByIDs(ctx context.Context, ids []string) ([]*models.Contact, error)
}

type VerificationRepository interface {
	datastore.BaseRepository[*models.Verification]

	GetAttempts(ctx context.Context, verificationID string) ([]*models.VerificationAttempt, error)
	SaveAttempt(ctx context.Context, verificationAttempt *models.VerificationAttempt) error
}

type RosterRepository interface {
	datastore.BaseRepository[*models.Roster]
	GetByContactAndProfileID(
		ctx context.Context,
		profileID, contactID string,
	) (*models.Roster, error)
	GetByContactIDsAndProfileID(
		ctx context.Context,
		contactIDs []string,
		profileID string,
	) ([]*models.Roster, error)
	GetByContactAndProfileIDAndName(
		ctx context.Context,
		profileID, contactID, name string,
	) (*models.Roster, error)
	GetByContactIDsAndProfileIDAndName(
		ctx context.Context,
		contactIDs []string,
		profileID, name string,
	) ([]*models.Roster, error)
	Search(
		ctx context.Context,
		query *data.SearchQuery,
	) (workerpool.JobResultPipe[[]*models.Roster], error)
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
	CountryGetByISO2(ctx context.Context, iso2 string) (*models.Country, error)
	CountryGetByAny(ctx context.Context, c string) (*models.Country, error)
	CountryGetByName(ctx context.Context, name string) (*models.Country, error)
}

type PropertyEntryRepository interface {
	datastore.BaseRepository[*models.PropertyEntry]
	AppendEntries(ctx context.Context, entries []*models.PropertyEntry) error
	LatestGlobalByProfile(ctx context.Context, profileID string) ([]*models.PropertyEntry, error)
	LatestScopedByProfileAndPartition(
		ctx context.Context,
		profileID, partitionID string,
	) ([]*models.PropertyEntry, error)
	HistoryByKey(
		ctx context.Context,
		profileID, key, callerTenantID string,
	) ([]*models.PropertyEntry, error)
}

type RelationshipRepository interface {
	datastore.BaseRepository[*models.Relationship]
	List(ctx context.Context,
		peerName string, peerID string,
		inverseRelation bool, relatedChildrenIDs []string,
		lastRelationshipID string, count int,
	) ([]*models.Relationship, error)

	RelationshipType(
		ctx context.Context,
		relationshipType profilev1.RelationshipType,
	) (*models.RelationshipType, error)
	RelationshipTypeByID(
		ctx context.Context,
		relationshipTypeID string,
	) (*models.RelationshipType, error)
}

type ProfileAccountRepository interface {
	datastore.BaseRepository[*models.ProfileAccount]
	ListByProfileID(ctx context.Context, profileID string) ([]*models.ProfileAccount, error)
	ListByAddresses(ctx context.Context, addresses [][]byte) ([]*models.ProfileAccount, error)
	CreateWithFact(
		ctx context.Context,
		account *models.ProfileAccount,
		fact func(*models.ProfileAccount) *outbox.Event,
	) (bool, error)
	// MoveToProfile runs on the caller's transaction and never commits.
	MoveToProfile(tx *gorm.DB, from, to *models.Profile) ([]*models.ProfileAccount, error)
	PersonProfilesWithoutAccount(
		ctx context.Context,
		profileTypeID, family string,
		accountVersion uint32,
		afterID string,
		limit int,
	) ([]*models.Profile, error)
}
