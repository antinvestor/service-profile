package queue

import (
	"context"
	"fmt"
	"strings"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/mssola/user_agent"
	"github.com/pitabwire/frame"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

type DeviceAnalysisQueueHandler struct {
	DeviceRepository    repository.DeviceRepository
	DeviceLogRepository repository.DeviceLogRepository
	SessionRepository   repository.DeviceSessionRepository
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

	var session *models.DeviceSession
	// The log is created with the device and session IDs, so we just need to update the session's LastSeen
	if deviceLog.DeviceSessionID == "" {
		dq.Service.Log(ctx).WithField("deviceLogID", deviceLogID).Warn("device log has no session ID, skipping")

		session, err = dq.createSessionFromLog(ctx, deviceLog)
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

	deviceID := session.DeviceID
	if deviceID == "" {

		device, err0 := dq.createDeviceFromSess(ctx, session)
		if err0 != nil {
			dq.Service.Log(ctx).WithError(err0).
				Warn("device could not be created from session")
			return nil
		}
		deviceID = device.ID
	}

	return nil
}

func (dq *DeviceAnalysisQueueHandler) createDeviceFromSess(ctx context.Context, session *models.DeviceSession) (*models.Device, error) {

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

func (dq *DeviceAnalysisQueueHandler) createSessionFromLog(ctx context.Context, deviceLog *models.DeviceLog) (*models.DeviceSession, error) {

	data := frame.DBPropertiesToMap(deviceLog.Data)

	sess := &models.DeviceSession{
		DeviceID:  deviceLog.DeviceID,
		UserAgent: data["userAgent"],
		IP:        data["ip"],
		LastSeen:  deviceLog.CreatedAt,
	}
	sess.GenID(ctx)

	geoIp, _ := QueryIPGeo(ctx, data["ip"])

	locale, err0 := dq.extractLocaleData(ctx, data, geoIp)
	if err0 != nil {
		return nil, err0
	}

	localeBytes, err := protojson.Marshal(locale)
	if err != nil {
		return nil, err
	}

	sess.Locale = localeBytes

	sess.Location = dq.extractLocationData(ctx, data, geoIp)

	err = dq.SessionRepository.Save(ctx, sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (dq *DeviceAnalysisQueueHandler) extractLocaleData(_ context.Context, data map[string]string, geoIP *GeoIP) (*devicev1.Locale, error) {

	var ok bool
	locale := devicev1.Locale{}
	locale.Timezone, ok = data["tz"]
	if !ok && geoIP != nil {
		locale.Timezone = geoIP.Timezone
	}

	languages, ok := data["lang"]
	if !ok && geoIP != nil {
		languages = geoIP.Languages
	}

	locale.Language = strings.Split(languages, ",")

	locale.Currency, ok = data["cur"]
	if !ok && geoIP != nil {
		locale.Currency = geoIP.Currency
	}

	locale.CurrencyName, ok = data["curNm"]
	if !ok && geoIP != nil {
		locale.CurrencyName = geoIP.CurrencyName
	}

	locale.Code, ok = data["code"]
	if !ok && geoIP != nil {
		locale.Code = geoIP.CountryCallingCode
	}

	return &locale, nil
}

func (dq *DeviceAnalysisQueueHandler) extractLocationData(_ context.Context, data map[string]string, geoIP *GeoIP) frame.JSONMap {

	locationData := map[string]string{}

	if geoIP != nil {
		locationData["country"] = geoIP.Country
		locationData["region"] = geoIP.Region
		locationData["city"] = geoIP.City
		locationData["latitude"] = fmt.Sprintf("%f", geoIP.Latitude)
		locationData["longitude"] = fmt.Sprintf("%f", geoIP.Longitude)
	}

	latitude, ok := data["lat"]
	if ok {
		locationData["latitude"] = latitude
	}

	longitude, ok := data["long"]
	if ok {
		locationData["longitude"] = longitude
	}

	return frame.DBPropertiesFromMap(locationData)
}
