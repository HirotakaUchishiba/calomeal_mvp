// backend/internal/service/log/service.go
package log

import (
	"context"
	"fmt"
)

// LogExerciseInputは運動記録の入力です
type LogExerciseInput struct {
	ExerciseName    string
	DurationMinutes int
	CaloriesBurned  float64
	Date            string
}

// TODO: LogFoodInputも同様に定義します

// Serviceは記録関連のビジネスロジックのインターフェースです
type Service interface {
	// TODO: LogFoodメソッドを定義します
	LogExercise(ctx context.Context, userID string, input LogExerciseInput) (int64, error)
}

type service struct {
	// TODO: ここにデータベース接続(リポジトリ)を保持します
}

// NewServiceは新しいlogサービスインスタンスを作成します
func NewService() Service {
	return &service{}
}

// LogExerciseは運動記録をデータベースに保存します
func (s *service) LogExercise(ctx context.Context, userID string, input LogExerciseInput) (int64, error) {
	// 【メンターズノート】
	// Lean MVPの原則に従い、消費カロリーの自動計算は行いません。
	// ユーザーが入力した値をそのままデータベースに保存するだけのシンプルな処理です。
	// これにより、METS表の統合といった「隠れた複雑性」を排除しています。

	// TODO: 実際のデータベースINSERTロジックを実装
	fmt.Printf("Saving exercise log for user %s: %+v\n", userID, input)

	// ダミーのIDを返します
	return 123, nil
}
