package business

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/pitabwire/frame/events"
	"github.com/pitabwire/util"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	"github.com/antinvestor/service-profile/apps/geolocation/service/observability"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
)

// GeoEventEmittedEventName is the internal frame event name for emitted geo events.
const GeoEventEmittedEventName = "geo.event.emitted"

type geofenceBusiness struct {
	eventsMan    events.Manager
	areaRepo     repository.AreaRepository
	stateRepo    repository.GeofenceStateRepository
	geoEventRepo repository.GeoEventRepository
	metrics      *observability.Metrics

	// Configurable parameters (injected from GeolocationConfig).
	hysteresisBufferM  float64
	dwellThreshold     time.Duration
	maxCandidateAreas  int
	maxAccuracyForEval float64
}

// GeofenceConfig holds tunable parameters for the geofence engine.
type GeofenceConfig struct {
	HysteresisBufferM  float64
	DwellThreshold     time.Duration
	MaxCandidateAreas  int
	MaxAccuracyForEval float64
}

// Geofence constants used when config values are zero (safety net).
const (
	defaultHysteresisBufferM  = 30.0
	defaultDwellThreshold     = 2 * time.Minute
	defaultMaxCandidateAreas  = 100
	defaultMaxAccuracyForEval = 500.0

	// confidenceDecayConstant is the decay constant used in the exponential
	// confidence calculation: confidence = e^(-accuracy / confidenceDecayConstant).
	confidenceDecayConstant = 100.0

	// minConfidence is the minimum confidence value returned by computeConfidence.
	minConfidence = 0.01
)

// NewGeofenceBusiness creates a new GeofenceBusiness with configurable parameters.
func NewGeofenceBusiness(
	eventsMan events.Manager,
	areaRepo repository.AreaRepository,
	stateRepo repository.GeofenceStateRepository,
	geoEventRepo repository.GeoEventRepository,
	metrics *observability.Metrics,
	cfg GeofenceConfig,
) GeofenceBusiness {
	hysteresis := cfg.HysteresisBufferM
	if hysteresis <= 0 {
		hysteresis = defaultHysteresisBufferM
	}
	dwell := cfg.DwellThreshold
	if dwell <= 0 {
		dwell = defaultDwellThreshold
	}
	maxCandidates := cfg.MaxCandidateAreas
	if maxCandidates <= 0 {
		maxCandidates = defaultMaxCandidateAreas
	}
	maxAccuracy := cfg.MaxAccuracyForEval
	if maxAccuracy <= 0 {
		maxAccuracy = defaultMaxAccuracyForEval
	}

	return &geofenceBusiness{
		eventsMan:          eventsMan,
		areaRepo:           areaRepo,
		stateRepo:          stateRepo,
		geoEventRepo:       geoEventRepo,
		metrics:            metrics,
		hysteresisBufferM:  hysteresis,
		dwellThreshold:     dwell,
		maxCandidateAreas:  maxCandidates,
		maxAccuracyForEval: maxAccuracy,
	}
}

