package business_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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
	"github.com/antinvestor/service-profile/apps/devices/service/models"
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
	return []testdef.TestResource{pg}
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

func (suite *DeviceBusinessTestSuite) createTestDeviceWithSession(
	ctx context.Context,
	svc *frame.Service,
	deviceID, sessionID string,
) error {
	device := &models.Device{
		Name: "Original Name",
		OS:   "Original OS",
	}
	device.GenID(ctx)
	device.ID = deviceID
	err := repository.NewDeviceRepository(svc).Save(ctx, device)
	if err != nil {
		return err
	}

	session := &models.DeviceSession{
		DeviceID:  deviceID,
		UserAgent: "Test Agent",
		IP:        "127.0.0.1",
	}
	session.GenID(ctx)
	if sessionID != "" {
		session.ID = sessionID
	}
	return repository.NewDeviceSessionRepository(svc).Save(ctx, session)
}

func (suite *DeviceBusinessTestSuite) verifyDeviceActivityLogged(
	ctx context.Context,
	biz business.DeviceBusiness,
	deviceID, sessionID string,
) error {
	if sessionID == "" || deviceID == "" {
		return nil
	}

	deviceLogsChan, err := biz.GetDeviceLogs(ctx, deviceID)
	if err != nil {
		return err
	}

	// Process the channel to get actual logs
	var deviceLogs []*devicev1.DeviceLog
	for result := range deviceLogsChan {
		if result.IsError() {
			return result.Error()
		}
		deviceLogs = append(deviceLogs, result.Item()...)
	}

	logCount := len(deviceLogs)
	if logCount == 0 {
		return errors.New("expected device activity to be logged but found no logs")
	}
	return nil
}

func (suite *DeviceBusinessTestSuite) runSaveDeviceTestCase(
	ctx context.Context,
	t *testing.T,
	svc *frame.Service,
	biz business.DeviceBusiness,
	tc struct {
		name        string
		id          string
		deviceName  string
		data        map[string]string
		expectError bool
		expectNil   bool
	},
) {
	// Setup existing device if needed
	if tc.id != "" && tc.name == "save device with existing ID" {
		sessionID := tc.data["session_id"]
		err := suite.createTestDeviceWithSession(ctx, svc, tc.id, sessionID)
		require.NoError(t, err)
	}

	// Execute SaveDevice
	result, err := biz.SaveDevice(ctx, tc.id, tc.deviceName, tc.data)

	// Verify error expectations
	if tc.expectError {
		require.Error(t, err)
		return
	}
	require.NoError(t, err)

	// Verify result expectations
	if tc.expectNil {
		assert.Nil(t, result)
	} else {
		assert.NotNil(t, result)
		assert.Equal(t, tc.id, result.GetId())
		assert.Equal(t, tc.deviceName, result.GetName())
	}

	// Verify activity logging only for successful cases with valid device ID
	sessionID := tc.data["session_id"]
	if !tc.expectError && tc.id != "" && sessionID != "" {
		err = suite.verifyDeviceActivityLogged(ctx, biz, tc.id, sessionID)
		assert.NoError(t, err)
	}
}

func (suite *DeviceBusinessTestSuite) TestSaveDevice() {
	t := suite.T()
	testCases := []struct {
		name        string
		id          string
		deviceName  string
		data        map[string]string
		expectError bool
		expectNil   bool
	}{
		{
			name:       "save device with existing ID",
			id:         "existing-device-id",
			deviceName: "Test Device",
			data: map[string]string{
				"profile_id": "profile-123",
				"os":         "Linux",
				"user_agent": "Mozilla/5.0",
				"ip":         "192.168.1.1",
				"session_id": "test-session-1",
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name:       "save device with empty ID returns error",
			id:         "",
			deviceName: "Test Device",
			data: map[string]string{
				"profile_id": "profile-456",
				"os":         "Windows",
				"session_id": "test-session-2",
			},
			expectError: true,
			expectNil:   false,
		},
		{
			name:       "save device with empty ID and session logs activity but returns error",
			id:         "",
			deviceName: "Minimal Device",
			data: map[string]string{
				"profile_id": "profile-789",
				"session_id": "test-session-3",
			},
			expectError: true,
			expectNil:   false,
		},
	}

	suite.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				suite.runSaveDeviceTestCase(ctx, t, svc, biz, tc)
			})
		}
	})
}

