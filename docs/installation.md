# Installation

OctoDoctor for Argo CD is a single self-contained binary with no runtime dependencies. It analyses YAML
files from a directory you give it, and can also
[export those files from a live cluster](export.md) for you first.

## Download a prebuilt binary

Grab the archive for your platform from the
[releases page](https://github.com/OctopusDeploy/octodoctor-argocd/releases), extract it, and put the
`octodoctor` binary somewhere on your `PATH`. Archives are named
`octodoctor_<version>_<os>_<arch>.tar.gz`, or `.zip` on Windows, and are published for Linux and macOS
on both `amd64` and `arm64`, and for Windows on `amd64`.

```bash
tar -xzf octodoctor_<version>_linux_amd64.tar.gz
./octodoctor --version
```

Every release also ships a `checksums.txt`. To verify a download before you trust it:

```bash
sha256sum --ignore-missing -c checksums.txt   # Linux
shasum -a 256 --ignore-missing -c checksums.txt   # macOS
```

## Build from source

You need Go 1.24.0 or later.

```bash
git clone https://github.com/OctopusDeploy/octodoctor-argocd.git
cd octodoctor-argocd
make octodoctor-build
```

This produces an `octodoctor` binary in the repository root. To cross-compile and archive every
release target at once, run `make snapshot`. That needs [GoReleaser](https://goreleaser.com) installed
and writes the archives to `dist/`, exactly as the release workflow does.

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
