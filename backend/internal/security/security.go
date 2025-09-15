package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

// SecurityConfig holds security configuration
type SecurityConfig struct {
	JWTSecret          string
	JWTExpiration      time.Duration
	PasswordSaltLength int
	Argon2Memory       uint32
	Argon2Iterations   uint32
	Argon2Parallelism  uint8
	Argon2KeyLength    uint32
	RateLimitEnabled   bool
	RateLimitRequests  int
	RateLimitWindow    time.Duration
}

// DefaultSecurityConfig returns a default security configuration
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		JWTSecret:          "your-secret-key",
		JWTExpiration:      24 * time.Hour,
		PasswordSaltLength: 16,
		Argon2Memory:       64 * 1024, // 64 MB
		Argon2Iterations:   3,
		Argon2Parallelism:  2,
		Argon2KeyLength:    32,
		RateLimitEnabled:   true,
		RateLimitRequests:  100,
		RateLimitWindow:    time.Minute,
	}
}

// ProductionSecurityConfig returns a production security configuration
func ProductionSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		JWTSecret:          "", // Must be set from environment
		JWTExpiration:      1 * time.Hour,
		PasswordSaltLength: 32,
		Argon2Memory:       128 * 1024, // 128 MB
		Argon2Iterations:   4,
		Argon2Parallelism:  4,
		Argon2KeyLength:    32,
		RateLimitEnabled:   true,
		RateLimitRequests:  60,
		RateLimitWindow:    time.Minute,
	}
}

// JWTClaims represents JWT claims
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"exp"`
	jwt.RegisteredClaims
}

// SecurityManager handles security operations
type SecurityManager struct {
	config *SecurityConfig
}

// NewSecurityManager creates a new security manager
func NewSecurityManager(config *SecurityConfig) *SecurityManager {
	if config == nil {
		config = DefaultSecurityConfig()
	}
	return &SecurityManager{config: config}
}

// HashPassword hashes a password using Argon2
func (sm *SecurityManager) HashPassword(password string) (string, error) {
	// Generate random salt
	salt := make([]byte, sm.config.PasswordSaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash password
	hash := argon2.IDKey([]byte(password), salt, sm.config.Argon2Iterations, sm.config.Argon2Memory, sm.config.Argon2Parallelism, sm.config.Argon2KeyLength)

	// Encode salt and hash
	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)
	hashEncoded := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		sm.config.Argon2Memory,
		sm.config.Argon2Iterations,
		sm.config.Argon2Parallelism,
		saltEncoded,
		hashEncoded), nil
}

// VerifyPassword verifies a password against a hash
func (sm *SecurityManager) VerifyPassword(password, hash string) (bool, error) {
	// Parse hash
	parts := strings.Split(hash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, fmt.Errorf("invalid hash format")
	}

	// Decode salt and hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	// Hash the provided password
	actualHash := argon2.IDKey([]byte(password), salt, sm.config.Argon2Iterations, sm.config.Argon2Memory, sm.config.Argon2Parallelism, sm.config.Argon2KeyLength)

	// Compare hashes
	return subtle.ConstantTimeCompare(expectedHash, actualHash) == 1, nil
}

// GenerateJWT generates a JWT token
func (sm *SecurityManager) GenerateJWT(userID, email, role string) (string, error) {
	claims := JWTClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		ExpiresAt: time.Now().Add(sm.config.JWTExpiration).Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(sm.config.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "calomeal",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sm.config.JWTSecret))
}

// ValidateJWT validates a JWT token
func (sm *SecurityManager) ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(sm.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ExtractTokenFromRequest extracts JWT token from HTTP request
func (sm *SecurityManager) ExtractTokenFromRequest(r *http.Request) (string, error) {
	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
	}

	// Check cookie
	cookie, err := r.Cookie("auth_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	return "", fmt.Errorf("no token found")
}

// SecurityHeadersMiddleware adds security headers to HTTP responses
func SecurityHeadersMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Security headers
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
			
			// HSTS (only for HTTPS)
			if r.TLS != nil {
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
			}
			
			// CSP (Content Security Policy)
			csp := "default-src 'self'; " +
				"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
				"style-src 'self' 'unsafe-inline'; " +
				"img-src 'self' data: https:; " +
				"font-src 'self' data:; " +
				"connect-src 'self' https:; " +
				"frame-ancestors 'none';"
			w.Header().Set("Content-Security-Policy", csp)

			next.ServeHTTP(w, r)
		})
	}
}

// CORSMiddleware handles CORS
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			
			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware implements rate limiting
func RateLimitMiddleware(requests int, window time.Duration) func(http.Handler) http.Handler {
	// Simple in-memory rate limiter (in production, use Redis or similar)
	clients := make(map[string][]time.Time)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)
			now := time.Now()
			
			// Clean old entries
			if times, exists := clients[clientIP]; exists {
				var validTimes []time.Time
				for _, t := range times {
					if now.Sub(t) < window {
						validTimes = append(validTimes, t)
					}
				}
				clients[clientIP] = validTimes
			}

			// Check rate limit
			if len(clients[clientIP]) >= requests {
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requests))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", now.Add(window).Unix()))
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Add current request
			clients[clientIP] = append(clients[clientIP], now)

			// Set rate limit headers
			remaining := requests - len(clients[clientIP])
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requests))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", now.Add(window).Unix()))

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fallback to RemoteAddr
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}

// ValidateInput validates user input
func ValidateInput(input string, maxLength int) error {
	if len(input) > maxLength {
		return fmt.Errorf("input too long: max %d characters", maxLength)
	}
	
	// Check for potentially dangerous characters
	dangerousChars := []string{"<", ">", "\"", "'", "&", "script", "javascript:", "data:"}
	for _, char := range dangerousChars {
		if strings.Contains(strings.ToLower(input), char) {
			return fmt.Errorf("input contains potentially dangerous content")
		}
	}
	
	return nil
}

// SanitizeInput sanitizes user input
func SanitizeInput(input string) string {
	// Remove HTML tags
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#x27;")
	input = strings.ReplaceAll(input, "&", "&amp;")
	
	return input
}
