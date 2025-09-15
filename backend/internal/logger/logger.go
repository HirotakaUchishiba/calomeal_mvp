package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"
)

// LogLevel represents the logging level
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LogFormat represents the logging format
type LogFormat string

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
)

// LoggerConfig holds the configuration for the logger
type LoggerConfig struct {
	Level      LogLevel
	Format     LogFormat
	Output     io.Writer
	AddSource  bool
	TimeFormat string
}

// DefaultLoggerConfig returns a default logger configuration
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:      LogLevelInfo,
		Format:     LogFormatText,
		Output:     os.Stdout,
		AddSource:  false,
		TimeFormat: time.RFC3339,
	}
}

// ProductionLoggerConfig returns a production-optimized logger configuration
func ProductionLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:      LogLevelInfo,
		Format:     LogFormatJSON,
		Output:     os.Stdout,
		AddSource:  false,
		TimeFormat: time.RFC3339,
	}
}

// DevelopmentLoggerConfig returns a development-optimized logger configuration
func DevelopmentLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:      LogLevelDebug,
		Format:     LogFormatText,
		Output:     os.Stdout,
		AddSource:  true,
		TimeFormat: time.RFC3339,
	}
}

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger
	config *LoggerConfig
}

// NewLogger creates a new logger with the given configuration
func NewLogger(config *LoggerConfig) *Logger {
	if config == nil {
		config = DefaultLoggerConfig()
	}

	var handler slog.Handler

	// Create handler based on format
	if config.Format == LogFormatJSON {
		handler = slog.NewJSONHandler(config.Output, &slog.HandlerOptions{
			AddSource: config.AddSource,
			Level:     parseLogLevel(config.Level),
		})
	} else {
		handler = slog.NewTextHandler(config.Output, &slog.HandlerOptions{
			AddSource: config.AddSource,
			Level:     parseLogLevel(config.Level),
		})
	}

	return &Logger{
		Logger: slog.New(handler),
		config: config,
	}
}

// NewProductionLogger creates a production logger
func NewProductionLogger() *Logger {
	return NewLogger(ProductionLoggerConfig())
}

// NewDevelopmentLogger creates a development logger
func NewDevelopmentLogger() *Logger {
	return NewLogger(DevelopmentLoggerConfig())
}

// parseLogLevel converts LogLevel to slog.Level
func parseLogLevel(level LogLevel) slog.Level {
	switch level {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// WithContext creates a logger with context values
func (l *Logger) WithContext(ctx context.Context) *Logger {
	attrs := []slog.Attr{}

	// Extract common context values
	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		attrs = append(attrs, slog.String("user_id", userID))
	}
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		attrs = append(attrs, slog.String("trace_id", traceID))
	}
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		attrs = append(attrs, slog.String("request_id", requestID))
	}
	if email, ok := ctx.Value(EmailKey).(string); ok && email != "" {
		attrs = append(attrs, slog.String("email", email))
	}

	if len(attrs) > 0 {
		// Convert slog.Attr to []any
		args := make([]any, len(attrs)*2)
		for i, attr := range attrs {
			args[i*2] = attr.Key
			args[i*2+1] = attr.Value.Any()
		}
		return &Logger{
			Logger: l.Logger.With(args...),
			config: l.config,
		}
	}

	return l
}

// WithFields creates a logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	attrs := make([]slog.Attr, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}

	// Convert slog.Attr to []any
	args := make([]any, len(attrs)*2)
	for i, attr := range attrs {
		args[i*2] = attr.Key
		args[i*2+1] = attr.Value.Any()
	}
	return &Logger{
		Logger: l.Logger.With(args...),
		config: l.config,
	}
}

// LogRequest logs an HTTP request
func (l *Logger) LogRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration, userAgent string) {
	l.WithContext(ctx).Info("HTTP request",
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status_code", statusCode),
		slog.Duration("duration", duration),
		slog.String("user_agent", userAgent),
	)
}

// LogGRPCRequest logs a gRPC request
func (l *Logger) LogGRPCRequest(ctx context.Context, service, method string, duration time.Duration, statusCode string) {
	l.WithContext(ctx).Info("gRPC request",
		slog.String("service", service),
		slog.String("method", method),
		slog.Duration("duration", duration),
		slog.String("status_code", statusCode),
	)
}

// LogDatabaseQuery logs a database query
func (l *Logger) LogDatabaseQuery(ctx context.Context, query string, duration time.Duration, rowsAffected int64) {
	l.WithContext(ctx).Debug("Database query",
		slog.String("query", query),
		slog.Duration("duration", duration),
		slog.Int64("rows_affected", rowsAffected),
	)
}

// LogError logs an error with context
func (l *Logger) LogError(ctx context.Context, err error, message string, fields ...interface{}) {
	attrs := []slog.Attr{
		slog.String("error", err.Error()),
	}

	// Add additional fields
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			attrs = append(attrs, slog.Any(fields[i].(string), fields[i+1]))
		}
	}

	// Convert slog.Attr to []any
	args := make([]any, len(attrs)*2)
	for i, attr := range attrs {
		args[i*2] = attr.Key
		args[i*2+1] = attr.Value.Any()
	}
	l.WithContext(ctx).Error(message, args...)
}

