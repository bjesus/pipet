package utils

import (
	"os"
	"testing"

	"github.com/bjesus/pipet/common"
	"github.com/stretchr/testify/assert" // Or use default `testing` if no external dependency needed
)

// Test GetSeparator
func TestGetSeparator(t *testing.T) {
	app := &common.PipetApp{
		Separator: []string{"-", ":"},
	}
	assert.Equal(t, "-", GetSeparator(app, 0))
	assert.Equal(t, ", ", GetSeparator(app, 5)) // Test default case
}

// Test FileExists
func TestFileExists(t *testing.T) {
	f, _ := os.CreateTemp("", "example")
	defer os.Remove(f.Name()) // Cleanup
	assert.True(t, FileExists(f.Name()))
	assert.False(t, FileExists("nonexistent.file"))
}

func TestRemoveUnnecessaryNesting(t *testing.T) {
	input := [][][]interface{}{{{"foo", "bar"}}}
	expected := []interface{}{"foo", "bar"}
	result := RemoveUnnecessaryNesting(input)
	assert.Equal(t, result, expected)
}
