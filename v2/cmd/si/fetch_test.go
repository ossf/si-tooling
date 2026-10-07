package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ossf/si-tooling/v2/internal/ghtest"
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
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"message":"API rate limit exceeded"}`))
	}))
	t.Cleanup(srv.Close)
	t.Setenv("GITHUB_API_URL", srv.URL)

	for _, target := range []string{"o/r", "o/r/security-insights.yml"} {
		r := fetch(ctx, target)
		assert.Equal(t, statusError, r.Status, target)
		assert.Contains(t, r.Error, "rate limit", target)
		assert.Nil(t, r.Insights)
	}

	srv500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv500.Close)
	t.Setenv("GITHUB_API_URL", srv500.URL)
	assert.Equal(t, statusError, fetch(ctx, "o/r").Status)
}

func TestFetchParentUnavailable(t *testing.T) {
	parent := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	t.Cleanup(parent.Close)
	ghtest.Serve(t, map[string]string{
		"security-insights.yml": "header:\n  schema-version: 2.0.0\n  project-si-source: " + parent.URL + "\n",
	}, nil)

	r := fetch(context.Background(), "o/r/security-insights.yml")
	assert.Equal(t, statusError, r.Status, "an unreachable parent says nothing about the file")
	assert.Equal(t, "2.0.0", r.SchemaVersion)
	assert.Nil(t, r.Insights)
}

func TestFetchBadTarget(t *testing.T) {
	r := fetch(context.Background(), "nonsense")
	assert.Equal(t, statusInvalid, r.Status)
	assert.Equal(t, "nonsense", r.Target)
	assert.Nil(t, r.Insights)
}

// Discovery itself is tested in package si; here the fake has no root listing,
// so an owner/repo target is not-found.
func TestFetch(t *testing.T) {
	ctx := context.Background()
	minimal, err := os.ReadFile("../../si/test_data/minimal-v2.2.0.yml")
	require.NoError(t, err)
	ghtest.Serve(t, map[string]string{
		"SECURITY-INSIGHTS.yml": string(minimal),
		"v1.yml":                "header:\n  schema-version: 1.0.0\n  commit-hash: abc\n",
	}, nil)

	exact := fetch(ctx, "https://github.com/o/r/blob/main/SECURITY-INSIGHTS.yml")
	assert.Equal(t, statusOK, exact.Status)
	assert.Equal(t, "o", exact.Owner)
	assert.Equal(t, "2.2.0", exact.SchemaVersion)
	require.NotNil(t, exact.Insights)

	invalid := fetch(ctx, "o/r/v1.yml")
	assert.Equal(t, statusInvalid, invalid.Status)
	assert.Equal(t, "1.0.0", invalid.SchemaVersion)
	assert.Nil(t, invalid.Insights)

	missing := fetch(ctx, "o/r/nope.yml")
	assert.Equal(t, statusNotFound, missing.Status)
	assert.Equal(t, "nope.yml", missing.Path)

	undiscoverable := fetch(ctx, "o/r")
	assert.Equal(t, statusNotFound, undiscoverable.Status)
	assert.Equal(t, "", undiscoverable.Path)
}
