# R07 – Reconciliation timeout without cache expiration

!!! warning "Severity: Medium"

Using `ARGOCD_RECONCILIATION_TIMEOUT=0` without proper `ARGOCD_DEFAULT_CACHE_EXPIRATION` settings can cause performance issues.

## Why this matters

`timeout.reconciliation` in `argocd-cm` controls how often the application controller goes back to
Git to see whether new manifests have been published. It defaults to 2 minutes plus jitter, and
as the Argo CD documentation puts it, "reconciliation by timeout is disabled if timeout is set to
0."

Setting it to `0` is a deliberate and reasonable thing to do when you have webhooks: polling every
application every 2 minutes is wasted work when your Git provider can tell you the instant
something changes. The problem is what happens when a webhook doesn't arrive. A provider outage, a
misconfigured hook, a network blip, or a firewall change is enough. Now nothing at all triggers a
refresh, because you switched off the only mechanism that didn't depend on the webhook.

How long that lasts is set by the repo-server's manifest cache. `reposerver.default.cache.expiration`
defaults to `24h0m0s`. So the combination of `timeout.reconciliation: "0"` and a default cache
expiry means a single missed webhook can leave an Application serving cached manifests for up to a
day, showing as synced the whole time. That's why this rule looks at both settings together and
only fails when polling is off *and* the cache is left at its default.

## How to fix it

Either keep a non-zero reconciliation timeout, or shorten the repo-server cache so it becomes the
backstop that polling used to be.

**Flagged by this rule**

```yaml
# argocd-cm
data:
  timeout.reconciliation: "0"
---
# argocd-cmd-params-cm (cache left at the default)
data:
  reposerver.default.cache.expiration: "24h0m0s"
```

**Passes**

```yaml
# argocd-cm
data:
  timeout.reconciliation: "0"
---
# argocd-cmd-params-cm (cache expiry shortened to bound the staleness)
data:
  reposerver.default.cache.expiration: "1h0m0s"
```

Setting `timeout.reconciliation` to anything other than `0` also passes, and the rule is inactive
when `argocd-cm` doesn't set the key at all. Note it examines both ConfigMaps together, so point the
tool at a directory containing `argocd-cm` and `argocd-cmd-params-cm` for it to reach a verdict.

## Exceptions

If you monitor your webhooks (you alert on a repository that hasn't delivered an event in longer
than expected, or on Applications whose `reconciledAt` has gone stale), you'll notice a broken
hook long before the cache matters, and the default expiry is fine. The finding is aimed at
installations that turned off polling and have nothing watching the replacement.

## Further reading

For more information, see the Argo CD documentation for
[`argocd-cm`](https://argo-cd.readthedocs.io/en/latest/operator-manual/argocd-cm-yaml/) (the
`timeout.reconciliation` key) and
[`argocd-cmd-params-cm`](https://argo-cd.readthedocs.io/en/latest/operator-manual/argocd-cmd-params-cm-yaml/)
(the `reposerver.default.cache.expiration` key).

[Back to all rules](index.md)
