package server

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	analyticspb "github.com/HirotakaUchishiba/calomeal_mvp/services/analytics/internal/proto"
	"github.com/HirotakaUchishiba/calomeal_mvp/services/analytics/internal/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AnalyticsService implements the analytics gRPC service
type AnalyticsService struct {
	analyticspb.UnimplementedAnalyticsServiceServer
	db *sql.DB
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(db *sql.DB) *AnalyticsService {
	return &AnalyticsService{
		db: db,
	}
}

// GetDailyNutritionSummary returns daily nutrition summary
func (s *AnalyticsService) GetDailyNutritionSummary(ctx context.Context, req *analyticspb.GetDailyNutritionSummaryRequest) (*analyticspb.GetDailyNutritionSummaryResponse, error) {
	// Validate metadata and user ID
	if err := middleware.ValidateUserID(ctx, req.UserId); err != nil {
		return nil, err
	}

	// Log request with metadata
	middleware.LogWithMetadata(ctx, "INFO", "GetDailyNutritionSummary request", "date", req.Date)

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		middleware.LogWithMetadata(ctx, "ERROR", "Invalid date format", "date", req.Date, "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid date format: %v", err)
	}

	// Get food logs for the day
	foodLogs, err := s.getFoodLogsForDate(ctx, req.UserId, date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get food logs: %v", err)
	}

	// Get exercise logs for the day
	exerciseLogs, err := s.getExerciseLogsForDate(ctx, req.UserId, date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get exercise logs: %v", err)
	}

	// Calculate nutrition summary
	summary := s.calculateDailyNutritionSummary(date, foodLogs, exerciseLogs)

	return &analyticspb.GetDailyNutritionSummaryResponse{
		Summary: summary,
	}, nil
}

// GetWeeklyNutritionTrends returns weekly nutrition trends
func (s *AnalyticsService) GetWeeklyNutritionTrends(ctx context.Context, req *analyticspb.GetWeeklyNutritionTrendsRequest) (*analyticspb.GetWeeklyNutritionTrendsResponse, error) {
	// Validate metadata and user ID
	if err := middleware.ValidateUserID(ctx, req.UserId); err != nil {
		return nil, err
	}

	// Log request with metadata
	middleware.LogWithMetadata(ctx, "INFO", "GetWeeklyNutritionTrends request", "startDate", req.StartDate, "endDate", req.EndDate)

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		middleware.LogWithMetadata(ctx, "ERROR", "Invalid start date format", "startDate", req.StartDate, "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid start date format: %v", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid end date format: %v", err)
	}

	// Get daily summaries for the week
	var dailySummaries []*analyticspb.DailyNutritionSummary
	currentDate := startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		foodLogs, _ := s.getFoodLogsForDate(ctx, req.UserId, currentDate)
		exerciseLogs, _ := s.getExerciseLogsForDate(ctx, req.UserId, currentDate)

		summary := s.calculateDailyNutritionSummary(currentDate, foodLogs, exerciseLogs)
		dailySummaries = append(dailySummaries, summary)

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Calculate trends
	trends := s.calculateNutritionTrends(dailySummaries)

	return &analyticspb.GetWeeklyNutritionTrendsResponse{
		DailySummaries: dailySummaries,
		Trends:         trends,
	}, nil
}

// GetMonthlyNutritionInsights returns monthly nutrition insights
func (s *AnalyticsService) GetMonthlyNutritionInsights(ctx context.Context, req *analyticspb.GetMonthlyNutritionInsightsRequest) (*analyticspb.GetMonthlyNutritionInsightsResponse, error) {
	// Parse year and month
	year, err := strconv.Atoi(req.Year)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid year: %v", err)
	}

	month, err := strconv.Atoi(req.Month)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid month: %v", err)
	}

	// Get start and end dates for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	// Get all food logs for the month
	foodLogs, err := s.getFoodLogsForDateRange(ctx, req.UserId, startDate, endDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get food logs: %v", err)
	}

	// Calculate monthly insights
	insights := s.calculateMonthlyInsights(year, month, foodLogs)

	return &analyticspb.GetMonthlyNutritionInsightsResponse{
		Insights: insights,
	}, nil
}

