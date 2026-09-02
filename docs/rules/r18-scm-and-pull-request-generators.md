# R18 – SCM and Pull Request generators combined

!!! warning "Severity: Medium"

ApplicationSets combine SCM Provider and Pull Request generators in matrix/merge configurations. This can cause exponential API calls and exhaust rate limits.

## Why this matters

The Argo CD documentation doesn't name this pairing as a problem. The argument below is assembled
from things it does state about the pieces.

Both of these generators work by calling your Git provider's API. The SCM Provider generator scans
an organization to discover repositories, and with `allBranches: true` it scans every branch of
every one of them. The Pull Request generator polls for open pull requests, by default every 30
minutes. The documentation warns for both that unauthenticated requests "have a lower rate limit and
can only see public repositories," and the ApplicationSet controller's rate-limiting guidance notes
that "ApplicationSet generators may call external APIs (SCM providers, Kubernetes secrets, cluster
API servers)" with deliberately conservative defaults to avoid hammering them.

A matrix combines exactly two child generators and produces their cross product. Put these two
inside one, and the API cost multiplies rather than adds: the SCM side enumerates repositories, and
the pull request side is evaluated across that enumeration. An organization with 50 repositories
turns one refresh into a large burst of calls, repeated on every requeue, forever.

When the budget runs out the failure is indirect. Generation starts erroring, the ApplicationSet
stops producing children, and preview environments quietly fail to appear or fail to be cleaned up.
It looks like an Argo CD fault and is actually a quota you exhausted, which makes it slow to
diagnose. Worse, the exhausted budget is shared, so unrelated ApplicationSets and CI jobs using the
same token start failing too.

## How to fix it

Use the Pull Request generator on its own, or combine it with a generator that doesn't call the
provider API — a `list` or `clusters` generator costs nothing extra per refresh.

**Flagged by this rule**

```yaml
spec:
  generators:
    - matrix:
        generators:
          - scmProvider:
              github:
                organization: myorg
                allBranches: true
          - pullRequest:
              github:
                owner: example-owner
                repo: preview-env-source-code
              requeueAfterSeconds: 180
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
          - pullRequest:
              github:
                owner: example-owner
                repo: preview-env-source-code
              requeueAfterSeconds: 180
```

The rule only reports a `matrix` or `merge` containing both an `scmProvider` and a `pullRequest`
child. An ApplicationSet with no Pull Request generator inside a combination is inactive, which
includes a standalone Pull Request generator at the top level.

If you need previews across many repositories, run one ApplicationSet per repository rather than
discovering them, and always supply credentials through `tokenRef` or `appSecretName` — the
authenticated rate limit is far higher than the anonymous one. Configuring pull request webhooks
removes most of the polling entirely. See also [R03](r03-low-requeue-after-seconds.md).

## Exceptions

A small organization with a handful of repositories may stay comfortably inside its rate limit with
this combination, particularly with an authenticated token and a long `requeueAfterSeconds`. Check
your provider's rate limit usage before deciding the finding doesn't apply — the failure mode only
shows up once you're near the ceiling.

## Further reading

For more information, see the Argo CD documentation on the
[Matrix generator restrictions](https://argo-cd.readthedocs.io/en/latest/operator-manual/applicationset/Generators-Matrix/#restrictions),
the [SCM Provider generator](https://argo-cd.readthedocs.io/en/latest/operator-manual/applicationset/Generators-SCM-Provider/),
the [Pull Request generator](https://argo-cd.readthedocs.io/en/latest/operator-manual/applicationset/Generators-Pull-Request/),
and [rate limiting ApplicationSet reconciliations](https://argo-cd.readthedocs.io/en/latest/operator-manual/high_availability/#rate-limiting-applicationset-reconciliations).

[Back to all rules](index.md)
