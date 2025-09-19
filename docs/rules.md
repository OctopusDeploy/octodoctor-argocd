# Argo CD detection rules

The tool will look at the extracted contents of a Codefresh support package and detect the following

| Rule Number | Description | Severity | Implemented |
|-------------|-------------|----------|-----------|
| R01 | Found Applications with dangerous sync policy settings (prune=true and allowEmpty=true) | High | YES |
| R02 | Found Applications without automated sync policy configured | Medium | YES |
| R03 | Application Sets with requeueAfterSeconds setting of less than 1 minute | Medium | NO |
| R04 | Found Applications with sync policy replace = true | Medium | NO |
| R05 | Found AppProjects with overly permissive sourceRepos settings (using wildcard '*') | High | YES |
| R06 | Not using `argocd.argoproj.io/manifest-generate-paths` annotation | Medium | YES |
| R07 | Using `ARGOCD_RECONCILIATION_TIMEOUT=0` without `--revision-cache-expiration` | Medium | NO |
| R08 | No Application Set resources while Application Resources are > 0 | Medium | NO |
| R09 | Found Application Sets with duplicate generators | Medium | NO |
| R10 | ArgoCD controller sync timeout set to 0 seconds | Medium | YES |
| R11 | ArgoCD anonymous user access is enabled | High | YES |
| R12 | Found Application Sets with 3 levels of generators | Medium | NO |
| R13 | Found Applications with auto-sync off | Medium | NO |
| R14 | Found Applications with prune = true | Medium | NO |
| R15 | Found Application CRDs with hardcoded values in helm: property | Medium | NO |
| R16 | Found Application CRDs with hardcoded values in kustomize: property | Medium | NO |
| R17 | Found Application CRDs with parameter overrides | Medium | NO |
| R18 | Found an Application Set that combines SMC generator and Pull Request generator | Medium | NO |
| R19 | Using Git hashes or branch names in Target Revision (anything else other than `HEAD`) | Medium | NO |
| R20 | Having more than 2 connected clusters and applicationsets don't use the cluster generator | Medium | NO |
| R21 | Having only a single Application Set | Medium | NO |
| R22 | No webhooks setup for Argo CD repo server | Medium | NO |
| R97 | Having exactly one application set (results in complex YAML) | Medium | YES |
| R98 | No Application Set resources while Application Resources are > 0 | Medium | YES |
| R99 | Having too many application sets (performance issues) | Medium | YES |
