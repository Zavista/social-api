package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret   string
	issuer   string
	audience string
	exp      time.Duration
}

func NewJWTAuthenticator(secret, issuer, audience string, exp time.Duration) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret:   secret,
		issuer:   issuer,
		audience: audience,
		exp:      exp,
	}
}

func (a *JWTAuthenticator) GenerateToken(userID int64) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": now.Add(a.exp).Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"iss": a.issuer,
		"aud": a.audience,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *JWTAuthenticator) ValidateToken(token string) (*Claims, error) {
	// TODO:
	return &Claims{}, nil
}
