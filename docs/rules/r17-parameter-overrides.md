# R17 – Parameter overrides

!!! danger "Severity: High"

Applications use parameter overrides in helm configuration. Consider using valueFiles instead of hardcoded parameters for better maintainability.

## Why this matters

`helm.parameters` is Argo CD's equivalent of `helm install --set`. It exists mainly so the UI and
the CLI can change a setting on a running Application without touching Git, and that's exactly the
problem: it's a documented way to bypass the source of truth.

After someone sets a parameter through the UI, the cluster's configuration and the repository's
configuration disagree, and nothing in Git records the difference. The usual GitOps questions stop
having answers. What changed? Who changed it? Can I roll it back? A `git log` won't tell you, and
the next person to redeploy the Application from scratch gets different behavior with no
explanation.

Parameters written into the manifest by hand are less dangerous but share the deeper flaw with
inline values: they aren't Helm's own mechanism, so a developer can't reproduce the deployment with
`helm template` locally. That's why this rule is High severity while
[R15](r15-hardcoded-helm-values.md) is Medium: parameters carry the extra risk of drift introduced
outside Git entirely.

## How to fix it

Express the same settings as value files in Git and reference them with `valueFiles`.

**Flagged by this rule**

```yaml
spec:
  source:
    repoURL: https://github.com/example-org/example-repo.git
    targetRevision: HEAD
    path: my-chart
    helm:
      parameters:
        - name: "my-example-setting-1"
          value: my-value1
        - name: "my-example-setting-2"
          value: "my-value2"
          forceString: true
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
        - values-production.yaml
```

The rule checks every entry of `spec.sources` as well as a single `spec.source`. Applications with
no `helm` block anywhere are inactive.

If your team uses the UI to override parameters, consider turning that off. Restricting who can
update Applications through AppProject roles (see
[R05](r05-permissive-source-repos.md)) removes the temptation as well as the capability.

## Exceptions

There isn't a good long-term reason to keep parameters in a manifest you control. The one place
they're hard to avoid is a value that only exists at generation time, such as an image tag templated
into a preview environment by an ApplicationSet.

## Further reading

For more information, see these Octopus Deploy blog posts:

- [Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops)
- [Using Helm hierarchies in multi-source Argo CD applications for promoting to different GitOps environments](https://octopus.com/blog/helm-values-argocd)

[Back to all rules](index.md)
