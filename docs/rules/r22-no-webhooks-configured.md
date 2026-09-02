# R22 – No webhooks configured

!!! warning "Severity: Medium"

No webhooks are configured in argocd-secret. Webhooks improve Argo CD performance by triggering immediate syncs instead of polling.

## Why this matters

Left to itself, Argo CD "polls Git/OCI/Helm repositories every three minutes to detect changes to
the manifests." That's the only way it learns a commit happened, and it has two costs.

The obvious one is latency. A developer merges a change and then waits — up to three minutes before
Argo CD even notices, plus however long the sync takes. It's not a long wait in absolute terms, but
it's dead time in every deployment, and it makes the feedback loop feel disconnected from the merge
that caused it.

The less obvious one is load. That polling happens per repository, continuously, whether or not
anything changed. Multiply it across every repository in a large installation and it's a steady
stream of Git provider API calls and repo-server work that produces nothing most of the time. It
also competes for rate limit budget with the ApplicationSet generators that need those same APIs.

A webhook inverts it: your Git provider tells Argo CD the moment a push lands, the affected
Applications refresh immediately, and the polling interval becomes a backstop rather than the
primary mechanism.

## How to fix it

Create a webhook in your Git provider pointing at `/api/webhook` on your Argo CD server, then store
the shared secret in the `argocd-secret` Secret so Argo CD can verify incoming events.

**Flagged by this rule**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: argocd-secret
type: Opaque
data:
  tls.crt:
  tls.key:
  admin.password:
  server.secretkey:
```

**Passes**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: argocd-secret
type: Opaque
stringData:
  webhook.github.secret: shhhh! it's a github secret
  webhook.gitlab.secret: shhhh! it's a gitlab secret
```

The rule passes if any key in the Secret's `data` starts with `webhook`, which covers every provider
Argo CD supports: `webhook.github.secret`, `webhook.gitlab.secret`, `webhook.bitbucket.uuid`,
`webhook.bitbucketserver.secret`, `webhook.gogs.secret`, `webhook.azuredevops.username` and
`.password`, and `webhook.harbor.secret`. It's inactive when `argocd-secret` isn't among the scanned
files.

## Exceptions

The secret is optional — Argo CD accepts unauthenticated webhook calls and will still refresh on
them, so you can have working webhooks and still be flagged here. Don't do that on an instance
reachable from the internet, though: without a shared secret anyone who finds the endpoint can force
refreshes.

The other case is a self-hosted Git server that can't reach your Argo CD instance, where polling is
the only option available.

## Further reading

For more information, see the Argo CD documentation on
[Git webhook configuration](https://argo-cd.readthedocs.io/en/latest/operator-manual/webhook/).

[Back to all rules](index.md)
