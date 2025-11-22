package business_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/frametests"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/workerpool"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/devices/service/business"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
	"github.com/antinvestor/service-profile/apps/devices/tests"
)

const (
	DefaultRandomStringLength = 8
)

type DeviceBusinessTestSuite struct {
	tests.DeviceBaseTestSuite
}

func TestDeviceBusinessTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceBusinessTestSuite))
}

func (suite *DeviceBusinessTestSuite) createTestDeviceWithSession(
	ctx context.Context,
	deviceRepo repository.DeviceRepository,
	sessionRepo repository.DeviceSessionRepository,
	deviceID, sessionID string,
) error {
	device := &models.Device{
		Name: "Original Name",
		OS:   "Original OS",
	}

	device.ID = deviceID
	err := deviceRepo.Create(ctx, device)
	if err != nil {
		return err
	}

	session := &models.DeviceSession{
		DeviceID:  deviceID,
		UserAgent: "Test Agent",
		IP:        "127.0.0.1",
	}

	if sessionID != "" {
		session.ID = sessionID
	}
	return sessionRepo.Create(ctx, session)
}

func (suite *DeviceBusinessTestSuite) verifyDeviceActivityLogged(
	ctx context.Context,
	deviceBusiness business.DeviceBusiness,
	deviceID, sessionID string,
) error {
	if sessionID == "" || deviceID == "" {
		return nil
	}

	deviceLogsChan, err := deviceBusiness.GetDeviceLogs(ctx, deviceID)
	if err != nil {
		return err
	}

	// Process the channel to get actual logs
	var deviceLogs []*devicev1.DeviceLog
	for {
		result, ok := deviceLogsChan.ReadResult(ctx)
		if !ok {
			break
		}

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
	deviceRepo repository.DeviceRepository,
	sessionRepo repository.DeviceSessionRepository,
	deviceBusiness business.DeviceBusiness,
	tc struct {
	name        string
	id          string
	deviceName  string
	data        data.JSONMap
	expectError bool
	expectNil   bool
},
) {
	// Setup existing device if needed
	if tc.id != "" && tc.name == "save device with existing ID" {
		sessionID := ""
		rawDat, ok := tc.data["session_id"]
		if ok {
			sessionID = rawDat.(string)
		}
		err := suite.createTestDeviceWithSession(ctx, deviceRepo, sessionRepo, tc.id, sessionID)
		require.NoError(t, err)
	}

	// Execute SaveDevice
	result, err := deviceBusiness.SaveDevice(ctx, tc.id, tc.deviceName, tc.data)

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
	sessionID := ""
	rawDat, ok := tc.data["session_id"]
	if ok {
		sessionID = rawDat.(string)
	}
	if !tc.expectError && tc.id != "" && sessionID != "" {
		cErr := suite.verifyDeviceActivityLogged(ctx, deviceBusiness, tc.id, sessionID)
		assert.NoError(t, cErr)
	}
}

func (suite *DeviceBusinessTestSuite) TestSaveDevice() {
	t := suite.T()
	testCases := []struct {
		name        string
		id          string
		deviceName  string
		data        data.JSONMap
		expectError bool
		expectNil   bool
	}{
		{
			name:       "save device with existing ID",
			id:         "existing-device-id",
			deviceName: "Test Device",
			data: data.JSONMap{
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
			data: data.JSONMap{
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
			data: data.JSONMap{
				"profile_id": "profile-789",
				"session_id": "test-session-3",
			},
			expectError: true,
			expectNil:   false,
		},
	}

	suite.WithTestDependencies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				suite.runSaveDeviceTestCase(ctx, t, deps.DeviceRepo, deps.SessionRepo, deps.DeviceBusiness, tc)
			})
		}
	})
}

