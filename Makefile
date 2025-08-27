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

.PHONY: test
test:
	@go test ./internal/...

.PHONY: gen.gorm
gen.gorm:
	@go run internal/infrastructure/postgres/gen/generator.go

.PHONY: format
format:
	$(call format)

define format
	@go fmt ./...
	@go tool goimports -w ./
	@go tool gofumpt -l -w .
	@go mod tidy
endef
