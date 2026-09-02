# R10 – Sync timeout of zero

**Severity:** Medium

ArgoCD controller sync timeout is set to 0 seconds, which means no timeout and can cause inefficient synchronization.

## Why this matters

`controller.sync.timeout.seconds` bounds how long a sync operation may run before the controller
gives up on it. The Argo CD documentation describes it plainly: "Specifies a sync timeout for
applications. `0` means no timeout (default `0`)."

Worth being clear about what this rule is flagging: `0` is Argo CD's own default, and the
documentation doesn't call it a mistake. The finding says you've written that default into
`argocd-cmd-params-cm` explicitly, which is a good moment to decide whether you meant to.

The case for a bounded timeout is that a sync with no timeout has no way to end on its own. A PreSync
hook Job that never completes, a resource waiting on a health check that will never pass, or a
cluster that stopped responding mid-sync all leave the operation running indefinitely. While it
runs, the Application holds its operation slot: you can't start a new sync, self-heal won't correct
drift, and the Application sits in `Progressing` rather than surfacing as failed. Nothing in your
alerting fires, because nothing has technically gone wrong yet.

With a timeout, the same situation ends in a failed sync you can see, retry, and alert on.

## How to fix it

Set a timeout comfortably longer than your slowest legitimate sync. Look at how long your biggest
Application takes, including hooks, and leave generous headroom.

**Flagged by this rule**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cmd-params-cm
data:
  controller.sync.timeout.seconds: "0"
```

**Passes**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cmd-params-cm
data:
  controller.sync.timeout.seconds: "600"
```

Omitting the key entirely also passes — the rule only reports a value it can parse as the integer
zero. It's inactive when `argocd-cmd-params-cm` isn't among the scanned files.

## Exceptions

Installations with genuinely long-running syncs — large Helm charts with slow hooks, or database
migrations run as PreSync Jobs — may find any timeout they pick either too short for the worst case
or too long to be useful. Leaving it unbounded is a defensible choice there, as long as something
else watches for Applications stuck in `Progressing`.

## Further reading

For more information, see the
[`argocd-cmd-params-cm` reference](https://argo-cd.readthedocs.io/en/latest/operator-manual/argocd-cmd-params-cm-yaml/)
in the Argo CD documentation, which documents the `controller.sync.timeout.seconds` key.

[Back to all rules](index.md)
