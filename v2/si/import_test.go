package si

import (
	"fmt"
	"testing"
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

func TestSchemaURL(t *testing.T) {
	testData := []struct {
		version string
		want    string
	}{
		{"2.0.0", "https://github.com/ossf/security-insights-spec/releases/download/v2.0.0/schema.cue"},
		{"1.0.0", "https://github.com/ossf/security-insights-spec/releases/download/v1.0.0/schema.cue"},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("SchemaURL(%s)", tt.version), func(t *testing.T) {
			got := schemaReleaseURL(tt.version)
			if got != tt.want {
				t.Errorf("SchemaURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
