package si

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v71/github"
)

// discoveryDirs and discoveryNames are the locations a repository may keep its
// Security Insights file, in priority order. The spec names security-insights.yml
// or .yaml at the root or under .github/; the other spellings are legacy variants seen in the wild.
var (
	discoveryDirs  = []string{"", ".github"}
	discoveryNames = append(SecurityInsightsFilenames(), "SECURITY-INSIGHTS.yml", "SECURITY_INSIGHTS.yml", "security_insights.yml")
)

// ErrNotFound is returned by Discover and Fetch when no Security Insights file exists at the requested location.
var ErrNotFound = errors.New("security insights file not found")

// githubClient builds the client for one GitHub operation. It authenticates
// with GITHUB_TOKEN when set (unauthenticated requests are limited to 60/hour)
// and talks to GITHUB_API_URL when set, which GitHub Actions exports and which
// points tests at a local server. A GITHUB_API_URL that is not an absolute URL
// is an error rather than a silent fallback to github.com, so a token meant
// for one host is never sent to another.
func githubClient() (*github.Client, error) {
	client := github.NewClient(httpClient)
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		client = client.WithAuthToken(token)
	}
	if raw := os.Getenv("GITHUB_API_URL"); raw != "" {
		base, err := url.Parse(raw)
		if err != nil || base.Host == "" {
			return nil, fmt.Errorf("GITHUB_API_URL %q is not an absolute URL", raw)
		}
		base.Path = strings.TrimSuffix(base.Path, "/") + "/"
		client.BaseURL = base
	}
	return client, nil
}

func isNotFound(err error) bool {
	var ghErr *github.ErrorResponse
	return errors.As(err, &ghErr) && ghErr.Response != nil && ghErr.Response.StatusCode == http.StatusNotFound
}

// Fetch returns the raw contents of a file in a public GitHub repository's
// default branch. It returns ErrNotFound when the path does not exist. The
// contents API does not return files over 1 MB; those yield an error.
func Fetch(ctx context.Context, owner, repo, filePath string) ([]byte, error) {
	client, err := githubClient()
	if err != nil {
		return nil, err
	}
	content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, nil)
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("%w: %s/%s/%s", ErrNotFound, owner, repo, filePath)
		}
		return nil, err
	}
	if content == nil {
		return nil, fmt.Errorf("%s/%s/%s is a directory, not a file", owner, repo, filePath)
	}
	s, err := content.GetContent()
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

// Discover returns the path of the repository's Security Insights file. It
// lists the root and .github directories (at most two API calls) rather than
// probing every candidate path, and returns ErrNotFound when none exist.
func Discover(ctx context.Context, owner, repo string) (string, error) {
	client, err := githubClient()
	if err != nil {
		return "", err
	}
	for _, dir := range discoveryDirs {
		names, err := listDir(ctx, client, owner, repo, dir)
		if err != nil {
			return "", err
		}
		for _, name := range discoveryNames {
			if names[name] {
				return path.Join(dir, name), nil
			}
		}
	}
	return "", fmt.Errorf("%w: %s/%s", ErrNotFound, owner, repo)
}

// listDir returns the set of file names in a repository directory; a missing
// directory yields an empty set rather than an error.
func listDir(ctx context.Context, client *github.Client, owner, repo, dir string) (map[string]bool, error) {
	_, entries, _, err := client.Repositories.GetContents(ctx, owner, repo, dir, nil)
	if err != nil {
		if isNotFound(err) {
			return map[string]bool{}, nil
		}
		return nil, err
	}
	names := make(map[string]bool, len(entries))
	for _, e := range entries {
		if e.GetType() == "file" {
			names[e.GetName()] = true
		}
	}
	return names, nil
}
