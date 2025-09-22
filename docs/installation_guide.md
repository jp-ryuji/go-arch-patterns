# Installation & Setup

## Clone the repository

```bash
git clone https://github.com/jp-ryuji/go-arch-patterns.git
cd go-arch-patterns
```

## Set up environment variables with direnv

This project uses `direnv` to manage environment variables. Direnv automatically loads environment variables from a `.envrc` file when you enter the project directory.

### Install direnv

First, install direnv following the [official installation instructions](https://direnv.net/docs/installation.html) for your platform.

For macOS with Homebrew:

```bash
brew install direnv
```

For other platforms, please refer to the [direnv documentation](https://direnv.net/docs/installation.html).

### Set up the environment file

Copy the example environment file:

```bash
cp .envrc.example .envrc
```

The `.envrc` file contains environment variables needed for the project, including database configurations.

### Customize the environment variables

You may want to modify the values in `.envrc` to match your local environment. For example, if you have a port conflict, you can change the external port numbers. Similarly, you can update any database credentials or other configuration values as needed for your setup.

Edit these values in your `.envrc` file as needed:

```bash
# Edit the .envrc file with your preferred editor
nano .envrc
```

### Allow direnv to load the environment

After creating and customizing your `.envrc` file, you need to allow direnv to load it:

```bash
direnv allow
```

Direnv will automatically load these environment variables whenever you enter the project directory.

## Set up git hooks

This project includes a pre-commit hook that automatically runs the Go linter before each commit. If the linter finds any issues, the commit will be aborted, and you'll need to fix the issues before committing again.

To install the pre-commit hook, run:

```bash
./scripts/install-hooks.sh
```

Note: Git hooks are local to your repository and are not shared through the repository. Each developer needs to run this script to install the hooks in their local environment.

## Install necessary tools

This project uses two categories of development tools:

**Linting tools (golangci-lint):**
We use `golangci-lint` as our primary linter. The golangci-lint project strongly recommends binary installation over installing from sources for performance and compatibility reasons. Follow [their binary installation instructions](https://golangci-lint.run/docs/welcome/install/#binaries) for your platform.

**Other development tools:**
Additional tools like `goimports` and `gofumpt` are managed using [Go 1.24's tool dependency management feature](https://tip.golang.org/doc/go1.24#tools). These tools are already declared in the `go.mod`.

For more information about tool management, see [Go Development Guide](golang.md).

## Start and stop services with Docker Compose

This project uses Docker Compose to manage its services (PostgreSQL and Redis) and `air` to start the main go app. The development environment has been separated into two components:

1. Infrastructure services (PostgreSQL and Redis)
2. The application itself with hot reload capabilities

To start the infrastructure services, run:

```bash
make dev.up
```

This command starts the Docker Compose services in the background. After the services are running, start the application with:

```bash
make dev.run
```

To stop the services, run:

```bash
make dev.down
```

The application uses Viper for configuration management, supporting all configuration options via environment variables. See `.envrc.example` for available options.

If you encounter port conflicts when starting the application, you can use the following Makefile commands:

- `make dev.check` - Check if ports 50051 (gRPC) or 8081 (HTTP) are already in use
- `make dev.kill` - Kill processes listening on these ports

### Development Workflow

For an efficient development workflow, you can use these commands in sequence:

1. `make dev.up` - Start infrastructure services (do this once)
2. `make dev.run` - Start the application with hot reload
3. When you need to restart just the application (keeping services running): `make dev.restart`

## Access the PostgreSQL database

You can access the PostgreSQL database using the following command:

```bash
make psql
```
