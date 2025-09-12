# AI Collaboration Guide

This document provides essential context for AI models interacting with this project. Adhering to these guidelines will ensure consistency and maintain code quality.

## 1. Project Overview & Purpose

* **Primary Goal:** This is a sample project built with Go. Based on the domain models defined under `internal/domain/model/`. The project is structured following Clean Architecture principles.
* **Business Domain:** The business domain appears to be related to Education, though the project's primary purpose is to serve as a Go sample application.

## 2. Core Technologies & Stack

* **Languages:** Go (version 1.25.0 as specified in `go.mod`).
* **Frameworks & Runtimes:** Standard Go runtime. No major web frameworks are currently in use.
* **Databases:** The infrastructure layer is not yet implemented, but the directory structure (`internal/infrastructure/database`) suggests that a database, Postgres, will be used.
* **Key Libraries/Dependencies:**
  * `github.com/oklog/ulid/v2`: For generating universally unique lexicographically sortable identifiers.
* **Package Manager(s):** Go Modules (`go mod`).

## 3. Architectural Patterns

* **Overall Architecture:** The project follows a **Clean Architecture** (or Hexagonal) pattern. This is evident from the strict separation of concerns between the `domain` and `infrastructure` layers within the `internal` directory. The domain contains core business logic and entities, while infrastructure handles external concerns like databases and HTTP servers.
* **Directory Structure Philosophy:**
  * `/internal`: Contains all the private application source code, which is not intended to be imported by other projects.
  * `/internal/domain`: The core of the application. It defines the business models, repositories (interfaces), and services. It has no dependencies on any other layer.
  * `/internal/infrastructure`: Contains the implementation of the interfaces defined in the domain layer. This includes database access, external API clients, and web handlers.
  * `/internal/pkg`: Provides shared utility packages for internal use.
  * `/docs`: Contains project documentation.
  * `/scripts`: Holds utility scripts for development and operational tasks.

## 4. Coding Conventions & Style Guide

* **Formatting:** The project uses `golangci-lint` for linting, as defined in `.golangci.yml`. This enforces standard Go formatting (`gofmt`) along with a variety of other static analysis checks to maintain code quality.
* **Naming Conventions:**
  * `structs`, `interfaces`, exported `functions`: `PascalCase` (e.g., `type School struct`, `func NewNote(...)`).
  * `variables`, unexported `functions`: `camelCase` (e.g., `schoolID`).
  * `files`: `snake_case` (e.g., `school_student.go`).
* **API Design:** (Inferred) An API is likely planned in `internal/infrastructure/http`, but is not yet implemented. It is expected to follow standard RESTful principles.
* **Error Handling:** Standard Go error handling is used, returning an `error` as the last value from functions that can fail. Errors are checked with `if err != nil`.

## 5. Key Files & Entrypoints

* **Main Entrypoint(s):** Not yet implemented. It is expected to be located within a `/cmd` directory in the future.
* **Configuration:** `go.mod` defines the project's module path and dependencies. `.golangci.yml` contains the configuration for the linter.
* **CI/CD Pipeline:** No CI/CD pipeline is configured at this time.

## 6. Development & Testing Workflow

* **Local Development Environment:** The `Makefile` provides commands for common development tasks. A pre-commit hook is available (`scripts/install-hooks.sh`) to automatically run the linter before each commit, ensuring code quality.
* **Testing:** Tests are run using the `make test` command, which executes `go test ./...`. All new features should be accompanied by corresponding unit tests.
* **Linting:** Run the linter using the `make lint` command. This is also integrated into a pre-commit hook.

## 7. Specific Instructions for AI Collaboration

* **Contribution Guidelines:** No formal `CONTRIBUTING.md` exists. However, any new code should adhere to the established architectural patterns, include unit tests, and pass the linter checks (`make lint`).
* **Infrastructure (IaC):** No Infrastructure as Code (IaC) is currently part of the project.
* **Security:** Be mindful of security best practices. Do not hardcode secrets or keys. Sanitize inputs and handle errors gracefully.
* **Dependencies:** To add a new dependency, use `go get`. After adding or updating dependencies, run `make tidy` to ensure `go.mod` and `go.sum` are consistent.
* **Commit Messages:** Commit messages should be concise, lowercase, and written in the imperative mood. They should describe the change clearly. Example: `feat: add user authentication service`.
