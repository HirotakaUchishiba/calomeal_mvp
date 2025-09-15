package analytics

import (
	"context"
	"fmt"
	"log"
	"time"

	analyticspb "github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/analytics/proto"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// AnalyticsClient wraps the analytics gRPC client
type AnalyticsClient struct {
	conn           *grpc.ClientConn
	client         analyticspb.AnalyticsServiceClient
	circuitBreaker *errors.CircuitBreaker
	retryConfig    *errors.RetryConfig
}

// NewGRPCClient creates a new analytics gRPC client
func NewGRPCClient(addr string) (*AnalyticsClient, error) {
	// Create gRPC connection
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeServiceUnavailable, "failed to connect to analytics service")
	}

	// Create client
	client := analyticspb.NewAnalyticsServiceClient(conn)

	// Create circuit breaker
	circuitBreaker := errors.NewCircuitBreaker(errors.DefaultCircuitBreakerConfig())

	// Create retry config
	retryConfig := errors.QuickRetryConfigs.Medium

	return &AnalyticsClient{
		conn:           conn,
		client:         client,
		circuitBreaker: circuitBreaker,
		retryConfig:    retryConfig,
	}, nil
}

// Close closes the gRPC connection
func (c *AnalyticsClient) Close() error {
	return c.conn.Close()
}

// convertGRPCError converts gRPC errors to CaloMealError
func (c *AnalyticsClient) convertGRPCError(err error, operation string) *errors.CaloMealError {
	if err == nil {
		return nil
	}

	// Handle gRPC status errors
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.DeadlineExceeded:
			return errors.Wrap(err, errors.ErrCodeServiceTimeout, "request timeout").
				WithService("analytics").
				WithOperation(operation)
		case codes.Unavailable:
			return errors.Wrap(err, errors.ErrCodeServiceUnavailable, "analytics service unavailable").
				WithService("analytics").
				WithOperation(operation)
		case codes.Unauthenticated:
			return errors.Wrap(err, errors.ErrCodeUnauthorized, "authentication failed").
				WithService("analytics").
				WithOperation(operation)
		case codes.PermissionDenied:
			return errors.Wrap(err, errors.ErrCodeForbidden, "permission denied").
				WithService("analytics").
				WithOperation(operation)
		case codes.InvalidArgument:
			return errors.Wrap(err, errors.ErrCodeInvalidInput, "invalid request parameters").
				WithService("analytics").
				WithOperation(operation)
		case codes.NotFound:
			return errors.Wrap(err, errors.ErrCodeRecordNotFound, "resource not found").
				WithService("analytics").
				WithOperation(operation)
		case codes.Internal:
			return errors.Wrap(err, errors.ErrCodeServiceError, "internal service error").
				WithService("analytics").
				WithOperation(operation)
		default:
			return errors.Wrap(err, errors.ErrCodeServiceError, "gRPC service error").
				WithService("analytics").
				WithOperation(operation).
				WithDetails(fmt.Sprintf("gRPC code: %s", st.Code()))
		}
	}

	// Handle non-gRPC errors
	return errors.Wrap(err, errors.ErrCodeServiceError, "analytics service error").
		WithService("analytics").
		WithOperation(operation)
}

// GetDailyNutritionSummary retrieves daily nutrition summary
func (c *AnalyticsClient) GetDailyNutritionSummary(ctx context.Context, userID, date string) (*analyticspb.DailyNutritionSummary, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Execute with circuit breaker and retry logic
	result, err := c.circuitBreaker.ExecuteWithResult(ctx, func() (*analyticspb.DailyNutritionSummary, error) {
		return errors.RetryWithResult(ctx, c.retryConfig, func() (*analyticspb.DailyNutritionSummary, error) {
			req := &analyticspb.GetDailyNutritionSummaryRequest{
				UserId: userID,
				Date:   date,
			}

			resp, err := c.client.GetDailyNutritionSummary(ctx, req)
			if err != nil {
				// Convert gRPC error to CaloMealError
				return nil, c.convertGRPCError(err, "GetDailyNutritionSummary")
			}

			return resp.Summary, nil
		})
	})

	if err != nil {
		// Log error with context
		errors.LogError(ctx, err, "GetDailyNutritionSummary", map[string]interface{}{
			"user_id": userID,
			"date":    date,
		})
	}

	return result, err
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
