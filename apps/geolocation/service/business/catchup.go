package business

import (
	"context"
	"fmt"
	"time"

	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/events"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

// CatchupConfig holds configuration for the ingestion catch-up mechanism.
type CatchupConfig struct {
	// LookbackDuration is how far back to look for unprocessed points.
	// Default 1 hour — points older than this are assumed processed or too stale to matter.
	LookbackDuration time.Duration
	// BatchSize is how many points to process per catch-up batch.
	BatchSize int
}

// Catchup defaults.
const (
	defaultCatchupLookback  = 1 * time.Hour
	defaultCatchupBatchSize = 500
)

// CatchupBusiness detects and re-emits events for location points
// that were persisted but never processed (e.g., due to a crash between INSERT and event emission).
type CatchupBusiness interface {
	// RunCatchup finds unprocessed location points and re-emits their events.
	// Safe to call at startup or periodically.
	RunCatchup(ctx context.Context) error
}

type catchupBusiness struct {
	dbPool    pool.Pool
	eventsMan events.Manager
	cfg       CatchupConfig
}

// NewCatchupBusiness creates a new CatchupBusiness.
func NewCatchupBusiness(
	dbPool pool.Pool,
	eventsMan events.Manager,
	cfg CatchupConfig,
) CatchupBusiness {
	if cfg.LookbackDuration <= 0 {
		cfg.LookbackDuration = defaultCatchupLookback
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = defaultCatchupBatchSize
	}
	return &catchupBusiness{dbPool: dbPool, eventsMan: eventsMan, cfg: cfg}
}

// RunCatchup finds location points ingested within the lookback window that have no
// corresponding latest_position record (or whose latest_position is older than the point).
// This detects the crash-recovery gap where points were INSERTed but events were never emitted.
func (b *catchupBusiness) RunCatchup(ctx context.Context) error {
	log := util.Log(ctx)
	db := b.dbPool.DB(ctx, true)

	cutoff := time.Now().Add(-b.cfg.LookbackDuration)

	// Find points that were ingested recently but appear unprocessed:
	// - The point's subject has no latest_position record, OR
	// - The latest_position.ts is older than the point's ts (meaning this point was never processed).
	// We limit to one point per subject (the most recent) to avoid flooding events.
	const query = `
		SELECT DISTINCT ON (lp.subject_id)
			lp.id AS point_id,
			lp.subject_id,
			lp.latitude,
			lp.longitude,
			lp.accuracy,
			lp.ts
		FROM location_points lp
		LEFT JOIN latest_positions lpos ON lp.subject_id = lpos.subject_id
		WHERE lp.ingested_at > $1
		  AND (lpos.subject_id IS NULL OR lpos.ts < lp.ts)
		ORDER BY lp.subject_id, lp.ts DESC
		LIMIT $2`

	type unprocessedPoint struct {
		PointID   string    `gorm:"column:point_id"`
		SubjectID string    `gorm:"column:subject_id"`
		Latitude  float64   `gorm:"column:latitude"`
		Longitude float64   `gorm:"column:longitude"`
		Accuracy  float64   `gorm:"column:accuracy"`
		TS        time.Time `gorm:"column:ts"`
	}

	var points []unprocessedPoint
	if err := db.Raw(query, cutoff, b.cfg.BatchSize).Scan(&points).Error; err != nil {
		return fmt.Errorf("catchup query: %w", err)
	}

	if len(points) == 0 {
		return nil
	}

	log.Info("catchup: found unprocessed location points",
		"count", len(points),
		"lookback", b.cfg.LookbackDuration,
	)

	var emitted int
	for _, pt := range points {
		event := &models.LocationPointIngestedEvent{
			PointID:   pt.PointID,
			SubjectID: pt.SubjectID,
			Latitude:  pt.Latitude,
			Longitude: pt.Longitude,
			Accuracy:  pt.Accuracy,
			Timestamp: pt.TS.UnixMilli(),
		}

		if err := b.eventsMan.Emit(
			ctx, LocationPointIngestedEventName, event,
		); err != nil {
			log.WithError(err).Error("catchup: failed to re-emit event",
				"point_id", pt.PointID,
				"subject_id", pt.SubjectID,
			)
			continue
		}
		emitted++
	}

	log.Info("catchup complete",
		"found", len(points),
		"emitted", emitted,
	)

	return nil
}
