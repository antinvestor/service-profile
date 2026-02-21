package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/devices/service/business"
)

// TURNHandler handles TURN credential requests via REST.
type TURNHandler struct {
	turnBusiness business.TURNBusiness
}

// NewTURNHandler creates a new TURNHandler.
func NewTURNHandler(turnBusiness business.TURNBusiness) *TURNHandler {
	return &TURNHandler{
		turnBusiness: turnBusiness,
	}
}

// GetTurnCredentials generates short-lived TURN credentials using the configured TURN provider.
func (h *TURNHandler) GetTurnCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	if h.turnBusiness == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":"TURN credentials provider is not configured"}`))
		return
	}

	credentials, err := h.turnBusiness.GetTurnCredentials(ctx)
	if err != nil {
		util.Log(ctx).WithError(err).Error("failed to get TURN credentials")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed to generate TURN credentials"}`))
		return
	}

	body, err := json.Marshal(credentials)
	if err != nil {
		util.Log(ctx).WithError(err).Error("failed to encode TURN credentials response")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"internal encoding error"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
