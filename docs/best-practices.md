# Best practices articles for Argo CD

We have published several best practices articles for Argo CD over the years.
Octodoctor is an implementation of the knowledge (and opinions) expressed in the articles. Here is the full list

1.  **[Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops)**

    Thirty distinct Argo CD anti-patterns, from misunderstanding the declarative model and abusing
    `targetRevision` to hardcoding Helm data in Application manifests and misusing finalizers. Most of
    the other posts on this page expand a single one of these items into a full article.

2.  **[How to model your GitOps environments](https://octopus.com/blog/how-to-model-your-gitops-environments)**

    A worked example of modelling five environments as folders on a single Git branch, covering initial
    deployment, diffing environments, and promotion.    

3.  **[Structuring Argo CD repositories](https://octopus.com/blog/how-to-structure-your-argo-cd-repositories-using-application-sets)**

    How to lay out the Git repositories behind Argo CD: four anti-patterns to avoid, then a recommended
    three-level structure with one repository per team.

4.  **[Distribute your Argo CD applications to different Kubernetes clusters using ApplicationSets](https://octopus.com/blog/argocd-clusters-labels-with-apps)**

    Labelling clusters by role, cloud, region and environment so ApplicationSets can select cluster
    groups instead of individual snowflakes. Nine worked scenarios cover selectors, exclusions, and
    per-environment Kustomize overlays and Helm values.

5.  **[Stop using branches for deploying to different GitOps environments](https://octopus.com/blog/stop-using-branches-deploying-different-gitops-environments)**

    Why a config repository with `staging`, `qa` and `production` branches is a trap. Promotion is
    never a simple merge, drift accumulates, and the model fights Helm and Kustomize.

6.  **[Using Helm hierarchies in multi-source Argo CD applications for promoting to different GitOps environments](https://octopus.com/blog/helm-values-argocd)**

    Building a hierarchy of Helm values files with shared defaults, environment overrides and how to use it
    with multi-source Applications and ApplicationSets.

7.  **[Securing Argo CD in a multi-tenant environment with Application Projects](https://octopus.com/blog/multi-tenant-argocd-with-application-projects)**

    Isolating teams that share one Argo CD installation with AppProjects, RBAC policies and roles, ending
    in a complete configuration for two developer teams.


8.  **[How to handle secrets like a pro using GitOps](https://octopus.com/blog/handle-secrets-like-pro-using-gitops)**

    The starting point for secrets: Bitnami Sealed Secrets, for when you want secrets in Git but not in
    plaintext.

9.  **[GitOps secrets with Argo CD, HashiCorp Vault, and External Secrets Operator](https://octopus.com/blog/gitops-secrets-hashicorp-vault-argo-cd)**

    The other approach that keeps secrets out of Git entirely and pull them from Vault at runtime, refreshing
    them without restarting pods.

10. **[Securing your GitOps secrets with the Secret Store CSI driver](https://octopus.com/blog/gitops-secrets-csi-driver)**

    The narrowest of the three secrets posts: mounting Vault secrets onto pods as files with the CSI
    driver, avoiding Kubernetes Secret objects and Argo CD syncs altogether.

11. **[Unlimited preview environments with Kubernetes namespaces](https://octopus.com/blog/unlimited-preview-environments)**

    The general namespace-per-pull-request pattern, independent of Argo CD that shows naming, ingress routing, PR
    comments, smoke tests and teardown.

12. **[Create temporary Argo CD preview environments based on pull requests](https://octopus.com/blog/creating-temporary-preview-environments-based-pull-requests-argo-cd)**

    The same pattern built on Argo CD's Pull Request generator, which creates and removes Applications
    as pull requests open and close.

13. **[How to simplify your Kubernetes Helm deployments](https://octopus.com/blog/simplify-kubernetes-helm-deployments)**

    Three techniques for a growing Helm setup: a chart repository, umbrella charts to stop version
    sprawl, and a deployment monorepo.

14. **[Using Helm to deploy Kubernetes to multiple environments](https://octopus.com/blog/helm-deployment-environments)**

    An introductory walkthrough of deploying one application to several environments with Helm values,
    then moving from predefined environments to dynamic ones.

15. **[How to preview and diff your Argo CD deployments](https://octopus.com/blog/argo-cd-preview-diff)**

    Four ways to see what a sync will change before it happens.  GUI diff, CLI local diff, pre-rendered
    manifests, and rendering on the fly (also attaching the diff to a Pull Request).

16. **[Argo CD application dependencies](https://octopus.com/blog/argo-cd-application-dependencies)**

    What to do when one Application must deploy before another: App-of-Apps, sync waves, and where
    ApplicationSets replace them.

17. **[Using GitOps for infrastructure and applications with Crossplane and Argo CD](https://octopus.com/blog/using-gitops-infrastructure-applications-crossplane-argo-cd)**

    Managing cloud infrastructure through the same GitOps workflow as your applications by pairing
    Crossplane with Argo CD.

18. **[Database migrations in the era of Kubernetes microservices](https://octopus.com/blog/database-migrations-kubernetes-microservices)**

    Giving schema migrations their own lifecycle instead of running them from Argo CD pre-sync hooks,
    with generated SQL previewed against a temporary database first.

19. **[Abusing the target revision field for Argo CD promotions](https://octopus.com/blog/argocd-application-target-revision-field)**

    An entire article about one field: point `targetRevision` at Git HEAD, use tags when you must, pin
    infrastructure charts and never use environment names, semantic ranges or commit hashes.

20. **[Argo CD finalizers: prevent accidental deletion issues](https://octopus.com/blog/argocd-application-deletion-finalizers)**

    The mysteries of Argo CD finalizers and how you can shoot yourself in the foot with them.

21. **[Kubernetes deployment anti-patterns – Part 1](https://octopus.com/blog/kubernetes-antipatterns-1)**

    Anti-patterns 1 to 5: `latest` tags, configuration baked into images, needless coupling to
    Kubernetes features, mixing application and infrastructure deployment, and hand-editing live
    resources with `kubectl edit` or `patch`.

22. **[Kubernetes deployment anti-patterns – Part 2](https://octopus.com/blog/kubernetes-antipatterns-2)**

    Anti-patterns 6 to 10, mostly about operations:  `kubectl` as a debugger, misunderstood networking,
    permanent staging instead of dynamic environments, production and non-production on one cluster, and
    workloads with no memory or CPU limits.

23. **[Kubernetes deployment anti-patterns – Part 3](https://octopus.com/blog/kubernetes-antipatterns-3)**

    Anti-patterns 11 to 15: misused health probes, skipping Helm, no deployment metrics, no strategy
    for secrets, and trying to solve every problem with Kubernetes.

24. **[10 Docker anti-patterns](https://octopus.com/blog/docker-anti-patterns)**

    Ten container-image bad practices with a fix for each, including treating containers as virtual
    machines, opaque images, a different image per environment, tracking Git hashes instead of image
    tags, and hardcoded secrets. This is the layer beneath everything else here and several Argo CD
    problems start as an image built the wrong way.

25. **[Enterprise CI/CD best practices – Part 1](https://octopus.com/blog/enterprise-ci-cd-best-practices-1)**

    Practices 1 to 7 of 23: everything in source control, one artifact promoted unchanged to every
    environment, artifacts rather than Git commits travelling the pipeline, short-lived branches,
    single-step and fast builds, and cached dependencies. Not Argo CD specific, but the layer your
    GitOps setup sits on top of.

26. **[Enterprise CI/CD best practices – Part 2](https://octopus.com/blog/enterprise-ci-cd-best-practices-part-2)**

    Practices 8 to 15, all about testing: automate everything, keep tests fast, clean up side effects,
    split into multiple suites, create test environments on demand, run suites concurrently, and treat
    security and quality scanning as pipeline stages.

27. **[Enterprise CI/CD best practices – Part 3](https://octopus.com/blog/enterprise-ci-cd-best-practices-part-3)**

    Practices 16 to 23, covering deployment and after: database changes on their own lifecycle and
    applied gradually, deploying only through the CD platform, progressive delivery, metrics that catch
    a bad release, automatic rollbacks, and staging that matches production.