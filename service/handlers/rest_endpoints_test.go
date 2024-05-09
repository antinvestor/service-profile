package handlers_test

import (
	"context"
	"crypto/sha256"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service/events"
	"github.com/antinvestor/service-profile/service/handlers"
	"github.com/gorilla/mux"
	"github.com/pitabwire/frame"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func getService() (context.Context, *frame.Service) {

	_ = os.Setenv("CONTACT_ENCRYPTION_KEY", "")
	_ = os.Setenv("CONTACT_ENCRYPTION_SALT", "")

	dbURL := frame.GetEnv("TEST_DATABASE_URL",
		"postgres://ant:secret@localhost:5434/service_profile?sslmode=disable")

	mainDB := frame.DatastoreCon(dbURL, false)

	var configProfile config.ProfileConfig
	err := frame.ConfigProcess("", &configProfile)
	if err != nil {
		logrus.WithError(err).Fatal("could not process configs")
	}

	verificationQueuePublisher := frame.RegisterPublisher(configProfile.QueueVerificationName, configProfile.QueueVerification)

	ctx, service := frame.NewService(
		"profile tests", mainDB,
		frame.Config(&configProfile), frame.NoopDriver())

	service.Init(verificationQueuePublisher,
		frame.RegisterEvents(
			&events.ClientConnectedSetupQueue{
				Service: service,
			},
		))

	_ = service.Run(ctx, "")
	return ctx, service
}

func getEncryptionKey() []byte {
	return pbkdf2.Key([]byte("ualgJEcb4GNXLn3jYV9TUGtgYrdTMg"), []byte("VufLmnycUCgz"), 4096, 32, sha256.New)
}

func TestRestUserInfoEndpoint(t *testing.T) {

	_, srv := getService()
	encKey := getEncryptionKey()

	ps := &handlers.ProfileServer{
		Service: srv,
		EncryptionKeyFunc: func() []byte {
			return encKey
		},
	}

	req, err := http.NewRequest("GET", "/user/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ps.RestUserInfo)

	handler.ServeHTTP(rr, req)

	//assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	//
	//var response map[string]interface{}
	//err = json.Unmarshal(rr.Body.Bytes(), &response)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//assert.Equal(t, "existing_profile_id", response["sub"], "expected profile ID not found")
	//assert.Equal(t, "Test User", response["name"], "expected profile name not found")
	//assert.Contains(t, response["contacts"], "email@example.com", "expected contact not found")
	//assert.Contains(t, response["contacts"], "phone:123456789", "expected contact not found")
	//assert.Equal(t, "https://example.com/profile_pic.jpg", response["url"], "expected profile pic URL not found")
}

func TestRestListRelationshipsEndpoint(t *testing.T) {
	_, srv := getService()
	encKey := getEncryptionKey()

	ps := &handlers.ProfileServer{
		Service: srv,
		EncryptionKeyFunc: func() []byte {
			return encKey
		},
	}
	req, err := http.NewRequest("GET", "/user/relations", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{
		"Count":          "10",
		"PeerObjectName": "TestObject",
		"PeerObjectID":   "existing_peer_id",
		"InvertRelation": "true",
	})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ps.RestListRelationshipsEndpoint)

	handler.ServeHTTP(rr, req)

	//assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	//
	//var response map[string]interface{}
	//err = json.Unmarshal(rr.Body.Bytes(), &response)
	//if err != nil {
	//	t.Fatal(err)
	//}

	//assert.Equal(t, "existing_peer_id", response["tenant_id"], "expected peer ID not found")
	//assert.Equal(t, "existing_peer_id", response["partition_id"], "expected peer ID not found")
	//assert.Equal(t, 2, response["count"], "expected count not found")
	//assert.NotEmpty(t, response["LastRelationshipID"], "expected LastRelationshipID not found")
}
