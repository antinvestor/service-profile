package queue

import (
	"context"

	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

type DeviceAnalysisQueueHandler struct {
	DeviceRepository    repository.DeviceRepository
	DeviceLogRepository repository.DeviceLogRepository
	SessionRepository   repository.DeviceSessionRepository
	Service             *frame.Service
}

func (dq *DeviceAnalysisQueueHandler) Handle(ctx context.Context, idPayload map[string]string, _ []byte) error {
	deviceLogID := idPayload["id"]
	if deviceLogID == "" {
		dq.Service.Log(ctx).WithField("payload", idPayload).Warn("no device log id found in payload")
		return nil
	}

	ctx = frame.SkipTenancyChecksOnClaims(ctx)

	// Fetch the device log
	deviceLog, err := dq.DeviceLogRepository.GetByID(ctx, deviceLogID)
	if err != nil {
		if frame.ErrorIsNoRows(err) {
			dq.Service.Log(ctx).WithField("deviceLogID", deviceLogID).Warn("device log not found")
			return nil
		}
		return err
	}

	// The log is created with the device and session IDs, so we just need to update the session's LastSeen
	if deviceLog.DeviceSessionID == "" {
		dq.Service.Log(ctx).WithField("deviceLogID", deviceLogID).Warn("device log has no session ID, skipping")
		return nil
	}

	session, err := dq.SessionRepository.GetByID(ctx, deviceLog.DeviceSessionID)
	if err != nil {
		if frame.ErrorIsNoRows(err) {
			dq.Service.Log(ctx).WithField("sessionID", deviceLog.DeviceSessionID).Warn("device session not found")
			return nil
		}
		return err
	}

	// Update session's last seen timestamp
	session.LastSeen = deviceLog.CreatedAt

	return dq.SessionRepository.Save(ctx, session)
}
