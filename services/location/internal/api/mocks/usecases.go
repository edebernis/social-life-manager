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

// CreateCategory creates a new category
func (u *LocationUsecaseMock) CreateCategory(ctx context.Context, cat *models.Category) error {
	args := u.Called(ctx, cat)
	return args.Error(0)
}

// GetCategories returns all categories
func (u *LocationUsecaseMock) GetCategories(ctx context.Context) (*models.Categories, error) {
	args := u.Called(ctx)
	cats := args.Get(0)
	if cats == nil {
		return nil, args.Error(1)
	}
	return cats.(*models.Categories), args.Error(1)
}

// FindCategoryByID returns category matching specified ID or nil
func (u *LocationUsecaseMock) FindCategoryByID(ctx context.Context, id models.ID) (*models.Category, error) {
	args := u.Called(ctx, id)
	cat := args.Get(0)
	if cat == nil {
		return nil, args.Error(1)
	}
	return cat.(*models.Category), args.Error(1)
}

// UpdateCategory update specified category
func (u *LocationUsecaseMock) UpdateCategory(ctx context.Context, cat *models.Category) error {
	args := u.Called(ctx, cat)
	return args.Error(0)
}

// DeleteCategory deletes specified category
func (u *LocationUsecaseMock) DeleteCategory(ctx context.Context, id models.ID) error {
	args := u.Called(ctx, id)
	return args.Error(0)
}

// CreateLocation creates a new user location
func (u *LocationUsecaseMock) CreateLocation(ctx context.Context, loc *models.Location) error {
	args := u.Called(ctx, loc)
	return args.Error(0)
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

// FindLocationByID returns location matching specified ID or nil
func (u *LocationUsecaseMock) FindLocationByID(ctx context.Context, id models.ID) (*models.Location, error) {
	args := u.Called(ctx, id)
	loc := args.Get(0)
	if loc == nil {
		return nil, args.Error(1)
	}
	return loc.(*models.Location), args.Error(1)
}

// FindLocationsByCategory returns locations matching specified category or nil
func (u *LocationUsecaseMock) FindLocationsByCategory(ctx context.Context, id models.ID) (*models.Locations, error) {
	args := u.Called(ctx, id)
	locs := args.Get(0)
	if locs == nil {
		return nil, args.Error(1)
	}
	return locs.(*models.Locations), args.Error(1)
}

// UpdateLocation update specified location
func (u *LocationUsecaseMock) UpdateLocation(ctx context.Context, loc *models.Location) error {
	args := u.Called(ctx, loc)
	return args.Error(0)
}

// DeleteLocation deletes specified location
func (u *LocationUsecaseMock) DeleteLocation(ctx context.Context, id models.ID) error {
	args := u.Called(ctx, id)
	return args.Error(0)
}
