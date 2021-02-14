package grpcapi

import (
	"context"

	pb "github.com/edebernis/social-life-manager/services/location/api/grpc/v1"
	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/edebernis/social-life-manager/services/location/internal/usecases"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateCategory creates a new category
func (s *GRPCServer) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	cat := models.NewCategory(models.NewID(), req.Name)

	err := s.api.LocationUsecase.CreateCategory(ctx, cat)
	if err != nil {
		switch err {
		case usecases.ErrCategoryAlreadyExists:
			return nil, status.Errorf(codes.AlreadyExists, "category %s already exists", req.Name)
		default:
			logger.Errorf("CreateCategory: failed to create category. %v", err)
			return nil, status.Error(codes.Internal, "failed to create category")
		}
	}

	return &pb.CreateCategoryResponse{
		Category: &pb.Category{
			Id:   cat.ID.String(),
			Name: cat.Name,
		},
	}, nil
}
