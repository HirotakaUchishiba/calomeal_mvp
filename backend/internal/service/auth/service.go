package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// AuthService defines the interface for authentication-related business logic
type Service interface {
	GenerateTokenPair(ctx context.Context, userID, email string) (*TokenPair, error)
	ValidateToken(ctx context.Context, tokenString string) (*JWTClaims, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	RevokeToken(ctx context.Context, tokenString string) error
}

type service struct {
	accessTokenSecret  string
	refreshTokenSecret string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

// NewService creates a new authentication service
func NewService() Service {
	return &service{
		accessTokenSecret:  getEnvOrDefault("JWT_ACCESS_SECRET", "your-access-secret-key"),
		refreshTokenSecret: getEnvOrDefault("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
		accessTokenExpiry:  time.Hour * 1,  // 1 hour
		refreshTokenExpiry: time.Hour * 24 * 7, // 7 days
	}
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GenerateTokenPair generates both access and refresh tokens
func (s *service) GenerateTokenPair(ctx context.Context, userID, email string) (*TokenPair, error) {
	now := time.Now()
	
	// Generate access token
	accessClaims := &JWTClaims{
		UserID:    userID,
		Email:     email,
		ExpiresAt: now.Add(s.accessTokenExpiry).Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenExpiry)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "calomeal-mvp",
			Subject:   userID,
		},
	}
	
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.accessTokenSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	
	// Generate refresh token
	refreshClaims := &JWTClaims{
		UserID:    userID,
		Email:     email,
		ExpiresAt: now.Add(s.refreshTokenExpiry).Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenExpiry)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "calomeal-mvp",
			Subject:   userID,
		},
	}
	
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.refreshTokenSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	
	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(s.accessTokenExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken validates and parses a JWT token
func (s *service) ValidateToken(ctx context.Context, tokenString string) (*JWTClaims, error) {
	// Try to parse as access token first
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.accessTokenSecret), nil
	})
	
	if err != nil {
		// If access token validation fails, try refresh token
		token, err = jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.refreshTokenSecret), nil
		})
		
		if err != nil {
			return nil, fmt.Errorf("invalid token: %w", err)
		}
	}
	
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	
	// Check if token is expired
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, errors.New("token has expired")
	}
	
	return claims, nil
}

// RefreshToken generates a new token pair using a refresh token
func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.ValidateToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	
	// Generate new token pair
	return s.GenerateTokenPair(ctx, claims.UserID, claims.Email)
}

// RevokeToken revokes a token (in a real implementation, you would add it to a blacklist)
func (s *service) RevokeToken(ctx context.Context, tokenString string) error {
	// In a production environment, you would:
	// 1. Parse the token to get its expiration time
	// 2. Add it to a Redis blacklist or database
	// 3. Check the blacklist in ValidateToken
	
	// For now, we'll just validate that the token is valid
	_, err := s.ValidateToken(ctx, tokenString)
	if err != nil {
		return fmt.Errorf("invalid token to revoke: %w", err)
	}
	
	fmt.Printf("AuthService: Token revoked for token ending in ...%s\n", tokenString[len(tokenString)-10:])
	return nil
}
