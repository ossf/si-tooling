package si

import (
	"cuelang.org/go/cue/cuecontext"
)

// SecurityInsightsFilename is the expected name of the YAML file containing the insights data. See https://github.com/ossf/security-insights-spec?tab=readme-ov-file#usage for more details.
const SecurityInsightsFilename = "security-insights.yml"

type SecurityInsights struct {
	Header     Header     `yaml:"header"`
	Project    Project    `yaml:"project" json:"project,omitempty"`
	Repository Repository `yaml:"repository" json:"repository,omitempty"`
}

type Header struct {
	LastReviewed    string `yaml:"last-reviewed"`
	LastUpdated     string `yaml:"last-updated"`
	SchemaVersion   string `yaml:"schema-version"`
	URL             string `yaml:"url"`
	Comment         string `yaml:"comment"`
	ProjectSISource string `yaml:"project-si-source" json:"project-si-source,omitempty"`
}

type Assessment struct {
	Comment  string `yaml:"comment"`
	Name     string `yaml:"name" json:"name,omitempty"`
	Evidence string `yaml:"evidence" json:"evidence,omitempty"`
	Date     string `yaml:"date" json:"date,omitempty"`
}

type Attestation struct {
	Name         string `yaml:"name"`
	Location     string `yaml:"location"`
	PredicateURI string `yaml:"predicate-uri"`
	Comment      string `yaml:"comment" json:"comment,omitempty"`
}

type Contact struct {
	Name        string `yaml:"name"`
	Primary     bool   `yaml:"primary"`
	Affiliation string `yaml:"affiliation" json:"affiliation,omitempty"`
	Email       string `yaml:"email" json:"email,omitempty"`
	Social      string `yaml:"social" json:"social,omitempty"`
}

type License struct {
	URL        string `yaml:"url"`
	Expression string `yaml:"expression"`
}

type Link struct {
	URI     string `yaml:"uri"`
	Comment string `yaml:"comment"`
}

type Project struct {
	Name           string               `yaml:"name"`
	Homepage       string               `yaml:"homepage" json:"homepage,omitempty"`
	Roadmap        string               `yaml:"roadmap" json:"roadmap,omitempty"`
	Funding        string               `yaml:"funding" json:"funding,omitempty"`
	Administrators []Contact            `yaml:"administrators"`
	Repositories   []Repo               `yaml:"repositories"`
	Vulnerability  VulnReport           `yaml:"vulnerability-reporting"`
	Documentation  ProjectDocumentation `yaml:"documentation" json:"documentation,omitempty"`
}

type Repo struct {
	Name    string `yaml:"name"`
	Comment string `yaml:"comment"`
	URL     string `yaml:"url"`
}

type VulnReport struct {
	ReportsAccepted    bool     `yaml:"reports-accepted"`
	BugBountyAvailable bool     `yaml:"bug-bounty-available"`
	BugBountyProgram   string   `yaml:"bug-bounty-program" json:"bug-bounty-program,omitempty"`
	Contact            Contact  `yaml:"contact" json:"contact,omitempty"`
	Comment            string   `yaml:"comment" json:"comment,omitempty"`
	SecurityPolicy     string   `yaml:"security-policy" json:"security-policy,omitempty"`
	PGPKey             string   `yaml:"pgp-key" json:"pgp-key,omitempty"`
	InScope            []string `yaml:"in-scope" json:"in-scope,omitempty"`
	OutOfScope         []string `yaml:"out-of-scope" json:"out-of-scope,omitempty"`
}

type ProjectDocumentation struct {
	DetailedGuide         string `yaml:"detailed-guide" json:"detailed-guide,omitempty"`
	CodeOfConduct         string `yaml:"code-of-conduct" json:"code-of-conduct,omitempty"`
	QuickstartGuide       string `yaml:"quickstart-guide" json:"quickstart-guide,omitempty"`
	ReleaseProcess        string `yaml:"release-process" json:"release-process,omitempty"`
	SignatureVerification string `yaml:"signature-verification" json:"signature-verification,omitempty"`
}

type Repository struct {
	Status                        string                  `yaml:"status"`
	URL                           string                  `yaml:"url"`
	AcceptsChangeRequest          bool                    `yaml:"accepts-change-request"`
	AcceptsAutomatedChangeRequest bool                    `yaml:"accepts-automated-change-request"`
	BugFixesOnly                  bool                    `yaml:"bug-fixes-only" json:"bug-fixes-only,omitempty"`
	NoThirdPartyPackages          bool                    `yaml:"no-third-party-packages" json:"no-third-party-packages,omitempty"`
	CoreTeam                      []Contact               `yaml:"core-team"`
	License                       License                 `yaml:"license"`
	Security                      SecurityInfo            `yaml:"security"`
	Release                       Release                 `yaml:"release" json:"release,omitempty"`
	Documentation                 RepositoryDocumentation `yaml:"documentation" json:"documentation,omitempty"`
}

type RepositoryDocumentation struct {
	Contributing         string `yaml:"contributing-guide"`
	DependencyManagement string `yaml:"dependency-management-policy"`
	Governance           string `yaml:"governance"`
	ReviewPolicy         string `yaml:"review-policy"`
	SecurityPolicy       string `yaml:"security-policy"`
}

type SecurityInfo struct {
	Assessments Assessments `yaml:"assessments"`
	Champions   []Contact   `yaml:"champions"`
	Tools       []Tool      `yaml:"tools"`
}

type Assessments struct {
	Self       Assessment   `yaml:"self"`
	ThirdParty []Assessment `yaml:"third-party"`
}

type Tool struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"`
	Version     string      `yaml:"version"`
	Comment     string      `yaml:"comment"`
	Rulesets    []string    `yaml:"rulesets"`
	Integration Integration `yaml:"integration"`
	Results     Results     `yaml:"results"`
}

type Integration struct {
	Adhoc   bool `yaml:"adhoc"`
	CI      bool `yaml:"ci"`
	Release bool `yaml:"release"`
}

type Results struct {
	Adhoc   Attestation `yaml:"adhoc"`
	CI      Attestation `yaml:"ci"`
	Release Attestation `yaml:"release"`
}

type Release struct {
	AutomatedPipeline  bool          `yaml:"automated-pipeline"`
	DistributionPoints []Link        `yaml:"distribution-points"`
	Changelog          string        `yaml:"changelog"`
	License            License       `yaml:"license"`
	Attestations       []Attestation `yaml:"attestations"`
}

type SIHeader struct {
	SchemaVersion string `yaml:"schema-version"`
	ChangeLogURL  string `yaml:"changelog"`
	LicenseURL    string `yaml:"license"`
}

// Validate checks if the SecurityInsights object conforms to the Cue schema referenced in the header and returns an error if it does not.
func (si *SecurityInsights) Validate() error {
	ctx := cuecontext.New()
	schemaBytes, err := getSchema(si.Header.SchemaVersion)
	if err != nil {
		return err
	}
	schema := ctx.CompileBytes(schemaBytes)
	if err := schema.Unify(ctx.Encode(si)).Validate(); err != nil {
		return err
	}
	return nil
}
