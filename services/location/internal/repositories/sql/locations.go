package sqlrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/edebernis/social-life-manager/services/location/internal/models"
)

// CreateCategory creates a new category in repository
func (r *SQLRepository) CreateCategory(ctx context.Context, cat *models.Category) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "INSERT INTO categories (id, name) VALUES ($1, $2)"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("CreateCategory: failed to prepare context for query %s. %w", query, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, cat.ID, cat.Name)
	if err != nil {
		return fmt.Errorf("CreateCategory: failed to exec context for query %s. %w", query, err)
	}

	return nil
}

// GetCategories fetches all categories in repository
func (r *SQLRepository) GetCategories(ctx context.Context) (*models.Categories, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "SELECT id, name FROM categories"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetCategories: failed to prepare context for query %s. %w", query, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetCategories: failed to query context for query %s. %w", query, err)
	}
	defer rows.Close()

	cats := make(models.Categories, 0)
	for rows.Next() {
		cat := new(models.Category)
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, fmt.Errorf("GetCategories: failed to scan SQL row. %w", err)
		}
		cats = append(cats, cat)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetCategories: rows failed. %w", err)
	}

	return &cats, nil
}

// FindCategoryByID returns category matching specified ID or nil
func (r *SQLRepository) FindCategoryByID(ctx context.Context, id models.ID) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "SELECT id, name FROM categories WHERE id = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("FindCategoryByID: failed to prepare context for query %s. %w", query, err)
	}
	defer stmt.Close()

	var cat models.Category
	err = stmt.QueryRowContext(ctx, id).Scan(&cat.ID, &cat.Name)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("FindCategoryByID: failed to query row for query %s. %w", query, err)
	default:
		return &cat, nil
	}
}

// FindCategoryByName returns category matching specified name or nil
func (r *SQLRepository) FindCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "SELECT id, name FROM categories WHERE name = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("FindCategoryByName: failed to prepare context for query %s. %w", query, err)
	}
	defer stmt.Close()

	var cat models.Category
	err = stmt.QueryRowContext(ctx, name).Scan(&cat.ID, &cat.Name)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("FindCategoryByName: failed to query row for query %s. %w", query, err)
	default:
		return &cat, nil
	}
}

// UpdateCategory updates category in repository
func (r *SQLRepository) UpdateCategory(ctx context.Context, cat *models.Category) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "UPDATE categories SET name = $1 WHERE id = $2"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("UpdateCategory: failed to prepare context for query %s. %w", query, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, cat.Name, cat.ID)
	if err != nil {
		return fmt.Errorf("UpdateCategory: failed to exec context for query %s. %w", query, err)
	}

	return nil
}

// DeleteCategory deletes category in repository
func (r *SQLRepository) DeleteCategory(ctx context.Context, id models.ID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "DELETE FROM categories WHERE id = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("DeleteCategory: failed to prepare context for query %s. %w", query, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("DeleteCategory: failed to exec context for query %s. %w", query, err)
	}

	return nil
}

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

	_, err = stmt.ExecContext(ctx, loc.ID, loc.Name, loc.Address, loc.Category, loc.User)
	if err != nil {
		return fmt.Errorf("CreateLocation: failed to exec context for query %s. %w", query, err)
	}

	return nil
}

// GetLocations returns all user locations stored in the repository
func (r *SQLRepository) GetLocations(ctx context.Context) (*models.Locations, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE user_id = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetLocations: failed to prepare context for query %s. %w", query, err)
	}
	defer stmt.Close()

	user, ok := models.NewUserFromContext(ctx)
	if !ok {
		return nil, errors.New("GetLocations: Failed to get user from context")
	}

	rows, err := stmt.QueryContext(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("GetLocations: failed to query context for query %s. %w", query, err)
	}
	defer rows.Close()

	locs := make(models.Locations, 0)
	for rows.Next() {
		loc := new(models.Location)
		if err := rows.Scan(&loc.ID, &loc.Name, &loc.Address, &loc.Category, &loc.User); err != nil {
			return nil, fmt.Errorf("GetLocations: failed to scan SQL row. %w", err)
		}
		locs = append(locs, loc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetLocations: rows failed. %w", err)
	}

	return &locs, nil
}

