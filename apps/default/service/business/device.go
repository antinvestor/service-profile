package business

import (
	"context"
	"errors"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"

	"github.com/pitabwire/frame"
)

type DeviceBusiness interface {
	GetByID(ctx context.Context, deviceID string) (*models.Device, error)
	GetByLinkID(ctx context.Context, linkID string) (*models.Device, error)
	GetByProfileID(ctx context.Context, linkID string) ([]*models.Device, error)
	UpdateProfileID(ctx context.Context, linkID, profileID string) (*models.Device, error)
	LogDevice(ctx context.Context, logData *models.DeviceLog) error
	GetDeviceLogByID(ctx context.Context, deviceLogID string) (*models.DeviceLog, error)
	GetDeviceLogByDeviceID(ctx context.Context, deviceID string) ([]*models.DeviceLog, error)
}

func NewDeviceBusiness(_ context.Context, service *frame.Service) DeviceBusiness {
	deviceRepo := repository.NewDeviceRepository(service)
	deviceLogRepo := repository.NewDeviceLogRepository(service)
	return &deviceBusiness{
		service:             service,
		deviceRepository:    deviceRepo,
		deviceLogRepository: deviceLogRepo,
	}
}

type deviceBusiness struct {
	service             *frame.Service
	deviceRepository    repository.DeviceRepository
	deviceLogRepository repository.DeviceLogRepository
}

func (dB *deviceBusiness) GetByID(ctx context.Context, deviceID string) (*models.Device, error) {
	return dB.deviceRepository.GetByID(ctx, deviceID)
}

func (dB *deviceBusiness) GetByLinkID(ctx context.Context, linkID string) (*models.Device, error) {
	device, err := dB.deviceRepository.GetByLinkID(ctx, linkID)
	if err != nil {
		if !frame.ErrorIsNoRows(err) {
			return nil, err
		}
	}

	if device != nil {
		return device, nil
	}

	deviceLog, err := dB.deviceLogRepository.GetByLinkID(ctx, linkID)
	if err != nil {
		return nil, err
	}

	if deviceLog.DeviceID == "" {
		return nil, errors.New("device log not yet successfully processed")
	}

	return dB.deviceRepository.GetByID(ctx, deviceLog.DeviceID)
}

func (dB *deviceBusiness) UpdateProfileID(ctx context.Context, linkID, profileID string) (*models.Device, error) {
	device, err := dB.GetByLinkID(ctx, linkID)
	if err != nil {
		return nil, err
	}

	if device.ProfileID == "" {
		device.ProfileID = profileID
		err = dB.deviceRepository.Save(ctx, device)
		if err != nil {
			return nil, err
		}
	}

	return device, nil
}

func (dB *deviceBusiness) GetByProfileID(ctx context.Context, profileID string) ([]*models.Device, error) {
	return dB.deviceRepository.List(ctx, profileID)
}

func (dB *deviceBusiness) LogDevice(ctx context.Context, logData *models.DeviceLog) error {
	logData.GenID(ctx)

	profileConfig := dB.service.Config().(*config.ProfileConfig)

	err := dB.deviceLogRepository.Save(ctx, logData)
	if err != nil {
		return err
	}

	payload := map[string]string{
		"id": logData.GetID(),
	}

	return dB.service.Publish(ctx, profileConfig.QueueDeviceAnalysisName, payload)
}

func (dB *deviceBusiness) GetDeviceLogByID(ctx context.Context, deviceLogID string) (*models.DeviceLog, error) {
	return dB.deviceLogRepository.GetByID(ctx, deviceLogID)
}

func (dB *deviceBusiness) GetDeviceLogByDeviceID(ctx context.Context, deviceID string) ([]*models.DeviceLog, error) {
	return dB.deviceLogRepository.ListByDeviceID(ctx, deviceID)
}
