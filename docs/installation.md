# Installation

OctoDoctor for Argo CD is a single self-contained binary with no runtime dependencies. It analyzes YAML
files from a directory you give it, and can also
[export those files from a live cluster](usage.md#mode-2-analyzing-a-live-cluster) for you first.

## Download a prebuilt binary

Download the archive for your platform from the
[releases page](https://github.com/OctopusDeploy/octodoctor-argocd/releases), extract it, and put the
`octodoctor` binary somewhere on your `PATH`. Archives are named
`octodoctor_<version>_<os>_<arch>.tar.gz`, or `.zip` on Windows, and are published for Linux and macOS
on both `amd64` and `arm64`, and for Windows on `amd64`.

```bash
tar -xzf octodoctor_<version>_linux_amd64.tar.gz
chmod +x octodoctor
./octodoctor --version
```

Every release also ships a `checksums.txt`. To verify a download before you trust it:

```bash
sha256sum --ignore-missing -c checksums.txt   # Linux
shasum -a 256 --ignore-missing -c checksums.txt   # macOS
```

## Build from source

You need Go 1.26.10 or later.

```bash
git clone https://github.com/OctopusDeploy/octodoctor-argocd.git
cd octodoctor-argocd
make octodoctor-build
```

This produces an `octodoctor` binary in the repository root. To cross-compile and archive every
release target at once, run `make snapshot`. This needs [GoReleaser](https://goreleaser.com) installed,
and writes the archives to `dist/` exactly as the release workflow does.

## Verify the installation

```bash
./octodoctor --directory /path/to/argocd/manifests
```

OctoDoctor writes the report to an `octodoctor-output` directory inside the directory you scanned. See
[Usage](usage.md) for every command-line option, and the [rules reference](rules/index.md) for what gets
checked.

