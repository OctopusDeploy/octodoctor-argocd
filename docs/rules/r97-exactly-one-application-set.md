# R97 – Exactly one ApplicationSet

!!! warning "Severity: Medium"

You have exactly one application set. This usually results in a very complex YAML file. Consider using multiple application sets.

## Why this matters

You aren't required to put every application into a single ApplicationSet, and teams that try
usually regret it. One ApplicationSet for everything has to cover every combination of application,
environment, and cluster you run, which means deeply nested generators, long lists of exceptions,
and template blocks full of conditionals to handle the cases that don't fit the general shape.

That file becomes a chokepoint. Everyone who deploys anything edits it, so pull requests queue up
behind each other and conflict. A change intended for one team's applications is one typo away from
regenerating everyone else's. Reviewers can't tell from a diff which Applications will be created or
destroyed. And because Argo CD gives you no dry run for generation, you find out by applying it.

Splitting doesn't cost you anything. ApplicationSets are cheap, and several small ones are easier to
read, review, and blame than one large one.

## How to fix it

Split along whatever axis your organization actually uses: by team, by environment, by application
type, or by the infrastructure-versus-workloads line. Separating infrastructure components from
developer workloads is usually the first and most valuable cut, since it stops developers debugging
failures in software they don't own.

**Flagged by this rule**

```
gitops-repo/
└── applicationsets/
    └── everything.yaml          # one ApplicationSet generating every app in every environment
```

**Passes**

```
gitops-repo/
└── applicationsets/
    ├── infrastructure.yaml      # cluster-wide components, deployed first
    ├── team-revenue-qa.yaml
    ├── team-revenue-prod.yaml
    └── team-platform-prod.yaml
```

This rule looks at the installation as a whole rather than at individual files. It fails when the
scanned directory contains exactly one ApplicationSet, and is inactive when it contains none (see
[R98](r98-no-application-sets.md)) or more than one.

## Exceptions

A small installation with a handful of applications on one cluster may genuinely need only one
ApplicationSet, and splitting it would be ceremony. The finding is worth acting on once that file
starts growing generators to accommodate special cases.

## Further reading

For more information, see these Octopus Deploy blog posts:

- [Structuring Argo CD repositories](https://octopus.com/blog/how-to-structure-your-argo-cd-repositories-using-application-sets)
- [Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops)

[Back to all rules](index.md)
