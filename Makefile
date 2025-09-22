# Development Makefile for Go Application

# Default target
.DEFAULT_GOAL := help

# Variables
GRPC_PORT := 50051
HTTP_PORT := 8081

#==============================================================================
# BUILD
#==============================================================================

.PHONY: build
build: ## Build the application
	go build -o ./.bin/app ./cmd/app

#==============================================================================
# DEVELOPMENT
#==============================================================================

.PHONY: dev.run
dev.run: ## Start the application with hot reload (assumes services are running)
	@go tool air -c .air.toml

.PHONY: dev.check
dev.check: ## Check if development ports are available
	@echo "Checking for processes on gRPC port ($(GRPC_PORT)) and HTTP port ($(HTTP_PORT))..."
	@if lsof -ti:$(GRPC_PORT) >/dev/null 2>&1; then \
		echo "Warning: Port $(GRPC_PORT) (gRPC) is already in use. Run 'make dev.kill' to kill the processes."; \
	fi
	@if lsof -ti:$(HTTP_PORT) >/dev/null 2>&1; then \
		echo "Warning: Port $(HTTP_PORT) (HTTP) is already in use. Run 'make dev.kill' to kill the processes."; \
	fi

.PHONY: dev.up
dev.up: ## Start infrastructure services only
	@echo "Starting infrastructure services..."
	@docker compose up -d
	@echo "Services started. Run 'make dev.run' to start the application."

.PHONY: dev.restart
dev.restart: ## Restart only the application (keep services running)
	@$(MAKE) --no-print-directory dev.kill || true
	@sleep 1
	@$(MAKE) --no-print-directory dev.run

.PHONY: dev.down
dev.down: ## Stop all services
	@docker compose down

.PHONY: dev.kill
dev.kill: ## Kill processes on development ports
	@echo "Killing processes listening on gRPC port ($(GRPC_PORT)) and HTTP port ($(HTTP_PORT))..."
	@lsof -ti:$(GRPC_PORT) | xargs kill -9 2>/dev/null || true
	@lsof -ti:$(HTTP_PORT) | xargs kill -9 2>/dev/null || true
	@echo "Processes killed."

.PHONY: dev.status
dev.status: ## Check service status and health
	@echo "Service Status:"
	@docker compose ps
	@echo "\nHealth Checks:"
	@docker compose exec postgres pg_isready -U ${DB_USER} -d ${DB_NAME} 2>/dev/null \
		&& echo "✅ PostgreSQL: Ready" || echo "❌ PostgreSQL: Not Ready"
	@docker compose exec redis redis-cli ping 2>/dev/null | grep -q PONG \
		&& echo "✅ Redis: Ready" || echo "❌ Redis: Not Ready" || true

#==============================================================================
# DATABASE
#==============================================================================

.PHONY: migrate
migrate: ## Run database migrations
	@go run internal/infrastructure/postgres/migrate/main.go

.PHONY: seed
seed: ## Seed database with test data
	@docker compose exec -T postgres psql -U ${DB_USER} -d ${DB_NAME} -f /seed/data.sql

.PHONY: psql
psql: ## Connect to PostgreSQL database
	@docker compose exec postgres psql -U ${DB_USER} -d ${DB_NAME}

#==============================================================================
# CODE GENERATION
#==============================================================================

.PHONY: gen.mocks
gen.mocks: ## Generate mocks
	@go generate ./...
	$(call format)

.PHONY: gen.ent
gen.ent: ## Generate Ent schema
	cd internal/infrastructure/postgres/ent && go generate ./...

.PHONY: gen.buf
gen.buf: ## Generate protobuf code
	@rm -rf ./api/generated
	@go tool buf generate
	@mv api/generated/api/proto/* api/generated/ 2>/dev/null || true
	@rm -rf api/generated/api
	@find api/generated -type f ! -name '*.pb.go' ! -name '*.pb.gw.go' -delete 2>/dev/null || true
	@find api/generated -type d -empty -delete 2>/dev/null || true
	$(call format)

.PHONY: gen.all
gen.all: gen.buf gen.ent gen.mocks ## Generate all code (protobuf, ent, mocks)

#==============================================================================
# LINTING & FORMATTING
#==============================================================================

.PHONY: lint.go
lint.go: ## Run Go linter
	@golangci-lint run
	$(call format)

.PHONY: lint.go.fix
lint.go.fix: ## Run Go linter with auto-fix
	@golangci-lint run --fix
	$(call format)

.PHONY: lint.proto
lint.proto: ## Run protobuf linter
	@go tool buf lint
	$(call format)

.PHONY: lint.all
lint.all: lint.go lint.proto ## Run all linters

.PHONY: format
format: ## Format all code
	$(call format)

#==============================================================================
# TESTING
#==============================================================================

.PHONY: test
test: ## Run tests
	@go test ./internal/...

.PHONY: test.coverage
test.coverage: ## Run tests with coverage report
	@go test -coverprofile=coverage.out ./internal/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

#==============================================================================
# MONITORING & DEBUGGING
#==============================================================================

.PHONY: logs
logs: ## Follow logs from all services
	docker compose logs -f

.PHONY: logs.app
logs.app: ## Follow application logs
	@tail -f tmp/air.log 2>/dev/null || echo "No air.log found. Start the app first."

.PHONY: logs.db
logs.db: ## Follow database logs
	@docker compose logs -f postgres

#==============================================================================
# CLEANUP
#==============================================================================

.PHONY: clean
clean: ## Stop services and remove volumes
	docker compose down -v

.PHONY: clean.build
clean.build: ## Clean build artifacts
	@rm -rf ./.bin
	@rm -rf ./tmp
	@echo "Build artifacts cleaned"

.PHONY: clean.all
clean.all: clean clean.build ## Clean everything

#==============================================================================
# UTILITY
#==============================================================================

.PHONY: tools
tools: ## Install development tools
	@echo "Installing development tools..."
	@go mod tidy
	@echo "Tools installed"

.PHONY: help
help: ## Show this help message
	@echo "Development Makefile for Go Application"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Common Workflows:"
	@echo "  make dev.up  # Start infrastructure"
	@echo "  make dev.run       # Start app (assumes services running)"
	@echo "  make dev.restart   # Restart app only"
	@echo ""
	@echo "Available Targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_.-]+:.*##/ {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

#==============================================================================
# HELPER FUNCTIONS
#==============================================================================

define format
	@go fmt ./internal/...
	@go tool goimports -w ./internal
	@go tool gofumpt -l -w .
	@go mod tidy
endef
