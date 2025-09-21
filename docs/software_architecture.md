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
  - **Primary/Driving Adapters**: Interface adapters that expose application functionality to external clients (gRPC services, HTTP handlers, etc.)

## Directory Structure

```plaintext
├── cmd
│   └── app
│       └── main.go
├── api                          # Protocol Buffers definitions
│   ├── proto
│   │   ├── car
│   │   │   └── v1
│   │   │       ├── car.proto
│   │   │       └── car_service.proto
│   │   └── common
│   │       └── v1
│   │           └── common.proto
│   └── generated                # All generated code
│       ├── car
│       │   └── v1
│       │       ├── car.pb.go
│       │       ├── car_grpc.pb.go
│       │       └── car.pb.gw.go
│       └── common
│           └── v1
│               └── common.pb.go
├── docker
├── docs
└── internal
    ├── config                   # Configuration
    ├── di                       # Dependency Injection
    ├── domain                   # Core Domain Layer (innermost)
    │   ├── entity               # Domain entities with identity
    │   ├── value                # Value objects (immutable)
    │   ├── service              # Domain services (business logic)
    │   └── repository           # Repository interfaces (ports)
    ├── application              # Application Layer
    │   ├── input                # Data transfer objects (input)
    │   ├── output               # Data transfer objects (output)
    │   └── service              # Application services (orchestration)
    ├── infrastructure           # Infrastructure Layer (outermost)
    │   ├── postgres             # PostgreSQL adapter
    │   │   ├── ent              # Ent schema design
    │   │   ├── entgen           # (Generated type-safe query code using Ent)
    │   │   ├── repository       # Repository implementations
    │   │   └── migration
    │   ├── redis                # Redis adapter
    │   ├── messaging            # Message queue adapters (if needed)
    │   └── config               # Configuration
    ├── interface                # Interface Adapters
    │   ├── grpc                 # gRPC service implementations
    │   │   ├── car
    │   │   │   └── v1
    │   │   │       └── service.go    # carServiceServer implementation
    │   │   ├── server.go             # gRPC server setup
    │   │   └── interceptor           # gRPC interceptors
    │   └── http                      # HTTP Gateway setup
    │       ├── server.go             # HTTP server + gRPC-Gateway
    │       ├── middleware.go         # HTTP middleware
    │       └── handler.go            # Custom HTTP handlers (if any)
    └── pkg                      # Shared utilities/libraries
```

## Go `internal` Directory

The `internal` directory is a special directory in Go that restricts access to its contents. Only code within the same module (in this case, `go-arch-patterns`) can import packages from `internal` directories. This prevents other projects from importing and depending on our internal implementation details, which helps maintain a clean public API and allows us to change internal implementations without breaking external dependencies.

In this project, all core business logic, domain models, use cases, and infrastructure implementations are placed under the `internal` directory to enforce this encapsulation and prevent accidental exposure of internal details as part of the public API.

Note: For most web applications, especially those with clear architectural boundaries like this one, the root `internal` directory is optional. The separation of concerns is already evident through the domain, application, and infrastructure packages. The `internal` directory provides an additional layer of enforcement but isn't strictly necessary when the boundaries are well-defined through naming conventions and architecture.

## Libraries Used

- [Ent](https://entgo.io/) for ORM and database operations
- [Dockertest](https://github.com/ory/dockertest) for booting up ephermal docker images for integration tests using Postgres or such
- [gomock](https://github.com/uber-go/mock) for mocking dependencies in tests
