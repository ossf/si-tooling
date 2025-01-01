package si

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func makeApiCall(endpoint, token string) (bytes []byte, err error) {
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	client := &http.Client{}
	response, err := client.Do(request)
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

func getGitHubSourceFile(endpoint string) (response FileAPIResponse, err error) {
	responseData, err := makeApiCall("https://api.github.com/"+endpoint, "")
	if err != nil {
		return
	}
	err = json.Unmarshal(responseData, &response)
	return
}

func Read(owner, repo, path string) (si SecurityInsights, err error) {
	var builder SIBuilder
	// Get Target SI
	response, err := getGitHubSourceFile(fmt.Sprintf("repos/%s/%s/contents/%s", owner, repo, path))
	if err != nil {
		err = fmt.Errorf("error reading target SI: %s", err.Error())
		return
	}

	err = yaml.Unmarshal(response.ByteContent, &builder.TargetSI)
	if err != nil {
		err = fmt.Errorf("error unmarshalling target SI: %s", err.Error())
		return
	}

	// check for parent SI, read if exists
	if builder.TargetSI.Header.ProjectSISource != "" {
		response, err = getGitHubSourceFile(builder.TargetSI.Header.ProjectSISource)
		if err != nil {
			err = fmt.Errorf("error reading parent SI: %s", err.Error())
			return
		}
		err = yaml.Unmarshal(response.ByteContent, &builder.ParentSI)
		if err != nil {
			err = fmt.Errorf("error unmarshalling parent SI: %s", err.Error())
			return
		}
	}

	// Override target SI project data with contents of parent SI project data
	builder.TargetSI.Project = builder.ParentSI.Project

	return builder.TargetSI, nil
}
