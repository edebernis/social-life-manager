package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/edebernis/social-life-manager/services/location/internal/utils"
)

var (
	// ErrLocationAlreadyExists specifies that a matching location
	// is already stored in repository
	ErrLocationAlreadyExists = errors.New("location already exists")
)

// LocationRepository describes how to create, get, find, update and delete
// locations in a repository
type LocationRepository interface {
	Open() error
	CreateLocation(context.Context, *models.Location) error
	GetLocations(context.Context) (*models.Locations, error)
	FindLocationByID(context.Context, models.ID) (*models.Location, error)
	FindLocationByName(context.Context, string) (*models.Location, error)
	FindLocationsByCategory(context.Context, *models.Category) (*models.Locations, error)
	UpdateLocation(context.Context, *models.Location) error
	DeleteLocation(context.Context, models.ID) error
	Close() error
}

// LocationUsecase represents a usecase around location handling
type LocationUsecase struct {
	repo LocationRepository
}

// NewLocationUsecase creates a new LocationUsecase object
func NewLocationUsecase(repo LocationRepository) *LocationUsecase {
	return &LocationUsecase{repo}
}

// CreateLocation stores a new user location in repository
func (u *LocationUsecase) CreateLocation(ctx context.Context, loc *models.Location) error {
	locWithSameName, err := u.repo.FindLocationByName(ctx, loc.Name)
	if err != nil {
		return fmt.Errorf("CreateLocation: failed to find location by name in repository : %s. %w", loc.Name, err)
	}

	// Location name must be unique
	if locWithSameName != nil {
		return ErrLocationAlreadyExists
	}

	if err := u.repo.CreateLocation(ctx, loc); err != nil {
		return fmt.Errorf("CreateLocation: failed to create location in repository : %v. %w", loc, err)
	}

	return nil
}

// FindLocationByName searches for location with the same name
func (*LocationUsecase) FindLocationByName(ctx context.Context, name string) (*models.Location, error) {
	return nil, utils.NotImplementedError("FindLocationByName")
}

// GetLocations returns all locations of a specific user
func (*LocationUsecase) GetLocations(ctx context.Context) (*models.Locations, error) {
	return nil, utils.NotImplementedError("GetLocations")
}
