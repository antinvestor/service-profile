package queue

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
	"github.com/mssola/user_agent"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/security"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	// ErrDeviceLogIDMissing is returned when the device log ID is missing from the payload.
	ErrDeviceLogIDMissing = errors.New("device log id missing from payload")
	// ErrDeviceLogNotFound is returned when the device log is not found in the database.
	ErrDeviceLogNotFound = errors.New("device log not found")
)

type DeviceAnalysisQueueHandler struct {
	DeviceRepository    repository.DeviceRepository
	DeviceLogRepository repository.DeviceLogRepository
	SessionRepository   repository.DeviceSessionRepository
	Service             *frame.Service
}

func NewDeviceAnalysisQueueHandler(
	svc *frame.Service,
) frame.SubscribeWorker {
	return &DeviceAnalysisQueueHandler{
		Service:             svc,
		DeviceRepository:    repository.NewDeviceRepository(svc),
		DeviceLogRepository: repository.NewDeviceLogRepository(svc),
		SessionRepository:   repository.NewDeviceSessionRepository(svc),
	}
}

func (dq *DeviceAnalysisQueueHandler) Handle(ctx context.Context, _ map[string]string, payload []byte) error {
	deviceLog, err := dq.getDeviceLog(ctx, payload)
	if err != nil {
		// Ignore expected errors (missing ID, not found)
		if errors.Is(err, ErrDeviceLogIDMissing) || errors.Is(err, ErrDeviceLogNotFound) {
			return nil
		}
		return err
	}

	ctx = security.SkipTenancyChecksOnClaims(ctx)

	session, err := dq.getOrCreateSession(ctx, deviceLog)
	if err != nil {
		return err
	}

	_, err = dq.getOrCreateDevice(ctx, session)
	return err
}

func (dq *DeviceAnalysisQueueHandler) getDeviceLog(
	ctx context.Context,
	payload []byte,
) (*models.DeviceLog, error) {
	var idPayload map[string]string
	err := json.Unmarshal(payload, &idPayload)
	if err != nil {
		return nil, err
	}

	deviceLogID := idPayload["id"]
	if deviceLogID == "" {
		dq.Service.Log(ctx).WithField("payload", idPayload).Warn("no device log id found in payload")
		return nil, ErrDeviceLogIDMissing
	}

	deviceLog, err := dq.DeviceLogRepository.GetByID(ctx, deviceLogID)
	if err != nil {
		if data.ErrorIsNoRows(err) {
			dq.Service.Log(ctx).WithField("deviceLogID", deviceLogID).Warn("device log not found")
			return nil, ErrDeviceLogNotFound
		}
		return nil, err
	}

	return deviceLog, nil
}

func (dq *DeviceAnalysisQueueHandler) getOrCreateSession(
	ctx context.Context,
	deviceLog *models.DeviceLog,
) (*models.DeviceSession, error) {
	if deviceLog.DeviceSessionID == "" {
		return dq.createSessionFromLog(ctx, deviceLog)
	}

	session, err := dq.SessionRepository.GetByID(ctx, deviceLog.DeviceSessionID)
	if err == nil {
		return dq.updateSessionLastSeen(ctx, session, deviceLog)
	}

	if !data.ErrorIsNoRows(err) {
		dq.Service.Log(ctx).WithField("sessionID", deviceLog.DeviceSessionID).WithError(err).
			Warn("error fetching device session")
		return nil, err
	}

	// Session ID provided but doesn't exist - create it
	return dq.createSessionFromLog(ctx, deviceLog)
}

