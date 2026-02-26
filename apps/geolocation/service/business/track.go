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

// TrackConfig holds tunable parameters for track/history queries.
type TrackConfig struct {
	DefaultTrackLimit int
	MaxTrackLimit     int
	DefaultEventLimit int
	MaxEventLimit     int
	MaxAreaSubjects   int
}

// Track defaults used when config values are zero.
const (
	defaultTrackQueryLimit = 100
	maxTrackQueryLimit     = 1000
	defaultEventQueryLimit = 100
	maxEventQueryLimit     = 500
	defaultMaxAreaSubjects = 1000
)

type trackBusiness struct {
	pointRepo    repository.LocationPointRepository
	geoEventRepo repository.GeoEventRepository
	stateRepo    repository.GeofenceStateRepository
	cfg          TrackConfig
}

// NewTrackBusiness creates a new TrackBusiness with configurable parameters.
func NewTrackBusiness(
	pointRepo repository.LocationPointRepository,
	geoEventRepo repository.GeoEventRepository,
	stateRepo repository.GeofenceStateRepository,
	cfg TrackConfig,
) TrackBusiness {
	if cfg.DefaultTrackLimit <= 0 {
		cfg.DefaultTrackLimit = defaultTrackQueryLimit
	}
	if cfg.MaxTrackLimit <= 0 {
		cfg.MaxTrackLimit = maxTrackQueryLimit
	}
	if cfg.DefaultEventLimit <= 0 {
		cfg.DefaultEventLimit = defaultEventQueryLimit
	}
	if cfg.MaxEventLimit <= 0 {
		cfg.MaxEventLimit = maxEventQueryLimit
	}
	if cfg.MaxAreaSubjects <= 0 {
		cfg.MaxAreaSubjects = defaultMaxAreaSubjects
	}

	return &trackBusiness{
		pointRepo:    pointRepo,
		geoEventRepo: geoEventRepo,
		stateRepo:    stateRepo,
		cfg:          cfg,
	}
}

// GetTrack retrieves the location history for a subject within a time range.
func (b *trackBusiness) GetTrack(ctx context.Context, req *models.GetTrackRequest) ([]*models.LocationPointAPI, error) {
	log := util.Log(ctx)

	if req == nil || req.SubjectID == "" {
		return nil, errors.New("subject_id is required")
	}

	from, to := resolveTimeRange(req.From, req.To)
	limit := clampLimit(int(req.Limit), b.cfg.DefaultTrackLimit, b.cfg.MaxTrackLimit)
	offset := max(int(req.Offset), 0)

	points, err := b.pointRepo.GetTrack(ctx, req.SubjectID, from, to, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get track: %w", err)
	}

	log.Debug("track query completed",
		"subject_id", req.SubjectID,
		"from", from,
		"to", to,
		"points", len(points),
	)

	result := make([]*models.LocationPointAPI, 0, len(points))
	for _, p := range points {
		result = append(result, p.ToAPI())
	}
	return result, nil
}

// GetSubjectEvents retrieves geo events (enter/exit/dwell) for a subject.
func (b *trackBusiness) GetSubjectEvents(
	ctx context.Context,
	req *models.GetSubjectEventsRequest,
) ([]*models.GeoEventAPI, error) {
	log := util.Log(ctx)

	if req == nil || req.SubjectID == "" {
		return nil, errors.New("subject_id is required")
	}

	var from, to *time.Time
	if req.From != nil {
		t := req.From.AsTime()
		from = &t
	}
	if req.To != nil {
		t := req.To.AsTime()
		to = &t
	}

	limit := clampLimit(int(req.Limit), b.cfg.DefaultEventLimit, b.cfg.MaxEventLimit)
	offset := max(int(req.Offset), 0)

	events, err := b.geoEventRepo.GetBySubject(ctx, req.SubjectID, from, to, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get subject events: %w", err)
	}

	log.Debug("subject events query completed",
		"subject_id", req.SubjectID,
		"events", len(events),
	)

	result := make([]*models.GeoEventAPI, 0, len(events))
	for _, e := range events {
		result = append(result, e.ToAPI())
	}
	return result, nil
}

// GetAreaSubjects retrieves subjects currently inside a given area.
func (b *trackBusiness) GetAreaSubjects(
	ctx context.Context,
	req *models.GetAreaSubjectsRequest,
) ([]*models.AreaSubjectAPI, error) {
	log := util.Log(ctx)

	if req == nil || req.AreaID == "" {
		return nil, errors.New("area_id is required")
	}

	states, err := b.stateRepo.GetInsideByArea(ctx, req.AreaID, b.cfg.MaxAreaSubjects)
	if err != nil {
		return nil, fmt.Errorf("get area subjects: %w", err)
	}

	log.Debug("area subjects query completed",
		"area_id", req.AreaID,
		"subjects", len(states),
	)

	result := make([]*models.AreaSubjectAPI, 0, len(states))
	for _, s := range states {
		api := &models.AreaSubjectAPI{
			SubjectID: s.SubjectID,
		}
		if s.EnterTS != nil {
			api.EnterTimestamp = timestamppb.New(*s.EnterTS)
		}
		result = append(result, api)
	}
	return result, nil
}

// resolveTimeRange converts optional protobuf timestamps into Go time values.
// Defaults: from = 24 hours ago, to = now.
func resolveTimeRange(from, to *timestamppb.Timestamp) (time.Time, time.Time) {
	now := time.Now()
	fromTime := now.Add(-24 * time.Hour)
	toTime := now

	if from != nil {
		fromTime = from.AsTime()
	}
	if to != nil {
		toTime = to.AsTime()
	}

	return fromTime, toTime
}
