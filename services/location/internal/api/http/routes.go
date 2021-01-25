package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) routes() {
	// Default middlewares used on every routes
	s.router.Use(loggerMiddleware())
	s.router.Use(errorMiddleware())
	s.router.Use(recoveryMiddleware())
	s.router.Use(newMetricsMiddleware(s.PrometheusRegistry).handlerFunc())

	// Healthchecks routes
	s.router.GET("/ping", s.handlePing)

	// Main APÌ routes group, versioned.
	api := s.router.Group(s.BaseURL)
	{
		// Applies authentication middleware on every API routes, for every API versions.
		// We may, in the future, use a different auth middleware for different API versions.
		authMW := newAuthMiddleware(s.Config.JWTAlgorithm, s.Config.JWTSecretKey)
		api.Use(authMW.handlerFunc())

		v1 := api.Group("/v1")
		{
			locations := v1.Group("/locations")
			{
				locations.GET("", s.handleLocationsGet)
				locations.POST("", s.handleLocationsCreate)
			}
		}
	}
}

// handlePing godoc
// @Summary Ping API
// @Description Basic check of HTTP API health. Ensure that HTTP serving is working correctly.
// @Tags healthchecks
// @Produce  json
// @Success 200
// @Router /ping [get]
func (s *HTTPServer) handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
