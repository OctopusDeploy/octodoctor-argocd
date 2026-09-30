# R01 – Dangerous sync policy

!!! danger "Severity: High"

Applications have dangerous sync policy settings (`prune=true` and `allowEmpty=true`). This combination can lead to unintended resource deletion.

## Why this matters

`prune: true` tells Argo CD to delete any cluster resource that no longer appears in the rendered
manifests. `allowEmpty: true` tells it to proceed even when the rendering produces nothing at all.
Together, they turn a bad render into a cluster-wide delete. A broken Helm values file, a mistyped
`path`, a renamed directory, or a repository that briefly returns an empty tree is enough for Argo CD
to conclude the desired state is "no resources" and prune everything the Application owns. This problem
is even more evident when your Git repository is down.

With `allowEmpty: false` (the default), Argo CD refuses to sync an empty result. It reports the
Application as out of sync instead, which is a loud signal that something upstream is broken.

The blast radius can be catastrophic if you have transient errors in your Git repository access. Argo CD deletes child resources in a cascading
manner, so one accidental empty render can take out Deployments, Services, PersistentVolumeClaims,
and everything else the Application tracks.

## How to fix it

Leave `allowEmpty` out so it takes its default of `false`, or set it explicitly. You can keep
`prune: true` if you want Argo CD to clean up resources you deliberately removed from Git.

**Flagged by this rule**

```yaml
spec:
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
      allowEmpty: true
```

**Passes**

```yaml
spec:
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
      allowEmpty: false
```

The rule is inactive for Applications with no `spec.syncPolicy.automated` block, because manual sync
never prunes on its own.

## Exceptions

A handful of Applications legitimately render to zero resources in some environments, such as an
optional component that's switched off per cluster. Those need `allowEmpty: true` to sync at all.
Keep them small and narrowly scoped so an accidental prune can't reach anything important, and make
sure every resource they own is in Git so you can recreate it.

## Further reading

For more information, see these Octopus Deploy blog posts:

- [Argo CD finalizers: prevent accidental deletion issues](https://octopus.com/blog/argocd-application-deletion-finalizers)
- [Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops)

The Argo CD documentation describes the setting itself under
[automatic pruning with allow-empty](https://argo-cd.readthedocs.io/en/latest/user-guide/auto_sync/#automatic-pruning-with-allow-empty-v18).

[Back to all rules](index.md)
