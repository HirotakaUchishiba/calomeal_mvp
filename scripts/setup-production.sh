#!/bin/bash

# Production Setup Script
# CaloMeal MVP - Initial Production Environment Setup

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ENV_FILE="${PROJECT_ROOT}/config/production.env"

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

# Error handling
error_exit() {
    log_error "$1"
    exit 1
}

# Generate secure password
generate_password() {
    openssl rand -base64 32 | tr -d "=+/" | cut -c1-25
}

# Generate JWT secret
generate_jwt_secret() {
    openssl rand -base64 64 | tr -d "=+/" | cut -c1-50
}

# Setup environment file
setup_environment() {
    log_info "Setting up production environment file..."
    
    if [ -f "$ENV_FILE" ]; then
        log_warning "Environment file already exists. Creating backup..."
        cp "$ENV_FILE" "${ENV_FILE}.backup.$(date +%Y%m%d_%H%M%S)"
    fi
    
    # Generate secure values
    DB_PASSWORD=$(generate_password)
    JWT_SECRET=$(generate_jwt_secret)
    
    # Create environment file
    cat > "$ENV_FILE" << EOF
# Production Environment Configuration
# CaloMeal MVP - Production Settings
# Generated on $(date)

# Database Configuration
DATABASE_URL=postgres://postgres:${DB_PASSWORD}@localhost:5432/calomeal?sslmode=require
DB_HOST=localhost
DB_PORT=5432
DB_NAME=calomeal
DB_USER=postgres
DB_PASSWORD=${DB_PASSWORD}

# Service Addresses (Production)
FOOD_SERVICE_ADDR=localhost:50051
LOGS_SERVICE_ADDR=localhost:50052
ANALYTICS_SERVICE_ADDR=localhost:50053

# BFF Configuration
BFF_PORT=8080
BFF_HOST=0.0.0.0

# Frontend Configuration
FRONTEND_PORT=5173
FRONTEND_HOST=0.0.0.0

# Security Configuration
JWT_SECRET=${JWT_SECRET}
CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json

# Performance Configuration
MAX_CONNECTIONS=100
CONNECTION_TIMEOUT=30s
REQUEST_TIMEOUT=60s

# Health Check Configuration
HEALTH_CHECK_INTERVAL=30s
HEALTH_CHECK_TIMEOUT=10s

# Monitoring Configuration
METRICS_ENABLED=true
METRICS_PORT=9090

# Error Handling Configuration
RETRY_MAX_ATTEMPTS=3
RETRY_INITIAL_DELAY=100ms
CIRCUIT_BREAKER_FAILURE_THRESHOLD=5
CIRCUIT_BREAKER_RESET_TIMEOUT=30s
EOF
    
    log_success "Environment file created: $ENV_FILE"
    log_warning "IMPORTANT: Update CORS_ORIGINS with your actual domain names!"
}

# Setup SSL certificates
setup_ssl_certificates() {
    log_info "Setting up SSL certificates..."
    
    SSL_DIR="${PROJECT_ROOT}/nginx/ssl"
    mkdir -p "$SSL_DIR"
    
    # Check if certificates already exist
    if [ -f "$SSL_DIR/cert.pem" ] && [ -f "$SSL_DIR/key.pem" ]; then
        log_warning "SSL certificates already exist. Skipping generation."
        return
    fi
    
    # Generate self-signed certificate for development
    log_info "Generating self-signed SSL certificate..."
    
    openssl req -x509 -newkey rsa:4096 -keyout "$SSL_DIR/key.pem" -out "$SSL_DIR/cert.pem" -days 365 -nodes \
        -subj "/C=US/ST=State/L=City/O=CaloMeal/CN=localhost" \
        -addext "subjectAltName=DNS:localhost,DNS:*.localhost,IP:127.0.0.1" || error_exit "Failed to generate SSL certificate"
    
    # Set proper permissions
    chmod 600 "$SSL_DIR/key.pem"
    chmod 644 "$SSL_DIR/cert.pem"
    
    log_success "SSL certificates generated in $SSL_DIR"
    log_warning "For production, replace with proper SSL certificates from a trusted CA"
}

# Setup directories
setup_directories() {
    log_info "Setting up production directories..."
    
    # Create necessary directories
    mkdir -p "${PROJECT_ROOT}/backups"
    mkdir -p "${PROJECT_ROOT}/logs"
    mkdir -p "${PROJECT_ROOT}/monitoring"
    mkdir -p "${PROJECT_ROOT}/nginx/ssl"
    
    # Set proper permissions
    chmod 755 "${PROJECT_ROOT}/backups"
    chmod 755 "${PROJECT_ROOT}/logs"
    chmod 755 "${PROJECT_ROOT}/monitoring"
    chmod 700 "${PROJECT_ROOT}/nginx/ssl"
    
    log_success "Production directories created"
}

# Setup systemd services (for Linux)
setup_systemd_services() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        log_info "Setting up systemd services..."
        
        # Create systemd service file
        cat > /tmp/calomeal.service << EOF
