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

	// Applies authentication middleware on every API routes, for every API versions.
	// We may, in the future, use a different auth middleware for different API versions.
	authMW := newAuthMiddleware(s.Config.JWTAlgorithm, s.Config.JWTSecretKey)

	// Main APÌ routes group, versioned.
	api := s.router.Group(s.BaseURL, authMW.handlerFunc())
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
