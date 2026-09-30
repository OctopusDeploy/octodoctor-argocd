# R03 – Low requeueAfterSeconds

!!! warning "Severity: Medium"

ApplicationSets have requeueAfterSeconds set to less than 180 seconds (3 minutes). This can cause excessive API calls and rate limiting issues,
especially for GitHub.

## Why this matters

The Pull Request and SCM Provider generators find their inputs by asking your Git provider, over
its API, on a timer. `requeueAfterSeconds` is that timer. Set it to 20 and every ApplicationSet
using the generator hits GitHub, GitLab, or Bitbucket 3 times a minute, forever, whether or not
anything changed.

Git providers meter those calls. GitHub's REST API gives an authenticated token a fixed budget per
hour across your whole organization, and Argo CD is rarely the only thing spending it. When the
budget runs out, the generator's calls start failing, the ApplicationSet stops producing children,
and preview environments quietly fail to appear or fail to be cleaned up. The failure looks like an
Argo CD bug and is really a quota problem, which makes it slow to diagnose.

Three minutes is the interval used in Octopus Deploy's own preview environments walkthrough, and
it's a sensible floor. Pull requests aren't opened often enough to justify anything tighter.

## How to fix it

Raise `requeueAfterSeconds` to 180 or more. This rule looks for the property anywhere in the
ApplicationSet, so it catches it inside a `matrix` or `merge` generator too, and it flags the file
if any occurrence is below the threshold.

**Flagged by this rule**

```yaml
spec:
  generators:
    - pullRequest:
        github:
          owner: myorg
          repo: myrepository
        requeueAfterSeconds: 20
```

**Passes**

```yaml
spec:
  generators:
    - pullRequest:
        github:
          owner: myorg
          repo: myrepository
        requeueAfterSeconds: 180
```

ApplicationSets that never set `requeueAfterSeconds` are inactive for this rule. They fall back to
the controller-wide default, which is already well above 3 minutes.

## Exceptions

If you need previews to appear faster than the poll interval, don't shorten the interval. Configure
a webhook from your Git provider to the Argo CD ApplicationSet controller so a new pull request
triggers generation immediately, and leave the timer as a slow backstop. See also
[R22](r22-no-webhooks-configured.md).

## Further reading

For more information, see the Octopus Deploy blog post
[Create temporary Argo CD preview environments based on pull requests](https://octopus.com/blog/creating-temporary-preview-environments-based-pull-requests-argo-cd).

[Back to all rules](index.md)
