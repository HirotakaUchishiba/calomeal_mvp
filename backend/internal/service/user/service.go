// backend/internal/service/user/service.go
package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// UserProfileInput represents the input for user profile creation
type UserProfileInput struct {
	Height        float64
	Weight        float64
	ActivityLevel string
}

// UserGoalInput represents the input for user goal creation
type UserGoalInput struct {
	TargetWeight float64
	TargetDate   string
}

// UserService defines the interface for user-related business logic.
type Service interface {
	CompleteOnboarding(ctx context.Context, userID string, profile UserProfileInput, goal UserGoalInput) error
}

// service is the concrete implementation of the Service interface.
type service struct {
	db *sql.DB
}

// NewService creates a new instance of the user service.
func NewService(db *sql.DB) Service {
	return &service{db: db}
}

// CompleteOnboarding handles the business logic for completing user onboarding.
func (s *service) CompleteOnboarding(ctx context.Context, userID string, profile UserProfileInput, goal UserGoalInput) error {
	// トランザクションを開始して、プロフィールと目標を同時に保存
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. ユーザープロフィールを保存
	profileQuery := `
		INSERT INTO user_profiles (user_id, height, weight, activity_level, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE SET
			height = EXCLUDED.height,
			weight = EXCLUDED.weight,
			activity_level = EXCLUDED.activity_level,
			updated_at = EXCLUDED.updated_at
	`
	now := time.Now()
	_, err = tx.ExecContext(ctx, profileQuery, userID, profile.Height, profile.Weight, profile.ActivityLevel, now, now)
	if err != nil {
		return fmt.Errorf("failed to save user profile: %w", err)
	}

	// 2. ユーザー目標を保存
	goalQuery := `
		INSERT INTO user_goals (user_id, target_weight, target_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET
			target_weight = EXCLUDED.target_weight,
			target_date = EXCLUDED.target_date,
			updated_at = EXCLUDED.updated_at
	`
	_, err = tx.ExecContext(ctx, goalQuery, userID, goal.TargetWeight, goal.TargetDate, now, now)
	if err != nil {
		return fmt.Errorf("failed to save user goal: %w", err)
	}

	// トランザクションをコミット
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("UserService: Successfully completed onboarding for user %s\n", userID)
	return nil
}