func (suite *DeviceBusinessTestSuite) TestGetDeviceByID() {
	t := suite.T()
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

	suite.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID string
				if tc.setupDevice {
					// Create a device and session first
					device := &models.Device{
						Name: "Test Device",
						OS:   "Linux",
					}
					device.GenID(ctx)
					err := repository.NewDeviceRepository(svc).Save(ctx, device)
					require.NoError(t, err)
					deviceID = device.GetID()

					// Create a session for the device (required by GetDeviceByID)
					session := &models.DeviceSession{
						DeviceID:  deviceID,
						UserAgent: "Test Agent",
						IP:        "127.0.0.1",
					}
					session.GenID(ctx)
					err = repository.NewDeviceSessionRepository(svc).Save(ctx, session)
					require.NoError(t, err)
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
	t := suite.T()
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

	suite.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var sessionID string
				if tc.setupDevice {
					// Create a device and session first
					device := &models.Device{
						Name: "Test Device",
						OS:   "Linux",
					}
					device.GenID(ctx)
					err := repository.NewDeviceRepository(svc).Save(ctx, device)
					require.NoError(t, err)

					// Create a session for the device
					session := &models.DeviceSession{
						DeviceID:  device.GetID(),
						UserAgent: "Test Agent",
						IP:        "127.0.0.1",
					}
					session.GenID(ctx)
					err = repository.NewDeviceSessionRepository(svc).Save(ctx, session)
					require.NoError(t, err)
					sessionID = session.GetID()
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
	t := suite.T()
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

	suite.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID, sessionID string
				if tc.setupDevice {
					// Create a device directly using repository
					device := &models.Device{
						Name: "Test Device",
						OS:   "Linux",
					}
					device.GenID(ctx)
					err := repository.NewDeviceRepository(svc).Save(ctx, device)
					require.NoError(t, err)
					deviceID = device.GetID()

					// Create a session for the device
					session := &models.DeviceSession{
						DeviceID:  deviceID,
						UserAgent: "Test Agent",
						IP:        "127.0.0.1",
					}
					session.GenID(ctx)
					err = repository.NewDeviceSessionRepository(svc).Save(ctx, session)
					require.NoError(t, err)
					sessionID = session.GetID()
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
	t := suite.T()
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

	suite.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID string
				if tc.setupDevice {
					// Create a device directly using repository
					device := &models.Device{
						Name: "Test Device",
						OS:   "Linux",
					}
					device.GenID(ctx)
					err := repository.NewDeviceRepository(svc).Save(ctx, device)
					require.NoError(t, err)
					deviceID = device.GetID()
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
	t := suite.T()
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

	suite.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID string
				if tc.setupDevice {
					// Create a device first
					device := &models.Device{
						Name: "Test Device",
						OS:   "Linux",
					}
					device.GenID(ctx)
					err := repository.NewDeviceRepository(svc).Save(ctx, device)
					require.NoError(t, err)
					deviceID = device.GetID()
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

func (suite *DeviceBusinessTestSuite) runSearchDevicesTestCase(
	ctx context.Context,
	t *testing.T,
	svc *frame.Service,
	biz business.DeviceBusiness,
	tc struct {
		name        string
		setupDevice bool
		profileID   string
		expectError bool
	},
) {
	if tc.setupDevice {
		err := suite.createDeviceWithProfile(ctx, svc, tc.profileID)
		require.NoError(t, err)
	}

	query := &devicev1.SearchRequest{
		Query: tc.profileID,
	}

	devicesChan, err := biz.SearchDevices(ctx, query)
	if tc.expectError {
		require.Error(t, err)
		return
	}

	require.NoError(t, err)
	assert.NotNil(t, devicesChan)

	devices, channelErr := suite.processSearchResults(devicesChan)
	require.NoError(t, channelErr)

	suite.verifySearchResults(t, tc, devices)
}

func (suite *DeviceBusinessTestSuite) verifySearchResults(t *testing.T, tc struct {
	name        string
	setupDevice bool
	profileID   string
	expectError bool
}, devices []*devicev1.DeviceObject) {
	if tc.setupDevice && tc.profileID != "" {
		assert.Len(t, devices, 1)
		if len(devices) > 0 {
			assert.NotEmpty(t, devices[0].GetId())
		}
	} else {
		assert.Empty(t, devices)
	}
}

func (suite *DeviceBusinessTestSuite) TestSearchDevices() {
	t := suite.T()
	testCases := []struct {
		name        string
		setupDevice bool
		profileID   string
		expectError bool
	}{
		{
			name:        "search with valid profile ID",
			setupDevice: true,
			profileID:   "profile-search-test",
			expectError: false,
		},
		{
			name:        "search with non-existent profile ID",
			setupDevice: false,
			profileID:   "non-existent-profile",
			expectError: false,
		},
		{
			name:        "search with empty profile ID",
			setupDevice: false,
			profileID:   "",
			expectError: false,
		},
	}

	suite.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				suite.runSearchDevicesTestCase(ctx, t, svc, biz, tc)
			})
		}
	})
}

func (suite *DeviceBusinessTestSuite) TestGetDeviceLogs() {
	t := suite.T()
	testCases := []struct {
		name        string
		setupDevice bool
		deviceID    string
		expectError bool
	}{
		{
			name:        "get logs for device with logs",
			setupDevice: true,
			deviceID:    "",
			expectError: false,
		},
		{
			name:        "get logs for device without logs",
			setupDevice: true,
			deviceID:    "",
			expectError: false,
		},
		{
			name:        "get logs for non-existent device",
			setupDevice: false,
			deviceID:    "non-existent-device",
			expectError: false,
		},
	}

	suite.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID string
				if tc.setupDevice {
					// Create a device first
					var err error
					deviceID, err = suite.createDeviceForLogs(ctx, svc, tc.name == "get logs for device with logs", biz)
					require.NoError(t, err)
				} else {
					deviceID = tc.deviceID
				}

				logsChan, logErr := biz.GetDeviceLogs(ctx, deviceID)
				require.NoError(t, logErr)
				assert.NotNil(t, logsChan)

				// Collect results from channel
				logs, err := suite.processDeviceLogsResults(logsChan)
				require.NoError(t, err)

				if tc.name == "get logs for device with logs" {
					assert.Len(t, logs, 1)
					assert.Equal(t, deviceID, logs[0].GetDeviceId())
				} else {
					assert.Empty(t, logs)
				}
			})
		}
	})
}

