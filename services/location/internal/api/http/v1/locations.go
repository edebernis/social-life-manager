package httpapi

import (
	"net/http"

	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/edebernis/social-life-manager/services/location/internal/usecases"
	"github.com/gin-gonic/gin"
)

// handleCategoriesCreate godoc
// @Summary Create categories
// @Description Create new categories.
// @Tags categories
// @Accept  json
// @Produce  json
// @Param category body models.CreateCategory true "New category"
// @Success 200 {object} models.Category "The created category"
// @Failure 400 {object} HTTPError "Bad Request"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /categories [post]
func (s *HTTPServer) handleCategoriesCreate(c *gin.Context) {
	var body models.CreateCategory
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Errorf("CategoriesCreate: invalid body. %v", err)
		abort(c, http.StatusBadRequest, "Invalid body")
		return
	}

	cat := models.NewCategory(models.NewID(), body.Name)

	err := s.api.LocationUsecase.CreateCategory(c.Request.Context(), cat)
	if err != nil {
		switch err {
		case usecases.ErrCategoryAlreadyExists:
			abort(c, http.StatusBadRequest, "Category already exists")
			return
		default:
			logger.Errorf("CategoriesCreate: failed to create category. %v", err)
			abort(c, http.StatusInternalServerError, "Failed to create category")
			return
		}
	}

	c.JSON(http.StatusOK, cat)
}

// handleCategoriesGet godoc
// @Summary Get categories
// @Description Get all categories.
// @Tags categories
// @Produce  json
// @Success 200 {object} models.Categories "The returned categories"
// @Failure 400 {object} HTTPError "Bad Request"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /categories [get]
func (s *HTTPServer) handleCategoriesGet(c *gin.Context) {
	cats, err := s.api.LocationUsecase.GetCategories(c.Request.Context())

	switch {
	case err != nil:
		logger.Errorf("CategoriesGet: failed to get categories. %v", err)
		abort(c, http.StatusInternalServerError, "Failed to get categories")
		return
	default:
		c.JSON(http.StatusOK, cats)
	}

}

// handleCategoriesGetByID godoc
// @Summary Get category with specified ID
// @Description Get one specific category using provided ID.
// @Tags categories
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {object} models.Category "The returned category"
// @Failure 400 {object} HTTPError "Bad Request"
// @Failure 404 {object} HTTPError "Not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /categories/{id} [get]
func (s *HTTPServer) handleCategoriesGetByID(c *gin.Context) {
	var query models.GetCategoryByID
	if err := c.ShouldBindUri(&query); err != nil {
		logger.Errorf("CategoriesGetByID: invalid query. %v", err)
		abort(c, http.StatusBadRequest, "Invalid query")
		return
	}

	id, err := models.ParseID(query.ID)
	if err != nil {
		logger.Errorf("CategoriesGetByID: invalid query. %v", err)
		abort(c, http.StatusBadRequest, "Invalid query")
		return
	}

	cat, err := s.api.LocationUsecase.FindCategoryByID(c.Request.Context(), id)
	switch {
	case err == usecases.ErrCategoryNotFound:
		abort(c, http.StatusNotFound, "Category not found")
		return
	case err != nil:
		logger.Errorf("CategoriesGetByID: failed to get category %s. %v", query.ID, err)
		abort(c, http.StatusInternalServerError, "Failed to get category")
		return
	default:
		c.JSON(http.StatusOK, cat)
	}
}

// handleCategoriesUpdate godoc
// @Summary Update category
// @Description Update specified category using provided values.
// @Tags categories
// @Produce  json
// @Param id path string true "Category ID"
// @Param name body models.UpdateCategory false "Category name"
// @Success 200 {object} models.Category "The updated category"
// @Failure 400 {object} HTTPError "Bad Request"
// @Failure 404 {object} HTTPError "Not Found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /categories/{id} [put]
func (s *HTTPServer) handleCategoriesUpdate(c *gin.Context) {
	var query models.UpdateCategoryQuery
	if err := c.ShouldBindUri(&query); err != nil {
		logger.Errorf("CategoriesUpdate: invalid query. %v", err)
		abort(c, http.StatusBadRequest, "Invalid query")
		return
	}

	var body models.UpdateCategory
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Errorf("CategoriesUpdate: invalid body. %v", err)
		abort(c, http.StatusBadRequest, "Invalid body")
		return
	}

	id, err := models.ParseID(query.ID)
	if err != nil {
		logger.Errorf("CategoriesUpdate: invalid query. %v", err)
		abort(c, http.StatusBadRequest, "Invalid query")
		return
	}

	cat := models.NewCategory(id, body.Name)

	err = s.api.LocationUsecase.UpdateCategory(c.Request.Context(), cat)
	switch {
	case err == usecases.ErrCategoryNotFound:
		abort(c, http.StatusNotFound, "Category not found")
		return
	case err != nil:
		logger.Errorf("CategoriesUpdate: failed to update category %s. %v", id, err)
		abort(c, http.StatusInternalServerError, "Failed to update category")
		return
	default:
		c.JSON(http.StatusOK, cat)
	}
}

