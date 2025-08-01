package business_test

import (
	"context"
	"testing"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/tests"
	"github.com/pitabwire/frame/tests/deps/testpostgres"
	"github.com/pitabwire/frame/tests/testdef"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/business"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

const (
	DefaultRandomStringLength = 8
)

type DeviceBusinessTestSuite struct {
	tests.FrameBaseTestSuite
}

func initResources(_ context.Context) []testdef.TestResource {
	pg := testpostgres.NewPGDepWithCred(testpostgres.PostgresqlDBImage, "ant", "s3cr3t", "service_profile")
	resources := []testdef.TestResource{pg}
	return resources
}

func (suite *DeviceBusinessTestSuite) SetupSuite() {
	suite.InitResourceFunc = initResources
	suite.FrameBaseTestSuite.SetupSuite()
}

func (suite *DeviceBusinessTestSuite) CreateService(
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

	svc.Init(ctx)

	err = repository.Migrate(ctx, svc, "../../migrations/0001")
	require.NoError(t, err)

	err = svc.Run(ctx, "")
	require.NoError(t, err)

	return svc, ctx
}

func (suite *DeviceBusinessTestSuite) TearDownSuite() {
	suite.FrameBaseTestSuite.TearDownSuite()
}

// WithTestDependancies Creates subtests with each known DependancyOption.
func (suite *DeviceBusinessTestSuite) WithTestDependancies(
	t *testing.T,
	testFn func(t *testing.T, dep *testdef.DependancyOption),
) {
	options := []*testdef.DependancyOption{
		testdef.NewDependancyOption("default", util.RandomString(DefaultRandomStringLength), suite.Resources()),
	}

	tests.WithTestDependancies(t, options, testFn)
}

func TestDeviceBusinessTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceBusinessTestSuite))
}

func (suite *DeviceBusinessTestSuite) TestSaveDevice() {
	testCases := []struct {
		name        string
		id          string
		deviceName  string
		data        map[string]string
		expectError bool
		errorMsg    string
	}{
		{
			name:       "valid device with all data",
			id:         "",
			deviceName: "Test Device",
			data: map[string]string{
				"profile_id": "profile-123",
				"os":         "Linux",
				"user_agent": "Mozilla/5.0",
				"ip":         "192.168.1.1",
				"locale":     "en-US",
				"location":   "US",
			},
			expectError: false,
		},
		{
			name:       "valid device with custom ID",
			id:         "custom-device-id",
			deviceName: "Custom Device",
			data: map[string]string{
				"profile_id": "profile-456",
				"os":         "Windows",
				"user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
				"ip":         "10.0.0.1",
				"locale":     "en-GB",
				"location":   "UK",
			},
			expectError: false,
		},
		{
			name:       "minimal device data",
			id:         "",
			deviceName: "Minimal Device",
			data: map[string]string{
				"profile_id": "profile-789",
			},
			expectError: false,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				device, err := biz.SaveDevice(ctx, tc.id, tc.deviceName, tc.data)

				if tc.expectError {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tc.errorMsg)
					assert.Nil(t, device)
				} else {
					require.NoError(t, err)
					assert.NotNil(t, device)
					assert.Equal(t, tc.deviceName, device.GetName())
					if tc.id != "" {
						assert.Equal(t, tc.id, device.GetId())
					} else {
						assert.NotEmpty(t, device.GetId())
					}
					if os, exists := tc.data["os"]; exists {
						assert.Equal(t, os, device.GetOs())
					}
				}
			})
		}
	})
}

func (suite *DeviceBusinessTestSuite) TestGetDeviceByID() {
	testCases := []struct {
		name        string
		setupDevice bool
		deviceID    string
		expectError bool
	}{
		{
			name:        "existing device",
			setupDevice: true,
			expectError: false,
		},
		{
			name:        "non-existent device",
			setupDevice: false,
			deviceID:    "non-existent-id",
			expectError: true,
		},
		{
			name:        "empty device ID",
			setupDevice: false,
			deviceID:    "",
			expectError: true,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID string
				if tc.setupDevice {
					// Create a device first
					device, err := biz.SaveDevice(ctx, "", "Test Device", map[string]string{
						"profile_id": "profile-test",
						"os":         "Linux",
						"user_agent": "Test Agent",
						"ip":         "127.0.0.1",
					})
					require.NoError(t, err)
					deviceID = device.GetId()
				} else {
					deviceID = tc.deviceID
				}

				device, err := biz.GetDeviceByID(ctx, deviceID)

				if tc.expectError {
					require.Error(t, err)
					assert.Nil(t, device)
				} else {
					require.NoError(t, err)
					assert.NotNil(t, device)
					assert.Equal(t, deviceID, device.GetId())
					assert.Equal(t, "Test Device", device.GetName())
				}
			})
		}
	})
}