func (suite *DeviceBusinessTestSuite) TestGetKeys() {
	t := suite.T()
	testCases := []struct {
		name        string
		setupDevice bool
		setupKey    bool
		deviceID    string
		keyType     devicev1.KeyType
		expectError bool
	}{
		{
			name:        "get keys for device with keys",
			setupDevice: true,
			setupKey:    true,
			keyType:     devicev1.KeyType_MATRIX_KEY,
			expectError: false,
		},
		{
			name:        "get keys for device without keys",
			setupDevice: true,
			setupKey:    false,
			keyType:     devicev1.KeyType_NOTIFICATION_KEY,
			expectError: false,
		},
		{
			name:        "get keys for non-existent device",
			setupDevice: false,
			setupKey:    false,
			deviceID:    "non-existent-device",
			keyType:     devicev1.KeyType_MATRIX_KEY,
			expectError: false,
		},
	}

	suite.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID string
				if tc.setupDevice {
					// Create a device first
					var err error
					deviceID, err = suite.createDeviceWithKey(ctx, svc, biz, tc.setupKey, tc.keyType)
					require.NoError(t, err)
				} else {
					deviceID = tc.deviceID
				}

				keysChan, err := biz.GetKeys(ctx, deviceID, tc.keyType)
				if tc.expectError {
					require.Error(t, err)
					return
				}

				require.NoError(t, err)
				assert.NotNil(t, keysChan)

				// Collect results from channel
				keys, channelErr := suite.processKeysResults(keysChan)
				require.NoError(t, channelErr)

				if tc.setupKey {
					assert.Len(t, keys, 1)
					assert.Equal(t, deviceID, keys[0].GetDeviceId())
					assert.Equal(t, []byte("test-key-data"), keys[0].GetKey())
				} else {
					assert.Empty(t, keys)
				}
			})
		}
	})
}

