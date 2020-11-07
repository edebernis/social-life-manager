package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/edebernis/social-life-manager/location/api/mocks"
	"github.com/edebernis/social-life-manager/location/models"
	"github.com/edebernis/social-life-manager/location/usecases"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestV1CreateLocationWithoutName(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"address":     "1 rue de la Poste, 75001 Paris",
		"category_id": "4b7a536e-7109-4a39-9549-f06f74f2093e",
	})

	server.handleLocationsCreate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1CreateLocationWithoutAddress(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"name":        "Test Location",
		"category_id": "4b7a536e-7109-4a39-9549-f06f74f2093e",
	})

	server.handleLocationsCreate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1CreateLocationWithoutCategory(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"name":    "Test Location",
		"address": "1 rue de la Poste, 75001 Paris",
	})

	server.handleLocationsCreate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1CreateLocationWithInvalidCategory(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"name":        "Test Location",
		"address":     "1 rue de la Poste, 75001 Paris",
		"category_id": "12",
	})

	server.handleLocationsCreate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1CreateLocationWithAlreadyExistsError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"name":        "Test Location",
		"address":     "1 rue de la Poste, 75001 Paris",
		"category_id": "4b7a536e-7109-4a39-9549-f06f74f2093e",
	})

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("CreateLocation", mockContextMatcher, mock.AnythingOfType("*models.Location")).
		Return(usecases.ErrLocationAlreadyExists)

	server.handleLocationsCreate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1CreateLocationWithError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"name":        "Test Location",
		"address":     "1 rue de la Poste, 75001 Paris",
		"category_id": "4b7a536e-7109-4a39-9549-f06f74f2093e",
	})

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("CreateLocation", mockContextMatcher, mock.AnythingOfType("*models.Location")).
		Return(errors.New("failed"))

	server.handleLocationsCreate(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
}

func TestV1CreateLocationWithSuccess(t *testing.T) {
	ctx, resp, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"name":        "Test Location",
		"address":     "1 rue de la Poste, 75001 Paris",
		"category_id": "4b7a536e-7109-4a39-9549-f06f74f2093e",
	})

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("CreateLocation", mockContextMatcher, mock.AnythingOfType("*models.Location")).
		Return(nil)

	server.handleLocationsCreate(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	var loc models.Location
	err := json.NewDecoder(resp.Result().Body).Decode(&loc)
	if assert.NoError(t, err) {
		assert.IsType(t, models.ID{}, loc.ID)
		assert.Equal(t, "Test Location", loc.Name)
		assert.Equal(t, "1 rue de la Poste, 75001 Paris", loc.Address)

		catID, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")
		assert.Equal(t, catID, loc.Category)
	}
}

func TestV1ListLocationsWithError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "GET", "/api/v1/locations", nil)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("GetLocations", mockContextMatcher).
		Return(nil, errors.New("failed"))

	server.handleLocationsGet(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
}

func TestV1ListLocationsWithSuccess(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "GET", "/api/v1/locations", nil)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("GetLocations", mockContextMatcher).
		Return(&models.Locations{}, nil)

	server.handleLocationsGet(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
}
