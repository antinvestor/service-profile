package handlers

import (
	"encoding/json"
	"fmt"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service/business"
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

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	relationshipBusiness := business.NewRelationshipBusiness(ctx, ps.Service, profileBusiness)

	request := &profilev1.ListRelationshipRequest{
		Parent:             "Profile",
		ParentId:           claims.ProfileID,
		LastRelationshipId: lastRelationshipID,
		Count:              int32(count),
	}

	relationships, err := relationshipBusiness.ListRelationships(ctx, request)
	if err != nil {
		ps.writeError(rw, err, 500)
		return
	}

	var relationshipObjectList []*profilev1.RelationshipObject

	for _, relationship := range relationships {

		relationshipObject, err1 := relationshipBusiness.ToAPI(ctx, request.GetParent(), request.GetParentId(), relationship)
		if err1 != nil {
			ps.writeError(rw, err, 500)
			return
		}

		relationshipObjectList = append(relationshipObjectList, relationshipObject)
	}

	if len(relationshipObjectList) > 0 {
		lastRelationshipID = relationshipObjectList[len(relationshipObjectList)-1].GetId()
	}

	response := map[string]interface{}{
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
	profile, err := profileBusiness.GetByID(ctx, claims.ProfileID)
	if err != nil {
		ps.writeError(rw, err, 500)
		return
	}

	response := map[string]interface{}{
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
		Name("UserInfoEndpoint").
		HandlerFunc(ps.RestListRelationshipsEndpoint).
		Methods("GET")

	return router
}
