package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service/business"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/pitabwire/frame"
	"net/http"
	"strconv"
)

func (ps *ProfileServer) writeError(ctx context.Context, w http.ResponseWriter, err error, code int) {

	w.Header().Set("Content-Type", "application/json")

	log := ps.Service.L(ctx).
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

	log := ps.Service.L(ctx).
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

	response["tenant_id"] = claims.GetTenantId()
	response["partition_id"] = claims.GetPartitionId()

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(response)
}

func (ps *ProfileServer) RestUserInfo(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	claims := frame.ClaimsFromContext(ctx)

	if claims == nil {
		ps.writeError(ctx, rw, errors.New("claims can not be empty"), http.StatusInternalServerError)
		return
	}

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service, ps.EncryptionKeyFunc)
	subject, _ := claims.GetSubject()
	profile, err := profileBusiness.GetByID(ctx, subject)
	if err != nil {
		ps.writeError(ctx, rw, err, 500)
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
	userServeMux.HandleFunc("/user/device/by_id", ps.RestGetDeviceByID)
	userServeMux.HandleFunc("/user/device/by_profile_id", ps.RestGetDevicesByProfileID)
	userServeMux.HandleFunc("/user/device/by_device_log_id/{deviceLogId}", ps.RestGetDeviceByDeviceLogID)
	userServeMux.HandleFunc("/user/device/log/{deviceLogId}", ps.RestGetDeviceLogByID)

	return userServeMux
}

func (ps *ProfileServer) NewInSecureRouterV1() *http.ServeMux {

	userServeMux := http.NewServeMux()
	userServeMux.HandleFunc("/device/log", ps.RestLogDeviceData)
	userServeMux.HandleFunc("/device/link", ps.RestDeviceLinkProfile)

	return userServeMux
}
