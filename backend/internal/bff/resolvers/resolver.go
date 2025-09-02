// backend/internal/bff/resolvers/resolver.go
package resolvers

import "github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/user" // この行を追加

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserService user.Service // この行を追加
}
