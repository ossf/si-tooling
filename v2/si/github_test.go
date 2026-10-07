package si

import (
	"context"
	"testing"

	"github.com/ossf/si-tooling/v2/internal/ghtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "x")
	srv := ghtest.Serve(t, map[string]string{"security-insights.yml": "header: {}"}, map[string][]string{".github": {"workflows"}})
	ctx := context.Background()

	got, err := Fetch(ctx, "o", "r", "security-insights.yml")
	require.NoError(t, err)
	assert.Equal(t, "header: {}", string(got))
	assert.Equal(t, "Bearer x", srv.Authorization)

	_, err = Fetch(ctx, "o", "r", "missing.yml")
	assert.ErrorIs(t, err, ErrNotFound)

	_, err = Fetch(ctx, "o", "r", ".github")
	require.Error(t, err, "a directory is not a file")
	assert.NotErrorIs(t, err, ErrNotFound)
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
			ghtest.Serve(t, nil, tt.dirs)
			got, err := Discover(context.Background(), "o", "r")
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReadUsesFetch(t *testing.T) {
	ghtest.Serve(t, map[string]string{"security-insights.yml": string(minimalV220TestData())}, nil)
	insights, err := Read("o", "r", "security-insights.yml")
	require.NoError(t, err)
	assert.Equal(t, "2.2.0", insights.Header.SchemaVersion.String())

	_, err = Read("o", "r", "missing.yml")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestGitHubClientEnv(t *testing.T) {
	t.Setenv("GITHUB_API_URL", "https://ghe.example.com/api/v3")
	client, err := githubClient()
	require.NoError(t, err)
	assert.Equal(t, "https://ghe.example.com/api/v3/", client.BaseURL.String())

	t.Setenv("GITHUB_API_URL", "")
	client, err = githubClient()
	require.NoError(t, err)
	assert.Equal(t, "https://api.github.com/", client.BaseURL.String())

	t.Setenv("GITHUB_API_URL", "ghe.example.com/api/v3")
	_, err = githubClient()
	assert.Error(t, err, "a scheme-less URL must not fall back to github.com")
}