// handleCategoriesDelete godoc
// @Summary Delete category
// @Description Delete one specific category using provided ID.
// @Tags categories
// @Param id path string true "Category ID"
// @Success 204 "OK"
// @Failure 400 {object} HTTPError "Bad Request"
// @Failure 404 {object} HTTPError "Not Found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /categories/{id} [delete]
func (s *HTTPServer) handleCategoriesDelete(c *gin.Context) {
	var query models.DeleteCategory
	if err := c.ShouldBindUri(&query); err != nil {
		logger.Errorf("CategoriesDelete: invalid query. %v", err)
		abort(c, http.StatusBadRequest, "Invalid query")
		return
	}

	id, err := models.ParseID(query.ID)
	if err != nil {
		logger.Errorf("CategoriesDelete: invalid query. %v", err)
		abort(c, http.StatusBadRequest, "Invalid query")
		return
	}

	err = s.api.LocationUsecase.DeleteCategory(c.Request.Context(), id)
	switch {
	case err == usecases.ErrCategoryNotFound:
		abort(c, http.StatusNotFound, "Category not found")
		return
	case err != nil:
		logger.Errorf("CategoriesDelete: failed to delete category %s. %v", id, err)
		abort(c, http.StatusInternalServerError, "Failed to delete category")
		return
	default:
		c.Status(http.StatusNoContent)
	}
}

// handleLocationsCreate godoc
// @Summary Create locations
// @Description Create new user locations.
// @Tags locations
// @Accept  json
// @Produce  json
// @Param location body models.CreateLocation true "New location"
// @Success 200 {object} models.Location "The created location"
// @Failure 400 {object} HTTPError "Bad Request"
// @Failure 404 {object} HTTPError "Not Found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /locations [post]
func (s *HTTPServer) handleLocationsCreate(c *gin.Context) {
	var body models.CreateLocation
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Errorf("LocationsCreate: invalid body. %v", err)
		abort(c, http.StatusBadRequest, "Invalid body")
		return
	}

	user, ok := models.NewUserFromContext(c.Request.Context())
	if !ok {
		logger.Error("LocationsCreate: failed to get user from request context.")
		abort(c, http.StatusInternalServerError, "Failed to get user data")
		return
	}

	loc := models.NewLocation(models.NewID(), body.Name, body.Address, body.Category, user.ID)

	err := s.api.LocationUsecase.CreateLocation(c.Request.Context(), loc)
	if err != nil {
		switch err {
		case usecases.ErrLocationAlreadyExists:
			abort(c, http.StatusBadRequest, "Location already exists")
			return
		case usecases.ErrCategoryNotFound:
			abort(c, http.StatusNotFound, "Category not found")
			return
		default:
			logger.Errorf("LocationsCreate: failed to create location. %v", err)
			abort(c, http.StatusInternalServerError, "Failed to create location")
			return
		}
	}

	c.JSON(http.StatusOK, loc)
}

// handleLocationsGet godoc
// @Summary Get locations
// @Description Get all user locations.
// @Tags locations
// @Produce  json
// @Param category_id query string false "Category ID"
// @Success 200 {object} models.Locations "The returned locations"
// @Failure 400 {object} HTTPError "Bad Request"
// @Failure 404 {object} HTTPError "Not Found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /locations [get]
func (s *HTTPServer) handleLocationsGet(c *gin.Context) {
	var query models.GetLocations
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Errorf("LocationsGet: invalid query. %v", err)
		abort(c, http.StatusBadRequest, "Invalid query")
		return
	}

	catID, err := models.ParseID(query.Category)
	if err != nil {
		logger.Errorf("LocationsGet: invalid query. %v", err)
		abort(c, http.StatusBadRequest, "Invalid query")
		return
	}

	var locations *models.Locations
	if catID != models.NilID {
		locations, err = s.api.LocationUsecase.FindLocationsByCategory(c.Request.Context(), catID)
	} else {
		locations, err = s.api.LocationUsecase.GetLocations(c.Request.Context())
	}

	switch {
	case err == usecases.ErrCategoryNotFound:
		abort(c, http.StatusNotFound, "Category not found")
		return
	case err != nil:
		logger.Errorf("LocationsGet: failed to get locations. %v", err)
		abort(c, http.StatusInternalServerError, "Failed to get locations")
		return
	default:
		c.JSON(http.StatusOK, locations)
	}

}

