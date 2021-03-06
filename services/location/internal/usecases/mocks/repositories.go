package mocks

import (
	"context"

	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/stretchr/testify/mock"
)

// LocationRepositoryMock mocks location repository
type LocationRepositoryMock struct {
	mock.Mock
}

// CreateCategory creates a new category in repository
func (r *LocationRepositoryMock) CreateCategory(ctx context.Context, cat *models.Category) error {
	args := r.Called(ctx, cat)
	return args.Error(0)
}

// GetCategories fetches all categories in repository
func (r *LocationRepositoryMock) GetCategories(ctx context.Context) (*models.Categories, error) {
	args := r.Called(ctx)

	cats := args.Get(0)
	if cats == nil {
		return nil, args.Error(1)
	}
	return cats.(*models.Categories), args.Error(1)
}

// FindCategoryByID returns category matching specified ID or nil
func (r *LocationRepositoryMock) FindCategoryByID(ctx context.Context, id models.ID) (*models.Category, error) {
	args := r.Called(ctx, id)
	cat := args.Get(0)
	if cat == nil {
		return nil, args.Error(1)
	}
	return cat.(*models.Category), args.Error(1)
}

// FindCategoryByName returns category matching specified name or nil
func (r *LocationRepositoryMock) FindCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	args := r.Called(ctx, name)
	cat := args.Get(0)
	if cat == nil {
		return nil, args.Error(1)
	}
	return cat.(*models.Category), args.Error(1)
}

// UpdateCategory updates category in repository
func (r *LocationRepositoryMock) UpdateCategory(ctx context.Context, cat *models.Category) error {
	args := r.Called(ctx, cat)
	return args.Error(0)
}

// DeleteCategory deletes category in repository
func (r *LocationRepositoryMock) DeleteCategory(ctx context.Context, id models.ID) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

// CreateLocation creates a new user location in repository
func (r *LocationRepositoryMock) CreateLocation(ctx context.Context, loc *models.Location) error {
	args := r.Called(ctx, loc)
	return args.Error(0)
}

// GetLocations returns all user locations stored in the repository
func (r *LocationRepositoryMock) GetLocations(ctx context.Context) (*models.Locations, error) {
	args := r.Called(ctx)

	locs := args.Get(0)
	if locs == nil {
		return nil, args.Error(1)
	}
	return locs.(*models.Locations), args.Error(1)
}

// FindLocationByID returns the user location that matches the requested ID or nil
func (r *LocationRepositoryMock) FindLocationByID(ctx context.Context, id models.ID) (*models.Location, error) {
	args := r.Called(ctx, id)
	loc := args.Get(0)
	if loc == nil {
		return nil, args.Error(1)
	}
	return loc.(*models.Location), args.Error(1)
}

// FindLocationByName returns the user location that matches the requested name or nil
func (r *LocationRepositoryMock) FindLocationByName(ctx context.Context, name string) (*models.Location, error) {
	args := r.Called(ctx, name)
	loc := args.Get(0)
	if loc == nil {
		return nil, args.Error(1)
	}
	return loc.(*models.Location), args.Error(1)
}

// FindLocationsByCategory returns all user locations filtered by specified category
func (r *LocationRepositoryMock) FindLocationsByCategory(ctx context.Context, cat *models.Category) (*models.Locations, error) {
	args := r.Called(ctx, cat)
	locs := args.Get(0)
	if locs == nil {
		return nil, args.Error(1)
	}
	return locs.(*models.Locations), args.Error(1)
}

// UpdateLocation updates specified location in repository
func (r *LocationRepositoryMock) UpdateLocation(ctx context.Context, loc *models.Location) error {
	args := r.Called(ctx, loc)
	return args.Error(0)
}

// DeleteLocation deletes location in repository
func (r *LocationRepositoryMock) DeleteLocation(ctx context.Context, id models.ID) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}
