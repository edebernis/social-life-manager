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
	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func newJWTToken(secret string, signingAlg string, claims jwt.Claims) (string, error) {
	alg := jwt.GetSigningMethod(signingAlg)
	token := jwt.NewWithClaims(alg, claims)

	out, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("Failed to sign JWT token: %v", err)
	}

	return out, nil
}

func TestLoggerMiddlewareWithOKRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()

	server := &HTTPServer{
		&Config{},
		"/api",
		nil,
		nil,
		gin.New(),
	}

	server.router.Use(loggerMiddleware())

	server.router.GET("/testLoggerMiddleware", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest("GET", "/testAuthMiddleware", nil)
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

func newTestAuthMiddleware() *HTTPServer {
	gin.SetMode(gin.TestMode)
	server := &HTTPServer{
		&Config{
			JWTAlgorithm: jwt.SigningMethodHS256.Name,
			JWTSecretKey: "secret",
		},
		"/api",
		nil,
		nil,
		gin.New(),
	}

	mw := newAuthMiddleware(server.Config.JWTAlgorithm, server.Config.JWTSecretKey)
	server.router.Use(mw.handlerFunc())

	return server
}

func TestAuthMiddlewareWithoutToken(t *testing.T) {
	resp := httptest.NewRecorder()

	server := newTestAuthMiddleware()

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

	server := newTestAuthMiddleware()

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

	server := newTestAuthMiddleware()

	server.router.GET("/testAuthMiddleware", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest("GET", "/testAuthMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	token, err := newJWTToken(server.Config.JWTSecretKey, jwt.SigningMethodHS512.Name, nil)
	if err != nil {
		t.FailNow()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthMiddlewareWithoutUserIDInToken(t *testing.T) {
	resp := httptest.NewRecorder()

	server := newTestAuthMiddleware()

	server.router.GET("/testAuthMiddleware", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest("GET", "/testAuthMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	token, err := newJWTToken(server.Config.JWTSecretKey, server.Config.JWTAlgorithm, nil)
	if err != nil {
		t.FailNow()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthMiddlewareWithExpiredToken(t *testing.T) {
	resp := httptest.NewRecorder()

	server := newTestAuthMiddleware()

	server.router.GET("/testAuthMiddleware", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest("GET", "/testAuthMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	claims := &userClaims{
		"test@no-reply.com",
		jwt.StandardClaims{
			NotBefore: time.Now().Add(time.Hour * -2).Unix(),
			ExpiresAt: time.Now().Add(time.Hour * -1).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "test",
			Subject:   models.NewID().String(),
		},
	}

	token, err := newJWTToken(server.Config.JWTSecretKey, server.Config.JWTAlgorithm, claims)
	if err != nil {
		t.FailNow()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthMiddlewareWithoutTokenSubject(t *testing.T) {
	resp := httptest.NewRecorder()

	server := newTestAuthMiddleware()

	server.router.GET("/testAuthMiddleware", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest("GET", "/testAuthMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	claims := &userClaims{
		"test@no-reply.com",
		jwt.StandardClaims{
			NotBefore: time.Now().Add(time.Hour * -1).Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "test",
		},
	}

	token, err := newJWTToken(server.Config.JWTSecretKey, server.Config.JWTAlgorithm, claims)
	if err != nil {
		t.FailNow()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthMiddlewareWithValidToken(t *testing.T) {
	resp := httptest.NewRecorder()

	server := newTestAuthMiddleware()

	userID := models.NewID()
	userEmail := "test@no-reply.com"

	server.router.GET("/testAuthMiddleware", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)

		user, exists := c.Get("user")

		assert.True(t, exists)
		assert.IsType(t, &models.User{}, user)

		userObj := user.(*models.User)

		assert.Equal(t, userObj.ID, userID)
		assert.Equal(t, userObj.Email, userEmail)
	})

	req, err := http.NewRequest("GET", "/testAuthMiddleware", nil)
	if err != nil {
		t.FailNow()
	}

	claims := &userClaims{
		userEmail,
		jwt.StandardClaims{
			NotBefore: time.Now().Add(time.Hour * -1).Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "test",
			Subject:   userID.String(),
		},
	}

	token, err := newJWTToken(server.Config.JWTSecretKey, server.Config.JWTAlgorithm, claims)
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
