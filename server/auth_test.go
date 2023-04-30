package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"skybluetrades.net/work-planning-demo/api"
	"skybluetrades.net/work-planning-demo/mocks"
	"skybluetrades.net/work-planning-demo/model"
	"skybluetrades.net/work-planning-demo/store"
)

const (
	adminEmail    = "worker1@test.com"
	adminName     = "worker1"
	adminPassword = "worker1pass"
)

func setupTestData(db *mocks.Store) {
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), 0)
	worker1 := model.Worker{
		ID:       1,
		Email:    adminEmail,
		Name:     adminName,
		IsAdmin:  true,
		Password: string(bcryptPassword),
	}
	db.
		On("Authenticate", adminEmail, adminPassword).Return(&worker1, nil).
		On("Authenticate", mock.Anything, mock.Anything).Return(nil, store.ErrWorkerNotFound)
	db.
		On("GetWorkerById", model.WorkerID(1)).Return(&worker1, nil).
		On("GetWorkerById", mock.Anything).Return(nil, store.ErrWorkerNotFound)
	db.
		On("GetWorkers").Return([]*model.Worker{&worker1}, nil)
}

func serverSetup(t *testing.T, testData bool) (*httpexpect.Expect, *httptest.Server) {
	// Test server configuration: we're going to inject a mock store
	// layer, and we set the token leases very short (1 second for the
	// access token and 60 seconds for the refresh token), so that we
	// can test token expiry handling in a reasonable time.
	cfg := &Config{
		DevMode:           true,
		StoreURL:          "mock",
		Port:              8080,
		AccessTokenLease:  1,  // second
		RefreshTokenLease: 60, // second
		AuthKey:           "test-key",
	}
	db := &mocks.Store{}
	serv := NewServer(cfg, db, nil)

	if testData {
		setupTestData(db)
	}

	srv := httptest.NewServer(serv)
	e := httpexpect.New(t, srv.URL)

	return e, srv
}

func TestAuthLogin(t *testing.T) {
	e, srv := serverSetup(t, true)
	defer srv.Close()

	// No login data in request body.
	e.POST("/auth/login").
		Expect().Status(http.StatusBadRequest).JSON().Object().
		HasValue("message", "request body has an error: value is required but missing")

	// Invalid login data.
	e.POST("/auth/login").WithJSON(api.Login{Email: "noname", Password: "nothing"}).
		Expect().Status(http.StatusForbidden).JSON().Object().
		HasValue("message", "Invalid login credentials")

	// Good login data.
	e.POST("/auth/login").WithJSON(api.Login{Email: adminEmail, Password: adminPassword}).
		Expect().Status(http.StatusOK).JSON().Object().
		ContainsKey("access_token").ContainsKey("refresh_token")
}

func TestAuthMiddleware(t *testing.T) {
	e, srv := serverSetup(t, true)
	defer srv.Close()

	// Missing Authorization header.
	e.GET("/worker").
		Expect().Status(http.StatusForbidden).JSON().Object().
		HasValue("message", "security requirements failed: no Authorization header")

	// Request with invalid token.
	e.GET("/worker").WithHeader("Authorization", "Bearer bad-token").
		Expect().Status(http.StatusForbidden).JSON().Object().
		HasValue("message", "security requirements failed: token contains an invalid number of segments")

	// Request with valid token.
	accessToken, refreshToken := getTokens(e)
	e.GET("/worker").WithHeader("Authorization", "Bearer "+accessToken).
		Expect().Status(http.StatusOK).JSON().Array()

	// Wait for the token to expire and try again. (Lease is set to 1
	// second for testing.)
	time.Sleep(2 * time.Second)
	e.GET("/worker").WithHeader("Authorization", "Bearer "+accessToken).
		Expect().Status(http.StatusForbidden).JSON().Object().
		Value("message").String().HasPrefix("security requirements failed: token is expired by ")

	// Refresh the token.
	ro := e.POST("/auth/refresh_token").
		WithJSON(api.CredentialsRefresh{RefreshToken: refreshToken}).
		Expect().Status(http.StatusOK).JSON().Object()
	accessToken = ro.Value("access_token").String().Raw()
	refreshToken = ro.Value("refresh_token").String().Raw()

	// Try again with new access token.
	e.GET("/worker").WithHeader("Authorization", "Bearer "+accessToken).
		Expect().Status(http.StatusOK).JSON().Array()
}

// Helper to get API tokens for tests.
func getTokens(e *httpexpect.Expect) (string, string) {
	login := &api.Login{Email: adminEmail, Password: adminPassword}
	ro := e.POST("/auth/login").WithJSON(login).
		Expect().Status(http.StatusOK).
		JSON().Object()
	accessToken := ro.Value("access_token").String().Raw()
	refreshToken := ro.Value("refresh_token").String().Raw()
	return accessToken, refreshToken
}
