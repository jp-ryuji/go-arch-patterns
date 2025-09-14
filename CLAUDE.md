# AI Collaboration Guide

This document provides essential context for AI models interacting with this project. Adhering to these guidelines will ensure consistency and maintain code quality.

## 1. Project Overview & Purpose

* **Primary Goal:** This repository serves as a sample project to demonstrate various architectural patterns and design principles for building scalable backend systems in Go. It uses a car rental platform as a practical example.
* **Business Domain:** The business domain is **Car Rental**.

## 2. Core Technologies & Stack

* **Languages:** Go (version 1.25.0 as specified in `go.mod`).
* **Frameworks & Runtimes:** Standard Go runtime. `ent` is used as the ORM.
* **Databases:** The project is configured to use `PostgreSQL 17` for primary data storage, `Redis 8` for caching, and `OpenSearch`. This is defined in the `compose.yml` file.
* **Key Libraries/Dependencies:**
  * `entgo.io/ent`: For ORM and database schema management.
  * `github.com/stretchr/testify`: For assertions in tests.
  * `github.com/oklog/ulid/v2`: For generating universally unique lexicographically sortable identifiers.
  * `go.uber.org/mock`: For generating mock implementations for interfaces.
* **Package Manager(s):** Go Modules (`go mod`).

## 3. Architectural Patterns

* **Overall Architecture:** The project follows a **Clean Architecture** (or Hexagonal) pattern. This is evident from the strict separation of concerns between the `domain` and `infrastructure` layers within the `internal` directory. The domain contains core business logic and entities, while infrastructure handles external concerns like databases.
* **Directory Structure Philosophy:**
  * `/internal`: Contains all the private application source code, not intended for import by other projects.
  * `/internal/domain`: The core of the application. It defines the business models, repository interfaces, and services. It has no dependencies on any other layer.
  * `/internal/infrastructure`: Contains the implementation of the interfaces defined in the domain layer. This includes database access (`postgres`) and other external services.
  * `/internal/pkg`: Provides shared utility packages for internal use (e.g., ID generation).
  * `/docs`: Contains project documentation.
  * `/scripts`: Holds utility scripts for development, such as git hooks.

## 4. Coding Conventions & Style Guide

* **Formatting:** The project uses `golangci-lint` for linting, as defined in `.golangci.yml`. The `Makefile` also enforces standard Go formatting via `gofmt`, `goimports`, and `gofumpt`.
* **Naming Conventions:**
  * `structs`, `interfaces`, exported `functions`: `PascalCase` (e.g., `type Car struct`).
  * `variables`, unexported `functions`: `camelCase` (e.g., `carID`).
  * `files`: `snake_case` (e.g., `car_repository.go`).
* **API Design:** (Inferred) An API is likely planned but not yet implemented. It is expected to follow standard RESTful principles.
* **Error Handling:** Standard Go error handling is used, returning an `error` as the last value from functions that can fail. Errors are checked with `if err != nil`.

## 5. Key Files & Entrypoints

* **Main Entrypoint(s):** The application does not have a main entrypoint yet. A database migration entrypoint exists at `internal/infrastructure/postgres/migrate/main.go`.
* **Configuration:**
  * `go.mod`: Defines the project's module path and dependencies.
  * `.golangci.yml`: Contains the configuration for the linter.
  * `compose.yml`: Defines the local development environment services.
* **CI/CD Pipeline:** The CI pipeline is defined in `.github/workflows/go.yml` using GitHub Actions.

## 6. Development & Testing Workflow

* **Local Development Environment:** The local environment is managed with Docker Compose, as defined in `compose.yml`. The `Makefile` provides commands for common development tasks like building, testing, and linting.
* **Testing:** Tests are run using the `make test` command, which executes `go test ./internal/...`. New features should be accompanied by corresponding unit tests.
* **Linting:** Run the linter using `make lint.go`. This is also available as a pre-commit hook (`scripts/install-hooks.sh`).
* **CI/CD Process:** On every pull request, GitHub Actions automatically runs jobs to build the application, run tests, and perform linting to ensure code quality.

## 7. Specific Instructions for AI Collaboration

* **Contribution Guidelines:** No formal `CONTRIBUTING.md` exists. However, any new code should adhere to the established architectural patterns, include unit tests, and pass the linter checks (`make lint.go`).
* **Infrastructure (IaC):** No Infrastructure as Code (IaC) is currently part of the project.
* **Security:** Be mindful of security best practices. The `README.md` notes that the current database connection does not use SSL and configuration management is basic; these should not be replicated. Do not hardcode secrets or keys.
* **Dependencies:** To add a new dependency, use `go get`. After adding or updating dependencies, run `make format` to ensure `go.mod` and `go.sum` are consistent.
* **Commit Messages:** (Inferred from best practices) Commit messages should be concise and written in the imperative mood. They should describe the change clearly. Example: `feat: add user authentication service`.
