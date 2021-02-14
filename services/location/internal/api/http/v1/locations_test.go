package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/edebernis/social-life-manager/services/location/internal/api/mocks"
	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/edebernis/social-life-manager/services/location/internal/usecases"
	"github.com/edebernis/social-life-manager/services/location/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestV1CreateCategoryWithInvalidQuery(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/categories", nil, nil)

	server.handleCategoriesCreate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1CreateCategoryWithAlreadyExistsError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/categories", &gin.H{
		"name": "Test Category",
	}, nil)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("CreateCategory", utils.MockContextMatcher, mock.AnythingOfType("*models.Category")).
		Return(usecases.ErrCategoryAlreadyExists)

	server.handleCategoriesCreate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1CreateCategoryWithError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/categories", &gin.H{
		"name": "Test Category",
	}, nil)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("CreateCategory", utils.MockContextMatcher, mock.AnythingOfType("*models.Category")).
		Return(errors.New("failed"))

	server.handleCategoriesCreate(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
}

func TestV1CreateCategoryWithSuccess(t *testing.T) {
	ctx, resp, server := newHandlerTestContext(t, "POST", "/api/v1/categories", &gin.H{
		"name": "Test Category",
	}, nil)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("CreateCategory", utils.MockContextMatcher, mock.AnythingOfType("*models.Category")).
		Return(nil)

	server.handleCategoriesCreate(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	var cat models.Category
	err := json.NewDecoder(resp.Result().Body).Decode(&cat)
	if assert.NoError(t, err) {
		assert.Equal(t, "Test Category", cat.Name)
	}
}

func TestV1GetCategoriesWithError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "GET", "/api/v1/categories", nil, nil)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("GetCategories", utils.MockContextMatcher).
		Return(nil, errors.New("failed"))

	server.handleCategoriesGet(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
}

func TestV1GetCategoriesWithSuccess(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "GET", "/api/v1/categories", nil, nil)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("GetCategories", utils.MockContextMatcher).
		Return(&models.Categories{}, nil)

	server.handleCategoriesGet(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
}

func TestV1GetCategoriesByIDWithError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"GET",
		"/api/v1/categories/4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("FindCategoryByID", utils.MockContextMatcher, id).
		Return(nil, errors.New("failed"))

	server.handleCategoriesGetByID(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
}

func TestV1GetCategoriesByIDWithCategoryNotFound(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"GET",
		"/api/v1/categories/4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("FindCategoryByID", utils.MockContextMatcher, id).
		Return(nil, usecases.ErrCategoryNotFound)

	server.handleCategoriesGetByID(ctx)

	assert.Equal(t, http.StatusNotFound, ctx.Writer.Status())
}

