# R11 – Anonymous access enabled

**Severity:** High ⚠️

ArgoCD anonymous user access is enabled, which poses a security risk as it allows unauthenticated access to the system.

## Why this matters

`users.anonymous.enabled` does exactly what the Argo CD documentation says: "Enabling anonymous
access to the Argo CD instance allows users to assume the default role permissions specified by
`policy.default` without being authenticated."

Two things make that riskier than it sounds. The first is that anonymous users inherit whatever
`policy.default` grants, and that setting is usually chosen with authenticated users in mind. If
`policy.default` is `role:readonly` — a common and otherwise sensible choice — then turning on
anonymous access hands every unauthenticated visitor a read-only tour of your deployments: every
Application, every repository URL, every cluster and namespace you deploy to, and the live manifests
of everything running. That's a detailed map of your infrastructure available without a login.

The second is that the documentation warns default-policy permissions "cannot be blocked by a `deny`
rule." You can't grant the default role broadly and then carve anonymous users back out. Whatever
the default gives, anonymous users get.

There's also no audit trail. Actions taken anonymously aren't attributable, so if anonymous users
can do anything beyond read, you lose the record of who did it.

## How to fix it

Set the key to `false`, or remove it — Argo CD disables anonymous access by default.

**Flagged by this rule**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cm
data:
  users.anonymous.enabled: "true"
```

**Passes**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cm
data:
  users.anonymous.enabled: "false"
```

Omitting the key passes too. The rule is inactive when `argocd-cm` isn't among the scanned files.

## Exceptions

If you want an unauthenticated read-only dashboard — a status wall in the office, say — the Argo CD
documentation recommends not reusing the default role for it. Create a dedicated
`role:unauthenticated` in `argocd-rbac-cm` granting the narrowest set of permissions that view
needs, and point `policy.default` at that instead of at `role:readonly`. Combine it with network
controls so the instance isn't reachable from outside your network, and remember that repository
URLs and cluster names are visible to anyone who can load the page.

## Further reading

For more information, see the Argo CD documentation on
[anonymous access](https://argo-cd.readthedocs.io/en/latest/operator-manual/rbac/#anonymous-access)
and the
[`argocd-cm` reference](https://argo-cd.readthedocs.io/en/latest/operator-manual/argocd-cm-yaml/)
for the `users.anonymous.enabled` key.

[Back to all rules](index.md)
