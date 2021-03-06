package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/edebernis/social-life-manager/services/location/internal/api"
	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/edebernis/social-life-manager/services/location/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestLoggerMiddlewareWithOKRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()

	server := &HTTPServer{
		&Config{},
		"/api",
		nil,
		nil,
		gin.New(),
		nil,
	}

	server.router.Use(loggerMiddleware())

	server.router.GET("/testLoggerMiddleware", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest("GET", "/testLoggerMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	server.router.ServeHTTP(resp, req)
}

func TestLoggerMiddlewareWithErrorRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()

	server := &HTTPServer{
		&Config{},
		"/api",
		nil,
		nil,
		gin.New(),
		nil,
	}

	server.router.Use(loggerMiddleware())

	server.router.GET("/testLoggerMiddleware", func(c *gin.Context) {
		abort(c, http.StatusInternalServerError, "failed")
	})

	req, err := http.NewRequest("GET", "/testLoggerMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	server.router.ServeHTTP(resp, req)
}

func TestErrorMiddlewareWithOKRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()

	server := &HTTPServer{
		&Config{},
		"/api",
		nil,
		nil,
		gin.New(),
		nil,
	}

	server.router.Use(errorMiddleware())

	server.router.GET("/testErrorMiddleware", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest("GET", "/testErrorMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	server.router.ServeHTTP(resp, req)
}

func TestErrorMiddlewareWithErrorRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()

	server := &HTTPServer{
		&Config{},
		"/api",
		nil,
		nil,
		gin.New(),
		nil,
	}

	server.router.Use(errorMiddleware())

	server.router.GET("/testErrorMiddleware", func(c *gin.Context) {
		abort(c, http.StatusInternalServerError, "failed")
	})

	req, err := http.NewRequest("GET", "/testErrorMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	var httpError HTTPError
	err = json.NewDecoder(resp.Result().Body).Decode(&httpError)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)
		assert.Equal(t, "failed", httpError.Message)
	}
}

func TestRecoveryMiddlewareWithPanicRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()

	server := &HTTPServer{
		&Config{},
		"/api",
		nil,
		nil,
		gin.New(),
		nil,
	}

	server.router.Use(recoveryMiddleware())

	server.router.GET("/testRecoveryMiddleware", func(c *gin.Context) {
		panic(errors.New("Test Panic"))
	})

	req, err := http.NewRequest("GET", "/testRecoveryMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	assert.NotPanics(t, func() {
		server.router.ServeHTTP(resp, req)
	})
}

func newTestAuthMiddleware() (*HTTPServer, *JWTAuthenticator) {
	gin.SetMode(gin.TestMode)
	server := &HTTPServer{
		&Config{},
		"/api",
		nil,
		nil,
		gin.New(),
		nil,
	}

	auth := NewJWTAuthenticator(
		jwt.SigningMethodHS256.Name,
		"secret",
	)

	server.router.Use(authMiddleware(auth))

	return server, auth
}

func TestAuthMiddlewareWithoutToken(t *testing.T) {
	resp := httptest.NewRecorder()

	server, _ := newTestAuthMiddleware()

	server.router.GET("/testAuthMiddleware", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest("GET", "/testAuthMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthMiddlewareWithInvalidHeader(t *testing.T) {
	resp := httptest.NewRecorder()

	server, _ := newTestAuthMiddleware()

	server.router.GET("/testAuthMiddleware", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest("GET", "/testAuthMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	req.Header.Set("Authorization", "Invalid Header")

	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthMiddlewareWithInvalidSignatureToken(t *testing.T) {
	resp := httptest.NewRecorder()

	server, auth := newTestAuthMiddleware()

	server.router.GET("/testAuthMiddleware", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest("GET", "/testAuthMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	token, err := utils.NewJWTToken(auth.SecretKey, jwt.SigningMethodHS512.Name, nil)
	if err != nil {
		t.FailNow()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthMiddlewareWithoutUserIDInToken(t *testing.T) {
	resp := httptest.NewRecorder()

	server, auth := newTestAuthMiddleware()

	server.router.GET("/testAuthMiddleware", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest("GET", "/testAuthMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	token, err := utils.NewJWTToken(auth.SecretKey, auth.Algorithm, nil)
	if err != nil {
		t.FailNow()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthMiddlewareWithExpiredToken(t *testing.T) {
	resp := httptest.NewRecorder()

	server, auth := newTestAuthMiddleware()

	server.router.GET("/testAuthMiddleware", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest("GET", "/testAuthMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	claims := &api.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Add(time.Hour * -2).Unix(),
			ExpiresAt: time.Now().Add(time.Hour * -1).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "test",
			Subject:   models.NewID().String(),
		},
		Email: "test@no-reply.com",
	}

	token, err := utils.NewJWTToken(auth.SecretKey, auth.Algorithm, claims)
	if err != nil {
		t.FailNow()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthMiddlewareWithoutTokenSubject(t *testing.T) {
	resp := httptest.NewRecorder()

	server, auth := newTestAuthMiddleware()

	server.router.GET("/testAuthMiddleware", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest("GET", "/testAuthMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	claims := &api.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Add(time.Hour * -1).Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "test",
		},
		Email: "test@no-reply.com",
	}

	token, err := utils.NewJWTToken(auth.SecretKey, auth.Algorithm, claims)
	if err != nil {
		t.FailNow()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthMiddlewareWithValidToken(t *testing.T) {
	resp := httptest.NewRecorder()

	server, auth := newTestAuthMiddleware()

	userID := models.NewID()
	userEmail := "test@no-reply.com"

	server.router.GET("/testAuthMiddleware", func(c *gin.Context) {
		user, ok := models.NewUserFromContext(c.Request.Context())

		assert.True(t, ok)
		assert.Equal(t, user.ID, userID)
		assert.Equal(t, user.Email, userEmail)

		c.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest("GET", "/testAuthMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	claims := &api.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Add(time.Hour * -1).Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "test",
			Subject:   userID.String(),
		},
		Email: userEmail,
	}

	token, err := utils.NewJWTToken(auth.SecretKey, auth.Algorithm, claims)
	if err != nil {
		t.FailNow()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestMetricsMiddlewareRequestsCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()

	server := &HTTPServer{
		&Config{},
		"/api",
		prometheus.NewRegistry(),
		nil,
		gin.New(),
		nil,
	}

	mw := newMetricsMiddleware(server.PrometheusRegistry)
	server.router.Use(mw.handlerFunc())

	server.router.GET("/test200", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})
	server.router.GET("/test500", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, nil)
	})

	req, err := http.NewRequest("GET", "/test200", nil)
	if err != nil {
		t.FailNow()
	}
	req2, err := http.NewRequest("GET", "/test500", nil)
	if err != nil {
		t.FailNow()
	}

	server.router.ServeHTTP(resp, req)
	server.router.ServeHTTP(resp, req)
	server.router.ServeHTTP(resp, req2)

	err = testutil.CollectAndCompare(mw.requestCount, strings.NewReader(`
	# HELP api_http_requests_total How many HTTP requests processed, partitioned by status code and HTTP method.
	# TYPE api_http_requests_total counter
	api_http_requests_total{code="200",host="",method="GET",path="/test200"} 2
	api_http_requests_total{code="500",host="",method="GET",path="/test500"} 1
	`), "api_http_requests_total")

	assert.NoError(t, err)
}
