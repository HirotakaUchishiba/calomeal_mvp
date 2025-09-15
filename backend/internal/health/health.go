package health

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc/health/grpc_health_v1"
)

// HealthStatus represents the health status of a service
type HealthStatus string

const (
	StatusHealthy   HealthStatus = "healthy"
	StatusUnhealthy HealthStatus = "unhealthy"
	StatusDegraded  HealthStatus = "degraded"
)

// HealthCheck represents a single health check
type HealthCheck struct {
	Name        string        `json:"name"`
	Status      HealthStatus  `json:"status"`
	Message     string        `json:"message,omitempty"`
	Duration    time.Duration `json:"duration"`
	LastChecked time.Time     `json:"last_checked"`
	Error       string        `json:"error,omitempty"`
}

// HealthResponse represents the overall health status
type HealthResponse struct {
	Status    HealthStatus  `json:"status"`
	Timestamp time.Time     `json:"timestamp"`
	Checks    []HealthCheck `json:"checks"`
	Version   string        `json:"version"`
	Uptime    time.Duration `json:"uptime"`
}

// HealthChecker defines the interface for health checking
type HealthChecker interface {
	Check(ctx context.Context) HealthCheck
}

// DatabaseHealthChecker checks database connectivity
type DatabaseHealthChecker struct {
	DB *sql.DB
}

func (h *DatabaseHealthChecker) Check(ctx context.Context) HealthCheck {
	start := time.Now()

	check := HealthCheck{
		Name:        "database",
		LastChecked: start,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := h.DB.PingContext(ctx); err != nil {
		check.Status = StatusUnhealthy
		check.Error = err.Error()
		check.Message = "Database connection failed"
	} else {
		check.Status = StatusHealthy
		check.Message = "Database connection successful"
	}

	check.Duration = time.Since(start)
	return check
}

// GRPCHealthChecker checks gRPC service connectivity
type GRPCHealthChecker struct {
	Name    string
	Address string
	Client  grpc_health_v1.HealthClient
}

func (h *GRPCHealthChecker) Check(ctx context.Context) HealthCheck {
	start := time.Now()

	check := HealthCheck{
		Name:        h.Name,
		LastChecked: start,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := h.Client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		check.Status = StatusUnhealthy
		check.Error = err.Error()
		check.Message = fmt.Sprintf("gRPC service %s is unhealthy", h.Name)
	} else {
		switch resp.Status {
		case grpc_health_v1.HealthCheckResponse_SERVING:
			check.Status = StatusHealthy
			check.Message = fmt.Sprintf("gRPC service %s is healthy", h.Name)
		case grpc_health_v1.HealthCheckResponse_NOT_SERVING:
			check.Status = StatusUnhealthy
			check.Message = fmt.Sprintf("gRPC service %s is not serving", h.Name)
		default:
			check.Status = StatusDegraded
			check.Message = fmt.Sprintf("gRPC service %s is in unknown state", h.Name)
		}
	}

	check.Duration = time.Since(start)
	return check
}

// HTTPHealthChecker checks HTTP service connectivity
type HTTPHealthChecker struct {
	Name   string
	URL    string
	Client *http.Client
}

func (h *HTTPHealthChecker) Check(ctx context.Context) HealthCheck {
	start := time.Now()

	check := HealthCheck{
		Name:        h.Name,
		LastChecked: start,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", h.URL, nil)
	if err != nil {
		check.Status = StatusUnhealthy
		check.Error = err.Error()
		check.Message = fmt.Sprintf("Failed to create request for %s", h.Name)
		check.Duration = time.Since(start)
		return check
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		check.Status = StatusUnhealthy
		check.Error = err.Error()
		check.Message = fmt.Sprintf("HTTP service %s is unhealthy", h.Name)
	} else {
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			check.Status = StatusHealthy
			check.Message = fmt.Sprintf("HTTP service %s is healthy", h.Name)
		} else {
			check.Status = StatusDegraded
			check.Message = fmt.Sprintf("HTTP service %s returned status %d", h.Name, resp.StatusCode)
		}
	}

	check.Duration = time.Since(start)
	return check
}

// HealthManager manages multiple health checks
type HealthManager struct {
	checkers  []HealthChecker
	version   string
	startTime time.Time
	mutex     sync.RWMutex
}

// NewHealthManager creates a new health manager
func NewHealthManager(version string) *HealthManager {
	return &HealthManager{
		checkers:  make([]HealthChecker, 0),
		version:   version,
		startTime: time.Now(),
	}
}

// AddChecker adds a health checker
func (hm *HealthManager) AddChecker(checker HealthChecker) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	hm.checkers = append(hm.checkers, checker)
}

