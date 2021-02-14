package utils

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func NewJWTToken(secret string, signingAlg string, claims jwt.Claims) (string, error) {
	alg := jwt.GetSigningMethod(signingAlg)
	token := jwt.NewWithClaims(alg, claims)

	out, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("Failed to sign JWT token: %v", err)
	}

	return out, nil
}
