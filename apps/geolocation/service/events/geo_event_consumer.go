package events

import (
	"context"
	"errors"

	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/geolocation/service/business"
	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

// GeoEventConsumer consumes GeoEventEmitted events.
// These are the enter/exit/dwell events emitted by the geofence engine.
// This consumer handles downstream effects like:
//   - Publishing to external queues for notification delivery.
//   - Updating analytics counters.
//   - Triggering webhook notifications.
type GeoEventConsumer struct{}

// NewGeoEventConsumer creates a new event consumer for geo events.
func NewGeoEventConsumer() *GeoEventConsumer {
	return &GeoEventConsumer{}
}

// Name returns the event name this consumer handles.
func (c *GeoEventConsumer) Name() string {
	return business.GeoEventEmittedEventName
}

// PayloadType returns the expected payload type for deserialization.
func (c *GeoEventConsumer) PayloadType() any {
	return &models.GeoEventEmitted{}
}

// Validate checks that the payload is the correct type and has required fields.
func (c *GeoEventConsumer) Validate(_ context.Context, payload any) error {
	event, ok := payload.(*models.GeoEventEmitted)
	if !ok {
		return errors.New("invalid payload type, expected *models.GeoEventEmitted")
	}
	if event.EventID == "" {
		return errors.New("event_id is required")
	}
	if event.SubjectID == "" {
		return errors.New("subject_id is required")
	}
	if event.AreaID == "" {
		return errors.New("area_id is required")
	}
	return nil
}

// Execute processes a geo event.
// Logs the event for observability. Future enhancements:
//   - Publish to external notification queue for webhook/push delivery.
//   - Update real-time analytics aggregates.
//   - Trigger area-specific automation rules.
func (c *GeoEventConsumer) Execute(ctx context.Context, payload any) error {
	event, ok := payload.(*models.GeoEventEmitted)
	if !ok {
		return errors.New("invalid payload type")
	}

	log := util.Log(ctx)
	log.Info("geo event emitted",
		"event_id", event.EventID,
		"subject_id", event.SubjectID,
		"area_id", event.AreaID,
		"event_type", event.EventType.String(),
		"confidence", event.Confidence,
	)

	return nil
}
