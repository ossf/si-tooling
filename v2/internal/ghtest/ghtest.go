// Package ghtest serves the subset of the GitHub contents API that package si
// uses, so tests in more than one package can share one fake.
package ghtest

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
)

// Server is a running fake GitHub for the repository o/r.
type Server struct {
	// Authorization is the Authorization header of the most recent request.
	Authorization string
}

// Serve starts a fake GitHub for repository o/r and points GITHUB_API_URL at
// it for the rest of the test. files maps a path to its content; dirs maps a
// directory to its entry names. Anything else is a 404.
func Serve(t *testing.T, files map[string]string, dirs map[string][]string) *Server {
	t.Helper()
	s := &Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/o/r/contents/", func(w http.ResponseWriter, r *http.Request) {
		s.Authorization = r.Header.Get("Authorization")
		p := strings.TrimPrefix(r.URL.Path, "/repos/o/r/contents/")
		if body, ok := files[p]; ok {
			_ = json.NewEncoder(w).Encode(map[string]string{
				"type": "file", "name": path.Base(p), "encoding": "base64",
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
	return s
}
