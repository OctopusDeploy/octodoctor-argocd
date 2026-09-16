# Installation

!!! note "This page is a placeholder"
    TODO: flesh out the installation instructions, including the exact release asset names and any
    checksum verification steps.

OctoDoctor for Argo CD is a single self-contained binary with no runtime dependencies. It analyses YAML
files from a directory you give it, and can also
[export those files from a live cluster](export.md) for you first.

## Download a prebuilt binary

Grab the archive for your platform from the
[releases page](https://github.com/OctopusDeploy/octodoctor-argocd/releases), then make it executable and
put it somewhere on your `PATH`. Binaries are published for Linux, macOS and Windows on both `amd64`
and `arm64`.

## Build from source

You need Go 1.24.0 or later.

```bash
git clone https://github.com/OctopusDeploy/octodoctor-argocd.git
cd octodoctor-argocd
make octodoctor-build
```

This produces an `octodoctor` binary in the repository root. To build every supported platform at once,
run `make release`.

## Verify the installation

```bash
./octodoctor --directory /path/to/argocd/manifests
```

The report is written to an `octodoctor-output` directory inside the directory you scanned. See the
[rules reference](rules/index.md) for what gets checked.

## Command line options

| Option | Default | What it does |
|---|---|---|
| `-d`, `--directory` | current directory | Directory to scan for Argo CD files |
| `-e`, `--export` | off | [Export the manifests of a live cluster](export.md) and analyse them |
| `-c`, `--context` | the current kubeconfig context | Kubeconfig context to export from, only used with `--export` |
| `-n`, `--namespace` | `argocd` | Namespace to export from, only used with `--export` |
| `-l`, `--limit` | `5` | How many files to list under each rule in the HTML report |
| `-v`, `--version` | | Print the version and exit |

OctoDoctor exits with `1` when at least one rule failed, and `0` otherwise.
