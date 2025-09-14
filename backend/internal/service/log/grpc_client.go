package log

import (
	"context"
	"strconv"

	logspb "github.com/HirotakaUchishiba/calomeal_mvp/proto/logs/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClient implements the log.Service interface using gRPC
type GRPCClient struct {
	client logspb.LogServiceClient
	conn   *grpc.ClientConn
}

// NewGRPCClient creates a new gRPC client for the log service
func NewGRPCClient(address string) (*GRPCClient, error) {
	// gRPC接続の設定
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := logspb.NewLogServiceClient(conn)

	return &GRPCClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection
func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

// LogFood logs a food entry via gRPC
func (c *GRPCClient) LogFood(ctx context.Context, userID string, input LogFoodInput) (int64, error) {
	req := &logspb.LogFoodRequest{
		FoodId:   input.FoodName, // 簡易実装: FoodNameをFoodIdとして使用
		Quantity: input.Quantity,
		Unit:     input.Unit,
		Date:     input.Date,
	}

	resp, err := c.client.LogFood(ctx, req)
	if err != nil {
		return 0, err
	}

	// 文字列IDをint64に変換
	logID, err := strconv.ParseInt(resp.LogId, 10, 64)
	if err != nil {
		return 0, err
	}

	return logID, nil
}

// LogExercise logs an exercise entry via gRPC
func (c *GRPCClient) LogExercise(ctx context.Context, userID string, input LogExerciseInput) (int64, error) {
	req := &logspb.LogExerciseRequest{
		ExerciseName:    input.ExerciseName,
		DurationMinutes: int32(input.DurationMinutes),
		CaloriesBurned:  input.CaloriesBurned,
		Date:            input.Date,
	}

	resp, err := c.client.LogExercise(ctx, req)
	if err != nil {
		return 0, err
	}

	// 文字列IDをint64に変換
	logID, err := strconv.ParseInt(resp.LogId, 10, 64)
	if err != nil {
		return 0, err
	}

	return logID, nil
}

// LogWeight logs a weight entry via gRPC
func (c *GRPCClient) LogWeight(ctx context.Context, userID string, input LogWeightInput) (int64, error) {
	req := &logspb.LogWeightRequest{
		Weight: input.Weight,
		Date:   input.Date,
	}

	resp, err := c.client.LogWeight(ctx, req)
	if err != nil {
		return 0, err
	}

	// 文字列IDをint64に変換
	logID, err := strconv.ParseInt(resp.LogId, 10, 64)
	if err != nil {
		return 0, err
	}

	return logID, nil
}

// GetDailySummary calculates daily summary via gRPC
func (c *GRPCClient) GetDailySummary(ctx context.Context, userID, date string) (DailySummary, error) {
	// 各ログタイプを取得
	foodLogs, err := c.GetFoodLogs(ctx, userID, date)
	if err != nil {
		return DailySummary{}, err
	}

	exerciseLogs, err := c.GetExerciseLogs(ctx, userID, date)
	if err != nil {
		return DailySummary{}, err
	}

	// 合計値を計算
	var totalCalories, totalProtein, totalCarbohydrate, totalFat, totalCaloriesBurned float64

	for _, log := range foodLogs {
		totalCalories += log.Calories
		totalProtein += log.Protein
		totalCarbohydrate += log.Carbohydrate
		totalFat += log.Fat
	}

	for _, log := range exerciseLogs {
		totalCaloriesBurned += log.CaloriesBurned
	}

	return DailySummary{
		CaloriesIntake: totalCalories,
		CaloriesBurned: totalCaloriesBurned,
		Protein:        totalProtein,
		Carbohydrate:   totalCarbohydrate,
		Fat:            totalFat,
	}, nil
}

// GetFoodLogs gets food logs for a specific date via gRPC
func (c *GRPCClient) GetFoodLogs(ctx context.Context, userID, date string) ([]FoodLog, error) {
	req := &logspb.ListFoodLogsByDateRequest{
		Date: date,
	}

	resp, err := c.client.ListFoodLogsByDate(ctx, req)
	if err != nil {
		return nil, err
	}

	var logs []FoodLog
	for _, log := range resp.Logs {
		// 文字列IDをint64に変換
		logID, err := strconv.ParseInt(log.Id, 10, 64)
		if err != nil {
			continue // エラーの場合はスキップ
		}

		logs = append(logs, FoodLog{
			ID:           logID,
			FoodName:     log.FoodName,
			Quantity:     log.Quantity,
			Unit:         log.Unit,
			Calories:     log.Calories,
			Protein:      log.Protein,
			Carbohydrate: log.Carbohydrate,
			Fat:          log.Fat,
			LoggedAt:     log.LoggedAt,
		})
	}

	return logs, nil
}

// GetExerciseLogs gets exercise logs for a specific date via gRPC
func (c *GRPCClient) GetExerciseLogs(ctx context.Context, userID, date string) ([]ExerciseLog, error) {
	req := &logspb.ListExerciseLogsByDateRequest{
		Date: date,
	}

	resp, err := c.client.ListExerciseLogsByDate(ctx, req)
	if err != nil {
		return nil, err
	}

	var logs []ExerciseLog
	for _, log := range resp.Logs {
		// 文字列IDをint64に変換
		logID, err := strconv.ParseInt(log.Id, 10, 64)
		if err != nil {
			continue // エラーの場合はスキップ
		}

		logs = append(logs, ExerciseLog{
			ID:              logID,
			ExerciseName:    log.ExerciseName,
			DurationMinutes: int(log.DurationMinutes),
			CaloriesBurned:  log.CaloriesBurned,
			LoggedAt:        log.LoggedAt,
		})
	}

	return logs, nil
}

// GetWeightLogs gets weight logs for a specific date via gRPC
func (c *GRPCClient) GetWeightLogs(ctx context.Context, userID, date string) ([]WeightLog, error) {
	req := &logspb.ListWeightLogsByDateRequest{
		Date: date,
	}

	resp, err := c.client.ListWeightLogsByDate(ctx, req)
	if err != nil {
		return nil, err
	}

	var logs []WeightLog
	for _, log := range resp.Logs {
		// 文字列IDをint64に変換
		logID, err := strconv.ParseInt(log.Id, 10, 64)
		if err != nil {
			continue // エラーの場合はスキップ
		}

		logs = append(logs, WeightLog{
			ID:       logID,
			Weight:   log.Weight,
			LoggedAt: log.LoggedAt,
		})
	}

	return logs, nil
}