package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT token generation and validation.
type JWTService struct {
	secretKey     string
	issuer        string
	expireMinutes int
}

// NewJWTService creates a new instance of JWTService.
func NewJWTService(secretKey, issuer string, expireMinutes int) *JWTService {
	return &JWTService{
		secretKey:     secretKey,
		issuer:        issuer,
		expireMinutes: expireMinutes,
	}
}

// GenerateToken generates a new JWT token for a given user ID.
func (s *JWTService) GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iss": s.issuer,
		"exp": time.Now().Add(time.Minute * time.Duration(s.expireMinutes)).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// Authenticate validates the given JWT token and implements the domain.Authenticator interface.
func (s *JWTService) Authenticate(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return false, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}
