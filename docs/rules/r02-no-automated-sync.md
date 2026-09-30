# R02 – No automated sync

!!! warning "Severity: Medium"

Applications do not have automated sync policy configured. Configuration drift is not eliminated.

## Why this matters

Correcting configuration drift is the single biggest thing Argo CD does for you. When auto-sync and
self-heal are on, the cluster converges on what Git says, and it keeps converging: someone who
edits a live Deployment with `kubectl` sees their change reverted, and Git stays the source of
truth.

By turning autosync off (and losing self-heal) you don't have the basic GitOps guarantee that what
is in Git is also in your cluster.

Teams usually turn it off to gate deployments, so nothing reaches production without a human
pressing a button. That gate belongs somewhere else. Put it in Git, with branch protection and
required reviews on the repository holding your manifests, and put it in the cluster, with RBAC that
controls who can change what. Both survive a restart of Argo CD and both leave an audit trail. A
sync button doesn't.

Manual sync also means the cluster can sit out of sync for hours without anyone noticing, so
"deployed" and "merged" drift apart, and incident response starts with the question of which one you
are actually looking at.

## How to fix it

Add an `automated` block to the Application's sync policy. `selfHeal: true` is what reverts
hand-edits made directly against the cluster.

**Flagged by this rule**

```yaml
spec:
  syncPolicy:
    syncOptions:
      - CreateNamespace=true
```

**Passes**

```yaml
spec:
  syncPolicy:
    automated:
      selfHeal: true
      prune: false
```

An Application with no `syncPolicy` at all is also flagged. This rule is never inactive: every
Application either has an automated policy or it doesn't.

## Exceptions

A deliberately manual Application is defensible in a couple of places, like a break-glass runbook
you sync by hand during an incident, or an early proof of concept. If you're turning off auto-sync
across the whole installation, that's the finding this rule exists to raise.

## Further reading

For more information, see the Octopus Deploy blog post
[Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops).

[Back to all rules](index.md)
