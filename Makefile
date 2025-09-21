.PHONY: build
build:
	go build -v ./...

.PHONY: dev.check
dev.check:
	@echo "Checking for processes on gRPC port (50051) and HTTP port (8081)..."
	@if lsof -ti:50051 >/dev/null 2>&1; then \
		echo "Warning: Port 50051 (gRPC) is already in use. Run 'make dev.kill' to kill the processes."; \
	fi
	@if lsof -ti:8081 >/dev/null 2>&1; then \
		echo "Warning: Port 8081 (HTTP) is already in use. Run 'make dev.kill' to kill the processes."; \
	fi

.PHONY: dev
dev: dev.check
	@docker compose up -d
	@go run cmd/app/main.go

.PHONY: dev.down
dev.down:
	@docker compose down

.PHONY: dev.kill
dev.kill:
	@echo "Killing processes listening on gRPC port (50051) and HTTP port (8081)..."
	@lsof -ti:50051 | xargs kill -9 2>/dev/null || true
	@lsof -ti:8081 | xargs kill -9 2>/dev/null || true
	@echo "Processes killed."

.PHONY: seed
seed:
	@docker compose exec -T postgres psql -U ${DB_USER} -d ${DB_NAME} -f /seed/data.sql

.PHONY: logs
logs:
	docker compose logs -f

.PHONY: clean
clean:
	docker compose down -v

.PHONY: lint.go
lint.go:
	@golangci-lint run
	$(call format)

.PHONY: lint.go.fix
lint.go.fix:
	@golangci-lint run --fix
	$(call format)

.PHONY: lint.proto
lint.proto:
	@go tool buf lint
	$(call format)

.PHONY: test
test:
	@go test ./internal/...

.PHONY: migrate
migrate:
	@go run internal/infrastructure/postgres/migrate/main.go

.PHONY: gen.mocks
gen.mocks:
	@go generate ./...
	$(call format)

.PHONY: gen.ent
gen.ent:
	cd internal/infrastructure/postgres/ent && go generate ./...

.PHONY: gen.buf
gen.buf:
	@rm -rf ./api/generated
	@go tool buf generate
	@mv api/generated/api/proto/* api/generated/ 2>/dev/null || true
	@rm -rf api/generated/api
	@find api/generated -type f ! -name '*.pb.go' ! -name '*.pb.gw.go' -delete 2>/dev/null || true
	@find api/generated -type d -empty -delete 2>/dev/null || true
	$(call format)

.PHONY: format
format:
	$(call format)

define format
	@go fmt ./internal/...
	@go tool goimports -w ./internal
	@go tool gofumpt -l -w .
	@go mod tidy
endef
