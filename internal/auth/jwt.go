package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims represents the JWT payload carried in every authenticated request.
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTService handles token generation and validation.
type JWTService struct {
	secretKey []byte
	issuer    string
	ttl       time.Duration
}

// NewJWTService creates a JWTService. ttl controls how long tokens stay valid.
func NewJWTService(secret string, ttl time.Duration) *JWTService {
	return &JWTService{
		secretKey: []byte(secret),
		issuer:    "todo-backend-go",
		ttl:       ttl,
	}
}

// GenerateToken creates a signed JWT for the given user.
func (s *JWTService) GenerateToken(userID int64, email string) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken parses a token string and returns the embedded claims.
func (s *JWTService) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
