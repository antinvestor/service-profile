package queue_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/queue"
	"github.com/antinvestor/service-profile/apps/devices/tests"
)

const (
	DefaultRandomStringLength = 8
)

type QueueTestSuite struct {
	tests.DeviceBaseTestSuite
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}

func (suite *QueueTestSuite) TestDeviceAnalysisQueueHandler_Handle() {
	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		// Create a device first
		device := &models.Device{
			Name: "Test Device",
			OS:   "Linux",
		}
		device.GenID(ctx)
		err := deps.DeviceRepo.Create(ctx, device)
		require.NoError(t, err)

		// Create a session
		session := &models.DeviceSession{
			DeviceID:  device.ID,
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
			IP:        "192.168.1.1",
			LastSeen:  time.Now(),
		}
		session.GenID(ctx)
		err = deps.SessionRepo.Create(ctx, session)
		require.NoError(t, err)

		// Create a device log
		deviceLog := &models.DeviceLog{
			DeviceID:        device.ID,
			DeviceSessionID: session.ID,
			Data: data.JSONMap{
				"action":    "page_view",
				"userAgent": session.UserAgent,
				"ip":        session.IP,
			},
		}
		deviceLog.GenID(ctx)
		err = deps.DeviceLogRepo.Create(ctx, deviceLog)
		require.NoError(t, err)

		// Test cases
		testCases := []struct {
			name        string
			payload     data.JSONMap
			expectError bool
		}{
			{
				name: "handle_valid_device_log",
				payload: data.JSONMap{
					"id": deviceLog.ID,
				},
				expectError: false,
			},
			{
				name: "handle_non-existent_log",
				payload: data.JSONMap{
					"id": "non-existent-log",
				},
				expectError: false, // Handler should handle gracefully
			},
			{
				name: "handle_empty_log_ID",
				payload: data.JSONMap{
					"id": "",
				},
				expectError: false, // Handler should handle gracefully
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Convert payload to expected format
				payloadBytes, _ := json.Marshal(tc.payload)
				handleErr := deps.AnalysisQueueHandler.Handle(ctx, nil, payloadBytes)
				if tc.expectError {
					assert.Error(t, handleErr)
				} else {
					assert.NoError(t, handleErr)
				}
			})
		}
	})
}

func (suite *QueueTestSuite) TestDeviceAnalysisQueueHandler_CreateSessionFromLog() {
	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		// Create a device first
		device := &models.Device{
			Name: "Test Device",
			OS:   "Linux",
		}
		device.GenID(ctx)
		err := deps.DeviceRepo.Create(ctx, device)
		require.NoError(t, err)

		// Create a device log with session data
		deviceLog := &models.DeviceLog{
			DeviceID: device.ID,
			Data: data.JSONMap{
				"userAgent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
				"ip":        "192.168.1.1",
				"tz":        "UTC",
				"lang":      "en-US",
				"cur":       "USD",
			},
		}
		deviceLog.GenID(ctx)

		// Test CreateSessionFromLog
		session, err := deps.AnalysisQueueHandler.CreateSessionFromLog(ctx, deviceLog)
		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, device.ID, session.DeviceID)
		assert.Equal(t, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36", session.UserAgent)
		assert.Equal(t, "192.168.1.1", session.IP)
		assert.NotEmpty(t, session.ID)
	})
}

func (suite *QueueTestSuite) TestDeviceAnalysisQueueHandler_CreateDeviceFromSess() {
	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		// Create a session
		session := &models.DeviceSession{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			IP:        "10.0.0.1",
			LastSeen:  time.Now(),
		}
		session.GenID(ctx)

		// Test CreateDeviceFromSess
		device, err := deps.AnalysisQueueHandler.CreateDeviceFromSess(ctx, session)
		require.NoError(t, err)
		assert.NotNil(t, device)
		assert.NotEmpty(t, device.ID)
		assert.NotEmpty(t, device.Name) // Should extract platform from user agent
		assert.NotEmpty(t, device.OS)   // Should extract OS from user agent
	})
}

