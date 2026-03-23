package business

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pitabwire/util"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
)

// ProximityConfig holds tunable parameters for proximity queries.
type ProximityConfig struct {
	DefaultRadiusM float64
	MaxRadiusM     float64
	StaleHours     int
	DefaultLimit   int
	MaxLimit       int
}

// Proximity defaults used when config values are zero.
const (
	defaultProximityDefaultRadiusM = 1000.0
	defaultProximityMaxRadiusM     = 50000.0
	defaultProximityStaleHours     = 1
	defaultProximityLimit          = 50
	maxProximityLimit              = 200
)

type proximityBusiness struct {
	latestPosRepo repository.LatestPositionRepository
	areaRepo      repository.AreaRepository
	cfg           ProximityConfig
}

// NewProximityBusiness creates a new ProximityBusiness with configurable parameters.
func NewProximityBusiness(
	latestPosRepo repository.LatestPositionRepository,
	areaRepo repository.AreaRepository,
	cfg ProximityConfig,
) ProximityBusiness {
	if cfg.DefaultRadiusM <= 0 {
		cfg.DefaultRadiusM = defaultProximityDefaultRadiusM
	}
	if cfg.MaxRadiusM <= 0 {
		cfg.MaxRadiusM = defaultProximityMaxRadiusM
	}
	if cfg.StaleHours <= 0 {
		cfg.StaleHours = defaultProximityStaleHours
	}
	if cfg.DefaultLimit <= 0 {
		cfg.DefaultLimit = defaultProximityLimit
	}
	if cfg.MaxLimit <= 0 {
		cfg.MaxLimit = maxProximityLimit
	}

	return &proximityBusiness{
		latestPosRepo: latestPosRepo,
		areaRepo:      areaRepo,
		cfg:           cfg,
	}
}

// GetNearbySubjects finds subjects near the requesting subject.
func (b *proximityBusiness) GetNearbySubjects(
	ctx context.Context,
	req *models.GetNearbySubjectsRequest,
) ([]*models.NearbySubjectAPI, error) {
	log := util.Log(ctx)

	if req.GetSubjectId() == "" {
		return nil, errors.New("subject_id is required")
	}

	radiusMeters := req.GetRadiusMeters()
	if radiusMeters <= 0 {
		radiusMeters = b.cfg.DefaultRadiusM
	}
	if radiusMeters > b.cfg.MaxRadiusM {
		return nil, fmt.Errorf("radius_meters %f exceeds maximum %f", radiusMeters, b.cfg.MaxRadiusM)
	}

	limit := clampLimit(int(req.GetLimit()), b.cfg.DefaultLimit, b.cfg.MaxLimit)

	// Get the requesting subject's latest position.
	pos, err := b.latestPosRepo.Get(ctx, req.GetSubjectId())
	if err != nil {
		return nil, fmt.Errorf("get position for subject %s: %w", req.GetSubjectId(), err)
	}

	// Find nearby subjects excluding the requester.
	results, err := b.latestPosRepo.GetNearbySubjects(
		ctx,
		pos.Latitude,
		pos.Longitude,
		radiusMeters,
		req.GetSubjectId(),
		b.cfg.StaleHours,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("get nearby subjects: %w", err)
	}

	log.Debug("nearby subjects query completed",
		"subject_id", req.GetSubjectId(),
		"radius_m", radiusMeters,
		"results", len(results),
	)

	apiResults := make([]*models.NearbySubjectAPI, 0, len(results))
	for _, r := range results {
		apiResults = append(apiResults, &models.NearbySubjectAPI{
			SubjectId:      r.SubjectID,
			DistanceMeters: r.DistanceMeters,
			LastSeen:       timestampFromTime(r.LastSeen),
		})
	}
	return apiResults, nil
}

// GetNearbyAreas finds areas near the given point.
func (b *proximityBusiness) GetNearbyAreas(
	ctx context.Context,
	req *models.GetNearbyAreasRequest,
) ([]*models.NearbyAreaAPI, error) {
	log := util.Log(ctx)

	if req == nil {
		return nil, errors.New("request is nil")
	}
	if err := models.ValidateLatLon(req.GetLatitude(), req.GetLongitude()); err != nil {
		return nil, fmt.Errorf("invalid coordinates: %w", err)
	}

	radiusMeters := req.GetRadiusMeters()
	if radiusMeters <= 0 {
		radiusMeters = b.cfg.DefaultRadiusM
	}
	if radiusMeters > b.cfg.MaxRadiusM {
		return nil, fmt.Errorf("radius_meters %f exceeds maximum %f", radiusMeters, b.cfg.MaxRadiusM)
	}

	limit := clampLimit(int(req.GetLimit()), b.cfg.DefaultLimit, b.cfg.MaxLimit)

	results, err := b.areaRepo.GetNearbyAreas(ctx, req.GetLatitude(), req.GetLongitude(), radiusMeters, limit)
	if err != nil {
		return nil, fmt.Errorf("get nearby areas: %w", err)
	}

	log.Debug("nearby areas query completed",
		"lat", req.GetLatitude(),
		"lon", req.GetLongitude(),
		"radius_m", radiusMeters,
		"results", len(results),
	)

	apiResults := make([]*models.NearbyAreaAPI, 0, len(results))
	for _, r := range results {
		apiResults = append(apiResults, &models.NearbyAreaAPI{
			AreaId:         r.Area.GetID(),
			Name:           r.Area.Name,
			AreaType:       models.ToProtoAreaType(r.Area.AreaType),
			DistanceMeters: r.DistanceMeters,
		})
	}
	return apiResults, nil
}

// UpdateLatestPosition updates the materialized latest position for a subject.
func (b *proximityBusiness) UpdateLatestPosition(ctx context.Context, event *models.LocationPointIngestedEvent) error {
	if event == nil {
		return errors.New("event is nil")
	}

	pos := &models.LatestPosition{
		SubjectID: event.SubjectID,
		DeviceID:  event.DeviceID,
		Latitude:  event.Latitude,
		Longitude: event.Longitude,
		Accuracy:  event.Accuracy,
		TS:        time.UnixMilli(event.Timestamp),
	}
	pos.TenantID = event.TenantID
	pos.PartitionID = event.PartitionID
	pos.AccessID = event.AccessID
	pos.GenID(ctx)

	if err := b.latestPosRepo.Upsert(ctx, pos); err != nil {
		return fmt.Errorf("upsert latest position for subject %s: %w", event.SubjectID, err)
	}

	return nil
}

// clampLimit ensures a limit is within the given bounds.
func clampLimit(limit, defaultLimit, maxLimit int) int {
	if limit <= 0 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

// timestampFromTime is a helper to create a protobuf timestamp from a Go time.
func timestampFromTime(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
