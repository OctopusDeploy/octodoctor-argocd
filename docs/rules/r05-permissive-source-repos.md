# R05 – Permissive sourceRepos

!!! danger "Severity: High"

AppProjects have overly permissive sourceRepos settings (using wildcard '*'). This lets applications deploy from any Git repository, which poses a security risk.

## Why this matters

AppProjects are the centerpiece of Argo CD security in a multi-tenant installation. `sourceRepos`
is the part that answers "where may manifests come from?", and `'*'` answers "anywhere".

With a bare wildcard, anyone who can create an Application in the project can point it at a
repository you've never reviewed, including one they control, and Argo CD will render whatever it
finds there and apply it to the project's destination clusters. Code review on your manifests repo
stops being a control, because nothing forces manifests to come from that repo. The same wildcard
also makes migrations and audits harder: you can't answer "what can this team deploy?" from the
project definition, only by reading every Application.

This matters most on the `default` project, which ships with `'*'` for both sources and
destinations. It's convenient for a first experiment and much too broad for anything a team shares.

## How to fix it

List the repositories or repository prefixes each project may use. Wildcards inside a path are
fine, and often what you want, so a team can add repositories under their own organization without
an operator editing the project.

**Flagged by this rule**

```yaml
spec:
  description: Example Project
  sourceRepos:
    - '*'
```

**Passes**

```yaml
spec:
  description: Example Project
  sourceRepos:
    - "https://github.com/acme/revenue/*"
```

Only the exact entry `*` is flagged, so a scoped pattern such as
`https://github.com/acme/revenue/*` passes. AppProjects with no `sourceRepos` list are inactive.

While you're editing the project, narrow `destinations` the same way, and give the project `roles`
so a team can sync its own Applications without seeing anyone else's.

## Exceptions

An operators-only project used to bootstrap the installation, locked down by RBAC so only cluster
admins can create Applications in it, is a reasonable place for a wildcard. Every project that
developers can reach should name its repositories.

## Further reading

For more information, see the Octopus Deploy blog post
[Securing Argo CD in a multi-tenant environment with application projects](https://octopus.com/blog/multi-tenant-argocd-with-application-projects),
and the Argo CD documentation on
[the default project](https://argo-cd.readthedocs.io/en/latest/user-guide/projects/#the-default-project).

[Back to all rules](index.md)
