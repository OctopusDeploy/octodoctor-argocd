# Argo CD detection rules

![Argo CD doctor logo](argo-cd-doctor-logo.png)

The tool will look at the extracted contents of a Codefresh support package and detect the following

| Rule Number | Description | Severity | Implemented |
|-------------|-------------|----------|-----------|
| R01 | Found Applications with dangerous sync policy settings (prune=true and allowEmpty=true) | High ⚠️ | ✅ |
| R02 | Found Applications without automated sync policy configured | Medium | ✅ |
| R03 | Application Sets with requeueAfterSeconds setting of less than 1 minute | Medium | ✅ |
| R04 | Found Applications with sync policy replace = true | Medium | ✅ |
| R05 | Found AppProjects with overly permissive sourceRepos settings (using wildcard '*') | High ⚠️ | ✅ |
| R06 | Not using `argocd.argoproj.io/manifest-generate-paths` annotation | Medium | ✅ |
| R07 | Using `ARGOCD_RECONCILIATION_TIMEOUT=0` without `ARGOCD_DEFAULT_CACHE_EXPIRATION ` | Medium | ✅ |
| R09 | Found Application Sets with duplicate generators | Medium | ✅ |
| R10 | ArgoCD controller sync timeout set to 0 seconds | Medium | ✅ |
| R11 | ArgoCD anonymous user access is enabled | High ⚠️ | ✅ |
| R12 | Found Application Sets with 3 levels of generators | Medium | ✅ |
| R14 | Found Applications with prune = true | Medium | ✅ |
| R15 | Found Application CRDs with hardcoded values in helm: property | Medium | ✅ |
| R16 | Found Application CRDs with hardcoded values in kustomize: property | Medium | ✅ |
| R17 | Found Application CRDs with parameter overrides | High ⚠️ | ✅ |
| R18 | Found an Application Set that combines SMC generator and Pull Request generator | Medium | ✅ |
| R19 | Using Git hashes or branch names in Target Revision (anything else other than `HEAD`) | Medium | ✅ |
| R20 | Having more than 2 connected clusters and applicationsets don't use the cluster generator | Medium | ✅ |
| R22 | No webhooks setup for Argo CD repo server/Git generator | Medium | ✅ |
| R25 | Having more than 2 sources in a multi-source app | Medium | ✅ |
| R26 | Having multiple sources with just Helm charts | Medium | ❌ |
| R30 | Using Image Updater without writing back to Git | High ⚠️ | ✅ |
| R97 | Having exactly one application set (results in complex YAML) | Medium | ✅ |
| R98 | No Application Set resources while Application Resources are > 0 | Medium | ✅ |
| R99 | Having too many application sets (performance issues) | Medium | ✅ |
