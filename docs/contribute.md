# Contribute

## Prerequisites for development

- Go 1.26.10 or later
- Make (for using the Makefile)
- `golangci-lint` (for code linting)

## Build commands

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

## Running unit tests

```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...
```

## Code quality and linting

```bash
# Run linter with auto-fix
make lint

# Manual linting
golangci-lint run --fix
```

## Adding new rules

To add a new validation rule:

1. Create a new file `pkg/rules/rXX.go` (where XX is the rule number) that implements the `Rule`
   interface. `Metadata()` returns the ID, the `DocSlug`, the pass and fail descriptions, and the
   severity, and `Execute` runs the check.
2. Add the rule to the registry in `pkg/rules/registry.go`.
3. Create test fixtures in `testdata/rXX/`, covering the pass, fail, and inactive cases.
4. Write a table-driven test in `pkg/rules/rXX_test.go`.
5. Add a `docs/rules/rXX-<slug>.md` page, a row in the table in `docs/rules/index.md`, and a `nav`
   entry in `mkdocs.yml`.

## Releasing a new version

To release a new version, push a Git tag matching `v*` (for example, `v1.2.3`):

```bash
git tag v1.2.3
git push origin v1.2.3
```

This triggers the [`release.yml`](https://github.com/OctopusDeploy/octodoctor-argocd/blob/main/.github/workflows/release.yml)
GitHub Actions workflow, which:

1. Verifies the tagged commit is an ancestor of `main`. The workflow fails if the tag points to a
   commit that hasn't been merged into `main`.
2. Runs [GoReleaser](https://goreleaser.com), which cross-compiles all 5 targets and packages each
   one as a `.tar.gz` (`.zip` on Windows) archive, alongside a `checksums.txt`.
3. Publishes a GitHub Release named `Release <tag>` with those 6 assets attached.

GoReleaser assembles the release body from the commit and pull request history since the previous
tag, followed by the static footer defined in `.goreleaser.yaml`. To change either of them, edit that
file, not the workflow.

You can also trigger the workflow manually with `workflow_dispatch` (the **Run workflow** button in the
Actions tab). A manual run has no tag to release, so it builds a snapshot and publishes nothing, which
makes it a safe dry run.

The build matrix, archive layout, version stamping, and release notes all live in `.goreleaser.yaml`.
The build injects `main.version` and deliberately gives it a `v` prefix, because it's embedded in
every `report.json` and `report.yaml`. Changing its format changes that output.
