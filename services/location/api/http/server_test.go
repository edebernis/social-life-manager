package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edebernis/social-life-manager/location/api"
	"github.com/edebernis/social-life-manager/location/api/mocks"
	"github.com/edebernis/social-life-manager/location/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	jwtSecretKey = "secretkey"
)

var (
	mockContextMatcher = mock.MatchedBy(func(ctx context.Context) bool { return true })
)

func newHandlerTestContext(t *testing.T, method, url string, payload *gin.H) (*gin.Context, *httptest.ResponseRecorder, *HTTPServer) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(resp)

	api := api.NewAPI(new(mocks.LocationUsecaseMock))
	server := &HTTPServer{
		&Config{},
		"/api",
		api,
		r,
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(payload)
	if err != nil {
		t.FailNow()
	}

	ctx.Request, err = http.NewRequest(method, url, &body)
	if err != nil {
		t.FailNow()
	}

	user := models.NewUser(models.NewID(), "testuser@no-reply.com")
	ctx.Set("user", user)

	return ctx, resp, server
}

func TestPingHTTPServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()

	api := api.NewAPI(new(mocks.LocationUsecaseMock))
	server := NewHTTPServer(api, &Config{
		JWTSecretKey: "",
	})

	req, _ := http.NewRequest("GET", "/ping", nil)

	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHTTPServerUnknownRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()

	api := api.NewAPI(new(mocks.LocationUsecaseMock))
	server := NewHTTPServer(api, &Config{
		JWTSecretKey: "",
	})

	req, _ := http.NewRequest("GET", "/pong", nil)

	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}