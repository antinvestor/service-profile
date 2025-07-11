package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/devices/service/business"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
)

type DevicesServer struct {
	Service *frame.Service
}

func (ps *DevicesServer) writeError(ctx context.Context, w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")

	log := ps.Service.Log(ctx).
		WithField("code", code)

	log.WithError(err).Error("internal service error")
	w.WriteHeader(code)

	encodeErr := json.NewEncoder(w).Encode(fmt.Sprintf(" internal processing err message: %v", err))
	if encodeErr != nil {
		log.WithError(encodeErr).Error("error encoding error response")
	}
}

func (ps *DevicesServer) RestLogDeviceData(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	deviceBusiness := business.NewDeviceBusiness(ctx, ps.Service)

	var deviceLog models.DeviceLog
	err := json.NewDecoder(req.Body).Decode(&deviceLog)
	if err != nil {
		ps.writeError(ctx, rw, err, http.StatusBadRequest)
		return
	}

	rawIPData := frame.GetIP(req)
	clientIPList := strings.Split(rawIPData, ",")
	clientIP := strings.TrimSpace(clientIPList[0])

	deviceLog.Data["ip"] = clientIP

	err = deviceBusiness.LogDevice(ctx, &deviceLog)
	if err != nil {
		ps.writeError(ctx, rw, err, http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"id": deviceLog.GetID(),
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(response)
}

func (ps *DevicesServer) RestDeviceLinkProfile(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	deviceBusiness := business.NewDeviceBusiness(ctx, ps.Service)

	var linkData map[string]string
	err := json.NewDecoder(req.Body).Decode(&linkData)
	if err != nil {
		ps.writeError(ctx, rw, err, http.StatusBadRequest)
		return
	}

	linkID, ok := linkData["link_id"]
	if !ok {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(rw).Encode("missing parameters")
		return
	}
	profileID, ok := linkData["profile_id"]
	if !ok {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(rw).Encode("missing parameters")
		return
	}

	device, err := deviceBusiness.UpdateProfileID(ctx, linkID, profileID)
	if err != nil {
		ps.writeError(ctx, rw, err, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(device)
}

func (ps *DevicesServer) RestGetDeviceLogByID(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	claims := frame.ClaimsFromContext(ctx)

	if claims == nil {
		ps.writeError(ctx, rw, errors.New("claims can not be empty"), http.StatusInternalServerError)
		return
	}

	deviceLogID := req.PathValue("deviceLogId")

	deviceBusiness := business.NewDeviceBusiness(ctx, ps.Service)

	deviceLog, err := deviceBusiness.GetDeviceLogByID(ctx, deviceLogID)
	if err != nil {
		ps.writeError(ctx, rw, err, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(deviceLog)
}

func (ps *DevicesServer) RestGetDeviceByDeviceLogID(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	claims := frame.ClaimsFromContext(ctx)

	if claims == nil {
		ps.writeError(ctx, rw, errors.New("claims can not be empty"), http.StatusInternalServerError)
		return
	}

	deviceLogID := req.PathValue("deviceLogId")

	deviceBusiness := business.NewDeviceBusiness(ctx, ps.Service)

	deviceLog, err := deviceBusiness.GetDeviceLogByID(ctx, deviceLogID)
	if err != nil {
		ps.writeError(ctx, rw, err, http.StatusInternalServerError)
		return
	}

	if deviceLog.DeviceID == "" {
		ps.writeError(ctx, rw, errors.New("device id has not yet been processed"), http.StatusInternalServerError)
		return
	}

	device, err := deviceBusiness.GetByID(ctx, deviceLog.DeviceID)
	if err != nil {
		ps.writeError(ctx, rw, err, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(device)
}

func (ps *DevicesServer) RestGetDeviceByID(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	claims := frame.ClaimsFromContext(ctx)

	if claims == nil {
		ps.writeError(ctx, rw, errors.New("claims can not be empty"), http.StatusInternalServerError)
		return
	}

	deviceBusiness := business.NewDeviceBusiness(ctx, ps.Service)

	deviceID := claims.GetDeviceID()

	device, err := deviceBusiness.GetByID(ctx, deviceID)
	if err != nil {
		ps.writeError(ctx, rw, err, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(device)
}

func (ps *DevicesServer) RestGetDevicesByProfileID(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	claims := frame.ClaimsFromContext(ctx)

	if claims == nil {
		ps.writeError(ctx, rw, errors.New("claims can not be empty"), http.StatusInternalServerError)
		return
	}

	profileID, err := claims.GetSubject()
	if err != nil {
		ps.writeError(ctx, rw, errors.New("could not obtain profile Id"), http.StatusInternalServerError)
		return
	}

	deviceBusiness := business.NewDeviceBusiness(ctx, ps.Service)
	deviceList, err := deviceBusiness.GetByProfileID(ctx, profileID)
	if err != nil {
		ps.writeError(ctx, rw, err, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(deviceList)
}

func (ps *DevicesServer) NewSecureRouterV1() *http.ServeMux {
	userServeMux := http.NewServeMux()

	userServeMux.HandleFunc("/user/device/by_id", ps.RestGetDeviceByID)
	userServeMux.HandleFunc("/user/device/by_profile_id", ps.RestGetDevicesByProfileID)
	userServeMux.HandleFunc("/user/device/by_device_log_id/{deviceLogId}", ps.RestGetDeviceByDeviceLogID)
	userServeMux.HandleFunc("/user/device/log/{deviceLogId}", ps.RestGetDeviceLogByID)

	return userServeMux
}

func (ps *DevicesServer) NewInSecureRouterV1() *http.ServeMux {
	userServeMux := http.NewServeMux()
	userServeMux.HandleFunc("/device/log", ps.RestLogDeviceData)
	userServeMux.HandleFunc("/device/link", ps.RestDeviceLinkProfile)

	return userServeMux
}
