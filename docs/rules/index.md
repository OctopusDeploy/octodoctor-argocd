# Detection rules

![OctoDoctor for Argo CD logo](../octodoctor-logo.png)

OctoDoctor reads your Argo CD manifests as YAML files and checks them against the following rules:

| Rule | Description | Severity | Implemented |
|-------------|-------------|----------|-----------|
| [R01](r01-dangerous-sync-policy.md) | Applications with a dangerous sync policy (`prune=true` and `allowEmpty=true`) | High ⚠️ | ✅ |
| [R02](r02-no-automated-sync.md) | Applications without an automated sync policy | Medium | ✅ |
| [R03](r03-low-requeue-after-seconds.md) | ApplicationSets with a `requeueAfterSeconds` of less than 3 minutes (180 seconds) | Medium | ✅ |
| [R04](r04-sync-policy-replace.md) | Applications with the sync option `Replace=true` | Medium | ✅ |
| [R05](r05-permissive-source-repos.md) | AppProjects with an overly permissive `sourceRepos` (using the wildcard `*`) | High ⚠️ | ✅ |
| [R06](r06-missing-manifest-generate-paths.md) | Applications without the `argocd.argoproj.io/manifest-generate-paths` annotation | Medium | ✅ |
| [R07](r07-reconciliation-timeout-without-cache-expiration.md) | `ARGOCD_RECONCILIATION_TIMEOUT=0` without `ARGOCD_DEFAULT_CACHE_EXPIRATION` | Medium | ✅ |
| [R09](r09-duplicate-generators.md) | ApplicationSets with duplicate generators | Medium | ✅ |
| [R10](r10-sync-timeout-of-zero.md) | Argo CD controller sync timeout set to 0 seconds | Medium | ✅ |
| [R11](r11-anonymous-access.md) | Argo CD anonymous user access turned on | High ⚠️ | ✅ |
| [R12](r12-nested-generators.md) | ApplicationSets with 3 levels of generators | Medium | ✅ |
| [R14](r14-prune-enabled.md) | Applications with `prune=true` | Medium | ✅ |
| [R15](r15-hardcoded-helm-values.md) | Applications with hardcoded values in the `helm` property | Medium | ✅ |
| [R16](r16-hardcoded-kustomize-values.md) | Applications with hardcoded values in the `kustomize` property | Medium | ✅ |
| [R17](r17-parameter-overrides.md) | Applications with parameter overrides | High ⚠️ | ✅ |
| [R18](r18-scm-and-pull-request-generators.md) | ApplicationSets that combine the SCM Provider and Pull Request generators | Medium | ✅ |
| [R19](r19-pinned-target-revision.md) | Git hashes or branch names in `targetRevision` (anything other than `HEAD`) | Medium | ✅ |
| [R20](r20-no-cluster-generator.md) | 2 or more connected clusters, and ApplicationSets that don't use the cluster generator | Medium | ✅ |
| [R22](r22-no-webhooks-configured.md) | No webhooks set up for the Argo CD repo server or Git generator | Medium | ✅ |
| [R25](r25-more-than-two-sources.md) | Multi-source Applications with more than 2 sources | Medium | ✅ |
| [R26](r26-multiple-helm-chart-sources.md) | Multi-source Applications with more than one Helm chart source | Medium | ✅ |
| [R30](r30-image-updater-without-git-writeback.md) | Image Updater without writing back to Git | High ⚠️ | ✅ |
| [R50](r50-invalid-resource-tracking-method.md) | Argo CD `resourceTrackingMethod` set to an invalid value | High ⚠️ | ✅ |
| [R97](r97-exactly-one-application-set.md) | Exactly one ApplicationSet (which leads to complex YAML) | Medium | ✅ |
| [R98](r98-no-application-sets.md) | No ApplicationSets while there are Applications | Medium | ✅ |
| [R99](r99-too-many-application-sets.md) | Too many ApplicationSets (performance issues) | Medium | ✅ |
