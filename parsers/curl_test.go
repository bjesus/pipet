package parsers

import (
	"fmt"
	"github.com/bjesus/pipet/common"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

// Test ExecuteCurlBlock
func TestExecuteCurlBlock(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }() // Restore exec.Command after test

	block := common.Block{
		Type:    BlockTypeCurl,
		Command: "curl http://example.com",
		Queries: []string{"h1"},
	}
	result, _, err := ExecuteCurlBlock(block)
	fmt.Printf("%v", result)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, []interface{}{"Example Domain"}, result) // Mocked JSON output
}
