# go-sample

This is a sample repository showcasing the Ports and Adapters architectural pattern (also known as Hexagonal Architecture) and Domain-Driven Design (DDD) methodology, with a focus on fundamental concepts such as value objects. For more details about the software architecture, see [Software Architecture](docs/software_architecture.md).

The system represents a SaaS platform for car rental companies. For more details about the system, see [Entity Relationship Diagram](docs/er-diagram.md).

## Key Implementation Examples

This repository demonstrates the following software engineering concepts:

- **Value Object**:
  - *Definition*: See [Email value object](internal/domain/model/value/email.go) with [tests](internal/domain/model/value/email_test.go)
  - *Usage*: See [Individual entity](internal/domain/model/individual.go) using the Email value object
- **Class Table Inheritance**: See [Renter model](internal/domain/model/renter.go) as the base class with [Company](internal/domain/model/company.go) and [Individual](internal/domain/model/individual.go) as specialized subclasses

## Documentation

Find specific documentation in the [docs/](docs/) folder:

- [Software Architecture](docs/software_architecture.md)
- [Entity Relationship Diagram](docs/er-diagram.md)
- [Installation Guide](docs/installation_guide.md)
- [Go Development Guide](docs/golang.md)
- [Database Schema Updates](docs/database_schema_updates.md)
- [Ent ORM Setup](docs/ent.md)

## Disclaimer

This repository is intended for educational and demonstration purposes. It is not recommended for production use without significant modifications.

1. **Database Schema Management**: The project uses Ent's auto migration feature to manage database schema changes for convenience. Proper migration strategies should be implemented for production environments. For more details on how database schema updates are handled, see [Database Schema Updates](docs/database_schema_updates.md).

2. **Database Security**: SSL mode is disabled for database access (see [internal/infrastructure/postgres/migrate/main.go](internal/infrastructure/postgres/migrate/main.go)). This configuration is insecure and should be enabled in production environments.

3. **Configuration Management**: Environment variable handling is implemented in a basic way in [internal/infrastructure/postgres/client.go](internal/infrastructure/postgres/client.go). For production use, consider implementing a more robust configuration management solution such as [Viper](https://github.com/spf13/viper) or similar libraries to handle configuration files, defaults, and validation in a more structured manner.