// EvaluatePoint runs the geofence state machine for a single ingested location point.
// Processing flow:
//  1. Validate the point's accuracy is sufficient for spatial decisions.
//  2. Bbox pre-filter: find candidate areas whose bounding box contains the point.
//  3. For each candidate area, run the containment test and state machine within a transaction.
//  4. Emit enter/exit/dwell events as transitions occur.
func (b *geofenceBusiness) EvaluatePoint(ctx context.Context, event *models.LocationPointIngestedEvent) error {
	log := util.Log(ctx)
	start := time.Now()
	ctx, span := b.metrics.StartSpan(ctx, "GeofenceEvaluatePoint")
	var spanErr error
	defer func() { b.metrics.EndSpan(ctx, span, spanErr) }()

	if event == nil {
		spanErr = errors.New("event is nil")
		return spanErr
	}

	// Skip points with poor accuracy — they would cause false transitions.
	if event.Accuracy > b.maxAccuracyForEval {
		log.Debug("skipping geofence evaluation due to poor accuracy",
			"subject_id", event.SubjectID,
			"accuracy", event.Accuracy,
		)
		return nil
	}

	// Skip points with zero accuracy — this is a sensor lie (no real GPS returns exactly 0).
	if event.Accuracy <= 0 {
		log.Debug("skipping geofence evaluation due to zero/negative accuracy",
			"subject_id", event.SubjectID,
			"accuracy", event.Accuracy,
		)
		return nil
	}

	pointTS := time.UnixMilli(event.Timestamp)

	// Step 1: Get candidate areas whose bounding box intersects this point.
	candidates, err := b.areaRepo.GetActiveByBoundingBox(ctx, event.Latitude, event.Longitude)
	if err != nil {
		return fmt.Errorf("get candidate areas: %w", err)
	}

	if len(candidates) > b.maxCandidateAreas {
		log.Warn("bbox pre-filter returned too many candidates, truncating",
			"subject_id", event.SubjectID,
			"candidates", len(candidates),
			"max", b.maxCandidateAreas,
		)
		candidates = candidates[:b.maxCandidateAreas]
	}

	// Step 2: For each candidate, test containment and run state machine.
	for _, area := range candidates {
		if evalErr := b.evaluateAreaForSubject(ctx, event, area, pointTS); evalErr != nil {
			// Log and continue — one area failure should not block evaluation of others.
			log.WithError(evalErr).Error("geofence evaluation failed for area",
				"subject_id", event.SubjectID,
				"area_id", area.GetID(),
			)
		}
	}

	b.metrics.RecordGeofenceEval(ctx, time.Since(start))
	return nil
}

// evaluateAreaForSubject runs the state machine for a single (subject, area) pair.
// The entire read-evaluate-write cycle runs within a database transaction with
// SELECT FOR UPDATE on the geofence state row to prevent concurrent modification.
func (b *geofenceBusiness) evaluateAreaForSubject(
	ctx context.Context,
	event *models.LocationPointIngestedEvent,
	area *models.Area,
	pointTS time.Time,
) error {
	log := util.Log(ctx)

	// Test actual geometry containment (not just bounding box).
	isInside, err := b.areaRepo.ContainsPoint(ctx, area.GetID(), event.Latitude, event.Longitude)
	if err != nil {
		return fmt.Errorf("containment test for area %s: %w", area.GetID(), err)
	}

	// Run the state machine within a transaction with row-level locking.
	db := b.stateRepo.Pool().DB(ctx, false)

	return db.Transaction(func(tx *gorm.DB) error {
		// Retrieve current state with row-level lock.
		// Returns (nil, nil) if no state exists yet.
		state, getErr := b.stateRepo.GetForUpdate(tx, event.SubjectID, area.GetID())
		if getErr != nil {
			return fmt.Errorf("get state for update: %w", getErr)
		}

		// No existing state — initialize.
		if state == nil {
			state = &models.GeofenceState{
				SubjectID: event.SubjectID,
				AreaID:    area.GetID(),
				Inside:    false,
			}
		}

		// Timestamp ordering guard: skip out-of-order points.
		// This prevents older events from flipping state set by newer events.
		if state.LastPointTS != nil && pointTS.Before(*state.LastPointTS) {
			log.Debug("skipping out-of-order point for geofence evaluation",
				"subject_id", event.SubjectID,
				"area_id", area.GetID(),
				"point_ts", pointTS,
				"last_point_ts", *state.LastPointTS,
			)
			return nil
		}

		// Compute confidence based on accuracy.
		confidence := computeConfidence(event.Accuracy)

		// Apply hysteresis: maintain current state when near boundary.
		effectiveInside := applyHysteresis(isInside, state.Inside, event.Accuracy, b.hysteresisBufferM)

		// Determine state transition.
		switch {
		case !state.Inside && effectiveInside:
			return b.handleEnter(ctx, tx, event, area, state, pointTS, confidence)

		case state.Inside && !effectiveInside:
			return b.handleExit(ctx, tx, event, area, state, pointTS, confidence)

		case state.Inside && effectiveInside:
			if dwellErr := b.checkDwell(ctx, tx, event, area, state, pointTS, confidence); dwellErr != nil {
				log.WithError(dwellErr).Error("dwell check failed",
					"subject_id", event.SubjectID,
					"area_id", area.GetID(),
				)
			}
			return b.updateStatePosition(tx, state, event, pointTS)

		default:
			return b.updateStatePosition(tx, state, event, pointTS)
		}
	})
}

