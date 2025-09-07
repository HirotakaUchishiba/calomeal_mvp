// backend/internal/bff/middleware/auth.go
package middleware

import (
	"context"
	"fmt"
	"os"

	"github.com/99designs/gqlgen/graphql"
)

func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	// 開発環境では認証をスキップ
	if os.Getenv("ENVIRONMENT") == "development" || os.Getenv("ENVIRONMENT") == "" {
		fmt.Println("Auth directive: Development mode - skipping authentication")
		return next(ctx)
	}

	// TODO: 本番環境では実際のJWT認証を実装
	fmt.Println("Auth directive middleware called (production auth not implemented yet).")
	return next(ctx)
}