// CheckAll runs all health checks
func (hm *HealthManager) CheckAll(ctx context.Context) HealthResponse {
	hm.mutex.RLock()
	checkers := make([]HealthChecker, len(hm.checkers))
	copy(checkers, hm.checkers)
	hm.mutex.RUnlock()

	response := HealthResponse{
		Status:    StatusHealthy,
		Timestamp: time.Now(),
		Checks:    make([]HealthCheck, 0, len(checkers)),
		Version:   hm.version,
		Uptime:    time.Since(hm.startTime),
	}

	// Run checks concurrently
	var wg sync.WaitGroup
	checkChan := make(chan HealthCheck, len(checkers))

	for _, checker := range checkers {
		wg.Add(1)
		go func(c HealthChecker) {
			defer wg.Done()
			checkChan <- c.Check(ctx)
		}(checker)
	}

	wg.Wait()
	close(checkChan)

	// Collect results
	for check := range checkChan {
		response.Checks = append(response.Checks, check)

		// Update overall status
		if check.Status == StatusUnhealthy {
			response.Status = StatusUnhealthy
		} else if check.Status == StatusDegraded && response.Status == StatusHealthy {
			response.Status = StatusDegraded
		}
	}

	return response
}

// GetHealthHandler returns an HTTP handler for health checks
func (hm *HealthManager) GetHealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		response := hm.CheckAll(ctx)

		w.Header().Set("Content-Type", "application/json")

		// Set appropriate HTTP status code
		switch response.Status {
		case StatusHealthy:
			w.WriteHeader(http.StatusOK)
		case StatusDegraded:
			w.WriteHeader(http.StatusOK) // Still OK, but degraded
		case StatusUnhealthy:
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		// Write JSON response
		fmt.Fprintf(w, `{
			"status": "%s",
			"timestamp": "%s",
			"version": "%s",
			"uptime": "%s",
			"checks": [`,
			response.Status,
			response.Timestamp.Format(time.RFC3339),
			response.Version,
			response.Uptime.String())

		for i, check := range response.Checks {
			if i > 0 {
				fmt.Fprint(w, ",")
			}
			fmt.Fprintf(w, `{
				"name": "%s",
				"status": "%s",
				"message": "%s",
				"duration": "%s",
				"last_checked": "%s"`,
				check.Name,
				check.Status,
				check.Message,
				check.Duration.String(),
				check.LastChecked.Format(time.RFC3339))

			if check.Error != "" {
				fmt.Fprintf(w, `,
				"error": "%s"`, check.Error)
			}
			fmt.Fprint(w, "}")
		}

		fmt.Fprint(w, "]}")
	}
}

// GetReadinessHandler returns an HTTP handler for readiness checks
func (hm *HealthManager) GetReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		response := hm.CheckAll(ctx)

		// Readiness means all critical services are healthy
		if response.Status == StatusHealthy {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "ready")
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, "not ready")
		}
	}
}

// GetLivenessHandler returns an HTTP handler for liveness checks
func (hm *HealthManager) GetLivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Liveness is simple - if the process is running, it's alive
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "alive")
	}
}