func (dq *DeviceAnalysisQueueHandler) updateSessionLastSeen(
	ctx context.Context,
	session *models.DeviceSession,
	deviceLog *models.DeviceLog,
) (*models.DeviceSession, error) {
	session.LastSeen = deviceLog.CreatedAt
	err := dq.SessionRepository.Save(ctx, session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (dq *DeviceAnalysisQueueHandler) createSessionFromLog(
	ctx context.Context,
	deviceLog *models.DeviceLog,
) (*models.DeviceSession, error) {
	dq.Service.Log(ctx).WithField("deviceLogID", deviceLog.GetID()).Info("creating session from device log")

	session, err := dq.CreateSessionFromLog(ctx, deviceLog)
	if err != nil {
		dq.Service.Log(ctx).WithField("sessionID", deviceLog.DeviceSessionID).WithError(err).
			Warn("could not create device session from log")
		return nil, err
	}
	return session, nil
}

func (dq *DeviceAnalysisQueueHandler) getOrCreateDevice(
	ctx context.Context,
	session *models.DeviceSession,
) (*models.Device, error) {
	if session.DeviceID == "" {
		return dq.createDeviceFromSession(ctx, session)
	}

	device, err := dq.DeviceRepository.GetByID(ctx, session.DeviceID)
	if err == nil {
		return device, nil
	}

	if !data.ErrorIsNoRows(err) {
		dq.Service.Log(ctx).WithField("deviceID", session.DeviceID).WithError(err).
			Warn("error fetching device")
		return nil, err
	}

	// Device ID provided but doesn't exist - create it
	return dq.createDeviceFromSession(ctx, session)
}

func (dq *DeviceAnalysisQueueHandler) createDeviceFromSession(
	ctx context.Context,
	session *models.DeviceSession,
) (*models.Device, error) {
	dq.Service.Log(ctx).WithField("sessionID", session.GetID()).Info("creating device from session")

	device, err := dq.CreateDeviceFromSess(ctx, session)
	if err != nil {
		dq.Service.Log(ctx).WithError(err).
			Warn("could not auto create device from session")
		return nil, err
	}
	return device, nil
}

func (dq *DeviceAnalysisQueueHandler) CreateDeviceFromSess(
	ctx context.Context,
	session *models.DeviceSession,
) (*models.Device, error) {
	ua := user_agent.New(session.UserAgent)

	dev := &models.Device{
		Name: ua.Platform(),
		OS:   ua.OSInfo().FullName,
	}

	// Use the device ID from the session if it was provided in the log
	dev.GenID(ctx)
	if session.DeviceID != "" {
		dev.ID = session.DeviceID
	}

	err := dq.DeviceRepository.Save(ctx, dev)
	if err != nil {
		return nil, err
	}

	// Update session to link it to the device if not already linked
	if session.DeviceID == "" {
		session.DeviceID = dev.ID
		err = dq.SessionRepository.Save(ctx, session)
		if err != nil {
			return nil, err
		}
	}

	return dev, nil
}

func (dq *DeviceAnalysisQueueHandler) CreateSessionFromLog(
	ctx context.Context,
	deviceLog *models.DeviceLog,
) (*models.DeviceSession, error) {
	data := deviceLog.Data

	sess := &models.DeviceSession{
		DeviceID: deviceLog.DeviceID,
		LastSeen: deviceLog.CreatedAt,
	}

	// Use the session ID from the log if provided, otherwise generate one
	sess.GenID(ctx)
	if deviceLog.DeviceSessionID != "" {
		sess.ID = deviceLog.DeviceSessionID
	}

	anyData, ok := data["userAgent"]
	if ok {
		sess.UserAgent, _ = anyData.(string)
	}

	anyData, ok = data["ip"]
	if ok {
		sess.IP, _ = anyData.(string)

		geoIP, _ := QueryIPGeo(ctx, dq.Service, sess.IP)

		locale, err0 := dq.ExtractLocaleData(ctx, data, geoIP)
		if err0 != nil {
			return nil, err0
		}

		localeBytes, err := protojson.Marshal(locale)
		if err != nil {
			return nil, err
		}

		sess.Locale = localeBytes

		sess.Location = dq.ExtractLocationData(ctx, data, geoIP)
	}

	err := dq.SessionRepository.Save(ctx, sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (dq *DeviceAnalysisQueueHandler) ExtractLocaleData(
	_ context.Context,
	data data.JSONMap,
	geoIP *GeoIP,
) (*devicev1.Locale, error) {
	locale := devicev1.Locale{}
	locale.Timezone = data.GetString("tz")

	if locale.GetTimezone() == "" && geoIP != nil {
		locale.Timezone = geoIP.Timezone
	}

	languages := data.GetString("lang")
	if languages == "" && geoIP != nil {
		languages = geoIP.Languages
	}

	locale.Language = strings.Split(languages, ",")

	locale.Currency = data.GetString("cur")
	if locale.GetCurrency() == "" && geoIP != nil {
		locale.Currency = geoIP.Currency
	}

	locale.CurrencyName = data.GetString("curNm")
	if locale.GetCurrencyName() == "" && geoIP != nil {
		locale.CurrencyName = geoIP.CurrencyName
	}

	locale.Code = data.GetString("code")
	if locale.GetCode() == "" && geoIP != nil {
		locale.Code = geoIP.CountryCallingCode
	}

	return &locale, nil
}

func (dq *DeviceAnalysisQueueHandler) ExtractLocationData(
	_ context.Context,
	data data.JSONMap,
	geoIP *GeoIP,
) data.JSONMap {
	locationData := data.JSONMap{}

	if geoIP != nil {
		locationData["country"] = geoIP.Country
		locationData["region"] = geoIP.Region
		locationData["city"] = geoIP.City
		locationData["latitude"] = geoIP.Latitude
		locationData["longitude"] = geoIP.Longitude
	}

	rawData, ok := data["lat"]
	if ok {
		locationData["latitude"] = rawData
	}

	rawData, ok = data["long"]
	if ok {
		locationData["longitude"] = rawData
	}

	return locationData
}
