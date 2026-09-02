# R16 – Hardcoded Kustomize values

**Severity:** Medium

Applications use hardcoded values in kustomize configuration. Consider using kustomize overlays instead of hardcoded values for better maintainability.

## Why this matters

Kustomize is built around bases and overlays. An overlay is a directory containing a
`kustomization.yaml` that patches a base, and it's a complete, self-contained description of one
variant of your application. You can render it with `kustomize build`, review it in a pull request,
and diff two of them against each other.

Argo CD's `kustomize` block in the Application manifest duplicates a subset of that: `namePrefix`,
`images`, `namespace`, `patches`, and friends. Anything you set there is invisible to Kustomize
itself. `kustomize build` on the overlay no longer produces what Argo CD applies, so local testing
stops matching reality, and the description of one environment is now split across two files in two
repositories owned by two teams.

That split is the real cost. Changing which image tag production runs should be a change to the
production overlay, reviewed by the people who own production manifests. Putting it in the
Application manifest turns it into an infrastructure change instead.

## How to fix it

Move every customization into an overlay directory and point the Application's `path` at that
overlay. The Application manifest should then say only *which* overlay, not *what* it does.

**Flagged by this rule**

```yaml
spec:
  source:
    repoURL: https://github.com/example-org/example-repo.git
    targetRevision: HEAD
    path: my-app
    kustomize:
      namePrefix: prod-
      images:
        - docker.io/example/my-app:0.2
      namespace: custom-namespace
```

**Passes**

```yaml
spec:
  source:
    repoURL: https://github.com/example-org/example-repo.git
    targetRevision: HEAD
    path: my-app/overlays/prod
```

This rule flags the presence of a `kustomize` block, not its contents, so any `kustomize:` section
in an Application is reported. It checks every entry of `spec.sources` as well as a single
`spec.source`. The rule is never inactive: an Application either has a `kustomize` block or it
doesn't.

## Exceptions

As with Helm values, a value that only the ApplicationSet template can know — a preview
environment's image tag built from the pull request SHA, for example — has to be set at this level.
Everything that's stable across syncs belongs in the overlay.

## Further reading

For more information, see the Octopus Deploy blog post
[Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops).

[Back to all rules](index.md)
