package sqlrepository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func newSQLMock(t *testing.T) (*SQLRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.FailNow()
	}
	repo := &SQLRepository{
		&Config{},
		db,
	}

	return repo, mock
}

func TestPostgresConnectionStringWithoutSSL(t *testing.T) {
	repo := NewSQLRepository(&Config{
		Driver:   PostgreSQLDriver,
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "test",
		SSL:      false,
	})

	connectionStr, err := repo.getConnectionString()

	assert.NoError(t, err)
	assert.Equal(t, connectionStr, "host=localhost port=5432 user=test password=test dbname=test sslmode=disable")
}

func TestPostgresConnectionStringWithSSL(t *testing.T) {
	repo := NewSQLRepository(&Config{
		Driver:   PostgreSQLDriver,
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "test",
		SSL:      true,
	})

	connectionStr, err := repo.getConnectionString()

	assert.NoError(t, err)
	assert.Equal(t, connectionStr, "host=localhost port=5432 user=test password=test dbname=test sslmode=enable")
}

func TestUnknownDriverConnectionString(t *testing.T) {
	repo := NewSQLRepository(&Config{
		Driver:   "unknown",
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "test",
		SSL:      false,
	})

	_, err := repo.getConnectionString()

	assert.Error(t, err)
}

func TestPingRepositoryWithError(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	mock.ExpectPing().WillReturnError(errors.New("failed"))

	err := repo.Ping()

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPingRepositoryWithSuccess(t *testing.T) {
	repo, mock := newSQLMock(t)
	defer repo.Close()

	mock.ExpectPing()

	err := repo.Ping()

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
