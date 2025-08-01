package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/gorilla/mux"
	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/devices/service/business"
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

func (ds *DevicesServer) LogDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	sessionId := vars["sessionId"]

	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		ds.writeError(ctx, w, err, http.StatusBadRequest)
		return
	}

	log, err := ds.Biz.LogDeviceActivity(ctx, deviceId, sessionId, data)
	if err != nil {
		ds.writeError(ctx, w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(log)
}

func (ds *DevicesServer) GetDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	deviceBusiness := business.NewDeviceBusiness(ctx, ds.Service)

	vars := mux.Vars(r)
	deviceId := vars["id"]

	device, err := deviceBusiness.GetDeviceByID(ctx, deviceId)
	if err != nil {
		if frame.ErrorIsNoRows(err) {
			ds.writeError(ctx, w, err, http.StatusNotFound)
			return
		}
		ds.writeError(ctx, w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(device)
}

func (ds *DevicesServer) SearchDevices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var query devicev1.SearchRequest

	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		ds.writeError(ctx, w, err, http.StatusBadRequest)
		return
	}

	devices, err := ds.Biz.SearchDevices(ctx, &query)
	if err != nil {
		ds.writeError(ctx, w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(devices)
}

func (ds *DevicesServer) RestLogDeviceData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var deviceData []byte
	_, err := r.Body.Read(deviceData)
	if err != nil {
		ds.writeError(ctx, w, err, http.StatusBadRequest)
		return
	}

	data := make(map[string]string)
	for k, v := range r.Header {
		data[k] = v[0]
	}
	data["data"] = string(deviceData)

	devLog, err := ds.Biz.LogDeviceActivity(ctx, "", "", data)
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

	var linkData map[string]string
	err := json.NewDecoder(r.Body).Decode(&linkData)
	if err != nil {
		ds.writeError(ctx, w, err, http.StatusBadRequest)
		return
	}

	sessionID, ok := linkData["session_id"]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("missing parameters")
		return
	}
	profileID, ok := linkData["profile_id"]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("missing parameters")
		return
	}

	device, err := ds.Biz.LinkDeviceToProfile(ctx, sessionID, profileID, linkData)
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
