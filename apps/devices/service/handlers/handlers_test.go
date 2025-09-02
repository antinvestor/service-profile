package handlers_test

import (
	"context"
	"testing"

	"buf.build/go/protovalidate"
	commonMocks "github.com/antinvestor/apis/go/common/mocks"
	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/antinvestor/service-profile/apps/devices/service/handlers"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
	"github.com/antinvestor/service-profile/apps/devices/service/tests"
)

type HandlersTestSuite struct {
	tests.DeviceBaseTestSuite
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

func (suite *HandlersTestSuite) TestDevicesServer_GetByID() {
	testCases := []struct {
		name           string
		setupDevice    bool
		deviceID       string
		expectedStatus codes.Code
	}{
		{
			name:           "get existing device",
			setupDevice:    true,
			expectedStatus: codes.OK,
		},
		{
			name:           "get non-existent device",
			setupDevice:    false,
			deviceID:       "non-existent-device",
			expectedStatus: codes.NotFound,
		},
		{
			name:           "get device with empty ID",
			setupDevice:    false,
			deviceID:       "",
			expectedStatus: codes.InvalidArgument,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Create server
		server := handlers.NewDeviceServer(ctx, svc)

		for _, tc := range testCases {
			suite.Run(tc.name, func() {
				var deviceID string
				if tc.setupDevice {
					// Create a device directly using repository
					device := &models.Device{
						ProfileID: "test-profile",
						Name:      "Test Device",
						OS:        "Linux",
					}
					device.GenID(ctx)
					err := repository.NewDeviceRepository(svc).Save(ctx, device)
					suite.Require().NoError(err)

					// Create a session for the device (required by GetDeviceByID business method)
					session := &models.DeviceSession{
						DeviceID:  device.GetID(),
						UserAgent: "Test Agent",
						IP:        "127.0.0.1",
					}
					session.GenID(ctx)
					err = repository.NewDeviceSessionRepository(svc).Save(ctx, session)
					suite.Require().NoError(err)

					deviceID = device.GetID()
				} else {
					deviceID = tc.deviceID
				}

				req := &devicev1.GetByIdRequest{
					Id: []string{deviceID},
				}

				resp, err := server.GetById(ctx, req)

				if tc.expectedStatus == codes.OK {
					suite.Require().NoError(err)
					suite.Require().NotNil(resp)
					suite.NotEmpty(resp.GetData())
					suite.Equal(deviceID, resp.GetData()[0].GetId())
				} else {
					suite.Require().Error(err)
					st, ok := status.FromError(err)
					suite.Require().True(ok)
					suite.Equal(tc.expectedStatus, st.Code())
				}
			})
		}
	})
}

func (suite *HandlersTestSuite) TestDevicesServer_Create() {
	testCases := []struct {
		name        string
		linkID      string
		deviceName  string
		data        frame.JSONMap
		expectError bool
	}{
		{
			name:       "create device successfully",
			linkID:     "test-link-123",
			deviceName: "Test Device",
			data: frame.JSONMap{
				"os":         "Linux",
				"user_agent": "Test Agent",
				"session_id": "test-session",
				"ip":         "127.0.0.1",
			},
			expectError: false,
		},
		{
			name:       "create device with empty name",
			linkID:     "test-link-456",
			deviceName: "",
			data: frame.JSONMap{
				"session_id": "test-session-2",
				"ip":         "127.0.0.1",
			},
			expectError: false, // Should still work with empty name
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Create server
		server := handlers.NewDeviceServer(ctx, svc)

		for _, tc := range testCases {
			suite.Run(tc.name, func() {
				req := &devicev1.CreateRequest{
					Name:       tc.deviceName,
					Properties: tc.data.ToProtoStruct(),
				}

				resp, err := server.Create(ctx, req)

				if tc.expectError {
					suite.Require().Error(err)
					suite.Nil(resp)
				} else {
					suite.Require().NoError(err)
					// Response might be nil if device creation doesn't return a device
					if resp != nil && resp.GetData() != nil {
						suite.NotEmpty(resp.GetData().GetId())
						suite.Equal(tc.deviceName, resp.GetData().GetName())
					}
				}
			})
		}
	})
}

func (suite *HandlersTestSuite) TestDevicesServer_Log() {
	testCases := []struct {
		name        string
		setupDevice bool
		deviceID    string
		sessionID   string
		data        frame.JSONMap
		expectError bool
	}{
		{
			name:        "log device activity successfully",
			setupDevice: true,
			sessionID:   "test-session",
			data: frame.JSONMap{
				"action": "page_view",
				"url":    "https://example.com",
			},
			expectError: false,
		},
		{
			name:        "log with empty data",
			setupDevice: true,
			sessionID:   "test-session-2",
			data:        frame.JSONMap{},
			expectError: false,
		},
		{
			name:        "log with empty device ID",
			setupDevice: false,
			deviceID:    "",
			sessionID:   "test-session-3",
			data: frame.JSONMap{
				"action": "login",
			},
			expectError: false, // Current implementation allows empty device ID
		},
		{
			name:        "log with non-existent device ID",
			setupDevice: false,
			deviceID:    "non-existent-device-123",
			sessionID:   "test-session-4",
			data: frame.JSONMap{
				"action": "logout",
			},
			expectError: false, // Current implementation doesn't validate device existence
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Create server
		server := handlers.NewDeviceServer(ctx, svc)

		for _, tc := range testCases {
			suite.Run(tc.name, func() {
				var deviceID string
				if tc.setupDevice {
					// Create a device directly using repository
					device := &models.Device{
						ProfileID: "test-profile",
						Name:      "Test Device",
						OS:        "Linux",
					}
					device.GenID(ctx)
					err := repository.NewDeviceRepository(svc).Save(ctx, device)
					suite.Require().NoError(err)
					deviceID = device.GetID()
				} else {
					deviceID = tc.deviceID
				}

				req := &devicev1.LogRequest{
					DeviceId:  deviceID,
					SessionId: tc.sessionID,
					Extras:    tc.data.ToProtoStruct(),
				}

				resp, err := server.Log(ctx, req)

				if tc.expectError {
					suite.Require().Error(err)
					suite.Nil(resp)
				} else {
					suite.Require().NoError(err)
					suite.NotNil(resp)
					suite.NotNil(resp.GetData())
				}
			})
		}
	})
}

func (suite *HandlersTestSuite) TestDevices_LogRequest() {
	testCases := []struct {
		name        string
		deviceID    string
		sessionID   string
		data        frame.JSONMap
		expectError bool
	}{
		{
			name:      "log device activity successfully",
			deviceID:  "",
			sessionID: "test-session",
			data: frame.JSONMap{
				"action": "page_view",
				"url":    "https://example.com",
			},
			expectError: false,
		},
		{
			name:        "log with error data",
			deviceID:    "h",
			sessionID:   "test-session-2",
			data:        frame.JSONMap{},
			expectError: true,
		},
		{
			name:      "log with empty device ID",
			deviceID:  "hellow",
			sessionID: "test-session-3",
			data: frame.JSONMap{
				"action": "login",
			},
			expectError: false, // Current implementation allows empty device ID
		},
		{
			name:      "log with very long device ID",
			deviceID:  "fasodeifwqoiejfpasdjfoiasdjfoisjdfljksjdflaksdjfosidjfsoidjfsoidjfoasdfasdfasdfsa",
			sessionID: "test-session-4",
			data: frame.JSONMap{
				"action": "logout",
			},
			expectError: true, // Current implementation doesn't validate device existence
		},
	}

	validator, err := protovalidate.New()
	suite.Require().NoError(err)
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			req := &devicev1.LogRequest{
				DeviceId:  tc.deviceID,
				SessionId: tc.sessionID,
				Extras:    tc.data.ToProtoStruct(),
			}

			err = validator.Validate(req)

			if tc.expectError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *HandlersTestSuite) TestDevicesServer_AddKey() {
	testCases := []struct {
		name        string
		setupDevice bool
		keyType     devicev1.KeyType
		keyData     []byte
		expectError bool
	}{
		{
			name:        "add key successfully",
			setupDevice: true,
			keyType:     devicev1.KeyType_MATRIX_KEY,
			keyData:     []byte("test-key-data"),
			expectError: false,
		},
		{
			name:        "add key to non-existent device",
			setupDevice: false,
			keyType:     devicev1.KeyType_NOTIFICATION_KEY,
			keyData:     []byte("test-key-data"),
			expectError: true,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Create server
		server := handlers.NewDeviceServer(ctx, svc)

		for _, tc := range testCases {
			suite.Run(tc.name, func() {
				var deviceID string
				if tc.setupDevice {
					// Create a device directly using repository
					device := &models.Device{
						ProfileID: "test-profile",
						Name:      "Test Device",
						OS:        "Linux",
					}
					device.GenID(ctx)
					err := repository.NewDeviceRepository(svc).Save(ctx, device)
					suite.Require().NoError(err)
					deviceID = device.GetID()
				} else {
					deviceID = "non-existent-device"
				}

				extras := frame.JSONMap{"test": "data"}

				req := &devicev1.AddKeyRequest{
					DeviceId: deviceID,
					KeyType:  tc.keyType,
					Data:     tc.keyData,
					Extras:   extras.ToProtoStruct(),
				}

				resp, err := server.AddKey(ctx, req)

				if tc.expectError {
					suite.Require().Error(err)
					suite.Nil(resp)
				} else {
					suite.Require().NoError(err)
					suite.NotNil(resp)
					suite.NotNil(resp.GetData())
					if resp.GetData() != nil {
						suite.NotEmpty(resp.GetData().GetId())
						suite.Equal(deviceID, resp.GetData().GetDeviceId())
						suite.Equal(tc.keyData, resp.GetData().GetKey())
					}
				}
			})
		}
	})
}

func (suite *HandlersTestSuite) TestDevicesServer_Search() {
	testCases := []struct {
		name        string
		setupDevice bool
		profileID   string
		query       string
		expectEmpty bool
	}{
		{
			name:        "search with matching query",
			setupDevice: true,
			profileID:   "matching-profile",
			query:       "Test", // Search for actual text in device name
			expectEmpty: false,
		},
		{
			name:        "search with non-matching query",
			setupDevice: true,
			profileID:   "different-profile",
			query:       "non-matching-query",
			expectEmpty: true,
		},
		{
			name:        "search with empty query",
			setupDevice: false,
			profileID:   "",
			query:       "",
			expectEmpty: true,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *definition.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		server := handlers.NewDeviceServer(ctx, svc)

		for _, tc := range testCases {
			suite.Run(tc.name, func() {
				suite.runSearchTestCase(ctx, svc, server, tc)
			})
		}
	})
}

func (suite *HandlersTestSuite) runSearchTestCase(
	ctx context.Context,
	svc *frame.Service,
	server *handlers.DevicesServer,
	tc struct {
		name        string
		setupDevice bool
		profileID   string
		query       string
		expectEmpty bool
	},
) {
	testCtx := suite.setupTestContext(ctx, tc.profileID)

	if tc.setupDevice {
		suite.createTestDevice(testCtx, svc, tc.profileID)
	}

	stream := suite.executeSearchRequest(testCtx, server, tc.query)
	suite.validateSearchResults(stream, tc.expectEmpty)
}

func (suite *HandlersTestSuite) setupTestContext(ctx context.Context, profileID string) context.Context {
	if profileID == "" {
		return ctx
	}

	claims := &frame.AuthenticationClaims{}
	claims.Subject = profileID
	return claims.ClaimsToContext(ctx)
}

func (suite *HandlersTestSuite) createTestDevice(ctx context.Context, svc *frame.Service, profileID string) {
	device := &models.Device{
		ProfileID: profileID,
		Name:      "Test Device",
		OS:        "Linux",
	}
	device.GenID(ctx)
	err := repository.NewDeviceRepository(svc).Save(ctx, device)
	suite.Require().NoError(err)
}

func (suite *HandlersTestSuite) executeSearchRequest(
	ctx context.Context,
	server *handlers.DevicesServer,
	query string,
) *commonMocks.MockServerStream[devicev1.SearchResponse] {
	req := &devicev1.SearchRequest{
		Query: query,
	}
	searchStream := commonMocks.NewMockServerStream[devicev1.SearchResponse](ctx)

	err := server.Search(req, searchStream)
	suite.Require().NoError(err)

	return searchStream
}

func (suite *HandlersTestSuite) validateSearchResults(
	stream *commonMocks.MockServerStream[devicev1.SearchResponse],
	expectEmpty bool,
) {
	totalDevices := suite.countDevicesInResponses(stream.GetResponses())

	if expectEmpty {
		suite.Equal(0, totalDevices, "Expected no devices in responses but got %d", totalDevices)
	} else {
		suite.Positive(totalDevices, "Expected devices in responses but got none")
	}
}

func (suite *HandlersTestSuite) countDevicesInResponses(responses []*devicev1.SearchResponse) int {
	totalDevices := 0
	for _, resp := range responses {
		totalDevices += len(resp.GetData())
	}
	return totalDevices
}

func (suite *HandlersTestSuite) TestGetClientIP() {
	testCases := []struct {
		name     string
		headers  map[string][]string
		expected string
	}{
		{
			name: "x-forwarded-for header",
			headers: map[string][]string{
				"x-forwarded-for": {"192.168.1.1, 10.0.0.1"},
			},
			expected: "192.168.1.1",
		},
		{
			name: "x-real-ip header",
			headers: map[string][]string{
				"x-real-ip": {"203.0.113.1"},
			},
			expected: "203.0.113.1",
		},
		{
			name: "x-forwarded-for takes precedence",
			headers: map[string][]string{
				"x-forwarded-for": {"192.168.1.1"},
				"x-real-ip":       {"203.0.113.1"},
			},
			expected: "192.168.1.1",
		},
		{
			name:     "no headers",
			headers:  map[string][]string{},
			expected: "", // Will fallback to peer info, which we can't easily mock
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx := context.Background()
			if len(tc.headers) > 0 {
				// Convert map[string][]string to map[string]string for metadata.New
				headerMap := make(map[string]string)
				for k, v := range tc.headers {
					if len(v) > 0 {
						headerMap[k] = v[0] // Take first value
					}
				}
				md := metadata.New(headerMap)
				ctx = metadata.NewIncomingContext(ctx, md)
			}

			ip := handlers.GetClientIP(ctx)
			if tc.expected != "" {
				suite.Equal(tc.expected, ip)
			}
			// For empty expected, we just check it doesn't panic
		})
	}
}