// handleEnter processes an enter transition within a transaction.
func (b *geofenceBusiness) handleEnter(
	ctx context.Context,
	tx *gorm.DB,
	event *models.LocationPointIngestedEvent,
	area *models.Area,
	state *models.GeofenceState,
	pointTS time.Time,
	confidence float64,
) error {
	log := util.Log(ctx)

	geoEvent := &models.GeoEvent{
		SubjectID:  event.SubjectID,
		AreaID:     area.GetID(),
		EventType:  models.GeoEventTypeEnter,
		TS:         pointTS,
		Confidence: confidence,
		PointID:    event.PointID,
	}
	geoEvent.GenID(ctx)

	if err := b.geoEventRepo.CreateTx(tx, geoEvent); err != nil {
		return fmt.Errorf("persist enter event: %w", err)
	}

	state.Inside = true
	state.LastTransition = &pointTS
	state.EnterTS = &pointTS
	state.LastPointTS = &pointTS
	state.LastLat = event.Latitude
	state.LastLon = event.Longitude

	if err := b.stateRepo.UpsertTx(tx, state); err != nil {
		return fmt.Errorf("upsert enter state: %w", err)
	}

	// Emit event for downstream consumers (outside transaction — non-fatal).
	b.emitGeoEvent(ctx, geoEvent)

	log.Info("geofence ENTER",
		"subject_id", event.SubjectID,
		"area_id", area.GetID(),
		"area_name", area.Name,
		"confidence", confidence,
	)

	return nil
}

// handleExit processes an exit transition within a transaction.
func (b *geofenceBusiness) handleExit(
	ctx context.Context,
	tx *gorm.DB,
	event *models.LocationPointIngestedEvent,
	area *models.Area,
	state *models.GeofenceState,
	pointTS time.Time,
	confidence float64,
) error {
	log := util.Log(ctx)

	geoEvent := &models.GeoEvent{
		SubjectID:  event.SubjectID,
		AreaID:     area.GetID(),
		EventType:  models.GeoEventTypeExit,
		TS:         pointTS,
		Confidence: confidence,
		PointID:    event.PointID,
	}
	geoEvent.GenID(ctx)

	if err := b.geoEventRepo.CreateTx(tx, geoEvent); err != nil {
		return fmt.Errorf("persist exit event: %w", err)
	}

	state.Inside = false
	state.LastTransition = &pointTS
	state.EnterTS = nil
	state.LastPointTS = &pointTS
	state.LastLat = event.Latitude
	state.LastLon = event.Longitude

	if err := b.stateRepo.UpsertTx(tx, state); err != nil {
		return fmt.Errorf("upsert exit state: %w", err)
	}

	b.emitGeoEvent(ctx, geoEvent)

	log.Info("geofence EXIT",
		"subject_id", event.SubjectID,
		"area_id", area.GetID(),
		"area_name", area.Name,
		"confidence", confidence,
	)

	return nil
}

