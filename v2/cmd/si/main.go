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
// to stderr. A file that is missing or does not parse is reported in-band in
// its result's "status" and does not affect the exit code. Exit codes:
//
//	0  every target was processed (statuses ok, not-found or invalid)
//	1  the JSON could not be written
//	2  usage error, or at least one target has status "error" (a transport
//	   failure, so the snapshot is incomplete and should not be trusted)
//
// Set GITHUB_TOKEN to raise the GitHub API rate limit.
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
	degraded := false
	for _, target := range os.Args[2:] {
		r := fetch(target)
		fmt.Fprintf(os.Stderr, "%-9s %s\n", r.Status, target)
		degraded = degraded || r.Status == statusError
		results = append(results, r)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if degraded {
		fmt.Fprintln(os.Stderr, "si fetch: one or more targets could not be fetched; the output is incomplete")
		os.Exit(2)
	}
}
