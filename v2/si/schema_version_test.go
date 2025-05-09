package si

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaVersionCheckVersion(t *testing.T) {
	tests := []struct {
		version   SchemaVersion
		expectErr bool
	}{
		{"2.0.0", false},
		{"2.0", false},
		{"2", false},
		{"1.0.0", true},
		{"3.0.0", true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("version %s", test.version), func(t *testing.T) {
			err := test.version.checkVersion()
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSchemaVersionParseVersion(t *testing.T) {
	tests := []struct {
		version  SchemaVersion
		expected [3]int
	}{
		{"2.0.0", [3]int{2, 0, 0}},
		{"2.0", [3]int{2, 0, 0}},
		{"2", [3]int{2, 0, 0}},
		{"1.0.0", [3]int{1, 0, 0}},
		{"3.0.0", [3]int{3, 0, 0}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("version %s", test.version), func(t *testing.T) {
			actualMajor, actualMinor, actualPatch := test.version.parseVersion()
			assert.Equal(t, test.expected[0], actualMajor)
			assert.Equal(t, test.expected[1], actualMinor)
			assert.Equal(t, test.expected[2], actualPatch)
		})
	}
}
