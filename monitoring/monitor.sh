#!/bin/bash

# CaloMeal Production Monitoring Script

LOG_FILE="/var/log/calomeal/monitor.log"
ALERT_EMAIL="admin@yourdomain.com"

# Create log directory if it doesn't exist
mkdir -p "$(dirname "$LOG_FILE")"

log_message() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

check_service() {
    local service_name=$1
    local check_command=$2
    
    if eval "$check_command" >/dev/null 2>&1; then
        log_message "OK: $service_name is running"
        return 0
    else
        log_message "ERROR: $service_name is not responding"
        return 1
    fi
}

# Check all services
check_service "Database" "docker exec calomeal-db-prod pg_isready -U postgres -d calomeal"
check_service "Foods Service" "nc -z localhost 50051"
check_service "Logs Service" "nc -z localhost 50052"
check_service "Analytics Service" "nc -z localhost 50053"
check_service "Backend Service" "curl -f http://localhost:8080/health"
check_service "Frontend Service" "curl -f http://localhost:80"

# Check disk space
DISK_USAGE=$(df / | awk 'NR==2 {print $5}' | sed 's/%//')
if [ "$DISK_USAGE" -gt 80 ]; then
    log_message "WARNING: Disk usage is ${DISK_USAGE}%"
fi

# Check memory usage
MEMORY_USAGE=$(free | awk 'NR==2{printf "%.0f", $3*100/$2}')
if [ "$MEMORY_USAGE" -gt 80 ]; then
    log_message "WARNING: Memory usage is ${MEMORY_USAGE}%"
fi