// checkDwell checks if a subject has been inside an area long enough to emit a DWELL event.
// Uses HasDwellEvent for an efficient targeted check instead of scanning all events.
func (b *geofenceBusiness) checkDwell(
	ctx context.Context,
	tx *gorm.DB,
	event *models.LocationPointIngestedEvent,
	area *models.Area,
	state *models.GeofenceState,
	pointTS time.Time,
	confidence float64,
) error {
	if state.EnterTS == nil {
		return nil
	}

	dwellDuration := pointTS.Sub(*state.EnterTS)
	if dwellDuration < b.dwellThreshold {
		return nil
	}

	// Check for existing dwell event within the transaction to prevent duplicates.
	// Using the tx ensures serialization with the SELECT FOR UPDATE lock on geofence_state.
	hasDwell, err := b.geoEventRepo.HasDwellEventTx(tx, event.SubjectID, area.GetID(), state.EnterTS)
	if err != nil {
		return fmt.Errorf("check existing dwell events: %w", err)
	}
	if hasDwell {
		return nil
	}

	log := util.Log(ctx)

	geoEvent := &models.GeoEvent{
		SubjectID:  event.SubjectID,
		AreaID:     area.GetID(),
		EventType:  models.GeoEventTypeDwell,
		TS:         pointTS,
		Confidence: confidence,
		PointID:    event.PointID,
	}
	geoEvent.GenID(ctx)

	if createErr := b.geoEventRepo.CreateTx(tx, geoEvent); createErr != nil {
		return fmt.Errorf("persist dwell event: %w", createErr)
	}

	b.emitGeoEvent(ctx, geoEvent)

	log.Info("geofence DWELL",
		"subject_id", event.SubjectID,
		"area_id", area.GetID(),
		"area_name", area.Name,
		"dwell_duration", dwellDuration.String(),
	)

	return nil
}

// updateStatePosition updates the geofence state with the latest point position
// without changing the inside/outside status. Runs within a transaction.
func (b *geofenceBusiness) updateStatePosition(
	tx *gorm.DB,
	state *models.GeofenceState,
	event *models.LocationPointIngestedEvent,
	pointTS time.Time,
) error {
	state.LastPointTS = &pointTS
	state.LastLat = event.Latitude
	state.LastLon = event.Longitude

	if err := b.stateRepo.UpsertTx(tx, state); err != nil {
		return fmt.Errorf("update state position: %w", err)
	}
	return nil
}

// emitGeoEvent emits a GeoEventEmitted internal event for downstream consumers.
func (b *geofenceBusiness) emitGeoEvent(ctx context.Context, geoEvent *models.GeoEvent) {
	b.metrics.RecordGeofenceTransition(ctx, geoEvent.EventType.String())

	emitted := &models.GeoEventEmitted{
		EventID:    geoEvent.GetID(),
		SubjectID:  geoEvent.SubjectID,
		AreaID:     geoEvent.AreaID,
		EventType:  geoEvent.EventType,
		Timestamp:  geoEvent.TS.UnixMilli(),
		Confidence: geoEvent.Confidence,
	}

	if err := b.eventsMan.Emit(ctx, GeoEventEmittedEventName, emitted); err != nil {
		util.Log(ctx).WithError(err).Error("failed to emit geo event",
			"event_id", geoEvent.GetID(),
			"event_type", geoEvent.EventType.String(),
		)
	}
}

// computeConfidence maps GPS accuracy (in meters) to a confidence score [0, 1].
// Lower accuracy value (more precise) = higher confidence.
// Uses an exponential decay curve centered around typical GPS accuracy ranges.
func computeConfidence(accuracyMeters float64) float64 {
	if accuracyMeters <= 0 {
		return 1.0
	}
	confidence := math.Exp(-accuracyMeters / confidenceDecayConstant)
	if confidence < minConfidence {
		return minConfidence
	}
	return confidence
}

// applyHysteresis determines the effective inside/outside state after accounting
// for the hysteresis buffer. When the subject is near a boundary (accuracy > buffer),
// we maintain the current state to prevent oscillation.
func applyHysteresis(rawInside, currentInside bool, accuracyMeters, hysteresisBufferM float64) bool {
	if rawInside == currentInside {
		return rawInside
	}

	// State transition requested. Only allow if accuracy is better than the buffer.
	if accuracyMeters <= hysteresisBufferM {
		return rawInside
	}

	// Accuracy too poor to confirm transition — maintain current state.
	return currentInside
}
