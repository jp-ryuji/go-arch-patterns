# Adding New Services

This document explains how to add new services to the car rental platform.

## Overview

When adding new services to the platform, you need to follow several steps to ensure proper integration with the existing architecture. The process involves creating new Protocol Buffer definitions, implementing gRPC services, and updating the dependency injection container.

## Steps to Add a New Service

### 1. Define the Service in Protocol Buffers

Create new Protocol Buffer files in the `api/proto` directory:

1. Create a new directory for your service under `api/proto`, e.g., `api/proto/rental/v1/`
2. Define your message types in a file like `rental.proto`
3. Define your service methods in a file like `rental_service.proto`

Example `api/proto/rental/v1/rental.proto`:

```protobuf
syntax = "proto3";

package rental.v1;

option go_package = "github.com/jp-ryuji/go-arch-patterns/api/generated/rental/v1;rentalv1";

import "google/protobuf/timestamp.proto";

message Rental {
  string id = 1;
  string tenant_id = 2;
  string car_id = 3;
  string renter_id = 4;
  google.protobuf.Timestamp start_date = 5;
  google.protobuf.Timestamp end_date = 6;
  google.protobuf.Timestamp created_at = 7;
  google.protobuf.Timestamp updated_at = 8;
}
```

Example `api/proto/rental/v1/rental_service.proto`:

```protobuf
syntax = "proto3";

package rental.v1;

import "api/proto/rental/v1/rental.proto";
import "google/api/annotations.proto";

option go_package = "github.com/jp-ryuji/go-arch-patterns/api/generated/rental/v1;rentalv1";

service RentalService {
  rpc CreateRental(CreateRentalRequest) returns (CreateRentalResponse) {
    option (google.api.http) = {
      post: "/v1/rentals"
      body: "*"
    };
  }

  rpc GetRental(GetRentalRequest) returns (GetRentalResponse) {
    option (google.api.http) = {
      get: "/v1/rentals/{id}"
    };
  }

  rpc ListRentals(ListRentalsRequest) returns (ListRentalsResponse) {
    option (google.api.http) = {
      get: "/v1/rentals"
    };
  }
}

message CreateRentalRequest {
  string tenant_id = 1;
  string car_id = 2;
  string renter_id = 3;
  int64 start_date = 4; // Unix timestamp
  int64 end_date = 5;   // Unix timestamp
}

message CreateRentalResponse {
  Rental rental = 1;
}

message GetRentalRequest {
  string id = 1;
}

message GetRentalResponse {
  Rental rental = 1;
}

message ListRentalsRequest {
  string tenant_id = 1;
  int32 page_size = 2;
  string page_token = 3;
}

message ListRentalsResponse {
  repeated Rental rentals = 1;
  string next_page_token = 2;
}
```

### 2. Lint Protocol Buffer Definitions

Before generating code, it's recommended to lint your Protocol Buffer definitions to ensure they follow best practices:

```bash
make lint.proto
```

This will check your `.proto` files for any linting issues according to the [Buf linting rules](https://docs.buf.build/lint/rules).

### 3. Generate Code from Protocol Buffers

After defining your service, generate the Go code:

```bash
make gen.buf
```

This will generate the necessary Go files in `api/generated/rental/v1/`.

### 4. Implement Domain Layer

Create the domain entities, repositories, and services:

1. Create entity in `internal/domain/entity/rental.go`
2. Create repository interface in `internal/domain/repository/rental.go`
3. Create any necessary value objects in `internal/domain/value/`

### 5. Implement Infrastructure Layer

Create the infrastructure implementations:

1. Create Ent schema in `internal/infrastructure/postgres/ent/schema/rental.go`
2. Generate Ent code: `make gen.ent`
3. Create repository implementation in `internal/infrastructure/postgres/repository/rental_repository.go`

### 6. Implement Application Layer

Create the application service in `internal/application/service/rental_service.go`.

### 7. Implement gRPC Service

Create the gRPC service implementation in `internal/interface/grpc/rental/v1/service.go`.

### 8. Update Dependency Injection Container

Update `internal/di/container.go` to include your new service:

```go
// Add new repository
rentalRepo := repository.NewRentalRepository(client)

// Add to transaction manager if needed
// Update transaction manager constructor if needed

// Create application service
rentalService := service.NewRentalService(rentalRepo, /* other dependencies */)

// Register gRPC service
rentalServiceServer := rental.NewRentalServiceServer(rentalService)
rentalv1.RegisterRentalServiceServer(s.grpcServer, rentalServiceServer)
```

### 9. Update HTTP Server (if needed)

If your service needs custom HTTP handling, update `internal/interface/http/server.go`.

### 10. Generate Mocks

Generate mocks for your new interfaces:

```bash
make gen.mocks
```

## Best Practices

1. **Follow Existing Patterns**: Look at how the car service is implemented as an example
2. **Use Proper Error Handling**: Follow the existing error handling patterns
3. **Implement Proper Logging**: Add logging where appropriate
4. **Write Tests**: Create unit tests for your new services
5. **Update Documentation**: Document your new service in the appropriate documentation files

## Testing Your New Service

1. **Unit Tests**: Write unit tests for your domain entities and services
2. **Integration Tests**: Write integration tests for your repository implementations
3. **API Tests**: Test your API endpoints manually or with automated tests
4. **End-to-End Tests**: Create end-to-end tests that exercise the full flow

## Example Checklist

When adding a new service, ensure you've completed these tasks:

- [ ] Protocol Buffer definitions created
- [ ] Protocol Buffer definitions linted
- [ ] Code generated from Protocol Buffers
- [ ] Domain entity created
- [ ] Repository interface defined
- [ ] Ent schema created
- [ ] Ent code generated
- [ ] Repository implementation created
- [ ] Application service created
- [ ] gRPC service implementation created
- [ ] Dependency injection container updated
- [ ] HTTP server updated (if needed)
- [ ] Mocks generated
- [ ] Unit tests written
- [ ] Integration tests written
- [ ] Documentation updated
