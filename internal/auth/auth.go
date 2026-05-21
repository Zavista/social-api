package auth

type Claims struct {
	UserID int64
}

type Authenticator interface {
	GenerateToken(userID int64) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}
