package queue

import (
	"context"
	"fmt"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
	"github.com/pitabwire/frame"
)

type DeviceAnalysisQueueHandler struct {
	DeviceRepository    repository.DeviceRepository
	DeviceLogRepository repository.DeviceLogRepository
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

	if deviceLog.DeviceID != "" {
		return nil
	}

	deviceLog, err = dq.processDeviceLog(ctx, deviceLog)
	if err != nil {
		return err
	}

	// Update device log with device ID
	err = dq.DeviceLogRepository.Save(ctx, deviceLog)
	if err != nil {
		return fmt.Errorf("failed to update device log: %w", err)
	}

	return nil
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

	// Update device properties
	device.LastSeen = deviceLog.CreatedAt

	// Check if link ID is available in device log
	if device.LinkID == "" && deviceLog.LinkID != "" {
		device.LinkID = deviceLog.LinkID
	}

	device.IP, _ = deviceLog.Data["ip"].(string)
	device.Location, _ = deviceLog.Data["location"].(map[string]any)

	system, _ := deviceLog.Data["system"].(map[string]any)
	device.OS, _ = system["platform"].(string)
	browser, _ := deviceLog.Data["browser"].(map[string]any)
	device.Browser, _ = browser["name"].(string)
	browserVersion, _ := browser["version"].(string)

	device.Name = fmt.Sprintf("%s_%s", browser["name"], browserVersion)

	locale, _ := deviceLog.Data["locale"].(map[string]any)
	device.Locale = frame.JSONMap{}
	for k, v := range locale {
		device.Locale[k] = v
	}

	// Save the device
	err := dq.DeviceRepository.Save(ctx, device)
	if err != nil {
		return deviceLog, err
	}
	deviceLog.DeviceID = device.ID

	return deviceLog, nil
}