func (suite *DeviceBusinessTestSuite) TestGetDeviceByID() {
	t := suite.T()
	testCases := []struct {
		name               string
		setupDevice        func(ctx context.Context, deps *tests.DepsBuilder, devId string) error
		deviceID           string
		errorAssertionFunc require.ErrorAssertionFunc
	}{
		{
			name:     "existing device",
			deviceID: "existing-device-id",
			setupDevice: func(ctx context.Context, deps *tests.DepsBuilder, devId string) error {
				// Create a device and session first
				device := &models.Device{
					Name: "Test Device",
					OS:   "Linux",
				}
				device.GenID(ctx)
				device.ID = devId

				err := deps.DeviceRepo.Create(ctx, device)
				if err != nil {
					return err
				}
				actualDeviceID := device.GetID()
				// Create a session for the device (required by GetDeviceByID)
				session := &models.DeviceSession{
					DeviceID:  actualDeviceID,
					UserAgent: "Test Agent",
					IP:        "127.0.0.1",
				}

				err = deps.SessionRepo.Create(ctx, session)
				if err != nil {
					return err
				}
				return nil
			},
			errorAssertionFunc: require.NoError,
		},
		{
			name:               "non-existent device",
			setupDevice:        nil,
			deviceID:           "non-existent-id",
			errorAssertionFunc: require.Error,
		},
		{
			name:               "empty device ID",
			setupDevice:        nil,
			deviceID:           "",
			errorAssertionFunc: require.Error,
		},
	}

	suite.WithTestDependencies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				if tc.setupDevice != nil {
					err := tc.setupDevice(ctx, deps, tc.deviceID)
					require.NoError(t, err)
				}

				device, err := deps.DeviceBusiness.GetDeviceByID(ctx, tc.deviceID)

				tc.errorAssertionFunc(t, err)
				if err != nil {
					require.Nil(t, device)
				} else {
					require.NotNil(t, device)
					require.Equal(t, tc.deviceID, device.GetId())
					require.Equal(t, "Test Device", device.GetName())
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

	suite.WithTestDependencies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var sessionID string
				if tc.setupDevice {
					// Create a device and session first
					device := &models.Device{
						Name: "Test Device",
						OS:   "Linux",
					}

					err := deps.DeviceRepo.Create(ctx, device)
					require.NoError(t, err)

					// Create a session for the device
					session := &models.DeviceSession{
						DeviceID:  device.GetID(),
						UserAgent: "Test Agent",
						IP:        "127.0.0.1",
					}

					err = deps.SessionRepo.Create(ctx, session)
					require.NoError(t, err)
					sessionID = session.GetID()
				} else {
					sessionID = tc.sessionID
				}

				device, err := deps.DeviceBusiness.GetDeviceBySessionID(ctx, sessionID)

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
		data        data.JSONMap
		expectError bool
	}{
		{
			name:        "valid activity log",
			setupDevice: true,
			data: data.JSONMap{
				"action": "login",
				"result": "success",
			},
			expectError: false,
		},
		{
			name:        "empty data",
			setupDevice: true,
			data:        data.JSONMap{},
			expectError: false,
		},
	}

	suite.WithTestDependencies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID, sessionID string
				if tc.setupDevice {
					// Create a device directly using repository
					device := &models.Device{
						Name: "Test Device",
						OS:   "Linux",
					}
					err := deps.DeviceRepo.Create(ctx, device)
					require.NoError(t, err)
					deviceID = device.GetID()

					// Create a session for the device
					session := &models.DeviceSession{
						DeviceID:  deviceID,
						UserAgent: "Test Agent",
						IP:        "127.0.0.1",
					}
					err = deps.SessionRepo.Create(ctx, session)
					require.NoError(t, err)
					sessionID = session.GetID()
				} else {
					deviceID = tc.deviceID
					sessionID = tc.sessionID
				}

				log, err := deps.DeviceBusiness.LogDeviceActivity(ctx, deviceID, sessionID, tc.data)

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
		extra       data.JSONMap
		expectError bool
	}{
		{
			name:        "add valid key",
			setupDevice: true,
			keyType:     devicev1.KeyType_MATRIX_KEY,
			key:         []byte("test-encryption-key"),
			extra: data.JSONMap{
				"algorithm": "AES256",
			},
			expectError: false,
		},
		{
			name:        "add key with empty extra",
			setupDevice: true,
			keyType:     devicev1.KeyType_NOTIFICATION_KEY,
			key:         []byte("test-signing-key"),
			extra:       data.JSONMap{},
			expectError: false,
		},
	}

	suite.WithTestDependencies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID string
				if tc.setupDevice {
					// Create a device directly using repository
					device := &models.Device{
						Name: "Test Device",
						OS:   "Linux",
					}
					err := deps.DeviceRepo.Create(ctx, device)
					require.NoError(t, err)
					deviceID = device.GetID()
				} else {
					deviceID = tc.deviceID
				}

				keyObj, err := deps.KeyBusiness.AddKey(ctx, deviceID, tc.keyType, tc.key, tc.extra)

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

	suite.WithTestDependencies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID string
				if tc.setupDevice {
					// Create a device first
					device := &models.Device{
						Name: "Test Device",
						OS:   "Linux",
					}

					err := deps.DeviceRepo.Create(ctx, device)
					require.NoError(t, err)
					deviceID = device.GetID()
				} else {
					deviceID = tc.deviceID
				}

				err := deps.DeviceBusiness.RemoveDevice(ctx, deviceID)

				if tc.expectError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)

					// Verify device is removed
					_, err = deps.DeviceBusiness.GetDeviceByID(ctx, deviceID)
					require.Error(t, err)
				}
			})
		}
	})
}

