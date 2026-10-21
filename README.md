# OctoDoctor for Argo CD

![OctoDoctor for Argo CD logo](docs/octodoctor-logo.png)

A comprehensive static analysis tool for Argo CD configurations that helps you identify potential issues, misconfigurations, and anti-patterns in your GitOps setup.

See available [rules](https://octodoctor-argocd.readthedocs.io/en/latest/rules/).

## Example report

![OctoDoctor for Argo CD report](docs/example-report.png)


The tool analyzes various Argo CD resource types including:
- **Applications** (`argoproj.io/v1alpha1/Application`)
- **ApplicationSets** (`argoproj.io/v1alpha1/ApplicationSet`) 
- **AppProjects** (`argoproj.io/v1alpha1/AppProject`)
- **ConfigMaps, Secrets** (Argo CD specific ConfigMaps like `argocd-cm`, `argocd-cmd-params-cm`, etc.)

## How to Run OctoDoctor

### Prerequisites

- A directory containing Argo CD YAML configuration files. The tool also will automatically list all files recursively 
- The tool works great when using as input Argo CD application manifests as yaml files
- Alternatively, a live cluster you can reach with your kubeconfig, and the tool will export the manifests for you

### Running 

Point the tool at the directory to scan:

```bash
# Scan current directory
./octodoctor

# Scan specific directory
./octodoctor -d /path/to/argocd/manifests
./octodoctor --directory /path/to/argocd/manifests
```

If you only have a cluster and no manifests, let the tool export them first and analyze them in
the same step:

```bash
# Export from the current kubeconfig context, namespace argocd
./octodoctor -e
./octodoctor --export

# Export from a specific context and namespace
./octodoctor -e -c production -n my-argocd
./octodoctor --export --context production --namespace my-argocd
```

The manifests are saved into `argocd-manifests-export-<context>` and kept there after the analysis.
Secret values are replaced with a placeholder, so no credentials are written to disk. See the
[export documentation](docs/export.md) for the details.

### Command line options

| Option | Default | What it does |
|---|---|---|
| `-d`, `--directory` | current directory | Directory to scan for Argo CD files |
| `-e`, `--export` | off | Export the manifests of a live cluster and analyze them |
| `-c`, `--context` | the current kubeconfig context | Kubeconfig context to export from, only used with `--export` |
| `-n`, `--namespace` | `argocd` | Namespace to export from, only used with `--export` |
| `-l`, `--limit` | `5` | How many files to list under each rule in the HTML report |
| `-v`, `--version` | | Print the version and exit |

The tool exits with `1` when at least one rule failed, and `0` otherwise.

## Using OctoDoctor from an AI agent

This repository ships an [agent skill](https://agentskills.io/specification) called `octodoctor-report`
in `.claude/skills/octodoctor-report/`. With the skill installed, your coding agent (Claude Code, GitHub
Copilot, Cursor, Codex and others) can run OctoDoctor for you and summarize the report in chat.

### Installing the skill

If you work inside a clone of this repository, Claude Code picks up the skill automatically and you
don't need to install anything. To use the skill in other projects, install it with one of the
following commands.

With `npx` (uses the [skills](https://github.com/vercel-labs/skills) CLI):

```bash
# Install into the current project
npx skills add OctopusDeploy/octodoctor-argocd --skill octodoctor-report

# Install for your user, available in every project
npx skills add OctopusDeploy/octodoctor-argocd --skill octodoctor-report -g
```

With the GitHub CLI (`gh skill` is in preview):

```bash
# Install into the current project
gh skill install OctopusDeploy/octodoctor-argocd octodoctor-report --allow-hidden-dirs --agent claude-code

# Install for your user, available in every project
gh skill install OctopusDeploy/octodoctor-argocd octodoctor-report --allow-hidden-dirs --agent claude-code --scope user
```

The `--allow-hidden-dirs` flag is required because the skill lives under `.claude/`. Replace
`claude-code` with your agent, for example `github-copilot`, `cursor` or `codex`. Run
`gh skill update --all` later to get newer versions of the skill.

### Using the skill

The skill needs the `octodoctor` binary on your `PATH`. Inside a clone of this repository it builds the
binary itself with `make octodoctor-build`.

Ask your agent in plain language. You don't have to mention the skill or OctoDoctor by name:

- "Check the Argo CD manifests in `./gitops` against best practices"
- "Audit the Argo CD installation in my `production` cluster"
- "Export Argo CD from context `staging`, namespace `my-argocd`, and analyze it"

The agent picks directory mode (`-d`) or cluster export mode (`-e`) based on your request, runs the
tool, and replies with:

- The number of failed, passed and inactive rules
- The failed rules ordered by severity, with the files that triggered them
- The path to the generated `index.html` report

A failed rule makes the tool exit with `1`. The skill treats this as a normal result, not an error. The
agent only edits your manifests to fix the failures if you ask it to.
