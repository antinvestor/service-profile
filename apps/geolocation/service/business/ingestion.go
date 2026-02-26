package business

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pitabwire/frame/events"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
)

// LocationPointIngestedEventName is the internal frame event name for ingested points.
const LocationPointIngestedEventName = "location.point.ingested"

// IngestionConfig holds tunable parameters for the ingestion pipeline.
type IngestionConfig struct {
	MaxBatchSize int
}

// Ingestion defaults.
const (
	defaultMaxBatchSize = 1000
	defaultSearchLimit  = 50
)

type ingestionBusiness struct {
	eventsMan events.Manager
	pointRepo repository.LocationPointRepository
	cfg       IngestionConfig
}

// NewIngestionBusiness creates a new IngestionBusiness with configurable parameters.
func NewIngestionBusiness(
	eventsMan events.Manager,
	pointRepo repository.LocationPointRepository,
	cfg IngestionConfig,
) IngestionBusiness {
	if cfg.MaxBatchSize <= 0 {
		cfg.MaxBatchSize = defaultMaxBatchSize
	}
	return &ingestionBusiness{
		eventsMan: eventsMan,
		pointRepo: pointRepo,
		cfg:       cfg,
	}
}

// IngestBatch validates, normalizes, persists, and emits events for a batch of location points.
// Uses batch INSERT for throughput rather than sequential per-point writes.
func (b *ingestionBusiness) IngestBatch(
	ctx context.Context,
	req *models.IngestLocationsRequest,
) (*models.IngestLocationsResponse, error) {
	log := util.Log(ctx)

	if req == nil {
		return nil, errors.New("request is nil")
	}
	if err := models.ValidateSubjectID(req.SubjectID); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	if len(req.Points) == 0 {
		return &models.IngestLocationsResponse{Accepted: 0, Rejected: 0}, nil
	}
	if len(req.Points) > b.cfg.MaxBatchSize {
		return nil, fmt.Errorf(
			"batch size %d exceeds maximum %d",
			len(req.Points),
			b.cfg.MaxBatchSize,
		)
	}

	now := time.Now()

	// Phase 1: Validate and build domain models for all valid points.
	validPoints := make([]*models.LocationPoint, 0, len(req.Points))
	validAPIs := make([]*models.LocationPointAPI, 0, len(req.Points))
	var rejected int32

	for _, pt := range req.Points {
		pt.SubjectID = req.SubjectID

		if err := models.ValidateLocationPoint(pt); err != nil {
			log.Debug(
				"rejecting location point",
				"reason",
				err.Error(),
				"subject_id",
				req.SubjectID,
			)
			rejected++
			continue
		}

		point := buildLocationPoint(ctx, pt, req.SubjectID, now)
		validPoints = append(validPoints, point)
		validAPIs = append(validAPIs, pt)
	}

	if len(validPoints) == 0 {
		return &models.IngestLocationsResponse{
			Accepted: 0,
			Rejected: rejected,
		}, nil
	}

	// Phase 2: Batch INSERT all valid points.
	if err := b.pointRepo.BulkCreate(ctx, validPoints); err != nil {
		return nil, fmt.Errorf("batch insert location points: %w", err)
	}

	// Phase 3: Emit events for all persisted points.
	// Event emission failures are non-fatal; the detection engine can catch up from the database.
	for i, point := range validPoints {
		pt := validAPIs[i]
		ts := point.TS

		ingestedEvent := &models.LocationPointIngestedEvent{
			PointID:   point.GetID(),
			SubjectID: req.SubjectID,
			Latitude:  pt.Latitude,
			Longitude: pt.Longitude,
			Accuracy:  pt.Accuracy,
			Timestamp: ts.UnixMilli(),
		}

		if emitErr := b.eventsMan.Emit(ctx, LocationPointIngestedEventName, ingestedEvent); emitErr != nil {
			log.WithError(emitErr).Error("failed to emit location point ingested event",
				"subject_id", req.SubjectID,
				"point_id", point.GetID(),
			)
		}
	}

	accepted := min(
		int32(len(validPoints)),   //nolint:gosec // bounded by MaxBatchSize (<< MaxInt32)
		int32(b.cfg.MaxBatchSize), //nolint:gosec // config default 1000, validated at startup
	)
	log.Info("batch ingestion complete",
		"subject_id", req.SubjectID,
		"accepted", accepted,
		"rejected", rejected,
	)

	return &models.IngestLocationsResponse{
		Accepted: accepted,
		Rejected: rejected,
	}, nil
}

// buildLocationPoint normalizes an API point into a domain model ready for persistence.
func buildLocationPoint(
	ctx context.Context,
	pt *models.LocationPointAPI,
	subjectID string,
	now time.Time,
) *models.LocationPoint {
	ts := now
	if pt.Timestamp != nil {
		ts = pt.Timestamp.AsTime()
	}

	point := &models.LocationPoint{
		SubjectID:  subjectID,
		TS:         ts,
		IngestedAt: now,
		Latitude:   pt.Latitude,
		Longitude:  pt.Longitude,
		Accuracy:   pt.Accuracy,
		Source:     pt.Source,
		Extras:     models.StructToJSONMap(pt.Extras),
	}

	if pt.Altitude != 0 {
		alt := pt.Altitude
		point.Altitude = &alt
	}
	if pt.Speed != 0 {
		spd := pt.Speed
		point.Speed = &spd
	}
	if pt.Bearing != 0 {
		brg := pt.Bearing
		point.Bearing = &brg
	}

	point.GenID(ctx)

	return point
}
