package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service/business"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/gorilla/mux"
	"github.com/pitabwire/frame"
	"net/http"
	"strconv"
)

func (ps *ProfileServer) writeError(w http.ResponseWriter, err error, code int) {

	w.Header().Set("Content-Type", "application/json")

	log := ps.Service.L().
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

	log := ps.Service.L().
		WithField("method", "RestListRelationshipsEndpoint")

	urlQuery := req.URL.Query()
	log.WithField("url params", urlQuery).Debug("listing relationships request")

	countStr := urlQuery.Get("Count")

	count, err := strconv.Atoi(countStr)
	if err != nil {
		count = 100
	}

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

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	relationshipBusiness := business.NewRelationshipBusiness(ctx, ps.Service, profileBusiness)

	request := &profilev1.ListRelationshipRequest{
		PeerName:           peerObject,
		PeerId:             peerObjectID,
		InvertRelation:     invertRelationship,
		LastRelationshipId: lastRelationshipID,
		Count:              int32(count),
	}

	relationships, err := relationshipBusiness.ListRelationships(ctx, request)
	if err != nil {

		if !frame.DBErrorIsRecordNotFound(err) {
			ps.writeError(rw, err, 500)
			return
		}
		relationships = []*models.Relationship{}
	}

	var relationshipObjectList []*profilev1.RelationshipObject
	relationshipObjectList = []*profilev1.RelationshipObject{}

	for _, relationship := range relationships {

		relationshipObject, err1 := relationshipBusiness.ToAPI(ctx, relationship, request.GetInvertRelation())
		if err1 != nil {
			ps.writeError(rw, err, 500)
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

	if claims != nil {
		response["tenant_id"] = claims.TenantId()
		response["partition_id"] = claims.PartitionId()
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(response)
}

func (ps *ProfileServer) RestUserInfo(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	claims := frame.ClaimsFromContext(ctx)

	if claims == nil {
		ps.writeError(rw, errors.New("claims can not be empty"), 500)
		return
	}

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	subject, _ := claims.GetSubject()
	profile, err := profileBusiness.GetByID(ctx, subject)
	if err != nil {
		ps.writeError(rw, err, 500)
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

func (ps *ProfileServer) NewRouterV1() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Path("/user/info").
		Name("UserInfoEndpoint").
		HandlerFunc(ps.RestUserInfo).
		Methods("GET")

	router.Path("/user/relations").
		Name("UserRelationsEndpoint").
		HandlerFunc(ps.RestListRelationshipsEndpoint).
		Methods("GET")

	return router
}
