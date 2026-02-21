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
	"golang.org/x/sync/singleflight"

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
	GetDevicesByIDs(ctx context.Context, ids []string) ([]*devicev1.DeviceObject, error)
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
	RemoveDevice(ctx context.Context, id string) (*devicev1.DeviceObject, error)

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

	// sfDevice collapses concurrent cache misses for the same device ID.
	sfDevice singleflight.Group
	// sfSession collapses concurrent cache misses for the same session ID.
	sfSession singleflight.Group
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
		if pubErr := b.qMan.Publish(ctx, b.cfg.QueueDeviceAnalysisName, payload, nil); pubErr != nil {
			util.Log(ctx).WithError(pubErr).WithField("device_log_id", log.GetID()).
				Warn("failed to publish device log to analysis queue")
		}
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
	// Validate ID before doing any work.
	if id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("device ID is required"))
	}

	sessionID := extra.GetString("session_id")

	_, logErr := b.LogDeviceActivity(ctx, id, sessionID, extra)
	if logErr != nil {
		return nil, logErr
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
	if obj, hit := b.tryDeviceCache(ctx, id); hit {
		return obj, nil
	}

	// Use singleflight to collapse concurrent misses for the same device.
	val, err, _ := b.sfDevice.Do(id, func() (any, error) {
		return b.fetchDeviceAndCache(ctx, id)
	})
	if err != nil {
		return nil, err
	}

	devObj, ok := val.(*devicev1.DeviceObject)
	if !ok {
		return nil, errors.New("unexpected type in singleflight result")
	}
	return devObj, nil
}

// tryDeviceCache attempts to read a device from cache.
// Returns the API object and true on hit, nil and false on miss.
func (b *deviceBusiness) tryDeviceCache(ctx context.Context, id string) (*devicev1.DeviceObject, bool) {
	if b.cache == nil {
		return nil, false
	}
	cached, found := b.cache.GetDevice(ctx, id)
	if !found {
		caching.RecordCacheMiss(ctx, "device")
		return nil, false
	}
	var result cachedDeviceResult
	if err := json.Unmarshal(cached, &result); err != nil {
		caching.RecordCacheMiss(ctx, "device")
		return nil, false
	}
	caching.RecordCacheHit(ctx, "device")
	return result.Device.ToAPI(result.Session), true
}

// fetchDeviceAndCache loads a device from the DB (with its latest session),
// populates the cache, and returns the API object. Used inside singleflight.
func (b *deviceBusiness) fetchDeviceAndCache(ctx context.Context, id string) (*devicev1.DeviceObject, error) {
	// Double-check cache in case another goroutine populated it.
	if obj, hit := b.tryDeviceCache(ctx, id); hit {
		return obj, nil
	}

	dev, err := b.deviceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	sess, sessErr := b.sessionRepo.GetLastByDeviceID(ctx, id)
	if sessErr != nil && !data.ErrorIsNoRows(sessErr) {
		return nil, sessErr
	}

	b.cacheDeviceResult(ctx, id, dev, sess)
	return dev.ToAPI(sess), nil
}

// GetDevicesByIDs retrieves multiple devices by ID in batch, avoiding N+1 queries.
func (b *deviceBusiness) GetDevicesByIDs(ctx context.Context, ids []string) ([]*devicev1.DeviceObject, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// For single ID, delegate to the cached path.
	if len(ids) == 1 {
		dev, err := b.GetDeviceByID(ctx, ids[0])
		if err != nil {
			return nil, err
		}
		return []*devicev1.DeviceObject{dev}, nil
	}

	results, missIDs := b.collectCachedDevices(ctx, ids)
	b.fetchAndMergeDeviceMisses(ctx, missIDs, results)

	// Preserve input order.
	ordered := make([]*devicev1.DeviceObject, 0, len(ids))
	for _, id := range ids {
		if obj, ok := results[id]; ok {
			ordered = append(ordered, obj)
		}
	}
	return ordered, nil
}

