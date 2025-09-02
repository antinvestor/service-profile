package repository_test

import (
	"context"
	"testing"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/framedata"
	"github.com/pitabwire/frame/frametests"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/frame/frametests/deps/testpostgres"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

const (
	DefaultRandomStringLength = 8
)

type DeviceRepositoryTestSuite struct {
	frametests.FrameBaseTestSuite
}

func initResources(_ context.Context) []definition.TestResource {
	pg := testpostgres.NewWithOpts("service_devices", definition.WithUserName("ant"))
	resources := []definition.TestResource{pg}
	return resources
}

func (suite *DeviceRepositoryTestSuite) SetupSuite() {
	suite.InitResourceFunc = initResources
	suite.FrameBaseTestSuite.SetupSuite()
}

func (suite *DeviceRepositoryTestSuite) CreateService(
	t *testing.T,
	depOpts *definition.DependancyOption,
) (*frame.Service, context.Context) {
	ctx := t.Context()
	t.Setenv("OTEL_TRACES_EXPORTER", "none")
	deviceConfig, err := frame.ConfigFromEnv[config.DevicesConfig]()
	require.NoError(t, err)

	deviceConfig.LogLevel = "debug"
	deviceConfig.RunServiceSecurely = false
	deviceConfig.ServerPort = ""

	res := depOpts.ByIsDatabase(ctx)
	testDS, cleanup, err0 := res.GetRandomisedDS(t.Context(), depOpts.Prefix())
	require.NoError(t, err0)

	t.Cleanup(func() {
		cleanup(t.Context())
	})

	deviceConfig.DatabasePrimaryURL = []string{testDS.String()}
	deviceConfig.DatabaseReplicaURL = []string{testDS.String()}

	ctx, svc := frame.NewServiceWithContext(t.Context(), "device tests",
		frame.WithConfig(&deviceConfig),
		frame.WithDatastore(),
		frametests.WithNoopDriver())

	svc.Init(ctx)

	err = repository.Migrate(ctx, svc, "../../migrations/0001")
	require.NoError(t, err)

	err = svc.Run(ctx, "")
	require.NoError(t, err)

	return svc, ctx
}

func (suite *DeviceRepositoryTestSuite) TearDownSuite() {
	suite.FrameBaseTestSuite.TearDownSuite()
}

// WithTestDependancies Creates subtests with each known DependancyOption.
func (suite *DeviceRepositoryTestSuite) WithTestDependancies(
	t *testing.T,
	testFn func(t *testing.T, dep *definition.DependancyOption),
) {
	options := []*definition.DependancyOption{
		definition.NewDependancyOption("default", util.RandomString(DefaultRandomStringLength), suite.Resources()),
	}

	frametests.WithTestDependancies(t, options, testFn)
}

func TestDeviceRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceRepositoryTestSuite))
}

func (suite *DeviceRepositoryTestSuite) TestDeviceRepository() {
	testCases := []struct {
		name        string
		profileID   string
		deviceName  string
		os          string
		expectError bool
	}{
		{
			name:        "save and retrieve device",
			profileID:   "profile-123",
			deviceName:  "Test Device",
			os:          "Linux",
			expectError: false,
		},
		{
			name:        "save device with empty name",
			profileID:   "profile-456",
			deviceName:  "",
			os:          "Windows",
			expectError: false,
		},
		{
			name:        "save device with long name",
			profileID:   "profile-789",
			deviceName:  util.RandomString(200),
			os:          "macOS",
			expectError: false,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		deviceRepo := repository.NewDeviceRepository(svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create device
				device := &models.Device{
					ProfileID: tc.profileID,
					Name:      tc.deviceName,
					OS:        tc.os,
				}

				err := deviceRepo.Save(ctx, device)
				if tc.expectError {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.NotEmpty(t, device.GetID())

				// Retrieve by ID
				retrievedDevice, err := deviceRepo.GetByID(ctx, device.GetID())
				require.NoError(t, err)
				assert.Equal(t, device.GetID(), retrievedDevice.GetID())
				assert.Equal(t, tc.profileID, retrievedDevice.ProfileID)
				assert.Equal(t, tc.deviceName, retrievedDevice.Name)
				assert.Equal(t, tc.os, retrievedDevice.OS)

				// Retrieve by ProfileID
				searchProperties := frame.JSONMap{
					"profile_id": tc.profileID,
				}
				q := framedata.NewSearchQuery("", searchProperties, 0, 50)
				devicesResult, err := deviceRepo.Search(ctx, q)
				require.NoError(t, err)

				// Read results from the pipe
				var devices []*models.Device
				for {
					res, ok := devicesResult.ReadResult(ctx)
					if !ok {
						break
					}
					if res.IsError() {
						require.NoError(t, res.Error())
					}
					devices = append(devices, res.Item()...)
				}
				assert.Len(t, devices, 1)
				assert.Equal(t, device.GetID(), devices[0].GetID())

				// Remove device
				removedDevice, err := deviceRepo.RemoveByID(ctx, device.GetID())
				require.NoError(t, err)
				assert.Equal(t, device.GetID(), removedDevice.GetID())

				// Verify removal
				_, err = deviceRepo.GetByID(ctx, device.GetID())
				assert.Error(t, err)
			})
		}
	})
}

