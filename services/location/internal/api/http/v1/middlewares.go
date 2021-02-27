package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/edebernis/social-life-manager/services/location/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func requestLogger(c *gin.Context, t time.Time) *logrus.Entry {
	return logrus.
		WithContext(c).
		WithFields(logrus.Fields{
			"package":   "httpapi",
			"latency":   time.Since(t).Milliseconds(),
			"status":    c.Writer.Status(),
			"path":      c.Request.URL.Path,
			"remote":    c.ClientIP(),
			"useragent": c.Request.UserAgent(),
		})
}

// loggerMiddleware handles all logging for each incoming requests.
// One log entry is generated for each request and additional entries
// may be logged in case of errors
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()

		if err := c.Errors.Last(); err != nil {
			requestLogger(c, t).Error(err)
		} else {
			requestLogger(c, t).Info()
		}
	}
}

// errorMiddleware handles error that happened during request processing.
// It returns a consistent message to the user describing the last error.
func errorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if err := c.Errors.Last(); err != nil {
			newError(c, c.Writer.Status(), err)
		}
	}
}

// recoveryMiddleware handles panics occurring during request handling.
// TODO: Replace with gin.CustomRecoveryWithWriter to specify custom handler
// to format HTTP Response as an HTTPError struct
func recoveryMiddleware() gin.HandlerFunc {
	return gin.RecoveryWithWriter(
		logrus.StandardLogger().WriterLevel(logrus.ErrorLevel),
	)
}

func authMiddleware(auth api.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		creds, err := auth.CredentialsFromContext(c)
		if err != nil {
			logger.Errorf("authMiddleware: unable to retrieve creds from context. %v", err)
			abort(c, http.StatusUnauthorized, "invalid auth")
			return
		}

		newCtx, err := auth.Authenticate(c, creds)
		if err != nil {
			logger.Errorf("authMiddleware: unable to authenticate user using provided credentials. %v", err)
			abort(c, http.StatusUnauthorized, "invalid auth")
			return
		}

		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}

// JWTAuthenticator is a middleware authenticating user using JWT tokens
type JWTAuthenticator struct {
	*api.JWTAuthenticator
}

// NewJWTAuthenticator creates a new JWTAuthenticator
func NewJWTAuthenticator(algorithm, secretKey string) *JWTAuthenticator {
	return &JWTAuthenticator{
		&api.JWTAuthenticator{
			Algorithm: algorithm,
			SecretKey: secretKey,
		},
	}
}

// CredentialsFromContext extracts JWT token from request headers
func (a *JWTAuthenticator) CredentialsFromContext(ctx context.Context) (interface{}, error) {
	c, ok := ctx.(*gin.Context)
	if !ok {
		return nil, errors.New("invalid context")
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("empty Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, fmt.Errorf("invalid Authorization header : %s", authHeader)
	}

	return parts[1], nil
}

// metricsMiddleware collects metrics about HTTP requests
type metricsMiddleware struct {
	requestCount    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	responseSize    prometheus.Summary

	namespace string
	subsystem string
}

func newMetricsMiddleware(registry prometheus.Registerer) *metricsMiddleware {
	mw := &metricsMiddleware{
		namespace: "api",
		subsystem: "http",
	}

	mw.requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: mw.namespace,
			Subsystem: mw.subsystem,
			Name:      "requests_total",
			Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method", "host", "path"},
	)
	registry.MustRegister(mw.requestCount)

	mw.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: mw.namespace,
			Subsystem: mw.subsystem,
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			Name:      "request_duration_seconds",
			Help:      "The HTTP request latency bucket",
		}, []string{"method", "path"},
	)
	registry.MustRegister(mw.requestDuration)

	mw.responseSize = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: mw.namespace,
			Subsystem: mw.subsystem,
			Name:      "response_size_bytes",
			Help:      "The HTTP response sizes in bytes.",
		},
	)
	registry.MustRegister(mw.responseSize)

	return mw
}

func (mw *metricsMiddleware) handlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		responseSize := float64(c.Writer.Size())

		mw.requestCount.WithLabelValues(status, c.Request.Method, c.Request.Host, c.Request.URL.Path).Inc()
		mw.requestDuration.WithLabelValues(c.Request.Method, c.Request.URL.Path).Observe(elapsed)
		mw.responseSize.Observe(responseSize)
	}
}
