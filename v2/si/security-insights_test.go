package si

import (
	"testing"
)

func TestValidate(t *testing.T) {
	// Test cases for validation
	testCases := []struct {
		name       string
		si         SecurityInsights
		valid      bool
		skip       bool
		skipReason string
	}{
		{
			name: "Valid Security Insights - barebones",
			si: SecurityInsights{
				Header: Header{
					SchemaVersion: "2.0.0",
					LastUpdated:   "2021-09-01",
					LastReviewed:  "2021-09-01",
					URL:           "https://example.com/security-insights",
				},
			},
			valid: true,
		},
		{
			name: "Valid Security Insights - project barebones",
			si: SecurityInsights{
				Header: Header{
					SchemaVersion: "2.0.0",
					LastUpdated:   "2021-09-01",
					LastReviewed:  "2021-09-01",
					URL:           "https://example.com/security-insights",
				},
				Project: Project{
					Name: "Example Project",
					Administrators: []Contact{
						{
							Name:    "Example Admin",
							Primary: true,
						},
					},
					Repositories: []Repo{
						{
							Name:    "Example Repo",
							Comment: "An example repository for the project",
							URL:     "https://example.com/example-repo",
						},
					},
					Vulnerability: VulnReport{
						ReportsAccepted:    true,
						BugBountyAvailable: false,
					},
				},
			},
			valid: true,
		},
		{
			name: "Valid Security Insights - project and repo barebones",
			si: SecurityInsights{
				Header: Header{
					SchemaVersion: "2.0.0",
					LastUpdated:   "2021-09-01",
					LastReviewed:  "2021-09-01",
					URL:           "https://example.com/security-insights",
				},
				Project: Project{
					Name: "Example Project",
					Administrators: []Contact{
						{
							Name:    "Example Admin",
							Primary: true,
						},
					},
					Repositories: []Repo{
						{
							Name:    "Example Repo",
							Comment: "An example repository for the project",
							URL:     "https://example.com/example-repo",
						},
					},
					Vulnerability: VulnReport{
						ReportsAccepted:    true,
						BugBountyAvailable: false,
					},
				},
				Repository: Repository{
					Status:                        "active",
					URL:                           "https://example.com/example-repo",
					AcceptsChangeRequest:          true,
					AcceptsAutomatedChangeRequest: true,
					CoreTeam: []Contact{
						{
							Name:    "Core Team Member",
							Primary: true,
						},
					},
					License: License{
						URL:        "https://example.com/license",
						Expression: "MIT",
					},
					Security: SecurityInfo{
						Assessments: Assessments{
							Self: Assessment{
								Comment: "Self-assessment of the repository's security",
							},
						},
					},
				},
			},
			valid: true,
		},
		{
			name: "Invalid Security Insights - invalid version in header",
			si: SecurityInsights{
				Header: Header{
					SchemaVersion: "invalid-version",
				},
			},
			valid: false,
		},
		{
			name: "Invalid Security Insights - missing header URL",
			si: SecurityInsights{
				Header: Header{
					SchemaVersion: "2.0.0",
					LastUpdated:   "2021-09-01",
					LastReviewed:  "2021-09-01",
				},
			},
			valid: false,
		},
		{
			name: "Invalid Security Insights - project without administrators",
			si: SecurityInsights{
				Header: Header{
					SchemaVersion: "2.0.0",
					LastUpdated:   "2021-09-01",
					LastReviewed:  "2021-09-01",
					URL:           "https://example.com/security-insights",
				},
				Project: Project{
					Name:           "Example Project",
					Administrators: []Contact{},
					Repositories: []Repo{
						{
							Name:    "Example Repo",
							Comment: "An example repository for the project",
							URL:     "https://example.com/example-repo",
						},
					},
					Vulnerability: VulnReport{
						ReportsAccepted:    true,
						BugBountyAvailable: false,
					},
				},
			},
			valid:      false,
			skip:       true,
			skipReason: "skipped until the new schema version is released that forbids empty administrators",
		},
		{
			name: "Invalid Security Insights - project without repositories",
			si: SecurityInsights{
				Header: Header{
					SchemaVersion: "2.0.0",
					LastUpdated:   "2021-09-01",
					LastReviewed:  "2021-09-01",
					URL:           "https://example.com/security-insights",
				},
				Project: Project{
					Name: "Example Project",
					Administrators: []Contact{
						{
							Name:    "Example Admin",
							Primary: true,
						},
					},
					Repositories: []Repo{},
					Vulnerability: VulnReport{
						ReportsAccepted:    true,
						BugBountyAvailable: false,
					},
				},
			},
			valid:      false,
			skip:       true,
			skipReason: "skipped until the new schema version is released that forbids empty repositories",
		},
		{
			name: "Invalid Security Insights - repository with empty core team",
			si: SecurityInsights{
				Header: Header{
					SchemaVersion: "2.0.0",
					LastUpdated:   "2021-09-01",
					LastReviewed:  "2021-09-01",
					URL:           "https://example.com/security-insights",
				},
				Project: Project{
					Name: "Example Project",
					Administrators: []Contact{
						{
							Name:    "Example Admin",
							Primary: true,
						},
					},
					Repositories: []Repo{
						{
							Name:    "Example Repo",
							Comment: "An example repository for the project",
							URL:     "https://example.com/example-repo",
						},
					},
					Vulnerability: VulnReport{
						ReportsAccepted:    true,
						BugBountyAvailable: false,
					},
				},
				Repository: Repository{
					Status:                        "active",
					URL:                           "https://example.com/example-repo",
					AcceptsChangeRequest:          true,
					AcceptsAutomatedChangeRequest: true,
					CoreTeam:                      []Contact{},
					License: License{
						URL:        "https://example.com/license",
						Expression: "MIT",
					},
					Security: SecurityInfo{
						Assessments: Assessments{
							Self: Assessment{
								Comment: "Self-assessment of the repository's security",
							},
						},
					},
				},
			},
			valid:      false,
			skip:       true,
			skipReason: "skipped until the new schema version is released that forbids empty coreteam",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skipf("Skipping test: %s, reason: %s", tc.name, tc.skipReason)
				return
			}
			err := tc.si.Validate()
			if (err == nil && !tc.valid) || (err != nil && tc.valid) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.valid)
			}
		})
	}
}
