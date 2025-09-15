#!/bin/bash

# Production Deployment Script
# CaloMeal MVP - Production Deployment Automation

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ENV_FILE="${PROJECT_ROOT}/config/production.env"
BACKUP_DIR="${PROJECT_ROOT}/backups"
LOG_FILE="${PROJECT_ROOT}/deployment.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Error handling
error_exit() {
    log_error "$1"
    exit 1
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if Docker is installed and running
    if ! command -v docker >/dev/null 2>&1; then
        error_exit "Docker is not installed"
    fi
    
    if ! docker info >/dev/null 2>&1; then
        error_exit "Docker is not running"
    fi
    
    # Check if Docker Compose is installed
    if ! command -v docker-compose >/dev/null 2>&1; then
        error_exit "Docker Compose is not installed"
    fi
    
    # Check if environment file exists
    if [ ! -f "$ENV_FILE" ]; then
        error_exit "Environment file not found: $ENV_FILE"
    fi
    
    # Check if required environment variables are set
    source "$ENV_FILE"
    if [ -z "$DB_PASSWORD" ] || [ -z "$JWT_SECRET" ]; then
        error_exit "Required environment variables not set in $ENV_FILE"
    fi
    
    log_success "Prerequisites check passed"
}

# Create backup
create_backup() {
    log_info "Creating backup..."
    
    # Create backup directory
    mkdir -p "$BACKUP_DIR"
    
    # Create timestamp
    TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
    BACKUP_PATH="${BACKUP_DIR}/backup_${TIMESTAMP}"
    
    # Backup database
    if docker ps | grep -q calomeal-db-prod; then
        log_info "Backing up database..."
        docker exec calomeal-db-prod pg_dump -U postgres calomeal > "${BACKUP_PATH}_database.sql"
        log_success "Database backup created: ${BACKUP_PATH}_database.sql"
    fi
    
    # Backup configuration files
    log_info "Backing up configuration files..."
    cp -r "$PROJECT_ROOT/config" "${BACKUP_PATH}_config"
    log_success "Configuration backup created: ${BACKUP_PATH}_config"
    
    # Keep only last 5 backups
    cd "$BACKUP_DIR"
    ls -t backup_* | tail -n +6 | xargs -r rm -rf
    cd "$PROJECT_ROOT"
    
    log_success "Backup completed"
}

# Build and deploy services
deploy_services() {
    log_info "Building and deploying services..."
    
    # Load environment variables
    source "$ENV_FILE"
    
    # Stop existing services
    log_info "Stopping existing services..."
    docker-compose -f docker-compose.prod.yml down --remove-orphans || true
    
    # Build and start services
    log_info "Building Docker images..."
    docker-compose -f docker-compose.prod.yml build --no-cache
    
    log_info "Starting services..."
    docker-compose -f docker-compose.prod.yml up -d
    
    log_success "Services deployed"
}

# Wait for services to be ready
wait_for_services() {
    log_info "Waiting for services to be ready..."
    
    # Wait for database
    log_info "Waiting for database..."
    timeout 60 bash -c 'until docker exec calomeal-db-prod pg_isready -U postgres -d calomeal; do sleep 2; done' || error_exit "Database failed to start"
    
    # Wait for gRPC services
    log_info "Waiting for gRPC services..."
    timeout 60 bash -c 'until nc -z localhost 50051; do sleep 2; done' || error_exit "Foods service failed to start"
    timeout 60 bash -c 'until nc -z localhost 50052; do sleep 2; done' || error_exit "Logs service failed to start"
    timeout 60 bash -c 'until nc -z localhost 50053; do sleep 2; done' || error_exit "Analytics service failed to start"
    
    # Wait for backend
    log_info "Waiting for backend service..."
    timeout 60 bash -c 'until curl -f http://localhost:8080/health >/dev/null 2>&1; do sleep 2; done' || error_exit "Backend service failed to start"
    
    # Wait for frontend
    log_info "Waiting for frontend service..."
    timeout 60 bash -c 'until curl -f http://localhost:80 >/dev/null 2>&1; do sleep 2; done' || error_exit "Frontend service failed to start"
    
    log_success "All services are ready"
}

# Run health checks
run_health_checks() {
    log_info "Running health checks..."
    
    # Run comprehensive health check
    if [ -f "${PROJECT_ROOT}/scripts/health-check.sh" ]; then
        chmod +x "${PROJECT_ROOT}/scripts/health-check.sh"
        if "${PROJECT_ROOT}/scripts/health-check.sh" check; then
            log_success "Health checks passed"
        else
            error_exit "Health checks failed"
        fi
    else
        log_warning "Health check script not found, skipping health checks"
    fi
}

