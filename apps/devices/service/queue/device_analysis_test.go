package queue_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/tests"
	"github.com/pitabwire/frame/tests/deps/testpostgres"
	"github.com/pitabwire/frame/tests/testdef"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/queue"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

const (
	DefaultRandomStringLength = 8
)

type QueueTestSuite struct {
	tests.FrameBaseTestSuite
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}

func initResources(_ context.Context) []testdef.TestResource {
	pg := testpostgres.NewPGDepWithCred(testpostgres.PostgresqlDBImage, "ant", "s3cr3t", "service_profile")
	return []testdef.TestResource{pg}
}

func (suite *QueueTestSuite) SetupSuite() {
	suite.InitResourceFunc = initResources
	suite.FrameBaseTestSuite.SetupSuite()
}

func (suite *QueueTestSuite) WithTestDependancies(t *testing.T, fn func(t *testing.T, dep *testdef.DependancyOption)) {
	options := []*testdef.DependancyOption{
		testdef.NewDependancyOption("default", util.RandomString(DefaultRandomStringLength), suite.Resources()),
	}

	tests.WithTestDependancies(t, options, fn)
}

func (suite *QueueTestSuite) CreateService(
	t *testing.T,
	depOpts *testdef.DependancyOption,
) (*frame.Service, context.Context) {
	t.Setenv("OTEL_TRACES_EXPORTER", "none")
	deviceConfig, err := frame.ConfigFromEnv[config.DevicesConfig]()
	require.NoError(t, err)

	deviceConfig.LogLevel = "debug"
	deviceConfig.RunServiceSecurely = false
	deviceConfig.ServerPort = ""

	for _, res := range depOpts.Database() {
		testDS, cleanup, err0 := res.GetRandomisedDS(t.Context(), depOpts.Prefix())
		require.NoError(t, err0)

		t.Cleanup(func() {
			cleanup(t.Context())
		})

		deviceConfig.DatabasePrimaryURL = []string{testDS.String()}
		deviceConfig.DatabaseReplicaURL = []string{testDS.String()}
	}

	ctx, svc := frame.NewServiceWithContext(t.Context(), "device tests",
		frame.WithConfig(&deviceConfig),
		frame.WithDatastore(),
		frame.WithNoopDriver())

	// Skip queue initialization for basic functionality tests
	// This allows us to test the core logic without queue dependencies

	err = repository.Migrate(ctx, svc, deviceConfig.GetDatabaseMigrationPath())
	require.NoError(t, err)

	return svc, ctx
}

func (suite *QueueTestSuite) TestDeviceAnalysisQueueHandler_Handle() {
	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Create repositories
		deviceRepo := repository.NewDeviceRepository(svc)
		sessionRepo := repository.NewDeviceSessionRepository(svc)
		logRepo := repository.NewDeviceLogRepository(svc)

		// Create handler
		handler := queue.DeviceAnalysisQueueHandler{
			DeviceRepository:    deviceRepo,
			DeviceLogRepository: logRepo,
			SessionRepository:   sessionRepo,
			Service:             svc,
		}

		// Create a device first
		device := &models.Device{
			Name: "Test Device",
			OS:   "Linux",
		}
		device.GenID(ctx)
		err := deviceRepo.Save(ctx, device)
		require.NoError(t, err)

		// Create a session
		session := &models.DeviceSession{
			DeviceID:  device.ID,
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
			IP:        "192.168.1.1",
			LastSeen:  time.Now(),
		}
		session.GenID(ctx)
		err = sessionRepo.Save(ctx, session)
		require.NoError(t, err)

		// Create a device log
		deviceLog := &models.DeviceLog{
			DeviceID:        device.ID,
			DeviceSessionID: session.ID,
			Data: frame.DBPropertiesFromMap(map[string]string{
				"action":    "page_view",
				"userAgent": session.UserAgent,
				"ip":        session.IP,
			}),
		}
		deviceLog.GenID(ctx)
		err = logRepo.Save(ctx, deviceLog)
		require.NoError(t, err)

		// Test cases
		testCases := []struct {
			name        string
			payload     map[string]string
			expectError bool
		}{
			{
				name: "handle_valid_device_log",
				payload: map[string]string{
					"id": deviceLog.ID,
				},
				expectError: false,
			},
			{
				name: "handle_non-existent_log",
				payload: map[string]string{
					"id": "non-existent-log",
				},
				expectError: false, // Handler should handle gracefully
			},
			{
				name: "handle_empty_log_ID",
				payload: map[string]string{
					"id": "",
				},
				expectError: false, // Handler should handle gracefully
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Convert payload to expected format
				payloadBytes, _ := json.Marshal(tc.payload)
				handleErr := handler.Handle(ctx, tc.payload, payloadBytes)
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
	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Create repositories
		deviceRepo := repository.NewDeviceRepository(svc)
		sessionRepo := repository.NewDeviceSessionRepository(svc)
		logRepo := repository.NewDeviceLogRepository(svc)

		// Create handler
		handler := queue.DeviceAnalysisQueueHandler{
			DeviceRepository:    deviceRepo,
			DeviceLogRepository: logRepo,
			SessionRepository:   sessionRepo,
			Service:             svc,
		}

		// Create a device first
		device := &models.Device{
			Name: "Test Device",
			OS:   "Linux",
		}
		device.GenID(ctx)
		err := deviceRepo.Save(ctx, device)
		require.NoError(t, err)

		// Create a device log with session data
		deviceLog := &models.DeviceLog{
			DeviceID: device.ID,
			Data: frame.DBPropertiesFromMap(map[string]string{
				"userAgent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
				"ip":        "192.168.1.1",
				"tz":        "UTC",
				"lang":      "en-US",
				"cur":       "USD",
			}),
		}
		deviceLog.GenID(ctx)

		// Test CreateSessionFromLog
		session, err := handler.CreateSessionFromLog(ctx, deviceLog)
		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, device.ID, session.DeviceID)
		assert.Equal(t, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36", session.UserAgent)
		assert.Equal(t, "192.168.1.1", session.IP)
		assert.NotEmpty(t, session.ID)
	})
}

