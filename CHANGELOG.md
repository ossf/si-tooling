# Changelog

## v2.3.0

This release expands on the Security Insights v2.1.0 specification by adding a new micro-CLI for non-library use cases. It also contains discovery and fetching improvements to the library.

### Features

- `si.Discover` finds a repository's Security Insights file among the known filename and location variants using at most two GitHub API calls
- `si.Fetch` returns the raw bytes of a file in a public GitHub repository; `si.ErrNotFound` identifies a missing file
- `si fetch` command (`v2/cmd/si`) runs discovery and parsing over many targets and prints machine-readable JSON, for catalogues and dashboards built on Security Insights data. Each result carries a status of `ok`, `not-found`, `invalid` or `error`; the command exits 2 when any target hit a transport `error`, so an incomplete batch is never mistaken for a clean one
- `si.Load` returns `si.ErrParentUnavailable` when the parent named in `project-si-source` cannot be fetched for a reason unrelated to its content; `si fetch` reports that as `error`, not `invalid`
- All GitHub API calls, including the existing `si.Read`, now authenticate with `GITHUB_TOKEN` when set and honour `GITHUB_API_URL`; all requests time out after 30 seconds

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
