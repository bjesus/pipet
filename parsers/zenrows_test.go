package parsers

import (
	"github.com/bjesus/pipet/common"
	"github.com/stretchr/testify/assert"
	scraperapi "github.com/zenrows/zenrows-go-sdk/service/api"
	"testing"
)

func TestExecuteZenRowsBlock(t *testing.T) {
	t.Setenv("ZENROWS_API_KEY", "")
	zenrows = scraperapi.NewClient()

	block := common.Block{
		Type:    BlockTypeZenRows,
		Command: "zenrows http://example.com",
		Queries: []string{"document.querySelector(\"h1\").innerText.split(\" \")", "document.querySelector(\"h1\") | wc -c"},
	}
	_, _, err := ExecuteZenRowsBlock(block)
	assert.ErrorIs(t, err, scraperapi.NotConfiguredError{})

	zenrows = scraperapi.NewClient(scraperapi.WithAPIKey("test"))
	_, _, err = ExecuteZenRowsBlock(block)
	assert.Error(t, err)
}
