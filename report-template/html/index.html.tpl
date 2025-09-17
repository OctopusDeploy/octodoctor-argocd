<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Argo CD Doctor report</title>
  <link rel="shortcut icon" href="favicon.png" type="image/png" />
  <link rel="stylesheet" href="style.css">
</head>
<body>
  <h1>Argo CD Doctor Analysis Report</h1>
  <p class="report-info">Generated at {{.Now.Format "02 Jan 2006 15:04 MST"}}</p>

  <section class="scan-results">
    <h2>Issues found that need attention</h2>
    <table class="rules-table triggered">
      <thead>
        <tr>
          <th>Rule</th>
          <th>Description</th>
          <th>Affected Files</th>
        </tr>
      </thead>
      <tbody>
        {{range .Analysis.RulesFailed}}
        <tr class="severity-high">
          <td class="rule-id">{{.ID}}</td>
          <td class="rule-description">{{.Description}}</td>
          <td class="affected-files">
            {{range .FilesTriggered}}
            <div class="file-item">{{.}}</div>
            {{end}}
          </td>
        </tr>
        {{end}}
      </tbody>
    </table>

    <h2>Best practices followed</h2>
    <table class="rules-table not-triggered">
      <thead>
        <tr>
          <th>Rule</th>
          <th>Description</th>
        </tr>
      </thead>
      <tbody>
        {{range .Analysis.RulesPassed}}
        <tr>
          <td class="rule-id">{{.ID}}</td>
          <td class="rule-description">{{.Description}}</td>
        </tr>
        {{end}}
      </tbody>
    </table>

    <h2>Inactive/Not applicable rules</h2>
    <table class="rules-table inactive">
      <thead>
        <tr>
          <th>Rule</th>
          <th>Description</th>
        </tr>
      </thead>
      <tbody>
        {{range .Analysis.RulesInactive}}
        <tr>
          <td class="rule-id">{{.ID}}</td>
          <td class="rule-description">{{.Description}}</td>
        </tr>
        {{end}}
      </tbody>
    </table>
  </section>
</body>
</html>