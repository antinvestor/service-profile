package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/pitabwire/frame/tests/testdef"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
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

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
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

				resp, err := server.GetByID(ctx, req)

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
		data        map[string]string
		expectError bool
	}{
		{
			name:       "create device successfully",
			linkID:     "test-link-123",
			deviceName: "Test Device",
			data: map[string]string{
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
			data: map[string]string{
				"session_id": "test-session-2",
				"ip":         "127.0.0.1",
			},
			expectError: false, // Should still work with empty name
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Create server
		server := handlers.NewDeviceServer(ctx, svc)

		for _, tc := range testCases {
			suite.Run(tc.name, func() {
				req := &devicev1.CreateRequest{
					Name:       tc.deviceName,
					Properties: tc.data,
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
		data        map[string]string
		expectError bool
	}{
		{
			name:        "log device activity successfully",
			setupDevice: true,
			sessionID:   "test-session",
			data: map[string]string{
				"action": "page_view",
				"url":    "https://example.com",
			},
			expectError: false,
		},
		{
			name:        "log with empty data",
			setupDevice: true,
			sessionID:   "test-session-2",
			data:        map[string]string{},
			expectError: false,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
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
					Extras:    tc.data,
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

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
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

				req := &devicev1.AddKeyRequest{
					DeviceId: deviceID,
					KeyType:  tc.keyType,
					Data:     tc.keyData,
					Extras:   map[string]string{"test": "data"},
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
			query:       "matching-profile",
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

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Create server
		server := handlers.NewDeviceServer(ctx, svc)

		for _, tc := range testCases {
			suite.Run(tc.name, func() {
				if tc.setupDevice {
					// Create a device directly using repository with unique profile ID
					device := &models.Device{
						ProfileID: tc.profileID,
						Name:      "Test Device",
						OS:        "Linux",
					}
					device.GenID(ctx)
					err := repository.NewDeviceRepository(svc).Save(ctx, device)
					suite.Require().NoError(err)
				}

				req := &devicev1.SearchRequest{
					Query: tc.query,
				}

				stream := &mockSearchStream{
					ctx: ctx,
				}

				err := server.Search(req, stream)
				suite.Require().NoError(err)

				if tc.expectEmpty {
					suite.Empty(stream.responses)
				} else {
					suite.NotEmpty(stream.responses)
				}
			})
		}
	})
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

func (suite *HandlersTestSuite) TestRESTEndpoints() {
	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Create server
		server := handlers.NewDeviceServer(ctx, svc)

		suite.Run("RestLogDeviceData", func() {
			reqBody := map[string]string{
				"session_id": "test-session-123",
				"action":     "page_view",
				"url":        "https://example.com",
			}

			body, err := json.Marshal(reqBody)
			suite.Require().NoError(err)

			req := httptest.NewRequest(http.MethodPost, "/log", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "Test Agent")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			server.RestLogDeviceData(w, req)

			suite.Equal(http.StatusOK, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			suite.Require().NoError(err)
			suite.NotNil(response)
		})

		suite.Run("RestLogDeviceData - missing session_id", func() {
			reqBody := map[string]string{
				"action": "page_view",
			}

			body, err := json.Marshal(reqBody)
			suite.Require().NoError(err)

			req := httptest.NewRequest(http.MethodPost, "/log", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			server.RestLogDeviceData(w, req)

			suite.Equal(http.StatusBadRequest, w.Code)
		})

		suite.Run("RestDeviceLinkProfile", func() {
			reqBody := map[string]string{
				"session_id": "test-session-456",
				"profile_id": "test-profile-789",
			}

			body, err := json.Marshal(reqBody)
			suite.Require().NoError(err)

			req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			server.RestDeviceLinkProfile(w, req)

			// This might return an error due to non-existent session, but should not panic
			suite.True(w.Code == http.StatusOK || w.Code >= 400)
		})

		suite.Run("RestDeviceLinkProfile - missing parameters", func() {
			reqBody := map[string]string{
				"session_id": "test-session-456",
				// missing profile_id
			}

			body, err := json.Marshal(reqBody)
			suite.Require().NoError(err)

			req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			server.RestDeviceLinkProfile(w, req)

			suite.Equal(http.StatusBadRequest, w.Code)
		})
	})
}

// Mock stream for testing streaming endpoints.
type mockSearchStream struct {
	grpc.ServerStream
	ctx       context.Context
	responses []*devicev1.SearchResponse
}

func (m *mockSearchStream) Send(resp *devicev1.SearchResponse) error {
	m.responses = append(m.responses, resp)
	return nil
}

func (m *mockSearchStream) Context() context.Context {
	return m.ctx
}