// handleLocationsGetByID godoc
// @Summary Get location with specified ID
// @Description Get one specific location using provided ID.
// @Tags locations
// @Produce  json
// @Param id path string true "Location ID"
// @Success 200 {object} models.Location "The returned location"
// @Failure 400 {object} HTTPError "Bad Request"
// @Failure 404 {object} HTTPError "Not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /locations/{id} [get]
func (s *HTTPServer) handleLocationsGetByID(c *gin.Context) {
	var query models.GetLocationByID
	if err := c.ShouldBindUri(&query); err != nil {
		logger.Errorf("LocationsGetByID: invalid query. %v", err)
		abort(c, http.StatusBadRequest, "Invalid query")
		return
	}

	id, err := models.ParseID(query.ID)
	if err != nil {
		logger.Errorf("LocationsGetByID: invalid query. %v", err)
		abort(c, http.StatusBadRequest, "Invalid query")
		return
	}

	loc, err := s.api.LocationUsecase.FindLocationByID(c.Request.Context(), id)
	switch {
	case err == usecases.ErrLocationNotFound:
		abort(c, http.StatusNotFound, "Location not found")
		return
	case err != nil:
		logger.Errorf("LocationsGetByID: failed to get location %s. %v", query.ID, err)
		abort(c, http.StatusInternalServerError, "Failed to get location")
		return
	default:
		c.JSON(http.StatusOK, loc)
	}
}

// handleLocationsUpdate godoc
// @Summary Update location
// @Description Update specified location using provided values.
// @Tags locations
// @Produce  json
// @Param id path string true "Location ID"
// @Param location body models.UpdateLocation true "Updated location"
// @Success 200 {object} models.Location "The updated location"
// @Failure 400 {object} HTTPError "Bad Request"
// @Failure 404 {object} HTTPError "Not Found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /locations/{id} [put]
func (s *HTTPServer) handleLocationsUpdate(c *gin.Context) {
	var query models.UpdateLocationQuery
	if err := c.ShouldBindUri(&query); err != nil {
		logger.Errorf("LocationsUpdate: invalid query. %v", err)
		abort(c, http.StatusBadRequest, "Invalid query")
		return
	}

	var body models.UpdateLocation
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Errorf("LocationsUpdate: invalid body. %v", err)
		abort(c, http.StatusBadRequest, "Invalid body")
		return
	}

	id, err := models.ParseID(query.ID)
	if err != nil {
		logger.Errorf("LocationsUpdate: invalid query. %v", err)
		abort(c, http.StatusBadRequest, "Invalid query")
		return
	}

	user, ok := models.NewUserFromContext(c.Request.Context())
	if !ok {
		logger.Errorf("LocationsUpdate: failed to get user from context")
		abort(c, http.StatusInternalServerError, "Failed to update location")
		return
	}

	loc := models.NewLocation(id, body.Name, body.Address, body.Category, user.ID)

	err = s.api.LocationUsecase.UpdateLocation(c.Request.Context(), loc)
	switch {
	case err == usecases.ErrLocationNotFound:
		abort(c, http.StatusNotFound, "Location not found")
		return
	case err == usecases.ErrCategoryNotFound:
		abort(c, http.StatusNotFound, "Category not found")
		return
	case err != nil:
		logger.Errorf("LocationsUpdate: failed to update location %s. %v", id, err)
		abort(c, http.StatusInternalServerError, "Failed to update location")
		return
	default:
		c.JSON(http.StatusOK, loc)
	}
}

// handleLocationsDelete godoc
// @Summary Delete location
// @Description Delete one specific location using provided ID.
// @Tags locations
// @Param id path string true "Location ID"
// @Success 204 "OK"
// @Failure 400 {object} HTTPError "Bad Request"
// @Failure 404 {object} HTTPError "Not Found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /locations/{id} [delete]
func (s *HTTPServer) handleLocationsDelete(c *gin.Context) {
	var query models.DeleteLocation
	if err := c.ShouldBindUri(&query); err != nil {
		logger.Errorf("LocationsDelete: invalid query. %v", err)
		abort(c, http.StatusBadRequest, "Invalid query")
		return
	}

	id, err := models.ParseID(query.ID)
	if err != nil {
		logger.Errorf("LocationsDelete: invalid query. %v", err)
		abort(c, http.StatusBadRequest, "Invalid query")
		return
	}

	err = s.api.LocationUsecase.DeleteLocation(c.Request.Context(), id)
	switch {
	case err == usecases.ErrLocationNotFound:
		abort(c, http.StatusNotFound, "Location not found")
		return
	case err != nil:
		logger.Errorf("LocationsDelete: failed to delete location %s. %v", id, err)
		abort(c, http.StatusInternalServerError, "Failed to delete location")
		return
	default:
		c.Status(http.StatusNoContent)
	}
}
