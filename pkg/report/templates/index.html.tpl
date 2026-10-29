<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>OctoDoctor for Argo CD report</title>
  <link rel="shortcut icon" href="favicon.ico" type="image/x-icon">
  <link rel="stylesheet" href="style.css">
</head>
<body>
  <header class="site-header">
    <div class="container">
      <img class="logo" src="octodoctor.png" alt="OctoDoctor for Argo CD logo" width="112" height="118">
      <div class="titles">
        <span class="eyebrow">OctoDoctor for Argo CD</span>
        <h1>Analysis report</h1>
        <p class="report-info">Generated {{.Now.Format "Monday, January 2, 2006 at 3:04pm MST"}}</p>
      </div>
      <nav class="header-actions" aria-label="Project links">
        <a class="btn btn-quiet" href="https://octodoctor-argocd.readthedocs.io/" target="_blank" rel="noopener">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/></svg>
          Documentation
        </a>
        <a class="btn btn-secondary" href="https://github.com/OctopusDeploy/octodoctor-argocd" target="_blank" rel="noopener">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"/><path d="M9 18c-4.51 2-5-2-7-2"/></svg>
          View OctoDoctor on GitHub
        </a>
      </nav>
    </div>
  </header>

  <main class="container">
    <div class="summary">
      <div class="card stat">
        <span class="stat-label">Issues found</span>
        <span class="stat-value is-danger">{{len .Analysis.RulesFailed}}</span>
        <span class="stat-note">{{sevCount .Analysis.RulesFailed "high"}} high, {{sevCount .Analysis.RulesFailed "medium"}} medium{{with sevCount .Analysis.RulesFailed "low"}}, {{.}} low{{end}}</span>
      </div>
      <div class="card stat">
        <span class="stat-label">Best practices followed</span>
        <span class="stat-value is-success">{{len .Analysis.RulesPassed}}</span>
        <span class="stat-note">Rules that passed</span>
      </div>
      <div class="card stat">
        <span class="stat-label">Not applicable</span>
        <span class="stat-value">{{len .Analysis.RulesInactive}}</span>
        <span class="stat-note">Rules with nothing to check</span>
      </div>
    </div>

    <section class="section scan-results">
      <div class="section-heading">
        <h2>Issues that need attention</h2>
        <p class="section-description">Select a rule ID to read what it checks and how to fix it.</p>
      </div>
      <div class="card table-wrap">
        <table class="rules-table issues triggered">
          <thead><tr>
            <th class="col-status" scope="col">Severity</th>
            <th class="col-rule" scope="col">Rule</th>
            <th scope="col">Description</th>
            <th class="col-files" scope="col">Affected files</th>
          </tr></thead>
          <tbody>
          {{range .Analysis.RulesFailed}}
          <tr class="severity-{{.Severity}}">
            <td>{{if eq .Severity "high"}}<span class="badge badge-danger">High</span>{{else if eq .Severity "medium"}}<span class="badge badge-warning">Medium</span>{{else}}<span class="badge badge-info">Low</span>{{end}}</td>
            <td><a class="rule-id" href="{{.DocumentationURL}}" target="_blank" rel="noopener" title="{{.Name}}">{{.ID}}</a></td>
            <td class="rule-description">{{.Description}}</td>
            <td><div class="affected-files">
              {{range head $.Limit .FilesTriggered}}<code class="file-item">{{.}}</code>
              {{end}}{{with more $.Limit .FilesTriggered}}<span class="more-files">+{{.}} more files</span>{{end}}
            </div></td>
          </tr>
          {{else}}
          <tr><td colspan="4">No rules in this category.</td></tr>
          {{end}}
          </tbody>
        </table>
      </div>
    </section>

    <section class="section">
      <h2>Best practices followed</h2>
      <div class="card table-wrap">
        <table class="rules-table not-triggered">
          <thead><tr>
            <th class="col-status" scope="col">Status</th>
            <th class="col-rule" scope="col">Rule</th>
            <th scope="col">Description</th>
          </tr></thead>
          <tbody>
          {{range .Analysis.RulesPassed}}
          <tr>
            <td><span class="badge badge-success">Passed</span></td>
            <td><a class="rule-id" href="{{.DocumentationURL}}" target="_blank" rel="noopener" title="{{.Name}}">{{.ID}}</a></td>
            <td class="rule-description">{{.Description}}</td>
          </tr>
          {{else}}
          <tr><td colspan="3">No rules in this category.</td></tr>
          {{end}}
          </tbody>
        </table>
      </div>
    </section>

    <section class="section">
      <h2>Inactive or not applicable rules</h2>
      <div class="card table-wrap">
        <table class="rules-table inactive">
          <thead><tr>
            <th class="col-status" scope="col">Status</th>
            <th class="col-rule" scope="col">Rule</th>
            <th scope="col">Description</th>
          </tr></thead>
          <tbody>
          {{range .Analysis.RulesInactive}}
          <tr>
            <td><span class="badge">Inactive</span></td>
            <td><a class="rule-id" href="{{.DocumentationURL}}" target="_blank" rel="noopener" title="{{.Name}}">{{.ID}}</a></td>
            <td class="rule-description">{{.Description}}</td>
          </tr>
          {{else}}
          <tr><td colspan="3">No rules in this category.</td></tr>
          {{end}}
          </tbody>
        </table>
      </div>
    </section>

    <footer class="site-footer">
      <div class="maintainer">
        <img src="argo-cd-maintainer-badge.png" alt="Argo CD Maintainer badge" width="72" height="77">
        <div>
          <strong>Built by Argo CD maintainers at Octopus Deploy</strong>
          <span>OctoDoctor for Argo CD is open source. Found a problem or have an idea for a rule?</span>
        </div>
      </div>
      <nav class="footer-links" aria-label="More from Octopus">
        <a class="external-link" href="https://github.com/OctopusDeploy/octodoctor-argocd" target="_blank" rel="noopener noreferrer">OctoDoctor on GitHub<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M7 17L17 7M9 7h8v8"/></svg></a>
        <a class="external-link" href="https://octodoctor-argocd.readthedocs.io/" target="_blank" rel="noopener noreferrer">Documentation<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M7 17L17 7M9 7h8v8"/></svg></a>
        <a class="external-link" href="https://github.com/OctopusDeploy/OpenSource" target="_blank" rel="noopener noreferrer">Learn about Octopus Deploy open source<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M7 17L17 7M9 7h8v8"/></svg></a>
      </nav>
    </footer>

    <div class="brand-row">
      <a class="brand" href="https://octopus.com" target="_blank" rel="noopener" aria-label="Octopus Deploy">
        <img src="octopus-mark.svg" alt="" width="32" height="32">
        <span>Octopus Deploy</span>
      </a>
      <span class="sign-off">Happy deployments!</span>
    </div>
  </main>
</body>
</html>
