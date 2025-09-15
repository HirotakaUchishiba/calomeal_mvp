package logger

// Context key types to avoid collisions
type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	TraceIDKey   contextKey = "trace_id"
	UserIDKey    contextKey = "user_id"
	EmailKey     contextKey = "email"
)
