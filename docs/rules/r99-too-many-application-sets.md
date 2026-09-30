# R99 – Too many ApplicationSets

!!! warning "Severity: Medium"

You have too many ApplicationSets. This can cause performance issues and make them hard to manage.

## Why this matters

ApplicationSets are cheap individually and not free in aggregate. Every one of them is reconciled on
its own schedule by the ApplicationSet controller, and generators that call out to a Git provider
(Git, SCM Provider, and Pull Request) turn that schedule into API traffic. Dozens of ApplicationSets
polling the same organization share one rate limit budget, and when it runs out, generation starts
failing across all of them at once.

The management cost grows too. Past a certain count you can no longer answer "where does this
application get deployed?" by reading one file, because the answer is spread across many that each
cover a slice. Duplication creeps in: the same template block, the same cluster selector, and the
same sync policy copied into 20 places, drifting apart over time. Changing a convention across
the installation stops being a single edit.

This is the counterpart to [R97](r97-exactly-one-application-set.md). One giant ApplicationSet is
unreadable; a large number of tiny ones is unmanageable. The healthy shape is a moderate number of
ApplicationSets split along an axis that means something (team, environment, or the
infrastructure-versus-workloads line) rather than one per application.

## How to fix it

Look for ApplicationSets that differ only in a value and merge them, using a generator to supply the
difference. Two ApplicationSets identical apart from a cluster name should be one, selecting both
clusters by label.

**Flagged by this rule**

```
gitops-repo/
└── applicationsets/
    ├── billing-qa.yaml
    ├── billing-prod.yaml          # one ApplicationSet per app per environment
    ├── checkout-qa.yaml
    ├── checkout-prod.yaml
    └── ...                        # more than 10 in total
```

**Passes**

```
gitops-repo/
└── applicationsets/
    ├── infrastructure.yaml
    ├── team-revenue-qa.yaml       # a Git generator picks up every app in the team's directory
    ├── team-revenue-prod.yaml
    └── team-platform-prod.yaml
```

This rule looks at the installation as a whole and fails when the scanned directory contains more
than 10 ApplicationSets.

## Exceptions

A large organization with many teams may pass 10 legitimately, and that's fine as long as each
ApplicationSet covers a group rather than a single application. Treat the finding as a prompt to
check that the split follows a deliberate axis, and pay attention to how many of them poll a Git
provider (see [R03](r03-low-requeue-after-seconds.md)).

## Further reading

For more information, see these Octopus Deploy blog posts:

- [Structuring Argo CD repositories](https://octopus.com/blog/how-to-structure-your-argo-cd-repositories-using-application-sets)
- [Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops)

[Back to all rules](index.md)
