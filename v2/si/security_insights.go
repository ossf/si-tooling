package si

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

// SecurityInsightsFilename is the expected name of the YAML file containing the insights data. See https://github.com/ossf/security-insights-spec?tab=readme-ov-file#usage for more details.
//
// Deprecated: the spec accepts both .yml and .yaml. Use SecurityInsightsFilenames.
const SecurityInsightsFilename = "security-insights.yml"

// SecurityInsightsFilenames returns every name a SecurityInsights YAML file may
// use. The historical .yml name comes first, so callers that stop at the first
// match keep their current behavior. A project should keep only one of them.
func SecurityInsightsFilenames() []string {
	return []string{SecurityInsightsFilename, "security-insights.yaml"}
}

// ErrParentUnavailable is returned by Load when the parent file named in
// header.project-si-source could not be fetched for a reason that says nothing
// about either file: a network failure, or any HTTP status other than 200 and
// 404. A 404 means the reference itself is broken and is reported as an
// ordinary error.
var ErrParentUnavailable = errors.New("parent security insights unavailable")

// httpClient is shared by every request the package makes.
// ponytail: fixed 30s timeout; export a setter if a caller needs to tune it.
var httpClient = &http.Client{Timeout: 30 * time.Second}

func fetchParentSecurityInsights(parentUrl string) ([]byte, error) {
	response, err := httpClient.Get(parentUrl)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParentUnavailable, err)
	}
	defer func() { _ = response.Body.Close() }()
	switch response.StatusCode {
	case http.StatusOK:
		return io.ReadAll(response.Body)
	case http.StatusNotFound:
		return nil, fmt.Errorf("unexpected response: %s", response.Status)
	default:
		return nil, fmt.Errorf("%w: %s", ErrParentUnavailable, response.Status)
	}
}

// Read reads a SecurityInsights YAML file from a public GitHub repository
// and returns an error if the file cannot be found or unmarshalled or returns
// a SecurityInsights resulting from the unmarshalling.
func Read(owner, repo, path string) (si SecurityInsights, err error) {
	response, err := Fetch(context.Background(), owner, repo, path)
	if err != nil {
		err = fmt.Errorf("error reading target SI: %w", err)
		return
	}
	insights, err := Load(response)
	if err != nil {
		return si, err
	}
	return *insights, nil
}

// Load loads a SecurityInsights struct from a byte slice. If the byte slice is not valid YAML, it will return an error. If the SecurityInsights data provided in contents refers to a schema version that is not supported, it will return an error. If the SecurityInsights data provided in contents is valid, it will return a pointer to the SecurityInsights struct. If the SecurityInsights data provided in contents is valid and refers to a parent SecurityInsights data source in Header.ProjectSISource, that data source will be loaded and the Project field of the returned SecurityInsights struct will be overridden with the Project field of the loaded data source.
func Load(contents []byte) (si *SecurityInsights, err error) {
	insights := &SecurityInsights{}
	err = yaml.UnmarshalWithOptions(contents, insights, yaml.Strict())
	if err != nil {
		err = fmt.Errorf("error unmarshalling SI: %s", err.Error())
		return nil, err
	}
	if (Header{}) == insights.Header {
		err = fmt.Errorf("data provided is not a valid SecurityInsights")
		return nil, err
	}

	err = insights.Header.SchemaVersion.checkVersion()
	if err != nil {
		return nil, err
	}
	mergeVulnerabilityReportingPolicy(insights)
	if insights.Header.ProjectSISource != nil {
		var raw []byte
		raw, err = fetchParentSecurityInsights(insights.Header.ProjectSISource.String())
		if err != nil {
			err = fmt.Errorf("error reading parent SI: %w", err)
			return
		}
		parent := &SecurityInsights{}
		err = yaml.UnmarshalWithOptions(raw, parent, yaml.Strict())
		if err != nil {
			err = fmt.Errorf("error unmarshalling parent SI: %s", err.Error())
			return
		}
		mergeVulnerabilityReportingPolicy(parent)
		insights.Project = parent.Project
	}
	return insights, nil
}

// mergeVulnerabilityReportingPolicy copies SecurityPolicy into Policy when Policy is nil so callers can rely on Policy for both v2.2 (policy) and older (security-policy) YAML.
func mergeVulnerabilityReportingPolicy(si *SecurityInsights) {
	if si == nil || si.Project == nil {
		return
	}
	vr := &si.Project.VulnerabilityReporting
	if vr.Policy == nil && vr.SecurityPolicy != nil {
		vr.Policy = vr.SecurityPolicy
	}
}

func (u URL) String() string {
	return string(u)
}

func NewURL(url string) URL {
	return URL(url)
}

func (e Email) String() string {
	return string(e)
}

func NewEmail(email string) Email {
	return Email(email)
}

func (d Date) String() string {
	return string(d)
}

func (sv SchemaVersion) String() string {
	return string(sv)
}

func NewSchemaVersion(version string) SchemaVersion {
	return SchemaVersion(version)
}

func (sv SchemaVersion) parseVersion() (major int, minor int, patch int) {
	splitVersion := strings.Split(sv.String(), ".")
	if len(splitVersion) == 3 {
		major, _ = strconv.Atoi(splitVersion[0])
		minor, _ = strconv.Atoi(splitVersion[1])
		patch, _ = strconv.Atoi(splitVersion[2])
		return
	}
	if len(splitVersion) == 2 {
		major, _ = strconv.Atoi(splitVersion[0])
		minor, _ = strconv.Atoi(splitVersion[1])
		return
	}
	if len(splitVersion) == 1 {
		major, _ = strconv.Atoi(splitVersion[0])
		return
	}
	return
}

func (sv SchemaVersion) checkVersion() error {
	// This is a placeholder to determine behavior for different schema versions
	// but currently only v2.0.0 is supported
	major, _, _ := sv.parseVersion()
	if major != 2 {
		return fmt.Errorf("unsupported schema version specified by target: %s", sv)
	}
	return nil
}
