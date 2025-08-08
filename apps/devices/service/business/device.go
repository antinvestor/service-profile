package business

import (
	"context"
	"errors"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/framedata"

	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

const defaultMaxLogsCount = 1000

// DeviceBusiness defines the interface for device-related business logic.
// It abstracts the underlying data storage and provides methods for interacting
// with device data in a consistent and transactional manner.
type DeviceBusiness interface {
	GetDeviceByID(ctx context.Context, id string) (*devicev1.DeviceObject, error)
	GetDeviceBySessionID(ctx context.Context, id string) (*devicev1.DeviceObject, error)
	SearchDevices(
		ctx context.Context,
		query *devicev1.SearchRequest,
	) (frame.JobResultPipe[[]*devicev1.DeviceObject], error)
	SaveDevice(ctx context.Context, id string, name string, data map[string]string) (*devicev1.DeviceObject, error)
	LinkDeviceToProfile(
		ctx context.Context,
		sessionID string,
		profileID string,
		data map[string]string,
	) (*devicev1.DeviceObject, error)
	RemoveDevice(ctx context.Context, id string) error

	AddKey(
		ctx context.Context,
		deviceID string,
		_ devicev1.KeyType,
		key []byte,
		extra map[string]string,
	) (*devicev1.KeyObject, error)
	GetKeys(
		ctx context.Context,
		deviceID string,
		_ devicev1.KeyType,
	) (<-chan frame.JobResult[[]*devicev1.KeyObject], error)
	RemoveKeys(ctx context.Context, id ...string) (<-chan frame.JobResult[[]*devicev1.KeyObject], error)

	LogDeviceActivity(
		ctx context.Context,
		deviceID, sessionID string,
		data map[string]string,
	) (*devicev1.DeviceLog, error)
	GetDeviceLogs(ctx context.Context, deviceID string) (frame.JobResultPipe[[]*devicev1.DeviceLog], error)
}

type deviceBusiness struct {
	cfg           *config.DevicesConfig
	deviceRepo    repository.DeviceRepository
	deviceLogRepo repository.DeviceLogRepository
	sessionRepo   repository.DeviceSessionRepository
	deviceKeyRepo repository.DeviceKeyRepository
	service       *frame.Service
}

// NewDeviceBusiness creates a new instance of DeviceBusiness.
func NewDeviceBusiness(_ context.Context, service *frame.Service) DeviceBusiness {
	cfg, _ := service.Config().(*config.DevicesConfig)
	return &deviceBusiness{
		cfg:           cfg,
		deviceRepo:    repository.NewDeviceRepository(service),
		deviceLogRepo: repository.NewDeviceLogRepository(service),
		sessionRepo:   repository.NewDeviceSessionRepository(service),
		deviceKeyRepo: repository.NewDeviceKeyRepository(service),
		service:       service,
	}
}

func (b *deviceBusiness) LogDeviceActivity(
	ctx context.Context,
	deviceID, sessionID string,
	extra map[string]string,
) (*devicev1.DeviceLog, error) {
	log := &models.DeviceLog{
		DeviceID:        deviceID,
		DeviceSessionID: sessionID,
		Data:            frame.DBPropertiesFromMap(extra),
	}

	if err := b.deviceLogRepo.Save(ctx, log); err != nil {
		return nil, err
	}

	// Publish to queue for further analysis
	if b.cfg.QueueDeviceAnalysisName != "" {
		payload := map[string]string{"id": log.GetID()}
		_ = b.service.Publish(ctx, b.cfg.QueueDeviceAnalysisName, payload, nil)
	}

	return log.ToAPI(), nil
}

func (b *deviceBusiness) GetDeviceLogs(
	ctx context.Context,
	deviceID string,
) (frame.JobResultPipe[[]*devicev1.DeviceLog], error) {
	resultPipe := frame.NewJob[[]*devicev1.DeviceLog](
		func(ctx context.Context, result frame.JobResultPipe[[]*devicev1.DeviceLog]) error {
			searchProperties := map[string]any{
				"device_id": deviceID,
			}

			q := framedata.NewSearchQuery("", searchProperties, 0, defaultMaxLogsCount)

			logsResult, err := b.deviceLogRepo.GetByDeviceID(ctx, q)
			if err != nil {
				return err
			}

			var apiDeviceLogs []*devicev1.DeviceLog
			for {
				res, ok := logsResult.ReadResult(ctx)
				if !ok {
					return nil
				}

				if res.IsError() {
					return res.Error()
				}

				for _, deviceLog := range res.Item() {
					apiDeviceLogs = append(apiDeviceLogs, deviceLog.ToAPI())
				}

				err = result.WriteResult(ctx, apiDeviceLogs)
				if err != nil {
					return err
				}
			}
		},
	)

	err := frame.SubmitJob(ctx, b.service, resultPipe)
	if err != nil {
		return nil, err
	}

	return resultPipe, nil
}

func (b *deviceBusiness) SaveDevice(
	ctx context.Context,
	id string,
	name string,
	data map[string]string,
) (*devicev1.DeviceObject, error) {
	sessionID := data["session_id"]

	_, err := b.LogDeviceActivity(ctx, id, sessionID, data)
	if err != nil {
		return nil, err
	}

	if id == "" {
		return nil, errors.New("device ID is required")
	}

	dev, err := b.deviceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	dev.Name = name
	err = b.deviceRepo.Save(ctx, dev)
	if err != nil {
		return nil, err
	}
	return b.GetDeviceByID(ctx, id)
}

func (b *deviceBusiness) GetDeviceByID(ctx context.Context, id string) (*devicev1.DeviceObject, error) {
	dev, err := b.deviceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	sess, err := b.sessionRepo.GetLastByDeviceID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dev.ToAPI(sess), nil
}

func (b *deviceBusiness) SearchDevices(
	ctx context.Context,
	request *devicev1.SearchRequest,
) (frame.JobResultPipe[[]*devicev1.DeviceObject], error) {
	resultPipe := frame.NewJob[[]*devicev1.DeviceObject](
		func(ctx context.Context, result frame.JobResultPipe[[]*devicev1.DeviceObject]) error {
			return b.processSearchRequest(ctx, request, result)
		},
	)

	err := frame.SubmitJob(ctx, b.service, resultPipe)
	if err != nil {
		return nil, err
	}

	return resultPipe, nil
}

// processSearchRequest handles the main search logic.
func (b *deviceBusiness) processSearchRequest(
	ctx context.Context,
	request *devicev1.SearchRequest,
	result frame.JobResultPipe[[]*devicev1.DeviceObject],
) error {
	searchQuery := b.buildSearchQuery(ctx, request)

	devicesResult, err := b.deviceRepo.Search(ctx, searchQuery)
	if err != nil {
		return err
	}

	return b.processSearchResults(ctx, devicesResult, result)
}

// buildSearchQuery creates the search query from the request.
func (b *deviceBusiness) buildSearchQuery(ctx context.Context, request *devicev1.SearchRequest) *framedata.SearchQuery {
	profileID := ""
	claims := frame.ClaimsFromContext(ctx)
	if claims != nil {
		profileID, _ = claims.GetSubject()
	}

	searchProperties := map[string]any{
		"profile_id": profileID,
		"start_date": request.GetStartDate(),
		"end_date":   request.GetEndDate(),
	}

	for _, p := range request.GetProperties() {
		searchProperties[p] = request.GetQuery()
	}

	return framedata.NewSearchQuery(request.GetQuery(), searchProperties, int(request.GetPage()),
		int(request.GetCount()))
}

// processSearchResults processes the search results and converts them to API objects.
func (b *deviceBusiness) processSearchResults(
	ctx context.Context,
	devicesResult frame.JobResultPipe[[]*models.Device],
	result frame.JobResultPipe[[]*devicev1.DeviceObject],
) error {
	var apiDevices []*devicev1.DeviceObject

	for {
		res, ok := devicesResult.ReadResult(ctx)
		if !ok {
			return nil
		}

		if res.IsError() {
			return res.Error()
		}

		for _, device := range res.Item() {
			apiDevice := b.convertDeviceToAPI(ctx, device)
			apiDevices = append(apiDevices, apiDevice)
		}

		err := result.WriteResult(ctx, apiDevices)
		if err != nil {
			return err
		}
	}
}

// convertDeviceToAPI converts a device model to API object with session data.
func (b *deviceBusiness) convertDeviceToAPI(ctx context.Context, device *models.Device) *devicev1.DeviceObject {
	sess, sessionErr := b.sessionRepo.GetLastByDeviceID(ctx, device.GetID())
	if sessionErr != nil {
		// Continue with nil session if not found
		sess = nil
	}
	return device.ToAPI(sess)
}

func (b *deviceBusiness) GetDeviceBySessionID(ctx context.Context, id string) (*devicev1.DeviceObject, error) {
	sess, err := b.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	dev, err := b.deviceRepo.GetByID(ctx, sess.DeviceID)
	if err != nil {
		return nil, err
	}

	return dev.ToAPI(sess), nil
}

func (b *deviceBusiness) LinkDeviceToProfile(
	ctx context.Context,
	sessionID string,
	profileID string,
	_ map[string]string,
) (*devicev1.DeviceObject, error) {
	session, err := b.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	device, err := b.deviceRepo.GetByID(ctx, session.DeviceID)
	if err != nil {
		return nil, err
	}

	if device.ProfileID == "" {
		device.ProfileID = profileID

		err = b.deviceRepo.Save(ctx, device)
		if err != nil {
			return nil, err
		}
	}

	return device.ToAPI(session), nil
}

func (b *deviceBusiness) RemoveDevice(ctx context.Context, id string) error {
	_, err := b.deviceRepo.RemoveByID(ctx, id)
	return err
}

func (b *deviceBusiness) AddKey(
	ctx context.Context,
	deviceID string,
	_ devicev1.KeyType,
	key []byte,
	extra map[string]string,
) (*devicev1.KeyObject, error) {
	// Validate that the device exists before adding a key
	_, err := b.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	deviceKey := &models.DeviceKey{
		DeviceID: deviceID,
		Key:      key,
		Extra:    frame.DBPropertiesFromMap(extra),
	}

	err = b.deviceKeyRepo.Save(ctx, deviceKey)
	if err != nil {
		return nil, err
	}

	return deviceKey.ToAPI(), nil
}

func (b *deviceBusiness) GetKeys(
	ctx context.Context,
	deviceID string,
	_ devicev1.KeyType,
) (<-chan frame.JobResult[[]*devicev1.KeyObject], error) {
	out := make(chan frame.JobResult[[]*devicev1.KeyObject])

	go func() {
		defer close(out)

		keys, err := b.deviceKeyRepo.GetByDeviceID(ctx, deviceID)
		if err != nil {
			out <- frame.ErrorResult[[]*devicev1.KeyObject](err)
			return
		}

		apiKeys := make([]*devicev1.KeyObject, len(keys))
		for i, key := range keys {
			apiKeys[i] = key.ToAPI()
		}

		out <- frame.Result[[]*devicev1.KeyObject](apiKeys)
	}()

	return out, nil
}

func (b *deviceBusiness) RemoveKeys(
	ctx context.Context,
	id ...string,
) (<-chan frame.JobResult[[]*devicev1.KeyObject], error) {
	out := make(chan frame.JobResult[[]*devicev1.KeyObject])

	go func() {
		defer close(out)

		var removedKeys []*devicev1.KeyObject

		for _, keyID := range id {
			removedKey, err := b.deviceKeyRepo.RemoveByID(ctx, keyID)
			if err != nil {
				out <- frame.ErrorResult[[]*devicev1.KeyObject](err)
				return
			}
			if removedKey != nil {
				removedKeys = append(removedKeys, removedKey.ToAPI())
			}
		}

		out <- frame.Result[[]*devicev1.KeyObject](removedKeys)
	}()

	return out, nil
}
