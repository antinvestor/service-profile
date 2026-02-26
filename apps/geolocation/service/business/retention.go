package business

import (
	"context"
	"fmt"
	"time"

	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/util"
)

// RetentionConfig holds data retention policy configuration.
type RetentionConfig struct {
	// LocationPointRetentionDays is the number of days to retain location points.
	// Points older than this are deleted. 0 disables retention.
	LocationPointRetentionDays int
	// GeoEventRetentionDays is the number of days to retain geo events.
	// 0 disables retention.
	GeoEventRetentionDays int
	// PartitionMaintenanceMonths is how many months ahead to create partitions.
	PartitionMaintenanceMonths int
	// RetentionBatchSize is the maximum number of rows to delete per batch.
	RetentionBatchSize int
	// GeofenceStateStaleDays is the number of days after which a geofence_state row
	// is considered stale and can be cleaned up. 0 disables cleanup.
	GeofenceStateStaleDays int
	// RetentionInterval is how often the retention loop runs. Default 24h.
	RetentionInterval time.Duration
}

// Retention defaults.
const (
	defaultLocationPointRetentionDays = 90
	defaultGeoEventRetentionDays      = 365
	defaultPartitionMaintenanceMonths = 3
	defaultRetentionBatchSize         = 10000
	defaultGeofenceStateStaleDays     = 30
	defaultRetentionInterval          = 24 * time.Hour
)

// RetentionBusiness handles data lifecycle management.
type RetentionBusiness interface {
	// RunRetention deletes expired data and maintains partitions.
	RunRetention(ctx context.Context) error
	// EnsurePartitions creates future partitions for partitioned tables.
	EnsurePartitions(ctx context.Context) error
	// StartScheduler runs RunRetention periodically. Blocks until ctx is cancelled.
	StartScheduler(ctx context.Context)
}

type retentionBusiness struct {
	dbPool pool.Pool
	cfg    RetentionConfig
}

// NewRetentionBusiness creates a new RetentionBusiness.
func NewRetentionBusiness(dbPool pool.Pool, cfg RetentionConfig) RetentionBusiness {
	if cfg.LocationPointRetentionDays <= 0 {
		cfg.LocationPointRetentionDays = defaultLocationPointRetentionDays
	}
	if cfg.GeoEventRetentionDays <= 0 {
		cfg.GeoEventRetentionDays = defaultGeoEventRetentionDays
	}
	if cfg.PartitionMaintenanceMonths <= 0 {
		cfg.PartitionMaintenanceMonths = defaultPartitionMaintenanceMonths
	}
	if cfg.RetentionBatchSize <= 0 {
		cfg.RetentionBatchSize = defaultRetentionBatchSize
	}
	if cfg.GeofenceStateStaleDays <= 0 {
		cfg.GeofenceStateStaleDays = defaultGeofenceStateStaleDays
	}
	if cfg.RetentionInterval <= 0 {
		cfg.RetentionInterval = defaultRetentionInterval
	}

	return &retentionBusiness{dbPool: dbPool, cfg: cfg}
}

// StartScheduler runs RunRetention on a periodic timer. Blocks until ctx is cancelled.
// Runs retention immediately on first call, then every RetentionInterval.
func (b *retentionBusiness) StartScheduler(ctx context.Context) {
	log := util.Log(ctx)

	// Run once immediately at startup.
	if err := b.RunRetention(ctx); err != nil {
		log.WithError(err).Error("initial retention run failed")
	}

	ticker := time.NewTicker(b.cfg.RetentionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("retention scheduler stopped")
			return
		case <-ticker.C:
			if err := b.RunRetention(ctx); err != nil {
				log.WithError(err).Error("scheduled retention run failed")
			}
		}
	}
}

// RunRetention deletes expired location points, geo events, and stale geofence states in batches.
// Designed to be called periodically via StartScheduler.
func (b *retentionBusiness) RunRetention(ctx context.Context) error {
	log := util.Log(ctx)

	// Delete expired location points.
	if b.cfg.LocationPointRetentionDays > 0 {
		cutoff := time.Now().AddDate(0, 0, -b.cfg.LocationPointRetentionDays)
		deleted, err := b.deleteExpired(ctx, "location_points", "ingested_at", cutoff)
		if err != nil {
			return fmt.Errorf("retain location_points: %w", err)
		}
		if deleted > 0 {
			log.Info("location_points retention complete",
				"deleted", deleted,
				"cutoff", cutoff,
			)
		}
	}

	// Delete expired geo events.
	if b.cfg.GeoEventRetentionDays > 0 {
		cutoff := time.Now().AddDate(0, 0, -b.cfg.GeoEventRetentionDays)
		deleted, err := b.deleteExpired(ctx, "geo_events", "ts", cutoff)
		if err != nil {
			return fmt.Errorf("retain geo_events: %w", err)
		}
		if deleted > 0 {
			log.Info("geo_events retention complete",
				"deleted", deleted,
				"cutoff", cutoff,
			)
		}
	}

	// Clean up stale geofence_states (subjects that stopped reporting).
	if b.cfg.GeofenceStateStaleDays > 0 {
		cutoff := time.Now().AddDate(0, 0, -b.cfg.GeofenceStateStaleDays)
		deleted, err := b.deleteExpired(ctx, "geofence_states", "updated_at", cutoff)
		if err != nil {
			return fmt.Errorf("retain geofence_states: %w", err)
		}
		if deleted > 0 {
			log.Info("geofence_states retention complete",
				"deleted", deleted,
				"cutoff", cutoff,
			)
		}
	}

	// Maintain partitions.
	if err := b.EnsurePartitions(ctx); err != nil {
		return fmt.Errorf("partition maintenance: %w", err)
	}

	return nil
}

// EnsurePartitions creates future monthly partitions for location_points.
// Safe to call repeatedly — uses IF NOT EXISTS.
func (b *retentionBusiness) EnsurePartitions(ctx context.Context) error {
	log := util.Log(ctx)
	db := b.dbPool.DB(ctx, false)

	result := db.Exec(
		"SELECT create_location_points_partitions($1)",
		b.cfg.PartitionMaintenanceMonths,
	)
	if result.Error != nil {
		// Non-fatal if the function doesn't exist (table not yet partitioned).
		log.WithError(result.Error).Warn(
			"partition maintenance skipped (function may not exist)",
		)
		return result.Error //nolint:wrapcheck // already logged with context
	}

	log.Info("partition maintenance complete",
		"months_ahead", b.cfg.PartitionMaintenanceMonths,
	)
	return nil
}

// deleteExpired deletes rows older than cutoff in batches.
// Returns the total number of deleted rows.
func (b *retentionBusiness) deleteExpired(
	ctx context.Context,
	table, column string,
	cutoff time.Time,
) (int64, error) {
	db := b.dbPool.DB(ctx, false)
	var totalDeleted int64

	for {
		// Delete in batches to avoid long-running transactions and lock contention.
		result := db.Exec(
			fmt.Sprintf(
				"DELETE FROM %s WHERE ctid IN (SELECT ctid FROM %s WHERE %s < $1 LIMIT $2)",
				table, table, column,
			),
			cutoff, b.cfg.RetentionBatchSize,
		)
		if result.Error != nil {
			return totalDeleted, fmt.Errorf("delete from %s: %w", table, result.Error)
		}

		totalDeleted += result.RowsAffected
		if result.RowsAffected < int64(b.cfg.RetentionBatchSize) {
			break
		}
	}

	return totalDeleted, nil
}
