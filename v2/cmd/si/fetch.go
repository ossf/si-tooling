package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/ossf/si-tooling/v2/si"
)

const (
	statusOK       = "ok"        // file fetched and parsed
	statusNotFound = "not-found" // no file at the path, or discovery found nothing
	statusInvalid  = "invalid"   // file fetched but si.Load rejected it, or the target string is malformed
	statusError    = "error"     // transport failure (rate limit, 5xx, network, bad token): the file's state is unknown
)

// result is one element of the JSON array `si fetch` prints.
type result struct {
	Target        string               `json:"target"`
	Owner         string               `json:"owner,omitempty"`
	Repo          string               `json:"repo,omitempty"`
	Path          string               `json:"path,omitempty"`
	Status        string               `json:"status"`
	Error         string               `json:"error,omitempty"`
	SchemaVersion string               `json:"schema-version,omitempty"`
	Insights      *si.SecurityInsights `json:"insights,omitempty"`
}

// parseTarget splits a target into owner, repo and path. An empty path means
// the caller should discover it.
func parseTarget(target string) (owner, repo, path string, err error) {
	var parts []string
	switch {
	case strings.HasPrefix(target, "https://github.com/"):
		// https://github.com/owner/repo/blob/<ref>/<path>
		u, perr := url.Parse(target)
		if perr != nil {
			return "", "", "", perr
		}
		parts = strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) < 5 || parts[2] != "blob" {
			return "", "", "", fmt.Errorf("expected https://github.com/owner/repo/blob/<ref>/<path>, got %s", target)
		}
		parts = append(parts[:2], parts[4:]...)
	case strings.HasPrefix(target, "https://raw.githubusercontent.com/"):
		// https://raw.githubusercontent.com/owner/repo/<ref>/<path>
		u, perr := url.Parse(target)
		if perr != nil {
			return "", "", "", perr
		}
		parts = strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) < 4 {
			return "", "", "", fmt.Errorf("expected https://raw.githubusercontent.com/owner/repo/<ref>/<path>, got %s", target)
		}
		parts = append(parts[:2], parts[3:]...)
	case strings.Contains(target, "://"):
		return "", "", "", fmt.Errorf("unsupported URL %s: only github.com blob and raw.githubusercontent.com URLs are accepted", target)
	default:
		parts = strings.Split(strings.Trim(target, "/"), "/")
	}
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", "", fmt.Errorf("expected owner/repo[/path], got %s", target)
	}
	return parts[0], parts[1], strings.Join(parts[2:], "/"), nil
}

func fetch(target string) result {
	r := result{Target: target}
	owner, repo, path, err := parseTarget(target)
	if err != nil {
		r.Status, r.Error = statusInvalid, err.Error()
		return r
	}
	r.Owner, r.Repo, r.Path = owner, repo, path
	if r.Path == "" {
		if r.Path, err = si.Discover(owner, repo); err != nil {
			return fail(r, err)
		}
	}
	raw, err := si.Fetch(owner, repo, r.Path)
	if err != nil {
		return fail(r, err)
	}
	r.SchemaVersion = schemaVersion(raw)
	insights, err := si.Load(raw)
	if err != nil {
		r.Status, r.Error = statusInvalid, err.Error()
		return r
	}
	r.Status, r.Insights = statusOK, insights
	return r
}

// fail records a Discover/Fetch failure. Only a 404 says anything about the
// file; every other failure is the transport's, and must not be mistaken for a
// bad or missing file by consumers that act on not-found/invalid.
func fail(r result, err error) result {
	r.Error = err.Error()
	if errors.Is(err, si.ErrNotFound) {
		r.Status = statusNotFound
	} else {
		r.Status = statusError
	}
	return r
}

// schemaVersion leniently reads header.schema-version so a file si.Load
// rejects (a v1 file, an unknown field) can still report which schema it claims.
func schemaVersion(raw []byte) string {
	var doc struct {
		Header struct {
			SchemaVersion string `yaml:"schema-version"`
		} `yaml:"header"`
	}
	_ = yaml.Unmarshal(raw, &doc)
	return doc.Header.SchemaVersion
}
