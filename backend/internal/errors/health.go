package errors

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// HealthStatus represents the health status of a service
type HealthStatus string

const (
	StatusHealthy   HealthStatus = "HEALTHY"
	StatusUnhealthy HealthStatus = "UNHEALTHY"
	StatusDegraded  HealthStatus = "DEGRADED"
	StatusUnknown   HealthStatus = "UNKNOWN"
)

// HealthCheck represents a health check for a service
type HealthCheck struct {
	Name        string        `json:"name"`
	Status      HealthStatus  `json:"status"`
	Message     string        `json:"message,omitempty"`
	ResponseTime time.Duration `json:"response_time,omitempty"`
	LastCheck   time.Time     `json:"last_check"`
	Error       error         `json:"error,omitempty"`
}

// HealthChecker defines the interface for health checkers
type HealthChecker interface {
	Check(ctx context.Context) *HealthCheck
	GetName() string
}

// HealthCheckFunc is a function that performs a health check
type HealthCheckFunc func(ctx context.Context) error

// FuncHealthChecker implements HealthChecker using a function
type FuncHealthChecker struct {
	name string
	fn   HealthCheckFunc
}

// NewFuncHealthChecker creates a new function-based health checker
func NewFuncHealthChecker(name string, fn HealthCheckFunc) *FuncHealthChecker {
	return &FuncHealthChecker{
		name: name,
		fn:   fn,
	}
}

// Check performs the health check
func (hc *FuncHealthChecker) Check(ctx context.Context) *HealthCheck {
	start := time.Now()
	
	check := &HealthCheck{
		Name:      hc.name,
		Status:    StatusUnknown,
		LastCheck: time.Now(),
	}

	err := hc.fn(ctx)
	check.ResponseTime = time.Since(start)

	if err != nil {
		check.Status = StatusUnhealthy
		check.Message = err.Error()
		check.Error = err
	} else {
		check.Status = StatusHealthy
		check.Message = "Service is healthy"
	}

	return check
}

// GetName returns the health checker name
func (hc *FuncHealthChecker) GetName() string {
	return hc.name
}

// HealthManager manages health checks for multiple services
type HealthManager struct {
	checkers map[string]HealthChecker
	mutex    sync.RWMutex
	timeout  time.Duration
}

// NewHealthManager creates a new health manager
func NewHealthManager(timeout time.Duration) *HealthManager {
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	return &HealthManager{
		checkers: make(map[string]HealthChecker),
		timeout:  timeout,
	}
}

// RegisterHealthChecker registers a health checker
func (hm *HealthManager) RegisterHealthChecker(checker HealthChecker) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	hm.checkers[checker.GetName()] = checker
}

// RegisterFuncHealthChecker registers a function-based health checker
func (hm *HealthManager) RegisterFuncHealthChecker(name string, fn HealthCheckFunc) {
	checker := NewFuncHealthChecker(name, fn)
	hm.RegisterHealthChecker(checker)
}

// CheckHealth performs health checks for all registered services
func (hm *HealthManager) CheckHealth(ctx context.Context) map[string]*HealthCheck {
	hm.mutex.RLock()
	checkers := make(map[string]HealthChecker)
	for name, checker := range hm.checkers {
		checkers[name] = checker
	}
	hm.mutex.RUnlock()

	results := make(map[string]*HealthCheck)
	var wg sync.WaitGroup

	for name, checker := range checkers {
		wg.Add(1)
		go func(name string, checker HealthChecker) {
			defer wg.Done()

			// Create a timeout context for this check
			checkCtx, cancel := context.WithTimeout(ctx, hm.timeout)
			defer cancel()

			check := checker.Check(checkCtx)
			results[name] = check

			// Log the health check result
			if check.Status == StatusUnhealthy {
				LogError(checkCtx, check.Error, "health_check", map[string]interface{}{
					"service":      name,
					"response_time": check.ResponseTime,
				})
			} else {
				LogInfo(checkCtx, check.Message, "health_check", map[string]interface{}{
					"service":      name,
					"status":       check.Status,
					"response_time": check.ResponseTime,
				})
			}
		}(name, checker)
	}

	wg.Wait()
	return results
}

// CheckServiceHealth performs a health check for a specific service
func (hm *HealthManager) CheckServiceHealth(ctx context.Context, serviceName string) *HealthCheck {
	hm.mutex.RLock()
	checker, exists := hm.checkers[serviceName]
	hm.mutex.RUnlock()

	if !exists {
		return &HealthCheck{
			Name:      serviceName,
			Status:    StatusUnknown,
			Message:   fmt.Sprintf("Health checker not found for service: %s", serviceName),
			LastCheck: time.Now(),
		}
	}

	// Create a timeout context for this check
	checkCtx, cancel := context.WithTimeout(ctx, hm.timeout)
	defer cancel()

	check := checker.Check(checkCtx)

	// Log the health check result
	if check.Status == StatusUnhealthy {
		LogError(checkCtx, check.Error, "health_check", map[string]interface{}{
			"service":      serviceName,
			"response_time": check.ResponseTime,
		})
	} else {
		LogInfo(checkCtx, check.Message, "health_check", map[string]interface{}{
			"service":      serviceName,
			"status":       check.Status,
			"response_time": check.ResponseTime,
		})
	}

	return check
}

