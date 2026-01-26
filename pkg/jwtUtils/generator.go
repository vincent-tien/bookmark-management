package jwtUtils

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	tokenLast = 24 * time.Hour
)

// JwtGenerator is an interface for generating JWT tokens.
//
// This interface encapsulates the logic for generating JWT tokens using a
// private key. It provides a single method, GenerateToken, which takes JWT
// content as input and returns a signed JWT token as a string.
//
// The private key is loaded from the file system using the privateKeyPath
// parameter passed to the NewJwtGenerator function.
type JwtGenerator interface {
	GenerateToken(jwtContent jwt.MapClaims) (string, error)
	GenerateContent(sub string) jwt.MapClaims
}

type jwtGenerator struct {
	privateKey *rsa.PrivateKey
}

func NewJwtGenerator(privateKeyPath string) (JwtGenerator, error) {
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, err
	}

	return &jwtGenerator{
		privateKey: privateKey,
	}, nil
}

func (j *jwtGenerator) GenerateContent(sub string) jwt.MapClaims {
	return jwt.MapClaims{
		"sub": sub,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(tokenLast).Unix(),
	}
}

func (j *jwtGenerator) GenerateToken(jwtContent jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtContent)
	return token.SignedString(j.privateKey)
}
