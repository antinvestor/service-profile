package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/queue"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/workerpool"
	"github.com/pitabwire/util"
	"go.opentelemetry.io/otel/attribute"

	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/caching"
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

	cache *caching.DeviceCacheService
}

// NewDeviceBusiness creates a new instance of DeviceBusiness.
func NewDeviceBusiness(_ context.Context, cfg *config.DevicesConfig, qMan queue.Manager,
	workMan workerpool.Manager, deviceRepo repository.DeviceRepository,
	deviceLogRepo repository.DeviceLogRepository, sessionRepo repository.DeviceSessionRepository,
	cacheSvc *caching.DeviceCacheService) DeviceBusiness {
	return &deviceBusiness{
		cfg:           cfg,
		qMan:          qMan,
		workMan:       workMan,
		deviceRepo:    deviceRepo,
		deviceLogRepo: deviceLogRepo,
		sessionRepo:   sessionRepo,
		cache:         cacheSvc,
	}
}

func (b *deviceBusiness) LogDeviceActivity(
	ctx context.Context,
	deviceID, sessionID string,
	logData data.JSONMap,
) (*devicev1.DeviceLog, error) {
	ctx, span := caching.StartSpan(ctx, "LogDeviceActivity",
		attribute.String("device_id", deviceID))
	defer caching.EndSpan(ctx, span, nil)

	// Rate limit log events per device.
	if b.cache != nil && b.cfg.RateLimitLogPerMinute > 0 {
		allowed, count := b.cache.CheckLogRateLimit(ctx, deviceID, b.cfg.RateLimitLogPerMinute)
		if !allowed {
			caching.RecordRateLimited(ctx, "log_device_activity")
			util.Log(ctx).WithField("device_id", deviceID).WithField("count", count).
				Warn("device log rate limited")
			return nil, connect.NewError(connect.CodeResourceExhausted,
				errors.New("device log rate limit exceeded"))
		}
	}

	log := &models.DeviceLog{
		DeviceID:        deviceID,
		DeviceSessionID: sessionID,
		Data:            logData,
	}
	log.GenID(ctx)

	if createErr := b.deviceLogRepo.Create(ctx, log); createErr != nil {
		return nil, createErr
	}

	// Publish to queue for further analysis.
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

	// Invalidate cache after mutation.
	b.invalidateDeviceCache(ctx, id)

	return b.GetDeviceByID(ctx, id)
}

// cachedDeviceResult is a serializable container for device + session data
// stored in the cache to avoid repeated DB lookups.
type cachedDeviceResult struct {
	Device  *models.Device        `json:"device"`
	Session *models.DeviceSession `json:"session,omitempty"`
}

