package errors

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"
)

// ErrorLogger provides structured logging for errors
type ErrorLogger struct {
	logger *slog.Logger
}

// NewErrorLogger creates a new error logger
func NewErrorLogger() *ErrorLogger {
	// Create a structured logger with JSON output
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		AddSource: true,
	}))

	return &ErrorLogger{
		logger: logger,
	}
}

// LogError logs an error with structured information
func (el *ErrorLogger) LogError(ctx context.Context, err error, operation string, additionalFields ...map[string]interface{}) {
	if err == nil {
		return
	}

	// Extract context information
	requestID := getContextValue(ctx, "request_id")
	traceID := getContextValue(ctx, "trace_id")
	userID := getContextValue(ctx, "user_id")
	service := getContextValue(ctx, "service")

	// Build log attributes
	attrs := []slog.Attr{
		slog.String("level", "ERROR"),
		slog.String("operation", operation),
		slog.String("error", err.Error()),
		slog.Time("timestamp", time.Now()),
	}

	// Add context information if available
	if requestID != "" {
		attrs = append(attrs, slog.String("request_id", requestID))
	}
	if traceID != "" {
		attrs = append(attrs, slog.String("trace_id", traceID))
	}
	if userID != "" {
		attrs = append(attrs, slog.String("user_id", userID))
	}
	if service != "" {
		attrs = append(attrs, slog.String("service", service))
	}

	// Add CaloMealError specific fields
	if caloMealErr, ok := err.(*CaloMealError); ok {
		attrs = append(attrs, 
			slog.String("error_code", string(caloMealErr.Code)),
			slog.String("severity", string(caloMealErr.Severity)),
			slog.Bool("retryable", caloMealErr.Retryable),
			slog.Int("http_status", caloMealErr.HTTPStatus),
		)
		
		if caloMealErr.Details != "" {
			attrs = append(attrs, slog.String("details", caloMealErr.Details))
		}
		if caloMealErr.Service != "" {
			attrs = append(attrs, slog.String("error_service", caloMealErr.Service))
		}
		if caloMealErr.Operation != "" {
			attrs = append(attrs, slog.String("error_operation", caloMealErr.Operation))
		}
	}

	// Add additional fields
	for _, fields := range additionalFields {
		for key, value := range fields {
			attrs = append(attrs, slog.Any(key, value))
		}
	}

	// Log the error
	el.logger.LogAttrs(ctx, slog.LevelError, "Error occurred", attrs...)
}

// LogWarning logs a warning with structured information
func (el *ErrorLogger) LogWarning(ctx context.Context, message string, operation string, additionalFields ...map[string]interface{}) {
	attrs := []slog.Attr{
		slog.String("level", "WARNING"),
		slog.String("operation", operation),
		slog.String("message", message),
		slog.Time("timestamp", time.Now()),
	}

	// Add context information
	if requestID := getContextValue(ctx, "request_id"); requestID != "" {
		attrs = append(attrs, slog.String("request_id", requestID))
	}
	if traceID := getContextValue(ctx, "trace_id"); traceID != "" {
		attrs = append(attrs, slog.String("trace_id", traceID))
	}
	if userID := getContextValue(ctx, "user_id"); userID != "" {
		attrs = append(attrs, slog.String("user_id", userID))
	}

	// Add additional fields
	for _, fields := range additionalFields {
		for key, value := range fields {
			attrs = append(attrs, slog.Any(key, value))
		}
	}

	el.logger.LogAttrs(ctx, slog.LevelWarn, "Warning", attrs...)
}