func (suite *DeviceBusinessTestSuite) TestGetDeviceBySessionID() {
	testCases := []struct {
		name        string
		setupDevice bool
		sessionID   string
		expectError bool
	}{
		{
			name:        "existing session",
			setupDevice: true,
			expectError: false,
		},
		{
			name:        "non-existent session",
			setupDevice: false,
			sessionID:   "non-existent-session",
			expectError: true,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var sessionID string
				if tc.setupDevice {
					// Create a device first which creates a session
					device, err := biz.SaveDevice(ctx, "", "Test Device", map[string]string{
						"profile_id": "profile-test",
						"os":         "Linux",
						"user_agent": "Test Agent",
						"ip":         "127.0.0.1",
					})
					require.NoError(t, err)
					sessionID = device.GetSessionId()
				} else {
					sessionID = tc.sessionID
				}

				device, err := biz.GetDeviceBySessionID(ctx, sessionID)

				if tc.expectError {
					require.Error(t, err)
					assert.Nil(t, device)
				} else {
					require.NoError(t, err)
					assert.NotNil(t, device)
					assert.Equal(t, sessionID, device.GetSessionId())
				}
			})
		}
	})
}

func (suite *DeviceBusinessTestSuite) TestLogDeviceActivity() {
	testCases := []struct {
		name        string
		setupDevice bool
		deviceID    string
		sessionID   string
		data        map[string]string
		expectError bool
	}{
		{
			name:        "valid activity log",
			setupDevice: true,
			data: map[string]string{
				"action": "login",
				"result": "success",
			},
			expectError: false,
		},
		{
			name:        "empty data",
			setupDevice: true,
			data:        map[string]string{},
			expectError: false,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID, sessionID string
				if tc.setupDevice {
					device, err := biz.SaveDevice(ctx, "", "Test Device", map[string]string{
						"profile_id": "profile-log-test",
						"os":         "Linux",
					})
					require.NoError(t, err)
					deviceID = device.GetId()
					sessionID = device.GetSessionId()
				} else {
					deviceID = tc.deviceID
					sessionID = tc.sessionID
				}

				log, err := biz.LogDeviceActivity(ctx, deviceID, sessionID, tc.data)

				if tc.expectError {
					require.Error(t, err)
					assert.Nil(t, log)
				} else {
					require.NoError(t, err)
					assert.NotNil(t, log)
					assert.Equal(t, deviceID, log.GetDeviceId())
					assert.Equal(t, sessionID, log.GetSessionId())
				}
			})
		}
	})
}

func (suite *DeviceBusinessTestSuite) TestAddKey() {
	testCases := []struct {
		name        string
		setupDevice bool
		deviceID    string
		keyType     devicev1.KeyType
		key         []byte
		extra       map[string]string
		expectError bool
	}{
		{
			name:        "add valid key",
			setupDevice: true,
			keyType:     devicev1.KeyType_MATRIX_KEY,
			key:         []byte("test-encryption-key"),
			extra: map[string]string{
				"algorithm": "AES256",
			},
			expectError: false,
		},
		{
			name:        "add key with empty extra",
			setupDevice: true,
			keyType:     devicev1.KeyType_NOTIFICATION_KEY,
			key:         []byte("test-signing-key"),
			extra:       map[string]string{},
			expectError: false,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID string
				if tc.setupDevice {
					device, err := biz.SaveDevice(ctx, "", "Test Device", map[string]string{
						"profile_id": "profile-key-test",
						"os":         "Linux",
					})
					require.NoError(t, err)
					deviceID = device.GetId()
				} else {
					deviceID = tc.deviceID
				}

				keyObj, err := biz.AddKey(ctx, deviceID, tc.keyType, tc.key, tc.extra)

				if tc.expectError {
					require.Error(t, err)
					assert.Nil(t, keyObj)
				} else {
					require.NoError(t, err)
					assert.NotNil(t, keyObj)
					assert.Equal(t, deviceID, keyObj.GetDeviceId())
					assert.Equal(t, tc.key, keyObj.GetKey())
				}
			})
		}
	})
}

func (suite *DeviceBusinessTestSuite) TestRemoveDevice() {
	testCases := []struct {
		name        string
		setupDevice bool
		deviceID    string
		expectError bool
	}{
		{
			name:        "remove existing device",
			setupDevice: true,
			expectError: false,
		},
		{
			name:        "remove non-existent device",
			setupDevice: false,
			deviceID:    "non-existent-device",
			expectError: true,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID string
				if tc.setupDevice {
					device, err := biz.SaveDevice(ctx, "", "Test Device", map[string]string{
						"profile_id": "profile-remove-test",
						"os":         "Linux",
					})
					require.NoError(t, err)
					deviceID = device.GetId()
				} else {
					deviceID = tc.deviceID
				}

				err := biz.RemoveDevice(ctx, deviceID)

				if tc.expectError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)

					// Verify device is removed
					_, err = biz.GetDeviceByID(ctx, deviceID)
					require.Error(t, err)
				}
			})
		}
	})
}
