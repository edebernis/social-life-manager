package httpapi

import (
	"errors"

	"github.com/edebernis/social-life-manager/location/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.WithFields(logrus.Fields{"package": "httpapi"})
)

// @title Locations Service REST API
// @version 1.0
// @description This REST API handles management of user locations. Locations can be saved in a local repository
// @description or fetched from third-party sources such as Google Maps "My Places".

// @contact.name Emeric de Bernis
// @contact.email emeric.debernis@gmail.com

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// HTTPServer instance. Serves HTTP endpoints for the API.
type HTTPServer struct {
	// Advanced configuration, retrieved at app startup.
	Config *Config
	// Base URL of the API, like "/api". Must starts with a slash.
	BaseURL string
	api     *api.API
	router  *gin.Engine
}

// Serve runs the server. This method blocks the current goroutine
func (s *HTTPServer) Serve(addr string) error {
	return s.router.Run(addr)
}

// NewHTTPServer creates a new HTTP server for the API
func NewHTTPServer(api *api.API, config *Config) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)

	s := &HTTPServer{
		config,
		"/api",
		api,
		gin.New(),
	}
	s.routes()

	return s
}

// Config advanced settings of HTTP API server
type Config struct {
	// Signing algorithm used for JWT
	JWTAlgorithm string
	// Key to check JWT signature
	JWTSecretKey string
}

// Abort current request and return consistent error to the user
func abort(c *gin.Context, code int, errorMsg string) {
	c.AbortWithError(code, errors.New(errorMsg))
}
