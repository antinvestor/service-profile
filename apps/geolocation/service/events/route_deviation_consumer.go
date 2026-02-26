package events

import (
	"context"
	"errors"

	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/geolocation/service/business"
	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

// RouteDeviationConsumer consumes RouteDeviationDetectedEvent events.
// These are the deviated/back_on_route events emitted by the route deviation engine.
// This consumer handles downstream effects like:
//   - Publishing to external queues for notification delivery.
//   - Triggering webhook notifications.
type RouteDeviationConsumer struct{}

// NewRouteDeviationConsumer creates a new event consumer for route deviation events.
func NewRouteDeviationConsumer() *RouteDeviationConsumer {
	return &RouteDeviationConsumer{}
}

// Name returns the event name this consumer handles.
func (c *RouteDeviationConsumer) Name() string {
	return business.RouteDeviationDetectedEventName
}

// PayloadType returns the expected payload type for deserialization.
func (c *RouteDeviationConsumer) PayloadType() any {
	return &models.RouteDeviationDetectedEvent{}
}

// Validate checks that the payload is the correct type and has required fields.
func (c *RouteDeviationConsumer) Validate(_ context.Context, payload any) error {
	event, ok := payload.(*models.RouteDeviationDetectedEvent)
	if !ok {
		return errors.New(
			"invalid payload type, expected *models.RouteDeviationDetectedEvent",
		)
	}
	if event.SubjectID == "" {
		return errors.New("subject_id is required")
	}
	if event.RouteID == "" {
		return errors.New("route_id is required")
	}
	if event.EventType == "" {
		return errors.New("event_type is required")
	}
	return nil
}

// Execute processes a route deviation event.
// Logs the event for observability. Future enhancements:
//   - Publish to external notification queue for webhook/push delivery.
//   - Update real-time analytics aggregates.
func (c *RouteDeviationConsumer) Execute(ctx context.Context, payload any) error {
	event, ok := payload.(*models.RouteDeviationDetectedEvent)
	if !ok {
		return errors.New("invalid payload type")
	}

	log := util.Log(ctx)
	log.Info("route deviation event emitted",
		"subject_id", event.SubjectID,
		"route_id", event.RouteID,
		"event_type", event.EventType,
		"distance_meters", event.DistanceMeters,
	)

	return nil
}
