package errors

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreakerConfig represents the configuration for a circuit breaker
type CircuitBreakerConfig struct {
	FailureThreshold    int           `json:"failure_threshold"`    // Number of failures to open circuit
	SuccessThreshold    int           `json:"success_threshold"`    // Number of successes to close circuit (half-open state)
	Timeout             time.Duration `json:"timeout"`              // Time to wait before trying half-open
	MaxRequests         int           `json:"max_requests"`         // Max requests in half-open state
	FailureRate         float64       `json:"failure_rate"`         // Failure rate threshold (0.0-1.0)
	WindowSize          time.Duration `json:"window_size"`          // Time window for failure rate calculation
	MinRequestThreshold int           `json:"min_request_threshold"` // Minimum requests before calculating failure rate
}

// DefaultCircuitBreakerConfig returns a default circuit breaker configuration
func DefaultCircuitBreakerConfig() *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		FailureThreshold:    5,
		SuccessThreshold:    3,
		Timeout:             30 * time.Second,
		MaxRequests:         3,
		FailureRate:         0.5,
		WindowSize:          60 * time.Second,
		MinRequestThreshold: 10,
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config     *CircuitBreakerConfig
	state      CircuitState
	failures   int
	successes  int
	requests   int
	lastFailure time.Time
	nextAttempt time.Time
	mutex      sync.RWMutex
	metrics    *CircuitBreakerMetrics
}

// CircuitBreakerMetrics tracks circuit breaker metrics
type CircuitBreakerMetrics struct {
	TotalRequests   int64     `json:"total_requests"`
	TotalFailures   int64     `json:"total_failures"`
	TotalSuccesses  int64     `json:"total_successes"`
	StateChanges    int64     `json:"state_changes"`
	LastStateChange time.Time `json:"last_state_change"`
	WindowStart     time.Time `json:"window_start"`
	WindowRequests  int64     `json:"window_requests"`
	WindowFailures  int64     `json:"window_failures"`
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config *CircuitBreakerConfig) *CircuitBreaker {
	if config == nil {
		config = DefaultCircuitBreakerConfig()
	}

	return &CircuitBreaker{
		config:  config,
		state:   StateClosed,
		metrics: &CircuitBreakerMetrics{},
	}
}

// Execute executes a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	// Check if circuit is open
	if !cb.allowRequest() {
		return New(ErrCodeServiceUnavailable, "circuit breaker is open").
			WithDetails(fmt.Sprintf("circuit state: %s, next attempt: %s", cb.getState(), cb.getNextAttempt()))
	}

	// Execute the function
	err := fn()
	
	// Record the result
	cb.recordResult(err)
	
	return err
}

// ExecuteWithResult executes a function with circuit breaker protection and returns a result
func (cb *CircuitBreaker) ExecuteWithResult(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	// Check if circuit is open
	if !cb.allowRequest() {
		return nil, New(ErrCodeServiceUnavailable, "circuit breaker is open").
			WithDetails(fmt.Sprintf("circuit state: %s, next attempt: %s", cb.getState(), cb.getNextAttempt()))
	}

	// Execute the function
	result, err := fn()
	
	// Record the result
	cb.recordResult(err)
	
	return result, err
}

// allowRequest checks if a request should be allowed
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		return time.Now().After(cb.nextAttempt)
	case StateHalfOpen:
		return cb.requests < cb.config.MaxRequests
	default:
		return false
	}
}

// recordResult records the result of a request
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.metrics.TotalRequests++
	cb.metrics.WindowRequests++

	if err != nil {
		cb.metrics.TotalFailures++
		cb.metrics.WindowFailures++
		cb.failures++
		cb.lastFailure = time.Now()
	} else {
		cb.metrics.TotalSuccesses++
		cb.successes++
	}

	// Update state based on result
	cb.updateState(err)
}

