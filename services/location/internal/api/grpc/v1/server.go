package grpcapi

import (
	"context"
	"fmt"
	"net"

	pb "github.com/edebernis/social-life-manager/services/location/api/grpc/v1"
	"github.com/edebernis/social-life-manager/services/location/internal/api"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

// Config holds gRPC server configuration parameters
type Config struct{}

// GRPCServer runs the gRPC service. It implements the Server interface.
type GRPCServer struct {
	pb.UnimplementedLocationServiceServer

	server *grpc.Server
}

// NewGRPCServer builds and register a new gRPC server
func NewGRPCServer(api *api.API, registry prometheus.Registerer, config *Config) *GRPCServer {
	s := &GRPCServer{
		server: grpc.NewServer(),
	}
	pb.RegisterLocationServiceServer(s.server, s)

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

// CreateCategory creates a new category
func (s *GRPCServer) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	return nil, nil
}
