package handlers

import (
	"encoding/json"
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

	ps.Service.L().
		WithField("code", code).
		WithError(err).Error("internal service error")
	w.WriteHeader(code)

	err = json.NewEncoder(w).Encode(fmt.Sprintf(" internal processing err message: %v", err))
	if err != nil {
		ps.Service.L().WithError(err).Error("could not write error to response")
	}
}

func (ps *ProfileServer) RestListRelationshipsEndpoint(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	claims := frame.ClaimsFromContext(ctx)

	params := mux.Vars(req)
	countStr := params["Count"]

	count, err := strconv.Atoi(countStr)
	if err != nil {
		count = 100
	}

	lastRelationshipID := params["LastRelationshipID"]

	peerObject, ok := params["ParentObjectName"]
	if !ok {
		peerObject = "Profile"
	}

	peerObjectID, ok1 := params["ParentObjectID"]
	if !ok1 || peerObject == "Profile" {
		subject, _ := claims.GetSubject()
		peerObjectID = subject
	}

	invertRelationship := false
	invertRelationshipStr, ok2 := params["InvertRelation"]
	if !ok2 {
		invertRelationship, _ = strconv.ParseBool(invertRelationshipStr)
	}

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
		relationshipObjectList = append(relationshipObjectList, relationship.ToAPI())
	}

	if len(relationshipObjectList) > 0 {
		lastRelationshipID = relationshipObjectList[len(relationshipObjectList)-1].GetId()
	}

	response := map[string]any{
		"tenant_id":          claims.TenantId(),
		"partition_id":       claims.PartitionId(),
		"relationships":      relationshipObjectList,
		"count":              len(relationshipObjectList),
		"LastRelationshipID": lastRelationshipID,
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(response)
}

func (ps *ProfileServer) RestUserInfo(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	claims := frame.ClaimsFromContext(ctx)

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
