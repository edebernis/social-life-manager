package sqlrepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dlmiddlecote/sqlstats"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// PostgreSQLDriver describes SQL driver for Postgres databases
	PostgreSQLDriver = "postgres"
)

// SQLRepository represents a repository using SQL to query
type SQLRepository struct {
	Config             *Config
	db                 *sql.DB
	prometheusRegistry prometheus.Registerer
}

// NewSQLRepository creates a new SQLRepository
func NewSQLRepository(config *Config, registry prometheus.Registerer) *SQLRepository {
	return &SQLRepository{
		config,
		nil,
		registry,
	}
}

// Config describes configs and options of SQL repository
type Config struct {
	Driver          string
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSL             bool
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

// Open opens DB handler
func (r *SQLRepository) Open() error {
	connectionStr, err := r.getConnectionString()
	if err != nil {
		return fmt.Errorf("Failed to set connection string. %w", err)
	}

	db, err := sql.Open(r.Config.Driver, connectionStr)
	if err != nil {
		return fmt.Errorf("Failed to open connection to repository. %w", err)
	}
	r.db = db

	r.db.SetConnMaxIdleTime(r.Config.ConnMaxIdleTime)
	r.db.SetConnMaxLifetime(r.Config.ConnMaxLifetime)
	r.db.SetMaxIdleConns(r.Config.MaxIdleConns)
	r.db.SetMaxOpenConns(r.Config.MaxOpenConns)

	collector := sqlstats.NewStatsCollector(r.Config.DBName, r.db)
	r.prometheusRegistry.MustRegister(collector)

	return nil
}

// Ping tests that DB is reachable over the network
func (r *SQLRepository) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := r.db.PingContext(ctx); err != nil {
		return fmt.Errorf("Failed to ping repository. %w", err)
	}

	return nil
}

func (r *SQLRepository) getConnectionString() (string, error) {
	switch r.Config.Driver {
	case PostgreSQLDriver:
		return r.newPostgresConnectionString(), nil
	default:
		return "", fmt.Errorf("Unknown SQL driver : %s", r.Config.Driver)
	}
}

func (r *SQLRepository) newPostgresConnectionString() string {
	var sslmode string
	if r.Config.SSL {
		sslmode = "enable"
	} else {
		sslmode = "disable"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		r.Config.Host,
		r.Config.Port,
		r.Config.User,
		r.Config.Password,
		r.Config.DBName,
		sslmode,
	)
}

// Close closes DB handler
func (r *SQLRepository) Close() error {
	if err := r.db.Close(); err != nil {
		return fmt.Errorf("Failed to close repository. %w", err)
	}

	return nil
}
