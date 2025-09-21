# Software Architecture

This project follows a Domain-Driven Design (DDD) approach with an Onion Architecture structure. This pattern emphasizes separation of concerns and dependency inversion, with dependencies pointing inward toward the domain core.

In this architecture:

- The **core domain** (business logic) resides at the center and is independent of external systems
- **Domain entities** represent the core business objects with identity
- **Value objects** are immutable objects that describe aspects of the domain
- **Repository interfaces** define the contracts for data access (ports)
- **Application services** orchestrate use cases and business flows
- **DTOs** (Data Transfer Objects) carry data between processes
- **Adapters** come in two types:
  - **Secondary/Driven Adapters**: Implementations of repository interfaces that connect to external systems (databases, message queues, etc.)
  - **Primary/Driving Adapters**: Interface adapters that expose application functionality to external clients (gRPC services, HTTP handlers, etc.)

## Data Flow

This section illustrates how data flows through all architectural layers during request/response cycles, showing the conversion points between different data representations.

For specific implementation details of the gRPC layer, see [API Documentation](./api-grpc-http.md).

### Request Path (Incoming)

#### gRPC Client Request

```plaintext
gRPC protobuf Request
    ↓
gRPC Server Handler
    ↓ converts protobuf → DTO
Application Service
    ↓ converts DTO → Entity
Domain Layer
    ↓ Entity → Repository Interface
Repository Implementation
    ↓ converts Entity → Database Model
Database
```

#### HTTP/REST Client Request

```plaintext
HTTP/REST JSON Request
    ↓
grpc-gateway (HTTP → gRPC conversion)
    ↓ converts JSON → protobuf
gRPC Server Handler
    ↓ converts protobuf → DTO
Application Service
    ↓ converts DTO → Entity
Domain Layer
    ↓ Entity → Repository Interface
Repository Implementation
    ↓ converts Entity → Database Model
Database
```

### Response Path (Outgoing)

#### gRPC Client Response

```plaintext
Database
    ↓ returns Database Model
Repository Implementation
    ↓ converts Database Model → Entity
Domain Layer
    ↓ returns Entity
Application Service
    ↓ converts Entity → DTO
gRPC Server Handler
    ↓ returns protobuf Response
```

#### HTTP/REST Client Response

```plaintext
Database
    ↓ returns Database Model
Repository Implementation
    ↓ converts Database Model → Entity
Domain Layer
    ↓ returns Entity
Application Service
    ↓ converts Entity → DTO
gRPC Server Handler
    ↓ returns protobuf Response
grpc-gateway (gRPC → HTTP conversion)
    ↓ converts protobuf → JSON
HTTP/REST JSON Response
```

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
    ├── domain                   # Core Domain Layer (innermost)
    │   ├── entity               # Domain entities with identity
    │   ├── value                # Value objects (immutable)
    │   ├── service              # Domain services (business logic)
    │   └── repository           # Repository interfaces (ports)
    ├── application              # Application Layer
    │   ├── service              # Application services (orchestration)
    │   └── dto                  # Data transfer objects
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
