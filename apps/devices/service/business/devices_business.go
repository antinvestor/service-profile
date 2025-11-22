package business

import (
	"context"
	"errors"
	"fmt"
	"time"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/queue"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/workerpool"

	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

const (
	// deviceLogRetentionDays defines the default retention period for device logs.
	deviceLogRetentionDays = 30
	// deviceLogFutureClockSkewDays accounts for clock skew in device timestamps.
	deviceLogFutureClockSkewDays = 1
)

// DeviceBusiness defines the interface for device-related business logic.
// It abstracts the underlying data storage and provides methods for interacting
// with device data in a consistent and transactional manner.
type DeviceBusiness interface {
	GetDeviceByID(ctx context.Context, id string) (*devicev1.DeviceObject, error)
	GetDeviceBySessionID(ctx context.Context, id string) (*devicev1.DeviceObject, error)
	SearchDevices(
		ctx context.Context,
		query *devicev1.SearchRequest,
	) (workerpool.JobResultPipe[[]*devicev1.DeviceObject], error)
	SaveDevice(ctx context.Context, id string, name string, data data.JSONMap) (*devicev1.DeviceObject, error)
	LinkDeviceToProfile(
		ctx context.Context,
		sessionID string,
		profileID string,
		data data.JSONMap,
	) (*devicev1.DeviceObject, error)
	RemoveDevice(ctx context.Context, id string) error

	LogDeviceActivity(
		ctx context.Context,
		deviceID, sessionID string,
		logData data.JSONMap,
	) (*devicev1.DeviceLog, error)
	GetDeviceLogs(
		ctx context.Context,
		deviceID string,
	) (workerpool.JobResultPipe[[]*devicev1.DeviceLog], error)
}

type deviceBusiness struct {
	cfg *config.DevicesConfig

	qMan    queue.Manager
	workMan workerpool.Manager

	deviceRepo    repository.DeviceRepository
	deviceLogRepo repository.DeviceLogRepository
	sessionRepo   repository.DeviceSessionRepository
}

// NewDeviceBusiness creates a new instance of DeviceBusiness.
func NewDeviceBusiness(_ context.Context, cfg *config.DevicesConfig, qMan queue.Manager,
	workMan workerpool.Manager, deviceRepo repository.DeviceRepository,
	deviceLogRepo repository.DeviceLogRepository, sessionRepo repository.DeviceSessionRepository) DeviceBusiness {
	return &deviceBusiness{
		cfg:           cfg,
		qMan:          qMan,
		workMan:       workMan,
		deviceRepo:    deviceRepo,
		deviceLogRepo: deviceLogRepo,
		sessionRepo:   sessionRepo,
	}
}

func (b *deviceBusiness) LogDeviceActivity(
	ctx context.Context,
	deviceID, sessionID string,
	logData data.JSONMap,
) (*devicev1.DeviceLog, error) {
	log := &models.DeviceLog{
		DeviceID:        deviceID,
		DeviceSessionID: sessionID,
		Data:            logData,
	}
	log.GenID(ctx)

	if err := b.deviceLogRepo.Create(ctx, log); err != nil {
		return nil, err
	}

	// Publish to queue for further analysis
	if b.cfg.QueueDeviceAnalysisName != "" {
		payload := data.JSONMap{"id": log.GetID()}
		_ = b.qMan.Publish(ctx, b.cfg.QueueDeviceAnalysisName, payload, nil)
	}

	return log.ToAPI(), nil
}

func (b *deviceBusiness) GetDeviceLogs(
	ctx context.Context,
	deviceID string,
) (workerpool.JobResultPipe[[]*devicev1.DeviceLog], error) {
	resultPipe := workerpool.NewJob[[]*devicev1.DeviceLog](
		func(ctx context.Context, result workerpool.JobResultPipe[[]*devicev1.DeviceLog]) error {
			logsResult, err := b.deviceLogRepo.GetByDeviceID(ctx, deviceID)
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

	err := workerpool.SubmitJob(ctx, b.workMan, resultPipe)
	if err != nil {
		return nil, err
	}

	return resultPipe, nil
}

func (b *deviceBusiness) SaveDevice(
	ctx context.Context,
	id string,
	name string,
	extra data.JSONMap,
) (*devicev1.DeviceObject, error) {
	sessionID := extra.GetString("session_id")

	_, logErr := b.LogDeviceActivity(ctx, id, sessionID, extra)
	if logErr != nil {
		return nil, logErr
	}

	if id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("device ID is required"))
	}

	dev, repoErr := b.deviceRepo.GetByID(ctx, id)
	if repoErr != nil {
		return nil, data.ErrorConvertToAPI(repoErr)
	}
	dev.Name = name
	_, updateErr := b.deviceRepo.Update(ctx, dev, "name")
	if updateErr != nil {
		return nil, data.ErrorConvertToAPI(updateErr)
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
	query *devicev1.SearchRequest,
) (workerpool.JobResultPipe[[]*devicev1.DeviceObject], error) {
	resultPipe := workerpool.NewJob[[]*devicev1.DeviceObject](
		func(ctx context.Context, result workerpool.JobResultPipe[[]*devicev1.DeviceObject]) error {
			return b.processSearchRequest(ctx, query, result)
		},
	)

	err := workerpool.SubmitJob(ctx, b.workMan, resultPipe)
	if err != nil {
		return nil, err
	}

	return resultPipe, nil
}

// processSearchRequest handles the main search logic.
func (b *deviceBusiness) processSearchRequest(
	ctx context.Context,
	query *devicev1.SearchRequest,
	result workerpool.JobResultPipe[[]*devicev1.DeviceObject],
) error {
	searchQuery := b.buildSearchQuery(ctx, query)

	devicesResult, err := b.deviceRepo.Search(ctx, searchQuery)
	if err != nil {
		return err
	}

	return b.processSearchResults(ctx, devicesResult, result)
}

// buildSearchQuery creates the search query from the request.
func (b *deviceBusiness) buildSearchQuery(ctx context.Context, query *devicev1.SearchRequest) *data.SearchQuery {
	profileID := ""
	claims := security.ClaimsFromContext(ctx)
	if claims != nil {
		profileID, _ = claims.GetSubject()
	}

	startDate, err := time.Parse(time.RFC3339, query.GetStartDate())
	if err != nil {
		startDate = time.Now().Add(-time.Duration(deviceLogRetentionDays) * 24 * time.Hour)
	}
	endDate, err := time.Parse(time.RFC3339, query.GetEndDate())
	if err != nil {
		endDate = time.Now().Add(time.Duration(deviceLogFutureClockSkewDays) * 24 * time.Hour)
	}

	searchProperties := map[string]any{}

	// Add additional properties from the request
	for _, p := range query.GetProperties() {
		searchProperties[fmt.Sprintf("%s = ? ", p)] = query.GetQuery()
	}

	// Build filters map, only add profile_id if it's not empty
	filters := map[string]any{}
	if profileID != "" {
		filters["profile_id"] = profileID
	}

	return data.NewSearchQuery(data.WithSearchLimit(int(query.GetCount())),
		data.WithSearchOffset(int(query.GetPage())), data.WithSearchByTimePeriod(&data.TimePeriod{
			Field:     "created_at",
			StartDate: &startDate,
			StopDate:  &endDate,
		}), data.WithSearchFiltersAndByValue(filters),
		data.WithSearchFiltersOrByValue(searchProperties))
}

// processSearchResults processes the search results and converts them to API objects.
func (b *deviceBusiness) processSearchResults(
	ctx context.Context,
	devicesResult workerpool.JobResultPipe[[]*models.Device],
	result workerpool.JobResultPipe[[]*devicev1.DeviceObject],
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
		// Continue with nil session if not found - this is not a critical error
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
	_ data.JSONMap,
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

		_, err = b.deviceRepo.Update(ctx, device, "profile_id")
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