// GetWeightProgressAnalysis returns weight progress analysis
func (s *AnalyticsService) GetWeightProgressAnalysis(ctx context.Context, req *analyticspb.GetWeightProgressAnalysisRequest) (*analyticspb.GetWeightProgressAnalysisResponse, error) {
	// Validate metadata and user ID
	if err := middleware.ValidateUserID(ctx, req.UserId); err != nil {
		return nil, err
	}

	// Log request with metadata
	middleware.LogWithMetadata(ctx, "INFO", "GetWeightProgressAnalysis request", "startDate", req.StartDate, "endDate", req.EndDate)

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		middleware.LogWithMetadata(ctx, "ERROR", "Invalid start date format", "startDate", req.StartDate, "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid start date format: %v", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid end date format: %v", err)
	}

	// Get weight logs for the period
	weightLogs, err := s.getWeightLogsForDateRange(ctx, req.UserId, startDate, endDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get weight logs: %v", err)
	}

	// Calculate weight progress analysis
	analysis := s.calculateWeightProgressAnalysis(startDate, endDate, weightLogs)

	return &analyticspb.GetWeightProgressAnalysisResponse{
		Analysis: analysis,
	}, nil
}

// GetCalorieBalanceAnalysis returns calorie balance analysis
func (s *AnalyticsService) GetCalorieBalanceAnalysis(ctx context.Context, req *analyticspb.GetCalorieBalanceAnalysisRequest) (*analyticspb.GetCalorieBalanceAnalysisResponse, error) {
	// Validate metadata and user ID
	if err := middleware.ValidateUserID(ctx, req.UserId); err != nil {
		return nil, err
	}

	// Log request with metadata
	middleware.LogWithMetadata(ctx, "INFO", "GetCalorieBalanceAnalysis request", "startDate", req.StartDate, "endDate", req.EndDate)

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		middleware.LogWithMetadata(ctx, "ERROR", "Invalid start date format", "startDate", req.StartDate, "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid start date format: %v", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid end date format: %v", err)
	}

	// Get food and exercise logs for the period
	foodLogs, err := s.getFoodLogsForDateRange(ctx, req.UserId, startDate, endDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get food logs: %v", err)
	}

	exerciseLogs, err := s.getExerciseLogsForDateRange(ctx, req.UserId, startDate, endDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get exercise logs: %v", err)
	}

	// Calculate calorie balance analysis
	analysis := s.calculateCalorieBalanceAnalysis(startDate, endDate, foodLogs, exerciseLogs)

	return &analyticspb.GetCalorieBalanceAnalysisResponse{
		Analysis: analysis,
	}, nil
}

// Helper methods for data retrieval and calculation

type FoodLog struct {
	ID       int
	UserID   string
	FoodID   string
	Quantity float64
	MealType string
	Date     time.Time
	Calories int
	Protein  int
	Carb     int
	Fat      int
	Fiber    int
	Sugar    int
	Sodium   int
}

type ExerciseLog struct {
	ID             int
	UserID         string
	ExerciseType   string
	Duration       int
	CaloriesBurned int
	Date           time.Time
}

type WeightLog struct {
	ID     int
	UserID string
	Weight float64
	Date   time.Time
}

