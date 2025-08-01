package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pitabwire/util"
)

func (ds *DevicesServer) writeError(ctx context.Context, w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")

	log := ds.Service.Log(ctx).
		WithField("code", code)

	log.WithError(err).Error("internal service error")
	w.WriteHeader(code)

	encodeErr := json.NewEncoder(w).Encode(fmt.Sprintf(" internal processing err message: %v", err))
	if encodeErr != nil {
		log.WithError(encodeErr).Error("error encoding error response")
	}
}

func (ds *DevicesServer) RestLogDeviceData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var deviceLogData map[string]string
	err := json.NewDecoder(r.Body).Decode(&deviceLogData)
	if err != nil {
		ds.writeError(ctx, w, err, http.StatusBadRequest)
		return
	}

	sessionID, ok := deviceLogData["session_id"]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("missing parameters")
		return
	}

	for k, v := range r.Header {
		deviceLogData[k] = v[0]
	}

	deviceLogData["ip"] = util.GetIP(r)

	devLog, err := ds.Biz.LogDeviceActivity(ctx, "", sessionID, deviceLogData)
	if err != nil {
		ds.writeError(ctx, w, err, http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"id": devLog.GetId(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (ds *DevicesServer) RestDeviceLinkProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var sessionData map[string]string
	err := json.NewDecoder(r.Body).Decode(&sessionData)
	if err != nil {
		ds.writeError(ctx, w, err, http.StatusBadRequest)
		return
	}

	sessionID, ok := sessionData["session_id"]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("missing parameters")
		return
	}
	profileID, ok := sessionData["profile_id"]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("missing parameters")
		return
	}

	device, err := ds.Biz.LinkDeviceToProfile(ctx, sessionID, profileID, sessionData)
	if err != nil {
		ds.writeError(ctx, w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(device)
}

func (ds *DevicesServer) NewInSecureRouterV1() *http.ServeMux {
	userServeMux := http.NewServeMux()
	userServeMux.HandleFunc("/device/log", ds.RestLogDeviceData)
	userServeMux.HandleFunc("/device/link", ds.RestDeviceLinkProfile)

	return userServeMux
}
