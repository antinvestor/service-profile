package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/geolocation/service/business"
	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	"github.com/antinvestor/service-profile/apps/geolocation/service/observability"
)

// LocationPointConsumer consumes LocationPointIngestedEvent events.
// It drives three downstream operations:
//  1. Update the materialized latest position (for proximity queries).
//  2. Run the geofence detection engine (enter/exit/dwell detection).
//
// These operations are executed sequentially: latest position must be updated
// before geofence evaluation can use it for any proximity-based lookups.
type LocationPointConsumer struct {
	proximityBiz      business.ProximityBusiness
	geofenceBiz       business.GeofenceBusiness
	routeDeviationBiz business.RouteDeviationBusiness
	metrics           *observability.Metrics
}

// NewLocationPointConsumer creates a new event consumer for ingested location points.
func NewLocationPointConsumer(
	proximityBiz business.ProximityBusiness,
	geofenceBiz business.GeofenceBusiness,
	routeDeviationBiz business.RouteDeviationBusiness,
	metrics *observability.Metrics,
) *LocationPointConsumer {
	return &LocationPointConsumer{
		proximityBiz:      proximityBiz,
		geofenceBiz:       geofenceBiz,
		routeDeviationBiz: routeDeviationBiz,
		metrics:           metrics,
	}
}

// Name returns the event name this consumer handles.
func (c *LocationPointConsumer) Name() string {
	return business.LocationPointIngestedEventName
}

// PayloadType returns the expected payload type for deserialization.
func (c *LocationPointConsumer) PayloadType() any {
	return &models.LocationPointIngestedEvent{}
}

// Validate checks that the payload is the correct type and has required fields.
func (c *LocationPointConsumer) Validate(_ context.Context, payload any) error {
	event, ok := payload.(*models.LocationPointIngestedEvent)
	if !ok {
		return errors.New("invalid payload type, expected *models.LocationPointIngestedEvent")
	}
	if event.SubjectID == "" {
		return errors.New("subject_id is required")
	}
	if event.PointID == "" {
		return errors.New("point_id is required")
	}
	return nil
}

// Execute processes an ingested location point event.
func (c *LocationPointConsumer) Execute(ctx context.Context, payload any) error {
	event, ok := payload.(*models.LocationPointIngestedEvent)
	if !ok {
		return errors.New("invalid payload type")
	}

	ctx, span := c.metrics.StartSpan(ctx, "LocationPointConsumer.Execute")
	var spanErr error
	defer func() { c.metrics.EndSpan(ctx, span, spanErr) }()

	log := util.Log(ctx)
	log.Debug("processing location point event",
		"point_id", event.PointID,
		"subject_id", event.SubjectID,
	)

	// Step 1: Update latest position for proximity queries.
	if err := c.proximityBiz.UpdateLatestPosition(ctx, event); err != nil {
		log.WithError(err).Error("failed to update latest position",
			"subject_id", event.SubjectID,
			"point_id", event.PointID,
		)
		spanErr = fmt.Errorf("update latest position: %w", err)
		return spanErr
	}

	// Step 2: Run geofence evaluation.
	if err := c.geofenceBiz.EvaluatePoint(ctx, event); err != nil {
		log.WithError(err).Error("failed to evaluate geofence",
			"subject_id", event.SubjectID,
			"point_id", event.PointID,
		)
		spanErr = fmt.Errorf("evaluate geofence: %w", err)
		return spanErr
	}

	// Step 3: Run route deviation evaluation.
	if err := c.routeDeviationBiz.EvaluatePoint(ctx, event); err != nil {
		log.WithError(err).Error("failed to evaluate route deviation",
			"subject_id", event.SubjectID,
			"point_id", event.PointID,
		)
		// Non-fatal: route deviation failure should not block the pipeline.
	}

	return nil
}
