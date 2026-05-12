package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenDuration    = 15 * time.Minute
	RefreshTokenDuration = 2 * time.Hour
	RememberMeDuration = 30 * 24 * time.Hour // 30 days
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID  string    `json:"user_id"`
	Email   string    `json:"email"`
	Role    string    `json:"role"`
	StoreID *string   `json:"store_id,omitempty"`
	Type    TokenType `json:"type"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret []byte
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{secret: []byte(secret)}
}

// GenerateAccessToken creates a short-lived access token (15 minutes).
func (m *JWTManager) GenerateAccessToken(userID, email, role string, storeID *string) (string, error) {
	return m.generateToken(userID, email, role, storeID, AccessToken, AccessTokenDuration)
}

// GenerateRefreshToken creates a longer-lived refresh token (2 hours).
func (m *JWTManager) GenerateRefreshToken(userID, email, role string, storeID *string) (string, error) {
	return m.generateToken(userID, email, role, storeID, RefreshToken, RefreshTokenDuration)
}

// GenerateRememberMeToken creates a long-lived remember token (30 days).
func (m *JWTManager) GenerateRememberMeToken(userID, email, role string, storeID *string) (string, error) {
	return m.generateToken(userID, email, role, storeID, RefreshToken, RememberMeDuration)
}

func (m *JWTManager) generateToken(userID, email, role string, storeID *string, tokenType TokenType, duration time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:  userID,
		Email:   email,
		Role:    role,
		StoreID: storeID,
		Type:    tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ValidateToken parses and validates a JWT string. Returns the claims if valid.
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

