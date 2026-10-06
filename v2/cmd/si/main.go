// Command si works with Security Insights files in public GitHub repositories.
//
//	si fetch <target>...
//
// Each target is one of:
//
//	owner/repo                                             discover the file's location
//	owner/repo/path/to/security-insights.yml               exact path
//	https://github.com/owner/repo/blob/<ref>/<path>        exact path (default branch is read; <ref> is ignored)
//	https://raw.githubusercontent.com/owner/repo/<ref>/<path>
//
// Results are written to stdout as a JSON array in input order; diagnostics go
// to stderr. The exit code is 0 whenever every target was processed, even if
// some could not be found or parsed — that outcome is reported in each result's
// "status". Set GITHUB_TOKEN to raise the GitHub API rate limit.
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "fetch" {
		fmt.Fprintln(os.Stderr, "usage: si fetch <target>...\n\ntarget: owner/repo | owner/repo/<path> | GitHub blob URL | raw.githubusercontent.com URL")
		os.Exit(2)
	}
	results := make([]result, 0, len(os.Args)-2)
	for _, target := range os.Args[2:] {
		r := fetch(target)
		fmt.Fprintf(os.Stderr, "%-9s %s\n", r.Status, target)
		results = append(results, r)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
