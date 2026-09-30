# R20 – No cluster generator

!!! warning "Severity: Medium"

ApplicationSets do not use cluster generators despite having multiple connected clusters. This can lead to performance issues and configuration complexity.

## Why this matters

When you have more than one cluster registered with Argo CD, you need a way to say where each
application goes. There are 2 ways to say it, and they age very differently.

The ad-hoc way is to name clusters directly, usually with a `list` generator holding one element per
cluster, or with per-cluster on/off flags scattered through the ApplicationSets. It
works on day one and gets worse every week. Adding a cluster means finding and editing every
ApplicationSet that should target it, and copying the same list into each one. Removing a cluster
means the same sweep in reverse, done carefully, in files that each look slightly different.
Answering "which clusters run this component?" means reading all of them. Each cluster becomes a
snowflake defined by a scattering of literal names.

The label way treats clusters as groups that share properties. You label the cluster secret once
(`env: prod`, `region: eu-west-1`, `tier: internal`) and each ApplicationSet selects the group it
cares about. Adding a production cluster then means labeling it, and every ApplicationSet matching
`env: prod` picks it up with no edits at all. It stays the same amount of work from 2 clusters
to 200.

## How to fix it

Add labels to your cluster secrets and use the `clusters` generator with a `selector`, on its own or
as one side of a `matrix`.

**Flagged by this rule**

```yaml
spec:
  generators:
    - list:
        elements:
          - cluster: engineering-dev
            url: https://1.2.3.4
          - cluster: engineering-prod
            url: https://2.4.6.8
```

**Passes**

```yaml
spec:
  generators:
    - clusters:
        selector:
          matchLabels:
            env: prod
```

The rule finds a `clusters` generator at any depth, including inside `matrix` and `merge`
generators, so the common Git-crossed-with-clusters pattern passes.

It counts clusters by looking for Secrets labeled `argocd.argoproj.io/secret-type: cluster` in the
scanned directory. With no such Secrets the rule is inactive, and with exactly one it passes
automatically: a single cluster doesn't need a generator to distribute over.

## Exceptions

An ApplicationSet that deliberately targets one named cluster, such as a bootstrap ApplicationSet
that only ever runs on the management cluster, doesn't need a cluster generator. If most of your
ApplicationSets are flagged, though, the finding is pointing at a topology that will get expensive
to maintain.

## Further reading

For more information, see these Octopus Deploy blog posts:

- [Distribute your Argo CD applications to different Kubernetes clusters using application sets](https://octopus.com/blog/argocd-clusters-labels-with-apps)
- [Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops)

[Back to all rules](index.md)