func (b *deviceBusiness) GetDeviceByID(ctx context.Context, id string) (*devicev1.DeviceObject, error) {
	ctx, span := caching.StartSpan(ctx, "GetDeviceByID",
		attribute.String("device_id", id))
	defer caching.EndSpan(ctx, span, nil)

	// Try cache first.
	if b.cache != nil {
		if cached, found := b.cache.GetDevice(ctx, id); found {
			var result cachedDeviceResult
			if unmarshalErr := json.Unmarshal(cached, &result); unmarshalErr == nil {
				caching.RecordCacheHit(ctx, "device")
				return result.Device.ToAPI(result.Session), nil
			}
		}
		caching.RecordCacheMiss(ctx, "device")
	}

	dev, err := b.deviceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	sess, err := b.sessionRepo.GetLastByDeviceID(ctx, id)
	if err != nil {
		// Session may not exist yet for newly created devices.
		if !data.ErrorIsNoRows(err) {
			return nil, err
		}
		sess = nil
	}

	// Populate cache.
	b.cacheDeviceResult(ctx, id, dev, sess)

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

	// Add additional properties from the request.
	for _, p := range query.GetProperties() {
		searchProperties[fmt.Sprintf("%s = ? ", p)] = query.GetQuery()
	}

	// Build filters map, only add profile_id if it's not empty.
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
// It batch-loads sessions for all devices in each chunk to avoid N+1 queries.
func (b *deviceBusiness) processSearchResults(
	ctx context.Context,
	devicesResult workerpool.JobResultPipe[[]*models.Device],
	result workerpool.JobResultPipe[[]*devicev1.DeviceObject],
) error {
	for {
		res, ok := devicesResult.ReadResult(ctx)
		if !ok {
			return nil
		}

		if res.IsError() {
			return res.Error()
		}

		devices := res.Item()

		// Collect all device IDs for batch session loading.
		deviceIDs := make([]string, 0, len(devices))
		for _, device := range devices {
			deviceIDs = append(deviceIDs, device.GetID())
		}

		// Single query for all sessions instead of one per device.
		sessionMap, sessErr := b.sessionRepo.GetLatestByDeviceIDs(ctx, deviceIDs)
		if sessErr != nil {
			sessionMap = map[string]*models.DeviceSession{}
		}

		apiDevices := make([]*devicev1.DeviceObject, 0, len(devices))
		for _, device := range devices {
			sess := sessionMap[device.GetID()]
			apiDevices = append(apiDevices, device.ToAPI(sess))

			// Opportunistically warm the cache for each device.
			b.cacheDeviceResult(ctx, device.GetID(), device, sess)
		}

		if err := result.WriteResult(ctx, apiDevices); err != nil {
			return err
		}
	}
}

func (b *deviceBusiness) GetDeviceBySessionID(ctx context.Context, id string) (*devicev1.DeviceObject, error) {
	// Try session cache.
	if b.cache != nil {
		if cached, found := b.cache.GetSession(ctx, id); found {
			var sess models.DeviceSession
			if err := json.Unmarshal(cached, &sess); err == nil {
				// Now get the device (which may also be cached).
				return b.getDeviceWithSession(ctx, &sess)
			}
		}
	}

	sess, err := b.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache session.
	b.cacheSession(ctx, sess)

	return b.getDeviceWithSession(ctx, sess)
}

// getDeviceWithSession retrieves the device for a session and returns the API object.
func (b *deviceBusiness) getDeviceWithSession(
	ctx context.Context,
	sess *models.DeviceSession,
) (*devicev1.DeviceObject, error) {
	dev, err := b.deviceRepo.GetByID(ctx, sess.DeviceID)
	if err != nil {
		return nil, err
	}

	// Cache the device+session pair.
	b.cacheDeviceResult(ctx, dev.GetID(), dev, sess)

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

		// Invalidate cache after profile link.
		b.invalidateDeviceCache(ctx, device.GetID())
	}

	return device.ToAPI(session), nil
}

func (b *deviceBusiness) RemoveDevice(ctx context.Context, id string) error {
	_, err := b.deviceRepo.RemoveByID(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate all caches for this device.
	b.invalidateDeviceCache(ctx, id)
	if b.cache != nil {
		b.cache.InvalidateDeviceKeys(ctx, id)
		b.cache.InvalidatePresence(ctx, id)
	}

	return nil
}

// --- Cache helpers ---

func (b *deviceBusiness) cacheDeviceResult(
	ctx context.Context,
	deviceID string,
	dev *models.Device,
	sess *models.DeviceSession,
) {
	if b.cache == nil {
		return
	}
	result := cachedDeviceResult{Device: dev, Session: sess}
	encoded, err := json.Marshal(result)
	if err != nil {
		return
	}
	b.cache.SetDevice(ctx, deviceID, encoded)

	// Also cache the latest session for this device.
	if sess != nil {
		b.cacheSession(ctx, sess)
		sessEncoded, sessErr := json.Marshal(sess)
		if sessErr != nil {
			return
		}
		b.cache.SetLatestSession(ctx, deviceID, sessEncoded)
	}
}

func (b *deviceBusiness) cacheSession(ctx context.Context, sess *models.DeviceSession) {
	if b.cache == nil || sess == nil {
		return
	}
	encoded, err := json.Marshal(sess)
	if err != nil {
		return
	}
	b.cache.SetSession(ctx, sess.GetID(), encoded)
}

func (b *deviceBusiness) invalidateDeviceCache(ctx context.Context, deviceID string) {
	if b.cache == nil {
		return
	}
	b.cache.InvalidateDevice(ctx, deviceID)
	b.cache.InvalidateLatestSession(ctx, deviceID)
}
