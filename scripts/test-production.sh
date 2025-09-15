#!/bin/bash

# Production Test Script
# CaloMeal MVP - Production Environment Testing

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ENV_FILE="${PROJECT_ROOT}/config/production.env"
TEST_RESULTS_DIR="${PROJECT_ROOT}/test-results"
LOG_FILE="${TEST_RESULTS_DIR}/production-test.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0
TESTS_TOTAL=0

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1" | tee -a "$LOG_FILE"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "$LOG_FILE"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$LOG_FILE"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"
}

# Test functions
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    TESTS_TOTAL=$((TESTS_TOTAL + 1))
    log_info "Running test: $test_name"
    
    if eval "$test_command" >/dev/null 2>&1; then
        log_success "PASS: $test_name"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        log_error "FAIL: $test_name"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

# Create test results directory
setup_test_environment() {
    log_info "Setting up test environment..."
    
    mkdir -p "$TEST_RESULTS_DIR"
    rm -f "$LOG_FILE"
    
    # Load environment variables
    if [ -f "$ENV_FILE" ]; then
        source "$ENV_FILE"
    else
        log_error "Environment file not found: $ENV_FILE"
        exit 1
    fi
    
    log_success "Test environment setup completed"
}

# Test Docker services
test_docker_services() {
    log_info "Testing Docker services..."
    
    # Test if all containers are running
    run_test "Database container running" "docker ps | grep -q calomeal-db-prod"
    run_test "Foods service container running" "docker ps | grep -q calomeal-foods-prod"
    run_test "Logs service container running" "docker ps | grep -q calomeal-logs-prod"
    run_test "Analytics service container running" "docker ps | grep -q calomeal-analytics-prod"
    run_test "Backend service container running" "docker ps | grep -q calomeal-backend-prod"
    run_test "Frontend service container running" "docker ps | grep -q calomeal-frontend-prod"
    
    # Test container health
    run_test "Database container healthy" "docker inspect calomeal-db-prod | grep -q '\"Health\": {\"Status\": \"healthy\"}'"
}

# Test network connectivity
test_network_connectivity() {
    log_info "Testing network connectivity..."
    
    # Test internal network connectivity
    run_test "Database port accessible" "nc -z localhost 5432"
    run_test "Foods service port accessible" "nc -z localhost 50051"
    run_test "Logs service port accessible" "nc -z localhost 50052"
    run_test "Analytics service port accessible" "nc -z localhost 50053"
    run_test "Backend service port accessible" "nc -z localhost 8080"
    run_test "Frontend service port accessible" "nc -z localhost 80"
}

# Test database connectivity
test_database() {
    log_info "Testing database connectivity..."
    
    # Test database connection
    run_test "Database connection" "docker exec calomeal-db-prod pg_isready -U postgres -d calomeal"
    
    # Test database queries
    run_test "Database query execution" "docker exec calomeal-db-prod psql -U postgres -d calomeal -c 'SELECT 1;'"
    
    # Test database tables exist
    run_test "Users table exists" "docker exec calomeal-db-prod psql -U postgres -d calomeal -c 'SELECT 1 FROM users LIMIT 1;'"
    run_test "Foods table exists" "docker exec calomeal-db-prod psql -U postgres -d calomeal -c 'SELECT 1 FROM foods LIMIT 1;'"
    run_test "Food logs table exists" "docker exec calomeal-db-prod psql -U postgres -d calomeal -c 'SELECT 1 FROM food_logs LIMIT 1;'"
}

# Test gRPC services
test_grpc_services() {
    log_info "Testing gRPC services..."
    
    # Test gRPC health checks
    if command -v grpc_health_probe >/dev/null 2>&1; then
        run_test "Foods gRPC health check" "grpc_health_probe -addr=localhost:50051"
        run_test "Logs gRPC health check" "grpc_health_probe -addr=localhost:50052"
        run_test "Analytics gRPC health check" "grpc_health_probe -addr=localhost:50053"
    else
        log_warning "grpc_health_probe not available, skipping gRPC health checks"
    fi
}

# Test HTTP services
test_http_services() {
    log_info "Testing HTTP services..."
    
    # Test backend health endpoint
    run_test "Backend health endpoint" "curl -f http://localhost:8080/health"
    
    # Test backend GraphQL endpoint
    run_test "Backend GraphQL endpoint" "curl -f -X POST http://localhost:8080/query -H 'Content-Type: application/json' -d '{\"query\": \"{ __schema { types { name } } }\"}'"
    
    # Test frontend
    run_test "Frontend accessibility" "curl -f http://localhost:80"
    
    # Test frontend static assets
    run_test "Frontend static assets" "curl -f http://localhost:80/assets/"
}