func (suite *DeviceRepositoryTestSuite) TestDeviceSessionRepository() {
	testCases := []struct {
		name        string
		deviceID    string
		userAgent   string
		ip          string
		expectError bool
	}{
		{
			name:        "save and retrieve session",
			deviceID:    "device-123",
			userAgent:   "Mozilla/5.0",
			ip:          "192.168.1.1",
			expectError: false,
		},
		{
			name:        "save session with empty user agent",
			deviceID:    "device-456",
			userAgent:   "",
			ip:          "10.0.0.1",
			expectError: false,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		sessionRepo := repository.NewDeviceSessionRepository(svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create session
				session := &models.DeviceSession{
					DeviceID:  tc.deviceID,
					UserAgent: tc.userAgent,
					IP:        tc.ip,
				}

				err := sessionRepo.Save(ctx, session)
				if tc.expectError {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.NotEmpty(t, session.GetID())

				// Retrieve by ID
				retrievedSession, err := sessionRepo.GetByID(ctx, session.GetID())
				require.NoError(t, err)
				assert.Equal(t, session.GetID(), retrievedSession.GetID())
				assert.Equal(t, tc.deviceID, retrievedSession.DeviceID)
				assert.Equal(t, tc.userAgent, retrievedSession.UserAgent)
				assert.Equal(t, tc.ip, retrievedSession.IP)

				// Retrieve last by device ID
				lastSession, err := sessionRepo.GetLastByDeviceID(ctx, tc.deviceID)
				require.NoError(t, err)
				assert.Equal(t, session.GetID(), lastSession.GetID())
			})
		}
	})
}

func (suite *DeviceRepositoryTestSuite) TestDeviceLogRepository() {
	testCases := []struct {
		name        string
		deviceID    string
		sessionID   string
		expectError bool
	}{
		{
			name:        "save and retrieve log",
			deviceID:    "device-123",
			sessionID:   "session-123",
			expectError: false,
		},
		{
			name:        "save log with empty session ID",
			deviceID:    "device-456",
			sessionID:   "",
			expectError: false,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		logRepo := repository.NewDeviceLogRepository(svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create log
				log := &models.DeviceLog{
					DeviceID:        tc.deviceID,
					DeviceSessionID: tc.sessionID,
					Data:            frame.JSONMap{"action": "test"},
				}

				err := logRepo.Save(ctx, log)
				if tc.expectError {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.NotEmpty(t, log.GetID())

				// Retrieve by ID
				retrievedLog, err := logRepo.GetByID(ctx, log.GetID())
				require.NoError(t, err)
				assert.Equal(t, log.GetID(), retrievedLog.GetID())
				assert.Equal(t, tc.deviceID, retrievedLog.DeviceID)
				assert.Equal(t, tc.sessionID, retrievedLog.DeviceSessionID)

				// Retrieve by device ID
				searchProperties := frame.JSONMap{
					"device_id": tc.deviceID,
				}
				q := framedata.NewSearchQuery("", searchProperties, 0, 50)
				logsResult, err := logRepo.GetByDeviceID(ctx, q)
				require.NoError(t, err)

				// Read results from the pipe
				var logs []*models.DeviceLog
				for {
					res, ok := logsResult.ReadResult(ctx)
					if !ok {
						break
					}
					if res.IsError() {
						require.NoError(t, res.Error())
					}
					logs = append(logs, res.Item()...)
				}
				assert.Len(t, logs, 1)
				assert.Equal(t, log.GetID(), logs[0].GetID())
			})
		}
	})
}

func (suite *DeviceRepositoryTestSuite) TestDeviceKeyRepository() {
	testCases := []struct {
		name        string
		deviceID    string
		key         []byte
		expectError bool
	}{
		{
			name:        "save and retrieve key",
			deviceID:    "device-123",
			key:         []byte("test-encryption-key"),
			expectError: false,
		},
		{
			name:        "save key with empty data",
			deviceID:    "device-456",
			key:         []byte{},
			expectError: false,
		},
		{
			name:        "save large key",
			deviceID:    "device-789",
			key:         []byte(util.RandomString(1000)),
			expectError: false,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		keyRepo := repository.NewDeviceKeyRepository(svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create key
				key := &models.DeviceKey{
					DeviceID: tc.deviceID,
					Key:      tc.key,
					Extra:    frame.JSONMap{"type": "test"},
				}

				err := keyRepo.Save(ctx, key)
				if tc.expectError {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.NotEmpty(t, key.GetID())

				// Retrieve by device ID
				keys, err := keyRepo.GetByDeviceID(ctx, tc.deviceID)
				require.NoError(t, err)
				assert.Len(t, keys, 1)
				assert.Equal(t, key.GetID(), keys[0].GetID())
				assert.Equal(t, tc.deviceID, keys[0].DeviceID)
				assert.Equal(t, tc.key, keys[0].Key)

				// Remove key
				removedKey, err := keyRepo.RemoveByID(ctx, key.GetID())
				require.NoError(t, err)
				assert.Equal(t, key.GetID(), removedKey.GetID())

				// Verify removal
				keys, err = keyRepo.GetByDeviceID(ctx, tc.deviceID)
				require.NoError(t, err)
				assert.Empty(t, keys)
			})
		}
	})
}
