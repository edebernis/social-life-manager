package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/edebernis/social-life-manager/location/api"
	httpapi "github.com/edebernis/social-life-manager/location/api/http"
	sqlrepo "github.com/edebernis/social-life-manager/location/repositories/sql"
	"github.com/edebernis/social-life-manager/location/usecases"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	configPath     string
	appAddress     string
	metricsAddress string

	logger = logrus.WithField("package", "main")
)

func setupFlagging() {
	flag.StringVar(&configPath, "C", "config.yml", "Path to config file")
	flag.StringVar(&appAddress, "H", ":8080", "Application HTTP network address")
	flag.StringVar(&metricsAddress, "M", ":2112", "Metrics HTTP network address")

	flag.Parse()
}

func setupLogging() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	if os.Getenv("DEBUG") == "1" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func setupInstrumenting() {
	http.Handle("/metrics", promhttp.Handler())

	logger.Infof("Start HTTP metrics service listening on address %s", metricsAddress)
	go func() {
		err := http.ListenAndServe(metricsAddress, nil)
		if err != nil {
			logger.Fatalf("Failed to start HTTP metrics server : %v", err)
		}
	}()
}

func setupConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Failed to load .env file. %w", err)
	}

	// If config file path, got from command-line, does not exist
	// try to get another config file path from env var
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = os.Getenv("CONFIG_FILE")
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open config file %s. %w", configPath, err)
	}
	defer file.Close()

	config, err := newConfig(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse config file %s. %w", configPath, err)
	}

	return config, nil
}

func newSQLConfig(config *Config) (*sqlrepo.Config, error) {
	port, err := strconv.Atoi(os.Getenv("POSTGRESQL_PORT"))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse POSTGRESQL_PORT integer. %w", err)
	}

	connMaxIdleTime, err := time.ParseDuration(config.SQL.ConnMaxIdleTime)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse SQL ConnMaxIdleTime config. %w", err)
	}

	ConnMaxLifetime, err := time.ParseDuration(config.SQL.ConnMaxLifetime)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse SQL ConnMaxLifetime config. %w", err)
	}

	return &sqlrepo.Config{
		Driver:          sqlrepo.PostgreSQLDriver,
		Host:            os.Getenv("POSTGRESQL_HOST"),
		Port:            port,
		User:            os.Getenv("POSTGRESQL_USER"),
		Password:        os.Getenv("POSTGRESQL_PASSWORD"),
		DBName:          os.Getenv("POSTGRESQL_DB"),
		ConnMaxIdleTime: connMaxIdleTime,
		ConnMaxLifetime: ConnMaxLifetime,
		MaxIdleConns:    config.SQL.maxIdleConns,
		MaxOpenConns:    config.SQL.maxOpenConns,
	}, nil
}

func setupSQLRepository(config *Config) (*sqlrepo.SQLRepository, error) {
	sqlConfig, err := newSQLConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to get SQL config. %w", err)
	}

	repo := sqlrepo.NewSQLRepository(sqlConfig)

	if err := repo.Open(); err != nil {
		return nil, fmt.Errorf("Failed to open SQL repository. %w", err)
	}

	if err := repo.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping SQL repository. %w", err)
	}

	return repo, nil
}

func setupHTTPAPI(repo *sqlrepo.SQLRepository, config *Config) *httpapi.HTTPServer {
	usecase := usecases.NewLocationUsecase(repo)
	api := api.NewAPI(usecase)

	server := httpapi.NewHTTPServer(api, &httpapi.Config{
		JWTAlgorithm: config.JWT.Algorithm,
		JWTSecretKey: os.Getenv("JWT_SECRET"),
	})

	return server
}

func setup() (*sqlrepo.SQLRepository, *httpapi.HTTPServer, error) {
	setupFlagging()
	setupLogging()
	setupInstrumenting()

	config, err := setupConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to setup configuration. %w", err)
	}

	repo, err := setupSQLRepository(config)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to setup SQL repository. %w", err)
	}

	server := setupHTTPAPI(repo, config)

	return repo, server, nil
}

func main() {
	repo, server, err := setup()
	if err != nil {
		logger.Fatalf("Failed to setup application. %v", err)
	}
	defer repo.Close()

	logger.Infof("Start HTTP API listening on address %s", appAddress)
	if err := server.Serve(appAddress); err != nil {
		logger.Fatalf("Failed to start HTTP API server : %v", err)
	}
}
