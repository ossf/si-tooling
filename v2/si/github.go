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

// DiscoveryPaths lists, in priority order, the locations a repository may keep
// its Security Insights file. The spec names security-insights.yml at the root
// or under .github/; the other spellings are legacy variants seen in the wild.
var DiscoveryPaths = []string{
	"security-insights.yml",
	"SECURITY-INSIGHTS.yml",
	"SECURITY_INSIGHTS.yml",
	"security_insights.yml",
	".github/security-insights.yml",
	".github/SECURITY-INSIGHTS.yml",
	".github/SECURITY_INSIGHTS.yml",
	".github/security_insights.yml",
}

// ErrNotFound is returned by Discover and Fetch when no Security Insights file exists at the requested location.
var ErrNotFound = errors.New("security insights file not found")

// githubClient builds the client for one GitHub API call. It authenticates
// with GITHUB_TOKEN when set (unauthenticated requests are limited to 60/hour)
// and talks to GITHUB_API_URL when set, which GitHub Actions exports and which
// points tests at a local server.
func githubClient() *github.Client {
	client := github.NewClient(http.DefaultClient)
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		client = client.WithAuthToken(token)
	}
	if base, err := url.Parse(os.Getenv("GITHUB_API_URL")); err == nil && base.Host != "" {
		base.Path = strings.TrimSuffix(base.Path, "/") + "/"
		client.BaseURL = base
	}
	return client
}

func isNotFound(err error) bool {
	var ghErr *github.ErrorResponse
	return errors.As(err, &ghErr) && ghErr.Response != nil && ghErr.Response.StatusCode == http.StatusNotFound
}

// Fetch returns the raw contents of a file in a public GitHub repository's
// default branch. It returns ErrNotFound when the path does not exist.
func Fetch(owner, repo, filePath string) ([]byte, error) {
	content, _, _, err := githubClient().Repositories.GetContents(context.Background(), owner, repo, filePath, nil)
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

// Discover returns the path of the repository's Security Insights file, trying
// each entry of DiscoveryPaths in order. It lists the root and .github
// directories once each (two API calls) rather than probing every path. It
// returns ErrNotFound when none of the candidate paths exist.
func Discover(owner, repo string) (string, error) {
	listed := map[string]map[string]bool{}
	for _, candidate := range DiscoveryPaths {
		dir := path.Dir(candidate)
		if dir == "." {
			dir = ""
		}
		names, seen := listed[dir]
		if !seen {
			var err error
			names, err = listDir(owner, repo, dir)
			if err != nil {
				return "", err
			}
			listed[dir] = names
		}
		if names[path.Base(candidate)] {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("%w: %s/%s", ErrNotFound, owner, repo)
}

// listDir returns the set of entry names in a repository directory; a missing
// directory yields an empty set rather than an error.
func listDir(owner, repo, dir string) (map[string]bool, error) {
	_, entries, _, err := githubClient().Repositories.GetContents(context.Background(), owner, repo, dir, nil)
	if err != nil {
		if isNotFound(err) {
			return map[string]bool{}, nil
		}
		return nil, err
	}
	names := make(map[string]bool, len(entries))
	for _, e := range entries {
		names[e.GetName()] = true
	}
	return names, nil
}
