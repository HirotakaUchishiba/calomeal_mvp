// backend/internal/bff/middleware/auth.go
package middleware

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	// This is a placeholder for your actual authentication logic.
	fmt.Println("Auth directive middleware called (not implemented yet).")
	return next(ctx)
}
