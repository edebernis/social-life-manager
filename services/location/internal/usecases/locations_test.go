package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/edebernis/social-life-manager/services/location/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategoryWithRepositoryFindCategoryByNameError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("CreateCategory", ctx, cat).Return(nil)
	repo.On("FindCategoryByName", ctx, "Test Category").Return(nil, errors.New("failed"))

	err := usecase.CreateCategory(ctx, cat)
	assert.Error(t, err)
}

func TestCreateCategoryWithRepositoryCreateCategoryError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("CreateCategory", ctx, cat).Return(errors.New("failed"))
	repo.On("FindCategoryByName", ctx, "Test Category").Return(nil, nil)

	err := usecase.CreateCategory(ctx, cat)
	assert.Error(t, err)
}

func TestCreateCategoryWithNameAlreadyExisting(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("CreateCategory", ctx, cat).Return(nil)
	repo.On("FindCategoryByName", ctx, "Test Category").Return(cat, nil)

	err := usecase.CreateCategory(ctx, cat)
	assert.Equal(t, err, ErrCategoryAlreadyExists)
}

func TestCreateCategoryWithSuccess(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("CreateCategory", ctx, cat).Return(nil)
	repo.On("FindCategoryByName", ctx, "Test Category").Return(nil, nil)

	err := usecase.CreateCategory(ctx, cat)
	assert.NoError(t, err)
}

func TestGetCategoriesWithRepositoryError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	repo.On("GetCategories", ctx).Return(nil, errors.New("failed"))

	cats, err := usecase.GetCategories(ctx)
	assert.Error(t, err)
	assert.Nil(t, cats)
}

func TestGetCategoriesWithSuccess(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cats := models.Categories{
		models.NewCategory(models.NewID(), "Test Category"),
	}
	repo.On("GetCategories", ctx).Return(&cats, nil)

	returnedCategories, err := usecase.GetCategories(ctx)
	assert.NoError(t, err)
	assert.Equal(t, cats, *returnedCategories)
}

func TestFindCategoryByIDWithRepositoryError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	id := models.NewID()
	repo.On("FindCategoryByID", ctx, id).Return(nil, errors.New("failed"))

	cat, err := usecase.FindCategoryByID(ctx, id)
	assert.Error(t, err)
	assert.Nil(t, cat)
}

func TestFindCategoryByIDWithSuccess(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)

	returnedCat, err := usecase.FindCategoryByID(ctx, cat.ID)
	assert.NoError(t, err)
	assert.Equal(t, cat, returnedCat)
}

func TestUpdateCategoryWithCategoryNotFoundError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("FindCategoryByID", ctx, cat.ID).Return(nil, nil)
	repo.On("UpdateCategory", ctx, cat).Return(nil)

	err := usecase.UpdateCategory(ctx, cat)
	assert.Equal(t, err, ErrCategoryNotFound)
}

func TestUpdateCategoryWithRepositoryUpdateCategoryError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)
	repo.On("UpdateCategory", ctx, cat).Return(errors.New("failed"))

	err := usecase.UpdateCategory(ctx, cat)
	assert.Error(t, err)
}

func TestUpdateCategoryWithRepositoryFindCategoryByIDError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("FindCategoryByID", ctx, cat.ID).Return(nil, errors.New("failed"))
	repo.On("UpdateCategory", ctx, cat).Return(nil)

	err := usecase.UpdateCategory(ctx, cat)
	assert.Error(t, err)
}

func TestUpdateCategoryWithSuccess(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)
	repo.On("UpdateCategory", ctx, cat).Return(nil)

	err := usecase.UpdateCategory(ctx, cat)
	assert.NoError(t, err)
}

func TestDeleteCategoryWithRepositoryFindCategoryByIDError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("FindCategoryByID", ctx, cat.ID).Return(nil, errors.New("failed"))
	repo.On("DeleteCategory", ctx, cat.ID).Return(nil)

	err := usecase.DeleteCategory(ctx, cat.ID)
	assert.Error(t, err)
}

func TestDeleteCategoryWithCategoryNotFoundError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("FindCategoryByID", ctx, cat.ID).Return(nil, nil)
	repo.On("DeleteCategory", ctx, cat.ID).Return(nil)

	err := usecase.DeleteCategory(ctx, cat.ID)
	assert.Equal(t, err, ErrCategoryNotFound)
}

