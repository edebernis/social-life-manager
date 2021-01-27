package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/edebernis/social-life-manager/services/location/cmd/location/config"
	"github.com/edebernis/social-life-manager/services/location/internal/api"
	httpapi "github.com/edebernis/social-life-manager/services/location/internal/api/http"
	sqlrepo "github.com/edebernis/social-life-manager/services/location/internal/repositories/sql"
	"github.com/edebernis/social-life-manager/services/location/internal/usecases"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func setupInstrumenting() *prometheus.Registry {
	registry := prometheus.NewRegistry()

	handler := promhttp.InstrumentMetricHandler(
		registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{
			ErrorLog:      logger,
			ErrorHandling: promhttp.HTTPErrorOnError,
			Registry:      registry,
		}),
	)

	http.Handle("/metrics", handler)

	logger.Infof("Start HTTP metrics service listening on address %s", config.Config.MetricsBindAddr)
	go func() {
		err := http.ListenAndServe(config.Config.MetricsBindAddr, nil)
		if err != nil {
			logger.Fatalf("Failed to start HTTP metrics server : %v", err)
		}
	}()

	return registry
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

func setup() (*sqlrepo.SQLRepository, *httpapi.HTTPServer, error) {
	if err := config.LoadConfig(); err != nil {
		return nil, nil, fmt.Errorf("Failed to load configuration. %w", err)
	}

	setupLogging()
	registry := setupInstrumenting()

	repo, err := setupSQLRepository()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to setup SQL repository. %w", err)
	}

	server := setupHTTPAPI(repo, registry)

	return repo, server, nil
}

func main() {
	repo, server, err := setup()
	if err != nil {
		logger.Fatalf("Failed to setup application. %v", err)
	}
	defer repo.Close()

	logger.Infof("Start HTTP API listening on address %s", config.Config.HTTPBindAddr)
	if err := server.Serve(config.Config.HTTPBindAddr); err != nil {
		logger.Fatalf("Failed to start HTTP API server : %v", err)
	}
}
