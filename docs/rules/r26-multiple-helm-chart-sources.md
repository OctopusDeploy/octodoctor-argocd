# R26 – Multiple Helm chart sources

**Severity:** Medium

Applications have 2 or more sources with Helm charts. Multiple Helm charts should be grouped differently.

## Why this matters

A source with a `chart` property pulls a packaged Helm chart from a Helm repository or an OCI
registry. Two of them in one Application means you're installing two separate pieces of software
under a single Argo CD Application, which is the multi-source feature being used as a folder rather
than as intended.

Multi-source exists so one chart can be combined with the value files that configure it. Two charts
have nothing to combine. What you get instead is one sync status covering both, so an unhealthy
Prometheus makes CoreDNS look broken too; one rollback that reverts both; one project, one
destination, and one sync policy for two things that will need different ones eventually; and
version upgrades that have to be coordinated because the Application syncs as a unit.

Each chart is its own deployable thing with its own release cadence, so each should be its own
Application.

## How to fix it

Give every chart its own Application, and generate them from an ApplicationSet if there are more
than a handful. Keep multi-source for the chart-plus-values pattern.

**Flagged by this rule**

```yaml
spec:
  sources:
    - repoURL: 'https://prometheus-community.github.io/helm-charts'
      chart: prometheus
      targetRevision: 15.7.1
    - repoURL: 'https://git.example.com/org/my-charts.git'
      chart: core-dns
      targetRevision: 1.9.3
```

**Passes**

```yaml
spec:
  sources:
    - repoURL: 'https://prometheus-community.github.io/helm-charts'
      chart: prometheus
      targetRevision: 15.7.1
      helm:
        valueFiles:
          - $values/charts/prometheus/values-production.yaml
    - repoURL: 'https://github.com/mycompany/gitops-values.git'
      targetRevision: HEAD
      ref: values
```

The rule counts sources that have a `chart` property, and fails at two or more. A multi-source
Application with one chart and one Git source passes, as does one with no charts at all.
Applications using the single `spec.source` field are inactive.

## Exceptions

There isn't a strong case for two charts in one Application. If you're doing it to keep two
components in lockstep, express that with sync waves across two Applications, or with an
app-of-apps parent, both of which keep the components independently observable.

## Further reading

For more information, see these Octopus Deploy blog posts:

- [Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops)
- [Using Helm hierarchies in multi-source Argo CD applications for promoting to different GitOps environments](https://octopus.com/blog/helm-values-argocd)

[Back to all rules](index.md)
