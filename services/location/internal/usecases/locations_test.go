package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/edebernis/social-life-manager/services/location/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateLocationWithRepositoryCreateLocationError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())
	repo.On("CreateLocation", ctx, loc).Return(errors.New("failed"))
	repo.On("FindLocationByName", ctx, "Test Location").Return(nil, nil)

	err := usecase.CreateLocation(ctx, loc)
	assert.Error(t, err)
}

func TestCreateLocationWithRepositoryFindLocationByNameError(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())
	repo.On("FindLocationByName", ctx, "Test Location").Return(nil, errors.New("failed"))

	err := usecase.CreateLocation(ctx, loc)
	assert.Error(t, err)
}

func TestCreateLocationWithSuccess(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())
	repo.On("CreateLocation", ctx, loc).Return(nil)
	repo.On("FindLocationByName", ctx, "Test Location").Return(nil, nil)

	err := usecase.CreateLocation(ctx, loc)
	assert.NoError(t, err)
}

func TestCreateLocationWithNameAlreadyExisting(t *testing.T) {
	repo := new(mocks.LocationRepositoryMock)
	usecase := NewLocationUsecase(repo)

	ctx := context.Background()
	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())
	repo.On("FindLocationByName", ctx, "Test Location").Return(loc, nil)

	err := usecase.CreateLocation(ctx, loc)
	assert.Equal(t, err, ErrLocationAlreadyExists)
}
