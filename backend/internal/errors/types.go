package errors

import (
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorCode represents a standardized error code
type ErrorCode string

const (
	// Authentication & Authorization
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeTokenExpired ErrorCode = "TOKEN_EXPIRED"
	ErrCodeInvalidToken ErrorCode = "INVALID_TOKEN"

	// Validation
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"
	ErrCodeMissingField ErrorCode = "MISSING_FIELD"

	// Database
	ErrCodeDatabase            ErrorCode = "DATABASE_ERROR"
	ErrCodeRecordNotFound      ErrorCode = "RECORD_NOT_FOUND"
	ErrCodeDuplicateRecord     ErrorCode = "DUPLICATE_RECORD"
	ErrCodeConstraintViolation ErrorCode = "CONSTRAINT_VIOLATION"

	// External Services
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrCodeServiceTimeout     ErrorCode = "SERVICE_TIMEOUT"
	ErrCodeServiceError       ErrorCode = "SERVICE_ERROR"
	ErrCodeNetworkError       ErrorCode = "NETWORK_ERROR"

	// Business Logic
	ErrCodeBusinessRule     ErrorCode = "BUSINESS_RULE_VIOLATION"
	ErrCodeInsufficientData ErrorCode = "INSUFFICIENT_DATA"
	ErrCodeRateLimit        ErrorCode = "RATE_LIMIT_EXCEEDED"

	// System
	ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrCodeNotImplemented ErrorCode = "NOT_IMPLEMENTED"
	ErrCodeConfiguration  ErrorCode = "CONFIGURATION_ERROR"
)

// ErrorSeverity represents the severity level of an error
type ErrorSeverity string

const (
	SeverityLow      ErrorSeverity = "LOW"
	SeverityMedium   ErrorSeverity = "MEDIUM"
	SeverityHigh     ErrorSeverity = "HIGH"
	SeverityCritical ErrorSeverity = "CRITICAL"
)

// CaloMealError represents a standardized error for the application
type CaloMealError struct {
	Code       ErrorCode     `json:"code"`
	Message    string        `json:"message"`
	Details    string        `json:"details,omitempty"`
	Severity   ErrorSeverity `json:"severity"`
	Timestamp  time.Time     `json:"timestamp"`
	RequestID  string        `json:"request_id,omitempty"`
	TraceID    string        `json:"trace_id,omitempty"`
	UserID     string        `json:"user_id,omitempty"`
	Service    string        `json:"service,omitempty"`
	Operation  string        `json:"operation,omitempty"`
	Cause      error         `json:"-"`
	Retryable  bool          `json:"retryable"`
	HTTPStatus int           `json:"http_status"`
}

// Error implements the error interface
func (e *CaloMealError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *CaloMealError) Unwrap() error {
	return e.Cause
}

// New creates a new CaloMealError
func New(code ErrorCode, message string) *CaloMealError {
	return &CaloMealError{
		Code:       code,
		Message:    message,
		Severity:   getDefaultSeverity(code),
		Timestamp:  time.Now(),
		Retryable:  isRetryable(code),
		HTTPStatus: getHTTPStatus(code),
	}
}

// Wrap wraps an existing error with CaloMealError
func Wrap(err error, code ErrorCode, message string) *CaloMealError {
	if err == nil {
		return nil
	}

	caloMealErr := New(code, message)
	caloMealErr.Cause = err

	// If the wrapped error is already a CaloMealError, preserve some fields
	if existingErr, ok := err.(*CaloMealError); ok {
		caloMealErr.RequestID = existingErr.RequestID
		caloMealErr.TraceID = existingErr.TraceID
		caloMealErr.UserID = existingErr.UserID
		caloMealErr.Service = existingErr.Service
		caloMealErr.Operation = existingErr.Operation
	}

	return caloMealErr
}

// WithDetails adds details to the error
func (e *CaloMealError) WithDetails(details string) *CaloMealError {
	e.Details = details
	return e
}

// WithRequestID adds request ID to the error
func (e *CaloMealError) WithRequestID(requestID string) *CaloMealError {
	e.RequestID = requestID
	return e
}

// WithTraceID adds trace ID to the error
func (e *CaloMealError) WithTraceID(traceID string) *CaloMealError {
	e.TraceID = traceID
	return e
}

// WithUserID adds user ID to the error
func (e *CaloMealError) WithUserID(userID string) *CaloMealError {
	e.UserID = userID
	return e
}

// WithService adds service name to the error
func (e *CaloMealError) WithService(service string) *CaloMealError {
	e.Service = service
	return e
}

// WithOperation adds operation name to the error
func (e *CaloMealError) WithOperation(operation string) *CaloMealError {
	e.Operation = operation
	return e
}

// WithSeverity sets the error severity
func (e *CaloMealError) WithSeverity(severity ErrorSeverity) *CaloMealError {
	e.Severity = severity
	return e
}

// ToGRPCStatus converts CaloMealError to gRPC status
func (e *CaloMealError) ToGRPCStatus() *status.Status {
	grpcCode := getGRPCCode(e.Code)
	st := status.New(grpcCode, e.Message)

	// Add error details as metadata (simplified version without protobuf)
	// Note: In a real implementation, you would define proper protobuf messages
	// for error details and use st.WithDetails() with those messages

	return st
}

// ErrorDetails represents error details for gRPC
type ErrorDetails struct {
	Code      string `json:"code"`
	Details   string `json:"details"`
	RequestId string `json:"request_id"`
	TraceId   string `json:"trace_id"`
	UserId    string `json:"user_id"`
	Service   string `json:"service"`
	Operation string `json:"operation"`
}

// getDefaultSeverity returns the default severity for an error code
func getDefaultSeverity(code ErrorCode) ErrorSeverity {
	switch code {
	case ErrCodeUnauthorized, ErrCodeForbidden, ErrCodeTokenExpired, ErrCodeInvalidToken:
		return SeverityMedium
	case ErrCodeValidation, ErrCodeInvalidInput, ErrCodeMissingField:
		return SeverityLow
	case ErrCodeDatabase, ErrCodeRecordNotFound, ErrCodeDuplicateRecord, ErrCodeConstraintViolation:
		return SeverityMedium
	case ErrCodeServiceUnavailable, ErrCodeServiceTimeout, ErrCodeServiceError, ErrCodeNetworkError:
		return SeverityHigh
	case ErrCodeBusinessRule, ErrCodeInsufficientData, ErrCodeRateLimit:
		return SeverityMedium
	case ErrCodeInternal, ErrCodeNotImplemented, ErrCodeConfiguration:
		return SeverityCritical
	default:
		return SeverityMedium
	}
}

// isRetryable determines if an error is retryable
func isRetryable(code ErrorCode) bool {
	switch code {
	case ErrCodeServiceUnavailable, ErrCodeServiceTimeout, ErrCodeNetworkError:
		return true
	case ErrCodeDatabase, ErrCodeInternal:
		return true
	default:
		return false
	}
}

// getHTTPStatus returns the appropriate HTTP status code
func getHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrCodeUnauthorized, ErrCodeTokenExpired, ErrCodeInvalidToken:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeValidation, ErrCodeInvalidInput, ErrCodeMissingField:
		return http.StatusBadRequest
	case ErrCodeRecordNotFound:
		return http.StatusNotFound
	case ErrCodeDuplicateRecord, ErrCodeConstraintViolation:
		return http.StatusConflict
	case ErrCodeServiceUnavailable, ErrCodeServiceTimeout, ErrCodeServiceError, ErrCodeNetworkError:
		return http.StatusServiceUnavailable
	case ErrCodeRateLimit:
		return http.StatusTooManyRequests
	case ErrCodeBusinessRule, ErrCodeInsufficientData:
		return http.StatusUnprocessableEntity
	case ErrCodeNotImplemented:
		return http.StatusNotImplemented
	case ErrCodeInternal, ErrCodeDatabase, ErrCodeConfiguration:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// getGRPCCode returns the appropriate gRPC status code
func getGRPCCode(code ErrorCode) codes.Code {
	switch code {
	case ErrCodeUnauthorized, ErrCodeTokenExpired, ErrCodeInvalidToken:
		return codes.Unauthenticated
	case ErrCodeForbidden:
		return codes.PermissionDenied
	case ErrCodeValidation, ErrCodeInvalidInput, ErrCodeMissingField:
		return codes.InvalidArgument
	case ErrCodeRecordNotFound:
		return codes.NotFound
	case ErrCodeDuplicateRecord, ErrCodeConstraintViolation:
		return codes.AlreadyExists
	case ErrCodeServiceUnavailable, ErrCodeServiceTimeout, ErrCodeServiceError, ErrCodeNetworkError:
		return codes.Unavailable
	case ErrCodeRateLimit:
		return codes.ResourceExhausted
	case ErrCodeBusinessRule, ErrCodeInsufficientData:
		return codes.FailedPrecondition
	case ErrCodeNotImplemented:
		return codes.Unimplemented
	case ErrCodeInternal, ErrCodeDatabase, ErrCodeConfiguration:
		return codes.Internal
	default:
		return codes.Internal
	}
}
