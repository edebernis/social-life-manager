package grpcapi

import (
	"fmt"
	"net"
	"time"

	pb "github.com/edebernis/social-life-manager/services/location/api/grpc/v1"
	"github.com/edebernis/social-life-manager/services/location/internal/api"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	logger = logrus.WithFields(logrus.Fields{"package": "grpcapi"})
)

// Config holds gRPC server configuration parameters
type Config struct {
	ConnectionTimeout time.Duration

	// Signing algorithm used for JWT
	JWTAlgorithm string
	// Key to check JWT signature
	JWTSecretKey string
}

// GRPCServer runs the gRPC service. It implements the Server interface.
type GRPCServer struct {
	pb.UnimplementedLocationServiceServer

	api      *api.API
	registry prometheus.Registerer
	server   *grpc.Server
}

// NewGRPCServer builds and register a new gRPC server
func NewGRPCServer(api *api.API, auth api.Authenticator, registry prometheus.Registerer, config *Config) *GRPCServer {
	metricsMW := newMetricsMiddleware(registry)

	s := &GRPCServer{
		api:      api,
		registry: registry,
		server: grpc.NewServer(
			grpc.ConnectionTimeout(config.ConnectionTimeout),
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				metricsMW.unaryServerInterceptor(),
				grpc_validator.UnaryServerInterceptor(),
				grpc_auth.UnaryServerInterceptor(newAuthHandlerFunc(auth)),
				grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(newRecoveryHandlerFunc())),
			)),
		),
	}

	pb.RegisterLocationServiceServer(s.server, s)
	metricsMW.initializeMetrics(s.server)

	return s
}

// Serve runs the server and listen for incoming requests
func (s *GRPCServer) Serve(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	return s.server.Serve(lis)
}

// Shutdown stops gracefully the server
func (s *GRPCServer) Shutdown() error {
	s.server.GracefulStop()
	return nil
}
