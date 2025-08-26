# go-sample

## Installation & Setup

1. **Clone the repository:**

```bash
git clone https://github.com/jp-ryuji/go-sample.git
cd go-sample
```

2. **Set up git hooks:**

This project includes a pre-commit hook that automatically runs the Go linter before each commit. If the linter finds any issues, the commit will be aborted, and you'll need to fix the issues before committing again.

To install the pre-commit hook, run:

```bash
./scripts/install-hooks.sh
```

Note: Git hooks are local to your repository and are not shared through the repository. Each developer needs to run this script to install the hooks in their local environment.

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
│   │   │   ├── value     # value object
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
│   ├── pkg               # library
│   └── usecase           # application service
├── schema
│   ├── openapi
│   └── proto
└── seed-data
```

### Go `internal` Directory

The `internal` directory is a special directory in Go that restricts access to its contents. Only code within the same module (in this case, `go-sample`) can import packages from `internal` directories. This prevents other projects from importing and depending on our internal implementation details, which helps maintain a clean public API and allows us to change internal implementations without breaking external dependencies.

In this project, all core business logic, domain models, use cases, and infrastructure implementations are placed under the `internal` directory to enforce this encapsulation and prevent accidental exposure of internal details as part of the public API.

Note: For most web applications, especially those with clear architectural boundaries like this one, the root `internal` directory is optional. The separation of concerns is already evident through the domain, usecase, and infrastructure packages. The `internal` directory provides an additional layer of enforcement but isn't strictly necessary when the boundaries are well-defined through naming conventions and architecture.