// GetOverallHealth returns the overall health status
func (hm *HealthManager) GetOverallHealth(ctx context.Context) HealthStatus {
	checks := hm.CheckHealth(ctx)
	
	if len(checks) == 0 {
		return StatusUnknown
	}

	healthyCount := 0
	unhealthyCount := 0
	degradedCount := 0

	for _, check := range checks {
		switch check.Status {
		case StatusHealthy:
			healthyCount++
		case StatusUnhealthy:
			unhealthyCount++
		case StatusDegraded:
			degradedCount++
		}
	}

	total := len(checks)
	
	// If any service is unhealthy, overall status is unhealthy
	if unhealthyCount > 0 {
		return StatusUnhealthy
	}
	
	// If any service is degraded, overall status is degraded
	if degradedCount > 0 {
		return StatusDegraded
	}
	
	// If all services are healthy, overall status is healthy
	if healthyCount == total {
		return StatusHealthy
	}
	
	return StatusUnknown
}

// StartPeriodicHealthChecks starts periodic health checks
func (hm *HealthManager) StartPeriodicHealthChecks(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			checks := hm.CheckHealth(ctx)
			
			// Log overall health status
			overallHealth := hm.GetOverallHealth(ctx)
			LogInfo(ctx, "Periodic health check completed", "periodic_health_check", map[string]interface{}{
				"overall_status": overallHealth,
				"total_services": len(checks),
			})
		}
	}
}

// DatabaseHealthChecker checks database connectivity
type DatabaseHealthChecker struct {
	name string
	ping func(ctx context.Context) error
}

// NewDatabaseHealthChecker creates a new database health checker
func NewDatabaseHealthChecker(name string, ping func(ctx context.Context) error) *DatabaseHealthChecker {
	return &DatabaseHealthChecker{
		name: name,
		ping: ping,
	}
}

// Check performs the database health check
func (dhc *DatabaseHealthChecker) Check(ctx context.Context) *HealthCheck {
	start := time.Now()
	
	check := &HealthCheck{
		Name:      dhc.name,
		Status:    StatusUnknown,
		LastCheck: time.Now(),
	}

	err := dhc.ping(ctx)
	check.ResponseTime = time.Since(start)

	if err != nil {
		check.Status = StatusUnhealthy
		check.Message = fmt.Sprintf("Database connection failed: %v", err)
		check.Error = err
	} else {
		check.Status = StatusHealthy
		check.Message = "Database connection is healthy"
	}

	return check
}

// GetName returns the health checker name
func (dhc *DatabaseHealthChecker) GetName() string {
	return dhc.name
}

// GRPCHealthChecker checks gRPC service connectivity
type GRPCHealthChecker struct {
	name    string
	check   func(ctx context.Context) error
	timeout time.Duration
}

// NewGRPCHealthChecker creates a new gRPC health checker
func NewGRPCHealthChecker(name string, check func(ctx context.Context) error, timeout time.Duration) *GRPCHealthChecker {
	if timeout == 0 {
		timeout = 3 * time.Second
	}

	return &GRPCHealthChecker{
		name:    name,
		check:   check,
		timeout: timeout,
	}
}

// Check performs the gRPC health check
func (ghc *GRPCHealthChecker) Check(ctx context.Context) *HealthCheck {
	start := time.Now()
	
	check := &HealthCheck{
		Name:      ghc.name,
		Status:    StatusUnknown,
		LastCheck: time.Now(),
	}

	// Create a timeout context for the gRPC call
	checkCtx, cancel := context.WithTimeout(ctx, ghc.timeout)
	defer cancel()

	err := ghc.check(checkCtx)
	check.ResponseTime = time.Since(start)

	if err != nil {
		check.Status = StatusUnhealthy
		check.Message = fmt.Sprintf("gRPC service check failed: %v", err)
		check.Error = err
	} else {
		check.Status = StatusHealthy
		check.Message = "gRPC service is healthy"
	}

	return check
}

// GetName returns the health checker name
func (ghc *GRPCHealthChecker) GetName() string {
	return ghc.name
}

// Global health manager instance
var DefaultHealthManager = NewHealthManager(5 * time.Second)

// RegisterDefaultHealthChecker registers a health checker with the default manager
func RegisterDefaultHealthChecker(checker HealthChecker) {
	DefaultHealthManager.RegisterHealthChecker(checker)
}

// RegisterDefaultFuncHealthChecker registers a function-based health checker with the default manager
func RegisterDefaultFuncHealthChecker(name string, fn HealthCheckFunc) {
	DefaultHealthManager.RegisterFuncHealthChecker(name, fn)
}

// CheckDefaultHealth performs health checks using the default manager
func CheckDefaultHealth(ctx context.Context) map[string]*HealthCheck {
	return DefaultHealthManager.CheckHealth(ctx)
}

// GetDefaultOverallHealth returns the overall health status using the default manager
func GetDefaultOverallHealth(ctx context.Context) HealthStatus {
	return DefaultHealthManager.GetOverallHealth(ctx)
}
