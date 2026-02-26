package events

import (
	"context"
	"errors"

	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/geolocation/service/business"
	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
)

// AreaChangeConsumer consumes AreaChangedEvent events.
// When an area is deleted, it cleans up related geofence states.
type AreaChangeConsumer struct {
	areaBiz   business.AreaBusiness
	stateRepo repository.GeofenceStateRepository
}

// NewAreaChangeConsumer creates a new event consumer for area changes.
func NewAreaChangeConsumer(
	areaBiz business.AreaBusiness,
	stateRepo repository.GeofenceStateRepository,
) *AreaChangeConsumer {
	return &AreaChangeConsumer{
		areaBiz:   areaBiz,
		stateRepo: stateRepo,
	}
}

// Name returns the event name this consumer handles.
func (c *AreaChangeConsumer) Name() string {
	return business.AreaChangedEventName
}

// PayloadType returns the expected payload type for deserialization.
func (c *AreaChangeConsumer) PayloadType() any {
	return &models.AreaChangedEvent{}
}

// Validate checks that the payload is the correct type and has required fields.
func (c *AreaChangeConsumer) Validate(_ context.Context, payload any) error {
	event, ok := payload.(*models.AreaChangedEvent)
	if !ok {
		return errors.New("invalid payload type, expected *models.AreaChangedEvent")
	}
	if event.AreaID == "" {
		return errors.New("area_id is required")
	}
	if event.Action == "" {
		return errors.New("action is required")
	}
	return nil
}

// Execute processes an area changed event.
// On "deleted": clean up geofence_states for the deleted area (defense in depth —
// the business layer also cleans up, but this handles edge cases like partial failures).
func (c *AreaChangeConsumer) Execute(ctx context.Context, payload any) error {
	event, ok := payload.(*models.AreaChangedEvent)
	if !ok {
		return errors.New("invalid payload type")
	}

	log := util.Log(ctx)
	log.Info("area change event processed",
		"area_id", event.AreaID,
		"action", event.Action,
		"owner_id", event.OwnerID,
	)

	// Clean up geofence states when an area is deleted.
	if event.Action == "deleted" {
		if err := c.stateRepo.DeleteByArea(ctx, event.AreaID); err != nil {
			log.WithError(err).Error("failed to clean up geofence states on area deletion",
				"area_id", event.AreaID,
			)
			// Non-fatal: the area is already deleted, new geofence evaluations will skip it.
		}
	}

	return nil
}
