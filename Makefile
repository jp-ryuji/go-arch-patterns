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
