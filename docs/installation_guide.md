# Installation & Setup

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

For more information about Go development, see [Go Development Guide](golang.md).