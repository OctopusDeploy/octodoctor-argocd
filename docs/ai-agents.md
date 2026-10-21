# AI agents

OctoDoctor ships an [agent skill](https://agentskills.io/specification) called `octodoctor-report`
in the `.claude/skills/octodoctor-report/` directory of its
[GitHub repository](https://github.com/OctopusDeploy/octodoctor-argocd). With the skill installed, your coding agent (Claude Code, GitHub
Copilot, Cursor, Codex and others) can run OctoDoctor for you and summarize the report in chat.

## Installing the skill

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

## Using the skill

The skill uses `octodoctor` from your `PATH` if it's there. If not, it downloads the latest
[release](https://github.com/OctopusDeploy/octodoctor-argocd/releases) for your operating system
(Linux, macOS or Windows) into a temporary directory and verifies its checksum. The download needs
`curl` and `tar` (PowerShell on Windows) and access to github.com. Nothing else needs to be installed.

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
