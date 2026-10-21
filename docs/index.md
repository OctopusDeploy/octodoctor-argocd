# OctoDoctor for Argo CD

![OctoDoctor for Argo CD logo](octodoctor-logo.png)

OctoDoctor for Argo CD is a static analysis tool for Argo CD configuration. Point it at a directory of
YAML manifests and it inspects your Applications, ApplicationSets, AppProjects, and Argo CD
ConfigMaps/Secrets, then reports misconfigurations and anti-patterns in your GitOps setup.

Every finding comes from a [rule](rules/index.md), and OctoDoctor writes the results as an HTML report
alongside machine-readable YAML and JSON.

You can use the tool in 2 modes:

- Read local Argo CD manifests (`--directory`)
- Export from a live Argo CD cluster (`--export`)

In export mode, OctoDoctor reaches the Argo CD cluster through your current kubeconfig.

[![Flow](octodoctor-flow.svg)](octodoctor-flow.svg)

OctoDoctor analyzes your files or cluster and creates a report in HTML, JSON, and YAML.

## Example report

Here's an example of the HTML report.

[![OctoDoctor for Argo CD report](example-report.png)](example-report.png)

## Next steps

- [Install OctoDoctor](installation.md)
- [Use OctoDoctor](usage.md)
- [Browse the detection rules](rules/index.md)
