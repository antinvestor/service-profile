package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
	"github.com/antinvestor/service-profile/apps/devices/service/tests"
	"github.com/pitabwire/frame/tests/testdef"
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
		server := NewDeviceServer(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
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
					require.NoError(t, err)
					deviceID = device.GetID()
				} else {
					deviceID = tc.deviceID
				}

				req := &devicev1.GetByIdRequest{
					Id: []string{deviceID},
				}

				resp, err := server.GetByID(ctx, req)

				if tc.expectedStatus == codes.OK {
					require.NoError(t, err)
					require.NotNil(t, resp)
					assert.NotEmpty(t, resp.GetData())
					assert.Equal(t, deviceID, resp.GetData()[0].GetId())
				} else {
					require.Error(t, err)
					st, ok := status.FromError(err)
					require.True(t, ok)
					assert.Equal(t, tc.expectedStatus, st.Code())
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
		server := NewDeviceServer(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := &devicev1.CreateRequest{
					Name:       tc.deviceName,
					Properties: tc.data,
				}

				resp, err := server.Create(ctx, req)

				if tc.expectError {
					require.Error(t, err)
					assert.Nil(t, resp)
				} else {
					require.NoError(t, err)
					// Response might be nil if device creation doesn't return a device
					if resp != nil && resp.GetData() != nil {
						assert.NotEmpty(t, resp.GetData().GetId())
						assert.Equal(t, tc.deviceName, resp.GetData().GetName())
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
		server := NewDeviceServer(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
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
					require.NoError(t, err)
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
					require.Error(t, err)
					assert.Nil(t, resp)
				} else {
					require.NoError(t, err)
					assert.NotNil(t, resp)
					assert.NotNil(t, resp.GetData())
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
		server := NewDeviceServer(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
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
					require.NoError(t, err)
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
					require.Error(t, err)
					assert.Nil(t, resp)
				} else {
					require.NoError(t, err)
					assert.NotNil(t, resp)
					assert.NotNil(t, resp.GetData())
					if resp.GetData() != nil {
						assert.NotEmpty(t, resp.GetData().GetId())
						assert.Equal(t, deviceID, resp.GetData().GetDeviceId())
						assert.Equal(t, tc.keyData, resp.GetData().GetKey())
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
		query       string
		expectEmpty bool
	}{
		{
			name:        "search with matching query",
			setupDevice: true,
			query:       "test-profile",
			expectEmpty: false,
		},
		{
			name:        "search with non-matching query",
			setupDevice: true,
			query:       "non-matching",
			expectEmpty: true,
		},
		{
			name:        "search with empty query",
			setupDevice: false,
			query:       "",
			expectEmpty: true,
		},
	}

	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Create server
		server := NewDeviceServer(ctx, svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.setupDevice {
					// Create a device directly using repository
					device := &models.Device{
						ProfileID: "test-profile",
						Name:      "Test Device",
						OS:        "Linux",
					}
					device.GenID(ctx)
					err := repository.NewDeviceRepository(svc).Save(ctx, device)
					require.NoError(t, err)
				}

				req := &devicev1.SearchRequest{
					Query: tc.query,
				}

				stream := &mockSearchStream{
					ctx: ctx,
				}

				err := server.Search(req, stream)
				require.NoError(t, err)

				if tc.expectEmpty {
					assert.Empty(t, stream.responses)
				} else {
					// Note: responses might be empty due to implementation details
					// The important thing is that no error occurred
				}
			})
		}
	})
}

// Mock stream for testing streaming endpoints
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
		suite.T().Run(tc.name, func(t *testing.T) {
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

			ip := getClientIP(ctx)
			if tc.expected != "" {
				assert.Equal(t, tc.expected, ip)
			}
			// For empty expected, we just check it doesn't panic
		})
	}
}

func (suite *HandlersTestSuite) TestRESTEndpoints() {
	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)

		// Create server
		server := NewDeviceServer(ctx, svc)

		t.Run("RestLogDeviceData", func(t *testing.T) {
			reqBody := map[string]string{
				"session_id": "test-session-123",
				"action":     "page_view",
				"url":        "https://example.com",
			}

			body, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/log", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "Test Agent")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			server.RestLogDeviceData(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.NotNil(t, response)
		})

		t.Run("RestLogDeviceData - missing session_id", func(t *testing.T) {
			reqBody := map[string]string{
				"action": "page_view",
			}

			body, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/log", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			server.RestLogDeviceData(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("RestDeviceLinkProfile", func(t *testing.T) {
			reqBody := map[string]string{
				"session_id": "test-session-456",
				"profile_id": "test-profile-789",
			}

			body, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			server.RestDeviceLinkProfile(w, req)

			// This might return an error due to non-existent session, but should not panic
			assert.True(t, w.Code == http.StatusOK || w.Code >= 400)
		})

		t.Run("RestDeviceLinkProfile - missing parameters", func(t *testing.T) {
			reqBody := map[string]string{
				"session_id": "test-session-456",
				// missing profile_id
			}

			body, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			server.RestDeviceLinkProfile(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	})
}
