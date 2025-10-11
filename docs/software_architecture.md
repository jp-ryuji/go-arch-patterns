# Software Architecture

This project follows a Domain-Driven Design (DDD) approach with an Onion Architecture structure. This pattern emphasizes separation of concerns and dependency inversion, with dependencies pointing inward toward the domain core.

In this architecture:

- **Domain entities** represent the core business objects with identity
- **Value objects** are immutable objects that describe aspects of the domain
- **Repository interfaces** define the contracts for data access (ports)
- **Application services** orchestrate use cases and business flows
- **DTOs** (Data Transfer Objects) carry data between processes
- **Adapters** come in two types:
  - **Secondary/Driven Adapters**: Implementations of repository interfaces that connect to external systems (databases, message queues, etc.)
  - **Primary/Driving Adapters**: Presentation layer that exposes application functionality to external clients (gRPC-Connect services, HTTP handlers, etc.)

## Directory Structure

```plaintext
├── cmd
│   └── app
│       └── main.go
├── api                          # Protocol Buffers definitions
│   ├── proto
│   │   ├── car/v1               # gRPC service definitions
│   │   └── common/v1            # Shared protobuf types
│   └── generated                # Generated gRPC/Connect/Go code
├── docker                       # Local development tooling
├── docs                         # Architecture and operational docs
└── internal
    ├── config                   # Configuration loading
    ├── di                       # Dependency injection container
    ├── domain
    │   └── repository           # Shared interfaces (outbox, transactions)
    ├── infrastructure           # Outer adapters (DB, messaging, etc.)
    │   └── postgres
    │       ├── ent              # Ent schema definitions
    │       ├── entgen           # Generated ent client code
    │       └── repository       # Persistence adapters per aggregate
    ├── presentation             # Driving adapters (HTTP / Connect)
    │   ├── connect/car/v1       # gRPC-Connect handlers
    │   └── http                 # HTTP gateway and middleware
    ├── rental                   # Rental bounded context
    │   ├── application
    │   │   ├── input            # Use-case input DTOs
    │   │   ├── output           # Use-case output DTOs
    │   │   └── service          # Application services + tests/mocks
    │   ├── car.go               # Car aggregate
    │   ├── rental.go            # Rental aggregate
    │   ├── rental_option.go     # Rental option entity
    │   ├── option.go            # Configurable rental option
    │   ├── repository.go        # Port definitions (CarRepository, etc.)
    │   └── mocks                # Generated mocks for repository interfaces
    ├── renter                   # Renter bounded context
    │   ├── entity_renter.go     # Base renter identity
    │   ├── entity_individual.go # Individual profile
    │   ├── entity_company.go    # Company profile + sizes
    │   ├── email.go             # Email value object
    │   ├── factory.go           # Fixture builders for tests
    │   ├── repository.go        # Renter/Company repository ports
    │   └── mocks                # Generated mocks for renter repositories
    ├── tenant                   # Tenant bounded context
    │   ├── tenant.go            # Tenant aggregate
    │   ├── factory.go           # Fixture builders
    │   ├── repository.go        # Tenant repository port
    │   └── mocks                # Generated mocks for tenant repository
    └── pkg                      # Shared utilities (IDs, nullable helpers)
```

Each bounded context keeps its own onion layers (domain → application → adapters) so business rules evolve independently. Shared cross-cutting abstractions such as the outbox and transaction manager interfaces remain under `internal/domain/repository`, while infrastructure and presentation layers wire the contexts together through dependency injection.

## Go `internal` Directory

The `internal` directory is a special directory in Go that restricts access to its contents. Only code within the same module (in this case, `go-arch-patterns`) can import packages from `internal` directories. This prevents other projects from importing and depending on our internal implementation details, which helps maintain a clean public API and allows us to change internal implementations without breaking external dependencies.

In this project, all core business logic, domain models, use cases, and infrastructure implementations are placed under the `internal` directory to enforce this encapsulation and prevent accidental exposure of internal details as part of the public API.

Note: For most web applications, especially those with clear architectural boundaries like this one, the root `internal` directory is optional. The separation of concerns is already evident through the domain, application, and infrastructure packages. The `internal` directory provides an additional layer of enforcement but isn't strictly necessary when the boundaries are well-defined through naming conventions and architecture.

## Libraries Used

- [Ent](https://entgo.io/) for ORM and database operations
- [Dockertest](https://github.com/ory/dockertest) for booting up ephermal docker images for integration tests using Postgres or such
- [gomock](https://github.com/uber-go/mock) for mocking dependencies in tests
- [Viper](https://github.com/spf13/viper) for configuration management
- [gRPC-Connect](https://connectrpc.com/) for exposing gRPC services over HTTP
