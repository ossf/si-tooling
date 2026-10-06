package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/ossf/si-tooling/v2/si"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTarget(t *testing.T) {
	tests := []struct {
		target            string
		owner, repo, path string
		wantErr           bool
	}{
		{target: "ossf/scorecard", owner: "ossf", repo: "scorecard"},
		{target: "ossf/scorecard/.github/security-insights.yml", owner: "ossf", repo: "scorecard", path: ".github/security-insights.yml"},
		{target: "https://github.com/ossf/scorecard/blob/main/.github/security-insights.yml", owner: "ossf", repo: "scorecard", path: ".github/security-insights.yml"},
		{target: "https://raw.githubusercontent.com/ossf/scorecard/main/.github/security-insights.yml", owner: "ossf", repo: "scorecard", path: ".github/security-insights.yml"},
		{target: "https://github.com/ossf/scorecard", wantErr: true},
		{target: "https://github.com/ossf/scorecard/tree/main/.github", wantErr: true},
		{target: "https://raw.githubusercontent.com/ossf/scorecard", wantErr: true},
		{target: "https://gitlab.com/o/r/-/blob/main/x.yml", wantErr: true},
		{target: "scorecard", wantErr: true},
		{target: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			owner, repo, path, err := parseTarget(tt.target)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, [3]string{tt.owner, tt.repo, tt.path}, [3]string{owner, repo, path})
		})
	}
}

func TestSchemaVersion(t *testing.T) {
	assert.Equal(t, "1.0.0", schemaVersion([]byte("header:\n  schema-version: 1.0.0\n  commit-hash: abc\n")))
	assert.Equal(t, "", schemaVersion([]byte("not: yaml: at: all")))
}

func TestFail(t *testing.T) {
	assert.Equal(t, statusNotFound, fail(result{}, si.ErrNotFound).Status)
	assert.Equal(t, statusError, fail(result{}, errors.New("boom")).Status)
}

func TestFetchTransportError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"message":"API rate limit exceeded"}`))
	}))
	t.Cleanup(srv.Close)
	t.Setenv("GITHUB_API_URL", srv.URL)

	for _, target := range []string{"o/r", "o/r/security-insights.yml"} {
		r := fetch(target)
		assert.Equal(t, statusError, r.Status, target)
		assert.Contains(t, r.Error, "rate limit", target)
		assert.Nil(t, r.Insights)
	}

	srv500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv500.Close)
	t.Setenv("GITHUB_API_URL", srv500.URL)
	assert.Equal(t, statusError, fetch("o/r").Status)
}

func TestFetchBadTarget(t *testing.T) {
	r := fetch("nonsense")
	assert.Equal(t, statusInvalid, r.Status)
	assert.Equal(t, "nonsense", r.Target)
	assert.Nil(t, r.Insights)
}

// fakeGitHub serves one repository, o/r, over the contents API: a root listing
// plus the given files.
func fakeGitHub(t *testing.T, files map[string]string) {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/o/r/contents/", func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/repos/o/r/contents/")
		if body, ok := files[p]; ok {
			_ = json.NewEncoder(w).Encode(map[string]string{
				"type": "file", "name": path.Base(p), "encoding": "base64",
				"content": base64.StdEncoding.EncodeToString([]byte(body)),
			})
			return
		}
		if p == "" {
			var entries []map[string]string
			for name := range files {
				if !strings.Contains(name, "/") {
					entries = append(entries, map[string]string{"type": "file", "name": name})
				}
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
	minimal, err := os.ReadFile("../../si/test_data/minimal-v2.2.0.yml")
	require.NoError(t, err)
	fakeGitHub(t, map[string]string{
		"SECURITY-INSIGHTS.yml": string(minimal),
		"v1.yml":                "header:\n  schema-version: 1.0.0\n  commit-hash: abc\n",
	})

	discovered := fetch("o/r")
	assert.Equal(t, statusOK, discovered.Status)
	assert.Equal(t, "SECURITY-INSIGHTS.yml", discovered.Path)
	assert.Equal(t, "2.2.0", discovered.SchemaVersion)
	require.NotNil(t, discovered.Insights)

	exact := fetch("https://github.com/o/r/blob/main/SECURITY-INSIGHTS.yml")
	assert.Equal(t, statusOK, exact.Status)
	assert.Equal(t, "o", exact.Owner)

	invalid := fetch("o/r/v1.yml")
	assert.Equal(t, statusInvalid, invalid.Status)
	assert.Equal(t, "1.0.0", invalid.SchemaVersion)
	assert.Nil(t, invalid.Insights)

	missing := fetch("o/r/nope.yml")
	assert.Equal(t, statusNotFound, missing.Status)
	assert.Equal(t, "nope.yml", missing.Path)

	undiscoverable := fetch("o/other")
	assert.Equal(t, statusNotFound, undiscoverable.Status)
	assert.Equal(t, "", undiscoverable.Path)
}
