package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type ctxKey int

const Key ctxKey = 1

type Auth struct {
	privateKey *rsa.PrivateKey
	publickey  *rsa.PublicKey
}

//go:generate mockgen -source=auth.go -destination=auth_mock.go -package=auth
type Authentication interface {
	GenerateAuthToken(claims jwt.RegisteredClaims) (string, error)
	ValidateToken(token string) (jwt.RegisteredClaims, error)
}

func (a *Auth) GenerateAuthToken(claims jwt.RegisteredClaims) (string, error) {
	// creates a new token with signing menthod and claims
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// signing our token with the private key
	token, err := tkn.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("error in signing the token : %w", err)
	}

	return token, nil
}

func (a *Auth) ValidateToken(token string) (jwt.RegisteredClaims, error) {
	// Parse the token with the registered claims.
	var c jwt.RegisteredClaims
	tkn, err := jwt.ParseWithClaims(token, &c, func(t *jwt.Token) (interface{}, error) {
		return a.publickey, nil
	})
	if err != nil {
		return jwt.RegisteredClaims{}, fmt.Errorf("error in parsing the token : %w", err)
	}

	// checking if the token is valid or not
	if !tkn.Valid {
		return jwt.RegisteredClaims{}, errors.New("token in not valid")
	}
	return c, nil
}

func NewAuth(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) (Authentication, error) {
	if privateKey == nil && publicKey == nil {
		return nil, errors.New("publickey and privatekey cannot be null")
	}
	return &Auth{
		privateKey: privateKey,
		publickey:  publicKey,
	}, nil
}
