# R09 – Duplicate generators

!!! warning "Severity: Medium"

ApplicationSets have duplicate generators in matrix/merge configurations. This creates complex and confusing configurations that are hard to maintain.

## Why this matters

`matrix` and `merge` exist to combine generators of *different* kinds. The canonical example is a
Git generator crossed with a cluster generator: the Git side discovers which applications exist, the
cluster side discovers where they should go, and the matrix produces every sensible combination.
Each child answers a different question, so the result is easy to describe in one sentence.

Two children of the same type break that. A matrix of 2 Git generators produces the cross product
of 2 directory listings, which is rarely what anyone wants and almost never what the next person
to read the file expects. Discovering 6 directories on one side and 4 on the other silently
generates 24 Applications. Nobody can tell from the YAML which combinations are meaningful, so the
ApplicationSet becomes something you change by trial and error.

There's usually a simpler shape hiding behind it. Two Git generators normally means "these are 2
different groups of applications", which is 2 ApplicationSets, not one matrix.

## How to fix it

Give each child of a `matrix` or `merge` a distinct generator type. If you genuinely have 2
independent sets of the same type, split them into separate ApplicationSets.

**Flagged by this rule**

```yaml
spec:
  generators:
    - matrix:
        generators:
          - git:
              repoURL: https://github.com/example-org/gitops-examples.git
              revision: HEAD
              directories:
                - path: application-sets/example-apps/*
          - git:
              repoURL: https://github.com/example-org/gitops-examples.git
              revision: HEAD
              directories:
                - path: application-sets/example-apps/values/*
```

**Passes**

```yaml
spec:
  generators:
    - matrix:
        generators:
          - git:
              repoURL: https://github.com/example-org/gitops-examples.git
              revision: HEAD
              directories:
                - path: application-sets/example-apps/*
          - clusters: {}
```

ApplicationSets with no `matrix` or `merge` generator are inactive for this rule, since a flat list
of generators has nothing to duplicate.

## Exceptions

You might have an exotic directory structure that really needs 2 generators of the same kind. It is better to
restructure your Git folders using best practices instead of trying to make Argo CD workarounds.

## Further reading

For more information, see these Octopus Deploy blog posts:

- [Structuring Argo CD repositories](https://octopus.com/blog/how-to-structure-your-argo-cd-repositories-using-application-sets)
- [Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops)

[Back to all rules](index.md)
