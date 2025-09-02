// backend/internal/service/user/service.go
package user

import (
	"context"
	"fmt"
)

// UserService defines the interface for user-related business logic.
type Service interface {
	// TODO: ここに引数を追加します (例: userID string, profileInput, goalInput)
	CompleteOnboarding(ctx context.Context) error
}

// service is the concrete implementation of the Service interface.
type service struct {
	// TODO: ここにデータベース接続(リポジトリ)などの依存関係を追加します
}

// NewService creates a new instance of the user service.
func NewService() Service {
	return &service{}
}

// CompleteOnboarding handles the business logic for completing user onboarding.
func (s *service) CompleteOnboarding(ctx context.Context) error {
	// 現時点では、実際のデータベース保存処理の代わりにコンソールにメッセージを表示します。
	// 今後のステップで、この部分に実際の保存ロジックを実装していきます。
	fmt.Println("UserService: CompleteOnboarding called. Saving user profile and goals to the database...")

	return nil
}
