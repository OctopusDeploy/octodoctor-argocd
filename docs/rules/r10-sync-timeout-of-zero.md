# R10 – Sync timeout of zero

!!! warning "Severity: Medium"

The Argo CD controller sync timeout is set to 0 seconds, which means there's no timeout. This can cause inefficient synchronization.

## Why this matters

`controller.sync.timeout.seconds` bounds how long a sync operation may run before the controller
gives up on it. The Argo CD documentation describes it plainly: "Specifies a sync timeout for
applications. `0` means no timeout (default `0`)."

It's worth being clear about what this rule is flagging: `0` is Argo CD's own default, and the
documentation doesn't call it a mistake. The finding says you've written that default into
`argocd-cmd-params-cm` explicitly, which is a good moment to decide whether you meant to.

So if you really want this value to be 0 make sure to understand why you want it that way.

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

Omitting the key entirely also passes: the rule only reports a value it can parse as the integer
zero. It's inactive when `argocd-cmd-params-cm` isn't among the scanned files.

## Exceptions

Installations with genuinely long-running syncs (large Helm charts with slow hooks, or database
migrations run as PreSync Jobs) may find any timeout they pick either too short for the worst case
or too long to be useful. Leaving it unbounded is a defensible choice there, as long as something
else watches for Applications stuck in `Progressing`.

## Further reading

For more information, see the
[`argocd-cmd-params-cm` reference](https://argo-cd.readthedocs.io/en/latest/operator-manual/argocd-cmd-params-cm-yaml/)
in the Argo CD documentation, which documents the `controller.sync.timeout.seconds` key.

[Back to all rules](index.md)
