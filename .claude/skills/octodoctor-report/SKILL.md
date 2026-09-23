---
name: octodoctor-report
description: Build/locate the `octodoctor` binary, run it against a directory of Argo CD manifests or export live resources from a Kubernetes cluster's Argo CD install and analyze those, then summarize the resulting best-practices report (pass/fail/inactive rules, severities, doc links) and point to the generated HTML report. Use this whenever the user asks to analyze, audit, review, or health-check an Argo CD installation, cluster, namespace or manifest directory against best practices — including phrases like "check my Argo CD setup", "run octodoctor", "audit my Argo CD cluster", "is my Argo CD config following best practices", "scan these Argo CD manifests", or "export Argo CD from context X and analyze it" — even when the user never says "octodoctor" by name.
---

# octodoctor best-practices report

`octodoctor` is a static analyzer for Argo CD manifests (Applications, ApplicationSets, AppProjects,
the `argocd-*` ConfigMaps, and cluster Secrets). It never mutates anything — the only network-facing
part is `--export`, which reads resources out of a live cluster and writes them to disk before the
same file-based analysis runs on them.

This skill has two entry points. Pick based on what the user actually has:

- **A directory of manifests already on disk** → directory mode (`-d`).
- **A running cluster** they want checked → export mode (`-e`), which pulls the manifests first.

Don't ask the user to pick a mode explicitly if their request already implies one ("check this folder"
vs. "check my cluster") — infer it and confirm briefly if genuinely ambiguous.

## Step 1: locate or build the binary

Try `octodoctor --version` first. If that resolves, use it as-is — don't rebuild it.

If it's not on PATH, check whether the current working directory is the `octodoctor-argocd` source
checkout (a `Makefile` with an `octodoctor-build` target, or a `main.go` importing `pkg/analysis`).
If so, run `make octodoctor-build` and use the resulting `./octodoctor` binary at the repo root.

Just try `make octodoctor-build` first — on a normal developer machine it works fine and none of this
applies. It only fails in a *sandboxed* Claude Code session where Go's default build cache
(`~/.cache/go-build`) and/or module cache (`~/go/pkg/mod`) have been made read-only by the sandbox. If
you hit that specific failure, don't give up on building — redirect the caches to somewhere writable
and build directly with the same flags the Makefile uses:

```
GOCACHE="$TMPDIR/gocache" GOMODCACHE="$TMPDIR/gomodcache" CGO_ENABLED=0 go build -buildvcs=false \
  -ldflags "-X main.version=$(git describe --tags --always 2>/dev/null || echo dev)" -o octodoctor .
```

If neither PATH nor a source checkout is available, stop and tell the user plainly: you need either
`octodoctor` on PATH or a checkout of its source to build from — don't try to fetch or install it
another way without asking.

## Step 2: directory mode

```
octodoctor -d <path-to-manifests>
```

Defaults to `.` (cwd) if `-d` is omitted. The directory must already exist. Never pass `--export`,
`--context`, or `--namespace` alongside `-d` — the tool refuses `-d` combined with `-e`, and silently
ignores `-c`/`-n` without `-e` set (with a warning), so mixing them just confuses the run.

## Step 3: cluster export mode

```
octodoctor -e -c <kube-context> -n <namespace>
```

- `-c/--context` defaults to whatever `kubectl config current-context` currently resolves to. Only
  ask the user which context to use if there are multiple plausible candidates (e.g. they mentioned a
  cluster name that doesn't match the current context) — otherwise just use the current one and say
  which one you used.
- `-n/--namespace` defaults to `argocd`. Only ask if the user's Argo CD is known to live somewhere else.
- This talks to a real cluster over the network and needs a working kubeconfig — if `kubectl` can't
  reach the cluster or the sandbox blocks the connection, report that plainly rather than working
  around it.
- Export only *reads* the five Argo CD resource kinds via `kubectl`/client-go; it writes them to a
  local `argocd-manifests-export-<context>/` directory. Secret **values** are replaced with a
  placeholder before analysis ever sees them — only key names are kept, because that's all the rules
  need. It's safe to mention this to the user if they're wondering whether secrets get exposed.
- A pre-existing export directory with the same name is deleted and recreated — if the user has one
  they care about from a previous run they want to keep, warn them before re-running.

## Step 4: run it and read the exit code

`octodoctor` exits **1** if any rule failed, **0** otherwise. This is expected, normal output for a
real installation — it is not a tool failure and you should not treat a nonzero exit code as an error
to work around. Just proceed to reading the report either way.

The report lands in `<target-directory>/octodoctor-output/`:
- `index.html` — the report the user will actually want to view.
- `report.json` — same data, structured, for you to summarize from.
- `report.yaml` — same data again, YAML.

(`<target-directory>` is the `-d` path in directory mode, or the export directory `octodoctor` printed
in export mode.)

## Step 5: summarize the report

Run `python3 scripts/summarize_report.py <target-directory>/octodoctor-output/report.json` (paths are
relative to this skill's own directory) and use its output
to write a short summary in chat: total rules failed/passed/inactive, and the failed rules ordered by
severity with what they mean and which files triggered them (the script already sorts and truncates
this for you — don't re-derive it by hand-parsing the JSON).

Then tell the user where `index.html` is so they can open it themselves — this skill doesn't need a
browser to do its job, so don't try to open one unless the user separately asks you to. If they ask
"can you fix these", switch to actually editing the flagged manifests to address the failures, but
don't do that unprompted.
