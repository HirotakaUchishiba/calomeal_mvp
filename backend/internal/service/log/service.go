// backend/internal/service/log/service.go
package log

import (
	"context"
	"database/sql"
)

// LogExerciseInputは運動記録の入力です
type LogExerciseInput struct {
	ExerciseName    string
	DurationMinutes int
	CaloriesBurned  float64
	Date            string
}

type LogFoodInput struct {
	FoodName     string
	Quantity     float64
	Unit         string
	Calories     float64
	Protein      float64
	Carbohydrate float64
	Fat          float64
	Date         string
  }

  type DailySummary struct {
	CaloriesIntake  float64
	CaloriesBurned  float64
	Protein         float64
	Carbohydrate    float64
	Fat             float64
  }

// FoodLogは食事記録を表します
type FoodLog struct {
	ID           int64
	FoodName     string
	Quantity     float64
	Unit         string
	Calories     float64
	Protein      float64
	Carbohydrate float64
	Fat          float64
	LoggedAt     string
}

// ExerciseLogは運動記録を表します
type ExerciseLog struct {
	ID              int64
	ExerciseName    string
	DurationMinutes int
	CaloriesBurned  float64
	LoggedAt        string
}

// Serviceは記録関連のビジネスロジックのインターフェースです
type Service interface {
	LogExercise(ctx context.Context, userID string, input LogExerciseInput) (int64, error)
	LogFood(ctx context.Context, userID string, input LogFoodInput) (int64, error)
	GetDailySummary(ctx context.Context, userID, date string) (DailySummary, error)
	GetFoodLogs(ctx context.Context, userID, date string) ([]FoodLog, error)
	GetExerciseLogs(ctx context.Context, userID, date string) ([]ExerciseLog, error)
}

type service struct {
	db *sql.DB
}

// NewServiceは新しいlogサービスインスタンスを作成します
func NewService(db *sql.DB) Service {
	return &service{db: db}
}

func (s *service) LogFood(ctx context.Context, userID string, in LogFoodInput) (int64, error) {
	const q = `
    INSERT INTO food_logs (user_id, food_name, quantity, unit, calories, protein, carbohydrate, fat, logged_at)
    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
    RETURNING id
	`
	var id int64
	if err := s.db.QueryRowContext(ctx, q, userID, in.FoodName, in.Quantity, in.Unit, in.Calories, in.Protein, in.Carbohydrate, in.Fat, in.Date).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
	}

// LogExerciseは運動記録をデータベースに保存します
func (s *service) LogExercise(ctx context.Context, userID string, in LogExerciseInput) (int64, error) {
	const q = `
		INSERT INTO exercise_logs (
			user_id, exercise_name, duration_minutes, calories_burned, logged_at
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id int64
	if err := s.db.QueryRowContext(ctx, q,
		userID,
		in.ExerciseName,
		in.DurationMinutes,
		in.CaloriesBurned,
		in.Date,
	).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *service) GetDailySummary(ctx context.Context, userID, date string) (DailySummary, error) {
	const q = `
		WITH f AS (
			SELECT
				COALESCE(SUM(calories), 0)     AS cal_in,
				COALESCE(SUM(protein), 0)      AS p_in,
				COALESCE(SUM(carbohydrate), 0) AS c_in,
				COALESCE(SUM(fat), 0)          AS f_in
			FROM food_logs
			WHERE user_id = $1
			  AND DATE_TRUNC('day', logged_at) = DATE_TRUNC('day', $2::timestamptz)
		), e AS (
			SELECT
				COALESCE(SUM(calories_burned), 0) AS cal_out
			FROM exercise_logs
			WHERE user_id = $1
			  AND DATE_TRUNC('day', logged_at) = DATE_TRUNC('day', $2::timestamptz)
		)
		SELECT f.cal_in, e.cal_out, f.p_in, f.c_in, f.f_in
		FROM f CROSS JOIN e
	`
	var ds DailySummary
	if err := s.db.
		QueryRowContext(ctx, q, userID, date).
		Scan(&ds.CaloriesIntake, &ds.CaloriesBurned, &ds.Protein, &ds.Carbohydrate, &ds.Fat); err != nil {
		return DailySummary{}, err
	}
	return ds, nil
}

// GetFoodLogsは指定された日付の食事記録を取得します
func (s *service) GetFoodLogs(ctx context.Context, userID, date string) ([]FoodLog, error) {
	const q = `
		SELECT id, food_name, quantity, unit, calories, protein, carbohydrate, fat, logged_at
		FROM food_logs
		WHERE user_id = $1
		  AND DATE_TRUNC('day', logged_at) = DATE_TRUNC('day', $2::timestamptz)
		ORDER BY logged_at DESC
	`
	rows, err := s.db.QueryContext(ctx, q, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []FoodLog
	for rows.Next() {
		var log FoodLog
		err := rows.Scan(
			&log.ID,
			&log.FoodName,
			&log.Quantity,
			&log.Unit,
			&log.Calories,
			&log.Protein,
			&log.Carbohydrate,
			&log.Fat,
			&log.LoggedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetExerciseLogsは指定された日付の運動記録を取得します
func (s *service) GetExerciseLogs(ctx context.Context, userID, date string) ([]ExerciseLog, error) {
	const q = `
		SELECT id, exercise_name, duration_minutes, calories_burned, logged_at
		FROM exercise_logs
		WHERE user_id = $1
		  AND DATE_TRUNC('day', logged_at) = DATE_TRUNC('day', $2::timestamptz)
		ORDER BY logged_at DESC
	`
	rows, err := s.db.QueryContext(ctx, q, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []ExerciseLog
	for rows.Next() {
		var log ExerciseLog
		err := rows.Scan(
			&log.ID,
			&log.ExerciseName,
			&log.DurationMinutes,
			&log.CaloriesBurned,
			&log.LoggedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}