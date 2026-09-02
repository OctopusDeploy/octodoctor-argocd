# R14 – Prune enabled

!!! warning "Severity: Medium"

Applications have prune=true in sync policy. This can lead to unintended resource deletion.

## Why this matters

`prune: true` completes the GitOps loop. Delete a manifest from Git and Argo CD deletes the
matching resource from the cluster, so Git really is the full description of what's running. Most
mature installations want this on.

It's on this list because the same setting makes a class of ordinary mistakes destructive. Renaming
an Application, moving a directory, changing the `path` a source points at, or splitting one
Application into two all look to Argo CD like "these resources are no longer declared", and it acts
on that. Refactoring your repository layout on a Friday afternoon with prune enabled everywhere is
how people lose PersistentVolumeClaims.

This rule is a lower-severity companion to [R01](r01-dangerous-sync-policy.md). R01 catches the
combination that can delete *everything* an Application owns in one sync. R14 flags prune on its own
so you can confirm it's a deliberate choice for each Application rather than something copied from
an example manifest.

## How to fix it

Decide per Application. Where you want the loop closed, leave `prune: true` and make it safe:
keep every resource in Git so you can recreate it, avoid renaming things casually, and watch
deletion events so an accidental removal surfaces quickly. Where an Application owns stateful
resources you'd rather delete by hand, set `prune: false` and clean up deliberately.

**Flagged by this rule**

```yaml
spec:
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

**Passes**

```yaml
spec:
  syncPolicy:
    automated:
      prune: false
      selfHeal: true
```

Omitting `prune` entirely also passes, because the default is `false`. Applications with no
`spec.syncPolicy.automated` block are inactive.

## Exceptions

Pruning is a legitimate and common choice, so treat this finding as a prompt rather than a defect.
Ignore it once you've checked that the Application's resources are all declared in Git, that the
`argocd.argoproj.io/sync-options: Prune=false` annotation protects any individual resource you
can't afford to lose, and that resource tracking is set up so Argo CD knows exactly which resources it owns
(see [R50](r50-invalid-resource-tracking-method.md)).

## Further reading

For more information, see these Octopus Deploy blog posts:

- [Argo CD finalizers: prevent accidental deletion issues](https://octopus.com/blog/argocd-application-deletion-finalizers)
- [Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops)

The Argo CD documentation describes the setting itself under
[automatic pruning](https://argo-cd.readthedocs.io/en/latest/user-guide/auto_sync/#automatic-pruning).

[Back to all rules](index.md)
