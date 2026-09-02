# OctoDoctor for Argo CD 🩺

![OctoDoctor for Argo CD logo](docs/octodoctor-logo.png)

A comprehensive static analysis tool for Argo CD configurations that helps you identify potential issues, misconfigurations, and anti-patterns in your GitOps setup.

See available [rules](docs/rules.md).


## Example report

![OctoDoctor for Argo CD report](docs/example-report.png)


The tool analyzes various Argo CD resource types including:
- **Applications** (`argoproj.io/v1alpha1/Application`)
- **ApplicationSets** (`argoproj.io/v1alpha1/ApplicationSet`) 
- **AppProjects** (`argoproj.io/v1alpha1/AppProject`)
- **ConfigMaps, Secrets** (Argo CD specific ConfigMaps like `argocd-cm`, `argocd-cmd-params-cm`, etc.)

## 🚀 How to Run This Tool

### Prerequisites

- A directory containing Argo CD YAML configuration files. The tool also will automatically list all files recursively 
- The tool works great when using as input Argo CD application manifests as yaml files

### Running 

The tool accepts a single parameter specifying the directory to scan:

```bash
# Scan current directory
./octodoctor

# Scan specific directory
./octodoctor -d /path/to/argocd/manifests
./octodoctor --directory /path/to/argocd/manifests
```


When you run OctoDoctor for Argo CD, it will:

1. **📁 Scan for YAML Files**: Recursively discovers all `.yaml` and `.yml` files in the specified directory
2. **🔎 Filter Argo CD Resources**: Identifies and categorizes Argo CD-specific resources
3. **✅ Execute Validation Rules**: Runs all configured rules against your configurations
4. **📊 Generate HTML Report**: Creates a comprehensive report in the `octodoctor-output` directory



## 🛠️ How to contribute


### Prerequisites for Development

- Go 1.23.0 or later
- Make (for using the Makefile)
- `golangci-lint` (for code linting)

### Build Commands

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

### Running Unit Tests

```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...
```

### Code Quality and Linting

```bash
# Run linter with auto-fix
make lint

# Manual linting
golangci-lint run --fix
```

### Adding New Rules

To add a new validation rule:

1. Create a new file `pkg/rules/rXX.go` (where XX is the rule number)
2. Implement the `Rule` interface with an `Execute` method
3. Add the rule to the registry in `pkg/rules/registry.go`
4. Create test data in `testdata/rXX/`
5. Update documentation as needed