func TestDeleteCategoryWithRepositoryDeleteCategoryError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)
	repo.On("DeleteCategory", ctx, cat.ID).Return(errors.New(("failed")))

	err := usecase.DeleteCategory(ctx, cat.ID)
	assert.Error(t, err)
}

func TestDeleteCategoryWithSuccess(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)
	repo.On("DeleteCategory", ctx, cat.ID).Return(nil)

	err := usecase.DeleteCategory(ctx, cat.ID)
	assert.NoError(t, err)
}

func TestCreateLocationWithRepositoryCreateLocationError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", cat.ID, models.NewID())
	repo.On("CreateLocation", ctx, loc).Return(errors.New("failed"))
	repo.On("FindLocationByName", ctx, "Test Location").Return(nil, nil)
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)

	err := usecase.CreateLocation(ctx, loc)
	assert.Error(t, err)
}

func TestCreateLocationWithRepositoryFindLocationByNameError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", cat.ID, models.NewID())
	repo.On("CreateLocation", ctx, loc).Return(nil)
	repo.On("FindLocationByName", ctx, "Test Location").Return(nil, errors.New("failed"))
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)

	err := usecase.CreateLocation(ctx, loc)
	assert.Error(t, err)
}

func TestCreateLocationWithRepositoryFindLocationByIDError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", cat.ID, models.NewID())
	repo.On("CreateLocation", ctx, loc).Return(nil)
	repo.On("FindLocationByName", ctx, "Test Location").Return(nil, nil)
	repo.On("FindCategoryByID", ctx, cat.ID).Return(nil, errors.New("failed"))

	err := usecase.CreateLocation(ctx, loc)
	assert.Error(t, err)
}

func TestCreateLocationWithNameAlreadyExisting(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", cat.ID, models.NewID())
	repo.On("CreateLocation", ctx, loc).Return(nil)
	repo.On("FindLocationByName", ctx, "Test Location").Return(loc, nil)
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)

	err := usecase.CreateLocation(ctx, loc)
	assert.Equal(t, err, ErrLocationAlreadyExists)
}

func TestCreateLocationWithCategoryNotFound(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", cat.ID, models.NewID())
	repo.On("CreateLocation", ctx, loc).Return(nil)
	repo.On("FindLocationByName", ctx, "Test Location").Return(nil, nil)
	repo.On("FindCategoryByID", ctx, cat.ID).Return(nil, nil)

	err := usecase.CreateLocation(ctx, loc)
	assert.Equal(t, err, ErrCategoryNotFound)
}

func TestCreateLocationWithSuccess(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", cat.ID, models.NewID())
	repo.On("CreateLocation", ctx, loc).Return(nil)
	repo.On("FindLocationByName", ctx, "Test Location").Return(nil, nil)
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)

	err := usecase.CreateLocation(ctx, loc)
	assert.NoError(t, err)
}

func TestGetLocationsWithRepositoryError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	repo.On("GetLocations", ctx).Return(nil, errors.New("failed"))

	locations, err := usecase.GetLocations(ctx)
	assert.Error(t, err)
	assert.Nil(t, locations)
}

func TestGetLocationsWithSuccess(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	locations := models.Locations{
		models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID()),
	}
	repo.On("GetLocations", ctx).Return(&locations, nil)

	returnedLocations, err := usecase.GetLocations(ctx)
	assert.NoError(t, err)
	assert.Equal(t, locations, *returnedLocations)
}

func TestFindLocationByIDWithRepositoryError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	id := models.NewID()
	repo.On("FindLocationByID", ctx, id).Return(nil, errors.New("failed"))

	location, err := usecase.FindLocationByID(ctx, id)
	assert.Error(t, err)
	assert.Nil(t, location)
}

func TestFindLocationByIDWithSuccess(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())
	repo.On("FindLocationByID", ctx, loc.ID).Return(loc, nil)

	returnedLoc, err := usecase.FindLocationByID(ctx, loc.ID)
	assert.NoError(t, err)
	assert.Equal(t, loc, returnedLoc)
}

func TestFindLocationsByCategoryWithRepositoryFindCategoryByIDError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("FindCategoryByID", ctx, cat.ID).Return(nil, errors.New("failed"))
	repo.On("FindLocationsByCategory", ctx, cat).Return(nil, nil)

	locations, err := usecase.FindLocationsByCategory(ctx, cat.ID)
	assert.Error(t, err)
	assert.Nil(t, locations)
}

