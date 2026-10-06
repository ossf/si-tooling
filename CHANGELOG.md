# Changelog

## Unreleased

### Features

- `si.Discover` finds a repository's Security Insights file among the known filename and location variants (`si.DiscoveryPaths`) using two GitHub API calls
- `si.Fetch` returns the raw bytes of a file in a public GitHub repository; `si.ErrNotFound` identifies a missing file
- `si fetch` command (`v2/cmd/si`) runs discovery and parsing over many targets and prints machine-readable JSON, for catalogues and dashboards built on Security Insights data
- GitHub API calls authenticate with `GITHUB_TOKEN` when set and honour `GITHUB_API_URL`

## v2.1.0

This release is made in concert with the v2.1.0 release of [Security Insights](https://github.com/ossf/security-insights-spec)

### Features

- Adds support for the `project.steward` field introduced in Insights v2.1.0
- `const SecurityInsightsFilename = "security-insights.yml"` now exported by `v2/si` for use by consuming applications
- `si.Read` now relies on `github.com/google/go-github/v71` for reading `security-insights.yml` files from public GitHub repositories

### Quality

- Added GitHub Actions CI workflow to lint all `go`
- Improved project documentation

### Security

- Replaced unmaintained dependency `gopkg.in/yaml.v3` with `github.com/goccy/go-yaml`
- Added `.github/security-insights.yml` to communicate the project's security posture
