package auth

import (
	"fmt"
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

func (a *JWTAuthenticator) ValidateToken(tokenString string) (*Claims, error) {
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(a.secret), nil
	},
		jwt.WithAudience(a.audience),
		jwt.WithIssuer(a.issuer),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)

	if err != nil {
		return nil, err
	}

	// Extract the JWT claims from the parsed token.
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract the user ID from the subject claim.
	sub, ok := claims["sub"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid subject")
	}

	return &Claims{
		UserID: int64(sub),
	}, nil
}
