# Data Flow

This doc illustrates how data flows through all architectural layers during request/response cycles, showing the conversion points between different data representations.

For specific implementation details of the gRPC layer, see [API Documentation](./api-grpc-http.md).

## Request Path (Incoming)

### gRPC Client Request

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

### HTTP/REST Client Request

```plaintext
HTTP/REST JSON Request
    ↓
grpc-connect (HTTP → gRPC conversion)
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

## Response Path (Outgoing)

### gRPC Client Response

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

### HTTP/REST Client Response

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
grpc-connect (gRPC → HTTP conversion)
    ↓ converts protobuf → JSON
HTTP/REST JSON Response
```
