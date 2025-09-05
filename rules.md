# Argo CD detection rules

The tool will look at the extracted contents of a Codefresh support package and detect the following


## High 

1. No webhooks setup for Argo CD repo server
2. No webhooks setup for Argo CD application Set server
3. Application Sets with reque setting of less than 1 minute
4. Found Applications with replace = true
5. Found Applications with allowEmpty = true


## Medium

1. No Application Set resources while Application Resources are > 0
2. Found Application Sets with duplicate generators
3. Found Application Sets with 3 levels of generators
4. Found Applications with auto-sync off
5. Found Applications with auto-sync on but self-heal of
6. Found Applications with prune = true
