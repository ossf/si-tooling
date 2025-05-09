package si

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/google/go-github/v71/github"
)

// SecurityInsightsFilename is the expected name of the YAML file containing the insights data. See https://github.com/ossf/security-insights-spec?tab=readme-ov-file#usage for more details.
const SecurityInsightsFilename = "security-insights.yml"

func fetchParentSecurityInsights(parentUrl string) (bytes []byte, err error) {
	request, err := http.NewRequest("GET", parentUrl, nil)
	if err != nil {
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("error making http call: %s", err.Error())
		return
	}
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected response: %s", response.Status)
		return
	}
	return io.ReadAll(response.Body)
}

func getGitHubSourceFile(owner, repo, path string) ([]byte, error) {
	client := github.NewClient(http.DefaultClient)
	content, _, _, err := client.Repositories.GetContents(context.Background(), owner, repo, path, nil)
	if err != nil {
		return nil, err
	}
	s, err := content.GetContent()
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

// Read reads a SecurityInsights YAML file from a public GitHub repository
// and returns an error if the file cannot be found or unmarshalled or returns
// a SecurityInsights resulting from the unmarshalling.
func Read(owner, repo, path string) (si SecurityInsights, err error) {
	response, err := getGitHubSourceFile(owner, repo, path)
	if err != nil {
		err = fmt.Errorf("error reading target SI: %s", err.Error())
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
	if insights.Header.ProjectSISource != nil {
		var raw []byte
		raw, err = fetchParentSecurityInsights(insights.Header.ProjectSISource.String())
		if err != nil {
			err = fmt.Errorf("error reading parent SI: %s", err.Error())
			return
		}
		parent := &SecurityInsights{}
		err = yaml.UnmarshalWithOptions(raw, parent, yaml.Strict())
		if err != nil {
			err = fmt.Errorf("error unmarshalling parent SI: %s", err.Error())
			return
		}
		insights.Project = parent.Project
	}
	return insights, nil
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
