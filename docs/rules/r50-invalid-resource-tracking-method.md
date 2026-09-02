# R50 – Invalid resourceTrackingMethod

!!! danger "Severity: High"

ArgoCD resourceTrackingMethod is set to an invalid value. It must be exactly one of: label, annotation, or annotation+label.

## Why this matters

Resource tracking is how Argo CD knows which Kubernetes resources belong to which Application. It's
the mechanism behind pruning, out-of-sync detection, the resource tree in the UI, and adopting
resources that already exist in a cluster. Get it wrong and Argo CD either can't find resources it
owns, or claims resources it doesn't.

`application.resourceTrackingMethod` in the `argocd-cm` ConfigMap accepts exactly three values:

- `label` — the legacy default. Argo CD writes `app.kubernetes.io/instance`, which is capped at 63
  characters and is a label other tools also write, so names can collide or be truncated.
- `annotation` — writes `argocd.argoproj.io/tracking-id`. No length limit and no collisions with
  other tooling. This is the recommended setting for new installations.
- `annotation+label` — writes both, for compatibility with tools that read the instance label while
  Argo CD relies on the annotation.

Anything else is not a fourth mode; it's a typo. Argo CD falls back to its default rather than
failing loudly, so the installation keeps running with tracking behavior nobody chose. That surfaces
later, at the worst time — during a migration between Argo CD instances, or when adopting resources
an older tool created — as resources Argo CD refuses to take ownership of, or duplicates it tries to
create alongside them.

The most common typo is spaces around the plus sign.

## How to fix it

Set the value to one of the three accepted strings, with no surrounding whitespace.

**Flagged by this rule**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cm
data:
  application.resourceTrackingMethod: annotation + label
```

**Passes**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cm
data:
  application.resourceTrackingMethod: annotation+label
```

`label` and `annotation` pass too. If the `argocd-cm` ConfigMap isn't in the scanned directory, or
it doesn't set the key at all, the rule is inactive — Argo CD then uses its own default.

## Exceptions

None. The three values are the complete set that Argo CD understands, so a flagged value is always
a mistake.

## Further reading

For more information, see the Octopus Deploy blog post
[Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops),
and the Argo CD documentation on
[choosing a tracking method](https://argo-cd.readthedocs.io/en/latest/user-guide/resource_tracking/#choosing-a-tracking-method).

[Back to all rules](index.md)