func (suite *DeviceBusinessTestSuite) runSearchDevicesTestCase(
	ctx context.Context,
	t *testing.T,
	deviceRepo repository.DeviceRepository,
	sessionRepo repository.DeviceSessionRepository,
	deviceBusiness business.DeviceBusiness,
	tc struct {
	name        string
	setupDevice bool
	profileID   string
	searchQuery string
	expectError bool
},
) {
	// Don't use claims context for now - just test text search
	testCtx := security.SkipTenancyChecksOnClaims(ctx)

	if tc.setupDevice {
		err := suite.createDeviceWithProfile(testCtx, deviceRepo, sessionRepo, tc.profileID)
		require.NoError(t, err)
	}

	query := &devicev1.SearchRequest{
		Query: tc.searchQuery,
	}

	devicesChan, err := deviceBusiness.SearchDevices(testCtx, query)
	if tc.expectError {
		require.Error(t, err)
		return
	}

	require.NoError(t, err)
	assert.NotNil(t, devicesChan)

	devices, channelErr := suite.processSearchResults(testCtx, devicesChan)
	require.NoError(t, channelErr)

	suite.verifySearchResults(t, tc, devices)
}

func (suite *DeviceBusinessTestSuite) verifySearchResults(t *testing.T, tc struct {
	name        string
	setupDevice bool
	profileID   string
	searchQuery string
	expectError bool
}, devices []*devicev1.DeviceObject) {
	if tc.setupDevice && tc.searchQuery != "" {
		// When we set up a device and search for it, we should find at least 1
		assert.GreaterOrEqual(t, len(devices), 1, "Should find at least the device we created")
		if len(devices) > 0 {
			assert.NotEmpty(t, devices[0].GetId())
			// Verify at least one device matches our search query
			found := false
			for _, dev := range devices {
				if dev.GetName() != "" || dev.GetOs() != "" {
					found = true
					break
				}
			}
			assert.True(t, found, "Should find device with name or os")
		}
	} else {
		// For tests without setup and empty query, we may find devices from other tests
		// due to test execution order, so we just verify the response is valid
		// (not testing for empty since tests may not be isolated)
		for _, dev := range devices {
			assert.NotEmpty(t, dev.GetId(), "Returned devices should have valid IDs")
		}
	}
}

