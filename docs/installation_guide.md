# Installation & Setup

## Clone the repository

```bash
git clone https://github.com/jp-ryuji/go-sample.git
cd go-sample
```

## Set up git hooks

This project includes a pre-commit hook that automatically runs the Go linter before each commit. If the linter finds any issues, the commit will be aborted, and you'll need to fix the issues before committing again.

To install the pre-commit hook, run:

```bash
./scripts/install-hooks.sh
```

Note: Git hooks are local to your repository and are not shared through the repository. Each developer needs to run this script to install the hooks in their local environment.

For more information about Go development, see [Go Development Guide](golang.md).

## Install necessary tools

This project uses two categories of development tools:

**Linting tools (golangci-lint):**
We use golangci-lint as our primary linter. The golangci-lint project strongly recommends binary installation over installing from sources for performance and compatibility reasons. Follow [their binary installation instructions](https://golangci-lint.run/docs/welcome/install/#binaries) for your platform.

**Other development tools:**
Additional tools like `goimports` and `gofumpt` are managed using [Go 1.24's tool dependency management feature](https://tip.golang.org/doc/go1.24#tools). These tools are already declared in the `go.mod`.

For more information about tool management, see [Go Development Guide](golang.md).
