.PHONY: lint.go
lint.go:
	@go tool golangci-lint run --fix
	$(call format)
