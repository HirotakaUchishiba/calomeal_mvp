package middleware

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// MetadataKeys defines the standard metadata keys used across services
const (
	UserIDKey    = "user_id"
	TraceIDKey   = "x-trace-id"
	RequestIDKey = "x-request-id"
	EmailKey     = "email"
)

// GetMetadataFromContext extracts metadata from incoming gRPC context
func GetMetadataFromContext(ctx context.Context) (userID, traceID, requestID, email string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", "", "", status.Errorf(codes.Unauthenticated, "metadata not found in context")
	}

	userIDs := md.Get(UserIDKey)
	if len(userIDs) == 0 {
		return "", "", "", "", status.Errorf(codes.Unauthenticated, "user ID not found in metadata")
	}
	userID = userIDs[0]

	traceIDs := md.Get(TraceIDKey)
	if len(traceIDs) > 0 {
		traceID = traceIDs[0]
	}

	requestIDs := md.Get(RequestIDKey)
	if len(requestIDs) > 0 {
		requestID = requestIDs[0]
	}

	emails := md.Get(EmailKey)
	if len(emails) > 0 {
		email = emails[0]
	}

	return userID, traceID, requestID, email, nil
}

// ValidateUserID validates that the user ID in metadata matches the request
func ValidateUserID(ctx context.Context, requestUserID string) error {
	metadataUserID, _, _, _, err := GetMetadataFromContext(ctx)
	if err != nil {
		return err
	}

	if metadataUserID != requestUserID {
		return status.Errorf(codes.PermissionDenied, "user ID mismatch: metadata=%s, request=%s", metadataUserID, requestUserID)
	}

	return nil
}

// LogWithMetadata creates a structured log entry with metadata
func LogWithMetadata(ctx context.Context, level, message string, fields ...interface{}) {
	userID, traceID, requestID, email, err := GetMetadataFromContext(ctx)
	if err != nil {
		// Fallback to basic logging if metadata is not available
		fmt.Printf("[%s] %s\n", level, message)
		return
	}

	logEntry := fmt.Sprintf("[%s] user=%s trace=%s request=%s email=%s %s", 
		level, userID, traceID, requestID, email, message)
	
	if len(fields) > 0 {
		logEntry += fmt.Sprintf(" %v", fields)
	}
	
	fmt.Println(logEntry)
}
