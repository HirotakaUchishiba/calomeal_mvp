package errors

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// RetryConfig represents the configuration for retry logic
type RetryConfig struct {
	MaxAttempts    int           `json:"max_attempts"`
	InitialDelay   time.Duration `json:"initial_delay"`
	MaxDelay       time.Duration `json:"max_delay"`
	BackoffFactor  float64       `json:"backoff_factor"`
	Jitter         bool          `json:"jitter"`
	RetryableCodes []ErrorCode   `json:"retryable_codes"`
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:   3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 2.0,
		Jitter:        true,
		RetryableCodes: []ErrorCode{
			ErrCodeServiceUnavailable,
			ErrCodeServiceTimeout,
			ErrCodeNetworkError,
			ErrCodeDatabase,
			ErrCodeInternal,
		},
	}
}

// RetryableFunc represents a function that can be retried
type RetryableFunc func() error

// RetryableFuncWithResult represents a function that returns a result and can be retried
type RetryableFuncWithResult[T any] func() (T, error)

// Retry executes a function with retry logic
func Retry(ctx context.Context, config *RetryConfig, fn RetryableFunc) error {
	if config == nil {
		config = DefaultRetryConfig()
	}

	var lastErr error
	delay := config.InitialDelay

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		default:
		}

		// Execute the function
		err := fn()
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Check if error is retryable
		if !isErrorRetryable(err, config.RetryableCodes) {
			return err
		}

		// Don't wait after the last attempt
		if attempt == config.MaxAttempts {
			break
		}

		// Calculate delay with exponential backoff
		delay = calculateDelay(delay, config)

		// Wait with jitter
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		case <-time.After(delay):
		}
	}

	return fmt.Errorf("retry failed after %d attempts: %w", config.MaxAttempts, lastErr)
}

// RetryWithResult executes a function with retry logic and returns a result
func RetryWithResult[T any](ctx context.Context, config *RetryConfig, fn RetryableFuncWithResult[T]) (T, error) {
	var zero T
	
	if config == nil {
		config = DefaultRetryConfig()
	}

	var lastErr error
	delay := config.InitialDelay

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return zero, fmt.Errorf("retry cancelled: %w", ctx.Err())
		default:
		}

		// Execute the function
		result, err := fn()
		if err == nil {
			return result, nil // Success
		}

		lastErr = err

		// Check if error is retryable
		if !isErrorRetryable(err, config.RetryableCodes) {
			return zero, err
		}

		// Don't wait after the last attempt
		if attempt == config.MaxAttempts {
			break
		}

		// Calculate delay with exponential backoff
		delay = calculateDelay(delay, config)

		// Wait with jitter
		select {
		case <-ctx.Done():
			return zero, fmt.Errorf("retry cancelled: %w", ctx.Err())
		case <-time.After(delay):
		}
	}

	return zero, fmt.Errorf("retry failed after %d attempts: %w", config.MaxAttempts, lastErr)
}

// isErrorRetryable checks if an error is retryable based on the configuration
func isErrorRetryable(err error, retryableCodes []ErrorCode) bool {
	// Check if it's a CaloMealError
	if caloMealErr, ok := err.(*CaloMealError); ok {
		// Check if the error is marked as retryable
		if caloMealErr.Retryable {
			return true
		}
		
		// Check if the error code is in the retryable codes list
		for _, code := range retryableCodes {
			if caloMealErr.Code == code {
				return true
			}
		}
		return false
	}

	// For non-CaloMealError, assume it's retryable if it's a network or service error
	// This is a fallback for external library errors
	return true
}

// calculateDelay calculates the delay for the next retry attempt
func calculateDelay(currentDelay time.Duration, config *RetryConfig) time.Duration {
	// Apply exponential backoff
	delay := time.Duration(float64(currentDelay) * config.BackoffFactor)
	
	// Cap at max delay
	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}
	
	// Add jitter if enabled
	if config.Jitter {
		jitter := time.Duration(rand.Float64() * float64(delay) * 0.1) // 10% jitter
		delay += jitter
	}
	
	return delay
}

// RetryConfigBuilder helps build retry configurations
type RetryConfigBuilder struct {
	config *RetryConfig
}

// NewRetryConfigBuilder creates a new retry config builder
func NewRetryConfigBuilder() *RetryConfigBuilder {
	return &RetryConfigBuilder{
		config: DefaultRetryConfig(),
	}
}

// WithMaxAttempts sets the maximum number of retry attempts
func (b *RetryConfigBuilder) WithMaxAttempts(attempts int) *RetryConfigBuilder {
	b.config.MaxAttempts = attempts
	return b
}

// WithInitialDelay sets the initial delay between retries
func (b *RetryConfigBuilder) WithInitialDelay(delay time.Duration) *RetryConfigBuilder {
	b.config.InitialDelay = delay
	return b
}

// WithMaxDelay sets the maximum delay between retries
func (b *RetryConfigBuilder) WithMaxDelay(delay time.Duration) *RetryConfigBuilder {
	b.config.MaxDelay = delay
	return b
}

// WithBackoffFactor sets the exponential backoff factor
func (b *RetryConfigBuilder) WithBackoffFactor(factor float64) *RetryConfigBuilder {
	b.config.BackoffFactor = factor
	return b
}

// WithJitter enables or disables jitter
func (b *RetryConfigBuilder) WithJitter(enabled bool) *RetryConfigBuilder {
	b.config.Jitter = enabled
	return b
}

// WithRetryableCodes sets the list of retryable error codes
func (b *RetryConfigBuilder) WithRetryableCodes(codes ...ErrorCode) *RetryConfigBuilder {
	b.config.RetryableCodes = codes
	return b
}

// Build returns the configured retry config
func (b *RetryConfigBuilder) Build() *RetryConfig {
	return b.config
}

// QuickRetryConfigs provides pre-configured retry configs for common scenarios
var QuickRetryConfigs = struct {
	Fast    *RetryConfig
	Medium  *RetryConfig
	Slow    *RetryConfig
	Network *RetryConfig
}{
	Fast: &RetryConfig{
		MaxAttempts:   2,
		InitialDelay:  50 * time.Millisecond,
		MaxDelay:      200 * time.Millisecond,
		BackoffFactor: 2.0,
		Jitter:        true,
		RetryableCodes: []ErrorCode{
			ErrCodeServiceTimeout,
			ErrCodeNetworkError,
		},
	},
	Medium: &RetryConfig{
		MaxAttempts:   3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      1 * time.Second,
		BackoffFactor: 2.0,
		Jitter:        true,
		RetryableCodes: []ErrorCode{
			ErrCodeServiceUnavailable,
			ErrCodeServiceTimeout,
			ErrCodeNetworkError,
			ErrCodeDatabase,
		},
	},
	Slow: &RetryConfig{
		MaxAttempts:   5,
		InitialDelay:  500 * time.Millisecond,
		MaxDelay:      10 * time.Second,
		BackoffFactor: 2.0,
		Jitter:        true,
		RetryableCodes: []ErrorCode{
			ErrCodeServiceUnavailable,
			ErrCodeServiceTimeout,
			ErrCodeNetworkError,
			ErrCodeDatabase,
			ErrCodeInternal,
		},
	},
	Network: &RetryConfig{
		MaxAttempts:   4,
		InitialDelay:  200 * time.Millisecond,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 1.5,
		Jitter:        true,
		RetryableCodes: []ErrorCode{
			ErrCodeNetworkError,
			ErrCodeServiceTimeout,
		},
	},
}
