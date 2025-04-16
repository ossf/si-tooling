package si

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validYAML = `header:
  schema-version: "2.0.0"
  last-updated: "2024-03-21"
  project-si-source: ""
project:
  name: "test-project"`

	validYAMLWithParent = `header:
  schema-version: "2.0.0"
  last-updated: "2024-03-21"
  project-si-source: "http://localhost/parent-si.json"
project:
  name: "test-project"`

	invalidSchemaYAML = `header:
  schema-version: "3.0.0"
  last-updated: "2024-03-21"
  project-si-source: ""
project:
  name: "test-project"`

	invalidYAML = `invalid:
yaml:content`

	parentSIJSON = `{
  "header": {
    "schema-version": "2.0.0",
    "last-updated": "2024-03-21"
  },
  "project": {
    "name": "parent-project",
    "homepage": "https://example.com"
  }
}`
)

// createGitHubResponse creates a GitHub API response with base64 encoded content
func createGitHubResponse(content string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(content))
	response := struct {
		Content string `json:"content"`
		SHA     string `json:"sha"`
	}{
		Content: encoded,
		SHA:     "abc123",
	}

	jsonBytes, _ := json.Marshal(response)
	return string(jsonBytes)
}

func TestRead(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    bool
		errMessage string
		responses  map[string]struct {
			status     string
			statusCode int
			body       string
		}
	}{
		{
			name:    "successful read without parent SI",
			wantErr: false,
			responses: map[string]struct {
				status     string
				statusCode int
				body       string
			}{
				"https://api.github.com/repos/test-owner/test-repo/contents/security-insights.yml": {
					status:     "200 OK",
					statusCode: http.StatusOK,
					body:       createGitHubResponse(validYAML),
				},
			},
		},
		{
			name:    "successful read with parent SI",
			wantErr: false,
			responses: map[string]struct {
				status     string
				statusCode int
				body       string
			}{
				"https://api.github.com/repos/test-owner/test-repo/contents/security-insights.yml": {
					status:     "200 OK",
					statusCode: http.StatusOK,
					body:       createGitHubResponse(validYAMLWithParent),
				},
				"http://localhost/parent-si.json": {
					status:     "200 OK",
					statusCode: http.StatusOK,
					body:       parentSIJSON,
				},
			},
		},
		{
			name:       "invalid schema version",
			wantErr:    true,
			errMessage: "unsupported schema version specified by target: 3.0.0",
			responses: map[string]struct {
				status     string
				statusCode int
				body       string
			}{
				"https://api.github.com/repos/test-owner/test-repo/contents/security-insights.yml": {
					status:     "200 OK",
					statusCode: http.StatusOK,
					body:       createGitHubResponse(invalidSchemaYAML),
				},
			},
		},
		{
			name:       "invalid GitHub API response",
			wantErr:    true,
			errMessage: "error reading target SI: unexpected response: 404 Not Found",
			responses: map[string]struct {
				status     string
				statusCode int
				body       string
			}{
				"https://api.github.com/repos/test-owner/test-repo/contents/security-insights.yml": {
					status:     "404 Not Found",
					statusCode: http.StatusNotFound,
					body:       `{"message": "Not Found"}`,
				},
			},
		},
		{
			name:       "invalid YAML content",
			wantErr:    true,
			errMessage: "error unmarshalling target SI:",
			responses: map[string]struct {
				status     string
				statusCode int
				body       string
			}{
				"https://api.github.com/repos/test-owner/test-repo/contents/security-insights.yml": {
					status:     "200 OK",
					statusCode: http.StatusOK,
					body:       createGitHubResponse(invalidYAML),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a custom transport that returns our mock responses
			transport := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
				response, exists := tt.responses[req.URL.String()]
				if !exists {
					return &http.Response{
						Status:     "404 Not Found",
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(strings.NewReader(`{"message": "Not Found"}`)),
					}, nil
				}

				return &http.Response{
					Status:     response.status,
					StatusCode: response.statusCode,
					Body:       io.NopCloser(strings.NewReader(response.body)),
				}, nil
			})

			// Override the default client with our custom transport
			originalClient := http.DefaultClient
			http.DefaultClient = &http.Client{Transport: transport}
			defer func() { http.DefaultClient = originalClient }()

			si, err := Read("test-owner", "test-repo", "security-insights.yml")
			t.Logf("si: %+v", si)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, si)
				assert.Equal(t, "2.0.0", si.Header.SchemaVersion)
				if tt.name == "successful read with parent SI" {
					assert.Equal(t, "parent-project", si.Project.Name)
				}
			}
		})
	}
}

// Helper type for modifying the HTTP client
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
