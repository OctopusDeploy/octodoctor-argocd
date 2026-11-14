# R08 – Self-heal disabled

!!! warning "Severity: Medium"

Applications have an automated sync policy without self-heal. Manual changes made directly against
the cluster aren't reverted.

## Why this matters

The big advantage of Argo CD over traditional pipelines is that it is bidirectional.
It applies changes from Git to the cluster but also makes sure that the cluster
always has what is in Git.

This is only true if `self-heal` is enabled. If not then anybody can make manual
changes to the cluster with `kubectl` or other tools/controllers and instantly
create configuration drift. This is especially important for production environments
where you want the hard guarantee that cluster state matches Git.

Without self-heal, the live state can differ from Git for days. A hotfix applied by hand during an
incident quietly survives, nobody knows it exists, and the next unrelated commit overwrites it at the
worst possible moment.

## How to fix it

Set `selfHeal: true` in the Application's `automated` block.

**Flagged by this rule**

```yaml
spec:
  syncPolicy:
    automated:
      prune: false
```

**Passes**

```yaml
spec:
  syncPolicy:
    automated:
      selfHeal: true
      prune: false
```

An `automated` block with no value, or with `selfHeal: false`, is also flagged. Applications with no
`automated` block at all are reported by [R02](r02-no-automated-sync.md) instead, and this rule is
inactive when no Application uses auto-sync.

## Exceptions

You may want to leave self-heal off briefly while debugging an Application live in the cluster. Turn
it back on when you're done. Fields that another controller legitimately changes, like `replicas`
managed by a HorizontalPodAutoscaler, belong in `ignoreDifferences`, not in a reason to turn
self-heal off.

## Further reading

For more information, see the Octopus Deploy blog post
[Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops).

[Back to all rules](index.md)
