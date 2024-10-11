package parsers

import (
	"testing"

	"github.com/bjesus/pipet/common"
	"github.com/stretchr/testify/assert"
)

func TestExecutePlaywrightBlock(t *testing.T) {
	block := common.Block{
		Type:    BlockTypePlaywright,
		Command: "playwright http://example.com",
		Queries: []string{"document.querySelector(\"h1\").innerText.split(\" \")", "document.querySelector(\"h1\") | wc -c"},
	}
	result, _, err := ExecutePlaywrightBlock(block)
	expected := []interface{}{[]interface{}{"Example", "Domain"}, "11\n"}

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected, result) // Mocked JSON output

}
