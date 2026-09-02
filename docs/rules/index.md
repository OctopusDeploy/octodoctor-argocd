# Detection rules

![OctoDoctor for Argo CD logo](../octodoctor-logo.png)

The tool will look at Argo CD application manifests as yaml files and detect the following

| Rule Number | Description | Severity | Implemented |
|-------------|-------------|----------|-----------|
| [R01](r01-dangerous-sync-policy.md) | Found Applications with dangerous sync policy settings (prune=true and allowEmpty=true) | High ⚠️ | ✅ |
| [R02](r02-no-automated-sync.md) | Found Applications without automated sync policy configured | Medium | ✅ |
| [R03](r03-low-requeue-after-seconds.md) | Application Sets with requeueAfterSeconds setting of less than 1 minute | Medium | ✅ |
| [R04](r04-sync-policy-replace.md) | Found Applications with sync policy replace = true | Medium | ✅ |
| [R05](r05-permissive-source-repos.md) | Found AppProjects with overly permissive sourceRepos settings (using wildcard '*') | High ⚠️ | ✅ |
| [R06](r06-missing-manifest-generate-paths.md) | Not using `argocd.argoproj.io/manifest-generate-paths` annotation | Medium | ✅ |
| [R07](r07-reconciliation-timeout-without-cache-expiration.md) | Using `ARGOCD_RECONCILIATION_TIMEOUT=0` without `ARGOCD_DEFAULT_CACHE_EXPIRATION ` | Medium | ✅ |
| [R09](r09-duplicate-generators.md) | Found Application Sets with duplicate generators | Medium | ✅ |
| [R10](r10-sync-timeout-of-zero.md) | ArgoCD controller sync timeout set to 0 seconds | Medium | ✅ |
| [R11](r11-anonymous-access.md) | ArgoCD anonymous user access is enabled | High ⚠️ | ✅ |
| [R12](r12-nested-generators.md) | Found Application Sets with 3 levels of generators | Medium | ✅ |
| [R14](r14-prune-enabled.md) | Found Applications with prune = true | Medium | ✅ |
| [R15](r15-hardcoded-helm-values.md) | Found Application CRDs with hardcoded values in helm: property | Medium | ✅ |
| [R16](r16-hardcoded-kustomize-values.md) | Found Application CRDs with hardcoded values in kustomize: property | Medium | ✅ |
| [R17](r17-parameter-overrides.md) | Found Application CRDs with parameter overrides | High ⚠️ | ✅ |
| [R18](r18-scm-and-pull-request-generators.md) | Found an Application Set that combines SMC generator and Pull Request generator | Medium | ✅ |
| [R19](r19-pinned-target-revision.md) | Using Git hashes or branch names in Target Revision (anything else other than `HEAD`) | Medium | ✅ |
| [R20](r20-no-cluster-generator.md) | Having more than 2 connected clusters and applicationsets don't use the cluster generator | Medium | ✅ |
| [R22](r22-no-webhooks-configured.md) | No webhooks setup for Argo CD repo server/Git generator | Medium | ✅ |
| [R25](r25-more-than-two-sources.md) | Having more than 2 sources in a multi-source app | Medium | ✅ |
| [R26](r26-multiple-helm-chart-sources.md) | Having multiple sources with just Helm charts | Medium | ✅ |
| [R30](r30-image-updater-without-git-writeback.md) | Using Image Updater without writing back to Git | High ⚠️ | ✅ |
| [R50](r50-invalid-resource-tracking-method.md) | ArgoCD resourceTrackingMethod is set to an invalid value | High ⚠️ | ✅ |
| [R97](r97-exactly-one-application-set.md) | Having exactly one application set (results in complex YAML) | Medium | ✅ |
| [R98](r98-no-application-sets.md) | No Application Set resources while Application Resources are > 0 | Medium | ✅ |
| [R99](r99-too-many-application-sets.md) | Having too many application sets (performance issues) | Medium | ✅ |
