package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/security/authorizer"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/geolocation/service/authz"
	"github.com/antinvestor/service-profile/apps/geolocation/service/business"
	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	"github.com/antinvestor/service-profile/apps/geolocation/service/observability"
)

const (
	defaultHandlerSearchLimit      = 50
	defaultHandlerTrackLimit       = 100
	defaultHandlerProximityRadiusM = 1000.0
	defaultMaxBodyBytes            = 2 << 20 // 2 MiB
)

// writeAuthzError writes an authorisation error as an HTTP response.
func writeAuthzError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if errors.Is(err, authorizer.ErrInvalidSubject) || errors.Is(err, authorizer.ErrInvalidObject) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthenticated"})
		return
	}

	var permErr *authorizer.PermissionDeniedError
	if errors.As(err, &permErr) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "permission denied"})
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal authorisation error"})
}

// GeolocationServer handles HTTP API requests for the geolocation service.
type GeolocationServer struct {
	Service      *frame.Service
	authz        authz.Middleware
	ingestionBiz business.IngestionBusiness
	areaBiz      business.AreaBusiness
	routeBiz     business.RouteBusiness
	proximityBiz business.ProximityBusiness
	trackBiz     business.TrackBusiness
	metrics      *observability.Metrics
	maxBodyBytes int64
}

// NewGeolocationServer creates a new GeolocationServer with all business dependencies.
func NewGeolocationServer(
	svc *frame.Service,
	authzMiddleware authz.Middleware,
	ingestionBiz business.IngestionBusiness,
	areaBiz business.AreaBusiness,
	routeBiz business.RouteBusiness,
	proximityBiz business.ProximityBusiness,
	trackBiz business.TrackBusiness,
	metrics *observability.Metrics,
	maxBodyBytes int64,
) *GeolocationServer {
	if maxBodyBytes <= 0 {
		maxBodyBytes = defaultMaxBodyBytes
	}
	return &GeolocationServer{
		Service:      svc,
		authz:        authzMiddleware,
		ingestionBiz: ingestionBiz,
		areaBiz:      areaBiz,
		routeBiz:     routeBiz,
		proximityBiz: proximityBiz,
		trackBiz:     trackBiz,
		metrics:      metrics,
		maxBodyBytes: maxBodyBytes,
	}
}

// NewRouter registers all geolocation REST API routes.
func (s *GeolocationServer) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check (unauthenticated).
	mux.HandleFunc("GET /healthz", s.HealthCheck)

	// Ingestion.
	mux.HandleFunc("POST /v1/locations/ingest", s.IngestLocations)

	// Areas CRUD.
	mux.HandleFunc("POST /v1/areas", s.CreateArea)
	mux.HandleFunc("GET /v1/areas/{id}", s.GetArea)
	mux.HandleFunc("PUT /v1/areas/{id}", s.UpdateArea)
	mux.HandleFunc("DELETE /v1/areas/{id}", s.DeleteArea)
	mux.HandleFunc("GET /v1/areas", s.SearchAreas)

	// Routes CRUD.
	mux.HandleFunc("POST /v1/routes", s.CreateRoute)
	mux.HandleFunc("GET /v1/routes/{id}", s.GetRoute)
	mux.HandleFunc("PUT /v1/routes/{id}", s.UpdateRoute)
	mux.HandleFunc("DELETE /v1/routes/{id}", s.DeleteRoute)
	mux.HandleFunc("GET /v1/routes", s.SearchRoutes)

	// Route assignments.
	mux.HandleFunc("POST /v1/routes/assignments", s.AssignRoute)
	mux.HandleFunc("DELETE /v1/routes/assignments/{id}", s.UnassignRoute)
	mux.HandleFunc(
		"GET /v1/subjects/{subjectId}/route-assignments",
		s.GetSubjectRouteAssignments,
	)

	// Track/History.
	mux.HandleFunc("GET /v1/track/{subjectId}", s.GetTrack)
	mux.HandleFunc("GET /v1/events/subject/{subjectId}", s.GetSubjectEvents)
	mux.HandleFunc("GET /v1/events/area/{areaId}/subjects", s.GetAreaSubjects)

	// Proximity.
	mux.HandleFunc("GET /v1/proximity/subjects/{subjectId}", s.GetNearbySubjects)
	mux.HandleFunc("GET /v1/proximity/areas", s.GetNearbyAreas)

	return mux
}

