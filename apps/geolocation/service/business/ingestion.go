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
	subjectID := req.GetSubjectId()
	if err := models.ValidateSubjectID(subjectID); err != nil {
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
	validInputs := make([]*models.LocationPointInput, 0, len(req.Points))
	var rejected int32

	for _, pt := range req.Points {
		if err := models.ValidateLocationPoint(pt); err != nil {
			log.Debug(
				"rejecting location point",
				"reason",
				err.Error(),
				"subject_id",
				subjectID,
			)
			rejected++
			continue
		}

		point := buildLocationPoint(ctx, pt, subjectID, now)
		validPoints = append(validPoints, point)
		validInputs = append(validInputs, pt)
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
		pt := validInputs[i]
		ts := point.TrueCreatedAt

		ingestedEvent := &models.LocationPointIngestedEvent{
			PointID: point.GetID(),
			EventTenancy: models.EventTenancy{
				TenantID:    point.TenantID,
				PartitionID: point.PartitionID,
				AccessID:    point.AccessID,
			},
			SubjectID: subjectID,
			DeviceID:  point.DeviceID,
			Latitude:  pt.Latitude,
			Longitude: pt.Longitude,
			Accuracy:  pt.Accuracy,
			Timestamp: ts.UnixMilli(),
		}

		if emitErr := b.eventsMan.Emit(ctx, LocationPointIngestedEventName, ingestedEvent); emitErr != nil {
			log.WithError(emitErr).Error("failed to emit location point ingested event",
				"subject_id", subjectID,
				"point_id", point.GetID(),
			)
		}
	}

	accepted := min(
		int32(len(validPoints)),   //nolint:gosec // bounded by MaxBatchSize (<< MaxInt32)
		int32(b.cfg.MaxBatchSize), //nolint:gosec // config default 1000, validated at startup
	)
	log.Info("batch ingestion complete",
		"subject_id", subjectID,
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
	pt *models.LocationPointInput,
	subjectID string,
	now time.Time,
) *models.LocationPoint {
	ts := now
	if pt.Timestamp != nil {
		ts = pt.Timestamp.AsTime()
	}

	point := &models.LocationPoint{
		SubjectID:       subjectID,
		DeviceID:        pt.GetDeviceId(),
		TrueCreatedAt:   ts,
		IngestedAt:      now,
		Latitude:        pt.Latitude,
		Longitude:       pt.Longitude,
		Accuracy:        pt.Accuracy,
		Source:          models.LocationSourceFromProto(pt.Source),
		Extras:          models.StructToJSONMap(pt.Extra),
		ProcessingState: models.LocationPointProcessingStatePending,
	}

	if pt.Altitude != nil {
		point.Altitude = pt.Altitude
	}
	if pt.Speed != nil {
		point.Speed = pt.Speed
	}
	if pt.Bearing != nil {
		point.Bearing = pt.Bearing
	}

	point.GenID(ctx)

	return point
}
