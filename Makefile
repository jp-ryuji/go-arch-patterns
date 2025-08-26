.PHONY: lint.go
lint.go:
	@go tool golangci-lint run
	$(call format)

.PHONY: lint.go.fix
lint.go.fix:
	@go tool golangci-lint run --fix
	$(call format)

.PHONY: test
test:
	@go test ./internal/...

.PHONY: format
format:
	$(call format)

define format
	@go fmt ./...
	@go run golang.org/x/tools/cmd/goimports -w ./
	@go run mvdan.cc/gofumpt -l -w .
	@go mod tidy
endef
