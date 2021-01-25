package sqlrepository

import (
	"context"
	"fmt"
	"time"

	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/edebernis/social-life-manager/services/location/internal/utils"
)

// CreateLocation creates a new user location in repository
func (r *SQLRepository) CreateLocation(ctx context.Context, loc *models.Location) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "INSERT INTO locations (id, name, address, category_id, user_id) VALUES ($1, $2, $3, $4, $5)"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("CreateLocation: failed to prepare context for query %s. %w", query, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, loc.ID.String(), loc.Name, loc.Address, loc.Category.String(), loc.User.String())
	if err != nil {
		return fmt.Errorf("CreateLocation: failed to exec context for query %s. %w", query, err)
	}

	return nil
}

// GetLocations returns all user locations stored in the repository
func (r *SQLRepository) GetLocations(ctx context.Context) (*models.Locations, error) {
	return nil, utils.NotImplementedError("GetLocations")
}

// FindLocationByID returns the user location that matches the requested ID or nil
func (r *SQLRepository) FindLocationByID(ctx context.Context, id models.ID) (*models.Location, error) {
	return nil, utils.NotImplementedError("FindLocationByID")
}

// FindLocationByName returns the user location that matches the requested name or nil
func (r *SQLRepository) FindLocationByName(ctx context.Context, name string) (*models.Location, error) {
	return nil, utils.NotImplementedError("FindLocationByName")
}

// FindLocationsByCategory returns all user locations filtered by specified category
func (r *SQLRepository) FindLocationsByCategory(ctx context.Context, cat *models.Category) (*models.Locations, error) {
	return nil, utils.NotImplementedError("FindLocationsByCategory")
}

// UpdateLocation updates specified location in repository
func (r *SQLRepository) UpdateLocation(ctx context.Context, loc *models.Location) error {
	return utils.NotImplementedError("UpdateLocation")
}

// DeleteLocation deletes location in repository
func (r *SQLRepository) DeleteLocation(ctx context.Context, id models.ID) error {
	return utils.NotImplementedError("DeleteLocation")
}