func (s *AnalyticsService) getFoodLogsForDate(ctx context.Context, userID string, date time.Time) ([]FoodLog, error) {
	query := `
		SELECT id, user_id, food_id, quantity, meal_type, date, calories, protein, carbohydrate, fat, fiber, sugar, sodium
		FROM food_logs 
		WHERE user_id = $1 AND DATE(date) = $2
		ORDER BY date
	`

	rows, err := s.db.QueryContext(ctx, query, userID, date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []FoodLog
	for rows.Next() {
		var log FoodLog
		var dateStr string
		err := rows.Scan(&log.ID, &log.UserID, &log.FoodID, &log.Quantity, &log.MealType, &dateStr,
			&log.Calories, &log.Protein, &log.Carb, &log.Fat, &log.Fiber, &log.Sugar, &log.Sodium)
		if err != nil {
			return nil, err
		}
		log.Date, _ = time.Parse("2006-01-02", dateStr)
		logs = append(logs, log)
	}

	return logs, nil
}

func (s *AnalyticsService) getFoodLogsForDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]FoodLog, error) {
	query := `
		SELECT id, user_id, food_id, quantity, meal_type, date, calories, protein, carbohydrate, fat, fiber, sugar, sodium
		FROM food_logs 
		WHERE user_id = $1 AND DATE(date) BETWEEN $2 AND $3
		ORDER BY date
	`

	rows, err := s.db.QueryContext(ctx, query, userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []FoodLog
	for rows.Next() {
		var log FoodLog
		var dateStr string
		err := rows.Scan(&log.ID, &log.UserID, &log.FoodID, &log.Quantity, &log.MealType, &dateStr,
			&log.Calories, &log.Protein, &log.Carb, &log.Fat, &log.Fiber, &log.Sugar, &log.Sodium)
		if err != nil {
			return nil, err
		}
		log.Date, _ = time.Parse("2006-01-02", dateStr)
		logs = append(logs, log)
	}

	return logs, nil
}

func (s *AnalyticsService) getExerciseLogsForDate(ctx context.Context, userID string, date time.Time) ([]ExerciseLog, error) {
	query := `
		SELECT id, user_id, exercise_type, duration, calories_burned, date
		FROM exercise_logs 
		WHERE user_id = $1 AND DATE(date) = $2
		ORDER BY date
	`

	rows, err := s.db.QueryContext(ctx, query, userID, date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []ExerciseLog
	for rows.Next() {
		var log ExerciseLog
		var dateStr string
		err := rows.Scan(&log.ID, &log.UserID, &log.ExerciseType, &log.Duration, &log.CaloriesBurned, &dateStr)
		if err != nil {
			return nil, err
		}
		log.Date, _ = time.Parse("2006-01-02", dateStr)
		logs = append(logs, log)
	}

	return logs, nil
}

func (s *AnalyticsService) getExerciseLogsForDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]ExerciseLog, error) {
	query := `
		SELECT id, user_id, exercise_type, duration, calories_burned, date
		FROM exercise_logs 
		WHERE user_id = $1 AND DATE(date) BETWEEN $2 AND $3
		ORDER BY date
	`

	rows, err := s.db.QueryContext(ctx, query, userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []ExerciseLog
	for rows.Next() {
		var log ExerciseLog
		var dateStr string
		err := rows.Scan(&log.ID, &log.UserID, &log.ExerciseType, &log.Duration, &log.CaloriesBurned, &dateStr)
		if err != nil {
			return nil, err
		}
		log.Date, _ = time.Parse("2006-01-02", dateStr)
		logs = append(logs, log)
	}

	return logs, nil
}

func (s *AnalyticsService) getWeightLogsForDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]WeightLog, error) {
	query := `
		SELECT id, user_id, weight, date
		FROM weight_logs 
		WHERE user_id = $1 AND DATE(date) BETWEEN $2 AND $3
		ORDER BY date
	`

	rows, err := s.db.QueryContext(ctx, query, userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []WeightLog
	for rows.Next() {
		var log WeightLog
		var dateStr string
		err := rows.Scan(&log.ID, &log.UserID, &log.Weight, &dateStr)
		if err != nil {
			return nil, err
		}
		log.Date, _ = time.Parse("2006-01-02", dateStr)
		logs = append(logs, log)
	}

	return logs, nil
}

func (s *AnalyticsService) calculateDailyNutritionSummary(date time.Time, foodLogs []FoodLog, exerciseLogs []ExerciseLog) *analyticspb.DailyNutritionSummary {
	// Calculate totals from food logs
	var totalCalories, totalProtein, totalCarb, totalFat, totalFiber, totalSugar, totalSodium int
	mealSummaries := make(map[string]*analyticspb.MealSummary)

	for _, log := range foodLogs {
		totalCalories += log.Calories
		totalProtein += log.Protein
		totalCarb += log.Carb
		totalFat += log.Fat
		totalFiber += log.Fiber
		totalSugar += log.Sugar
		totalSodium += log.Sodium

		// Group by meal type
		if mealSummary, exists := mealSummaries[log.MealType]; exists {
			mealSummary.Calories += int32(log.Calories)
			mealSummary.Protein += int32(log.Protein)
			mealSummary.Carbohydrate += int32(log.Carb)
			mealSummary.Fat += int32(log.Fat)
		} else {
			mealSummaries[log.MealType] = &analyticspb.MealSummary{
				MealType:     log.MealType,
				Calories:     int32(log.Calories),
				Protein:      int32(log.Protein),
				Carbohydrate: int32(log.Carb),
				Fat:          int32(log.Fat),
			}
		}
	}

	// Calculate calories burned from exercise logs
	var totalCaloriesBurned int
	for _, log := range exerciseLogs {
		totalCaloriesBurned += log.CaloriesBurned
	}

	// Convert meal summaries to slice
	var meals []*analyticspb.MealSummary
	for _, meal := range mealSummaries {
		meals = append(meals, meal)
	}

	return &analyticspb.DailyNutritionSummary{
		Date:            date.Format("2006-01-02"),
		CaloriesIntake:  int32(totalCalories),
		CaloriesBurned:  int32(totalCaloriesBurned),
		CaloriesBalance: int32(totalCalories - totalCaloriesBurned),
		Protein:         int32(totalProtein),
		Carbohydrate:    int32(totalCarb),
		Fat:             int32(totalFat),
		Fiber:           int32(totalFiber),
		Sugar:           int32(totalSugar),
		Sodium:          int32(totalSodium),
		Meals:           meals,
	}
}

func (s *AnalyticsService) calculateNutritionTrends(summaries []*analyticspb.DailyNutritionSummary) *analyticspb.NutritionTrends {
	if len(summaries) == 0 {
		return &analyticspb.NutritionTrends{}
	}

	// Calculate averages
	var totalCalories, totalProtein, totalCarb, totalFat float64
	for _, summary := range summaries {
		totalCalories += float64(summary.CaloriesIntake)
		totalProtein += float64(summary.Protein)
		totalCarb += float64(summary.Carbohydrate)
		totalFat += float64(summary.Fat)
	}

	count := float64(len(summaries))
	avgCalories := totalCalories / count
	avgProtein := totalProtein / count
	avgCarb := totalCarb / count
	avgFat := totalFat / count

	// Calculate trends (simplified - comparing first half to second half)
	var caloriesTrend, proteinTrend, carbTrend, fatTrend float64
	if len(summaries) >= 2 {
		mid := len(summaries) / 2
		firstHalf := summaries[:mid]
		secondHalf := summaries[mid:]

		firstHalfAvg := s.calculateAverage(firstHalf, "calories")
		secondHalfAvg := s.calculateAverage(secondHalf, "calories")
		caloriesTrend = ((secondHalfAvg - firstHalfAvg) / firstHalfAvg) * 100

		firstHalfAvg = s.calculateAverage(firstHalf, "protein")
		secondHalfAvg = s.calculateAverage(secondHalf, "protein")
		proteinTrend = ((secondHalfAvg - firstHalfAvg) / firstHalfAvg) * 100

		firstHalfAvg = s.calculateAverage(firstHalf, "carb")
		secondHalfAvg = s.calculateAverage(secondHalf, "carb")
		carbTrend = ((secondHalfAvg - firstHalfAvg) / firstHalfAvg) * 100

		firstHalfAvg = s.calculateAverage(firstHalf, "fat")
		secondHalfAvg = s.calculateAverage(secondHalf, "fat")
		fatTrend = ((secondHalfAvg - firstHalfAvg) / firstHalfAvg) * 100
	}

	return &analyticspb.NutritionTrends{
		CaloriesAvg:       avgCalories,
		ProteinAvg:        avgProtein,
		CarbohydrateAvg:   avgCarb,
		FatAvg:            avgFat,
		CaloriesTrend:     caloriesTrend,
		ProteinTrend:      proteinTrend,
		CarbohydrateTrend: carbTrend,
		FatTrend:          fatTrend,
	}
}

func (s *AnalyticsService) calculateAverage(summaries []*analyticspb.DailyNutritionSummary, nutrient string) float64 {
	if len(summaries) == 0 {
		return 0
	}

	var total float64
	for _, summary := range summaries {
		switch nutrient {
		case "calories":
			total += float64(summary.CaloriesIntake)
		case "protein":
			total += float64(summary.Protein)
		case "carb":
			total += float64(summary.Carbohydrate)
		case "fat":
			total += float64(summary.Fat)
		}
	}

	return total / float64(len(summaries))
}

func (s *AnalyticsService) calculateMonthlyInsights(year, month int, foodLogs []FoodLog) *analyticspb.MonthlyNutritionInsights {
	var totalCalories, totalProtein, totalCarb, totalFat int
	foodCounts := make(map[string]int)

	for _, log := range foodLogs {
		totalCalories += log.Calories
		totalProtein += log.Protein
		totalCarb += log.Carb
		totalFat += log.Fat
		foodCounts[log.FoodID]++
	}

	// Calculate daily averages
	daysInMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
	avgDailyCalories := totalCalories / daysInMonth
	avgDailyProtein := totalProtein / daysInMonth
	avgDailyCarb := totalCarb / daysInMonth
	avgDailyFat := totalFat / daysInMonth

	// Get top foods
	var topFoods []string
	for foodID, count := range foodCounts {
		if count >= 3 { // Foods eaten 3+ times
			topFoods = append(topFoods, foodID)
		}
	}

	// Generate recommendations
	var recommendations []string
	if avgDailyCalories < 1200 {
		recommendations = append(recommendations, "Consider increasing calorie intake for better nutrition")
	}
	if avgDailyProtein < 50 {
		recommendations = append(recommendations, "Increase protein intake for muscle health")
	}
	if avgDailyFat < 20 {
		recommendations = append(recommendations, "Consider adding healthy fats to your diet")
	}

	return &analyticspb.MonthlyNutritionInsights{
		Year:                 fmt.Sprintf("%d", year),
		Month:                fmt.Sprintf("%02d", month),
		TotalCalories:        int32(totalCalories),
		AvgDailyCalories:     int32(avgDailyCalories),
		TotalProtein:         int32(totalProtein),
		AvgDailyProtein:      int32(avgDailyProtein),
		TotalCarbohydrate:    int32(totalCarb),
		AvgDailyCarbohydrate: int32(avgDailyCarb),
		TotalFat:             int32(totalFat),
		AvgDailyFat:          int32(avgDailyFat),
		TopFoods:             topFoods,
		Recommendations:      recommendations,
	}
}

func (s *AnalyticsService) calculateWeightProgressAnalysis(startDate, endDate time.Time, weightLogs []WeightLog) *analyticspb.WeightProgressAnalysis {
	if len(weightLogs) == 0 {
		return &analyticspb.WeightProgressAnalysis{
			StartDate: startDate.Format("2006-01-02"),
			EndDate:   endDate.Format("2006-01-02"),
			Trend:     "no_data",
		}
	}

	// Sort by date
	for i := 0; i < len(weightLogs)-1; i++ {
		for j := i + 1; j < len(weightLogs); j++ {
			if weightLogs[i].Date.After(weightLogs[j].Date) {
				weightLogs[i], weightLogs[j] = weightLogs[j], weightLogs[i]
			}
		}
	}

	startWeight := weightLogs[0].Weight
	endWeight := weightLogs[len(weightLogs)-1].Weight
	weightChange := endWeight - startWeight
	weightChangePercentage := (weightChange / startWeight) * 100

	// Calculate weekly average change
	weeks := int(endDate.Sub(startDate).Hours() / (24 * 7))
	var avgWeeklyChange float64
	if weeks > 0 {
		avgWeeklyChange = weightChange / float64(weeks)
	}

	// Convert to data points
	var dataPoints []*analyticspb.WeightDataPoint
	for _, log := range weightLogs {
		dataPoints = append(dataPoints, &analyticspb.WeightDataPoint{
			Date:   log.Date.Format("2006-01-02"),
			Weight: log.Weight,
		})
	}

	// Determine trend
	var trend string
	if weightChange > 0.5 {
		trend = "increasing"
	} else if weightChange < -0.5 {
		trend = "decreasing"
	} else {
		trend = "stable"
	}

	return &analyticspb.WeightProgressAnalysis{
		StartDate:              startDate.Format("2006-01-02"),
		EndDate:                endDate.Format("2006-01-02"),
		StartWeight:            startWeight,
		EndWeight:              endWeight,
		WeightChange:           weightChange,
		WeightChangePercentage: weightChangePercentage,
		AvgWeeklyChange:        avgWeeklyChange,
		DataPoints:             dataPoints,
		Trend:                  trend,
	}
}

func (s *AnalyticsService) calculateCalorieBalanceAnalysis(startDate, endDate time.Time, foodLogs []FoodLog, exerciseLogs []ExerciseLog) *analyticspb.CalorieBalanceAnalysis {
	// Calculate total calories intake
	var totalCaloriesIntake int
	for _, log := range foodLogs {
		totalCaloriesIntake += log.Calories
	}

	// Calculate total calories burned
	var totalCaloriesBurned int
	for _, log := range exerciseLogs {
		totalCaloriesBurned += log.CaloriesBurned
	}

	totalCalorieBalance := totalCaloriesIntake - totalCaloriesBurned

	// Calculate daily averages
	days := int(endDate.Sub(startDate).Hours()/24) + 1
	avgDailyBalance := float64(totalCalorieBalance) / float64(days)

	// Count days in deficit and surplus
	daysInDeficit := 0
	daysInSurplus := 0
	currentDate := startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		dailyFoodLogs, _ := s.getFoodLogsForDate(context.Background(), "", currentDate)
		dailyExerciseLogs, _ := s.getExerciseLogsForDate(context.Background(), "", currentDate)

		var dailyIntake, dailyBurned int
		for _, log := range dailyFoodLogs {
			dailyIntake += log.Calories
		}
		for _, log := range dailyExerciseLogs {
			dailyBurned += log.CaloriesBurned
		}

		balance := dailyIntake - dailyBurned
		if balance < 0 {
			daysInDeficit++
		} else if balance > 0 {
			daysInSurplus++
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	deficitPercentage := float64(daysInDeficit) / float64(days) * 100

	// Create daily balances
	var dailyBalances []*analyticspb.DailyBalance
	currentDate = startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		dailyFoodLogs, _ := s.getFoodLogsForDate(context.Background(), "", currentDate)
		dailyExerciseLogs, _ := s.getExerciseLogsForDate(context.Background(), "", currentDate)

		var dailyIntake, dailyBurned int
		for _, log := range dailyFoodLogs {
			dailyIntake += log.Calories
		}
		for _, log := range dailyExerciseLogs {
			dailyBurned += log.CaloriesBurned
		}

		dailyBalances = append(dailyBalances, &analyticspb.DailyBalance{
			Date:           currentDate.Format("2006-01-02"),
			CaloriesIntake: int32(dailyIntake),
			CaloriesBurned: int32(dailyBurned),
			Balance:        int32(dailyIntake - dailyBurned),
		})

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return &analyticspb.CalorieBalanceAnalysis{
		StartDate:           startDate.Format("2006-01-02"),
		EndDate:             endDate.Format("2006-01-02"),
		TotalCaloriesIntake: int32(totalCaloriesIntake),
		TotalCaloriesBurned: int32(totalCaloriesBurned),
		TotalCalorieBalance: int32(totalCalorieBalance),
		AvgDailyBalance:     avgDailyBalance,
		DaysInDeficit:       int32(daysInDeficit),
		DaysInSurplus:       int32(daysInSurplus),
		DeficitPercentage:   deficitPercentage,
		DailyBalances:       dailyBalances,
	}
}
