package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.WithField("package", "metrics")
)

// Server serves metrics to allow fetching by external tools like Prometheus
type Server struct {
	Registry   *prometheus.Registry
	httpServer *http.Server
}

// Serve listens to specified address
func (s *Server) Serve(addr string) error {
	s.httpServer.Addr = addr
	return s.httpServer.ListenAndServe()
}

// Shutdown closes the server gracefully
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("Server forced to shutdown: %w", err)
	}

	return nil
}

// NewMetricsServer creates a new metrics server
func NewMetricsServer(path string) *Server {
	registry := prometheus.NewRegistry()
	handler := promhttp.InstrumentMetricHandler(
		registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{
			ErrorLog:      logger,
			ErrorHandling: promhttp.HTTPErrorOnError,
			Registry:      registry,
		}),
	)

	// Setup HTTP server to serve metrics
	mux := http.NewServeMux()
	mux.Handle(path, handler)

	return &Server{
		Registry: registry,
		httpServer: &http.Server{
			Handler: mux,
		},
	}
}
