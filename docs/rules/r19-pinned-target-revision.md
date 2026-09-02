# R19 – Pinned targetRevision

**Severity:** Medium

Applications have targetRevision set to a specific value (not `HEAD`) for Git repositories. This can lead to stale deployments.

## Why this matters

`targetRevision: HEAD` means "deploy whatever is on the default branch of this repository". It's a
dull setting that quietly does three valuable things.

Developers get self-service. Merging a change to the manifests repository deploys it, everywhere
that repository is referenced, without anyone editing an Argo CD manifest first. Auditing gets
easy, because deployment history and Git history become the same history: `git log` answers what
changed, when, and who approved it. And rollback is `git revert`, a command every engineer already
knows, rather than a coordinated edit of Application YAML.

Every other value gives one of those up:

- **Branch names** (`dev`, `staging`, `prod`) look like promotion but break auditing. Reconstructing
  what production ran last Tuesday means correlating two repositories' histories by timestamp.
- **Semantic version ranges** (`3.4.x`) remove rollback entirely. Argo CD only moves forward to newer
  matching versions, so during an incident you can't go back to the release that worked.
- **Git hashes** turn every deployment into a manual edit of the Application manifest by someone
  with access to it. Developer self-service is gone and incident response gets slower.
- **Git tags** are the least bad of the four and are a reasonable gate for production, but they still
  need a human to move them, and tags are mutable unless you've configured your Git provider to
  protect them.

## How to fix it

Use `HEAD` and promote changes by promoting *content* — Helm value files or Kustomize overlays —
between directories or repositories, rather than by moving a pointer.

**Flagged by this rule**

```yaml
spec:
  source:
    repoURL: https://github.com/example-org/example-repo.git
    targetRevision: dev
    path: my-staging-app
```

**Passes**

```yaml
spec:
  source:
    repoURL: https://github.com/example-org/example-repo.git
    targetRevision: HEAD
    path: my-staging-app
```

Omitting `targetRevision` also passes. The rule checks every entry of `spec.sources` as well as a
single `spec.source`.

## Exceptions

The rule is inactive for any source with a `chart` property, because a source that pulls an
off-the-shelf chart from a Helm repository or OCI registry has to name a chart version:

```yaml
spec:
  source:
    chart: sealed-secrets
    repoURL: https://bitnami-labs.github.io/sealed-secrets
    targetRevision: 1.16.1
```

That's the legitimate case — third-party infrastructure software you upgrade deliberately and never
promote between environments. Your own applications should track `HEAD`.

## Further reading

For more information, see these Octopus Deploy blog posts:

- [Abusing the target revision field for Argo CD promotions](https://octopus.com/blog/argocd-application-target-revision-field)
- [Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops)

The Argo CD documentation compares the options under
[tracking and deployment strategies](https://argo-cd.readthedocs.io/en/latest/user-guide/tracking_strategies/#head-branch-tracking).

[Back to all rules](index.md)
