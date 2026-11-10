# Usage

OctoDoctor for Argo CD analyzes Argo CD manifests in one of two ways:

- **Manifests on disk**: point it at a directory that holds your Argo CD YAML files, for example a Git
  checkout, with `--directory`.
- **A live cluster**: let it export the manifests from a running Argo CD installation first, with
  `--export`, and then analyze them.

In both modes, the tool is read-only (it doesn't change anything in your files or cluster) and
produces the same result.

## Command-line options

| Option | Default | What it does |
|---|---|---|
| `-d`, `--directory` | current directory | Directory to scan for Argo CD files |
| `-e`, `--export` | off | [Export the manifests of a live cluster](#mode-2-analyzing-a-live-cluster) and analyze them |
| `-c`, `--context` | the current kubeconfig context | Kubeconfig context to export from, only used with `--export` |
| `-n`, `--namespace` | `argocd` | Namespace to export from, only used with `--export` |
| `-l`, `--limit` | `5` | How many files to list under each rule in the HTML report |
| `--fail-on` | `medium` | Lowest severity of a failed rule that makes OctoDoctor exit with `1`: `high`, `medium` or `none` |
| `--config` | `.octodoctor.yaml` in the current directory, if it exists | [Configuration file with the failures to ignore](#ignoring-rules) |
| `-v`, `--version` | | Print the version and exit |

!!! note "`--export` and `--directory` are mutually exclusive"
    `--export` analyzes a cluster and `--directory` analyzes a local folder, so passing both is an
    error. If you pass neither, OctoDoctor analyzes the current directory. Without `--export`,
    OctoDoctor ignores `--context` and `--namespace` and logs a line saying so.

OctoDoctor exits with `1` when at least one rule failed, and `0` otherwise, so you can use it to fail a CI
pipeline. Use `--fail-on` to choose which failures count:

- `medium` (default): any failed rule exits with `1`.
- `high`: only failed high severity rules exit with `1`. Medium failures still appear in the reports.
- `none`: always exits with `0`. Use this to try OctoDoctor in a pipeline before you make it blocking.

```bash
./octodoctor -d /path/to/argocd/manifests --fail-on high
```

## Mode 1 - Analyzing manifests on disk

Pass the directory that holds your Argo CD manifests. OctoDoctor searches it and all of its
subdirectories for `.yaml` and `.yml` files:

```bash
./octodoctor --directory /path/to/argocd/manifests
./octodoctor -d /path/to/argocd/manifests

# Analyze the current directory
./octodoctor
```

OctoDoctor only looks at the files that describe Argo CD: Applications, ApplicationSets, AppProjects, the
Argo CD ConfigMaps and Secrets listed under [what gets exported](#what-gets-exported), and cluster
Secrets. It ignores every other file, so you can point it at a repository that also contains other
Kubernetes manifests.

The report is written to an `octodoctor-output` directory inside the directory you scanned:

```
/path/to/argocd/manifests/
└── octodoctor-output/
    ├── index.html     # open this in a browser
    ├── report.json
    └── report.yaml
```

`index.html` lists at most `--limit` files under each rule. `report.json` and `report.yaml` always contain
every file, which makes them the better choice for scripts.

```bash
./octodoctor -d /path/to/argocd/manifests --limit 20
```

## Mode 2 - Analyzing a live cluster

Sometimes you only have access to an Argo CD cluster. Instead of first exporting the manifests
yourself with
[argocd admin export](https://argo-cd.readthedocs.io/en/latest/user-guide/commands/argocd_admin_export/)
or [kubexporter](https://github.com/bakito/kubexporter) and then running OctoDoctor on the result,
use `--export` to do both in one step:

```bash
./octodoctor --export
./octodoctor -e

# A specific context
./octodoctor --export --context production
./octodoctor --export -c production

# A specific namespace
./octodoctor --export --namespace my-argocd
./octodoctor --export -n my-argocd
```

OctoDoctor connects to your cluster with the standard kubeconfig, saves the Argo CD manifests it cares
about into a local directory, and then analyzes that directory exactly as if you had passed it with
`--directory`. The report ends up inside the exported directory.

If your kubeconfig has no current context, OctoDoctor stops and asks you to choose one with `--context`.

### Where the files go

OctoDoctor writes the manifests to `argocd-manifests-export-<context>` in your current directory, with one
subfolder per resource type and one file per resource:

```
argocd-manifests-export-production/
├── applications/
│   ├── team-a-billing.yaml
│   └── team-b-billing.yaml
├── applicationsets/
│   └── billing-application-set.yaml
├── appprojects/
│   └── default.yaml
├── configmaps/
│   ├── argocd-cm.yaml
│   └── argocd-rbac-cm.yaml
├── secrets/
│   └── argocd-secret.yaml
└── octodoctor-output/
    ├── index.html
    ├── report.json
    └── report.yaml
```

OctoDoctor replaces characters that aren't safe in a directory name with a dash, so an EKS context called
`arn:aws:eks:eu-west-1:123456789:cluster/prod` becomes
`argocd-manifests-export-arn-aws-eks-eu-west-1-123456789-cluster-prod`.

OctoDoctor does **not** delete the exported files when the analysis finishes. They stay on disk so
you can read them, commit them, or analyze them again later with `--directory`. If the directory already exists from
an earlier export, OctoDoctor deletes it and creates it again, so that manifests of resources that no
longer exist in the cluster are never analyzed. A log line tells you when this happens.

### What gets exported

Only the resources OctoDoctor already knows how to analyze:

| Resource type | What is exported |
|---|---|
| `Application` | all of them, including the ones generated by an ApplicationSet |
| `ApplicationSet` | all of them |
| `AppProject` | all of them |
| `ConfigMap` | only the Argo CD settings: `argocd-cm`, `argocd-cmd-params-cm`, `argocd-rbac-cm`, `argocd-tls-certs-cm`, and `argocd-ssh-known-hosts-cm` |
| `Secret` | only `argocd-secret` and the external cluster secrets labeled `argocd.argoproj.io/secret-type: cluster` |

OctoDoctor ignores every other ConfigMap and Secret in the namespace.

OctoDoctor saves the files in the simplest form possible. It removes everything that only describes
the live object in the cluster, so what's left looks like a manifest somebody would commit to Git: `status`, the
`metadata` fields `uid`, `resourceVersion`, `generation`, `creationTimestamp`, `selfLink`,
`managedFields`, `deletionTimestamp`, `deletionGracePeriodSeconds`, `ownerReferences`, and `finalizers`,
and the `kubectl.kubernetes.io/last-applied-configuration` annotation.

!!! warning "Secret values are redacted"
    OctoDoctor keeps the **keys** of a Secret but replaces every **value** with
    `REDACTED-BY-OCTODOCTOR`, so no password, token, or private key is ever written to your disk. The
    rules only ever look at the key names of a Secret, so this doesn't change any result.

Even with the Secret values redacted, the exported manifests still describe how your Argo CD is
configured, for example, your RBAC policy and your single sign-on settings. That's why only
you can read the exported directory and the report inside it. Copy the report somewhere else if you
want to share it.

### Permissions

OctoDoctor only reads from your cluster and never writes to it. It needs `list` permission on Applications,
ApplicationSets, AppProjects, ConfigMaps, and Secrets in the namespace you point it at.

If you aren't allowed to list one of these, or the Argo CD CRDs aren't installed on the cluster at
all, OctoDoctor logs a warning, skips that resource type, and carries on with the rest. Only a cluster it
can't reach, or credentials it can't authenticate with, stop the run.

## Ignoring rules

Some failures are deliberate decisions, for example a component pinned to a Git commit on purpose. To
stop reporting them as issues, list them in a `.octodoctor.yaml` file:

```yaml
ignore:
  # Ignore a whole rule
  - rule: R99
    reason: One ApplicationSet per application is our convention

  # Ignore a rule only for some files
  - rule: R19
    resources:
      - "applications/*mesh-*"
    reason: Pinned on purpose, upgraded one cluster at a time
```

| Field | Required | What it does |
|---|---|---|
| `rule` | yes | ID of the rule, for example `R19` |
| `resources` | no | Paths of the files to ignore, exactly as the report lists them under the rule. Leave it out to ignore the whole rule |
| `reason` | no | Text shown next to the ignored failure in the reports |

In `resources`, `*` matches any characters, including `/`, and a pattern must match the whole path.
Paths differ between the two modes: `--directory` lists them relative to the scanned directory, and
`--export` lists them as `applications/<name>.yaml`. A pattern such as `*mesh-*` matches in both.

OctoDoctor reads `.octodoctor.yaml` from the current directory, or the file you pass with `--config`.
It never looks inside the scanned directory, because `--export` creates that directory again on every
run. With `--export`, run OctoDoctor outside your Git checkout, so the exported manifests don't end up
in it, and point `--config` at the file in the checkout:

```bash
./octodoctor --export --config ~/repos/argocd-config/.octodoctor.yaml
```

Ignored failures:

- appear in an **Ignored failures** section of `index.html`, and under `results.ignored` in
  `report.json` and `report.yaml`, so nothing is hidden;
- don't count towards `--fail-on`;
- only apply to failures. A rule that passes is still reported as passed.

If a rule fails on more files than you ignored, it stays failed with the remaining files.

OctoDoctor stops with an error when the file has an unknown key, an unknown rule ID, an empty
`resources` list, or ignores a rule both as a whole and by resources. It logs a warning for every
entry or pattern that ignored nothing, so you notice ignores that are no longer needed.

## Running in GitHub Actions

Use the OctoDoctor action to check the Argo CD manifests in your repository on every push or pull
request. The action downloads the OctoDoctor release that matches the action version, and it only
runs on Linux runners.

```yaml
name: OctoDoctor

on:
  pull_request:

jobs:
  octodoctor:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7

      - uses: OctopusDeploy/octodoctor-argocd@v0.5.0
        with:
          directory: argocd
          fail-on: high

      - uses: actions/upload-artifact@v4
        if: always()
        with:
          name: octodoctor-report
          path: argocd/octodoctor-output/
```

| Input | Default | What it does |
|---|---|---|
| `directory` | `.` | Directory to scan for Argo CD files |
| `fail-on` | `medium` | Lowest severity of a failed rule that fails the step: `high`, `medium` or `none` |
| `version` | the action version | OctoDoctor release to download, for example `v0.5.0` |

The `if: always()` line uploads the report even when the OctoDoctor step fails, so you can see which
rules failed.
