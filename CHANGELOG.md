# Changelog

## v2.2.1

### Features

- `SecurityInsightsFilenames()` returns both accepted names, `security-insights.yml` and `security-insights.yaml`, matching the spec

### Deprecations

- `SecurityInsightsFilename` names only the `.yml` file. Use `SecurityInsightsFilenames()` instead

### Bug Fixes

- `in-scope` and `out-of-scope` load as string lists. Files that set either field no longer fail to load ([#68](https://github.com/ossf/si-tooling/issues/68))

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
