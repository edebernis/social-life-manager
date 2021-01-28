package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/edebernis/social-life-manager/services/location/cmd/location/config"
	"github.com/edebernis/social-life-manager/services/location/cmd/location/metrics"
	"github.com/edebernis/social-life-manager/services/location/internal/api"
	httpapi "github.com/edebernis/social-life-manager/services/location/internal/api/http"
	sqlrepo "github.com/edebernis/social-life-manager/services/location/internal/repositories/sql"
	"github.com/edebernis/social-life-manager/services/location/internal/usecases"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.WithField("package", "main")
)

func setupLogging() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	if config.Config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func setupSQLRepository() (*sqlrepo.SQLRepository, error) {
	repo := sqlrepo.NewSQLRepository(&sqlrepo.Config{
		Driver:          sqlrepo.PostgreSQLDriver,
		Host:            config.Config.SQL.Host,
		Port:            config.Config.SQL.Port,
		User:            config.Config.SQL.User,
		Password:        config.Config.SQL.Password,
		DBName:          config.Config.SQL.DB,
		ConnMaxIdleTime: config.Config.SQL.ConnMaxIdleTime,
		ConnMaxLifetime: config.Config.SQL.ConnMaxLifeTime,
		MaxIdleConns:    config.Config.SQL.MaxIdleConns,
		MaxOpenConns:    config.Config.SQL.MaxOpenConns,
	})

	if err := repo.Open(); err != nil {
		return nil, fmt.Errorf("Failed to open SQL repository. %w", err)
	}

	if err := repo.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping SQL repository. %w", err)
	}

	return repo, nil
}

func setupHTTPAPI(repo *sqlrepo.SQLRepository, registry *prometheus.Registry) *httpapi.HTTPServer {
	usecase := usecases.NewLocationUsecase(repo)
	api := api.NewAPI(usecase)

	server := httpapi.NewHTTPServer(api, registry, &httpapi.Config{
		JWTAlgorithm: config.Config.JWT.Algorithm,
		JWTSecretKey: config.Config.JWT.Secret,
	})

	return server
}

func setup() (*sqlrepo.SQLRepository, *httpapi.HTTPServer, *metrics.Server, error) {
	if err := config.LoadConfig(); err != nil {
		return nil, nil, nil, fmt.Errorf("Failed to load configuration. %w", err)
	}

	setupLogging()

	repo, err := setupSQLRepository()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Failed to setup SQL repository. %w", err)
	}

	metricsServer := metrics.NewMetricsServer(config.Config.Metrics.Path)
	httpServer := setupHTTPAPI(repo, metricsServer.Registry)

	return repo, httpServer, metricsServer, nil
}

func main() {
	repo, httpServer, metricsServer, err := setup()
	if err != nil {
		logger.Fatalf("Failed to setup application. %v", err)
	}

	logger.Infof("Start HTTP Metrics server listening on address %s", config.Config.Metrics.BindAddr)
	go func() {
		if err := metricsServer.Serve(config.Config.Metrics.BindAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Failed to start HTTP Metrics server : %v", err)
		}
	}()

	logger.Infof("Start HTTP API server listening on address %s", config.Config.API.HTTPBindAddr)
	go func() {
		if httpServer.Serve(config.Config.API.HTTPBindAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Failed to start HTTP API server : %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down")

	if err := metricsServer.Shutdown(); err != nil {
		logger.Errorf("Failed to shutdown gracefully metrics server. %s", err)
	}
	if err := httpServer.Shutdown(); err != nil {
		logger.Errorf("Failed to shutdown gracefully API HTTP server. %s", err)
	}
	if err := repo.Close(); err != nil {
		logger.Errorf("Failed to close repository. %s", err)
	}
}
