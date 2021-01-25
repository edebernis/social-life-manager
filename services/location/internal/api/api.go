package api

import (
	"context"

	"github.com/edebernis/social-life-manager/services/location/internal/models"
)

// ILocationUsecase describes functions available in location usecase
type ILocationUsecase interface {
	CreateLocation(context.Context, *models.Location) error
	GetLocations(context.Context) (*models.Locations, error)
	FindLocationByName(context.Context, string) (*models.Location, error)
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
}
