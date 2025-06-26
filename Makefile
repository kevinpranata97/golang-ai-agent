# Golang AI Agent Makefile

# Variables
BINARY_NAME=golang-ai-agent
DOCKER_IMAGE=golang-ai-agent
VERSION?=latest
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# Default target
.PHONY: all
all: clean test build

# Build the application
.PHONY: build
build:
	@echo "Building ${BINARY_NAME}..."
	go build ${LDFLAGS} -o ${BINARY_NAME} .

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-windows-amd64.exe .

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage: test
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Lint code
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -f ${BINARY_NAME}
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run the application
.PHONY: run
run: build
	@echo "Running ${BINARY_NAME}..."
	./${BINARY_NAME}

# Run with development settings
.PHONY: dev
dev:
	@echo "Running in development mode..."
	export LOG_LEVEL=debug && ./${BINARY_NAME}

# Docker build
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t ${DOCKER_IMAGE}:${VERSION} .

# Docker run
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -d \
		-p 8080:8080 \
		-e GITHUB_TOKEN="${GITHUB_TOKEN}" \
		-e WEBHOOK_SECRET="${WEBHOOK_SECRET}" \
		-v $(PWD)/data:/root/data \
		${DOCKER_IMAGE}:${VERSION}

# Docker push
.PHONY: docker-push
docker-push:
	@echo "Pushing Docker image..."
	docker push ${DOCKER_IMAGE}:${VERSION}

# Setup development environment
.PHONY: setup
setup:
	@echo "Setting up development environment..."
	go mod download
	cp config.json.example config.json
	mkdir -p data
	@echo "Setup complete! Edit config.json with your settings."

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Security scan
.PHONY: security
security:
	@echo "Running security scan..."
	gosec ./...

# Generate documentation
.PHONY: docs
docs:
	@echo "Generating documentation..."
	godoc -http=:6060

# Health check
.PHONY: health
health:
	@echo "Checking application health..."
	curl -f http://localhost:8080/health || echo "Application not running"

# Status check
.PHONY: status
status:
	@echo "Checking application status..."
	curl -s http://localhost:8080/status | jq .

# Deploy to production
.PHONY: deploy
deploy: test build docker-build
	@echo "Deploying to production..."
	# Add your deployment commands here

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean, test, and build"
	@echo "  build        - Build the application"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  bench        - Run benchmarks"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  tidy         - Tidy dependencies"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Build and run the application"
	@echo "  dev          - Run in development mode"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  docker-push  - Push Docker image"
	@echo "  setup        - Setup development environment"
	@echo "  deps         - Install dependencies"
	@echo "  security     - Run security scan"
	@echo "  docs         - Generate documentation"
	@echo "  health       - Check application health"
	@echo "  status       - Check application status"
	@echo "  deploy       - Deploy to production"
	@echo "  help         - Show this help message"

