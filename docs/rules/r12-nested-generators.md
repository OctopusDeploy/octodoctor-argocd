# R12 – Deeply nested generators

**Severity:** Medium

ApplicationSets have 3 or more levels of generators by nesting matrix/merge combinations. This creates overly complex and hard-to-maintain configurations.

## Why this matters

Argo CD lets you put a `matrix` inside a `merge`, or a `merge` inside a `matrix`, as deep as you
like. That doesn't mean you should. There are only about three levels of abstraction in a real
deployment topology — the application, the environment, and the cluster — and teams that build four
or five levels of generators end up with a mental model nobody can hold.

The cost shows up when something goes wrong. A deeply nested ApplicationSet generates a set of
Applications you can't predict by reading it, so you find out what it does by applying it. A change
to the innermost generator can multiply out across every level above it and create or destroy
dozens of Applications at once. Parameter names collide between levels, and working out which level
supplied a given value means tracing the whole chain by hand.

The fix is almost always to move a dimension out of the generators and into Git, where you can see
it. Directory layout, cluster labels, and Helm value files each express one of those dimensions more
legibly than another nesting level does.

## How to fix it

Keep to a single `matrix` or `merge` with plain generators as its children — that's the two-level
shape this rule considers healthy. If you need a third dimension, express it as directory structure
for the Git generator or as labels on your cluster secrets, or split the ApplicationSet in two.

**Flagged by this rule**

```yaml
spec:
  generators:
    - merge:
        mergeKeys:
          - path
        generators:
          - matrix:                     # a combination generator inside a combination generator
              generators:
                - git:
                    directories:
                      - path: aws-ebs-csi-driver
                      - path: cluster-autoscaler
                - clusters: {}
          - list:
              elements:
                - path: aws-ebs-csi-driver
                  namespace: system
```

**Passes**

```yaml
spec:
  generators:
    - matrix:
        generators:
          - list:
              elements:
                - myapp: blog-client
                  myenv: stage
          - clusters: {}
```

ApplicationSets with no `matrix` or `merge` at all are inactive. A single `matrix` or `merge` whose
children are all plain generators passes.

## Exceptions

Rarely, an infrastructure ApplicationSet that fans out over applications, clusters, and regions
really does need a third level, and no directory layout expresses it better. If you keep it, add a
comment naming each level so the next reader doesn't have to reverse-engineer the cross product.

## Further reading

For more information, see these Octopus Deploy blog posts:

- [Structuring Argo CD repositories](https://octopus.com/blog/how-to-structure-your-argo-cd-repositories-using-application-sets)
- [Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops)

[Back to all rules](index.md)
