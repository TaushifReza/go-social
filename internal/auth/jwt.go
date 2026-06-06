package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

// jwtClaims is an unexported helper to work directly with the jwt package mechanisms
type jwtClaims struct {
	Claims
	jwt.RegisteredClaims
}

// jwtAuthenticator is unexported. It can only be constructed via NewJWTAuthenticator
type jwtAuthenticator struct {
	jwtSecret string
	aud       string
	iss       string
}

// NewJWTAuthenticator instantiates the unexported struct but returns the exported Authenticator interface.
func NewJWTAuthenticator(jwtSecret, aud, iss string) Authenticator {
	return &jwtAuthenticator{
		jwtSecret: jwtSecret,
		aud:       aud,
		iss:       iss,
	}
}

func (a *jwtAuthenticator) GenerateAccessToken(id int64, email string) (string, error) {
	tokenUUID := uuid.New()
	claims := &jwtClaims{
		Claims: Claims{
			ID:        id,
			Email:     email,
			TokenType: AccessTokenType,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", id),
			ID:        tokenUUID.String(),
			Issuer:    a.iss,
			Audience:  []string{a.aud},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecret))
}

func (a *jwtAuthenticator) GenerateRefreshToken(id int64, email string) (string, error) {
	tokenUUID := uuid.New()
	claims := &jwtClaims{
		Claims: Claims{
			ID:        id,
			Email:     email,
			TokenType: RefreshTokenType,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", id),
			ID:        tokenUUID.String(),
			Issuer:    a.iss,
			Audience:  []string{a.aud},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // 30 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecret))
}

func (a *jwtAuthenticator) VerifyToken(tokenStr, tokenType string) (*Claims, error) {
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Validate the signing method matches what we expect
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.jwtSecret), nil
	},
		jwt.WithExpirationRequired(),                                // Enforces that 'exp' must exist and be valid
		jwt.WithAudience(a.aud),                                     // Validates the 'aud' matches
		jwt.WithIssuer(a.iss),                                       // Validates the 'iss' matches
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}), // Restricts algorithms strictly to HS256
	)

	if err != nil {
		fmt.Println("TOKEN VALIDATION ERROR:", err)
		return nil, errors.New("unauthorized: invalid or expired token")
	}

	wrapper, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("unauthorized: claims are invalid")
	}

	if wrapper.TokenType != tokenType {
		return nil, errors.New("unauthorized: token type mismatch")
	}

	// Safely map down to app's core data structure and return it
	return &wrapper.Claims, nil
}
