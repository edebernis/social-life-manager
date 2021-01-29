package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/edebernis/social-life-manager/services/location/internal/models"
)

var (
	// ErrLocationAlreadyExists is raised when a matching location
	// is already stored in repository
	ErrLocationAlreadyExists = errors.New("location already exists")
	// ErrLocationNotFound is raised when specified location has not been found
	// in repository
	ErrLocationNotFound = errors.New("location not found")
	// ErrCategoryAlreadyExists is raised when a matching category
	// is already stored in repository
	ErrCategoryAlreadyExists = errors.New("category already exists")
	// ErrCategoryNotFound is raised when specified category has not been found
	// in repository
	ErrCategoryNotFound = errors.New("category not found")
)

// LocationRepository describes how to create, get, find, update and delete
// locations and categories in a repository
type LocationRepository interface {
	CreateCategory(context.Context, *models.Category) error
	GetCategories(context.Context) (*models.Categories, error)
	FindCategoryByID(context.Context, models.ID) (*models.Category, error)
	FindCategoryByName(context.Context, string) (*models.Category, error)
	UpdateCategory(context.Context, *models.Category) error
	DeleteCategory(context.Context, models.ID) error

	CreateLocation(context.Context, *models.Location) error
	GetLocations(context.Context) (*models.Locations, error)
	FindLocationByID(context.Context, models.ID) (*models.Location, error)
	FindLocationByName(context.Context, string) (*models.Location, error)
	FindLocationsByCategory(context.Context, *models.Category) (*models.Locations, error)
	UpdateLocation(context.Context, *models.Location) error
	DeleteLocation(context.Context, models.ID) error
}

// LocationUsecase represents a usecase around location handling
type LocationUsecase struct {
	repo LocationRepository
}

// NewLocationUsecase creates a new LocationUsecase object
func NewLocationUsecase(repo LocationRepository) *LocationUsecase {
	return &LocationUsecase{repo}
}

// CreateCategory stores a new location category in repository
func (u *LocationUsecase) CreateCategory(ctx context.Context, cat *models.Category) error {
	catWithSameName, err := u.repo.FindCategoryByName(ctx, cat.Name)
	if err != nil {
		return fmt.Errorf("CreateCategory: failed to find category by name in repository : %s. %w", cat.Name, err)
	}
	if catWithSameName != nil {
		return ErrCategoryAlreadyExists
	}

	if err := u.repo.CreateCategory(ctx, cat); err != nil {
		return fmt.Errorf("CreateCategory: failed to create category in repository : %v. %w", cat, err)
	}

	return nil
}

// GetCategories returns all categories of a specific user
func (u *LocationUsecase) GetCategories(ctx context.Context) (*models.Categories, error) {
	cats, err := u.repo.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetCategories: failed to get categories from repository. %w", err)
	}

	return cats, nil
}

// FindCategoryByID returns category matching specified ID or nil
func (u *LocationUsecase) FindCategoryByID(ctx context.Context, id models.ID) (*models.Category, error) {
	cat, err := u.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("FindCategoryByID: failed to get category by id, %s. %w", id, err)
	}
	if cat == nil {
		return nil, ErrCategoryNotFound
	}

	return cat, nil
}

// UpdateCategory updates specified location
func (u *LocationUsecase) UpdateCategory(ctx context.Context, cat *models.Category) error {
	catByID, err := u.repo.FindCategoryByID(ctx, cat.ID)
	if err != nil {
		return fmt.Errorf("UpdateCategory: failed to find category by id, %s. %w", cat.ID, err)
	}
	if catByID == nil {
		return ErrCategoryNotFound
	}

	if err := u.repo.UpdateCategory(ctx, cat); err != nil {
		return fmt.Errorf("UpdateCategory: failed to update category, %s. %w", cat.ID, err)
	}

	return nil
}

// DeleteCategory deletes specified category
func (u *LocationUsecase) DeleteCategory(ctx context.Context, id models.ID) error {
	cat, err := u.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("DeleteCategory: failed to find category by id, %s. %w", id, err)
	}
	if cat == nil {
		return ErrCategoryNotFound
	}

	if err := u.repo.DeleteCategory(ctx, id); err != nil {
		return fmt.Errorf("DeleteCategory: failed to delete category, %s. %w", id, err)
	}

	return nil
}

