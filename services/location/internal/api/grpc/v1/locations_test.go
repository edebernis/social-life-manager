package grpcapi

import (
	"context"
	"os"
	"testing"

	pb "github.com/edebernis/social-life-manager/services/location/api/grpc/v1"
	"github.com/edebernis/social-life-manager/services/location/internal/api/mocks"
	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/edebernis/social-life-manager/services/location/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

var (
	client pb.LocationServiceClient
	server *GRPCServer
)

func TestMain(m *testing.M) {
	var conn *grpc.ClientConn
	server, conn = newTestGRPCClientConnection()
	defer conn.Close()

	client = pb.NewLocationServiceClient(conn)

	retCode := m.Run()
	os.Exit(retCode)
}

func TestCreateCategoryWithSuccess(t *testing.T) {
	server.api.LocationUsecase.(*mocks.LocationUsecaseMock).
		On("CreateCategory", utils.MockContextMatcher, mock.AnythingOfType("*models.Category")).
		Return(nil)

	request := &pb.CreateCategoryRequest{Name: "Test Category"}
	response, err := client.CreateCategory(context.Background(), request)

	assert.NoError(t, err)
	if assert.NotNil(t, response) {
		_, err := models.ParseID(response.Category.Id)
		assert.NoError(t, err)
		assert.Equal(t, request.Name, response.Category.Name)
	}
}