// LogInfo logs an info message with structured information
func (el *ErrorLogger) LogInfo(ctx context.Context, message string, operation string, additionalFields ...map[string]interface{}) {
	attrs := []slog.Attr{
		slog.String("level", "INFO"),
		slog.String("operation", operation),
		slog.String("message", message),
		slog.Time("timestamp", time.Now()),
	}

	// Add context information
	if requestID := getContextValue(ctx, "request_id"); requestID != "" {
		attrs = append(attrs, slog.String("request_id", requestID))
	}
	if traceID := getContextValue(ctx, "trace_id"); traceID != "" {
		attrs = append(attrs, slog.String("trace_id", traceID))
	}
	if userID := getContextValue(ctx, "user_id"); userID != "" {
		attrs = append(attrs, slog.String("user_id", userID))
	}

	// Add additional fields
	for _, fields := range additionalFields {
		for key, value := range fields {
			attrs = append(attrs, slog.Any(key, value))
		}
	}

	el.logger.LogAttrs(ctx, slog.LevelInfo, "Info", attrs...)
}

// LogCircuitBreakerState logs circuit breaker state changes
func (el *ErrorLogger) LogCircuitBreakerState(ctx context.Context, name string, oldState, newState CircuitState, metrics *CircuitBreakerMetrics) {
	attrs := []slog.Attr{
		slog.String("level", "INFO"),
		slog.String("operation", "circuit_breaker_state_change"),
		slog.String("circuit_breaker_name", name),
		slog.String("old_state", oldState.String()),
		slog.String("new_state", newState.String()),
		slog.Time("timestamp", time.Now()),
	}

	if metrics != nil {
		attrs = append(attrs,
			slog.Int64("total_requests", metrics.TotalRequests),
			slog.Int64("total_failures", metrics.TotalFailures),
			slog.Int64("total_successes", metrics.TotalSuccesses),
			slog.Int64("state_changes", metrics.StateChanges),
			slog.Time("last_state_change", metrics.LastStateChange),
		)
	}

	el.logger.LogAttrs(ctx, slog.LevelInfo, "Circuit breaker state changed", attrs...)
}

// LogRetryAttempt logs a retry attempt
func (el *ErrorLogger) LogRetryAttempt(ctx context.Context, operation string, attempt int, maxAttempts int, err error, delay time.Duration) {
	attrs := []slog.Attr{
		slog.String("level", "WARNING"),
		slog.String("operation", "retry_attempt"),
		slog.String("retry_operation", operation),
		slog.Int("attempt", attempt),
		slog.Int("max_attempts", maxAttempts),
		slog.Duration("delay", delay),
		slog.Time("timestamp", time.Now()),
	}

	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}

	el.logger.LogAttrs(ctx, slog.LevelWarn, "Retry attempt", attrs...)
}

// LogRetryFailure logs a retry failure
func (el *ErrorLogger) LogRetryFailure(ctx context.Context, operation string, maxAttempts int, finalErr error) {
	attrs := []slog.Attr{
		slog.String("level", "ERROR"),
		slog.String("operation", "retry_failure"),
		slog.String("retry_operation", operation),
		slog.Int("max_attempts", maxAttempts),
		slog.String("final_error", finalErr.Error()),
		slog.Time("timestamp", time.Now()),
	}

	el.logger.LogAttrs(ctx, slog.LevelError, "Retry failed", attrs...)
}

// LogRetrySuccess logs a retry success
func (el *ErrorLogger) LogRetrySuccess(ctx context.Context, operation string, attempts int) {
	attrs := []slog.Attr{
		slog.String("level", "INFO"),
		slog.String("operation", "retry_success"),
		slog.String("retry_operation", operation),
		slog.Int("attempts", attempts),
		slog.Time("timestamp", time.Now()),
	}

	el.logger.LogAttrs(ctx, slog.LevelInfo, "Retry succeeded", attrs...)
}

// LogServiceHealth logs service health information
func (el *ErrorLogger) LogServiceHealth(ctx context.Context, serviceName string, healthy bool, responseTime time.Duration, err error) {
	attrs := []slog.Attr{
		slog.String("level", "INFO"),
		slog.String("operation", "service_health_check"),
		slog.String("service_name", serviceName),
		slog.Bool("healthy", healthy),
		slog.Duration("response_time", responseTime),
		slog.Time("timestamp", time.Now()),
	}

	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}

	el.logger.LogAttrs(ctx, slog.LevelInfo, "Service health check", attrs...)
}

