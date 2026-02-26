package business

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pitabwire/frame/events"
	"github.com/pitabwire/util"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	"github.com/antinvestor/service-profile/apps/geolocation/service/observability"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
)

// RouteDeviationDetectedEventName is the event name for route deviation events.
const RouteDeviationDetectedEventName = "route.deviation.detected"

// RouteDeviationConfig holds the global GPS accuracy filter.
// Per-route thresholds come from the Route model.
type RouteDeviationConfig struct {
	MaxAccuracyForEval float64
}

const (
	defaultRouteDeviationMaxAccuracy = 500.0
	defaultConsecutiveCount          = 3
	defaultCooldownSec               = 300
)

type routeDeviationBusiness struct {
	eventsMan     events.Manager
	routeRepo     repository.RouteRepository
	deviationRepo repository.RouteDeviationStateRepository
	metrics       *observability.Metrics
	maxAccuracyM  float64
}

// NewRouteDeviationBusiness creates a new RouteDeviationBusiness.
func NewRouteDeviationBusiness(
	eventsMan events.Manager,
	routeRepo repository.RouteRepository,
	deviationRepo repository.RouteDeviationStateRepository,
	metrics *observability.Metrics,
	cfg RouteDeviationConfig,
) RouteDeviationBusiness {
	maxAcc := cfg.MaxAccuracyForEval
	if maxAcc <= 0 {
		maxAcc = defaultRouteDeviationMaxAccuracy
	}
	return &routeDeviationBusiness{
		eventsMan:     eventsMan,
		routeRepo:     routeRepo,
		deviationRepo: deviationRepo,
		metrics:       metrics,
		maxAccuracyM:  maxAcc,
	}
}

// EvaluatePoint runs the route deviation state machine for a single ingested point.
func (b *routeDeviationBusiness) EvaluatePoint(
	ctx context.Context,
	event *models.LocationPointIngestedEvent,
) error {
	log := util.Log(ctx)
	start := time.Now()
	ctx, span := b.metrics.StartSpan(ctx, "RouteDeviationEvaluatePoint")
	var spanErr error
	defer func() { b.metrics.EndSpan(ctx, span, spanErr) }()

	if event == nil {
		spanErr = errors.New("event is nil")
		return spanErr
	}

	if event.Accuracy > b.maxAccuracyM || event.Accuracy <= 0 {
		log.Debug("skipping route deviation evaluation due to accuracy",
			"subject_id", event.SubjectID,
			"accuracy", event.Accuracy,
		)
		return nil
	}

	pointTS := time.UnixMilli(event.Timestamp)

	assignments, err := b.routeRepo.GetActiveAssignmentsForSubject(
		ctx, event.SubjectID, pointTS,
	)
	if err != nil {
		return fmt.Errorf("get active route assignments: %w", err)
	}

	for _, rwa := range assignments {
		if evalErr := b.evaluateRouteForSubject(
			ctx, event, rwa.Route, pointTS,
		); evalErr != nil {
			log.WithError(evalErr).Error("route deviation evaluation failed",
				"subject_id", event.SubjectID,
				"route_id", rwa.Route.GetID(),
			)
		}
	}

	b.metrics.RecordRouteDeviationEval(ctx, time.Since(start))
	return nil
}

// routeThresholds holds resolved per-route deviation parameters.
type routeThresholds struct {
	thresholdM       float64
	consecutiveCount int
	cooldownSec      int
}

func resolveThresholds(route *models.Route) routeThresholds {
	t := routeThresholds{
		thresholdM:       *route.DeviationThresholdM,
		consecutiveCount: defaultConsecutiveCount,
		cooldownSec:      defaultCooldownSec,
	}
	if route.DeviationConsecutiveCount != nil {
		t.consecutiveCount = *route.DeviationConsecutiveCount
	}
	if route.DeviationCooldownSec != nil {
		t.cooldownSec = *route.DeviationCooldownSec
	}
	return t
}

func (b *routeDeviationBusiness) evaluateRouteForSubject(
	ctx context.Context,
	event *models.LocationPointIngestedEvent,
	route *models.Route,
	pointTS time.Time,
) error {
	if !route.HasDeviationConfig() {
		return nil
	}

	thresholds := resolveThresholds(route)

	distance, err := b.routeRepo.DistanceToRouteMeters(
		ctx, route.GetID(), event.Latitude, event.Longitude,
	)
	if err != nil {
		return fmt.Errorf("distance to route %s: %w", route.GetID(), err)
	}

	offRoute := distance > thresholds.thresholdM
	db := b.deviationRepo.Pool().DB(ctx, false)

	return db.Transaction(func(tx *gorm.DB) error {
		state, getErr := b.deviationRepo.GetForUpdate(
			tx, event.SubjectID, route.GetID(),
		)
		if getErr != nil {
			return fmt.Errorf("get deviation state for update: %w", getErr)
		}

		if state == nil {
			state = &models.RouteDeviationState{
				SubjectID: event.SubjectID,
				RouteID:   route.GetID(),
			}
		}

		if state.LastPointTS != nil && pointTS.Before(*state.LastPointTS) {
			return nil
		}

		return b.processDeviationState(
			ctx, tx, event, route, state, offRoute, distance, pointTS, thresholds,
		)
	})
}

