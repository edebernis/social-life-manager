package httpapi

import (
	"net/http"

	"github.com/edebernis/social-life-manager/services/location/internal/api"
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) routes(auth api.Authenticator) {
	// Default middlewares used on every routes
	s.router.Use(loggerMiddleware())
	s.router.Use(errorMiddleware())
	s.router.Use(recoveryMiddleware())
	s.router.Use(newMetricsMiddleware(s.PrometheusRegistry).handlerFunc())

	// Healthchecks routes
	s.router.GET("/ping", s.handlePing)

	// Main APÌ routes group, versioned.
	api := s.router.Group(s.BaseURL, authMiddleware(auth))
	{
		v1 := api.Group("/v1")
		{
			categories := v1.Group("/categories")
			{
				categories.POST("", s.handleCategoriesCreate)
				categories.GET("", s.handleCategoriesGet)
				categories.GET(":id", s.handleCategoriesGetByID)
				categories.PUT(":id", s.handleCategoriesUpdate)
				categories.DELETE(":id", s.handleCategoriesDelete)
			}

			locations := v1.Group("/locations")
			{
				locations.POST("", s.handleLocationsCreate)
				locations.GET("", s.handleLocationsGet)
				locations.GET(":id", s.handleLocationsGetByID)
				locations.PUT(":id", s.handleLocationsUpdate)
				locations.DELETE(":id", s.handleLocationsDelete)
			}
		}
	}
}

// handlePing godoc
// @Summary Ping API
// @Description Basic check of HTTP API health. Ensure that HTTP service is working correctly.
// @Tags healthchecks
// @Produce  json
// @Success 200
// @Router /ping [get]
func (s *HTTPServer) handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