// LogSecurityEvent logs a security-related event
func (l *Logger) LogSecurityEvent(ctx context.Context, event string, severity string, details map[string]interface{}) {
	attrs := []slog.Attr{
		slog.String("event", event),
		slog.String("severity", severity),
	}

	for k, v := range details {
		attrs = append(attrs, slog.Any(k, v))
	}

	// Convert slog.Attr to []any
	args := make([]any, len(attrs)*2)
	for i, attr := range attrs {
		args[i*2] = attr.Key
		args[i*2+1] = attr.Value.Any()
	}
	l.WithContext(ctx).Warn("Security event", args...)
}

// LogPerformance logs a performance metric
func (l *Logger) LogPerformance(ctx context.Context, operation string, duration time.Duration, metrics map[string]interface{}) {
	attrs := []slog.Attr{
		slog.String("operation", operation),
		slog.Duration("duration", duration),
	}

	for k, v := range metrics {
		attrs = append(attrs, slog.Any(k, v))
	}

	// Convert slog.Attr to []any
	args := make([]any, len(attrs)*2)
	for i, attr := range attrs {
		args[i*2] = attr.Key
		args[i*2+1] = attr.Value.Any()
	}
	l.WithContext(ctx).Info("Performance metric", args...)
}

// LogBusinessEvent logs a business-related event
func (l *Logger) LogBusinessEvent(ctx context.Context, event string, details map[string]interface{}) {
	attrs := []slog.Attr{
		slog.String("event", event),
	}

	for k, v := range details {
		attrs = append(attrs, slog.Any(k, v))
	}

	// Convert slog.Attr to []any
	args := make([]any, len(attrs)*2)
	for i, attr := range attrs {
		args[i*2] = attr.Key
		args[i*2+1] = attr.Value.Any()
	}
	l.WithContext(ctx).Info("Business event", args...)
}

// LogStartup logs application startup information
func (l *Logger) LogStartup(version, buildTime, gitCommit string, config map[string]interface{}) {
	attrs := []slog.Attr{
		slog.String("version", version),
		slog.String("build_time", buildTime),
		slog.String("git_commit", gitCommit),
	}

	for k, v := range config {
		attrs = append(attrs, slog.Any(k, v))
	}

	// Convert slog.Attr to []any
	args := make([]any, len(attrs)*2)
	for i, attr := range attrs {
		args[i*2] = attr.Key
		args[i*2+1] = attr.Value.Any()
	}
	l.Info("Application started", args...)
}

// LogShutdown logs application shutdown information
func (l *Logger) LogShutdown(reason string, duration time.Duration) {
	l.Info("Application shutting down",
		slog.String("reason", reason),
		slog.Duration("duration", duration),
	)
}

// LogHealthCheck logs health check results
func (l *Logger) LogHealthCheck(ctx context.Context, service string, status string, duration time.Duration, details map[string]interface{}) {
	attrs := []slog.Attr{
		slog.String("service", service),
		slog.String("status", status),
		slog.Duration("duration", duration),
	}

	for k, v := range details {
		attrs = append(attrs, slog.Any(k, v))
	}

	// Convert slog.Attr to []any
	args := make([]any, len(attrs)*2)
	for i, attr := range attrs {
		args[i*2] = attr.Key
		args[i*2+1] = attr.Value.Any()
	}
	l.WithContext(ctx).Info("Health check", args...)
}

// LogAudit logs an audit event
func (l *Logger) LogAudit(ctx context.Context, action string, resource string, result string, details map[string]interface{}) {
	attrs := []slog.Attr{
		slog.String("action", action),
		slog.String("resource", resource),
		slog.String("result", result),
	}

	for k, v := range details {
		attrs = append(attrs, slog.Any(k, v))
	}

	// Convert slog.Attr to []any
	args := make([]any, len(attrs)*2)
	for i, attr := range attrs {
		args[i*2] = attr.Key
		args[i*2+1] = attr.Value.Any()
	}
	l.WithContext(ctx).Info("Audit event", args...)
}

// SetLevel dynamically changes the log level
func (l *Logger) SetLevel(level LogLevel) {
	l.config.Level = level
	// Note: slog doesn't support dynamic level changes, so we'd need to recreate the logger
	// This is a placeholder for future implementation
}

// IsDebugEnabled returns true if debug logging is enabled
func (l *Logger) IsDebugEnabled() bool {
	return l.config.Level == LogLevelDebug
}

// IsProduction returns true if the logger is configured for production
func (l *Logger) IsProduction() bool {
	return l.config.Format == LogFormatJSON && l.config.Level == LogLevelInfo
}

// GetConfig returns the logger configuration
func (l *Logger) GetConfig() *LoggerConfig {
	return l.config
}
