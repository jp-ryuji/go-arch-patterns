.PHONY: build
build:
	go build -v ./...

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
	$(call format)

.PHONY: gen.proto
gen.proto:
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
	@go fmt ./...
	@go tool goimports -w ./
	@go tool gofumpt -l -w .
	@go mod tidy
endef
