# R06 – Missing manifest-generate-paths

!!! warning "Severity: Medium"

Applications do not have the argocd.argoproj.io/manifest-generate-paths annotation configured.

## Why this matters

Every time Argo CD sees a new commit in a repository, it re-renders the manifests for every
Application pointing at that repository, then compares the result to the cluster. In a repository
holding one application that's cheap. In a monorepo holding two hundred, a one-line change to a
single service triggers two hundred renders, and each one means the repo-server fetching the
revision, running Helm or Kustomize, and diffing the output.

The `argocd.argoproj.io/manifest-generate-paths` annotation tells Argo CD which paths in the
repository actually affect this Application's output. When a new commit arrives, Argo CD compares
the changed files against that list, and if nothing matches it reuses the cached manifests instead
of regenerating them. The Application is marked as synced to the new revision without any work.

The saving compounds. Fewer renders means a smaller refresh queue, less repo-server CPU, and
noticeably faster feedback for the person whose change *did* match. This is the main lever Argo CD
gives you for monorepo scaling, and it costs one annotation per Application.

## How to fix it

Add the annotation and point it at the paths the Application renders from. The value is a
semicolon-separated list, and paths can be relative to the Application's `source.path`, absolute
from the repository root, or globs.

**Flagged by this rule**

```yaml
metadata:
  name: demo
  namespace: argocd
spec:
  source:
    path: ./my-app
```

**Passes**

```yaml
metadata:
  name: demo
  namespace: argocd
  annotations:
    argocd.argoproj.io/manifest-generate-paths: .
spec:
  source:
    path: ./my-app
```

`.` is the common case: only changes under the Application's own `source.path` matter. Use
`.;../shared` when the Application also renders a shared directory, `/guestbook` for a path from the
repository root, or a glob like `/shared/*-secret.yaml`.

This rule is never inactive — every Application either carries the annotation or doesn't.

## Exceptions

The annotation buys you little on an Application whose repository contains only that Application,
since every commit is relevant anyway. It's also worth skipping if you can't state the paths
accurately: getting the list wrong is worse than omitting it, because Argo CD will skip
regeneration for a change that genuinely mattered and the Application will sit on stale manifests.

## Further reading

For more information, see the Argo CD documentation on the
[manifest paths annotation](https://argo-cd.readthedocs.io/en/latest/operator-manual/high_availability/#manifest-paths-annotation)
and its entry in the
[annotations reference](https://argo-cd.readthedocs.io/en/latest/user-guide/annotations-and-labels/#manifest-generate-paths).

[Back to all rules](index.md)
