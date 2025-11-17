.PHONY: help test test-unit test-integration build build-frontend build-all run run-frontend docker-build docker-run clean deps

# Variables
BACKEND_DIR=backend
FRONTEND_DIR=frontend
BINARY_NAME=app
DOCKER_IMAGE=appdirect-workshop
DOCKER_TAG=latest

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Install dependencies
	@echo "Installing backend dependencies..."
	cd $(BACKEND_DIR) && go mod download
	@echo "Installing frontend dependencies..."
	cd $(FRONTEND_DIR) && npm install

test: ## Run all tests
	@echo "Running all tests..."
	cd $(BACKEND_DIR) && go test ./... -v

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	cd $(BACKEND_DIR) && go test ./internal/handlers ./internal/middleware ./internal/config -v

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	cd $(BACKEND_DIR) && go test -v -run TestIntegration

build: ## Build backend binary
	@echo "Building backend..."
	cd $(BACKEND_DIR) && go build -o $(BINARY_NAME) ./main.go

build-frontend: ## Build frontend
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm run build

build-all: build build-frontend ## Build both frontend and backend

run: ## Run backend locally
	@echo "Running backend..."
	cd $(BACKEND_DIR) && ./$(BINARY_NAME)

run-frontend: ## Run frontend dev server
	@echo "Running frontend dev server..."
	cd $(FRONTEND_DIR) && npm run dev

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 \
		-e SUBSCOLLECTION_ID=workshop-2024 \
		-e ADMIN_PASSWORD=change-me \
		-e PORT=8080 \
		-e CORS_ORIGIN=http://localhost:8080 \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f $(BACKEND_DIR)/$(BINARY_NAME)
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/node_modules/.vite
	cd $(BACKEND_DIR) && go clean

