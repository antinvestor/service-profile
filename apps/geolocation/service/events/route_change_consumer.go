package events

import (
	"context"
	"errors"

	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/geolocation/service/business"
	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
)

// RouteChangeConsumer consumes RouteChangedEvent events.
// When a route is deleted, it defensively cleans up any dangling assignments
// and deviation state that may remain after partial failures.
type RouteChangeConsumer struct {
	assignmentRepo repository.RouteAssignmentRepository
	deviationRepo  repository.RouteDeviationStateRepository
}

// NewRouteChangeConsumer creates a new event consumer for route changes.
func NewRouteChangeConsumer(
	assignmentRepo repository.RouteAssignmentRepository,
	deviationRepo repository.RouteDeviationStateRepository,
) *RouteChangeConsumer {
	return &RouteChangeConsumer{
		assignmentRepo: assignmentRepo,
		deviationRepo:  deviationRepo,
	}
}

// Name returns the event name this consumer handles.
func (c *RouteChangeConsumer) Name() string {
	return business.RouteChangedEventName
}

// PayloadType returns the expected payload type for deserialization.
func (c *RouteChangeConsumer) PayloadType() any {
	return &models.RouteChangedEvent{}
}

// Validate checks that the payload is the correct type and has required fields.
func (c *RouteChangeConsumer) Validate(_ context.Context, payload any) error {
	event, ok := payload.(*models.RouteChangedEvent)
	if !ok {
		return errors.New("invalid payload type, expected *models.RouteChangedEvent")
	}
	if event.RouteID == "" {
		return errors.New("route_id is required")
	}
	if event.Action == "" {
		return errors.New("action is required")
	}
	return nil
}

// Execute processes a route changed event.
func (c *RouteChangeConsumer) Execute(ctx context.Context, payload any) error {
	event, ok := payload.(*models.RouteChangedEvent)
	if !ok {
		return errors.New("invalid payload type")
	}
	ctx = models.ContextWithEventTenancy(ctx, event.EventTenancy, event.OwnerID)

	log := util.Log(ctx)
	log.Info("route change event processed",
		"route_id", event.RouteID,
		"action", event.Action,
		"owner_id", event.OwnerID,
	)

	if event.Action != "deleted" {
		return nil
	}

	if err := c.assignmentRepo.DeleteByRoute(ctx, event.RouteID); err != nil {
		log.WithError(err).Error("failed to clean up route assignments on route deletion",
			"route_id", event.RouteID,
		)
	}
	if err := c.deviationRepo.DeleteByRoute(ctx, event.RouteID); err != nil {
		log.WithError(err).Error("failed to clean up route deviation states on route deletion",
			"route_id", event.RouteID,
		)
	}

	return nil
}
