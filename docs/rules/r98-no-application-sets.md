# R98 – No ApplicationSets

!!! warning "Severity: Medium"

You aren't using ApplicationSets. ApplicationSets are a powerful feature of Argo CD that lets you deploy applications to multiple clusters and namespaces.

## Why this matters

Hand-written Application manifests are fine while you're learning, and they stop scaling almost
immediately after. Every new service, environment, or cluster means another near-identical file,
copied from the last one and edited in 3 places. The copies drift. One of them keeps a
`targetRevision` someone pinned during an incident, another never got the sync policy change that
the rest received, and nobody notices until the difference matters.

ApplicationSets replace the copying with generation. A generator discovers the inputs (directories
in a repository, labeled clusters, open pull requests, or a literal list), and one template turns
each input into an Application. A single manifest can produce a dozen Applications, and adding the
thirteenth means adding a directory or labeling a cluster, not writing more YAML. Consistency stops
being something you maintain by hand.

The other benefit is that the ApplicationSet becomes the readable description of your topology. One
file tells you which applications run where, instead of that fact being distributed across 50
Application manifests.

## How to fix it

Pick the generator that matches how your repository is already organized, and template the
Application you're currently copying. A Git directory generator is the usual starting point.

**Flagged by this rule**

```
gitops-repo/
└── applications/
    ├── billing-qa.yaml          # hand-written, near-identical
    ├── billing-prod.yaml
    ├── checkout-qa.yaml
    └── checkout-prod.yaml
```

**Passes**

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: team-revenue
spec:
  generators:
    - matrix:
        generators:
          - git:
              repoURL: https://github.com/acme/gitops.git
              revision: HEAD
              directories:
                - path: apps/*
          - clusters:
              selector:
                matchLabels:
                  env: prod
  template:
    metadata:
      name: '{{path.basename}}-{{name}}'
    spec:
      project: revenue
      source:
        repoURL: https://github.com/acme/gitops.git
        targetRevision: HEAD
        path: '{{path}}'
      destination:
        server: '{{server}}'
        namespace: '{{path.basename}}'
```

This rule looks at the installation as a whole rather than at individual files. It fails when the
scanned directory contains one or more Application resources but no ApplicationSet resources. It is
inactive when the directory contains neither, for example an export that holds only AppProjects or
configuration, because there are no Applications to generate.

## Exceptions

A directory holding only a couple of Applications, or one scanned in isolation from the rest of your
GitOps repository, will trigger this without anything being wrong. Check what the tool actually
scanned before acting on it.

## Further reading

For more information, see these Octopus Deploy blog posts:

- [Structuring Argo CD repositories](https://octopus.com/blog/how-to-structure-your-argo-cd-repositories-using-application-sets)
- [Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops)

[Back to all rules](index.md)
