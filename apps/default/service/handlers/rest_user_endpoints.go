package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
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

func (ps *ProfileServer) RestListRelationshipsEndpoint(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	claims := frame.ClaimsFromContext(ctx)

	log := ps.Service.Log(ctx).
		WithField("method", "RestListRelationshipsEndpoint")

	if claims == nil {
		ps.writeError(ctx, rw, errors.New("claims can not be empty"), http.StatusInternalServerError)
		return
	}

	urlQuery := req.URL.Query()
	log.WithField("url params", urlQuery).Debug("listing relationships request")

	countStr := urlQuery.Get("Count")

	count, err := strconv.Atoi(countStr)
	if err != nil {
		count = 20
	}

	// Safe casting with bounds check
	lastRelationshipID := ""
	if urlQuery.Has("LastRelationshipID") {
		lastRelationshipID = urlQuery.Get("LastRelationshipID")
	}

	peerObject := "Profile"
	if urlQuery.Has("PeerObjectName") {
		peerObject = urlQuery.Get("PeerObjectName")
	}

	peerObjectID := urlQuery.Get("PeerObjectID")
	if !urlQuery.Has("PeerObjectID") || peerObject == "Profile" {
		subject, _ := claims.GetSubject()
		peerObjectID = subject
	}

	invertRelationshipStr := "false"
	if urlQuery.Has("InvertRelation") {
		invertRelationshipStr = urlQuery.Get("InvertRelation")
	}

	invertRelationship, _ := strconv.ParseBool(invertRelationshipStr)

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	relationshipBusiness := business.NewRelationshipBusiness(ctx, ps.Service, profileBusiness)

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

	relationships, err := relationshipBusiness.ListRelationships(ctx, request)
	if err != nil {
		if !frame.ErrorIsNoRows(err) {
			ps.writeError(ctx, rw, err, http.StatusInternalServerError)
			return
		}
		relationships = []*models.Relationship{}
	}

	var relationshipObjectList []*profilev1.RelationshipObject
	relationshipObjectList = []*profilev1.RelationshipObject{}

	for _, relationship := range relationships {
		relationshipObject, err1 := relationshipBusiness.ToAPI(ctx, relationship, request.GetInvertRelation())
		if err1 != nil {
			ps.writeError(ctx, rw, err, http.StatusInternalServerError)
			return
		}

		relationshipObjectList = append(relationshipObjectList, relationshipObject)
	}

	if len(relationshipObjectList) > 0 {
		lastRelationshipID = relationshipObjectList[len(relationshipObjectList)-1].GetId()
	}

	response := map[string]any{
		"relationships":      relationshipObjectList,
		"count":              len(relationshipObjectList),
		"LastRelationshipID": lastRelationshipID,
	}

	response["tenant_id"] = claims.GetTenantID()
	response["partition_id"] = claims.GetPartitionID()

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

	response := map[string]any{
		"sub":      profile.GetId(),
		"name":     profile.GetProperties()["name"],
		"contacts": profile.GetContacts(),
		"url":      profile.GetProperties()["profile_pic"],
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
