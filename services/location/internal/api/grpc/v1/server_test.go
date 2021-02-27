package grpcapi

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/edebernis/social-life-manager/services/location/internal/api"
	"github.com/edebernis/social-life-manager/services/location/internal/api/mocks"
	"github.com/edebernis/social-life-manager/services/location/internal/models"
	"github.com/edebernis/social-life-manager/services/location/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func newTestGRPCClientConnection() (*GRPCServer, *grpc.ClientConn) {
	api := api.NewAPI(new(mocks.LocationUsecaseMock))
	auth := NewJWTAuthenticator(
		"bearer",
		jwt.SigningMethodHS256.Name,
		"secret",
	)
	s := NewGRPCServer(api, auth, prometheus.NewRegistry(), &Config{})

	listener := bufconn.Listen(1024 * 1024)
	go func() {
		if err := s.server.Serve(listener); err != nil {
			logger.Fatal(err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "",
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(
			&testJWTCredentials{
				auth.Algorithm,
				auth.SecretKey,
			},
		),
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return listener.Dial()
			},
		),
	)
	if err != nil {
		logger.Fatal(err)
	}

	return s, conn
}

type testJWTCredentials struct {
	Algorithm string
	SecretKey string
}

func (t *testJWTCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	claims := &api.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Add(time.Hour * -1).Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "test",
			Subject:   models.NewID().String(),
		},

		Email: "test@no-reply.com",
	}

	token, err := utils.NewJWTToken(t.SecretKey, t.Algorithm, claims)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", token),
	}, nil
}

func (testJWTCredentials) RequireTransportSecurity() bool {
	return false
}
