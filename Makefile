.PHONY: help test test-unit test-integration build build-frontend build-backend docker-build docker-run clean deps lint

# Default target
help:
	@echo "Available targets:"
	@echo "  make test              - Run all tests"
	@echo "  make test-unit         - Run unit tests only"
	@echo "  make test-integration  - Run integration tests only"
	@echo "  make build             - Build backend binary"
	@echo "  make build-frontend    - Build frontend"
	@echo "  make build-backend     - Build backend"
	@echo "  make docker-build      - Build Docker image"
	@echo "  make docker-run        - Run Docker container locally"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make deps              - Install dependencies"
	@echo "  make lint              - Run linters"

# Variables
BACKEND_DIR := backend
FRONTEND_DIR := frontend
DOCKER_IMAGE := appdirect-workshop
DOCKER_TAG := latest

# Test targets
test:
	@echo "Running all tests..."
	cd $(BACKEND_DIR) && go test ./... -v

test-unit:
	@echo "Running unit tests..."
	cd $(BACKEND_DIR) && go test ./internal/... -v -short

test-integration:
	@echo "Running integration tests..."
	cd $(BACKEND_DIR) && go test ./integration/... -v

# Build targets
build: build-backend build-frontend

build-backend:
	@echo "Building backend..."
	cd $(BACKEND_DIR) && go build -o app ./main.go

build-frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm install && npm run build

# Docker targets
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 \
		-e FIREBASE_SERVICE_ACCOUNT="$(shell cat backend/.env | grep FIREBASE_SERVICE_ACCOUNT | cut -d'=' -f2-)" \
		-e SUBSCOLLECTION_ID="$(shell cat backend/.env | grep SUBSCOLLECTION_ID | cut -d'=' -f2)" \
		-e ADMIN_PASSWORD="$(shell cat backend/.env | grep ADMIN_PASSWORD | cut -d'=' -f2)" \
		-e PORT=8080 \
		-e CORS_ORIGIN=http://localhost:8080 \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

# Clean targets
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BACKEND_DIR)/app
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/node_modules
	cd $(BACKEND_DIR) && go clean -cache -testcache

# Dependencies
deps:
	@echo "Installing dependencies..."
	cd $(BACKEND_DIR) && go mod download && go mod tidy
	cd $(FRONTEND_DIR) && npm install

# Lint targets
lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	cd $(BACKEND_DIR) && golangci-lint run
	cd $(FRONTEND_DIR) && npm run lint || true

