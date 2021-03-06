package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/edebernis/social-life-manager/services/location/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
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
	BaseURL            string
	PrometheusRegistry prometheus.Registerer

	api    *api.API
	router *gin.Engine
	server *http.Server
}

// Serve runs the server. This method blocks the current goroutine.
func (s *HTTPServer) Serve(addr string) error {
	s.server.Addr = addr
	return s.server.ListenAndServe()
}

// Shutdown stops the server gracefully
func (s *HTTPServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	return nil
}

// NewHTTPServer creates a new HTTP server for the API
func NewHTTPServer(api *api.API, auth api.Authenticator, registry prometheus.Registerer, config *Config) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	s := &HTTPServer{
		config,
		"/api",
		registry,
		api,
		router,
		&http.Server{
			Handler:           router,
			ReadHeaderTimeout: config.ReadHeaderTimeout,
			ReadTimeout:       config.ReadTimeout,
			WriteTimeout:      config.WriteTimeout,
		},
	}
	s.routes(auth)

	return s
}

// Config advanced settings of HTTP API server
type Config struct {
	// HTTP server timeouts
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
}

// Abort current request and return consistent error to the user
func abort(c *gin.Context, code int, errorMsg string) {
	c.AbortWithStatusJSON(code, HTTPError{code, errorMsg})
}
