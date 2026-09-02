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

## Releasing a New Version

Releases are triggered by pushing a Git tag matching `v*` (e.g. `v1.2.3`):

```bash
git tag v1.2.3
git push origin v1.2.3
```

This triggers the [`release.yml`](https://github.com/OctopusDeploy/octodoctor-argocd/blob/main/.github/workflows/release.yml)
GitHub Actions workflow, which:

1. Verifies the tagged commit is an ancestor of `main` — the workflow fails if the tag points to a
   commit that hasn't been merged into `main`.
2. Runs `make release` to build binaries for macOS (AMD64/ARM64), Linux (AMD64/ARM64), and Windows
   (AMD64).
3. Publishes a GitHub Release named `Release <tag>` with all five binaries attached.

The workflow can also be triggered manually via `workflow_dispatch` (the "Run workflow" button in the
Actions tab) for testing, without pushing a tag.