[Unit]
Description=CaloMeal MVP Production Services
After=docker.service
Requires=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=${PROJECT_ROOT}
ExecStart=/usr/bin/docker-compose -f docker-compose.prod.yml up -d
ExecStop=/usr/bin/docker-compose -f docker-compose.prod.yml down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
EOF
        
        log_info "Systemd service file created at /tmp/calomeal.service"
        log_warning "To install: sudo cp /tmp/calomeal.service /etc/systemd/system/ && sudo systemctl enable calomeal"
    else
        log_info "Skipping systemd setup (not on Linux)"
    fi
}

# Setup monitoring
setup_monitoring() {
    log_info "Setting up monitoring..."
    
    MONITORING_DIR="${PROJECT_ROOT}/monitoring"
    
    # Create monitoring script
    cat > "${MONITORING_DIR}/monitor.sh" << 'EOF'
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
EOF
    
    chmod +x "${MONITORING_DIR}/monitor.sh"
    
    # Create logrotate configuration
    cat > "${MONITORING_DIR}/logrotate.conf" << EOF
/var/log/calomeal/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 root root
}
EOF
    
    log_success "Monitoring setup completed"
}

# Setup backup script
setup_backup() {
    log_info "Setting up backup script..."
    
    BACKUP_DIR="${PROJECT_ROOT}/backups"
    
    cat > "${BACKUP_DIR}/backup.sh" << 'EOF'
#!/bin/bash

# CaloMeal Production Backup Script

BACKUP_DIR="/opt/calomeal/backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/calomeal_backup_${TIMESTAMP}.sql"

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Backup database
docker exec calomeal-db-prod pg_dump -U postgres calomeal > "$BACKUP_FILE"

# Compress backup
gzip "$BACKUP_FILE"

# Keep only last 7 days of backups
find "$BACKUP_DIR" -name "calomeal_backup_*.sql.gz" -mtime +7 -delete

echo "Backup completed: ${BACKUP_FILE}.gz"
EOF
    
    chmod +x "${BACKUP_DIR}/backup.sh"
    
    log_success "Backup script created"
}

# Setup firewall rules (for Linux)
setup_firewall() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        log_info "Setting up firewall rules..."
        
        # Check if ufw is available
        if command -v ufw >/dev/null 2>&1; then
            log_info "Configuring UFW firewall..."
            
            # Allow SSH
            sudo ufw allow ssh
            
            # Allow HTTP and HTTPS
            sudo ufw allow 80/tcp
            sudo ufw allow 443/tcp
            
            # Allow internal gRPC ports (only from localhost)
            sudo ufw allow from 127.0.0.1 to any port 50051
            sudo ufw allow from 127.0.0.1 to any port 50052
            sudo ufw allow from 127.0.0.1 to any port 50053
            sudo ufw allow from 127.0.0.1 to any port 8080
            
            log_success "Firewall rules configured"
        else
            log_warning "UFW not available, skipping firewall setup"
        fi
    else
        log_info "Skipping firewall setup (not on Linux)"
    fi
}

# Validate setup
validate_setup() {
    log_info "Validating production setup..."
    
    # Check environment file
    if [ ! -f "$ENV_FILE" ]; then
        error_exit "Environment file not found: $ENV_FILE"
    fi
    
    # Check SSL certificates
    if [ ! -f "${PROJECT_ROOT}/nginx/ssl/cert.pem" ] || [ ! -f "${PROJECT_ROOT}/nginx/ssl/key.pem" ]; then
        error_exit "SSL certificates not found"
    fi
    
    # Check Docker
    if ! command -v docker >/dev/null 2>&1; then
        error_exit "Docker not installed"
    fi
    
    if ! command -v docker-compose >/dev/null 2>&1; then
        error_exit "Docker Compose not installed"
    fi
    
    log_success "Production setup validation passed"
}

# Show setup summary
show_summary() {
    log_success "Production setup completed!"
    echo "=================================="
    echo ""
    echo "Next steps:"
    echo "1. Update CORS_ORIGINS in $ENV_FILE with your actual domain names"
    echo "2. Replace self-signed SSL certificates with proper certificates"
    echo "3. Run deployment: ./scripts/deploy-production.sh"
    echo "4. Setup monitoring: ${PROJECT_ROOT}/monitoring/monitor.sh"
    echo "5. Setup automated backups: ${PROJECT_ROOT}/backups/backup.sh"
    echo ""
    echo "Important files:"
    echo "  Environment: $ENV_FILE"
    echo "  SSL Certificates: ${PROJECT_ROOT}/nginx/ssl/"
    echo "  Monitoring: ${PROJECT_ROOT}/monitoring/"
    echo "  Backups: ${PROJECT_ROOT}/backups/"
    echo ""
    echo "Security notes:"
    echo "  - Change default passwords in production"
    echo "  - Use proper SSL certificates"
    echo "  - Configure firewall rules"
    echo "  - Enable log monitoring"
    echo "  - Setup automated backups"
}

# Main function
main() {
    log_info "Starting production setup..."
    echo "=============================="
    
    # Setup steps
    setup_directories
    setup_environment
    setup_ssl_certificates
    setup_monitoring
    setup_backup
    setup_systemd_services
    setup_firewall
    
    # Validate setup
    validate_setup
    
    # Show summary
    show_summary
}

# Run main function
main "$@"
