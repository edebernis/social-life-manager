package sqlrepository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func newSQLMock(t *testing.T) (*SQLRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(
		sqlmock.MonitorPingsOption(true),
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		t.FailNow()
	}
	repo := &SQLRepository{
		&Config{QueryTimeout: 5 * time.Second},
		db,
		prometheus.NewRegistry(),
	}

	return repo, mock
}

func newTestContext() context.Context {
	user := models.NewUser(models.NewID(), "user@test.fr")
	return models.NewContextWithUser(context.Background(), user)
}

func TestPostgresDSNWithoutSSL(t *testing.T) {
	repo := NewSQLRepository(&Config{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "test",
		SSL:      false,
	}, nil)

	assert.Equal(t, repo.dsn(), "host=localhost port=5432 user=test password=test dbname=test sslmode=disable")
}

func TestPostgresDSNWithSSL(t *testing.T) {
	repo := NewSQLRepository(&Config{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "test",
		SSL:      true,
	}, nil)

	assert.Equal(t, repo.dsn(), "host=localhost port=5432 user=test password=test dbname=test sslmode=enable")
}

func TestPingRepositoryWithError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	mock.ExpectPing().WillReturnError(errors.New("failed"))

	err := repo.Ping(context.Background())

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPingRepositoryWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	mock.ExpectPing()

	err := repo.Ping(context.Background())

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseRepositoryWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)

	mock.ExpectClose()

	err := repo.Close()

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseRepositoryWithError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	mock.ExpectClose().WillReturnError(errors.New("failed"))

	err := repo.Close()

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