// HealthCheck returns 200 if the service is healthy.
func (s *GeolocationServer) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// IngestLocations handles batch location point ingestion.
func (s *GeolocationServer) IngestLocations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "IngestLocations")
	start := time.Now()

	var req models.IngestLocationsRequest
	if err := s.decodeBody(r, &req); err != nil {
		s.metrics.EndSpan(ctx, span, err)
		s.writeClientError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Authorization: the authenticated subject must match the ingestion subject,
	// or the caller must have the ingest_location permission.
	if authErr := s.authz.CanIngestLocationSelf(ctx, req.SubjectID); authErr != nil {
		s.metrics.EndSpan(ctx, span, authErr)
		writeAuthzError(w, authErr)
		return
	}

	resp, err := s.ingestionBiz.IngestBatch(ctx, &req)
	s.metrics.EndSpan(ctx, span, err)
	if err != nil {
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.metrics.RecordIngestBatch(ctx, resp.Accepted, resp.Rejected, time.Since(start))
	s.writeJSON(w, http.StatusOK, resp)
}

// CreateArea handles area creation.
func (s *GeolocationServer) CreateArea(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "CreateArea")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	if authErr := s.authz.CanManageGeolocation(ctx); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	var req models.CreateAreaRequest
	if err := s.decodeBody(r, &req); err != nil {
		spanErr = err
		s.writeClientError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	area, err := s.areaBiz.CreateArea(ctx, &req)
	if err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.writeJSON(w, http.StatusCreated, &models.CreateAreaResponse{Data: area})
}

// GetArea retrieves an area by ID.
func (s *GeolocationServer) GetArea(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "GetArea")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	if authErr := s.authz.CanViewGeolocation(ctx); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	areaID := r.PathValue("id")
	if areaID == "" {
		s.writeClientError(w, "area id is required", http.StatusBadRequest)
		return
	}

	area, err := s.areaBiz.GetArea(ctx, areaID)
	if err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, area)
}

func (s *GeolocationServer) UpdateArea(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "UpdateArea")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	if authErr := s.authz.CanManageGeolocation(ctx); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	areaID := r.PathValue("id")
	if areaID == "" {
		s.writeClientError(w, "area id is required", http.StatusBadRequest)
		return
	}

	var req models.UpdateAreaRequest
	if err := s.decodeBody(r, &req); err != nil {
		spanErr = err
		s.writeClientError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.ID = areaID

	area, err := s.areaBiz.UpdateArea(ctx, &req)
	if err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, &models.UpdateAreaResponse{Data: area})
}

// DeleteArea handles area soft deletion.
func (s *GeolocationServer) DeleteArea(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "DeleteArea")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	if authErr := s.authz.CanManageGeolocation(ctx); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	areaID := r.PathValue("id")
	if areaID == "" {
		s.writeClientError(w, "area id is required", http.StatusBadRequest)
		return
	}

	if err := s.areaBiz.DeleteArea(ctx, areaID); err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchAreas handles area search by query text or owner ID.
func (s *GeolocationServer) SearchAreas(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "SearchAreas")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	if authErr := s.authz.CanViewGeolocation(ctx); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	query := r.URL.Query()

	q := query.Get("query")
	ownerID := query.Get("owner_id")
	limit := parseInt32Param(query.Get("limit"), defaultHandlerSearchLimit)

	areas, err := s.areaBiz.SearchAreas(ctx, q, ownerID, int(limit))
	if err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, areas)
}

// GetTrack retrieves location history for a subject.
func (s *GeolocationServer) GetTrack(w http.ResponseWriter, r *http.Request) {
	s.handleSubjectList(
		w,
		r,
		"GetTrack",
		func(ctx context.Context, subjectID string, limit, offset int32) (any, error) {
			return s.trackBiz.GetTrack(ctx, &models.GetTrackRequest{
				SubjectID: subjectID, Limit: limit, Offset: offset,
			})
		},
	)
}

// GetSubjectEvents retrieves geo events for a subject.
func (s *GeolocationServer) GetSubjectEvents(w http.ResponseWriter, r *http.Request) {
	s.handleSubjectList(
		w,
		r,
		"GetSubjectEvents",
		func(ctx context.Context, subjectID string, limit, offset int32) (any, error) {
			return s.trackBiz.GetSubjectEvents(ctx, &models.GetSubjectEventsRequest{
				SubjectID: subjectID, Limit: limit, Offset: offset,
			})
		},
	)
}

// handleSubjectList is a shared handler for subject-scoped paginated list endpoints.
func (s *GeolocationServer) handleSubjectList(
	w http.ResponseWriter, r *http.Request, spanName string,
	fetch func(ctx context.Context, subjectID string, limit, offset int32) (any, error),
) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, spanName)
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	subjectID := r.PathValue("subjectId")
	if subjectID == "" {
		s.writeClientError(w, "subject_id is required", http.StatusBadRequest)
		return
	}

	if authErr := s.authz.CanViewGeolocationSelf(ctx, subjectID); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	limit := parseInt32Param(r.URL.Query().Get("limit"), defaultHandlerTrackLimit)
	offset := parseInt32Param(r.URL.Query().Get("offset"), 0)

	result, err := fetch(ctx, subjectID, limit, offset)
	if err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, result)
}

