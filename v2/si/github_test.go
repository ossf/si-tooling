package si

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeGitHub serves the subset of the contents API that Fetch and Discover use:
// files maps "path" => content; dirs maps "dir" => entry names. Anything else is a 404.
func fakeGitHub(t *testing.T, files map[string]string, dirs map[string][]string) {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/o/r/contents/", func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path[len("/repos/o/r/contents/"):]
		if body, ok := files[p]; ok {
			_ = json.NewEncoder(w).Encode(map[string]string{
				"type": "file", "name": p, "encoding": "base64",
				"content": base64.StdEncoding.EncodeToString([]byte(body)),
			})
			return
		}
		if names, ok := dirs[p]; ok {
			entries := make([]map[string]string, 0, len(names))
			for _, n := range names {
				entries = append(entries, map[string]string{"type": "file", "name": n})
			}
			_ = json.NewEncoder(w).Encode(entries)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Not Found"}`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	t.Setenv("GITHUB_API_URL", srv.URL)
}

func TestFetch(t *testing.T) {
	fakeGitHub(t, map[string]string{"security-insights.yml": "header: {}"}, nil)

	got, err := Fetch("o", "r", "security-insights.yml")
	require.NoError(t, err)
	assert.Equal(t, "header: {}", string(got))

	_, err = Fetch("o", "r", "missing.yml")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestDiscover(t *testing.T) {
	tests := []struct {
		name string
		dirs map[string][]string
		want string
		err  error
	}{
		{
			name: "root wins over .github",
			dirs: map[string][]string{"": {"README.md", "SECURITY-INSIGHTS.yml"}, ".github": {"security-insights.yml"}},
			want: "SECURITY-INSIGHTS.yml",
		},
		{
			name: "falls back to .github",
			dirs: map[string][]string{"": {"README.md"}, ".github": {"workflows", "security-insights.yml"}},
			want: ".github/security-insights.yml",
		},
		{
			name: "no .github directory",
			dirs: map[string][]string{"": {"README.md"}},
			err:  ErrNotFound,
		},
		{
			name: "empty repository",
			dirs: map[string][]string{},
			err:  ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGitHub(t, nil, tt.dirs)
			got, err := Discover("o", "r")
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReadUsesFetch(t *testing.T) {
	fakeGitHub(t, map[string]string{"security-insights.yml": string(minimalV220TestData())}, nil)
	insights, err := Read("o", "r", "security-insights.yml")
	require.NoError(t, err)
	assert.Equal(t, "2.2.0", insights.Header.SchemaVersion.String())
}

func TestGitHubClientEnv(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "x")
	t.Setenv("GITHUB_API_URL", "https://ghe.example.com/api/v3")
	assert.Equal(t, "https://ghe.example.com/api/v3/", githubClient().BaseURL.String())
	t.Setenv("GITHUB_API_URL", "")
	assert.Equal(t, "https://api.github.com/", githubClient().BaseURL.String())
}
