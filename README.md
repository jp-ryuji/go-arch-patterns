# go-arch-patterns

This repository demonstrates various architectural patterns and design principles for building scalable backend systems in Go. The codebase showcases architectural patterns, database design techniques, and SaaS-specific implementations through a practical car rental platform example.

For more details about the architecture, see [Software Architecture](docs/software_architecture.md).
For the data model and relationships, see [Entity Relationship Diagram](docs/er-diagram.md).

## Key Implementation Examples

This repository demonstrates the following software engineering concepts:

### Domain-Driven Design & Architecture

- **Value Object**:
  - *Definition*: See [Email value object](internal/domain/value/email.go) with [tests](internal/domain/value/email_test.go)
  - *Usage*: See [Individual entity](internal/domain/entity/individual.go) using the Email value object

### Database Design Patterns

- **Class Table Inheritance**: See [Renter entity](internal/domain/entity/renter.go) as the base class with [Company](internal/domain/entity/company.go) and [Individual](internal/domain/entity/individual.go) as specialized subclasses

### Microservices Patterns

- **Outbox Pattern**: Reliable event publishing for distributed systems. See [documentation](docs/outbox_pattern.md) and [implementation](internal/application/service/car_impl.go)

### SaaS Patterns

- **PostgreSQL Row-Level Security**: Multi-tenant data isolation *(Coming Soon)*

## Documentation

Find specific documentation in the [docs/](docs/) folder:

- [Software Architecture](docs/software_architecture.md)
- [Entity Relationship Diagram](docs/er-diagram.md)
- [Installation Guide](docs/installation_guide.md)
- [Go Development Guide](docs/golang.md)
- [Database Schema Updates](docs/database_schema_updates.md)
- [Ent ORM Setup](docs/ent.md)
  - [Go ORM/Query Builder Selection Summary](docs/orm-selection-summary.md)
- [Outbox Pattern Implementation](docs/outbox_pattern.md)
- [API (gRPC with gRPC-Gateway) Documentation](docs/api-grpc-http.md)
- [Adding New Services](docs/adding_new_services.md)

## Disclaimer

This repository is intended for educational and demonstration purposes. It is not recommended for production use without significant modifications.

1. **Database Schema Management**: The project uses Ent's auto migration feature to manage database schema changes for convenience. Proper migration strategies should be implemented for production environments. For more details on how database schema updates are handled, see [Database Schema Updates](docs/database_schema_updates.md).

2. **Database Security**: SSL mode is disabled for database access (see [internal/infrastructure/postgres/migrate/main.go](internal/infrastructure/postgres/migrate/main.go)). This configuration is insecure and should be enabled in production environments.
