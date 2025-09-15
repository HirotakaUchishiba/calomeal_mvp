package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"google.golang.org/grpc/metadata"
)

// MetadataKeys defines the standard metadata keys used across services
const (
	TraceIDKey   = "x-trace-id"
	RequestIDKey = "x-request-id"
)

// AddMetadataToContext adds standard metadata to the context for gRPC calls
func AddMetadataToContext(ctx context.Context, userID, email string) context.Context {
	traceID := generateTraceID()
	requestID := generateRequestID()
	
	md := metadata.Pairs(
		string(UserIDKey), userID,
		TraceIDKey, traceID,
		RequestIDKey, requestID,
		string(EmailKey), email,
	)
	
	return metadata.NewOutgoingContext(ctx, md)
}

// GetMetadataFromContext extracts metadata from incoming gRPC context
func GetMetadataFromContext(ctx context.Context) (userID, traceID, requestID, email string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", "", "", fmt.Errorf("metadata not found in context")
	}

	userIDs := md.Get(string(UserIDKey))
	if len(userIDs) == 0 {
		return "", "", "", "", fmt.Errorf("user ID not found in metadata")
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

	emails := md.Get(string(EmailKey))
	if len(emails) > 0 {
		email = emails[0]
	}

	return userID, traceID, requestID, email, nil
}

// ValidateUserID validates that the user ID in metadata matches the request
func ValidateUserID(ctx context.Context, requestUserID string) error {
	metadataUserID, _, _, _, err := GetMetadataFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user ID from metadata: %w", err)
	}

	if metadataUserID != requestUserID {
		return fmt.Errorf("user ID mismatch: metadata=%s, request=%s", metadataUserID, requestUserID)
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

// generateTraceID generates a unique trace ID for request tracking
func generateTraceID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%d-%s", timestamp, hex.EncodeToString(bytes))
}
