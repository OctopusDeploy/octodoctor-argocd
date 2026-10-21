# Usage

OctoDoctor for Argo CD analyzes Argo CD manifests in one of two ways:

- **Manifests on disk**: point it at a directory that holds your Argo CD YAML files, for example a Git
  checkout, with `--directory`.
- **A live cluster**: let it export the manifests from a running Argo CD installation first, with
  `--export`, and then analyze them.

In both modes the tool is read-only (doesn't change anything in your files or cluster) and produces the same result.

## Command-line options

| Option | Default | What it does |
|---|---|---|
| `-d`, `--directory` | current directory | Directory to scan for Argo CD files |
| `-e`, `--export` | off | [Export the manifests of a live cluster](#mode-2-analyzing-a-live-cluster) and analyze them |
| `-c`, `--context` | the current kubeconfig context | Kubeconfig context to export from, only used with `--export` |
| `-n`, `--namespace` | `argocd` | Namespace to export from, only used with `--export` |
| `-l`, `--limit` | `5` | How many files to list under each rule in the HTML report |
| `-v`, `--version` | | Print the version and exit |

!!! note "`--export` and `--directory` are mutually exclusive"
    `--export` analyzes a cluster and `--directory` analyzes a local folder, so passing both is an
    error. If you pass neither, OctoDoctor analyzes the current directory. Without `--export`,
    OctoDoctor ignores `--context` and `--namespace` and logs a line saying so.

OctoDoctor exits with `1` when at least one rule failed, and `0` otherwise, so you can use it to fail a CI
pipeline.

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

Sometimes you only have access to an Argo CD cluster. Instead of first exporting the manifests yourself with
[argocd admin export](https://argo-cd.readthedocs.io/en/latest/user-guide/commands/argocd_admin_export/) or [kubexporter](https://github.com/bakito/kubexporter) and then running OctoDoctor on the result, use `--export` to do
both in one step:

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
`--directory`. The report ends up inside the exported directory

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

OctoDoctor does **not** delete the exported files when the analysis finishes. They stay on disk so you can read
them, commit them, or analyze them again later with `--directory`. If the directory already exists from
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
| `Secret` | only `argocd-secret` and the external cluster secrets labelled `argocd.argoproj.io/secret-type: cluster` |

OctoDoctor ignores every other ConfigMap and Secret in the namespace.

OctoDoctor saves the files in the simplest form possible. It removes everything that only describes the live object in the
cluster, so what is left looks like a manifest somebody would commit to Git: `status`, the
`metadata` fields `uid`, `resourceVersion`, `generation`, `creationTimestamp`, `selfLink`,
`managedFields`, `deletionTimestamp`, `deletionGracePeriodSeconds`, `ownerReferences`, and `finalizers`,
and the `kubectl.kubernetes.io/last-applied-configuration` annotation.

!!! warning "Secret values are redacted"
    OctoDoctor keeps the **keys** of a Secret but replaces every **value** with
    `REDACTED-BY-OCTODOCTOR`, so no password, token, or private key is ever written to your disk. The
    rules only ever look at the key names of a Secret, so this doesn't change any result.

Even with the Secret values redacted, the exported manifests still describe how your Argo CD is
configured, for example, your RBAC policy and your single sign-on settings. That's why only you can read the
exported directory and the report inside it. Copy the report somewhere else if you
want to share it.

### Permissions

OctoDoctor only reads from your cluster and never writes to it. It needs `list` permission on Applications,
ApplicationSets, AppProjects, ConfigMaps and Secrets in the namespace you point it at.

If you aren't allowed to list one of these, or the Argo CD CRDs aren't installed on the cluster at
all, OctoDoctor logs a warning, skips that resource type and carries on with the rest. Only a cluster it
can't reach, or credentials it can't authenticate with, stop the run.
