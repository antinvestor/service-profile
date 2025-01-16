package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/pitabwire/frame"
	"gorm.io/datatypes"
)

type DeviceAnalysisQueueHandler struct {
	Service             *frame.Service
	DeviceRepository    repository.DeviceRepository
	DeviceLogRepository repository.DeviceLogRepository
}

func (dq *DeviceAnalysisQueueHandler) Handle(ctx context.Context, _ map[string]string, payload []byte) error {

	var idPayload map[string]string
	err := json.Unmarshal(payload, &idPayload)
	if err != nil {
		return err
	}

	deviceLogId := idPayload["id"]

	ctx = frame.SkipTenancyChecksOnClaims(ctx)

	deviceLog, err := dq.DeviceLogRepository.GetByID(ctx, deviceLogId)
	if err != nil {

		if frame.DBErrorIsRecordNotFound(err) {
			return nil
		}

		return err
	}

	if deviceLog.DeviceID != "" {
		return nil
	}

	deviceLog = enrichDeviceLog(ctx, deviceLog)

	deviceEmbeddings := getEmbeddingsFromDeviceLog(ctx, deviceLog)

	var deviceList []*models.Device
	if len(deviceEmbeddings) > 0 {
		deviceList, err = dq.DeviceRepository.ListByEmbedding(ctx, deviceEmbeddings)
		if err != nil {
			return err
		}
	}

	var device *models.Device
	if len(deviceList) == 0 {
		device = &models.Device{}
		device.GenID(ctx)

	} else if len(deviceList) == 1 {

		device = deviceList[0]

	} else {
		device = narrowSimilarityChecks(ctx, deviceList, deviceLog)
	}

	err = updateDeviceProperties(ctx, device, deviceLog)
	if err != nil {
		return err
	}
	err = dq.DeviceRepository.Save(ctx, device)
	if err != nil {
		return err
	}

	deviceLog.DeviceID = device.ID
	err = dq.DeviceLogRepository.Save(ctx, deviceLog)
	if err != nil {
		return err
	}

	return err
}

func enrichDeviceLog(ctx context.Context, deviceLog *models.DeviceLog) *models.DeviceLog {
	// Add geo location information

	return deviceLog
}

func narrowSimilarityChecks(ctx context.Context, list []*models.Device, deviceLog *models.DeviceLog) *models.Device {
	return list[0]
}

func updateDeviceProperties(ctx context.Context, device *models.Device, deviceLog *models.DeviceLog) error {
	device.LastSeen = deviceLog.CreatedAt

	device.IP, _ = deviceLog.Data["ip"].(string)

	system, _ := deviceLog.Data["system"].(map[string]any)
	device.OS, _ = system["platform"].(string)

	browser, _ := system["browser"].(map[string]any)
	device.Browser, _ = browser["name"].(string)
	browserVersion, _ := browser["version"].(string)

	device.Name = fmt.Sprintf("%s_%s", browser["name"], browserVersion)

	locales, _ := deviceLog.Data["locales"].(map[string]any)
	device.Locale = datatypes.JSONMap{}
	for k, v := range locales {
		device.Locale[k] = v
	}
	device.Location = datatypes.JSONMap{}

	device.Embedding = deviceLog.Embedding

	return nil
}

func getEmbeddingsFromDeviceLog(ctx context.Context, deviceLog *models.DeviceLog) []float32 {

	if deviceLog.Embedding == nil {
		return nil
	}

	return deviceLog.Embedding.Slice()
}
