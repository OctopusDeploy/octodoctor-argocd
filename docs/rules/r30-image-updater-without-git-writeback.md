# R30 – Image Updater without Git writeback

!!! danger "Severity: High"

Applications using Image Updater without Git writeback. This is against GitOps principles.

## Why this matters

Argo CD Image Updater watches your container registry and rolls out new image tags as they appear.
It can record the new tag in one of 2 places. With `write-back-method: git` it commits back to
your repository, so Git keeps describing what's running. With the default, `argocd`, it patches the
live Application through the Argo CD API and writes nothing down.

The default breaks the single thing GitOps guarantees. Your repository now says one image tag and
your cluster runs another, with no commit explaining the difference and no history of the tags that
came before. Roll the Application back from Git and you don't get the version you were running an
hour ago, you get whatever the manifest says. Recreate the Argo CD installation from scratch, which
is meant to be a routine disaster-recovery step, and every image silently reverts. Auditing which
image was live during an incident stops being a `git log` and becomes an archaeology exercise.

It's the same failure as changing an Application through the UI, just automated and running
continuously, which is why this rule is High severity.

## How to fix it

Add the `write-back-method: git` annotation so every update becomes a commit. Image Updater then
needs write credentials for the repository, and it's worth pointing it at a dedicated branch or
using `write-back-target` so its commits are easy to distinguish from human ones.

**Flagged by this rule**

```yaml
metadata:
  annotations:
    argocd-image-updater.argoproj.io/image-list: example=ghcr.io/acme/something
    argocd-image-updater.argoproj.io/example.update-strategy: latest
```

**Passes**

```yaml
metadata:
  annotations:
    argocd-image-updater.argoproj.io/image-list: example=ghcr.io/acme/something
    argocd-image-updater.argoproj.io/example.update-strategy: latest
    argocd-image-updater.argoproj.io/write-back-method: git
```

The rule looks for any annotation starting with `argocd-image-updater.argoproj.io` and then checks
for `write-back-method` set to exactly `git`. Applications with no Image Updater annotations are
inactive, so this only reports Applications that have opted into the tool.

## Exceptions

A short-lived preview environment that's deleted along with its pull request has nothing to roll
back to and no history worth keeping, so writeback matters less there. For anything you'd want to
recover after losing the cluster, commit the tag.

## Further reading

For more information, see the Octopus Deploy blog post
[Top 30 Argo CD anti-patterns to avoid when adopting GitOps](https://octopus.com/blog/30-argo-cd-antipatterns-for-gitops).

[Back to all rules](index.md)