// CreateLocation stores a new user location in repository
func (u *LocationUsecase) CreateLocation(ctx context.Context, loc *models.Location) error {
	locWithSameName, err := u.repo.FindLocationByName(ctx, loc.Name)
	if err != nil {
		return fmt.Errorf("CreateLocation: failed to find location by name in repository : %s. %w", loc.Name, err)
	}
	if locWithSameName != nil {
		return ErrLocationAlreadyExists
	}

	cat, err := u.repo.FindCategoryByID(ctx, loc.Category)
	if err != nil {
		return fmt.Errorf("CreateLocation: failed to find category by ID, %s. %w", loc.Category, err)
	}
	if cat == nil {
		return ErrCategoryNotFound
	}

	if err := u.repo.CreateLocation(ctx, loc); err != nil {
		return fmt.Errorf("CreateLocation: failed to create location in repository : %v. %w", loc, err)
	}

	return nil
}

// GetLocations returns all locations of a specific user
func (u *LocationUsecase) GetLocations(ctx context.Context) (*models.Locations, error) {
	locations, err := u.repo.GetLocations(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetLocations: failed to get locations from repository. %w", err)
	}

	return locations, nil
}

// FindLocationByID returns location matching specified ID or nil
func (u *LocationUsecase) FindLocationByID(ctx context.Context, id models.ID) (*models.Location, error) {
	location, err := u.repo.FindLocationByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("FindLocationByID: failed to get location by id, %s. %w", id, err)
	}
	if location == nil {
		return nil, ErrLocationNotFound
	}

	return location, nil
}

// FindLocationsByCategory returns locations matching specified category or nil
func (u *LocationUsecase) FindLocationsByCategory(ctx context.Context, catID models.ID) (*models.Locations, error) {
	cat, err := u.repo.FindCategoryByID(ctx, catID)
	if err != nil {
		return nil, fmt.Errorf("FindLocationByID: failed to find category by ID, %s. %w", catID, err)
	}
	if cat == nil {
		return nil, ErrCategoryNotFound
	}

	locations, err := u.repo.FindLocationsByCategory(ctx, cat)
	if err != nil {
		return nil, fmt.Errorf("FindLocationsByCategory: failed to find locations by category, %s. %w", cat, err)
	}

	return locations, nil
}

// UpdateLocation updates specified location
func (u *LocationUsecase) UpdateLocation(ctx context.Context, loc *models.Location) error {
	locByID, err := u.repo.FindLocationByID(ctx, loc.ID)
	if err != nil {
		return fmt.Errorf("UpdateLocation: failed to find location by id, %s. %w", loc.ID, err)
	}
	if locByID == nil {
		return ErrLocationNotFound
	}

	cat, err := u.repo.FindCategoryByID(ctx, loc.Category)
	if err != nil {
		return fmt.Errorf("UpdateLocation: failed to find category by ID, %s. %w", loc.Category, err)
	}
	if cat == nil {
		return ErrCategoryNotFound
	}

	// Merge existing location with new fields to update
	if loc.Name == "" {
		loc.Name = locByID.Name
	}
	if loc.Address == "" {
		loc.Address = locByID.Address
	}
	if loc.Category == models.NilID {
		loc.Category = locByID.Category
	}

	if err := u.repo.UpdateLocation(ctx, loc); err != nil {
		return fmt.Errorf("UpdateLocation: failed to update location, %s. %w", loc.ID, err)
	}

	return nil
}

// DeleteLocation deletes specified location
func (u *LocationUsecase) DeleteLocation(ctx context.Context, id models.ID) error {
	loc, err := u.repo.FindLocationByID(ctx, id)
	if err != nil {
		return fmt.Errorf("DeleteLocation: failed to find location by id, %s. %w", id, err)
	}
	if loc == nil {
		return ErrLocationNotFound
	}

	if err := u.repo.DeleteLocation(ctx, id); err != nil {
		return fmt.Errorf("DeleteLocation: failed to delete location, %s. %w", id, err)
	}

	return nil
}
