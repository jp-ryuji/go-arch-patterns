# go-sample

## Software Architecture

This project follows the Ports and Adapters Architecture (also known as Hexagonal Architecture). This pattern, along with Onion Architecture and Clean Architecture, share common principles of separation of concerns and dependency inversion, though each has its own emphasis:

- **Ports and Adapters** emphasizes defining interfaces (ports) at the system boundaries and implementing adapters for external systems
- **Onion Architecture** focuses on layers with dependencies pointing inward toward the domain core
- **Clean Architecture** emphasizes dependency rules where inner layers should not depend on outer layers

All these patterns promote loose coupling between the application core and external systems.

In this architecture:

- The **core domain** (business logic) resides at the center and is independent of external systems
- **Ports** are interfaces defined in the core that allow communication with external systems
- **Adapters** are implementations of those ports that connect to external systems (databases, HTTP APIs, etc.)

Structure:

```plaintext
├── cmd
│   └── app
│       └── main.go
├── docker
├── docs
├── internal
│   ├── domain
│   │   ├── model         # domain model
│   │   │   └── factory
│   │   ├── repository    # repository interface (port)
│   │   └── service       # domain service
│   ├── infrastructure
│   │   ├── cmd
│   │   │   ├── internal
│   │   │   └── root.go
│   │   ├── database      # database adapter implementation
│   │   ├── http          # HTTP adapter implementation
│   │   ├── redis         # Redis adapter implementation
│   │   └── s3            # S3 adapter implementation
│   └── usecase           # application service
├── schema
│   ├── openapi
│   └── proto
├── seed-data
└── tools
```

### Go `internal` Directory

The `internal` directory is a special directory in Go that restricts access to its contents. Only code within the same module (in this case, `go-sample`) can import packages from `internal` directories. This prevents other projects from importing and depending on our internal implementation details, which helps maintain a clean public API and allows us to change internal implementations without breaking external dependencies.

In this project, all core business logic, domain models, use cases, and infrastructure implementations are placed under the `internal` directory to enforce this encapsulation and prevent accidental exposure of internal details as part of the public API.

Note: For most web applications, especially those with clear architectural boundaries like this one, the root `internal` directory is optional. The separation of concerns is already evident through the domain, usecase, and infrastructure packages. The `internal` directory provides an additional layer of enforcement but isn't strictly necessary when the boundaries are well-defined through naming conventions and architecture.

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

For more detailed information about Go tooling and linting in this project, see [docs/golang.md](docs/golang.md).

More information about Go 1.24's tool dependency management can be found in the [Go 1.24 release notes](https://go.dev/doc/go1.24#go-command).
