package auth

// Authenticator is the public interface application injects
type Authenticator interface {
	GenerateAccessToken(id int64, email string) (string, error)
	GenerateRefreshToken(id int64, email string) (string, error)
	VerifyToken(tokenStr, tokenType string) (*Claims, error)
}

type Claims struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
}