func TestV1GetCategoriesByIDWithInvalidID(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"GET",
		"/api/v1/categories/123",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "123",
			},
		},
	)

	server.handleCategoriesGetByID(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1GetCategoriesByIDWithInvalidQuery(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"GET",
		"/api/v1/categories/",
		nil,
		nil,
	)

	server.handleCategoriesGetByID(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1GetCategoriesByIDWithSuccess(t *testing.T) {
	ctx, resp, server := newHandlerTestContext(
		t,
		"GET",
		"/api/v1/categories/4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")
	cat := models.NewCategory(id, "Test Category")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("FindCategoryByID", utils.MockContextMatcher, id).
		Return(cat, nil)

	server.handleCategoriesGetByID(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	var returnedCat models.Category
	err := json.NewDecoder(resp.Result().Body).Decode(&returnedCat)
	if assert.NoError(t, err) {
		assert.Equal(t, *cat, returnedCat)
	}
}

func TestV1UpdateCategoryWithError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"PUT",
		"/api/v1/categories/4b7a536e-7109-4a39-9549-f06f74f2093e",
		&gin.H{
			"name": "Test Category",
		},
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")
	cat := models.NewCategory(id, "Test Category")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("UpdateCategory", utils.MockContextMatcher, cat).
		Return(errors.New("failed"))

	server.handleCategoriesUpdate(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
}

func TestV1UpdateCategoryWithCategoryNotFound(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"PUT",
		"/api/v1/categories/4b7a536e-7109-4a39-9549-f06f74f2093e",
		&gin.H{
			"name": "Test Category",
		},
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")
	cat := models.NewCategory(id, "Test Category")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("UpdateCategory", utils.MockContextMatcher, cat).
		Return(usecases.ErrCategoryNotFound)

	server.handleCategoriesUpdate(ctx)

	assert.Equal(t, http.StatusNotFound, ctx.Writer.Status())
}

func TestV1UpdateCategoryWithInvalidID(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"PUT",
		"/api/v1/categories/123",
		&gin.H{
			"name": "Test Category",
		},
		&[]gin.Param{
			{
				Key:   "id",
				Value: "123",
			},
		},
	)

	server.handleCategoriesUpdate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1UpdateCategoryWithInvalidQuery(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"PUT",
		"/api/v1/categories/",
		&gin.H{
			"name": "Test Category",
		},
		nil,
	)

	server.handleCategoriesUpdate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1UpdateCategoryWithInvalidBody(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"PUT",
		"/api/v1/categories/4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	server.handleCategoriesUpdate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1UpdateCategoryWithSuccess(t *testing.T) {
	ctx, resp, server := newHandlerTestContext(
		t,
		"PUT",
		"/api/v1/categories/4b7a536e-7109-4a39-9549-f06f74f2093e",
		&gin.H{
			"name": "Test Category",
		},
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")
	cat := models.NewCategory(id, "Test Category")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("UpdateCategory", utils.MockContextMatcher, cat).
		Return(nil)

	server.handleCategoriesUpdate(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	var returnedCat models.Category
	err := json.NewDecoder(resp.Result().Body).Decode(&returnedCat)
	if assert.NoError(t, err) {
		assert.Equal(t, *cat, returnedCat)
	}
}

func TestV1DeleteCategoryWithError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"DELETE",
		"/api/v1/categories/4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("DeleteCategory", utils.MockContextMatcher, id).
		Return(errors.New(("failed")))

	server.handleCategoriesDelete(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
}

func TestV1DeleteCategoryWithCategoryNotFound(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"DELETE",
		"/api/v1/categories/4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("DeleteCategory", utils.MockContextMatcher, id).
		Return(usecases.ErrCategoryNotFound)

	server.handleCategoriesDelete(ctx)

	assert.Equal(t, http.StatusNotFound, ctx.Writer.Status())
}

func TestV1DeleteCategoryWithInvalidID(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"DELETE",
		"/api/v1/categories/123",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "123",
			},
		},
	)

	server.handleCategoriesDelete(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1DeleteCategoryWithInvalidQuery(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"DELETE",
		"/api/v1/categories/",
		nil,
		nil,
	)

	server.handleCategoriesDelete(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1DeleteCategoryWithSuccess(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"DELETE",
		"/api/v1/categories/4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("DeleteCategory", utils.MockContextMatcher, id).
		Return(nil)

	server.handleCategoriesDelete(ctx)

	assert.Equal(t, http.StatusNoContent, ctx.Writer.Status())
}

func TestV1CreateLocationWithoutName(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"address":     "1 rue de la Poste, 75001 Paris",
		"category_id": "4b7a536e-7109-4a39-9549-f06f74f2093e",
	}, nil)

	server.handleLocationsCreate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1CreateLocationWithoutAddress(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"name":        "Test Location",
		"category_id": "4b7a536e-7109-4a39-9549-f06f74f2093e",
	}, nil)

	server.handleLocationsCreate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1CreateLocationWithoutCategory(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"name":    "Test Location",
		"address": "1 rue de la Poste, 75001 Paris",
	}, nil)

	server.handleLocationsCreate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1CreateLocationWithInvalidCategory(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"name":        "Test Location",
		"address":     "1 rue de la Poste, 75001 Paris",
		"category_id": "12",
	}, nil)

	server.handleLocationsCreate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1CreateLocationWithAlreadyExistsError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"name":        "Test Location",
		"address":     "1 rue de la Poste, 75001 Paris",
		"category_id": "4b7a536e-7109-4a39-9549-f06f74f2093e",
	}, nil)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("CreateLocation", utils.MockContextMatcher, mock.AnythingOfType("*models.Location")).
		Return(usecases.ErrLocationAlreadyExists)

	server.handleLocationsCreate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1CreateLocationWithCategoryNotFoundError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"name":        "Test Location",
		"address":     "1 rue de la Poste, 75001 Paris",
		"category_id": "4b7a536e-7109-4a39-9549-f06f74f2093e",
	}, nil)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("CreateLocation", utils.MockContextMatcher, mock.AnythingOfType("*models.Location")).
		Return(usecases.ErrCategoryNotFound)

	server.handleLocationsCreate(ctx)

	assert.Equal(t, http.StatusNotFound, ctx.Writer.Status())
}

