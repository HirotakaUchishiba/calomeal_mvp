// backend/internal/bff/middleware/auth.go
package middleware

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

// Auth is a directive middleware that checks for authentication.
func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	// TODO: この部分に実際の認証チェックロジックを実装します
	// （例：コンテキストからJWTトークンを取得し、検証する）
	fmt.Println("Auth middleware is called, but not implemented yet.")

	// 認証チェックが成功したとして、次のリゾルバを呼び出す
	return next(ctx)
}