// FindLocationByID returns the user location that matches the requested ID or nil
func (r *SQLRepository) FindLocationByID(ctx context.Context, id models.ID) (*models.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE id = $1 AND user_id = $2"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("FindLocationByID: failed to prepare context for query %s. %w", query, err)
	}
	defer stmt.Close()

	user, ok := models.NewUserFromContext(ctx)
	if !ok {
		return nil, errors.New("FindLocationByID: Failed to get user from context")
	}

	var loc models.Location
	err = stmt.QueryRowContext(ctx, id, user.ID).Scan(&loc.ID, &loc.Name, &loc.Address, &loc.Category, &loc.User)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("FindLocationByID: failed to query row for query %s. %w", query, err)
	default:
		return &loc, nil
	}
}

// FindLocationByName returns the user location that matches the requested name or nil
func (r *SQLRepository) FindLocationByName(ctx context.Context, name string) (*models.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE name = $1 AND user_id = $2"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("FindLocationByName: failed to prepare context for query %s. %w", query, err)
	}
	defer stmt.Close()

	user, ok := models.NewUserFromContext(ctx)
	if !ok {
		return nil, errors.New("FindLocationByName: Failed to get user from context")
	}

	var loc models.Location
	err = stmt.QueryRowContext(ctx, name, user.ID).Scan(&loc.ID, &loc.Name, &loc.Address, &loc.Category, &loc.User)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("FindLocationByName: failed to query row for query %s. %w", query, err)
	default:
		return &loc, nil
	}
}

// FindLocationsByCategory returns all user locations filtered by specified category
func (r *SQLRepository) FindLocationsByCategory(ctx context.Context, cat *models.Category) (*models.Locations, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE category_id = $1 AND user_id = $2"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("FindLocationsByCategory: failed to prepare context for query %s. %w", query, err)
	}
	defer stmt.Close()

	user, ok := models.NewUserFromContext(ctx)
	if !ok {
		return nil, errors.New("FindLocationsByCategory: Failed to get user from context")
	}

	rows, err := stmt.QueryContext(ctx, cat.ID, user.ID)
	if err != nil {
		return nil, fmt.Errorf("FindLocationsByCategory: failed to query context for query %s. %w", query, err)
	}
	defer rows.Close()

	locs := make(models.Locations, 0)
	for rows.Next() {
		loc := new(models.Location)
		if err := rows.Scan(&loc.ID, &loc.Name, &loc.Address, &loc.Category, &loc.User); err != nil {
			return nil, fmt.Errorf("FindLocationsByCategory: failed to scan SQL row. %w", err)
		}
		locs = append(locs, loc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("FindLocationsByCategory: rows failed. %w", err)
	}

	return &locs, nil
}

// UpdateLocation updates specified location in repository
func (r *SQLRepository) UpdateLocation(ctx context.Context, loc *models.Location) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "UPDATE locations SET name = $1, address = $2, category_id = $3 WHERE id = $4 AND user_id = $5"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("UpdateLocation: failed to prepare context for query %s. %w", query, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, loc.Name, loc.Address, loc.Category, loc.ID, loc.User)
	if err != nil {
		return fmt.Errorf("UpdateLocation: failed to exec context for query %s. %w", query, err)
	}

	return nil
}

// DeleteLocation deletes location in repository
func (r *SQLRepository) DeleteLocation(ctx context.Context, id models.ID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "DELETE FROM locations WHERE id = $1 AND user_id = $2"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("DeleteLocation: failed to prepare context for query %s. %w", query, err)
	}
	defer stmt.Close()

	user, ok := models.NewUserFromContext(ctx)
	if !ok {
		return errors.New("DeleteLocation: Failed to get user from context")
	}

	_, err = stmt.ExecContext(ctx, id, user.ID)
	if err != nil {
		return fmt.Errorf("DeleteLocation: failed to exec context for query %s. %w", query, err)
	}

	return nil
}
