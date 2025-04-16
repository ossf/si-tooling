package si

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type FileAPIResponse struct {
	ByteContent []byte `json:"content"`
	SHA         string `json:"sha"`
}

type SIBuilder struct {
	TargetSI SecurityInsights
	ParentSI SecurityInsights
}

func getSecurityInsightFile(endpoint string) (bytes []byte, err error) {
	response, err := http.Get(endpoint)
	if err != nil {
		err = fmt.Errorf("error making http call: %s", err.Error())
		return
	}
	if response.StatusCode != 200 {
		err = fmt.Errorf("unexpected response: %s", response.Status)
		return
	}
	return io.ReadAll(response.Body)
}

func getGitHubFile(endpoint string) (response FileAPIResponse, err error) {
	responseData, err := getSecurityInsightFile("https://api.github.com/" + endpoint)
	if err != nil {
		return
	}
	err = json.Unmarshal(responseData, &response)
	return
}

func parseVersion(version string) (major int, minor int, patch int) {
	splitVersion := strings.Split(version, ".")
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

func checkVersion(version string) error {
	// This is a placeholder to determine behavior for different schema versions
	// but currently only v2.0.0 is supported
	major, minor, patch := parseVersion(version)
	if major != 2 || minor+patch != 0 {
		return fmt.Errorf("unsupported schema version specified by target: %s", version)
	}
	return nil
}

func Read(owner, repo, path string) (si SecurityInsights, err error) {
	var builder SIBuilder
	// Get Target SI
	response, err := getGitHubFile(fmt.Sprintf("repos/%s/%s/contents/%s", owner, repo, path))
	if err != nil {
		err = fmt.Errorf("error reading target SI: %s", err.Error())
		return
	}

	err = yaml.Unmarshal(response.ByteContent, &builder.TargetSI)
	if err != nil {
		err = fmt.Errorf("error unmarshalling target SI: %s", err.Error())
		return
	}

	err = checkVersion(builder.TargetSI.Header.SchemaVersion)
	if err != nil {
		return
	}

	// check for parent SI, read if exists
	if builder.TargetSI.Header.ProjectSISource != "" {
		var raw []byte
		raw, err = getSecurityInsightFile(builder.TargetSI.Header.ProjectSISource)
		if err != nil {
			err = fmt.Errorf("error reading parent SI: %s", err.Error())
			return
		}
		err = yaml.Unmarshal(raw, &builder.ParentSI)
		if err != nil {
			err = fmt.Errorf("error unmarshalling parent SI: %s", err.Error())
			return
		}
	}

	// Override target SI project data with contents of parent SI project data
	builder.TargetSI.Project = builder.ParentSI.Project

	return builder.TargetSI, nil
}
