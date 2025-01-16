package handlers

import (
	"encoding/json"
	"errors"
	"github.com/antinvestor/service-profile/service/business"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/gorilla/mux"
	"github.com/pitabwire/frame"
	"net/http"
	"strings"
)

func (ps *ProfileServer) RestLogDeviceData(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	deviceBusiness := business.NewDeviceBusiness(ctx, ps.Service)

	var deviceLog models.DeviceLog
	err := json.NewDecoder(req.Body).Decode(&deviceLog)
	if err != nil {
		ps.writeError(ctx, rw, err, http.StatusBadRequest)
		return
	}

	rawIpData := frame.GetIp(req)
	clientIPList := strings.Split(rawIpData, ",")
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

func (ps *ProfileServer) RestGetDeviceLogByID(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	claims := frame.ClaimsFromContext(ctx)

	if claims == nil {
		ps.writeError(ctx, rw, errors.New("claims can not be empty"), http.StatusInternalServerError)
		return
	}

	params := mux.Vars(req)
	deviceLogID := params["deviceLogId"]

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

func (ps *ProfileServer) RestGetDeviceByDeviceLogID(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	claims := frame.ClaimsFromContext(ctx)

	if claims == nil {
		ps.writeError(ctx, rw, errors.New("claims can not be empty"), http.StatusInternalServerError)
		return
	}

	params := mux.Vars(req)
	deviceLogID := params["deviceLogId"]

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

func (ps *ProfileServer) RestGetDeviceByID(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	claims := frame.ClaimsFromContext(ctx)

	if claims == nil {
		ps.writeError(ctx, rw, errors.New("claims can not be empty"), http.StatusInternalServerError)
		return
	}

	deviceBusiness := business.NewDeviceBusiness(ctx, ps.Service)

	deviceID := claims.GetDeviceId()

	device, err := deviceBusiness.GetByID(ctx, deviceID)
	if err != nil {
		ps.writeError(ctx, rw, err, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(device)
}

func (ps *ProfileServer) RestGetDevicesByProfileID(rw http.ResponseWriter, req *http.Request) {
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
