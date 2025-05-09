# Changelog

## v3.0.0

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