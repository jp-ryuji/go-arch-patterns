# golang

## Development Tools

This project uses [Go 1.24's new tool dependency management feature](https://tip.golang.org/doc/go1.24#tools). Tool dependencies are declared in the `go.mod` file using `tool` directives and can be managed using the `-tool` flag with `go get`.

For example, to add a new tool dependency:

```bash
go get -tool github.com/some/tool/cmd/tool-name
```

To upgrade all tools:

```bash
go get tool
```

To install all tools:

```bash
go install tool
```

## Testing

To run all tests in the project:

```bash
make test
```

### Running Tests in Short Mode

Some tests, particularly integration tests that require external dependencies like Docker, can be skipped when running in short mode. This is useful for quick local testing or CI pipelines where you want to avoid lengthy setup processes.

The project's Makefile target for testing (`make test`) runs all tests by default. To take advantage of the short mode skipping feature, you can modify the Makefile or run the Go command directly with the `-short` flag:

```bash
go test -short ./internal/...
```

Integration tests that use Docker (such as repository tests) will automatically be skipped when running in short mode thanks to the `SkipIfShort` utility function in the test utilities package.

To run integration tests specifically:

```bash
go test -tags=integration ./internal/...
```

To run integration tests in verbose mode:

```bash
go test -tags=integration ./internal/... -v
```

### Integration Testing with Docker

This project uses [Dockertest](https://github.com/ory/dockertest) for integration testing with a PostgreSQL database. The integration tests are tagged with `//go:build integration` and are only run when the `-tags=integration` flag is provided.

When adding new database models, they must be registered in the `allModels` slice in `internal/infrastructure/postgres/repository/testutil/testutil.go` to ensure they are migrated during test setup.

### Mocking

This project uses [gomock](https://github.com/uber-go/mock) for mocking dependencies in tests.

Mocks are generated using `go:generate` directives placed directly before interface definitions. To generate mocks, use:

```bash
make gen.mocks
```

## lint

There are two linting targets available:

```bash
# Run linter to detect issues (used by pre-commit hook)
make lint.go

# Run linter and automatically fix issues where possible
make lint.go.fix
```

For information about the Git pre-commit hook, see [Installation Guide](installation_guide.md).

## Database Schema Updates

For information about database schema updates, including Ent migration and code generation, see [Database Schema Updates](database_schema_updates.md).
