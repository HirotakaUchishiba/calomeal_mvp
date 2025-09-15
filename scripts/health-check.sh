#!/bin/bash

# Production Health Check Script
# CaloMeal MVP - Comprehensive Health Monitoring

set -e

# Configuration
HEALTH_CHECK_TIMEOUT=${HEALTH_CHECK_TIMEOUT:-30}
HEALTH_CHECK_INTERVAL=${HEALTH_CHECK_INTERVAL:-10}
MAX_RETRIES=${MAX_RETRIES:-3}

# Service URLs
DB_URL="postgres://postgres:${DB_PASSWORD:-password}@localhost:5432/calomeal?sslmode=disable"
FOODS_URL="localhost:50051"
LOGS_URL="localhost:50052"
ANALYTICS_URL="localhost:50053"
BACKEND_URL="http://localhost:8080"
FRONTEND_URL="http://localhost:80"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Health check functions
check_database() {
    log_info "Checking database connection..."
    
    if command -v psql >/dev/null 2>&1; then
        if psql "$DB_URL" -c "SELECT 1;" >/dev/null 2>&1; then
            log_success "Database connection: OK"
            return 0
        else
            log_error "Database connection: FAILED"
            return 1
        fi
    else
        # Use Docker exec if psql is not available
        if docker exec calomeal-db-prod psql -U postgres -d calomeal -c "SELECT 1;" >/dev/null 2>&1; then
            log_success "Database connection: OK"
            return 0
        else
            log_error "Database connection: FAILED"
            return 1
        fi
    fi
}

check_grpc_service() {
    local service_name=$1
    local service_url=$2
    
    log_info "Checking $service_name gRPC service..."
    
    if command -v grpc_health_probe >/dev/null 2>&1; then
        if grpc_health_probe -addr="$service_url" -timeout="${HEALTH_CHECK_TIMEOUT}s" >/dev/null 2>&1; then
            log_success "$service_name service: OK"
            return 0
        else
            log_error "$service_name service: FAILED"
            return 1
        fi
    else
        # Fallback: try to connect using netcat
        if timeout 5 bash -c "</dev/tcp/${service_url/:/ }" >/dev/null 2>&1; then
            log_success "$service_name service: OK (port open)"
            return 0
        else
            log_error "$service_name service: FAILED (port not accessible)"
            return 1
        fi
    fi
}

check_http_service() {
    local service_name=$1
    local service_url=$2
    local endpoint=${3:-"/health"}
    
    log_info "Checking $service_name HTTP service..."
    
    if curl -f -s --max-time "$HEALTH_CHECK_TIMEOUT" "$service_url$endpoint" >/dev/null 2>&1; then
        log_success "$service_name service: OK"
        return 0
    else
        log_error "$service_name service: FAILED"
        return 1
    fi
}

check_graphql_schema() {
    log_info "Checking GraphQL schema..."
    
    local response
    response=$(curl -s --max-time "$HEALTH_CHECK_TIMEOUT" \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{"query": "{ __schema { types { name } } }"}' \
        "$BACKEND_URL/query" 2>/dev/null)
    
    if echo "$response" | grep -q "__schema"; then
        log_success "GraphQL schema: OK"
        return 0
    else
        log_error "GraphQL schema: FAILED"
        return 1
    fi
}

check_frontend_assets() {
    log_info "Checking frontend assets..."
    
    if curl -f -s --max-time "$HEALTH_CHECK_TIMEOUT" "$FRONTEND_URL" >/dev/null 2>&1; then
        log_success "Frontend assets: OK"
        return 0
    else
        log_error "Frontend assets: FAILED"
        return 1
    fi
}

# Comprehensive health check
run_health_check() {
    local overall_status=0
    local failed_services=()
    
    log_info "Starting comprehensive health check..."
    echo "=========================================="
    
    # Database check
    if ! check_database; then
        overall_status=1
        failed_services+=("database")
    fi
    
    # gRPC services check
    if ! check_grpc_service "Foods" "$FOODS_URL"; then
        overall_status=1
        failed_services+=("foods")
    fi
    
    if ! check_grpc_service "Logs" "$LOGS_URL"; then
        overall_status=1
        failed_services+=("logs")
    fi
    
    if ! check_grpc_service "Analytics" "$ANALYTICS_URL"; then
        overall_status=1
        failed_services+=("analytics")
    fi
    
    # Backend service check
    if ! check_http_service "Backend" "$BACKEND_URL"; then
        overall_status=1
        failed_services+=("backend")
    fi
    
    # GraphQL schema check
    if ! check_graphql_schema; then
        overall_status=1
        failed_services+=("graphql")
    fi
    
    # Frontend check
    if ! check_frontend_assets; then
        overall_status=1
        failed_services+=("frontend")
    fi
    
    echo "=========================================="
    
    if [ $overall_status -eq 0 ]; then
        log_success "All services are healthy!"
        return 0
    else
        log_error "Health check failed for: ${failed_services[*]}"
        return 1
    fi
}

# Continuous monitoring
monitor_services() {
    log_info "Starting continuous monitoring (interval: ${HEALTH_CHECK_INTERVAL}s)"
    
    while true; do
        echo ""
        echo "$(date): Running health check..."
        
        if ! run_health_check; then
            log_warning "Health check failed, will retry in ${HEALTH_CHECK_INTERVAL}s"
        fi
        
        sleep "$HEALTH_CHECK_INTERVAL"
    done
}

# Wait for services to be ready
wait_for_services() {
    log_info "Waiting for all services to be ready..."
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        log_info "Attempt $attempt/$max_attempts"
        
        if run_health_check >/dev/null 2>&1; then
            log_success "All services are ready!"
            return 0
        fi
        
        log_warning "Services not ready yet, waiting ${HEALTH_CHECK_INTERVAL}s..."
        sleep "$HEALTH_CHECK_INTERVAL"
        attempt=$((attempt + 1))
    done
    
    log_error "Services failed to become ready within timeout"
    return 1
}

# Main function
main() {
    case "${1:-check}" in
        "check")
            run_health_check
            ;;
        "monitor")
            monitor_services
            ;;
        "wait")
            wait_for_services
            ;;
        *)
            echo "Usage: $0 {check|monitor|wait}"
            echo "  check   - Run a single health check"
            echo "  monitor - Run continuous monitoring"
            echo "  wait    - Wait for services to be ready"
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
