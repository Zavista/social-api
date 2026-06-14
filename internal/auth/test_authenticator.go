package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TestAuthenticator struct{}

const (
	testSecret   = "test-secret"
	testAudience = "test-aud"
	testIssuer   = "test-aud"
)

func (a *TestAuthenticator) GenerateToken(userID int64) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": now.Add(time.Hour).Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"iss": testIssuer,
		"aud": testAudience,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(testSecret))
}

func (a *TestAuthenticator) ValidateToken(tokenString string) (*Claims, error) {
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(testSecret), nil
	},
		jwt.WithAudience(testAudience),
		jwt.WithIssuer(testIssuer),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid subject")
	}

	return &Claims{UserID: int64(sub)}, nil
}
