package auth

type Claims struct {
	UserID int64
	// TODO: add roles in the claims later
}

type Authenticator interface {
	GenerateToken(userID int64) (string, error)
	ValidateToken(token string) (*Claims, error)
}
