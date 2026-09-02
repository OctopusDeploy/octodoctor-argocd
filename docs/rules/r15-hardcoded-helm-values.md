# R15 – Hardcoded Helm values

!!! warning "Severity: Medium"

Applications use hardcoded values in helm configuration. Consider using valueFiles instead of hardcoded values for better maintainability.

## Why this matters

Helm already has a way to layer configuration: an ordered list of value files, where later files
override earlier ones. It's well understood, it works the same on a laptop as it does in a cluster,
and the files live in Git next to the chart they configure.

Inlining the same settings into the Application manifest, through `helm.values` or
`helm.valuesObject`, throws that away. A developer who wants to know what their service is
configured with now has to find the Application manifest, which usually lives in a different
repository owned by a different team. They can't reproduce the deployment locally, because
`helm template` against the chart no longer renders what the cluster runs. And every settings
change, however small, becomes a pull request against operator-owned infrastructure YAML.

It also erases the value hierarchy. With files, you can see that a production cluster gets
`common-values.yaml`, then `all-prod-envs.yaml`, then a cluster-specific file, and reason about
which layer set what. With one inline blob, there are no layers left to reason about.

## How to fix it

Move the values into files in Git and reference them with `valueFiles`. Order matters: the last
file wins.

**Flagged by this rule**

```yaml
spec:
  source:
    repoURL: https://github.com/example-org/example-repo.git
    targetRevision: HEAD
    path: my-chart
    helm:
      valuesObject:
        image:
          repository: docker.io/example/my-app
          tag: 0.1
```

**Passes**

```yaml
spec:
  source:
    repoURL: https://github.com/example-org/example-repo.git
    targetRevision: HEAD
    path: my-chart
    helm:
      valueFiles:
        - common-values.yaml
        - values-production.yaml
```

Both `helm.values` (a YAML string) and `helm.valuesObject` (a nested mapping) are flagged. The rule
checks every entry of `spec.sources` as well as a single `spec.source`, so a multi-source
Application is flagged if any of its sources inlines values. Applications with no `helm` block
anywhere are inactive.

If your chart comes from a Helm repository or an OCI registry rather than your own Git repo, the
multi-source feature is how you combine it with your own value files: one source for the chart, one
source for the Git repo holding the values, with `valueFiles` using the `$ref` syntax to point at
it. See [R25](r25-more-than-two-sources.md) for the limits of that pattern.

## Exceptions

A single value that must be computed by the ApplicationSet template, such as an image tag derived
from a pull request SHA in a preview environment, has nowhere else to live. Keep those to a minimum
and put everything static in files.

## Further reading

For more information, see these Octopus Deploy blog posts:

- [Using Helm hierarchies in multi-source Argo CD applications for promoting to different GitOps environments](https://octopus.com/blog/helm-values-argocd)
- [Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops)

[Back to all rules](index.md)
