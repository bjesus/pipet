package parsers

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock exec.Command for testing ExecuteCurlBlock and ExecutePipe
var execCommand = exec.Command

func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// Test helper for mock exec
func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	cmd := os.Args[3]
	switch cmd {
	case "curl":
		os.Stdout.Write([]byte(`{"success":true}`)) // Mock JSON response
	default:
		os.Stdout.Write([]byte("bash output"))
	}
	os.Exit(0)
}

// Test ExecutePipe
func TestExecutePipe(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }() // Restore exec.Command after test

	input := "pipet"
	command := "wc -c"
	output, err := ExecutePipe(input, command)
	assert.NoError(t, err)
	assert.Equal(t, "5\n", output) // Mocked output
}

// Test CalculateIndentation
func TestCalculateIndentation(t *testing.T) {
	assert.Equal(t, 4, CalculateIndentation("    indented"))
	assert.Equal(t, 0, CalculateIndentation("not indented"))
}

// Test commandExists
func TestCommandExists(t *testing.T) {
	assert.True(t, commandExists("bash"))  // Assuming bash exists on your system
	assert.False(t, commandExists("fake")) // Fake command
}
