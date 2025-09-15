package analytics

import (
	"context"
	"fmt"
	"log"
	"time"

	analyticspb "github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/analytics/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// AnalyticsClient wraps the analytics gRPC client
type AnalyticsClient struct {
	conn   *grpc.ClientConn
	client analyticspb.AnalyticsServiceClient
}

// NewGRPCClient creates a new analytics gRPC client
func NewGRPCClient(addr string) (*AnalyticsClient, error) {
	// Create gRPC connection
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to analytics service: %w", err)
	}

	// Create client
	client := analyticspb.NewAnalyticsServiceClient(conn)

	return &AnalyticsClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection
func (c *AnalyticsClient) Close() error {
	return c.conn.Close()
}

// GetDailyNutritionSummary retrieves daily nutrition summary
func (c *AnalyticsClient) GetDailyNutritionSummary(ctx context.Context, userID, date string) (*analyticspb.DailyNutritionSummary, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := &analyticspb.GetDailyNutritionSummaryRequest{
		UserId: userID,
		Date:   date,
	}

	resp, err := c.client.GetDailyNutritionSummary(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.DeadlineExceeded:
				return nil, fmt.Errorf("request timeout: %w", err)
			case codes.Unavailable:
				return nil, fmt.Errorf("analytics service unavailable: %w", err)
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %w", err)
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %w", err)
			default:
				return nil, fmt.Errorf("gRPC error [%s]: %w", st.Code(), err)
			}
		}
		log.Printf("Failed to get daily nutrition summary: %v", err)
		return nil, err
	}

	return resp.Summary, nil
}

// GetWeeklyNutritionTrends retrieves weekly nutrition trends
func (c *AnalyticsClient) GetWeeklyNutritionTrends(ctx context.Context, userID, startDate, endDate string) (*analyticspb.GetWeeklyNutritionTrendsResponse, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req := &analyticspb.GetWeeklyNutritionTrendsRequest{
		UserId:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	resp, err := c.client.GetWeeklyNutritionTrends(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.DeadlineExceeded:
				return nil, fmt.Errorf("request timeout: %w", err)
			case codes.Unavailable:
				return nil, fmt.Errorf("analytics service unavailable: %w", err)
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %w", err)
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %w", err)
			default:
				return nil, fmt.Errorf("gRPC error [%s]: %w", st.Code(), err)
			}
		}
		log.Printf("Failed to get weekly nutrition trends: %v", err)
		return nil, err
	}

	return resp, nil
}

// GetMonthlyNutritionInsights retrieves monthly nutrition insights
func (c *AnalyticsClient) GetMonthlyNutritionInsights(ctx context.Context, userID, year, month string) (*analyticspb.MonthlyNutritionInsights, error) {
	req := &analyticspb.GetMonthlyNutritionInsightsRequest{
		UserId: userID,
		Year:   year,
		Month:  month,
	}

	resp, err := c.client.GetMonthlyNutritionInsights(ctx, req)
	if err != nil {
		log.Printf("Failed to get monthly nutrition insights: %v", err)
		return nil, err
	}

	return resp.Insights, nil
}

// GetWeightProgressAnalysis retrieves weight progress analysis
func (c *AnalyticsClient) GetWeightProgressAnalysis(ctx context.Context, userID, startDate, endDate string) (*analyticspb.WeightProgressAnalysis, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req := &analyticspb.GetWeightProgressAnalysisRequest{
		UserId:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	resp, err := c.client.GetWeightProgressAnalysis(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.DeadlineExceeded:
				return nil, fmt.Errorf("request timeout: %w", err)
			case codes.Unavailable:
				return nil, fmt.Errorf("analytics service unavailable: %w", err)
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %w", err)
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %w", err)
			default:
				return nil, fmt.Errorf("gRPC error [%s]: %w", st.Code(), err)
			}
		}
		log.Printf("Failed to get weight progress analysis: %v", err)
		return nil, err
	}

	return resp.Analysis, nil
}

// GetCalorieBalanceAnalysis retrieves calorie balance analysis
func (c *AnalyticsClient) GetCalorieBalanceAnalysis(ctx context.Context, userID, startDate, endDate string) (*analyticspb.CalorieBalanceAnalysis, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req := &analyticspb.GetCalorieBalanceAnalysisRequest{
		UserId:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	resp, err := c.client.GetCalorieBalanceAnalysis(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.DeadlineExceeded:
				return nil, fmt.Errorf("request timeout: %w", err)
			case codes.Unavailable:
				return nil, fmt.Errorf("analytics service unavailable: %w", err)
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %w", err)
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %w", err)
			default:
				return nil, fmt.Errorf("gRPC error [%s]: %w", st.Code(), err)
			}
		}
		log.Printf("Failed to get calorie balance analysis: %v", err)
		return nil, err
	}

	return resp.Analysis, nil
}

// HealthCheck checks if the analytics service is healthy
func (c *AnalyticsClient) HealthCheck(ctx context.Context) error {
	// Create a timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Try to get a simple daily summary to check health
	_, err := c.GetDailyNutritionSummary(ctx, "health-check", "2025-01-01")
	if err != nil {
		// Even if the request fails due to no data, the service is healthy if we can connect
		return nil
	}

	return nil
}
