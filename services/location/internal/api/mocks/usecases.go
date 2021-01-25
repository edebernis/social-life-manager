package mocks

import (
	"context"

	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/stretchr/testify/mock"
)

// LocationUsecaseMock mocks location usecase
type LocationUsecaseMock struct {
	mock.Mock
}

// CreateLocation creates a new user location
func (u *LocationUsecaseMock) CreateLocation(ctx context.Context, loc *models.Location) error {
	args := u.Called(ctx, loc)
	return args.Error(0)
}

// FindLocationByName returns the unique location with that name
func (u *LocationUsecaseMock) FindLocationByName(ctx context.Context, name string) (*models.Location, error) {
	args := u.Called(ctx, name)
	loc := args.Get(0)
	if loc == nil {
		return nil, args.Error(1)
	}
	return loc.(*models.Location), args.Error(1)
}

// GetLocations returns all locations of a specific user
func (u *LocationUsecaseMock) GetLocations(ctx context.Context) (*models.Locations, error) {
	args := u.Called(ctx)
	locs := args.Get(0)
	if locs == nil {
		return nil, args.Error(1)
	}
	return locs.(*models.Locations), args.Error(1)
}
