package api

import (
	"context"

	"github.com/edebernis/social-life-manager/services/location/internal/models"
)

// ILocationUsecase describes functions available in location usecase
type ILocationUsecase interface {
	CreateCategory(context.Context, *models.Category) error
	GetCategories(context.Context) (*models.Categories, error)
	FindCategoryByID(context.Context, models.ID) (*models.Category, error)
	UpdateCategory(context.Context, *models.Category) error
	DeleteCategory(context.Context, models.ID) error

	CreateLocation(context.Context, *models.Location) error
	GetLocations(context.Context) (*models.Locations, error)
	FindLocationByID(context.Context, models.ID) (*models.Location, error)
	FindLocationsByCategory(context.Context, models.ID) (*models.Locations, error)
	UpdateLocation(context.Context, *models.Location) error
	DeleteLocation(context.Context, models.ID) error
}

// API lists usecases of this service
type API struct {
	LocationUsecase ILocationUsecase
}

// NewAPI builds a new API
func NewAPI(locationUsecase ILocationUsecase) *API {
	return &API{locationUsecase}
}

// Server describes a service to listen and serve requests
type Server interface {
	Serve(addr string) error
	Shutdown() error
}