// getContextValue extracts a value from context
func getContextValue(ctx context.Context, key string) string {
	if value := ctx.Value(key); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// ErrorMetrics tracks error metrics
type ErrorMetrics struct {
	TotalErrors     int64                    `json:"total_errors"`
	ErrorsByCode    map[ErrorCode]int64      `json:"errors_by_code"`
	ErrorsByService map[string]int64         `json:"errors_by_service"`
	ErrorsBySeverity map[ErrorSeverity]int64 `json:"errors_by_severity"`
	LastError       time.Time                `json:"last_error"`
}

// ErrorMetricsCollector collects and tracks error metrics
type ErrorMetricsCollector struct {
	metrics ErrorMetrics
	mutex   sync.RWMutex
}

// NewErrorMetricsCollector creates a new error metrics collector
func NewErrorMetricsCollector() *ErrorMetricsCollector {
	return &ErrorMetricsCollector{
		metrics: ErrorMetrics{
			ErrorsByCode:     make(map[ErrorCode]int64),
			ErrorsByService:  make(map[string]int64),
			ErrorsBySeverity: make(map[ErrorSeverity]int64),
		},
	}
}

// RecordError records an error in the metrics
func (emc *ErrorMetricsCollector) RecordError(err error) {
	emc.mutex.Lock()
	defer emc.mutex.Unlock()

	emc.metrics.TotalErrors++
	emc.metrics.LastError = time.Now()

	if caloMealErr, ok := err.(*CaloMealError); ok {
		emc.metrics.ErrorsByCode[caloMealErr.Code]++
		emc.metrics.ErrorsBySeverity[caloMealErr.Severity]++
		
		if caloMealErr.Service != "" {
			emc.metrics.ErrorsByService[caloMealErr.Service]++
		}
	}
}

// GetMetrics returns the current metrics
func (emc *ErrorMetricsCollector) GetMetrics() ErrorMetrics {
	emc.mutex.RLock()
	defer emc.mutex.RUnlock()

	// Return a copy to avoid race conditions
	metrics := ErrorMetrics{
		TotalErrors:     emc.metrics.TotalErrors,
		LastError:       emc.metrics.LastError,
		ErrorsByCode:    make(map[ErrorCode]int64),
		ErrorsByService: make(map[string]int64),
		ErrorsBySeverity: make(map[ErrorSeverity]int64),
	}

	for code, count := range emc.metrics.ErrorsByCode {
		metrics.ErrorsByCode[code] = count
	}
	for service, count := range emc.metrics.ErrorsByService {
		metrics.ErrorsByService[service] = count
	}
	for severity, count := range emc.metrics.ErrorsBySeverity {
		metrics.ErrorsBySeverity[severity] = count
	}

	return metrics
}

// Reset resets the metrics
func (emc *ErrorMetricsCollector) Reset() {
	emc.mutex.Lock()
	defer emc.mutex.Unlock()

	emc.metrics = ErrorMetrics{
		ErrorsByCode:     make(map[ErrorCode]int64),
		ErrorsByService:  make(map[string]int64),
		ErrorsBySeverity: make(map[ErrorSeverity]int64),
	}
}

// Global instances
var (
	DefaultErrorLogger        = NewErrorLogger()
	DefaultErrorMetricsCollector = NewErrorMetricsCollector()
)

// LogError is a convenience function to log errors using the default logger
func LogError(ctx context.Context, err error, operation string, additionalFields ...map[string]interface{}) {
	DefaultErrorLogger.LogError(ctx, err, operation, additionalFields...)
	DefaultErrorMetricsCollector.RecordError(err)
}

// LogWarning is a convenience function to log warnings using the default logger
func LogWarning(ctx context.Context, message string, operation string, additionalFields ...map[string]interface{}) {
	DefaultErrorLogger.LogWarning(ctx, message, operation, additionalFields...)
}

// LogInfo is a convenience function to log info messages using the default logger
func LogInfo(ctx context.Context, message string, operation string, additionalFields ...map[string]interface{}) {
	DefaultErrorLogger.LogInfo(ctx, message, operation, additionalFields...)
}