func TestV1CreateLocationWithError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"name":        "Test Location",
		"address":     "1 rue de la Poste, 75001 Paris",
		"category_id": "4b7a536e-7109-4a39-9549-f06f74f2093e",
	}, nil)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("CreateLocation", utils.MockContextMatcher, mock.AnythingOfType("*models.Location")).
		Return(errors.New("failed"))

	server.handleLocationsCreate(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
}

func TestV1CreateLocationWithSuccess(t *testing.T) {
	ctx, resp, server := newHandlerTestContext(t, "POST", "/api/v1/locations", &gin.H{
		"name":        "Test Location",
		"address":     "1 rue de la Poste, 75001 Paris",
		"category_id": "4b7a536e-7109-4a39-9549-f06f74f2093e",
	}, nil)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("CreateLocation", utils.MockContextMatcher, mock.AnythingOfType("*models.Location")).
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

func TestV1GetLocationsWithError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "GET", "/api/v1/locations", nil, nil)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("GetLocations", utils.MockContextMatcher).
		Return(nil, errors.New("failed"))

	server.handleLocationsGet(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
}

func TestV1GetLocationsByCategoryWithCategoryNotFound(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"GET",
		"/api/v1/locations?category_id=4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		nil,
	)

	catID, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("FindLocationsByCategory", utils.MockContextMatcher, catID).
		Return(nil, usecases.ErrCategoryNotFound)

	server.handleLocationsGet(ctx)

	assert.Equal(t, http.StatusNotFound, ctx.Writer.Status())
}

