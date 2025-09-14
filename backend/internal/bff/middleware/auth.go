// backend/internal/bff/middleware/auth.go
package middleware

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/auth"
)

// ContextKey is used for context keys
type ContextKey string

const (
	UserIDKey  ContextKey = "user_id"
	EmailKey   ContextKey = "email"
	ClaimsKey  ContextKey = "claims"
)

// AuthService is injected into the middleware
var authService auth.Service

// InitAuthMiddleware initializes the auth middleware with the auth service
func InitAuthMiddleware(service auth.Service) {
	authService = service
}

func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	// 開発環境では認証をスキップ（環境変数で制御）
	if os.Getenv("ENVIRONMENT") == "development" || os.Getenv("ENVIRONMENT") == "" {
		// 開発環境でも認証ヘッダーがある場合は認証を実行
		reqCtx := graphql.GetOperationContext(ctx)
		if reqCtx != nil {
			authHeader := reqCtx.Headers.Get("Authorization")
			if authHeader != "" {
				// 認証ヘッダーがある場合は認証を実行
				fmt.Println("Auth directive: Development mode with auth header - performing authentication")
			} else {
				// 認証ヘッダーがない場合はダミーユーザーを設定
				fmt.Println("Auth directive: Development mode - using dummy user")
				ctx = context.WithValue(ctx, UserIDKey, "dev-user-123")
				ctx = context.WithValue(ctx, EmailKey, "dev@example.com")
				return next(ctx)
			}
		} else {
			// コンテキストがない場合はダミーユーザーを設定
			fmt.Println("Auth directive: Development mode - using dummy user")
			ctx = context.WithValue(ctx, UserIDKey, "dev-user-123")
			ctx = context.WithValue(ctx, EmailKey, "dev@example.com")
			return next(ctx)
		}
	}
	
	// 本番環境ではJWT認証を実行
	if authService == nil {
		return nil, fmt.Errorf("auth service not initialized")
	}
	
	// GraphQL contextからHTTP requestを取得
	reqCtx := graphql.GetOperationContext(ctx)
	if reqCtx == nil {
		return nil, fmt.Errorf("no operation context found")
	}
	
	// HTTP requestからAuthorization headerを取得
	authHeader := reqCtx.Headers.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is required")
	}
	
	// Bearer tokenを抽出
	tokenString := extractBearerToken(authHeader)
	if tokenString == "" {
		return nil, fmt.Errorf("invalid authorization header format")
	}
	
	// JWT tokenを検証
	claims, err := authService.ValidateToken(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	
	// ユーザー情報をcontextに追加
	ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
	ctx = context.WithValue(ctx, EmailKey, claims.Email)
	ctx = context.WithValue(ctx, ClaimsKey, claims)
	
	fmt.Printf("Auth directive: User authenticated - ID: %s, Email: %s\n", claims.UserID, claims.Email)
	return next(ctx)
}

// extractBearerToken extracts the Bearer token from the Authorization header
func extractBearerToken(authHeader string) string {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return ""
	}
	return strings.TrimSpace(authHeader[len(bearerPrefix):])
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetEmailFromContext extracts email from context
func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailKey).(string)
	return email, ok
}

// GetClaimsFromContext extracts JWT claims from context
func GetClaimsFromContext(ctx context.Context) (*auth.JWTClaims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*auth.JWTClaims)
	return claims, ok
}
