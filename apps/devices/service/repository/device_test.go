package repository_test

import (
	"testing"

	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/tests"
)

const (
	DefaultRandomStringLength = 8
)

type DeviceRepositoryTestSuite struct {
	tests.DeviceBaseTestSuite
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

	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create device
				device := &models.Device{
					ProfileID: tc.profileID,
					Name:      tc.deviceName,
					OS:        tc.os,
				}

				err := deps.DeviceRepo.Create(ctx, device)
				if tc.expectError {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.NotEmpty(t, device.GetID())

				// Retrieve by ID
				retrievedDevice, err := deps.DeviceRepo.GetByID(ctx, device.GetID())
				require.NoError(t, err)
				assert.Equal(t, device.GetID(), retrievedDevice.GetID())
				assert.Equal(t, tc.profileID, retrievedDevice.ProfileID)
				assert.Equal(t, tc.deviceName, retrievedDevice.Name)
				assert.Equal(t, tc.os, retrievedDevice.OS)

				// Retrieve by ProfileID
				searchProperties := data.JSONMap{
					"profile_id": tc.profileID,
				}
				q := data.NewSearchQuery("", data.WithSearchFiltersAndByValue(searchProperties),
					data.WithSearchOffset(0), data.WithSearchLimit(50))
				devicesResult, err := deps.DeviceRepo.Search(ctx, q)
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
				removedDevice, err := deps.DeviceRepo.RemoveByID(ctx, device.GetID())
				require.NoError(t, err)
				assert.Equal(t, device.GetID(), removedDevice.GetID())

				// Verify removal
				_, err = deps.DeviceRepo.GetByID(ctx, device.GetID())
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

	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create session
				session := &models.DeviceSession{
					DeviceID:  tc.deviceID,
					UserAgent: tc.userAgent,
					IP:        tc.ip,
				}

				err := deps.SessionRepo.Create(ctx, session)
				if tc.expectError {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.NotEmpty(t, session.GetID())

				// Retrieve by ID
				retrievedSession, err := deps.SessionRepo.GetByID(ctx, session.GetID())
				require.NoError(t, err)
				assert.Equal(t, session.GetID(), retrievedSession.GetID())
				assert.Equal(t, tc.deviceID, retrievedSession.DeviceID)
				assert.Equal(t, tc.userAgent, retrievedSession.UserAgent)
				assert.Equal(t, tc.ip, retrievedSession.IP)

				// Retrieve last by device ID
				lastSession, err := deps.SessionRepo.GetLastByDeviceID(ctx, tc.deviceID)
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

	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create log
				log := &models.DeviceLog{
					DeviceID:        tc.deviceID,
					DeviceSessionID: tc.sessionID,
					Data:            data.JSONMap{"action": "test"},
				}

				err := deps.DeviceLogRepo.Create(ctx, log)
				if tc.expectError {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.NotEmpty(t, log.GetID())

				// Retrieve by ID
				retrievedLog, err := deps.DeviceLogRepo.GetByID(ctx, log.GetID())
				require.NoError(t, err)
				assert.Equal(t, log.GetID(), retrievedLog.GetID())
				assert.Equal(t, tc.deviceID, retrievedLog.DeviceID)
				assert.Equal(t, tc.sessionID, retrievedLog.DeviceSessionID)

				// Retrieve by device ID
				logsResult, err := deps.DeviceLogRepo.GetByDeviceID(ctx, tc.deviceID)
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

	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create key
				key := &models.DeviceKey{
					DeviceID: tc.deviceID,
					Key:      tc.key,
					Extra:    data.JSONMap{"type": "test"},
				}

				err := deps.KeyRepo.Create(ctx, key)
				if tc.expectError {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.NotEmpty(t, key.GetID())

				// Retrieve by device ID
				keys, err := deps.KeyRepo.GetByDeviceID(ctx, tc.deviceID)
				require.NoError(t, err)
				assert.Len(t, keys, 1)
				assert.Equal(t, key.GetID(), keys[0].GetID())
				assert.Equal(t, tc.deviceID, keys[0].DeviceID)
				assert.Equal(t, tc.key, keys[0].Key)

				// Remove key
				removedKey, err := deps.KeyRepo.RemoveByID(ctx, key.GetID())
				require.NoError(t, err)
				assert.Equal(t, key.GetID(), removedKey.GetID())

				// Verify removal
				keys, err = deps.KeyRepo.GetByDeviceID(ctx, tc.deviceID)
				require.NoError(t, err)
				assert.Empty(t, keys)
			})
		}
	})
}
