#!/bin/bash

# Install pre-commit hook
echo "Installing pre-commit hook..."

# Create the hooks directory if it doesn't exist
mkdir -p .git/hooks

# Copy the pre-commit hook script
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/sh

# Run go linter before commit
echo "Running go linter..."
make lint.go

# If make fails, exit with error code
if [ $? -ne 0 ]; then
  echo "Linting failed. Please fix the issues before committing."
  echo "You can run 'make lint.go.fix' to automatically fix some issues."
  exit 1
fi

echo "Linting passed. Continuing with commit..."
exit 0
EOF

# Make the hook executable
chmod +x .git/hooks/pre-commit

echo "Pre-commit hook installed successfully!"