# Test SSL/TLS
test_ssl() {
    log_info "Testing SSL/TLS..."
    
    # Test HTTPS endpoint (if nginx is running)
    if docker ps | grep -q calomeal-nginx-prod; then
        run_test "HTTPS endpoint accessible" "curl -k -f https://localhost:443"
        run_test "SSL certificate valid" "openssl s_client -connect localhost:443 -servername localhost < /dev/null 2>/dev/null | grep -q 'Verify return code: 0'"
    else
        log_warning "Nginx container not running, skipping SSL tests"
    fi
}

# Test security
test_security() {
    log_info "Testing security..."
    
    # Test security headers
    run_test "Security headers present" "curl -I http://localhost:8080/health | grep -q 'X-Content-Type-Options: nosniff'"
    
    # Test CORS headers
    run_test "CORS headers present" "curl -I -H 'Origin: http://localhost:3000' http://localhost:8080/health | grep -q 'Access-Control-Allow-Origin'"
    
    # Test rate limiting (if implemented)
    # This would require multiple rapid requests to test
    log_info "Rate limiting test would require multiple rapid requests"
}

# Test performance
test_performance() {
    log_info "Testing performance..."
    
    # Test response times
    local start_time
    local end_time
    local response_time
    
    start_time=$(date +%s%N)
    curl -f http://localhost:8080/health >/dev/null 2>&1
    end_time=$(date +%s%N)
    response_time=$(( (end_time - start_time) / 1000000 )) # Convert to milliseconds
    
    if [ $response_time -lt 1000 ]; then
        log_success "PASS: Backend response time < 1s (${response_time}ms)"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        log_error "FAIL: Backend response time > 1s (${response_time}ms)"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
    TESTS_TOTAL=$((TESTS_TOTAL + 1))
    
    # Test frontend response time
    start_time=$(date +%s%N)
    curl -f http://localhost:80 >/dev/null 2>&1
    end_time=$(date +%s%N)
    response_time=$(( (end_time - start_time) / 1000000 ))
    
    if [ $response_time -lt 2000 ]; then
        log_success "PASS: Frontend response time < 2s (${response_time}ms)"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        log_error "FAIL: Frontend response time > 2s (${response_time}ms)"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
    TESTS_TOTAL=$((TESTS_TOTAL + 1))
}

# Test logging
test_logging() {
    log_info "Testing logging..."
    
    # Test if log files are being created
    run_test "Backend logs being generated" "docker logs calomeal-backend-prod 2>&1 | grep -q 'Application started'"
    
    # Test log format (JSON in production)
    run_test "Backend logs in JSON format" "docker logs calomeal-backend-prod 2>&1 | head -1 | grep -q '{'"
}

# Test monitoring
test_monitoring() {
    log_info "Testing monitoring..."
    
    # Test if monitoring script exists and is executable
    run_test "Monitoring script exists" "[ -f ${PROJECT_ROOT}/monitoring/monitor.sh ]"
    run_test "Monitoring script executable" "[ -x ${PROJECT_ROOT}/monitoring/monitor.sh ]"
    
    # Test if backup script exists and is executable
    run_test "Backup script exists" "[ -f ${PROJECT_ROOT}/backups/backup.sh ]"
    run_test "Backup script executable" "[ -x ${PROJECT_ROOT}/backups/backup.sh ]"
}

# Test configuration
test_configuration() {
    log_info "Testing configuration..."
    
    # Test environment file
    run_test "Environment file exists" "[ -f $ENV_FILE ]"
    run_test "Environment file readable" "[ -r $ENV_FILE ]"
    
    # Test SSL certificates
    run_test "SSL certificate exists" "[ -f ${PROJECT_ROOT}/nginx/ssl/cert.pem ]"
    run_test "SSL private key exists" "[ -f ${PROJECT_ROOT}/nginx/ssl/key.pem ]"
    
    # Test Docker Compose file
    run_test "Docker Compose file exists" "[ -f ${PROJECT_ROOT}/docker-compose.prod.yml ]"
}

