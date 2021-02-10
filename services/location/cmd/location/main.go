package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/edebernis/social-life-manager/services/location/cmd/location/config"
	"github.com/edebernis/social-life-manager/services/location/cmd/location/metrics"
	"github.com/edebernis/social-life-manager/services/location/internal/api"
	grpcapi "github.com/edebernis/social-life-manager/services/location/internal/api/grpc/v1"
	httpapi "github.com/edebernis/social-life-manager/services/location/internal/api/http/v1"
	sqlrepo "github.com/edebernis/social-life-manager/services/location/internal/repositories/sql"
	"github.com/edebernis/social-life-manager/services/location/internal/usecases"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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

func setupSQLRepository(registry *prometheus.Registry) (*sqlrepo.SQLRepository, error) {
	repo := sqlrepo.NewSQLRepository(&sqlrepo.Config{
		Host:            config.Config.SQL.Host,
		Port:            config.Config.SQL.Port,
		User:            config.Config.SQL.User,
		Password:        config.Config.SQL.Password,
		DBName:          config.Config.SQL.DB,
		ConnMaxIdleTime: config.Config.SQL.ConnMaxIdleTime,
		ConnMaxLifetime: config.Config.SQL.ConnMaxLifeTime,
		MaxIdleConns:    config.Config.SQL.MaxIdleConns,
		MaxOpenConns:    config.Config.SQL.MaxOpenConns,
		QueryTimeout:    config.Config.SQL.QueryTimeout,
	}, registry)

	if err := repo.Open(); err != nil {
		return nil, fmt.Errorf("Failed to open SQL repository. %w", err)
	}

	if err := repo.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("Failed to ping SQL repository. %w", err)
	}

	return repo, nil
}

func setupHTTPAPI(api *api.API, registry *prometheus.Registry) *httpapi.HTTPServer {
	return httpapi.NewHTTPServer(api, registry, &httpapi.Config{
		ReadHeaderTimeout: config.Config.API.HTTPReadHeaderTimeout,
		ReadTimeout:       config.Config.API.HTTPReadTimeout,
		WriteTimeout:      config.Config.API.HTTPWriteTimeout,
		JWTAlgorithm:      config.Config.JWT.Algorithm,
		JWTSecretKey:      config.Config.JWT.Secret,
	})
}

func setupGRPCAPI(api *api.API, registry *prometheus.Registry) *grpcapi.GRPCServer {
	return grpcapi.NewGRPCServer(api, registry, &grpcapi.Config{})
}

func setup() (*sqlrepo.SQLRepository, *httpapi.HTTPServer, *grpcapi.GRPCServer, *metrics.Server, error) {
	if err := config.LoadConfig(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Failed to load configuration. %w", err)
	}

	setupLogging()
	metricsServer := metrics.NewMetricsServer(config.Config.Metrics.Path)

	repo, err := setupSQLRepository(metricsServer.Registry)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Failed to setup SQL repository. %w", err)
	}

	usecase := usecases.NewLocationUsecase(repo)
	api := api.NewAPI(usecase)

	httpServer := setupHTTPAPI(api, metricsServer.Registry)
	grpcServer := setupGRPCAPI(api, metricsServer.Registry)

	return repo, httpServer, grpcServer, metricsServer, nil
}

func main() {
	repo, httpServer, grpcServer, metricsServer, err := setup()
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
		if err := httpServer.Serve(config.Config.API.HTTPBindAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Failed to start HTTP API server : %v", err)
		}
	}()

	logger.Infof("Start GRPC API server listening on address %s", config.Config.API.GRPCBindAddr)
	go func() {
		if err := httpServer.Serve(config.Config.API.GRPCBindAddr); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logger.Fatalf("Failed to start GRPC API server : %v", err)
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
	if err := grpcServer.Shutdown(); err != nil {
		logger.Errorf("Failed to shutdown gracefully API GRPC server. %s", err)
	}
	if err := repo.Close(); err != nil {
		logger.Errorf("Failed to close repository. %s", err)
	}
}