func (suite *DeviceBusinessTestSuite) TestRemoveKeys() {
	t := suite.T()
	testCases := []struct {
		name        string
		setupDevice bool
		setupKeys   int
		removeCount int
		expectError bool
	}{
		{
			name:        "remove existing keys",
			setupDevice: true,
			setupKeys:   2,
			removeCount: 2,
			expectError: false,
		},
		{
			name:        "remove non-existent keys",
			setupDevice: false,
			setupKeys:   0,
			removeCount: 1,
			expectError: true,
		},
		{
			name:        "remove subset of keys",
			setupDevice: true,
			setupKeys:   3,
			removeCount: 2,
			expectError: false,
		},
	}

	suite.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var keysToRemove []string
				if tc.setupDevice {
					// Create a device first
					_, keyIDs, err := suite.createDeviceWithKeys(ctx, svc, biz, tc.setupKeys)
					require.NoError(t, err)

					// Select keys to remove
					for i := 0; i < tc.removeCount && i < len(keyIDs); i++ {
						keysToRemove = append(keysToRemove, keyIDs[i])
					}
				} else {
					// Use non-existent key IDs
					keysToRemove = suite.generateNonExistentKeyIDs(tc.removeCount)
				}

				keysChan, err := biz.RemoveKeys(ctx, keysToRemove...)
				require.NoError(t, err)
				assert.NotNil(t, keysChan)

				// Process results from channel
				removedKeys, channelError := suite.processRemoveKeysResults(keysChan)

				if tc.expectError {
					require.Error(t, channelError)
				} else {
					require.NoError(t, channelError)
					assert.Len(t, removedKeys, tc.removeCount)
				}
			})
		}
	})
}

func (suite *DeviceBusinessTestSuite) createDeviceWithKeys(
	ctx context.Context,
	svc *frame.Service,
	biz business.DeviceBusiness,
	keyCount int,
) (string, []string, error) {
	// Create a device directly using repository
	device := &models.Device{
		Name: "Remove Key Test Device",
		OS:   "Linux",
	}
	device.GenID(ctx)
	err := repository.NewDeviceRepository(svc).Save(ctx, device)
	if err != nil {
		return "", nil, err
	}
	deviceID := device.GetID()

	var keyIDs []string
	for i := range keyCount {
		keyResult, keyErr := biz.AddKey(ctx, deviceID, devicev1.KeyType_MATRIX_KEY,
			[]byte(fmt.Sprintf("test-key-data-%d", i)), map[string]string{
				"index": strconv.Itoa(i),
			})
		if keyErr != nil {
			return "", nil, keyErr
		}
		keyIDs = append(keyIDs, keyResult.GetId())
	}
	return deviceID, keyIDs, nil
}

func (suite *DeviceBusinessTestSuite) generateNonExistentKeyIDs(count int) []string {
	var keyIDs []string
	for i := range count {
		keyIDs = append(keyIDs, fmt.Sprintf("non-existent-key-%d", i))
	}
	return keyIDs
}

func (suite *DeviceBusinessTestSuite) processRemoveKeysResults(
	keysChan <-chan frame.JobResult[[]*devicev1.KeyObject],
) ([]*devicev1.KeyObject, error) {
	var removedKeys []*devicev1.KeyObject
	for result := range keysChan {
		if result.IsError() {
			return nil, result.Error()
		}
		removedKeys = append(removedKeys, result.Item()...)
	}
	return removedKeys, nil
}

func (suite *DeviceBusinessTestSuite) TestLinkDeviceToProfile() {
	t := suite.T()
	testCases := []struct {
		name        string
		setupDevice bool
		sessionID   string
		profileID   string
		expectError bool
	}{
		{
			name:        "link device to profile successfully",
			setupDevice: true,
			profileID:   "new-profile-123",
			expectError: false,
		},
		{
			name:        "link device that already has profile",
			setupDevice: true,
			profileID:   "another-profile-456",
			expectError: false,
		},
		{
			name:        "link with non-existent session",
			setupDevice: false,
			sessionID:   "non-existent-session",
			profileID:   "profile-789",
			expectError: true,
		},
	}

	suite.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		biz := business.NewDeviceBusiness(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var sessionID string
				if tc.setupDevice {
					// Create a device and session first
					device := &models.Device{
						ProfileID: "original-profile",
						Name:      "Link Test Device",
						OS:        "Linux",
					}
					device.GenID(ctx)
					err := repository.NewDeviceRepository(svc).Save(ctx, device)
					require.NoError(t, err)

					// Create a session for the device
					session := &models.DeviceSession{
						DeviceID:  device.GetID(),
						UserAgent: "Test Agent",
						IP:        "127.0.0.1",
					}
					session.GenID(ctx)
					err = repository.NewDeviceSessionRepository(svc).Save(ctx, session)
					require.NoError(t, err)
					sessionID = session.GetID()
				} else {
					sessionID = tc.sessionID
				}

				linkedDevice, err := biz.LinkDeviceToProfile(ctx, sessionID, tc.profileID, map[string]string{
					"link_reason": "test",
				})

				if tc.expectError {
					require.Error(t, err)
					assert.Nil(t, linkedDevice)
				} else {
					require.NoError(t, err)
					assert.NotNil(t, linkedDevice)
					// Note: LinkDeviceToProfile only updates if profile is empty, so existing profile remains
					assert.NotEmpty(t, linkedDevice.GetId())
				}
			})
		}
	})
}

