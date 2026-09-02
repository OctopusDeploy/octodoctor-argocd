# Installation

!!! note "This page is a placeholder"
    TODO: flesh out the installation instructions, including the exact release asset names and any
    checksum verification steps.

OctoDoctor for Argo CD is a single self-contained binary with no runtime dependencies. It never talks to
a cluster — it only reads YAML files from a directory you give it.

## Download a prebuilt binary

Grab the archive for your platform from the
[releases page](https://github.com/OctopusDeploy/octodoctor-argocd/releases), then make it executable and
put it somewhere on your `PATH`. Binaries are published for Linux, macOS and Windows on both `amd64`
and `arm64`.

## Build from source

You need Go 1.23.0 or later.

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
[rules reference](rules.md) for what gets checked.