func (suite *QueueTestSuite) TestDeviceAnalysisQueueHandler_ExtractLocaleData() {
	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		testCases := []struct {
			name    string
			data    data.JSONMap
			geoIP   *queue.GeoIP
			wantTz  string
			wantCur string
		}{
			{
				name: "extract from data map",
				data: data.JSONMap{
					"tz":    "America/New_York",
					"lang":  "en-US,en",
					"cur":   "USD",
					"curNm": "US Dollar",
					"code":  "+1",
				},
				geoIP:   nil,
				wantTz:  "America/New_York",
				wantCur: "USD",
			},
			{
				name: "extract from geoIP fallback",
				data: data.JSONMap{},
				geoIP: &queue.GeoIP{
					Timezone:           "Europe/London",
					Languages:          "en-GB,en",
					Currency:           "GBP",
					CurrencyName:       "British Pound",
					CountryCallingCode: "+44",
				},
				wantTz:  "Europe/London",
				wantCur: "GBP",
			},
			{
				name: "data overrides geoIP",
				data: data.JSONMap{
					"tz":  "Asia/Tokyo",
					"cur": "JPY",
				},
				geoIP: &queue.GeoIP{
					Timezone: "Europe/London",
					Currency: "GBP",
				},
				wantTz:  "Asia/Tokyo",
				wantCur: "JPY",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				locale, err := deps.AnalysisQueueHandler.ExtractLocaleData(ctx, tc.data, tc.geoIP)
				require.NoError(t, err)
				assert.NotNil(t, locale)
				assert.Equal(t, tc.wantTz, locale.GetTimezone())
				assert.Equal(t, tc.wantCur, locale.GetCurrency())
			})
		}
	})
}

func (suite *QueueTestSuite) TestDeviceAnalysisQueueHandler_ExtractLocationData() {
	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		testCases := []struct {
			name     string
			data     data.JSONMap
			geoIP    *queue.GeoIP
			wantLat  float64
			wantLong float64
		}{
			{
				name: "extract from data map",
				data: data.JSONMap{
					"lat":  40.7128,
					"long": -74.0060,
				},
				geoIP:    nil,
				wantLat:  40.7128,
				wantLong: -74.006,
			},
			{
				name: "extract from geoIP",
				data: data.JSONMap{},
				geoIP: &queue.GeoIP{
					Country:   "United States",
					Region:    "New York",
					City:      "New York",
					Latitude:  40.7128,
					Longitude: -74.0060,
				},
				wantLat:  40.712800,
				wantLong: -74.006,
			},
			{
				name: "data overrides geoIP",
				data: data.JSONMap{
					"lat":  35.6762,
					"long": 139.6503,
				},
				geoIP: &queue.GeoIP{
					Latitude:  40.7128,
					Longitude: -74.0060,
				},
				wantLat:  35.6762,
				wantLong: 139.6503,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				locationData := deps.AnalysisQueueHandler.ExtractLocationData(ctx, tc.data, tc.geoIP)
				assert.NotNil(t, locationData)

				assert.InDelta(t, tc.wantLat, locationData["latitude"], 0.0001)
				assert.InDelta(t, tc.wantLong, locationData["longitude"], 0.0001)
			})
		}
	})
}

func (suite *QueueTestSuite) TestQueryIPGeo() {
	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc, _ := suite.CreateService(t, dep)

		// Test QueryIPGeo function
		geoIP, err := queue.QueryIPGeo(ctx, svc.HTTPClientManager(), "8.8.8.8")

		// Note: This test may fail due to external API rate limiting
		// The important thing is that the function executes without panic
		if err != nil {
			t.Logf("QueryIPGeo failed (likely due to rate limiting): %v", err)
			return
		}

		if geoIP != nil {
			assert.Equal(t, "8.8.8.8", geoIP.IP)
			assert.NotEmpty(t, geoIP.Country)
		}
	})
}
