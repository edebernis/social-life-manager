package grpcapi

import (
	"context"
	"fmt"

	"github.com/edebernis/social-life-manager/services/location/internal/api"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newRecoveryHandlerFunc() grpc_recovery.RecoveryHandlerFuncContext {
	return func(ctx context.Context, p interface{}) (err error) {
		logger.Errorf("recoveryMiddleware: panic recovered. %v", p)
		return status.Error(codes.Internal, "internal server error")
	}
}

func newAuthHandlerFunc(auth Authenticator) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		creds, err := auth.CredentialsFromContext(ctx)
		if err != nil {
			logger.Errorf("authMiddleware: unable to retrieve creds from context. %v", err)
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth: %v", err)
		}

		newCtx, err := auth.Authenticate(ctx, creds)
		if err != nil {
			logger.Errorf("authMiddleware: unable to authenticate user using provided credentials. %v", err)
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth: %v", err)
		}

		return newCtx, nil
	}
}

type metricsMiddleware struct {
	metrics   *grpc_prometheus.ServerMetrics
	namespace string
	subsystem string
}

func newMetricsMiddleware(registry prometheus.Registerer) *metricsMiddleware {
	m := &metricsMiddleware{
		metrics:   grpc_prometheus.NewServerMetrics(),
		namespace: "api",
		subsystem: "grpc",
	}

	registry.MustRegister(m.metrics)
	return m
}

func (m *metricsMiddleware) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return m.metrics.UnaryServerInterceptor()
}

func (m *metricsMiddleware) initializeMetrics(s *grpc.Server) {
	m.metrics.InitializeMetrics(s)
}

// Authenticator describes authentication mechanism for GRPC API
type Authenticator interface {
	CredentialsFromContext(context.Context) (interface{}, error)
	Authenticate(ctx context.Context, credentials interface{}) (context.Context, error)
}

// JWTAuthenticator is a middleware authenticating user using JWT tokens
type JWTAuthenticator struct {
	*api.JWTAuthenticator

	// Auth scheme (e.g. `basic`, `bearer`), in a case-insensitive format (see rfc2617, sec 1.2)
	Scheme string
}

// NewJWTAuthenticator creates a new JWTAuthenticator
func NewJWTAuthenticator(scheme, algorithm, secretKey string) *JWTAuthenticator {
	return &JWTAuthenticator{
		&api.JWTAuthenticator{
			Algorithm: algorithm,
			SecretKey: secretKey,
		},
		scheme,
	}
}

// CredentialsFromContext extracts JWT token from request metadata
func (a *JWTAuthenticator) CredentialsFromContext(ctx context.Context) (interface{}, error) {
	token, err := grpc_auth.AuthFromMD(ctx, a.Scheme)
	if err != nil {
		return nil, fmt.Errorf("JWTAuthenticator: unable to retrieve token from metadata. %w", err)
	}

	return token, nil
}
