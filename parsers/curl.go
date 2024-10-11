package parsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bjesus/pipet/common"
	"github.com/google/shlex"
	"log"
	"os/exec"
)

var BlockTypeCurl = common.BlockType{
	Name:       "curl",
	LinePrefix: "curl",
	Handler:    ExecuteCurlBlock,
}

func ExecuteCurlBlock(block common.Block) (interface{}, string, error) {
	if !commandExists("curl") {
		return nil, "", fmt.Errorf("curl command not found. Please install curl and try again")
	}

	var parts []string
	switch cmd := block.Command.(type) {
	case string:
		parts, _ = shlex.Split(cmd)
	case []string:
		parts = cmd
	default:
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return nil, "", fmt.Errorf("curl command failed: %w\nOutput: %s", err, string(output))
	}

	isJSON := json.Valid(bytes.TrimSpace(output))

	if isJSON {
		log.Println("JSON detected")
		return ParseJSONQueries(output, block.Queries)
	} else {
		log.Println("HTML detected")
		return ParseHTMLQueries(output, block.Queries, block.NextPage)
	}
}
