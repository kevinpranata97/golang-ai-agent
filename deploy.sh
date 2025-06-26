#!/bin/bash

# Golang AI Agent Deployment Script
# This script helps deploy the AI agent in various environments

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="golang-ai-agent"
DOCKER_IMAGE="golang-ai-agent"
PORT=${PORT:-8080}

# Functions
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

check_requirements() {
    log_info "Checking requirements..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed. Please install Go 1.21 or later."
        exit 1
    fi
    
    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_info "Go version: $GO_VERSION"
    
    # Check if Docker is installed (optional)
    if command -v docker &> /dev/null; then
        log_info "Docker is available"
        DOCKER_AVAILABLE=true
    else
        log_warning "Docker is not installed. Docker deployment will not be available."
        DOCKER_AVAILABLE=false
    fi
    
    # Check if git is installed
    if ! command -v git &> /dev/null; then
        log_error "Git is not installed. Please install Git."
        exit 1
    fi
}

build_application() {
    log_info "Building application..."
    
    # Clean previous builds
    rm -f $APP_NAME
    
    # Build the application
    go mod tidy
    go build -v -ldflags="-s -w" -o $APP_NAME .
    
    if [ -f "$APP_NAME" ]; then
        log_success "Application built successfully"
        chmod +x $APP_NAME
    else
        log_error "Failed to build application"
        exit 1
    fi
}

run_tests() {
    log_info "Running tests..."
    
    # Run unit tests
    if go test -v ./...; then
        log_success "All tests passed"
    else
        log_warning "Some tests failed, but continuing deployment"
    fi
    
    # Run static analysis
    log_info "Running static analysis..."
    go vet ./...
    
    # Check formatting
    if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
        log_warning "Code formatting issues found:"
        gofmt -s -l .
        log_info "Run 'go fmt ./...' to fix formatting"
    fi
}

deploy_local() {
    log_info "Deploying locally..."
    
    # Create necessary directories
    mkdir -p data generated_apps logs
    
    # Set environment variables
    export PORT=$PORT
    export STORAGE_DIR="./data"
    
    log_info "Starting application on port $PORT..."
    log_info "Available endpoints:"
    log_info "  - Health: http://localhost:$PORT/health"
    log_info "  - Status: http://localhost:$PORT/status"
    log_info "  - Generate App: POST http://localhost:$PORT/generate-app"
    log_info "  - Test App: POST http://localhost:$PORT/test-app"
    log_info "  - Generate & Test: POST http://localhost:$PORT/generate-and-test"
    log_info "  - Webhook: POST http://localhost:$PORT/webhook"
    
    # Start the application
    ./$APP_NAME
}

deploy_docker() {
    if [ "$DOCKER_AVAILABLE" = false ]; then
        log_error "Docker is not available"
        exit 1
    fi
    
    log_info "Building Docker image..."
    
    # Build Docker image
    docker build -t $DOCKER_IMAGE .
    
    if [ $? -eq 0 ]; then
        log_success "Docker image built successfully"
    else
        log_error "Failed to build Docker image"
        exit 1
    fi
    
    log_info "Running Docker container..."
    
    # Create necessary directories for volume mounts
    mkdir -p data generated_apps
    
    # Run Docker container
    docker run -d \
        --name $APP_NAME \
        -p $PORT:8080 \
        -e PORT=8080 \
        -e STORAGE_DIR=/app/data \
        -v $(pwd)/data:/app/data \
        -v $(pwd)/generated_apps:/app/generated_apps \
        --restart unless-stopped \
        $DOCKER_IMAGE
    
    if [ $? -eq 0 ]; then
        log_success "Docker container started successfully"
        log_info "Container name: $APP_NAME"
        log_info "Port: $PORT"
        log_info "Health check: curl http://localhost:$PORT/health"
    else
        log_error "Failed to start Docker container"
        exit 1
    fi
}

deploy_docker_compose() {
    if [ "$DOCKER_AVAILABLE" = false ]; then
        log_error "Docker is not available"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed"
        exit 1
    fi
    
    log_info "Deploying with Docker Compose..."
    
    # Create environment file if it doesn't exist
    if [ ! -f ".env" ]; then
        log_info "Creating .env file..."
        cat > .env << EOF
PORT=8080
GEMINI_API_KEY=your_gemini_api_key_here
GITHUB_TOKEN=your_github_token_here
STORAGE_DIR=/app/data
EOF
        log_warning "Please update .env file with your API keys"
    fi
    
    # Start services
    docker-compose up -d
    
    if [ $? -eq 0 ]; then
        log_success "Services started successfully"
        log_info "Application: http://localhost:$PORT"
        log_info "Nginx proxy: http://localhost:80"
        log_info "View logs: docker-compose logs -f"
        log_info "Stop services: docker-compose down"
    else
        log_error "Failed to start services"
        exit 1
    fi
}

cleanup() {
    log_info "Cleaning up..."
    
    # Stop and remove Docker containers
    if [ "$DOCKER_AVAILABLE" = true ]; then
        docker stop $APP_NAME 2>/dev/null || true
        docker rm $APP_NAME 2>/dev/null || true
        docker-compose down 2>/dev/null || true
    fi
    
    # Clean build artifacts
    rm -f $APP_NAME
    
    log_success "Cleanup completed"
}

show_help() {
    echo "Golang AI Agent Deployment Script"
    echo ""
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  build           Build the application"
    echo "  test            Run tests and static analysis"
    echo "  local           Deploy locally (default)"
    echo "  docker          Deploy using Docker"
    echo "  compose         Deploy using Docker Compose"
    echo "  cleanup         Clean up containers and build artifacts"
    echo "  help            Show this help message"
    echo ""
    echo "Options:"
    echo "  PORT=8080       Set the port number (default: 8080)"
    echo ""
    echo "Examples:"
    echo "  $0 local                    # Deploy locally"
    echo "  PORT=9000 $0 local         # Deploy locally on port 9000"
    echo "  $0 docker                  # Deploy using Docker"
    echo "  $0 compose                 # Deploy using Docker Compose"
    echo "  $0 cleanup                 # Clean up everything"
}

# Main script logic
case "${1:-local}" in
    "build")
        check_requirements
        build_application
        ;;
    "test")
        check_requirements
        run_tests
        ;;
    "local")
        check_requirements
        build_application
        run_tests
        deploy_local
        ;;
    "docker")
        check_requirements
        deploy_docker
        ;;
    "compose")
        check_requirements
        deploy_docker_compose
        ;;
    "cleanup")
        cleanup
        ;;
    "help"|"-h"|"--help")
        show_help
        ;;
    *)
        log_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac

