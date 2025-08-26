# golang

## Development Tools

This project uses Go 1.24's new tool dependency management feature. Tool dependencies are declared in the `go.mod` file using `tool` directives and can be managed using the `-tool` flag with `go get`.

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

## lint

There are two linting targets available:

```bash
# Run linter to detect issues (used by pre-commit hook)
make lint.go

# Run linter and automatically fix issues where possible
make lint.go.fix
```

For information about the Git pre-commit hook, see the main README file.
