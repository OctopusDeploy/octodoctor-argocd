# Contribute

## Prerequisites for Development

- Go 1.23.0 or later
- Make (for using the Makefile)
- `golangci-lint` (for code linting)

## Build Commands

```bash
# Build for current platform
make octodoctor-build

# Build for specific platform
make BIN_NAME=octodoctor-linux-amd64 GOOS=linux GOARCH=amd64 octodoctor-build

# Build for all supported platforms
make release
```

The `make release` command builds binaries for:
- macOS (AMD64 and ARM64)
- Linux (AMD64 and ARM64) 
- Windows (AMD64)

## Running Unit Tests

```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...
```

## Code Quality and Linting

```bash
# Run linter with auto-fix
make lint

# Manual linting
golangci-lint run --fix
```

## Adding New Rules

To add a new validation rule:

1. Create a new file `pkg/rules/rXX.go` (where XX is the rule number)
2. Implement the `Rule` interface with an `Execute` method
3. Add the rule to the registry in `pkg/rules/registry.go`
4. Create test data in `testdata/rXX/`
5. Update documentation as needed
