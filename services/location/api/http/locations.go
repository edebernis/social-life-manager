package httpapi

import (
	"net/http"

	"github.com/edebernis/social-life-manager/location/models"
	"github.com/edebernis/social-life-manager/location/usecases"
	"github.com/gin-gonic/gin"
)

// handleLocationsCreate godoc
// @Summary Create locations
// @Description Create new user locations.
// @Tags locations
// @Accept  json
// @Produce  json
// @Param location body models.CreateLocation true "New location"
// @Success 200 {object} models.Location
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /locations [post]
func (s *HTTPServer) handleLocationsCreate(c *gin.Context) {
	var loc models.CreateLocation
	if err := c.ShouldBindJSON(&loc); err != nil {
		logger.Errorf("createLocation: invalid body. %v", err)
		abort(c, http.StatusBadRequest, "Invalid body")
		return
	}

	catID, err := models.ParseID(loc.CategoryID)
	if err != nil {
		logger.Errorf("createLocation: invalid category ID : %s. %v", loc.CategoryID, err)
		abort(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	newLoc := models.NewLocation(models.NewID(), loc.Name, loc.Address, catID, models.NewUserFromContext(c).ID)

	err = s.api.LocationUsecase.CreateLocation(c.Request.Context(), newLoc)
	if err != nil {
		switch err {
		case usecases.ErrLocationAlreadyExists:
			abort(c, http.StatusBadRequest, "Location already exists")
			return
		default:
			logger.Errorf("createLocation: failed to create location. %v", err)
			abort(c, http.StatusInternalServerError, "Failed to create location")
			return
		}
	}

	c.JSON(http.StatusOK, newLoc)
}

// handleLocationsGet godoc
// @Summary Get locations
// @Description Get all user locations.
// @Tags locations
// @Produce  json
// @Success 200 {object} models.Locations
// @Failure 500 {object} HTTPError
// @Router /locations [get]
func (s *HTTPServer) handleLocationsGet(c *gin.Context) {
	locations, err := s.api.LocationUsecase.GetLocations(c.Request.Context())
	if err != nil {
		logger.Errorf("getLocations: failed to get locations. %v", err)
		abort(c, http.StatusInternalServerError, "Failed to get locations")
		return
	}

	c.JSON(http.StatusOK, locations)
}
