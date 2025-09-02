package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/models"
)

func (ps *ProfileServer) writeError(ctx context.Context, w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")

	log := ps.Service.Log(ctx).
		WithField("code", code)

	log.WithError(err).Error("internal service error")
	w.WriteHeader(code)

	err = json.NewEncoder(w).Encode(fmt.Sprintf(" internal processing err message: %v", err))
	if err != nil {
		log.WithError(err).Error("could not write error to response")
	}
}

// RestListRelationshipsEndpoint handles listing relationships via REST API.
func (ps *ProfileServer) RestListRelationshipsEndpoint(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := ps.Service.Log(ctx).WithField("method", "RestListRelationshipsEndpoint")

	// Validate claims
	claims := frame.ClaimsFromContext(ctx)
	if claims == nil {
		ps.writeError(ctx, rw, errors.New("claims can not be empty"), http.StatusInternalServerError)
		return
	}

	// Parse request parameters
	urlQuery := req.URL.Query()
	log.WithField("url params", urlQuery).Debug("listing relationships request")

	// Extract parameters and build request
	request := ps.buildRelationshipListRequest(urlQuery, claims)

	// Fetch relationships
	relationshipObjectList, lastRelID, err := ps.fetchRelationships(ctx, request)
	if err != nil {
		ps.writeError(ctx, rw, err, http.StatusInternalServerError)
		return
	}

	// Prepare and send response
	ps.sendRelationshipListResponse(ctx, rw, relationshipObjectList, lastRelID, claims)
}

// buildRelationshipListRequest extracts URL parameters and constructs a ListRelationshipRequest.
func (ps *ProfileServer) buildRelationshipListRequest(
	urlQuery url.Values,
	claims *frame.AuthenticationClaims,
) *profilev1.ListRelationshipRequest {
	// Parse count parameter
	countStr := urlQuery.Get("Count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		count = 20 // Default count
	}

	// Get last relationship ID for pagination
	lastRelationshipID := ""
	if urlQuery.Has("LastRelationshipID") {
		lastRelationshipID = urlQuery.Get("LastRelationshipID")
	}

	// Determine peer object
	peerObject := "Profile"
	if urlQuery.Has("PeerObjectName") {
		peerObject = urlQuery.Get("PeerObjectName")
	}

	// Determine peer object ID
	peerObjectID := urlQuery.Get("PeerObjectID")
	if !urlQuery.Has("PeerObjectID") || peerObject == "Profile" {
		subject, _ := claims.GetSubject()
		peerObjectID = subject
	}

	// Parse invert relationship flag
	invertRelationshipStr := "false"
	if urlQuery.Has("InvertRelation") {
		invertRelationshipStr = urlQuery.Get("InvertRelation")
	}
	invertRelationship, _ := strconv.ParseBool(invertRelationshipStr)

	// Create the request object
	request := &profilev1.ListRelationshipRequest{
		PeerName:           peerObject,
		PeerId:             peerObjectID,
		InvertRelation:     invertRelationship,
		LastRelationshipId: lastRelationshipID,
	}

	// Apply count limits with bounds check
	if count > math.MaxInt32 {
		request.Count = math.MaxInt32
	} else {
		// Safe conversion with bounds check
		request.Count = int32(count) // #nosec G109,G115 -- bounds checked above
	}

	return request
}

// fetchRelationships fetches relationships based on the request.
func (ps *ProfileServer) fetchRelationships(
	ctx context.Context,
	request *profilev1.ListRelationshipRequest,
) ([]*profilev1.RelationshipObject, string, error) {
	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	relationshipBusiness := business.NewRelationshipBusiness(ctx, ps.Service, profileBusiness)

	// Get relationships from business layer
	relationships, err := relationshipBusiness.ListRelationships(ctx, request)
	if err != nil {
		if !frame.ErrorIsNoRows(err) {
			return nil, "", err
		}
		relationships = []*models.Relationship{}
	}

	// Convert relationships to API objects
	relationshipObjectList := make([]*profilev1.RelationshipObject, 0, len(relationships))
	lastRelationshipID := ""

	for _, relationship := range relationships {
		relationshipObject, err0 := relationshipBusiness.ToAPI(ctx, relationship, request.GetInvertRelation())
		if err0 != nil {
			return nil, "", err0
		}
		relationshipObjectList = append(relationshipObjectList, relationshipObject)
	}

	// Set the last relationship ID for pagination
	if len(relationshipObjectList) > 0 {
		lastRelationshipID = relationshipObjectList[len(relationshipObjectList)-1].GetId()
	}

	return relationshipObjectList, lastRelationshipID, nil
}

// sendRelationshipListResponse sends the formatted response to the client.
func (ps *ProfileServer) sendRelationshipListResponse(
	_ context.Context,
	rw http.ResponseWriter,
	relationshipObjects []*profilev1.RelationshipObject,
	lastRelationshipID string,
	claims *frame.AuthenticationClaims,
) {
	// Create response object
	response := frame.JSONMap{
		"relationships":      relationshipObjects,
		"count":              len(relationshipObjects),
		"LastRelationshipID": lastRelationshipID,
		"tenant_id":          claims.GetTenantID(),
		"partition_id":       claims.GetPartitionID(),
	}

	// Write response
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(response)
}

func (ps *ProfileServer) RestUserInfo(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	claims := frame.ClaimsFromContext(ctx)

	if claims == nil {
		ps.writeError(ctx, rw, errors.New("claims can not be empty"), http.StatusForbidden)
		return
	}

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	subject, _ := claims.GetSubject()
	profile, err := profileBusiness.GetByID(ctx, subject)
	if err != nil {
		ps.writeError(ctx, rw, err, http.StatusInternalServerError)
		return
	}

	properties := profile.GetProperties().AsMap()
	response := frame.JSONMap{
		"sub":      profile.GetId(),
		"name":     properties["name"],
		"contacts": profile.GetContacts(),
		"url":      properties["profile_pic"],
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(response)
}

func (ps *ProfileServer) NewSecureRouterV1() *http.ServeMux {
	userServeMux := http.NewServeMux()

	userServeMux.HandleFunc("/user/info", ps.RestUserInfo)

	userServeMux.HandleFunc("/user/relations", ps.RestListRelationshipsEndpoint)

	return userServeMux
}
