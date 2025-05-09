package si

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		{
			name:          "minimal - v2.1.0",
			contents:      minimalV210TestData(),
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

func minimalV210TestData() []byte {
	data, err := os.ReadFile("test_data/minimal-v2.1.0.yml")
	if err != nil {
		panic(fmt.Sprintf("failed to read test data: %v", err))
	}
	return data
}

func TestNewURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected URL
	}{
		{
			name:     "valid URL",
			url:      "https://example.com",
			expected: URL("https://example.com"),
		},
		{
			name:     "empty URL",
			url:      "",
			expected: URL(""),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := NewURL(test.url)
			assert.Equal(t, test.expected, actual)
		})
	}
}
func TestNewEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected Email
	}{
		{
			name:     "valid email",
			email:    "foo@example.com",
			expected: Email("foo@example.com"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := NewEmail(test.email)
			assert.Equal(t, test.expected, actual)
		})
	}
}
