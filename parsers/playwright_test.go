package parsers

import (
	"fmt"
	"testing"

	"github.com/bjesus/pipet/common"
	"github.com/stretchr/testify/assert"
)

func TestExecutePlaywrightBlock(t *testing.T) {

	block := common.Block{
		Type:    "playwright",
		Command: "playwright http://example.com",
		Queries: []string{"document.querySelector(\"h1\").innerText.split(\" \")", "document.querySelector(\"h1\") | wc -c"},
	}
	result, err := ExecutePlaywrightBlock(block)
	expected := []interface{}{[]interface{}{"Example", "Domain"}, "11\n"}
	fmt.Printf("%v", result)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected, result) // Mocked JSON output

}