# Test end-to-end functionality
test_e2e_functionality() {
    log_info "Testing end-to-end functionality..."
    
    # Test GraphQL schema introspection
    local schema_response
    schema_response=$(curl -s -X POST http://localhost:8080/query \
        -H 'Content-Type: application/json' \
        -d '{"query": "{ __schema { types { name } } }"}')
    
    if echo "$schema_response" | grep -q "__schema"; then
        log_success "PASS: GraphQL schema introspection"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        log_error "FAIL: GraphQL schema introspection"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
    TESTS_TOTAL=$((TESTS_TOTAL + 1))
    
    # Test frontend HTML structure
    local frontend_response
    frontend_response=$(curl -s http://localhost:80)
    
    if echo "$frontend_response" | grep -q "<!DOCTYPE html>"; then
        log_success "PASS: Frontend HTML structure"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        log_error "FAIL: Frontend HTML structure"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
    TESTS_TOTAL=$((TESTS_TOTAL + 1))
}

# Generate test report
generate_test_report() {
    log_info "Generating test report..."
    
    local report_file="${TEST_RESULTS_DIR}/production-test-report.html"
    
    cat > "$report_file" << EOF
<!DOCTYPE html>
<html>
<head>
    <title>CaloMeal Production Test Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f0f0f0; padding: 20px; border-radius: 5px; }
        .summary { margin: 20px 0; }
        .test-result { margin: 10px 0; padding: 10px; border-radius: 3px; }
        .pass { background-color: #d4edda; color: #155724; }
        .fail { background-color: #f8d7da; color: #721c24; }
        .stats { display: flex; gap: 20px; }
        .stat { text-align: center; padding: 10px; border-radius: 5px; }
        .stat.passed { background-color: #d4edda; }
        .stat.failed { background-color: #f8d7da; }
        .stat.total { background-color: #d1ecf1; }
    </style>
</head>
<body>
    <div class="header">
        <h1>CaloMeal Production Test Report</h1>
        <p>Generated on: $(date)</p>
    </div>
    
    <div class="summary">
        <h2>Test Summary</h2>
        <div class="stats">
            <div class="stat passed">
                <h3>$TESTS_PASSED</h3>
                <p>Passed</p>
            </div>
            <div class="stat failed">
                <h3>$TESTS_FAILED</h3>
                <p>Failed</p>
            </div>
            <div class="stat total">
                <h3>$TESTS_TOTAL</h3>
                <p>Total</p>
            </div>
        </div>
    </div>
    
    <div class="test-results">
        <h2>Test Results</h2>
        <pre>$(cat "$LOG_FILE")</pre>
    </div>
</body>
</html>
EOF
    
    log_success "Test report generated: $report_file"
}

# Main test function
run_all_tests() {
    log_info "Starting production tests..."
    echo "=============================="
    
    # Setup
    setup_test_environment
    
    # Run all test suites
    test_configuration
    test_docker_services
    test_network_connectivity
    test_database
    test_grpc_services
    test_http_services
    test_ssl
    test_security
    test_performance
    test_logging
    test_monitoring
    test_e2e_functionality
    
    # Generate report
    generate_test_report
    
    # Show summary
    echo ""
    echo "=============================="
    log_info "Test Summary:"
    echo "  Total Tests: $TESTS_TOTAL"
    echo "  Passed: $TESTS_PASSED"
    echo "  Failed: $TESTS_FAILED"
    echo "  Success Rate: $(( (TESTS_PASSED * 100) / TESTS_TOTAL ))%"
    
    if [ $TESTS_FAILED -eq 0 ]; then
        log_success "All tests passed! Production environment is ready."
        return 0
    else
        log_error "Some tests failed. Please check the logs and fix the issues."
        return 1
    fi
}

# Main function
main() {
    case "${1:-all}" in
        "all")
            run_all_tests
            ;;
        "docker")
            setup_test_environment
            test_docker_services
            ;;
        "network")
            setup_test_environment
            test_network_connectivity
            ;;
        "database")
            setup_test_environment
            test_database
            ;;
        "http")
            setup_test_environment
            test_http_services
            ;;
        "security")
            setup_test_environment
            test_security
            ;;
        "performance")
            setup_test_environment
            test_performance
            ;;
        *)
            echo "Usage: $0 {all|docker|network|database|http|security|performance}"
            echo "  all        - Run all tests"
            echo "  docker     - Test Docker services"
            echo "  network    - Test network connectivity"
            echo "  database   - Test database"
            echo "  http       - Test HTTP services"
            echo "  security   - Test security"
            echo "  performance - Test performance"
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
