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
We use golangci-lint as our primary linter. The golangci-lint project strongly recommends binary installation over installing from sources for performance and compatibility reasons. Follow [their binary installation instructions](https://golangci-lint.run/docs/welcome/install/#binaries) for your platform.

**Other development tools:**
Additional tools like `goimports` and `gofumpt` are managed using [Go 1.24's tool dependency management feature](https://tip.golang.org/doc/go1.24#tools). These tools are already declared in the `go.mod`.

For more information about tool management, see [Go Development Guide](golang.md).

## Start and stop services with Docker Compose

This project uses Docker Compose to manage its services (PostgreSQL and Redis). To start the services, run:

```bash
docker compose up -d
```

To stop the services, run:

```bash
docker compose down
```

## Access the PostgreSQL database

You can access the PostgreSQL database using the following command:

```bash
PGPASSWORD=$POSTGRES_PASSWORD docker exec -it go-arch-patterns-postgres-1 psql -U $POSTGRES_USERNAME -d $POSTGRES_DBNAME
```

This command uses the environment variables loaded by direnv to connect to the database.