// GetAreaSubjects retrieves subjects currently inside an area.
func (s *GeolocationServer) GetAreaSubjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "GetAreaSubjects")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	if authErr := s.authz.CanViewGeolocation(ctx); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	areaID := r.PathValue("areaId")
	if areaID == "" {
		s.writeClientError(w, "area_id is required", http.StatusBadRequest)
		return
	}

	req := &models.GetAreaSubjectsRequest{AreaID: areaID}

	subjects, err := s.trackBiz.GetAreaSubjects(ctx, req)
	if err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, subjects)
}

// GetNearbySubjects finds subjects near the given subject.
func (s *GeolocationServer) GetNearbySubjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "GetNearbySubjects")
	start := time.Now()

	subjectID := r.PathValue("subjectId")
	if subjectID == "" {
		s.metrics.EndSpan(ctx, span, nil)
		s.writeClientError(w, "subject_id is required", http.StatusBadRequest)
		return
	}

	if authErr := s.authz.CanViewGeolocationSelf(ctx, subjectID); authErr != nil {
		s.metrics.EndSpan(ctx, span, authErr)
		writeAuthzError(w, authErr)
		return
	}

	req := &models.GetNearbySubjectsRequest{
		SubjectID: subjectID,
		RadiusMeters: parseFloatParam(
			r.URL.Query().Get("radius_meters"),
			defaultHandlerProximityRadiusM,
		),
		Limit: parseInt32Param(r.URL.Query().Get("limit"), defaultHandlerSearchLimit),
	}

	subjects, err := s.proximityBiz.GetNearbySubjects(ctx, req)
	s.metrics.EndSpan(ctx, span, err)
	if err != nil {
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.metrics.RecordProximityQuery(ctx, time.Since(start), len(subjects))
	s.writeJSON(w, http.StatusOK, subjects)
}

// GetNearbyAreas finds areas near the given coordinates.
func (s *GeolocationServer) GetNearbyAreas(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "GetNearbyAreas")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	if authErr := s.authz.CanViewGeolocation(ctx); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	q := r.URL.Query()

	latStr := q.Get("latitude")
	lonStr := q.Get("longitude")
	if latStr == "" || lonStr == "" {
		s.writeClientError(
			w,
			"latitude and longitude query parameters are required",
			http.StatusBadRequest,
		)
		return
	}

	lat, latErr := strconv.ParseFloat(latStr, 64)
	lon, lonErr := strconv.ParseFloat(lonStr, 64)
	if latErr != nil || lonErr != nil {
		s.writeClientError(w, "latitude and longitude must be valid numbers", http.StatusBadRequest)
		return
	}

	req := &models.GetNearbyAreasRequest{
		Latitude:     lat,
		Longitude:    lon,
		RadiusMeters: parseFloatParam(q.Get("radius_meters"), defaultHandlerProximityRadiusM),
		Limit:        parseInt32Param(q.Get("limit"), defaultHandlerSearchLimit),
	}

	areas, err := s.proximityBiz.GetNearbyAreas(ctx, req)
	if err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, areas)
}

// decodeBody reads and decodes JSON from the request body with a size limit.
func (s *GeolocationServer) decodeBody(r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, s.maxBodyBytes)
	return json.NewDecoder(r.Body).Decode(dst)
}

// writeJSON writes a JSON response with the given status code.
func (s *GeolocationServer) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// writeClientError writes a safe, generic error message to the client.
// Internal error details are never exposed.
func (s *GeolocationServer) writeClientError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// handleBusinessError logs the full error server-side and returns a safe message to the client.
func (s *GeolocationServer) handleBusinessError(
	ctx context.Context,
	w http.ResponseWriter,
	err error,
) {
	log := util.Log(ctx)

	// Map known error patterns to appropriate HTTP status codes.
	switch {
	case errors.Is(err, errNotFound):
		s.writeClientError(w, "resource not found", http.StatusNotFound)
	case isValidationError(err):
		// Validation errors are safe to expose — they contain field names, not internal details.
		s.writeClientError(w, err.Error(), http.StatusBadRequest)
	default:
		log.WithError(err).Error("internal error processing request")
		s.writeClientError(w, "internal server error", http.StatusInternalServerError)
	}
}

// Sentinel errors for error classification.
var errNotFound = errors.New("not found")

