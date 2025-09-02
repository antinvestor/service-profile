package queue

import (
	"context"
	"encoding/json"
	"strings"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/mssola/user_agent"
	"github.com/pitabwire/frame"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

type DeviceAnalysisQueueHandler struct {
	DeviceRepository    repository.DeviceRepository
	DeviceLogRepository repository.DeviceLogRepository
	SessionRepository   repository.DeviceSessionRepository
	Service             *frame.Service
}

func (dq *DeviceAnalysisQueueHandler) Handle(ctx context.Context, _ map[string]string, payload []byte) error {
	var idPayload map[string]string
	err := json.Unmarshal(payload, &idPayload)
	if err != nil {
		return err
	}

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

	var session *models.DeviceSession
	// The log is created with the device and session IDs, so we just need to update the session's LastSeen
	if deviceLog.DeviceSessionID == "" {
		dq.Service.Log(ctx).WithField("deviceLogID", deviceLogID).Warn("device log has no session ID, skipping")

		session, err = dq.CreateSessionFromLog(ctx, deviceLog)
		if err != nil {
			dq.Service.Log(ctx).WithField("sessionID", deviceLog.DeviceSessionID).WithError(err).
				Warn("could not extract device session from log")
			return nil
		}
	} else {
		session, err = dq.SessionRepository.GetByID(ctx, deviceLog.DeviceSessionID)
		if err != nil {
			dq.Service.Log(ctx).WithField("sessionID", deviceLog.DeviceSessionID).WithError(err).
				Warn("device session not found")
			return nil
		}

		// Update session's last seen timestamp
		session.LastSeen = deviceLog.CreatedAt
		err = dq.SessionRepository.Save(ctx, session)
		if err != nil {
			return err
		}
	}

	var deviceID string
	if session.DeviceID != "" {
		deviceID = session.DeviceID
	} else {
		device, err0 := dq.CreateDeviceFromSess(ctx, session)
		if err0 != nil {
			dq.Service.Log(ctx).WithError(err0).
				Warn("could not extract device from session")
			return nil
		}
		deviceID = device.ID
	}

	// Use deviceID for any device-related operations here
	_ = deviceID // TODO: Implement device-specific logic

	return nil
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

	dev.GenID(ctx)

	err := dq.DeviceRepository.Save(ctx, dev)
	if err != nil {
		return nil, err
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
	sess.GenID(ctx)

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
	data frame.JSONMap,
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
	data frame.JSONMap,
	geoIP *GeoIP,
) frame.JSONMap {
	locationData := frame.JSONMap{}

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
