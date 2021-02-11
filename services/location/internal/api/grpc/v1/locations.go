package grpcapi

import (
	"context"

	pb "github.com/edebernis/social-life-manager/services/location/api/grpc/v1"
)

// CreateCategory creates a new category
func (s *GRPCServer) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	return nil, nil
}