func (b *routeDeviationBusiness) processDeviationState(
	ctx context.Context,
	tx *gorm.DB,
	event *models.LocationPointIngestedEvent,
	route *models.Route,
	state *models.RouteDeviationState,
	offRoute bool,
	distance float64,
	pointTS time.Time,
	t routeThresholds,
) error {
	updatePosition(state, event, pointTS)

	switch {
	case offRoute:
		return b.handleOffRoute(ctx, tx, event, route, state, distance, pointTS, t)
	case state.Deviated:
		return b.handleBackOnRoute(ctx, tx, event, route, state, distance, pointTS)
	default:
		state.ConsecutiveOffRoute = 0
		return b.deviationRepo.UpsertTx(tx, state)
	}
}

func (b *routeDeviationBusiness) handleOffRoute(
	ctx context.Context,
	tx *gorm.DB,
	event *models.LocationPointIngestedEvent,
	route *models.Route,
	state *models.RouteDeviationState,
	distance float64,
	pointTS time.Time,
	t routeThresholds,
) error {
	state.ConsecutiveOffRoute++

	if !b.shouldTriggerDeviation(state, pointTS, t) {
		return b.deviationRepo.UpsertTx(tx, state)
	}

	state.Deviated = true
	state.LastDeviationEventAt = &pointTS

	if uErr := b.deviationRepo.UpsertTx(tx, state); uErr != nil {
		return fmt.Errorf("upsert deviated state: %w", uErr)
	}

	b.emitDeviationEvent(ctx, event, route, "deviated", distance, pointTS)

	util.Log(ctx).Info("route DEVIATED",
		"subject_id", event.SubjectID,
		"route_id", route.GetID(),
		"route_name", route.Name,
		"distance_m", distance,
	)
	return nil
}

func (b *routeDeviationBusiness) shouldTriggerDeviation(
	state *models.RouteDeviationState,
	pointTS time.Time,
	t routeThresholds,
) bool {
	if state.ConsecutiveOffRoute < t.consecutiveCount {
		return false
	}
	if !state.Deviated {
		return true
	}
	// Re-trigger check: already deviated, check cooldown.
	if state.LastDeviationEventAt == nil {
		return true
	}
	elapsed := pointTS.Sub(*state.LastDeviationEventAt)
	return elapsed >= time.Duration(t.cooldownSec)*time.Second
}

func (b *routeDeviationBusiness) handleBackOnRoute(
	ctx context.Context,
	tx *gorm.DB,
	event *models.LocationPointIngestedEvent,
	route *models.Route,
	state *models.RouteDeviationState,
	distance float64,
	pointTS time.Time,
) error {
	state.Deviated = false
	state.ConsecutiveOffRoute = 0

	if uErr := b.deviationRepo.UpsertTx(tx, state); uErr != nil {
		return fmt.Errorf("upsert back-on-route state: %w", uErr)
	}

	b.emitDeviationEvent(ctx, event, route, "back_on_route", distance, pointTS)

	util.Log(ctx).Info("route BACK_ON_ROUTE",
		"subject_id", event.SubjectID,
		"route_id", route.GetID(),
		"route_name", route.Name,
	)
	return nil
}

func updatePosition(
	state *models.RouteDeviationState,
	event *models.LocationPointIngestedEvent,
	pointTS time.Time,
) {
	state.LastPointTS = &pointTS
	state.LastLat = event.Latitude
	state.LastLon = event.Longitude
}

func (b *routeDeviationBusiness) emitDeviationEvent(
	ctx context.Context,
	event *models.LocationPointIngestedEvent,
	route *models.Route,
	eventType string,
	distance float64,
	pointTS time.Time,
) {
	b.metrics.RecordRouteDeviationTransition(ctx, eventType)

	payload := &models.RouteDeviationDetectedEvent{
		SubjectID:      event.SubjectID,
		RouteID:        route.GetID(),
		EventType:      eventType,
		DistanceMeters: distance,
		Latitude:       event.Latitude,
		Longitude:      event.Longitude,
		Timestamp:      pointTS.UnixMilli(),
	}

	if err := b.eventsMan.Emit(
		ctx, RouteDeviationDetectedEventName, payload,
	); err != nil {
		util.Log(ctx).WithError(err).Error("failed to emit route deviation event",
			"subject_id", event.SubjectID,
			"route_id", route.GetID(),
			"event_type", eventType,
		)
	}
}