# Run database migrations
run_migrations() {
    log_info "Running database migrations..."
    
    # Wait for database to be ready
    timeout 30 bash -c 'until docker exec calomeal-db-prod pg_isready -U postgres -d calomeal; do sleep 1; done' || error_exit "Database not ready for migrations"
    
    # Run migrations
    if [ -d "${PROJECT_ROOT}/database/migrations" ]; then
        for migration in "${PROJECT_ROOT}/database/migrations"/*.sql; do
            if [ -f "$migration" ]; then
                log_info "Running migration: $(basename "$migration")"
                docker exec -i calomeal-db-prod psql -U postgres -d calomeal < "$migration" || error_exit "Migration failed: $(basename "$migration")"
            fi
        done
        log_success "Database migrations completed"
    else
        log_warning "No migrations found"
    fi
}

# Setup SSL certificates
setup_ssl() {
    log_info "Setting up SSL certificates..."
    
    SSL_DIR="${PROJECT_ROOT}/nginx/ssl"
    mkdir -p "$SSL_DIR"
    
    # Check if certificates exist
    if [ ! -f "$SSL_DIR/cert.pem" ] || [ ! -f "$SSL_DIR/key.pem" ]; then
        log_warning "SSL certificates not found, generating self-signed certificates..."
        
        # Generate self-signed certificate (for development/testing only)
        openssl req -x509 -newkey rsa:4096 -keyout "$SSL_DIR/key.pem" -out "$SSL_DIR/cert.pem" -days 365 -nodes -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost" || error_exit "Failed to generate SSL certificates"
        
        log_warning "Self-signed certificates generated. For production, use proper SSL certificates."
    else
        log_success "SSL certificates found"
    fi
}

# Configure monitoring
setup_monitoring() {
    log_info "Setting up monitoring..."
    
    # Create monitoring directory
    MONITORING_DIR="${PROJECT_ROOT}/monitoring"
    mkdir -p "$MONITORING_DIR"
    
    # Create basic monitoring script
    cat > "${MONITORING_DIR}/monitor.sh" << 'EOF'
#!/bin/bash
# Basic monitoring script

while true; do
    echo "$(date): Checking services..."
    
    # Check database
    if ! docker exec calomeal-db-prod pg_isready -U postgres -d calomeal >/dev/null 2>&1; then
        echo "ERROR: Database is not responding"
    fi
    
    # Check gRPC services
    for port in 50051 50052 50053; do
        if ! nc -z localhost $port; then
            echo "ERROR: Service on port $port is not responding"
        fi
    done
    
    # Check HTTP services
    if ! curl -f http://localhost:8080/health >/dev/null 2>&1; then
        echo "ERROR: Backend service is not responding"
    fi
    
    if ! curl -f http://localhost:80 >/dev/null 2>&1; then
        echo "ERROR: Frontend service is not responding"
    fi
    
    sleep 60
done
EOF
    
    chmod +x "${MONITORING_DIR}/monitor.sh"
    log_success "Monitoring setup completed"
}

# Cleanup old resources
cleanup() {
    log_info "Cleaning up old resources..."
    
    # Remove unused Docker images
    docker image prune -f
    
    # Remove unused Docker volumes
    docker volume prune -f
    
    # Remove unused Docker networks
    docker network prune -f
    
    log_success "Cleanup completed"
}

# Rollback deployment
rollback() {
    log_info "Rolling back deployment..."
    
    # Stop current services
    docker-compose -f docker-compose.prod.yml down
    
    # Find latest backup
    LATEST_BACKUP=$(ls -t "$BACKUP_DIR"/backup_*_database.sql 2>/dev/null | head -n1)
    
    if [ -n "$LATEST_BACKUP" ]; then
        log_info "Restoring database from backup: $LATEST_BACKUP"
        
        # Start database only
        docker-compose -f docker-compose.prod.yml up -d db
        
        # Wait for database
        timeout 30 bash -c 'until docker exec calomeal-db-prod pg_isready -U postgres -d calomeal; do sleep 1; done'
        
        # Restore database
        docker exec -i calomeal-db-prod psql -U postgres -d calomeal < "$LATEST_BACKUP"
        
        log_success "Database restored from backup"
    else
        log_warning "No backup found for rollback"
    fi
    
    log_success "Rollback completed"
}

# Show deployment status
show_status() {
    log_info "Deployment Status:"
    echo "=================="
    
    # Show running containers
    echo "Running Containers:"
    docker-compose -f docker-compose.prod.yml ps
    
    echo ""
    echo "Service URLs:"
    echo "  Frontend: http://localhost:80"
    echo "  Backend:  http://localhost:8080"
    echo "  Health:   http://localhost:8080/health"
    
    echo ""
    echo "Logs:"
    echo "  View logs: docker-compose -f docker-compose.prod.yml logs -f"
    echo "  Monitor:   ${PROJECT_ROOT}/monitoring/monitor.sh"
}

# Main deployment function
deploy() {
    log_info "Starting production deployment..."
    echo "=================================="
    
    # Pre-deployment checks
    check_prerequisites
    
    # Create backup
    create_backup
    
    # Setup SSL
    setup_ssl
    
    # Deploy services
    deploy_services
    
    # Wait for services
    wait_for_services
    
    # Run migrations
    run_migrations
    
    # Run health checks
    run_health_checks
    
    # Setup monitoring
    setup_monitoring
    
    # Cleanup
    cleanup
    
    # Show status
    show_status
    
    log_success "Production deployment completed successfully!"
}

# Main function
main() {
    case "${1:-deploy}" in
        "deploy")
            deploy
            ;;
        "rollback")
            rollback
            ;;
        "status")
            show_status
            ;;
        "health")
            run_health_checks
            ;;
        "backup")
            create_backup
            ;;
        *)
            echo "Usage: $0 {deploy|rollback|status|health|backup}"
            echo "  deploy   - Deploy to production"
            echo "  rollback - Rollback to previous version"
            echo "  status   - Show deployment status"
            echo "  health   - Run health checks"
            echo "  backup   - Create backup"
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
