package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/edebernis/social-life-manager/services/location/internal/models"
)

// JWTAuthenticator authenticates request using JWT tokens
type JWTAuthenticator struct {
	Algorithm string
	SecretKey string
}

// JWTClaims describes that can be set in JWT token
type JWTClaims struct {
	StdClaims jwt.StandardClaims

	Email string `json:"email,omitempty"`
}

// Valid implements Claims interface
func (c JWTClaims) Valid() error {
	return c.StdClaims.Valid()
}

// Authenticate user using supplied JWT token
func (a *JWTAuthenticator) Authenticate(ctx context.Context, credentials interface{}) (context.Context, error) {
	tokenString, ok := credentials.(string)
	if !ok {
		return ctx, errors.New("credentials is not a JWT token string")
	}

	token, err := a.parseToken(tokenString)
	if err != nil {
		return ctx, fmt.Errorf("invalid JWT token. %w", err)
	}

	newCtx, err := a.newContextWithClaims(ctx, token)
	if err != nil {
		return ctx, fmt.Errorf("invalid JWT token. %w", err)
	}

	return newCtx, nil
}

func (a *JWTAuthenticator) parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(a.Algorithm) != token.Method {
			return nil, fmt.Errorf("Invalid signing algorithm for JWT token : %s", token.Header["alg"])
		}
		return []byte(a.SecretKey), nil
	})
}

func (a *JWTAuthenticator) newContextWithClaims(ctx context.Context, token *jwt.Token) (context.Context, error) {
	claims := token.Claims.(*JWTClaims)

	userID, err := models.ParseID(claims.StdClaims.Subject)
	if err != nil {
		return ctx, fmt.Errorf("Invalid user ID in JWT token subject : %s. %w", claims.StdClaims.Subject, err)
	}
	if userID == models.NilID {
		return ctx, errors.New("User ID cannot be nil")
	}

	return models.NewContextWithUser(
		ctx,
		models.NewUser(userID, claims.Email),
	), nil
}
