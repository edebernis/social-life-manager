package sqlrepository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/edebernis/social-life-manager/location/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateLocationWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())

	query := "INSERT INTO locations \\(id, name, address, category_id, user_id\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\, \\$5\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().
		WithArgs(loc.ID.String(), loc.Name, loc.Address, loc.Category.String(), loc.User.String()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.CreateLocation(context.Background(), loc)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateLocationWithPrepareError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())

	query := "INSERT INTO locations \\(id, name, address, category_id, user_id\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\, \\$5\\)"
	mock.ExpectPrepare(query).WillReturnError(errors.New("failed"))

	err := repo.CreateLocation(context.Background(), loc)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateLocationWithExecError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	loc := models.NewLocation(models.NewID(), "Test Location", "1 rue de la Poste, 75001 Paris", models.NewID(), models.NewID())

	query := "INSERT INTO locations \\(id, name, address, category_id, user_id\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\, \\$5\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(loc.ID.String(), loc.Name, loc.Address, loc.Category.String(), loc.User.String()).WillReturnError(errors.New("failed"))

	err := repo.CreateLocation(context.Background(), loc)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
