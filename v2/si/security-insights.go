package si

// SecurityInsightsFilename is the expected name of the YAML file containing the Security Insights data. See https://github.com/ossf/security-insights-spec?tab=readme-ov-file#usage-by-project-maintainers for more details.
const SecurityInsightsFilename = "security-insights.yml"

func (u URL) String() string {
	return string(u)
}

func (e Email) String() string {
	return string(e)
}

func (d Date) String() string {
	return string(d)
}

func (sv SchemaVersion) String() string {
	return string(sv)
}
