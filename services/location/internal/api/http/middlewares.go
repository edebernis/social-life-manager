package httpapi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/edebernis/social-life-manager/services/location/internal/models"
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
			return
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

// authenticationMiddleware handles request authentication
// and add user data into request context
type authenticationMiddleware struct {
	jwtAlgorithm string
	jwtSecretKey string
}

type userClaims struct {
	Email string `json:"email,omitempty"`
	jwt.StandardClaims
}

func newAuthMiddleware(jwtAlgorithm, jwtSecretKey string) *authenticationMiddleware {
	return &authenticationMiddleware{
		jwtAlgorithm,
		jwtSecretKey,
	}
}

func (mw *authenticationMiddleware) handlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := mw.parseJWTToken(c, &userClaims{})
		if err != nil {
			logger.Errorf("authenticationMiddleware: invalid JWT token. %v", err)
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("Invalid authentication token"))
			return
		}

		claims := token.Claims.(*userClaims)

		if err := mw.setContextDataFromClaims(c, claims); err != nil {
			logger.Errorf("authenticationMiddleware: failed to set context data from JWT token. %v", err)
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("Invalid authentication token"))
			return
		}

		c.Next()
	}
}

func (mw *authenticationMiddleware) setContextDataFromClaims(c *gin.Context, claims *userClaims) error {
	userID, err := models.ParseID(claims.StandardClaims.Subject)
	if err != nil {
		return fmt.Errorf("Invalid user ID in JWT token subject : %s. %w", claims.StandardClaims.Subject, err)
	}

	user := models.NewUser(userID, claims.Email)
	c.Set("user", user)

	return nil
}

func (mw *authenticationMiddleware) parseJWTToken(c *gin.Context, claims *userClaims) (*jwt.Token, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("Empty Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, fmt.Errorf("Invalid Authorization header : %s", authHeader)
	}

	encodedToken := parts[1]

	return jwt.ParseWithClaims(encodedToken, claims, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.jwtAlgorithm) != token.Method {
			return nil, fmt.Errorf("Invalid signing algorithm for JWT token : %s", token.Header["alg"])
		}
		return []byte(mw.jwtSecretKey), nil
	})
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