// updateState updates the circuit breaker state based on the result
func (cb *CircuitBreaker) updateState(err error) {
	oldState := cb.state

	switch cb.state {
	case StateClosed:
		if cb.shouldOpen() {
			cb.state = StateOpen
			cb.nextAttempt = time.Now().Add(cb.config.Timeout)
		}
	case StateOpen:
		if time.Now().After(cb.nextAttempt) {
			cb.state = StateHalfOpen
			cb.requests = 0
			cb.successes = 0
		}
	case StateHalfOpen:
		if err != nil {
			cb.state = StateOpen
			cb.nextAttempt = time.Now().Add(cb.config.Timeout)
		} else if cb.successes >= cb.config.SuccessThreshold {
			cb.state = StateClosed
			cb.failures = 0
			cb.successes = 0
		}
	}

	// Update metrics if state changed
	if oldState != cb.state {
		cb.metrics.StateChanges++
		cb.metrics.LastStateChange = time.Now()
	}

	// Reset window if needed
	if time.Since(cb.metrics.WindowStart) > cb.config.WindowSize {
		cb.metrics.WindowStart = time.Now()
		cb.metrics.WindowRequests = 0
		cb.metrics.WindowFailures = 0
	}
}

// shouldOpen determines if the circuit should be opened
func (cb *CircuitBreaker) shouldOpen() bool {
	// Check failure count threshold
	if cb.failures >= cb.config.FailureThreshold {
		return true
	}

	// Check failure rate threshold
	if cb.metrics.WindowRequests >= int64(cb.config.MinRequestThreshold) {
		failureRate := float64(cb.metrics.WindowFailures) / float64(cb.metrics.WindowRequests)
		if failureRate >= cb.config.FailureRate {
			return true
		}
	}

	return false
}

// getState returns the current state
func (cb *CircuitBreaker) getState() CircuitState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// getNextAttempt returns the next attempt time
func (cb *CircuitBreaker) getNextAttempt() time.Time {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.nextAttempt
}

// GetMetrics returns the current metrics
func (cb *CircuitBreaker) GetMetrics() *CircuitBreakerMetrics {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	
	// Return a copy to avoid race conditions
	return &CircuitBreakerMetrics{
		TotalRequests:   cb.metrics.TotalRequests,
		TotalFailures:   cb.metrics.TotalFailures,
		TotalSuccesses:  cb.metrics.TotalSuccesses,
		StateChanges:    cb.metrics.StateChanges,
		LastStateChange: cb.metrics.LastStateChange,
		WindowStart:     cb.metrics.WindowStart,
		WindowRequests:  cb.metrics.WindowRequests,
		WindowFailures:  cb.metrics.WindowFailures,
	}
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	cb.state = StateClosed
	cb.failures = 0
	cb.successes = 0
	cb.requests = 0
	cb.nextAttempt = time.Time{}
}

// GetState returns the current state as a string
func (cb *CircuitBreaker) GetState() string {
	return cb.getState().String()
}

// IsOpen returns true if the circuit is open
func (cb *CircuitBreaker) IsOpen() bool {
	return cb.getState() == StateOpen
}

// IsClosed returns true if the circuit is closed
func (cb *CircuitBreaker) IsClosed() bool {
	return cb.getState() == StateClosed
}

// IsHalfOpen returns true if the circuit is half-open
func (cb *CircuitBreaker) IsHalfOpen() bool {
	return cb.getState() == StateHalfOpen
}

// CircuitBreakerManager manages multiple circuit breakers
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	mutex    sync.RWMutex
}

// NewCircuitBreakerManager creates a new circuit breaker manager
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// GetOrCreate gets an existing circuit breaker or creates a new one
func (m *CircuitBreakerManager) GetOrCreate(name string, config *CircuitBreakerConfig) *CircuitBreaker {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if breaker, exists := m.breakers[name]; exists {
		return breaker
	}

	breaker := NewCircuitBreaker(config)
	m.breakers[name] = breaker
	return breaker
}

// Get gets an existing circuit breaker
func (m *CircuitBreakerManager) Get(name string) (*CircuitBreaker, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	breaker, exists := m.breakers[name]
	return breaker, exists
}

// GetAll returns all circuit breakers
func (m *CircuitBreakerManager) GetAll() map[string]*CircuitBreaker {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[string]*CircuitBreaker)
	for name, breaker := range m.breakers {
		result[name] = breaker
	}
	return result
}

// Reset resets a specific circuit breaker
func (m *CircuitBreakerManager) Reset(name string) {
	m.mutex.RLock()
	breaker, exists := m.breakers[name]
	m.mutex.RUnlock()

	if exists {
		breaker.Reset()
	}
}

// ResetAll resets all circuit breakers
func (m *CircuitBreakerManager) ResetAll() {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, breaker := range m.breakers {
		breaker.Reset()
	}
}
