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

# Build and archive every release target (needs GoReleaser)
make snapshot
```

`make snapshot` runs GoReleaser exactly as CI does, and writes archives plus a `checksums.txt` into
`dist/` without publishing anything. It covers:
- macOS (AMD64 and ARM64)
- Linux (AMD64 and ARM64)
- Windows (AMD64)

Install GoReleaser with `go install github.com/goreleaser/goreleaser/v2@latest`, and run
`make release-check` to validate `.goreleaser.yaml` without building.

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
2. Runs [GoReleaser](https://goreleaser.com), which cross-compiles all five targets and packages each
   one as a `.tar.gz` (`.zip` on Windows) archive, alongside a `checksums.txt`.
3. Publishes a GitHub Release named `Release <tag>` with those six assets attached.

The release body is assembled by GoReleaser from the commit and pull request history since the previous
tag, followed by the static footer defined in `.goreleaser.yaml` — that file is where you change either
of them, not the workflow.

The workflow can also be triggered manually via `workflow_dispatch` (the "Run workflow" button in the
Actions tab). A manual run has no tag to release, so it builds a snapshot and publishes nothing, which
makes it a safe dry run.

The build matrix, archive layout, version stamping and release notes all live in `.goreleaser.yaml`.
Note that `main.version` is injected at build time and is deliberately `v`-prefixed, because it is
embedded in every `report.json` and `report.yaml`; changing its format changes that output.
