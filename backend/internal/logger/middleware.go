package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// LoggingMiddleware creates HTTP middleware for request logging
func LoggingMiddleware(logger *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Generate request ID
			requestID := uuid.New().String()
			
			// Add request ID to context
			ctx := context.WithValue(r.Context(), "request_id", requestID)
			ctx = context.WithValue(ctx, "trace_id", requestID) // Use same ID for trace
			
			// Add request ID to response headers
			w.Header().Set("X-Request-ID", requestID)
			
			// Create response writer wrapper to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			// Process request
			next.ServeHTTP(wrapped, r.WithContext(ctx))
			
			// Log request
			duration := time.Since(start)
			logger.LogRequest(ctx, r.Method, r.URL.Path, wrapped.statusCode, duration, r.UserAgent())
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// GRPCLoggingInterceptor creates gRPC interceptor for request logging
func GRPCLoggingInterceptor(logger *Logger) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		
		// Generate request ID if not present
		if _, ok := ctx.Value("request_id").(string); !ok {
			requestID := uuid.New().String()
			ctx = context.WithValue(ctx, "request_id", requestID)
			ctx = context.WithValue(ctx, "trace_id", requestID)
		}
		
		// Process request
		resp, err := handler(ctx, req)
		
		// Log request
		duration := time.Since(start)
		statusCode := "OK"
		if err != nil {
			statusCode = "ERROR"
		}
		
		service := "unknown"
		method := "unknown"
		if info != nil {
			service = info.Server
			method = info.FullMethod
		}
		
		logger.LogGRPCRequest(ctx, service, method, duration, statusCode)
		
		return resp, err
	}
}

// ContextLogger creates a logger with context values
func ContextLogger(ctx context.Context, logger *Logger) *Logger {
	return logger.WithContext(ctx)
}

// RequestLogger creates a logger for a specific request
func RequestLogger(r *http.Request, logger *Logger) *Logger {
	ctx := r.Context()
	return logger.WithContext(ctx)
}

// GRPCContextLogger creates a logger for a specific gRPC context
func GRPCContextLogger(ctx context.Context, logger *Logger) *Logger {
	return logger.WithContext(ctx)
}
