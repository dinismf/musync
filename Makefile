.PHONY: build run test clean docker-up docker-down docker-logs lint help dev-frontend dev-backend dev

# Variables
BINARY_NAME=musync
MAIN_PATH=./backend/cmd/api/main.go
DOCKER_COMPOSE=docker compose

# Colors for help message
BLUE=\033[0;34m
NC=\033[0m # No Color

help: ## Display this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "  ${BLUE}%-20s${NC} %s\n", $$1, $$NF }' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building..."
	@cd backend && go build -o $(BINARY_NAME) $(MAIN_PATH)

run: ## Run the application
	@echo "Running..."
	@cd backend && go run $(MAIN_PATH)

test: ## Run tests
	@echo "Running tests..."
	@cd backend && go test -v ./...

clean: ## Clean build files
	@echo "Cleaning..."
	@rm -f backend/$(BINARY_NAME)
	@go clean

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	@$(DOCKER_COMPOSE) up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@$(DOCKER_COMPOSE) down

docker-logs: ## View Docker container logs
	@echo "Viewing Docker logs..."
	@$(DOCKER_COMPOSE) logs -f

docker-clean: ## Stop containers and remove volumes
	@echo "Cleaning Docker environment..."
	@$(DOCKER_COMPOSE) down -v

lint: ## Run linter
	@echo "Running linter..."
	@cd backend && go vet ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@cd backend && go mod tidy
	@cd frontend && npm install

migrate: ## Run database migrations
	@echo "Running migrations..."
	@cd backend && go run $(MAIN_PATH) migrate

dev-frontend: ## Start frontend development server
	@echo "Starting frontend development server..."
	@cd frontend && npm start

dev-backend: ## Start backend development server
	@echo "Starting backend development server..."
	@cd backend && go run $(MAIN_PATH)

dev: docker-up ## Start both frontend and backend development servers
	@echo "Starting development servers..."
	@make dev-backend & make dev-frontend

# Development shortcuts
dev-setup: deps docker-up ## Setup development environment
	@echo "Development environment ready!"

dev-clean: docker-clean clean ## Clean development environment
	@echo "Development environment cleaned!" 