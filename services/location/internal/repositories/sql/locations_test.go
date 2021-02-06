package sqlrepository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategoryWithPrepareError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	cat := models.NewCategory(models.NewID(), "Test Category")

	query := "INSERT INTO categories (id, name) VALUES ($1, $2)"
	mock.ExpectPrepare(query).WillReturnError(errors.New("failed"))

	err := repo.CreateCategory(newTestContext(), cat)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateCategoryWithExecError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	cat := models.NewCategory(models.NewID(), "Test Category")

	query := "INSERT INTO categories (id, name) VALUES ($1, $2)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(cat.ID, cat.Name).WillReturnError(errors.New("failed"))

	err := repo.CreateCategory(newTestContext(), cat)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateCategoryWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	cat := models.NewCategory(models.NewID(), "Test Category")

	query := "INSERT INTO categories (id, name) VALUES ($1, $2)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().
		WithArgs(cat.ID, cat.Name).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.CreateCategory(newTestContext(), cat)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCategoriesWithPrepareError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	query := "SELECT id, name FROM categories"
	mock.ExpectPrepare(query).WillReturnError(errors.New("failed"))

	_, err := repo.GetCategories(newTestContext())
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCategoriesWithQueryError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	query := "SELECT id, name FROM categories"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WillReturnError(errors.New("failed"))

	_, err := repo.GetCategories(newTestContext())
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCategoriesWithRowError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	cat := models.NewCategory(models.NewID(), "Test Category 1")

	rows := sqlmock.NewRows([]string{"id", "name"}).
		RowError(0, errors.New("failed")).
		AddRow(cat.ID, cat.Name)

	query := "SELECT id, name FROM categories"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WillReturnRows(rows)

	_, err := repo.GetCategories(newTestContext())
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCategoriesWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	cat1 := models.NewCategory(models.NewID(), "Test Category 1")
	cat2 := models.NewCategory(models.NewID(), "Test Category 2")

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(cat1.ID, cat1.Name).
		AddRow(cat2.ID, cat2.Name)

	query := "SELECT id, name FROM categories"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WillReturnRows(rows)

	cats, err := repo.GetCategories(newTestContext())
	assert.NoError(t, err)
	assert.ElementsMatch(t, *cats, models.Categories{cat1, cat2})
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindCategoryByIDWithPrepareError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	id := models.NewID()

	query := "SELECT id, name FROM categories WHERE id = $1"
	mock.ExpectPrepare(query).WillReturnError(errors.New("failed"))

	_, err := repo.FindCategoryByID(newTestContext(), id)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindCategoryByIDWithQueryError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	id := models.NewID()

	query := "SELECT id, name FROM categories WHERE id = $1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(id).WillReturnError(errors.New("failed"))

	_, err := repo.FindCategoryByID(newTestContext(), id)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindCategoryByIDWithRowError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	cat := models.NewCategory(models.NewID(), "Test Category 1")

	rows := sqlmock.NewRows([]string{"id", "name"}).
		RowError(0, errors.New("failed")).
		AddRow(cat.ID, cat.Name)

	query := "SELECT id, name FROM categories WHERE id = $1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(cat.ID).WillReturnRows(rows)

	_, err := repo.FindCategoryByID(newTestContext(), cat.ID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindCategoryByIDWithNoResult(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	id := models.NewID()

	query := "SELECT id, name FROM categories WHERE id = $1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(id).WillReturnRows(sqlmock.NewRows(nil))

	returnedCat, err := repo.FindCategoryByID(newTestContext(), id)
	assert.NoError(t, err)
	assert.Nil(t, returnedCat)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindCategoryByIDWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	cat := models.NewCategory(models.NewID(), "Test Category 1")

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(cat.ID, cat.Name)

	query := "SELECT id, name FROM categories WHERE id = $1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(cat.ID).WillReturnRows(rows)

	returnedCat, err := repo.FindCategoryByID(newTestContext(), cat.ID)
	assert.NoError(t, err)
	assert.Equal(t, cat, returnedCat)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindCategoryByNameWithPrepareError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	query := "SELECT id, name FROM categories WHERE name = $1"
	mock.ExpectPrepare(query).WillReturnError(errors.New("failed"))

	_, err := repo.FindCategoryByName(newTestContext(), "Test Category")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindCategoryByNameWithQueryError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	name := "Test Category"

	query := "SELECT id, name FROM categories WHERE name = $1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(name).WillReturnError(errors.New("failed"))

	_, err := repo.FindCategoryByName(newTestContext(), name)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindCategoryByNameWithRowError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	cat := models.NewCategory(models.NewID(), "Test Category 1")

	rows := sqlmock.NewRows([]string{"id", "name"}).
		RowError(0, errors.New("failed")).
		AddRow(cat.ID, cat.Name)

	query := "SELECT id, name FROM categories WHERE name = $1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(cat.Name).WillReturnRows(rows)

	_, err := repo.FindCategoryByName(newTestContext(), cat.Name)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindCategoryByNameWithNoResult(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	name := "Test Category"

	query := "SELECT id, name FROM categories WHERE name = $1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(name).WillReturnRows(sqlmock.NewRows(nil))

	returnedCat, err := repo.FindCategoryByName(newTestContext(), name)
	assert.NoError(t, err)
	assert.Nil(t, returnedCat)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindCategoryByNameWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	cat := models.NewCategory(models.NewID(), "Test Category 1")

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(cat.ID, cat.Name)

	query := "SELECT id, name FROM categories WHERE name = $1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(cat.Name).WillReturnRows(rows)

	returnedCat, err := repo.FindCategoryByName(newTestContext(), cat.Name)
	assert.NoError(t, err)
	assert.Equal(t, cat, returnedCat)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCategoryWithPrepareError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	cat := models.NewCategory(models.NewID(), "Test Category")

	query := "UPDATE categories SET name = $1 WHERE id = $2"
	mock.ExpectPrepare(query).WillReturnError(errors.New("failed"))

	err := repo.UpdateCategory(newTestContext(), cat)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCategoryWithExecError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	cat := models.NewCategory(models.NewID(), "Test Category")

	query := "UPDATE categories SET name = $1 WHERE id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(cat.Name, cat.ID).WillReturnError(errors.New("failed"))

	err := repo.UpdateCategory(newTestContext(), cat)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCategoryWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	cat := models.NewCategory(models.NewID(), "Test Category")

	query := "UPDATE categories SET name = $1 WHERE id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().
		WithArgs(cat.Name, cat.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateCategory(newTestContext(), cat)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCategoryWithPrepareError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	id := models.NewID()

	query := "DELETE FROM categories WHERE id = $1"
	mock.ExpectPrepare(query).WillReturnError(errors.New("failed"))

	err := repo.DeleteCategory(newTestContext(), id)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCategoryWithExecError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	id := models.NewID()

	query := "DELETE FROM categories WHERE id = $1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(id).WillReturnError(errors.New("failed"))

	err := repo.DeleteCategory(newTestContext(), id)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCategoryWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	id := models.NewID()

	query := "DELETE FROM categories WHERE id = $1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteCategory(newTestContext(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateLocationWithPrepareError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())

	query := "INSERT INTO locations (id, name, address, category_id, user_id) VALUES ($1, $2, $3, $4, $5)"
	mock.ExpectPrepare(query).WillReturnError(errors.New("failed"))

	err := repo.CreateLocation(newTestContext(), loc)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateLocationWithExecError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())

	query := "INSERT INTO locations (id, name, address, category_id, user_id) VALUES ($1, $2, $3, $4, $5)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(loc.ID, loc.Name, loc.Address, loc.Category, loc.User).WillReturnError(errors.New("failed"))

	err := repo.CreateLocation(newTestContext(), loc)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateLocationWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())

	query := "INSERT INTO locations (id, name, address, category_id, user_id) VALUES ($1, $2, $3, $4, $5)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().
		WithArgs(loc.ID, loc.Name, loc.Address, loc.Category, loc.User).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.CreateLocation(newTestContext(), loc)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLocationsWithPrepareError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE user_id = $1"
	mock.ExpectPrepare(query).WillReturnError(errors.New("failed"))

	_, err := repo.GetLocations(newTestContext())
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLocationsWithQueryError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE user_id = $1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WillReturnError(errors.New("failed"))

	_, err := repo.GetLocations(newTestContext())
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLocationsWithRowError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), user.ID)

	rows := sqlmock.NewRows([]string{"id", "name", "address", "category_id", "user_id"}).
		RowError(0, errors.New("failed")).
		AddRow(loc.ID, loc.Name, loc.Address, loc.Category, loc.User)

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE user_id = $1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WillReturnRows(rows)

	_, err := repo.GetLocations(ctx)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLocationsWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	loc1 := models.NewLocation(models.NewID(), "Test Location 1", "1 rue de la Poste, 75001 Paris", models.NewID(), user.ID)
	loc2 := models.NewLocation(models.NewID(), "Test Location 2", "2 rue de la Poste, 75001 Paris", models.NewID(), user.ID)

	rows := sqlmock.NewRows([]string{"id", "name", "address", "category_id", "user_id"}).
		AddRow(loc1.ID, loc1.Name, loc1.Address, loc1.Category, loc1.User).
		AddRow(loc2.ID, loc2.Name, loc2.Address, loc2.Category, loc2.User)

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE user_id = $1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WillReturnRows(rows)

	locs, err := repo.GetLocations(ctx)
	assert.NoError(t, err)
	assert.ElementsMatch(t, *locs, models.Locations{loc1, loc2})
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationByIDWithPrepareError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	id := models.NewID()

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE id = $1 AND user_id = $2"
	mock.ExpectPrepare(query).WillReturnError(errors.New("failed"))

	_, err := repo.FindLocationByID(newTestContext(), id)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationByIDWithQueryError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)
	id := models.NewID()

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE id = $1 AND user_id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(id, user.ID).WillReturnError(errors.New("failed"))

	_, err := repo.FindLocationByID(ctx, id)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationByIDWithRowError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), user.ID)

	rows := sqlmock.NewRows([]string{"id", "name", "address", "category_id", "user_id"}).
		RowError(0, errors.New("failed")).
		AddRow(loc.ID, loc.Name, loc.Address, loc.Category, loc.User)

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE id = $1 AND user_id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(loc.ID, user.ID).WillReturnRows(rows)

	_, err := repo.FindLocationByID(ctx, loc.ID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationByIDWithNoResult(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)
	id := models.NewID()

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE id = $1 AND user_id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(id, user.ID).WillReturnRows(sqlmock.NewRows(nil))

	returnedCat, err := repo.FindLocationByID(ctx, id)
	assert.NoError(t, err)
	assert.Nil(t, returnedCat)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationByIDWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), user.ID)

	rows := sqlmock.NewRows([]string{"id", "name", "address", "category_id", "user_id"}).
		AddRow(loc.ID, loc.Name, loc.Address, loc.Category, loc.User)

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE id = $1 AND user_id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(loc.ID, user.ID).WillReturnRows(rows)

	returnedLoc, err := repo.FindLocationByID(ctx, loc.ID)
	assert.NoError(t, err)
	assert.Equal(t, loc, returnedLoc)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationByNameWithPrepareError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE name = $1 AND user_id = $2"
	mock.ExpectPrepare(query).WillReturnError(errors.New("failed"))

	_, err := repo.FindLocationByName(newTestContext(), "test")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationByNameWithQueryError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)
	name := "test"

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE name = $1 AND user_id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(name, user.ID).WillReturnError(errors.New("failed"))

	_, err := repo.FindLocationByName(ctx, name)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationByNameWithRowError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), user.ID)

	rows := sqlmock.NewRows([]string{"id", "name", "address", "category_id", "user_id"}).
		RowError(0, errors.New("failed")).
		AddRow(loc.ID, loc.Name, loc.Address, loc.Category, loc.User)

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE name = $1 AND user_id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(loc.Name, user.ID).WillReturnRows(rows)

	_, err := repo.FindLocationByName(ctx, loc.Name)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationByNameWithNoResult(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)
	name := "test"

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE name = $1 AND user_id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(name, user.ID).WillReturnRows(sqlmock.NewRows(nil))

	returnedCat, err := repo.FindLocationByName(ctx, name)
	assert.NoError(t, err)
	assert.Nil(t, returnedCat)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationByNameWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), user.ID)

	rows := sqlmock.NewRows([]string{"id", "name", "address", "category_id", "user_id"}).
		AddRow(loc.ID, loc.Name, loc.Address, loc.Category, loc.User)

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE name = $1 AND user_id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(loc.Name, user.ID).WillReturnRows(rows)

	returnedLoc, err := repo.FindLocationByName(ctx, loc.Name)
	assert.NoError(t, err)
	assert.Equal(t, loc, returnedLoc)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationsByCategoryWithPrepareError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	cat := models.NewCategory(models.NewID(), "Test Category")

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE category_id = $1 AND user_id = $2"
	mock.ExpectPrepare(query).WillReturnError(errors.New("failed"))

	_, err := repo.FindLocationsByCategory(newTestContext(), cat)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationsByCategoryWithQueryError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	cat := models.NewCategory(models.NewID(), "Test Category")

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE category_id = $1 AND user_id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(cat.ID, user.ID).WillReturnError(errors.New("failed"))

	_, err := repo.FindLocationsByCategory(ctx, cat)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationsByCategoryWithRowError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	cat := models.NewCategory(models.NewID(), "Test Category")
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), user.ID)

	rows := sqlmock.NewRows([]string{"id", "name", "address", "category_id", "user_id"}).
		RowError(0, errors.New("failed")).
		AddRow(loc.ID, loc.Name, loc.Address, loc.Category, loc.User)

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE category_id = $1 AND user_id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(cat.ID, user.ID).WillReturnRows(rows)

	_, err := repo.FindLocationsByCategory(ctx, cat)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationsByCategoryWithNoResult(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	cat := models.NewCategory(models.NewID(), "Test Category")

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE category_id = $1 AND user_id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(cat.ID, user.ID).WillReturnRows(sqlmock.NewRows(nil))

	locs, err := repo.FindLocationsByCategory(ctx, cat)
	assert.NoError(t, err)
	assert.Empty(t, locs)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLocationsByCategoryWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	cat := models.NewCategory(models.NewID(), "Test Category")
	loc1 := models.NewLocation(models.NewID(), "Test Location 1", "1 rue de la Poste, 75001 Paris", cat.ID, user.ID)
	loc2 := models.NewLocation(models.NewID(), "Test Location 2", "2 rue de la Poste, 75001 Paris", cat.ID, user.ID)

	rows := sqlmock.NewRows([]string{"id", "name", "address", "category_id", "user_id"}).
		AddRow(loc1.ID, loc1.Name, loc1.Address, loc1.Category, loc1.User).
		AddRow(loc2.ID, loc2.Name, loc2.Address, loc2.Category, loc2.User)

	query := "SELECT id, name, address, category_id, user_id FROM locations WHERE category_id = $1 AND user_id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(cat.ID, user.ID).WillReturnRows(rows)

	locs, err := repo.FindLocationsByCategory(ctx, cat)
	assert.NoError(t, err)
	assert.ElementsMatch(t, *locs, models.Locations{loc1, loc2})
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateLocationWithPrepareError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), user.ID)

	query := "UPDATE locations SET name = $1, address = $2, category_id = $3 WHERE id = $4 AND user_id = $5"
	mock.ExpectPrepare(query).WillReturnError(errors.New("failed"))

	err := repo.UpdateLocation(ctx, loc)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateLocationWithExecError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), user.ID)

	query := "UPDATE locations SET name = $1, address = $2, category_id = $3 WHERE id = $4 AND user_id = $5"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().
		WithArgs(loc.Name, loc.Address, loc.Category, loc.ID, user.ID).
		WillReturnError(errors.New("failed"))

	err := repo.UpdateLocation(ctx, loc)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateLocationWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), user.ID)

	query := "UPDATE locations SET name = $1, address = $2, category_id = $3 WHERE id = $4 AND user_id = $5"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().
		WithArgs(loc.Name, loc.Address, loc.Category, loc.ID, user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateLocation(ctx, loc)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteLocationWithPrepareError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	id := models.NewID()

	query := "DELETE FROM locations WHERE id = $1 AND user_id = $2"
	mock.ExpectPrepare(query).WillReturnError(errors.New("failed"))

	err := repo.DeleteLocation(newTestContext(), id)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteLocationWithExecError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	id := models.NewID()

	query := "DELETE FROM locations WHERE id = $1 AND user_id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().
		WithArgs(id, user.ID).
		WillReturnError(errors.New("failed"))

	err := repo.DeleteLocation(ctx, id)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteLocationWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	ctx := newTestContext()
	user, _ := models.NewUserFromContext(ctx)

	id := models.NewID()

	query := "DELETE FROM locations WHERE id = $1 AND user_id = $2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().
		WithArgs(id, user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteLocation(ctx, id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
