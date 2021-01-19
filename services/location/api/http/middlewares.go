package httpapi

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/edebernis/social-life-manager/location/models"
	"github.com/gin-gonic/gin"
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
