# Argo CD detection rules

The tool will look at the extracted contents of a Codefresh support package and detect the following


## High 

1. No webhooks setup for Argo CD repo server
1. No webhooks setup for Argo CD application Set server when using Application Sets
1. Application Sets with requeueAfterSeconds setting of less than 1 minute
1. Found Applications with sync policy replace = true
1. Found Applications with sync policy allowEmpty = true
1. Using `ARGOCD_RECONCILIATION_TIMEOUT=0` without `--revision-cache-expiration`
1. Not using `argocd.argoproj.io/manifest-generate-paths` 


## Medium

1. No Application Set resources while Application Resources are > 0
1. Found Application Sets with duplicate generators
1. Found Application Sets with 3 levels of generators
1. Having only a single Application Set
1. Found Applications with auto-sync off
1. Found Applications with auto-sync on but self-heal off
1. Found Applications with prune = true
1. Found Application CRDs with hardcoded values in helm: property
1. Found Application CRDs with hardcoded values in kustomize: property
1. Found Application CRDs with parameter overrides
1. Found an Application Set that combines SMC generator and Pull Request generator
1. Using Git hashes or branch names in Target Revision (anything else other than `HEAD`)
1. Having more than 2 connected clusters and applicationsets don't use the cluster generator
