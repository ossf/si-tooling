package si

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The files in test_data/spec-v2.2.0 are unmodified copies of the examples
// published with Security Insights v2.2.0. Run `make check-spec-examples` to
// confirm they still match.
const specExamplesDir = "test_data/spec-v2.2.0/"

// specParentURL is the project-si-source in the reuse example.
const specParentURL = "https://raw.githubusercontent.com/example/repo/refs/heads/main/security-insights.yml"

func specExample(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(specExamplesDir + name)
	require.NoError(t, err)
	return data
}

func TestLoadSpecExamples(t *testing.T) {
	t.Run("minimum", func(t *testing.T) {
		si, err := Load(specExample(t, "example-minimum.yml"))
		require.NoError(t, err)
		require.NotNil(t, si.Project)
		require.NotNil(t, si.Repository)
		assert.Equal(t, "FooBar", si.Project.Name)
		assert.Len(t, si.Project.Repositories, 1)
		assert.Equal(t, "active", si.Repository.Status)
	})

	t.Run("full", func(t *testing.T) {
		si, err := Load(specExample(t, "example-full.yml"))
		require.NoError(t, err)
		require.NotNil(t, si.Project)
		require.NotNil(t, si.Repository)
		assert.Equal(t, "FooBar", si.Project.Name)
		assert.Len(t, si.Project.Repositories, 2)
		assert.Equal(t, []string{"broken access control", "other"}, si.Project.VulnerabilityReporting.InScope)
		assert.Equal(t, []string{"other"}, si.Project.VulnerabilityReporting.OutOfScope)
		assert.Equal(t, "https://github.com/kubernetes/kubernetes", si.Repository.Url.String())
	})

	t.Run("multi-repository project", func(t *testing.T) {
		si, err := Load(specExample(t, "example-multi-repository-project.yml"))
		require.NoError(t, err)
		require.NotNil(t, si.Project)
		assert.Equal(t, "FooBar", si.Project.Name)
		assert.Len(t, si.Project.Repositories, 2)
	})

	t.Run("multi-repository project reuse", func(t *testing.T) {
		parent := specExample(t, "example-multi-repository-project.yml")
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			_, _ = w.Write(parent)
		}))
		defer server.Close()

		child := string(specExample(t, "example-multi-repository-project-reuse.yml"))
		require.Contains(t, child, specParentURL)
		child = strings.Replace(child, specParentURL, server.URL, 1)

		si, err := Load([]byte(child))
		require.NoError(t, err)
		require.NotNil(t, si.Project, "project should be inherited from project-si-source")
		require.NotNil(t, si.Repository)
		assert.Equal(t, "FooBar", si.Project.Name)
		assert.Equal(t, "https://example.com/foobar/bar", si.Repository.Url.String())
	})
}
