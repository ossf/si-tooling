package ghtest

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServe(t *testing.T) {
	srv := Serve(t, map[string]string{"a.yml": "hi"}, map[string][]string{"": {"a.yml"}})
	base := os.Getenv("GITHUB_API_URL") + "/repos/o/r/contents/"

	get := func(p string) (int, string) {
		req, err := http.NewRequest(http.MethodGet, base+p, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer x")
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		return resp.StatusCode, string(body)
	}

	code, body := get("a.yml")
	assert.Equal(t, http.StatusOK, code)
	assert.Contains(t, body, `"encoding":"base64"`)
	assert.Equal(t, "Bearer x", srv.Authorization)

	code, body = get("")
	assert.Equal(t, http.StatusOK, code)
	assert.Contains(t, body, `"name":"a.yml"`)

	code, _ = get("nope")
	assert.Equal(t, http.StatusNotFound, code)
}
