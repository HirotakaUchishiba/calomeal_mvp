// backend/internal/bff/resolvers/resolver.go
package resolvers

import (
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/analytics"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/fooddata"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/log"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserService      user.Service
	FoodDataService  fooddata.Service // 追加
	LogService       log.Service      // 追加
	AnalyticsService *analytics.AnalyticsClient // 追加
}