func TestFindLocationsByCategoryWithCategoryNotFound(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("FindCategoryByID", ctx, cat.ID).Return(nil, nil)
	repo.On("FindLocationsByCategory", ctx, cat).Return(nil, nil)

	locations, err := usecase.FindLocationsByCategory(ctx, cat.ID)
	assert.Equal(t, err, ErrCategoryNotFound)
	assert.Nil(t, locations)
}

func TestFindLocationsByCategoryWithRepositoryFindLocationsByCategoryError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)
	repo.On("FindLocationsByCategory", ctx, cat).Return(nil, errors.New("failed"))

	locations, err := usecase.FindLocationsByCategory(ctx, cat.ID)
	assert.Error(t, err)
	assert.Nil(t, locations)
}

func TestFindLocationsByCategoryWithSuccess(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	locs := models.Locations{
		models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", cat.ID, models.NewID()),
	}
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)
	repo.On("FindLocationsByCategory", ctx, cat).Return(&locs, nil)

	returnedLocs, err := usecase.FindLocationsByCategory(ctx, cat.ID)
	assert.NoError(t, err)
	assert.Equal(t, locs, *returnedLocs)
}

func TestUpdateLocationWithRepositoryFindLocationByIDError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", cat.ID, models.NewID())
	repo.On("FindLocationByID", ctx, loc.ID).Return(nil, errors.New("failed"))
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)
	repo.On("UpdateLocation", ctx, loc).Return(nil)

	err := usecase.UpdateLocation(ctx, loc)
	assert.Error(t, err)
}

func TestUpdateLocationWithLocationNotFoundError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", cat.ID, models.NewID())
	repo.On("FindLocationByID", ctx, loc.ID).Return(nil, nil)
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)
	repo.On("UpdateLocation", ctx, loc).Return(nil)

	err := usecase.UpdateLocation(ctx, loc)
	assert.Equal(t, err, ErrLocationNotFound)
}

func TestUpdateLocationWithRepositoryUpdateLocationError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", cat.ID, models.NewID())
	repo.On("FindLocationByID", ctx, loc.ID).Return(loc, nil)
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)
	repo.On("UpdateLocation", ctx, loc).Return(errors.New("failed"))

	err := usecase.UpdateLocation(ctx, loc)
	assert.Error(t, err)
}

func TestUpdateLocationWithRepositoryFindCategoryByIDError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", cat.ID, models.NewID())
	repo.On("FindLocationByID", ctx, loc.ID).Return(loc, nil)
	repo.On("FindCategoryByID", ctx, cat.ID).Return(nil, errors.New("failed"))
	repo.On("UpdateLocation", ctx, loc).Return(nil)

	err := usecase.UpdateLocation(ctx, loc)
	assert.Error(t, err)
}

func TestUpdateLocationWithSuccess(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	cat := models.NewCategory(models.NewID(), "Test Category")
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", cat.ID, models.NewID())
	repo.On("FindLocationByID", ctx, loc.ID).Return(loc, nil)
	repo.On("FindCategoryByID", ctx, cat.ID).Return(cat, nil)
	repo.On("UpdateLocation", ctx, loc).Return(nil)

	err := usecase.UpdateLocation(ctx, loc)
	assert.NoError(t, err)
}

func TestDeleteLocationWithRepositoryFindLocationByIDError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())
	repo.On("FindLocationByID", ctx, loc.ID).Return(nil, errors.New("failed"))
	repo.On("DeleteLocation", ctx, loc.ID).Return(nil)

	err := usecase.DeleteLocation(ctx, loc.ID)
	assert.Error(t, err)
}

func TestDeleteLocationWithLocationNotFoundError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())
	repo.On("FindLocationByID", ctx, loc.ID).Return(nil, nil)
	repo.On("DeleteLocation", ctx, loc.ID).Return(nil)

	err := usecase.DeleteLocation(ctx, loc.ID)
	assert.Equal(t, err, ErrLocationNotFound)
}

func TestDeleteLocationWithRepositoryDeleteLocationError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())
	repo.On("FindLocationByID", ctx, loc.ID).Return(loc, nil)
	repo.On("DeleteLocation", ctx, loc.ID).Return(errors.New("failed"))

	err := usecase.DeleteLocation(ctx, loc.ID)
	assert.Error(t, err)
}

func TestDeleteLocationWithSuccess(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())
	repo.On("FindLocationByID", ctx, loc.ID).Return(loc, nil)
	repo.On("DeleteLocation", ctx, loc.ID).Return(nil)

	err := usecase.DeleteLocation(ctx, loc.ID)
	assert.NoError(t, err)
}
