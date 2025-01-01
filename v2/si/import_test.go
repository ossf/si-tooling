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