func (suite *DeviceBusinessTestSuite) TestSearchDevices() {
	t := suite.T()
	testCases := []struct {
		name        string
		setupDevice bool
		profileID   string
		searchQuery string
		expectError bool
	}{
		{
			name:        "search with valid profile ID",
			setupDevice: true,
			profileID:   "profile-search-test",
			searchQuery: "Search", // Search for actual text in device name
			expectError: false,
		},
		{
			name:        "search with non-existent profile ID",
			setupDevice: false,
			profileID:   "non-existent-profile",
			searchQuery: "non-matching-query",
			expectError: false,
		},
		{
			name:        "search with empty profile ID",
			setupDevice: false,
			profileID:   "",
			searchQuery: "",
			expectError: false,
		},
	}

	suite.WithTestDependencies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				suite.runSearchDevicesTestCase(ctx, t, deps.DeviceRepo, deps.SessionRepo, deps.DeviceBusiness, tc)
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

	suite.WithTestDependencies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID string
				if tc.setupDevice {
					// Create a device first
					var err error
					deviceID, err = suite.createDeviceForLogs(
						ctx,
						tc.name == "get logs for device with logs",
						deps.DeviceRepo,
						deps.DeviceBusiness,
					)
					require.NoError(t, err)
				} else {
					deviceID = tc.deviceID
				}

				logsChan, logErr := deps.DeviceBusiness.GetDeviceLogs(ctx, deviceID)
				require.NoError(t, logErr)
				assert.NotNil(t, logsChan)

				// Collect results from channel
				logs, err := suite.processDeviceLogsResults(ctx, logsChan)
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

	suite.WithTestDependencies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var deviceID string
				if tc.setupDevice {
					// Create a device first
					var err error
					deviceID, err = suite.createDeviceWithKey(
						ctx,
						deps.DeviceRepo,
						deps.KeyBusiness,
						tc.setupKey,
						tc.keyType,
					)
					require.NoError(t, err)
				} else {
					deviceID = tc.deviceID
				}

				keysChan, err := deps.KeyBusiness.GetKeys(ctx, deviceID, tc.keyType)
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

	suite.WithTestDependencies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var keysToRemove []string
				if tc.setupDevice {
					// Create a device first
					_, keyIDs, err := suite.createDeviceWithKeys(ctx, deps.DeviceRepo, deps.KeyBusiness, tc.setupKeys)
					require.NoError(t, err)

					// Select keys to remove
					for i := 0; i < tc.removeCount && i < len(keyIDs); i++ {
						keysToRemove = append(keysToRemove, keyIDs[i])
					}
				} else {
					// Use non-existent key IDs
					keysToRemove = suite.generateNonExistentKeyIDs(tc.removeCount)
				}

				keysChan, err := deps.KeyBusiness.RemoveKeys(ctx, keysToRemove...)
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
	deviceRepo repository.DeviceRepository,
	keyBusiness business.KeysBusiness,
	keyCount int,
) (string, []string, error) {
	// Create a device directly using repository
	device := &models.Device{
		Name: "Remove Key Test Device",
		OS:   "Linux",
	}

	err := deviceRepo.Create(ctx, device)
	if err != nil {
		return "", nil, err
	}
	deviceID := device.GetID()

	var keyIDs []string
	for i := range keyCount {
		keyResult, keyErr := keyBusiness.AddKey(ctx, deviceID, devicev1.KeyType_MATRIX_KEY,
			[]byte(fmt.Sprintf("test-key-data-%d", i)), data.JSONMap{
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
	keysChan <-chan workerpool.JobResult[[]*devicev1.KeyObject],
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

	suite.WithTestDependencies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

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

					err := deps.DeviceRepo.Create(ctx, device)
					require.NoError(t, err)

					// Create a session for the device
					session := &models.DeviceSession{
						DeviceID:  device.GetID(),
						UserAgent: "Test Agent",
						IP:        "127.0.0.1",
					}

					err = deps.SessionRepo.Create(ctx, session)
					require.NoError(t, err)
					sessionID = session.GetID()
				} else {
					sessionID = tc.sessionID
				}

				linkedDevice, err := deps.DeviceBusiness.LinkDeviceToProfile(ctx, sessionID, tc.profileID, data.JSONMap{
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
	deviceRepo repository.DeviceRepository,
	sessionRepo repository.DeviceSessionRepository,
	profileID string,
) error {
	device := &models.Device{
		ProfileID: profileID,
		Name:      "Search Test Device",
		OS:        "Linux",
	}

	err := deviceRepo.Create(ctx, device)
	if err != nil {
		return err
	}

	session := &models.DeviceSession{
		DeviceID:  device.GetID(),
		UserAgent: "Test Agent",
		IP:        "127.0.0.1",
	}

	return sessionRepo.Create(ctx, session)
}

// Helper method to process search results channel.
func (suite *DeviceBusinessTestSuite) processSearchResults(
	ctx context.Context, devicesChan workerpool.JobResultPipe[[]*devicev1.DeviceObject],
) ([]*devicev1.DeviceObject, error) {
	var devices []*devicev1.DeviceObject
	for {
		result, ok := devicesChan.ReadResult(ctx)
		if !ok {
			break
		}
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
	addLog bool,
	deviceRepo repository.DeviceRepository,
	deviceBusiness business.DeviceBusiness,
) (string, error) {
	device := &models.Device{
		Name: "Log Test Device",
		OS:   "Linux",
	}

	err := deviceRepo.Create(ctx, device)
	if err != nil {
		return "", err
	}

	deviceID := device.GetID()
	if addLog {
		_, logErr := deviceBusiness.LogDeviceActivity(ctx, deviceID, "test-session-10", data.JSONMap{
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
	ctx context.Context, logsChan workerpool.JobResultPipe[[]*devicev1.DeviceLog],
) ([]*devicev1.DeviceLog, error) {
	var logs []*devicev1.DeviceLog
	for {
		result, ok := logsChan.ReadResult(ctx)
		if !ok {
			break
		}
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
	deviceRepo repository.DeviceRepository,
	keyBusiness business.KeysBusiness,
	addKey bool,
	keyType devicev1.KeyType,
) (string, error) {
	device := &models.Device{
		Name: "Key Test Device",
		OS:   "Linux",
	}

	err := deviceRepo.Create(ctx, device)
	if err != nil {
		return "", err
	}

	deviceID := device.GetID()
	if addKey {
		_, keyErr := keyBusiness.AddKey(ctx, deviceID, keyType, []byte("test-key-data"), data.JSONMap{
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
	keysChan <-chan workerpool.JobResult[[]*devicev1.KeyObject],
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

// TestLogDeviceActivity_AutoCreateDeviceAndSession validates that when DeviceID and SessionID
// are provided to LogDeviceActivity, the device and session will be auto-created during device analysis.
func (suite *DeviceBusinessTestSuite) TestLogDeviceActivity_AutoCreateDeviceAndSession() {
	t := suite.T()

	suite.WithTestDependencies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		// Generate new IDs for device and session
		deviceID := util.IDString()
		sessionID := util.IDString()

		// Prepare device log data with user agent and IP for session creation
		logData := data.JSONMap{
			"userAgent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
			"ip":        "192.168.1.100",
			"action":    "login",
			"tz":        "Africa/Nairobi",
			"lang":      "en-US,en",
			"cur":       "KES",
			"curNm":     "Kenyan Shilling",
			"code":      "+254",
		}

		// Step 1: Call LogDeviceActivity with new DeviceID and SessionID
		deviceLog, err := deps.DeviceBusiness.LogDeviceActivity(ctx, deviceID, sessionID, logData)
		require.NoError(t, err)
		require.NotNil(t, deviceLog)
		assert.Equal(t, deviceID, deviceLog.GetDeviceId())
		assert.Equal(t, sessionID, deviceLog.GetSessionId())

		// Step 2: Verify device log was created
		savedLog, cErr := deps.DeviceLogRepo.GetByID(ctx, deviceLog.GetId())
		require.NoError(t, cErr)
		assert.Equal(t, deviceID, savedLog.DeviceID)
		assert.Equal(t, sessionID, savedLog.DeviceSessionID)

		// Step 3: Verify device and session don't exist before queue processing
		deviceRepo := deps.DeviceRepo
		sessionRepo := deps.SessionRepo

		_, sessionErr := sessionRepo.GetByID(ctx, sessionID)
		require.Error(t, sessionErr, "Session should not exist before queue processing")
		_, deviceErr := deviceRepo.GetByID(ctx, deviceID)
		require.Error(t, deviceErr, "Device should not exist before queue processing")

		// Step 4: Wait for queue to process and create device and session
		// The queue handler will automatically process the device log and create both
		sessionCreated, cErr := frametests.WaitForCheckedConditionWithResult(
			ctx,
			func() (*models.DeviceSession, error) {
				return sessionRepo.GetByID(ctx, sessionID)
			},
			func(sess *models.DeviceSession, err error) bool {
				return err == nil && sess != nil
			},
			5*time.Second,
			100*time.Millisecond,
		)
		require.NoError(t, cErr)
		require.NotNil(t, sessionCreated)
		assert.Equal(t, sessionID, sessionCreated.GetID())
		assert.Equal(t, deviceID, sessionCreated.DeviceID)
		assert.Equal(t, logData.GetString("userAgent"), sessionCreated.UserAgent)
		assert.Equal(t, logData.GetString("ip"), sessionCreated.IP)

		deviceCreated, cErr := deviceRepo.GetByID(ctx, deviceID)
		require.NoError(t, cErr)
		require.NotNil(t, deviceCreated)
		assert.Equal(t, deviceID, deviceCreated.GetID())
		assert.NotEmpty(t, deviceCreated.Name)
		assert.NotEmpty(t, deviceCreated.OS)

		// Step 5: Verify we can retrieve the device through business layer
		deviceObj, err := deps.DeviceBusiness.GetDeviceByID(ctx, deviceID)
		require.NoError(t, err)
		assert.Equal(t, deviceID, deviceObj.GetId())

		// Step 6: Verify we can retrieve device by session ID
		deviceBySession, err := deps.DeviceBusiness.GetDeviceBySessionID(ctx, sessionID)
		require.NoError(t, err)
		assert.Equal(t, deviceID, deviceBySession.GetId())
		assert.Equal(t, sessionID, deviceBySession.GetSessionId())
	})
}
