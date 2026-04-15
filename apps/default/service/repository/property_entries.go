package repository

import (
	"context"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/workerpool"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)

type propertyEntryRepository struct {
	datastore.BaseRepository[*models.PropertyEntry]
}

func NewPropertyEntryRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) PropertyEntryRepository {
	return &propertyEntryRepository{
		BaseRepository: datastore.NewBaseRepository[*models.PropertyEntry](
			ctx, dbPool, workMan, func() *models.PropertyEntry { return &models.PropertyEntry{} },
		),
	}
}

// AppendEntries writes property entries with raw ctx to stamp tenant provenance.
func (r *propertyEntryRepository) AppendEntries(ctx context.Context, entries []*models.PropertyEntry) error {
	return r.Pool().DB(ctx, false).Create(&entries).Error
}

// LatestGlobalByProfile returns the latest global (scoped=false) entry per key.
func (r *propertyEntryRepository) LatestGlobalByProfile(ctx context.Context, profileID string) ([]*models.PropertyEntry, error) {
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	var entries []*models.PropertyEntry
	err := r.Pool().DB(unscopedCtx, true).
		Raw(`SELECT DISTINCT ON (key) * FROM property_entries
			 WHERE profile_id = ? AND scoped = FALSE AND deleted_at IS NULL
			 ORDER BY key, created_at DESC`, profileID).
		Scan(&entries).Error
	return entries, err
}

// LatestScopedByProfileAndPartition returns latest scoped entries per key for a partition.
func (r *propertyEntryRepository) LatestScopedByProfileAndPartition(ctx context.Context, profileID, partitionID string) ([]*models.PropertyEntry, error) {
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	var entries []*models.PropertyEntry
	err := r.Pool().DB(unscopedCtx, true).
		Raw(`SELECT DISTINCT ON (key) * FROM property_entries
			 WHERE profile_id = ? AND scoped = TRUE AND partition_id = ? AND deleted_at IS NULL
			 ORDER BY key, created_at DESC`, profileID, partitionID).
		Scan(&entries).Error
	return entries, err
}

// HistoryByKey returns all entries for a profile+key, most recent first.
// Scoped entries are filtered to the caller's tenant.
func (r *propertyEntryRepository) HistoryByKey(ctx context.Context, profileID, key, callerTenantID string) ([]*models.PropertyEntry, error) {
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	var entries []*models.PropertyEntry
	err := r.Pool().DB(unscopedCtx, true).
		Where("profile_id = ? AND key = ? AND (scoped = FALSE OR tenant_id = ?) AND deleted_at IS NULL",
			profileID, key, callerTenantID).
		Order("created_at DESC").
		Find(&entries).Error
	return entries, err
}
