package si

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	testData := []struct {
		owner string
		repo  string
		path  string
	}{
		{"ossf", "security-insights-spec", ".github/security-insights.yml"},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Read(%s, %s, %s)", tt.owner, tt.repo, tt.path), func(t *testing.T) {
			// TODO: Add real test cases
			out, err := Read(tt.owner, tt.repo, tt.path)
			if err != nil {
				t.Errorf("Read() error = %v", err)
				return
			}
			fmt.Print(out)
		})
	}
}

func TestLoad(t *testing.T) {
	type testCase struct {
		name          string
		contents      []byte
		errorExpected bool
		want          *SecurityInsights
	}
	testCases := []testCase{
		{
			name:          "invalid YAML",
			contents:      []byte("invalid YAML"),
			errorExpected: true,
			want:          nil,
		},
		{
			name:          "empty header",
			contents:      []byte("header:\n"),
			errorExpected: true,
			want:          nil,
		},
		{
			name:          "invalid schema version",
			contents:      []byte("header:\n  schemaVersion: invalid"),
			errorExpected: true,
			want:          nil,
		},
		{
			name:          "minimal",
			contents:      minimalTestData(),
			errorExpected: false,
			want:          nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(tt.contents)
			assert.Equal(t, tt.errorExpected, err != nil)
		})
	}
}

func minimalTestData() []byte {
	data, err := os.ReadFile("test_data/minimal.yml")
	if err != nil {
		panic(fmt.Sprintf("failed to read test data: %v", err))
	}
	return data
}
