package queue

import (
	"context"
	"fmt"

	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

type DeviceAnalysisQueueHandler struct {
	Service             *frame.Service
	DeviceRepository    repository.DeviceRepository
	DeviceLogRepository repository.DeviceLogRepository
}

func (dq *DeviceAnalysisQueueHandler) Handle(ctx context.Context, idPayload map[string]string, _ []byte) error {
	deviceLogID := idPayload["id"]
	if deviceLogID == "" {
		dq.Service.Log(ctx).WithField("payload", idPayload).Warn("no device log id found in payload")
		return nil
	}

	ctx = frame.SkipTenancyChecksOnClaims(ctx)

	deviceLog, err := dq.DeviceLogRepository.GetByID(ctx, deviceLogID)
	if err != nil {
		if frame.ErrorIsNoRows(err) {
			return nil
		}

		return err
	}

	if deviceLog.DeviceID != "" {
		return nil
	}

	deviceLog, err = dq.processDeviceLog(ctx, deviceLog)
	if err != nil {
		return err
	}

	err = dq.DeviceLogRepository.Save(ctx, deviceLog)
	if err != nil {
		return err
	}

	return err
}

func (dq *DeviceAnalysisQueueHandler) processDeviceLog(
	ctx context.Context,
	deviceLog *models.DeviceLog,
) (*models.DeviceLog, error) {
	var device *models.Device

	if deviceLog.DeviceID != "" {
		var err error
		device, err = dq.DeviceRepository.GetByID(ctx, deviceLog.DeviceID)
		if err != nil {
			if !frame.ErrorIsNoRows(err) {
				return deviceLog, err
			}
		}
	}

	if device == nil {
		device = &models.Device{}
		device.GenID(ctx)
	}

	err := updateDeviceProperties(ctx, device, deviceLog)
	if err != nil {
		return deviceLog, err
	}
	err = dq.DeviceRepository.Save(ctx, device)
	if err != nil {
		return deviceLog, err
	}
	deviceLog.DeviceID = device.ID

	return deviceLog, nil
}

func updateDeviceProperties(_ context.Context, device *models.Device, deviceLog *models.DeviceLog) error {
	device.LastSeen = deviceLog.CreatedAt

	if device.LinkID == "" {
		device.LinkID = deviceLog.LinkID
	}

	device.IP, _ = deviceLog.Data["ip"].(string)
	device.Location, _ = deviceLog.Data["location"].(map[string]any)

	system, _ := deviceLog.Data["system"].(map[string]any)
	device.OS, _ = system["platform"].(string)

	browser, _ := system["browser"].(map[string]any)
	device.Browser, _ = browser["name"].(string)
	browserVersion, _ := browser["version"].(string)

	device.Name = fmt.Sprintf("%s_%s", browser["name"], browserVersion)

	locales, _ := deviceLog.Data["locales"].(map[string]any)
	device.Locale = frame.JSONMap{}
	for k, v := range locales {
		device.Locale[k] = v
	}
	device.Location = frame.JSONMap{}

	return nil
}
