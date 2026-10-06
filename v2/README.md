# v2/si-tooling [![Go Reference](https://pkg.go.dev/badge/github.com/ossf/si-tooling/v2.svg)](https://pkg.go.dev/github.com/ossf/si-tooling/v2)

This is a go module for working with [Security Insights](https://github.com/ossf/security-insights-spec) data in YAML `security-insights.yml` and Go `si.SecurityInsights`.

## Usage

Unmarshal the `security-insights.yml` data in [ossf/security-insights-spec](https://github.com/ossf/security-insights-spec)

```go
import (
    "fmt"

    "github.com/ossf/si-tooling/v2/si"
)

func main() {
    insights, err := si.Read("ossf", "security-insights-spec", ".github/security-insights.yml")
    message = fmt.Sprintf("Repository license is: %s", insights.Repository.License.Expression)
}
```

### Discovering the file

Repositories keep the file under a handful of names and locations. `si.Discover` lists the root and `.github/` directories (two API calls) and returns the first match from `si.DiscoveryPaths`; `si.Fetch` returns the raw bytes of a known path.

```go
path, err := si.Discover("ossf", "scorecard")   // ".github/security-insights.yml"
raw, err := si.Fetch("ossf", "scorecard", path)
insights, err := si.Load(raw)
```

Set `GITHUB_TOKEN` to authenticate GitHub API calls (unauthenticated requests are limited to 60 per hour). `GITHUB_API_URL`, which GitHub Actions exports, is honoured for GitHub Enterprise.

## Command line

`si fetch` runs discovery and parsing over any number of targets and prints a JSON array, in input order, for other tools to consume:

```sh
cd v2
go run ./cmd/si fetch ossf/scorecard mindersec/minder \
  https://github.com/openbao/openbao/blob/main/.github/security-insights.yml
```

Targets may be `owner/repo` (discover the path), `owner/repo/<path>`, a `github.com/.../blob/...` URL, or a `raw.githubusercontent.com` URL. URLs are read from the default branch; the ref in the URL is ignored.

Each result looks like:

```json
{
  "target": "ossf/scorecard",
  "owner": "ossf",
  "repo": "scorecard",
  "path": ".github/security-insights.yml",
  "status": "ok",
  "schema-version": "2.0.0",
  "insights": { "header": { ... }, "project": { ... }, "repository": { ... } }
}
```

`status` is `ok`, `not-found` (no file at the path, or discovery found nothing) or `invalid` (the file exists but does not parse as Security Insights v2; `error` says why and `schema-version` is still reported when the header declares one). The exit code is 0 whenever every target was processed; only a usage error exits non-zero. Progress is written to stderr.

## Schema version support

The module supports Security Insights schema version 2.x, including v2.2.0:

> [!WARNING]
> Security Insights **v2.2.0** `vulnerability-reporting.policy` replaces the former `security-policy` field under vulnerability reporting. This backwards compatibility violation was tolerated by the Security Insights maintainers due to the lack of evidence that the former field had been adopted by end users. Issues may arise if users of SI Tooling fail to update to the latest version _and_ the users of Security Insights specification begin to use the new field.
