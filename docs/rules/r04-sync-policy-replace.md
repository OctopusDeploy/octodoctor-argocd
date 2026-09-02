# R04 – Sync policy replace

**Severity:** Medium

Applications have sync policy replace=true. This can lead to unintended resource deletion.

## Why this matters

By default Argo CD syncs with `kubectl apply`, which merges your manifest into whatever is already
in the cluster. The `Replace=true` sync option swaps that for `kubectl replace` or `kubectl create`,
which overwrites the object wholesale.

The Argo CD documentation is blunt about the consequence: this sync option "has the potential to be
destructive and might lead to resources having to be recreated, which could cause an outage for
your application."

Recreation is the part that bites. Fields you never set in your manifest — ones a controller, a
mutating webhook, or a default filled in — are not merged, they're dropped. A resource that can't be
updated in place gets deleted and made again, which for a Service means a new cluster IP, for a
PersistentVolumeClaim means the bound volume is at risk, and for a Deployment means every Pod
restarts at once rather than rolling. None of that is visible in the diff you review before syncing;
it happens at apply time.

## How to fix it

Remove `Replace=true` and let Argo CD use `kubectl apply`. If you turned it on to work around a
specific resource, scope it to that resource with the
`argocd.argoproj.io/sync-options: Replace=true` annotation instead of setting it for the whole
Application.

**Flagged by this rule**

```yaml
spec:
  syncPolicy:
    syncOptions:
      - Replace=true
```

**Passes**

```yaml
spec:
  syncPolicy:
    syncOptions:
      - Replace=false
      - CreateNamespace=true
```

The rule reads `spec.syncPolicy.syncOptions` and looks for an entry starting with `Replace=`. It's
inactive when there's no `syncOptions` list, or when the list has no `Replace=` entry at all — the
default behavior is what you want, so not mentioning the option is fine.

## Exceptions

There's one legitimate reason the documentation gives: a resource whose spec is too large to fit in
the `kubectl.kubernetes.io/last-applied-configuration` annotation that `kubectl apply` maintains.
Large CustomResourceDefinitions are the usual culprit. Set the option on that one resource with the
annotation rather than on every resource the Application owns.

## Further reading

For more information, see the Argo CD documentation on
[replacing resources instead of applying changes](https://argo-cd.readthedocs.io/en/latest/user-guide/sync-options/#replace-resource-instead-of-applying-changes).

[Back to all rules](index.md)