// isValidationError checks if an error is a validation error (safe to expose to clients).
func isValidationError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	// Validation errors from the models layer contain these prefixes.
	for _, prefix := range []string{
		"invalid", "required", "must be", "exceeds maximum",
		"batch size", "either query or owner_id",
	} {
		if len(msg) >= len(prefix) && msg[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// parseInt32Param parses a string to int32, returning defaultVal on failure.
func parseInt32Param(s string, defaultVal int32) int32 {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultVal
	}
	if v > math.MaxInt32 {
		return math.MaxInt32
	}
	if v < math.MinInt32 {
		return math.MinInt32
	}
	return int32(v)
}

// parseFloatParam parses a string to float64, returning defaultVal on failure.
func parseFloatParam(s string, defaultVal float64) float64 {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultVal
	}
	return v
}

// --- Route handlers ---

// CreateRoute handles route creation.
func (s *GeolocationServer) CreateRoute(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "CreateRoute")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	if authErr := s.authz.CanManageGeolocation(ctx); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	var req models.CreateRouteRequest
	if err := s.decodeBody(r, &req); err != nil {
		spanErr = err
		s.writeClientError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	route, err := s.routeBiz.CreateRoute(ctx, &req)
	if err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.writeJSON(w, http.StatusCreated, &models.CreateRouteResponse{Data: route})
}

// GetRoute retrieves a route by ID.
func (s *GeolocationServer) GetRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "GetRoute")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	if authErr := s.authz.CanViewGeolocation(ctx); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	routeID := r.PathValue("id")
	if routeID == "" {
		s.writeClientError(w, "route id is required", http.StatusBadRequest)
		return
	}

	route, err := s.routeBiz.GetRoute(ctx, routeID)
	if err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, route)
}

func (s *GeolocationServer) UpdateRoute(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "UpdateRoute")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	if authErr := s.authz.CanManageGeolocation(ctx); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	routeID := r.PathValue("id")
	if routeID == "" {
		s.writeClientError(w, "route id is required", http.StatusBadRequest)
		return
	}

	var req models.UpdateRouteRequest
	if err := s.decodeBody(r, &req); err != nil {
		spanErr = err
		s.writeClientError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.ID = routeID

	route, err := s.routeBiz.UpdateRoute(ctx, &req)
	if err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, &models.UpdateRouteResponse{Data: route})
}

// DeleteRoute handles route soft deletion.
func (s *GeolocationServer) DeleteRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "DeleteRoute")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	if authErr := s.authz.CanManageGeolocation(ctx); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	routeID := r.PathValue("id")
	if routeID == "" {
		s.writeClientError(w, "route id is required", http.StatusBadRequest)
		return
	}

	if err := s.routeBiz.DeleteRoute(ctx, routeID); err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchRoutes handles route search by owner ID.
func (s *GeolocationServer) SearchRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "SearchRoutes")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	if authErr := s.authz.CanViewGeolocation(ctx); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	ownerID := r.URL.Query().Get("owner_id")
	limit := parseInt32Param(r.URL.Query().Get("limit"), defaultHandlerSearchLimit)

	routes, err := s.routeBiz.SearchRoutes(ctx, ownerID, int(limit))
	if err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, routes)
}

// AssignRoute handles assigning a subject to a route.
func (s *GeolocationServer) AssignRoute(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "AssignRoute")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	if authErr := s.authz.CanManageGeolocation(ctx); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	var req models.AssignRouteRequest
	if err := s.decodeBody(r, &req); err != nil {
		spanErr = err
		s.writeClientError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	assignment, err := s.routeBiz.AssignRoute(ctx, &req)
	if err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.writeJSON(w, http.StatusCreated, &models.AssignRouteResponse{Data: assignment})
}

// UnassignRoute handles removing a route assignment.
func (s *GeolocationServer) UnassignRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "UnassignRoute")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	if authErr := s.authz.CanManageGeolocation(ctx); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	assignmentID := r.PathValue("id")
	if assignmentID == "" {
		s.writeClientError(w, "assignment id is required", http.StatusBadRequest)
		return
	}

	if err := s.routeBiz.UnassignRoute(ctx, assignmentID); err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetSubjectRouteAssignments retrieves a subject's active route assignments.
func (s *GeolocationServer) GetSubjectRouteAssignments(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	ctx, span := s.metrics.StartSpan(ctx, "GetSubjectRouteAssignments")
	var spanErr error
	defer func() { s.metrics.EndSpan(ctx, span, spanErr) }()

	subjectID := r.PathValue("subjectId")
	if subjectID == "" {
		s.writeClientError(w, "subject_id is required", http.StatusBadRequest)
		return
	}

	if authErr := s.authz.CanViewGeolocationSelf(ctx, subjectID); authErr != nil {
		spanErr = authErr
		writeAuthzError(w, authErr)
		return
	}

	assignments, err := s.routeBiz.GetSubjectAssignments(ctx, subjectID)
	if err != nil {
		spanErr = err
		s.handleBusinessError(ctx, w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, assignments)
}