func TestV1GetLocationsWithInvalidCategory(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "GET", "/api/v1/locations?category_id=123", nil, nil)

	server.handleLocationsGet(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1GetLocationsWithSuccess(t *testing.T) {
	ctx, _, server := newHandlerTestContext(t, "GET", "/api/v1/locations", nil, nil)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("GetLocations", utils.MockContextMatcher).
		Return(&models.Locations{}, nil)

	server.handleLocationsGet(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
}

func TestV1GetLocationsByCategoryWithSuccess(t *testing.T) {
	ctx, resp, server := newHandlerTestContext(
		t,
		"GET",
		"/api/v1/locations?category_id=4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		nil,
	)

	catID, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")
	user, _ := models.NewUserFromContext(ctx.Request.Context())
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", catID, user.ID)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("FindLocationsByCategory", utils.MockContextMatcher, catID).
		Return(&models.Locations{loc}, nil)

	server.handleLocationsGet(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	var returnedLocs models.Locations
	err := json.NewDecoder(resp.Result().Body).Decode(&returnedLocs)
	if assert.NoError(t, err) {
		assert.Equal(t, loc, returnedLocs[0])
	}
}

func TestV1GetLocationsByIDWithError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"GET",
		"/api/v1/locations/4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("FindLocationByID", utils.MockContextMatcher, id).
		Return(nil, errors.New("failed"))

	server.handleLocationsGetByID(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
}

func TestV1GetLocationsByIDWithLocationNotFound(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"GET",
		"/api/v1/locations/4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("FindLocationByID", utils.MockContextMatcher, id).
		Return(nil, usecases.ErrLocationNotFound)

	server.handleLocationsGetByID(ctx)

	assert.Equal(t, http.StatusNotFound, ctx.Writer.Status())
}

func TestV1GetLocationsByIDWithInvalidID(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"GET",
		"/api/v1/locations/123",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "123",
			},
		},
	)

	server.handleLocationsGetByID(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1GetLocationsByIDWithInvalidQuery(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"GET",
		"/api/v1/locations/",
		nil,
		nil,
	)

	server.handleLocationsGetByID(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1GetLocationsByIDWithSuccess(t *testing.T) {
	ctx, resp, server := newHandlerTestContext(
		t,
		"GET",
		"/api/v1/locations/4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")
	user, _ := models.NewUserFromContext(ctx.Request.Context())
	loc := models.NewLocation(id, "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), user.ID)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("FindLocationByID", utils.MockContextMatcher, id).
		Return(loc, nil)

	server.handleLocationsGetByID(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	var returnedLoc models.Location
	err := json.NewDecoder(resp.Result().Body).Decode(&returnedLoc)
	if assert.NoError(t, err) {
		assert.Equal(t, *loc, returnedLoc)
	}
}

func TestV1UpdateLocationWithError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"PUT",
		"/api/v1/locations/4b7a536e-7109-4a39-9549-f06f74f2093e",
		&gin.H{
			"name":        "Test Location",
			"address":     "1 rue de la Poste, 75001 Paris",
			"category_id": "550e8400-e29b-41d4-a716-446655440000",
		},
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")
	catID, _ := models.ParseID("550e8400-e29b-41d4-a716-446655440000")
	user, _ := models.NewUserFromContext(ctx.Request.Context())
	loc := models.NewLocation(id, "Test Location", "1 rue de la Poste, 75001 Paris", catID, user.ID)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("UpdateLocation", utils.MockContextMatcher, loc).
		Return(errors.New("failed"))

	server.handleLocationsUpdate(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
}

func TestV1UpdateLocationWithCategoryNotFound(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"PUT",
		"/api/v1/locations/4b7a536e-7109-4a39-9549-f06f74f2093e",
		&gin.H{
			"name":        "Test Location",
			"address":     "1 rue de la Poste, 75001 Paris",
			"category_id": "550e8400-e29b-41d4-a716-446655440000",
		},
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")
	catID, _ := models.ParseID("550e8400-e29b-41d4-a716-446655440000")
	user, _ := models.NewUserFromContext(ctx.Request.Context())
	loc := models.NewLocation(id, "Test Location", "1 rue de la Poste, 75001 Paris", catID, user.ID)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("UpdateLocation", utils.MockContextMatcher, loc).
		Return(usecases.ErrCategoryNotFound)

	server.handleLocationsUpdate(ctx)

	assert.Equal(t, http.StatusNotFound, ctx.Writer.Status())
}

func TestV1UpdateLocationWithLocationNotFound(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"PUT",
		"/api/v1/locations/4b7a536e-7109-4a39-9549-f06f74f2093e",
		&gin.H{
			"name":        "Test Location",
			"address":     "1 rue de la Poste, 75001 Paris",
			"category_id": "550e8400-e29b-41d4-a716-446655440000",
		},
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")
	catID, _ := models.ParseID("550e8400-e29b-41d4-a716-446655440000")
	user, _ := models.NewUserFromContext(ctx.Request.Context())
	loc := models.NewLocation(id, "Test Location", "1 rue de la Poste, 75001 Paris", catID, user.ID)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("UpdateLocation", utils.MockContextMatcher, loc).
		Return(usecases.ErrLocationNotFound)

	server.handleLocationsUpdate(ctx)

	assert.Equal(t, http.StatusNotFound, ctx.Writer.Status())
}

func TestV1UpdateLocationWithInvalidID(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"PUT",
		"/api/v1/locations/123",
		&gin.H{
			"name":        "Test Location",
			"address":     "1 rue de la Poste, 75001 Paris",
			"category_id": "550e8400-e29b-41d4-a716-446655440000",
		},
		&[]gin.Param{
			{
				Key:   "id",
				Value: "123",
			},
		},
	)

	server.handleLocationsUpdate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1UpdateLocationWithInvalidCategory(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"PUT",
		"/api/v1/locations/4b7a536e-7109-4a39-9549-f06f74f2093e",
		&gin.H{
			"name":        "Test Location",
			"address":     "1 rue de la Poste, 75001 Paris",
			"category_id": "123",
		},
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	server.handleLocationsUpdate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1UpdateLocationWithInvalidQuery(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"PUT",
		"/api/v1/locations/",
		&gin.H{
			"name":        "Test Location",
			"address":     "1 rue de la Poste, 75001 Paris",
			"category_id": "550e8400-e29b-41d4-a716-446655440000",
		},
		nil,
	)

	server.handleLocationsUpdate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1UpdateLocationWithInvalidBody(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"PUT",
		"/api/v1/locations/4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	server.handleLocationsUpdate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1UpdateLocationWithSuccess(t *testing.T) {
	ctx, resp, server := newHandlerTestContext(
		t,
		"PUT",
		"/api/v1/locations/4b7a536e-7109-4a39-9549-f06f74f2093e",
		&gin.H{
			"name":        "Test Location",
			"address":     "1 rue de la Poste, 75001 Paris",
			"category_id": "550e8400-e29b-41d4-a716-446655440000",
		},
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")
	catID, _ := models.ParseID("550e8400-e29b-41d4-a716-446655440000")
	user, _ := models.NewUserFromContext(ctx.Request.Context())
	loc := models.NewLocation(id, "Test Location", "1 rue de la Poste, 75001 Paris", catID, user.ID)

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("UpdateLocation", utils.MockContextMatcher, loc).
		Return(nil)

	server.handleLocationsUpdate(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	var returnedLoc models.Location
	err := json.NewDecoder(resp.Result().Body).Decode(&returnedLoc)
	if assert.NoError(t, err) {
		assert.Equal(t, *loc, returnedLoc)
	}
}

func TestV1DeleteLocationWithError(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"DELETE",
		"/api/v1/locations/4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("DeleteLocation", utils.MockContextMatcher, id).
		Return(errors.New(("failed")))

	server.handleLocationsDelete(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
}

func TestV1DeleteLocationWithLocationNotFound(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"DELETE",
		"/api/v1/locations/4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("DeleteLocation", utils.MockContextMatcher, id).
		Return(usecases.ErrLocationNotFound)

	server.handleLocationsDelete(ctx)

	assert.Equal(t, http.StatusNotFound, ctx.Writer.Status())
}

func TestV1DeleteLocationWithInvalidID(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"DELETE",
		"/api/v1/locations/123",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "123",
			},
		},
	)

	server.handleLocationsDelete(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1DeleteLocationWithInvalidQuery(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"DELETE",
		"/api/v1/locations/",
		nil,
		nil,
	)

	server.handleLocationsDelete(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestV1DeleteLocationWithSuccess(t *testing.T) {
	ctx, _, server := newHandlerTestContext(
		t,
		"DELETE",
		"/api/v1/locations/4b7a536e-7109-4a39-9549-f06f74f2093e",
		nil,
		&[]gin.Param{
			{
				Key:   "id",
				Value: "4b7a536e-7109-4a39-9549-f06f74f2093e",
			},
		},
	)

	id, _ := models.ParseID("4b7a536e-7109-4a39-9549-f06f74f2093e")

	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("DeleteLocation", utils.MockContextMatcher, id).
		Return(nil)

	server.handleLocationsDelete(ctx)

	assert.Equal(t, http.StatusNoContent, ctx.Writer.Status())
}