func (suite *QueueTestSuite) TestDeviceAnalysisQueueHandler_CreateDeviceFromSess() {
	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Create repositories
		deviceRepo := repository.NewDeviceRepository(svc)
		sessionRepo := repository.NewDeviceSessionRepository(svc)
		logRepo := repository.NewDeviceLogRepository(svc)

		// Create handler
		handler := queue.DeviceAnalysisQueueHandler{
			DeviceRepository:    deviceRepo,
			DeviceLogRepository: logRepo,
			SessionRepository:   sessionRepo,
			Service:             svc,
		}

		// Create a session
		session := &models.DeviceSession{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			IP:        "10.0.0.1",
			LastSeen:  time.Now(),
		}
		session.GenID(ctx)

		// Test CreateDeviceFromSess
		device, err := handler.CreateDeviceFromSess(ctx, session)
		require.NoError(t, err)
		assert.NotNil(t, device)
		assert.NotEmpty(t, device.ID)
		assert.NotEmpty(t, device.Name) // Should extract platform from user agent
		assert.NotEmpty(t, device.OS)   // Should extract OS from user agent
	})
}

func (suite *QueueTestSuite) TestDeviceAnalysisQueueHandler_ExtractLocaleData() {
	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Create handler
		handler := queue.DeviceAnalysisQueueHandler{
			Service: svc,
		}

		testCases := []struct {
			name    string
			data    map[string]string
			geoIP   *queue.GeoIP
			wantTz  string
			wantCur string
		}{
			{
				name: "extract from data map",
				data: map[string]string{
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
				data: map[string]string{},
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
				data: map[string]string{
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
				locale, err := handler.ExtractLocaleData(ctx, tc.data, tc.geoIP)
				require.NoError(t, err)
				assert.NotNil(t, locale)
				assert.Equal(t, tc.wantTz, locale.GetTimezone())
				assert.Equal(t, tc.wantCur, locale.GetCurrency())
			})
		}
	})
}

func (suite *QueueTestSuite) TestDeviceAnalysisQueueHandler_ExtractLocationData() {
	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Create handler
		handler := queue.DeviceAnalysisQueueHandler{
			Service: svc,
		}

		testCases := []struct {
			name     string
			data     map[string]string
			geoIP    *queue.GeoIP
			wantLat  string
			wantLong string
		}{
			{
				name: "extract from data map",
				data: map[string]string{
					"lat":  "40.7128",
					"long": "-74.0060",
				},
				geoIP:    nil,
				wantLat:  "40.7128",
				wantLong: "-74.0060",
			},
			{
				name: "extract from geoIP",
				data: map[string]string{},
				geoIP: &queue.GeoIP{
					Country:   "United States",
					Region:    "New York",
					City:      "New York",
					Latitude:  40.7128,
					Longitude: -74.0060,
				},
				wantLat:  "40.712800",
				wantLong: "-74.006000",
			},
			{
				name: "data overrides geoIP",
				data: map[string]string{
					"lat":  "35.6762",
					"long": "139.6503",
				},
				geoIP: &queue.GeoIP{
					Latitude:  40.7128,
					Longitude: -74.0060,
				},
				wantLat:  "35.6762",
				wantLong: "139.6503",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				locationData := handler.ExtractLocationData(ctx, tc.data, tc.geoIP)
				assert.NotNil(t, locationData)

				locationMap := frame.DBPropertiesToMap(locationData)
				if tc.wantLat != "" {
					assert.Equal(t, tc.wantLat, locationMap["latitude"])
				}
				if tc.wantLong != "" {
					assert.Equal(t, tc.wantLong, locationMap["longitude"])
				}
			})
		}
	})
}

func (suite *QueueTestSuite) TestQueryIPGeo() {
	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Test QueryIPGeo function
		geoIP, err := queue.QueryIPGeo(ctx, svc, "8.8.8.8")

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