// Helper method to create device with profile for search testing.
func (suite *DeviceBusinessTestSuite) createDeviceWithProfile(
	ctx context.Context,
	svc *frame.Service,
	profileID string,
) error {
	device := &models.Device{
		ProfileID: profileID,
		Name:      "Search Test Device",
		OS:        "Linux",
	}
	device.GenID(ctx)
	err := repository.NewDeviceRepository(svc).Save(ctx, device)
	if err != nil {
		return err
	}

	session := &models.DeviceSession{
		DeviceID:  device.GetID(),
		UserAgent: "Test Agent",
		IP:        "127.0.0.1",
	}
	session.GenID(ctx)
	return repository.NewDeviceSessionRepository(svc).Save(ctx, session)
}

// Helper method to process search results channel.
func (suite *DeviceBusinessTestSuite) processSearchResults(
	devicesChan <-chan frame.JobResult[[]*devicev1.DeviceObject],
) ([]*devicev1.DeviceObject, error) {
	var devices []*devicev1.DeviceObject
	for result := range devicesChan {
		if result.IsError() {
			return nil, result.Error()
		}
		devices = append(devices, result.Item()...)
	}
	return devices, nil
}

// Helper method to create device for logs testing.
func (suite *DeviceBusinessTestSuite) createDeviceForLogs(
	ctx context.Context,
	svc *frame.Service,
	addLog bool,
	biz business.DeviceBusiness,
) (string, error) {
	device := &models.Device{
		Name: "Log Test Device",
		OS:   "Linux",
	}
	device.GenID(ctx)
	err := repository.NewDeviceRepository(svc).Save(ctx, device)
	if err != nil {
		return "", err
	}

	deviceID := device.GetID()
	if addLog {
		_, logErr := biz.LogDeviceActivity(ctx, deviceID, "test-session-10", map[string]string{
			"action": "test_action",
		})
		if logErr != nil {
			return "", logErr
		}
	}
	return deviceID, nil
}

// Helper method to process device logs results channel.
func (suite *DeviceBusinessTestSuite) processDeviceLogsResults(
	logsChan <-chan frame.JobResult[[]*devicev1.DeviceLog],
) ([]*devicev1.DeviceLog, error) {
	var logs []*devicev1.DeviceLog
	for result := range logsChan {
		if result.IsError() {
			return nil, result.Error()
		}
		logs = append(logs, result.Item()...)
	}
	return logs, nil
}

// Helper method to create device with key for keys testing.
func (suite *DeviceBusinessTestSuite) createDeviceWithKey(
	ctx context.Context,
	svc *frame.Service,
	biz business.DeviceBusiness,
	addKey bool,
	keyType devicev1.KeyType,
) (string, error) {
	device := &models.Device{
		Name: "Key Test Device",
		OS:   "Linux",
	}
	device.GenID(ctx)
	err := repository.NewDeviceRepository(svc).Save(ctx, device)
	if err != nil {
		return "", err
	}

	deviceID := device.GetID()
	if addKey {
		_, keyErr := biz.AddKey(ctx, deviceID, keyType, []byte("test-key-data"), map[string]string{
			"test": "data",
		})
		if keyErr != nil {
			return "", keyErr
		}
	}
	return deviceID, nil
}

// Helper method to process keys results channel.
func (suite *DeviceBusinessTestSuite) processKeysResults(
	keysChan <-chan frame.JobResult[[]*devicev1.KeyObject],
) ([]*devicev1.KeyObject, error) {
	var keys []*devicev1.KeyObject
	for result := range keysChan {
		if result.IsError() {
			return nil, result.Error()
		}
		keys = append(keys, result.Item()...)
	}
	return keys, nil
}
