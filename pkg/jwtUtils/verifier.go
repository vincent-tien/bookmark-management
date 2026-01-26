package jwtUtils

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

// JwtValidator is an interface for validating JWT tokens.
//
// This interface encapsulates the logic for validating JWT tokens using a
// public key. It provides a single method, ValidateToken, which takes a JWT
// token string as input and returns the JWT claims as a map and an error if the
// token is invalid.
type JwtValidator interface {
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type jwtValidator struct {
	publicKey *rsa.PublicKey
}

func NewJwtValidator(publicKeyPath string) (JwtValidator, error) {
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}

	return &jwtValidator{
		publicKey: publicKey,
	}, nil
}

var errInvalidToken = errors.New("invalid token")

func (j *jwtValidator) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errInvalidToken
	}

	return token.Claims.(jwt.MapClaims), nil
}
