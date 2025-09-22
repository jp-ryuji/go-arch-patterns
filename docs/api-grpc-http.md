# API Documentation

## Architecture

The API is built using gRPC with a gRPC-Gateway that translates HTTP/JSON requests into gRPC calls. This architecture provides:

- Automatic API documentation generation from protobuf definitions
- Easy integration between HTTP/JSON and gRPC
- Type-safe API definitions

## Configuration Files

This project uses [buf](https://buf.build) to manage Protocol Buffer files and generation:

- **`buf.yaml`** - Main configuration file that sets up the buf module, including linting standards, breaking change detection rules, and external dependencies
- **`buf.lock`** - Dependency lock file that ensures reproducible builds by locking the versions of external dependencies (similar to `go.sum` in Go modules)
- **`buf.gen.yaml`** - Code generation configuration that specifies which plugins to run, their output locations, and generation options for converting Protocol Buffer definitions into source code

## Endpoints (examples)

### Create Car

Creates a new car entity.

- **URL**: `/v1/cars`
- **Method**: `POST`
- **Request Body**:

  ```json
  {
    "tenant_id": "string",
    "model": "string"
  }
  ```

- **Response**:

  ```json
  {
    "car": {
      "id": "string",
      "tenant_id": "string",
      "model": "string",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  }
  ```

### Get Car

Retrieves a car by its ID.

- **URL**: `/v1/cars/{id}`
- **Method**: `GET`
- **Response**:

  ```json
  {
    "car": {
      "id": "string",
      "tenant_id": "string",
      "model": "string",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  }
  ```

### List Cars

Retrieves a list of cars for a tenant with pagination support.

- **URL**: `/v1/cars`
- **Method**: `GET`
- **Query Parameters**:
  - `tenant_id` (required): The tenant ID to filter cars
  - `page_size` (optional): Number of cars to return per page (default: 10)
  - `page_token` (optional): Token for pagination

- **Response**:

  ```json
  {
    "cars": [
      {
        "id": "string",
        "tenant_id": "string",
        "model": "string",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ],
    "next_page_token": "string"
  }
  ```

## Protocol Buffers

The API is defined using Protocol Buffers in the following files:

- `api/proto/car/v1/car.proto` - Defines the Car message structure
- `api/proto/car/v1/car_service.proto` - Defines the gRPC service and methods:
  - `CreateCar` - Creates a new car
  - `GetCar` - Retrieves a car by ID
  - `ListCars` - Retrieves a list of cars with pagination

### Dependency Management

This project uses buf to manage Protocol Buffer dependencies. External dependencies are declared in `buf.yaml` and locked in `buf.lock` to ensure reproducible builds.

To update dependencies:

```bash
buf mod update
```

This will update `buf.lock` with the latest versions of dependencies while respecting the constraints defined in `buf.yaml`.

## Implementation

The API is implemented in the following directories:

- `internal/interface/grpc/car/v1/service.go` - gRPC service implementation
- `internal/interface/http/server.go` - HTTP server with gRPC-Gateway

### Dependency Injection

The project uses a dependency injection container located at `internal/di/container.go` to manage service dependencies. When adding new services, you'll need to update this container to:

1. Create new repositories for the service
2. Create new application services
3. Register the gRPC service implementations
4. Update the HTTP server configuration if needed

The container provides a centralized place to manage dependencies and ensures that services are properly wired together. It accepts an existing database client and port configurations, making it flexible and testable.

For detailed instructions on how to add new services, see [Adding New Services](./adding_new_services.md).

## Code Generation

To generate code from Protocol Buffer definitions, follow these steps:

1. **Modify Protocol Buffer files**:
   - Update `api/proto/car/v1/car.proto` for message definitions
   - Update `api/proto/car/v1/car_service.proto` for service definitions
   - Update `api/proto/common/v1/common.proto` for common message types

2. **Lint Protocol Buffer files**:
   Run the following command to check for any issues in the Protocol Buffer definitions:

   ```bash
   make lint.proto
   ```

3. **Update dependencies (if needed)**:
   If you've added or updated external dependencies in `buf.yaml`, run:

   ```bash
   buf mod update
   ```

   This will update `buf.lock` with the latest versions of dependencies.

4. **Generate Go code**:
   Run the following command to generate Go code from the Protocol Buffer definitions:

   ```bash
   make gen.buf
   ```

5. **Generated files**:
   The following files are generated:
   - `api/generated/car/v1/car.pb.go` - Message types
   - `api/generated/car/v1/car_service_grpc.pb.go` - gRPC service interface
   - `api/generated/car/v1/car_service.pb.gw.go` - gRPC-Gateway handlers
   - `api/generated/car/v1/car_service.pb.go` - Service message types
   - `api/generated/common/v1/common.pb.go` - Common message types

### buf.gen.yaml Configuration

The `buf.gen.yaml` file configures how Protocol Buffer files are processed to generate code:

- **go plugin** - Generates Go message types (`*.pb.go`)
- **go-grpc plugin** - Generates gRPC service interfaces (`*_grpc.pb.go`)
- **grpc-gateway plugin** - Generates HTTP/JSON to gRPC gateway code (`*.pb.gw.go`)

All generated code is placed in the `api/generated` directory with paths relative to the source files.

## Testing API Endpoints

### Prerequisites

Before testing the API endpoints, ensure that:

1. The database and other services are running:

   ```bash
   make dev.up
   ```

2. The database migrations have been applied:

   ```bash
   make migrate
   ```

3. The database has been seeded with sample data:

   ```bash
   make seed
   ```

4. The application server is running:

   ```bash
   make dev.run
   ```

### Testing with curl

#### List Cars

Retrieve a list of cars for a specific tenant:

```bash
curl -X GET "http://localhost:8081/v1/cars?tenant_id=TENANT_ID"
```

Replace `TENANT_ID` with an actual tenant ID from your database.

#### Get Car by ID

Retrieve a specific car by its ID:

```bash
curl -X GET "http://localhost:8081/v1/cars/CAR_ID"
```

Replace `CAR_ID` with an actual car ID from your database.

#### Create Car

Create a new car:

```bash
curl -X POST "http://localhost:8081/v1/cars" \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": "TENANT_ID", "model": "New Car Model"}'
```

Replace `TENANT_ID` with an actual tenant ID from your database.
