package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/pitabwire/frame/events"
	"github.com/pitabwire/util"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
)

// RouteChangedEventName is the internal frame event name for route changes.
const RouteChangedEventName = "route.changed"

type routeBusiness struct {
	eventsMan      events.Manager
	routeRepo      repository.RouteRepository
	assignmentRepo repository.RouteAssignmentRepository
	deviationRepo  repository.RouteDeviationStateRepository
}

// NewRouteBusiness creates a new RouteBusiness.
func NewRouteBusiness(
	eventsMan events.Manager,
	routeRepo repository.RouteRepository,
	assignmentRepo repository.RouteAssignmentRepository,
	deviationRepo repository.RouteDeviationStateRepository,
) RouteBusiness {
	return &routeBusiness{
		eventsMan:      eventsMan,
		routeRepo:      routeRepo,
		assignmentRepo: assignmentRepo,
		deviationRepo:  deviationRepo,
	}
}

func (b *routeBusiness) CreateRoute(
	ctx context.Context,
	req *models.CreateRouteRequest,
) (*models.RouteAPI, error) {
	log := util.Log(ctx)

	if req == nil || req.Data == nil {
		return nil, errors.New("create route request data is nil")
	}

	apiData := req.Data

	if err := models.ValidateRouteName(apiData.Name); err != nil {
		return nil, fmt.Errorf("invalid route name: %w", err)
	}
	if err := models.ValidateRouteGeoJSON(apiData.GeometryJSON); err != nil {
		return nil, fmt.Errorf("invalid geometry: %w", err)
	}
	if apiData.OwnerID == "" {
		return nil, errors.New("owner_id is required")
	}

	route := &models.Route{
		OwnerID:                   apiData.OwnerID,
		Name:                      apiData.Name,
		Description:               apiData.Description,
		GeometryJSON:              apiData.GeometryJSON,
		State:                     StateActive,
		Extras:                    models.StructToJSONMap(apiData.Extras),
		DeviationThresholdM:       apiData.DeviationThresholdM,
		DeviationConsecutiveCount: apiData.DeviationConsecutiveCount,
		DeviationCooldownSec:      apiData.DeviationCooldownSec,
	}
	route.GenID(ctx)

	db := b.routeRepo.Pool().DB(ctx, false)
	txErr := db.Transaction(func(tx *gorm.DB) error {
		if createErr := tx.Create(route).Error; createErr != nil {
			return fmt.Errorf("create route: %w", createErr)
		}
		if geomErr := b.routeRepo.UpdateGeometryTx(
			tx, route.GetID(), apiData.GeometryJSON,
		); geomErr != nil {
			return fmt.Errorf("set route geometry: %w", geomErr)
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	persisted, err := b.routeRepo.GetByID(ctx, route.GetID())
	if err != nil {
		return nil, fmt.Errorf("read back created route: %w", err)
	}

	b.emitRouteChanged(ctx, persisted.GetID(), persisted.OwnerID, "created")

	log.Info("route created", "route_id", persisted.GetID(), "name", persisted.Name)
	return persisted.ToAPI(), nil
}

func (b *routeBusiness) UpdateRoute(
	ctx context.Context,
	req *models.UpdateRouteRequest,
) (*models.RouteAPI, error) {
	log := util.Log(ctx)

	if req == nil || req.ID == "" {
		return nil, errors.New("update route request requires an ID")
	}

	route, err := b.routeRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("route not found: %w", err)
	}

	if req.Name != "" {
		if vErr := models.ValidateRouteName(req.Name); vErr != nil {
			return nil, fmt.Errorf("invalid route name: %w", vErr)
		}
		route.Name = req.Name
	}
	if req.Description != "" {
		route.Description = req.Description
	}
	if req.DeviationThresholdM != nil {
		route.DeviationThresholdM = req.DeviationThresholdM
	}
	if req.DeviationConsecutiveCount != nil {
		route.DeviationConsecutiveCount = req.DeviationConsecutiveCount
	}
	if req.DeviationCooldownSec != nil {
		route.DeviationCooldownSec = req.DeviationCooldownSec
	}
	if req.Extras != nil {
		route.Extras = models.StructToJSONMap(req.Extras)
	}

	if _, err = b.routeRepo.Update(ctx, route); err != nil {
		return nil, fmt.Errorf("update route: %w", err)
	}

	if req.Geometry != "" {
		if vErr := models.ValidateRouteGeoJSON(req.Geometry); vErr != nil {
			return nil, fmt.Errorf("invalid geometry: %w", vErr)
		}
		if gErr := b.routeRepo.UpdateGeometry(
			ctx, route.GetID(), req.Geometry,
		); gErr != nil {
			return nil, fmt.Errorf("update route geometry: %w", gErr)
		}
	}

	persisted, err := b.routeRepo.GetByID(ctx, route.GetID())
	if err != nil {
		return nil, fmt.Errorf("read back updated route: %w", err)
	}

	b.emitRouteChanged(ctx, persisted.GetID(), persisted.OwnerID, "updated")

	log.Info("route updated", "route_id", persisted.GetID())
	return persisted.ToAPI(), nil
}

func (b *routeBusiness) DeleteRoute(ctx context.Context, routeID string) error {
	log := util.Log(ctx)

	route, err := b.routeRepo.GetByID(ctx, routeID)
	if err != nil {
		return fmt.Errorf("route not found: %w", err)
	}

	route.State = StateDeleted
	if _, err = b.routeRepo.Update(ctx, route); err != nil {
		return fmt.Errorf("soft delete route: %w", err)
	}

	// Clean up assignments.
	if delErr := b.assignmentRepo.DeleteByRoute(ctx, routeID); delErr != nil {
		log.WithError(delErr).Error("failed to clean up assignments for deleted route",
			"route_id", routeID,
		)
	}

	// Clean up deviation states.
	if delErr := b.deviationRepo.DeleteByRoute(ctx, routeID); delErr != nil {
		log.WithError(delErr).Error("failed to clean up deviation states for deleted route",
			"route_id", routeID,
		)
	}

	b.emitRouteChanged(ctx, route.GetID(), route.OwnerID, "deleted")

	log.Info("route deleted", "route_id", route.GetID())
	return nil
}

func (b *routeBusiness) GetRoute(
	ctx context.Context,
	routeID string,
) (*models.RouteAPI, error) {
	route, err := b.routeRepo.GetByID(ctx, routeID)
	if err != nil {
		return nil, fmt.Errorf("get route: %w", err)
	}
	return route.ToAPI(), nil
}

func (b *routeBusiness) SearchRoutes(
	ctx context.Context,
	ownerID string,
	limit int,
) ([]*models.RouteAPI, error) {
	if limit <= 0 {
		limit = defaultSearchLimit
	}
	if ownerID == "" {
		return nil, errors.New("owner_id is required for route search")
	}

	routes, err := b.routeRepo.SearchByOwner(ctx, ownerID, limit)
	if err != nil {
		return nil, fmt.Errorf("search routes: %w", err)
	}

	result := make([]*models.RouteAPI, 0, len(routes))
	for _, r := range routes {
		result = append(result, r.ToAPI())
	}
	return result, nil
}

func (b *routeBusiness) AssignRoute(
	ctx context.Context,
	req *models.AssignRouteRequest,
) (*models.RouteAssignmentAPI, error) {
	log := util.Log(ctx)

	if req == nil {
		return nil, errors.New("assign route request is nil")
	}
	if req.SubjectID == "" {
		return nil, errors.New("subject_id is required")
	}
	if req.RouteID == "" {
		return nil, errors.New("route_id is required")
	}

	// Verify route exists.
	if _, err := b.routeRepo.GetByID(ctx, req.RouteID); err != nil {
		return nil, fmt.Errorf("route not found: %w", err)
	}

	assignment := &models.RouteAssignment{
		SubjectID: req.SubjectID,
		RouteID:   req.RouteID,
		State:     StateActive,
	}
	if req.ValidFrom != nil {
		t := req.ValidFrom.AsTime()
		assignment.ValidFrom = &t
	}
	if req.ValidUntil != nil {
		t := req.ValidUntil.AsTime()
		assignment.ValidUntil = &t
	}
	assignment.GenID(ctx)

	if err := b.assignmentRepo.Create(ctx, assignment); err != nil {
		return nil, fmt.Errorf("create route assignment: %w", err)
	}

	log.Info("route assigned",
		"assignment_id", assignment.GetID(),
		"subject_id", req.SubjectID,
		"route_id", req.RouteID,
	)
	return assignment.ToAPI(), nil
}

func (b *routeBusiness) UnassignRoute(
	ctx context.Context,
	assignmentID string,
) error {
	log := util.Log(ctx)

	assignment, err := b.assignmentRepo.GetByID(ctx, assignmentID)
	if err != nil {
		return fmt.Errorf("assignment not found: %w", err)
	}

	assignment.State = StateDeleted
	if _, err = b.assignmentRepo.Update(ctx, assignment); err != nil {
		return fmt.Errorf("soft delete assignment: %w", err)
	}

	log.Info("route unassigned", "assignment_id", assignmentID)
	return nil
}

func (b *routeBusiness) GetSubjectAssignments(
	ctx context.Context,
	subjectID string,
) ([]*models.RouteAssignmentAPI, error) {
	assignments, err := b.assignmentRepo.GetBySubject(ctx, subjectID)
	if err != nil {
		return nil, fmt.Errorf("get subject assignments: %w", err)
	}

	result := make([]*models.RouteAssignmentAPI, 0, len(assignments))
	for _, a := range assignments {
		result = append(result, a.ToAPI())
	}
	return result, nil
}

func (b *routeBusiness) emitRouteChanged(
	ctx context.Context,
	routeID, ownerID, action string,
) {
	event := &models.RouteChangedEvent{
		RouteID: routeID,
		OwnerID: ownerID,
		Action:  action,
	}

	if err := b.eventsMan.Emit(ctx, RouteChangedEventName, event); err != nil {
		util.Log(ctx).WithError(err).Error("failed to emit route changed event",
			"route_id", routeID,
			"action", action,
		)
	}
}