// collectCachedDevices partitions IDs into cache hits (populated in results) and misses.
func (b *deviceBusiness) collectCachedDevices(
	ctx context.Context,
	ids []string,
) (map[string]*devicev1.DeviceObject, []string) {
	results := make(map[string]*devicev1.DeviceObject, len(ids))
	if b.cache == nil {
		return results, ids
	}

	var missIDs []string
	for _, id := range ids {
		if obj, hit := b.tryDeviceCache(ctx, id); hit {
			results[id] = obj
		} else {
			missIDs = append(missIDs, id)
		}
	}
	return results, missIDs
}

// fetchAndMergeDeviceMisses fetches devices and sessions from DB for cache misses,
// merges sessions, populates cache, and updates results in-place.
func (b *deviceBusiness) fetchAndMergeDeviceMisses(
	ctx context.Context,
	missIDs []string,
	results map[string]*devicev1.DeviceObject,
) {
	if len(missIDs) == 0 {
		return
	}

	// Fetch devices from DB.
	devices := make(map[string]*models.Device, len(missIDs))
	for _, id := range missIDs {
		dev, err := b.deviceRepo.GetByID(ctx, id)
		if err != nil {
			continue
		}
		devices[id] = dev
		results[id] = dev.ToAPI(nil)
	}

	// Batch load sessions.
	sessionMap, sessErr := b.sessionRepo.GetLatestByDeviceIDs(ctx, missIDs)
	if sessErr != nil {
		return
	}

	// Merge sessions and populate cache.
	for id, dev := range devices {
		sess := sessionMap[id]
		if sess != nil {
			results[id] = dev.ToAPI(sess)
		}
		b.cacheDeviceResult(ctx, id, dev, sess)
	}
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
				// Use GetDeviceByID (which checks the device cache) instead of
				// going directly to the repo, completing the cache-aside pattern.
				return b.getDeviceWithCachedSession(ctx, &sess)
			}
		}
	}

	// Use singleflight to collapse concurrent session lookups.
	val, err, _ := b.sfSession.Do(id, func() (any, error) {
		sess, sessErr := b.sessionRepo.GetByID(ctx, id)
		if sessErr != nil {
			return nil, sessErr
		}

		// Cache session.
		b.cacheSession(ctx, sess)

		return b.getDeviceWithCachedSession(ctx, sess)
	})
	if err != nil {
		return nil, err
	}
	devObj, ok := val.(*devicev1.DeviceObject)
	if !ok {
		return nil, errors.New("unexpected type in singleflight result")
	}
	return devObj, nil
}

// getDeviceWithCachedSession retrieves the device for a session using the cache-aware path.
func (b *deviceBusiness) getDeviceWithCachedSession(
	ctx context.Context,
	sess *models.DeviceSession,
) (*devicev1.DeviceObject, error) {
	// Use GetDeviceByID which checks the device cache and uses singleflight.
	devObj, err := b.GetDeviceByID(ctx, sess.DeviceID)
	if err != nil {
		return nil, err
	}

	// Overlay session-specific fields onto the device object since GetDeviceByID
	// may have used a different (or no) session.
	devObj.SessionId = sess.GetID()
	devObj.UserAgent = sess.UserAgent
	devObj.Ip = sess.IP
	devObj.LastSeen = sess.LastSeen.String()

	return devObj, nil
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

// RemoveDevice deletes a device and returns the device data that was removed.
func (b *deviceBusiness) RemoveDevice(ctx context.Context, id string) (*devicev1.DeviceObject, error) {
	dev, err := b.deviceRepo.RemoveByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Invalidate all caches for this device.
	b.invalidateDeviceCache(ctx, id)
	if b.cache != nil {
		b.cache.InvalidateDeviceKeys(ctx, id)
		b.cache.InvalidatePresence(ctx, id)
	}

	return dev.ToAPI(nil), nil
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
		// Marshal session once and reuse for both caches.
		sessEncoded, sessErr := json.Marshal(sess)
		if sessErr != nil {
			return
		}
		b.cache.SetSession(ctx, sess.GetID(), sessEncoded)
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
