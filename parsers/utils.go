package parsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/bjesus/pipet/common"
	"github.com/google/shlex"
)

func ExecuteCurlBlock(block common.Block) (interface{}, error) {
	if !commandExists("curl") {
		return nil, fmt.Errorf("curl command not found. Please install curl and try again")
	}

	// Split the command into curl and its arguments
	parts, err := shlex.Split(block.Command)
	log.Println("Execute curl:", block.Command)
	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("curl command failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Println(string(output))
	isJSON := json.Valid(bytes.TrimSpace(output))
	log.Println(isJSON)
	if isJSON {
		log.Println("it's json")
		return ParseJSONQueries(output, block.Queries)
	} else {
		log.Println("it's HTM")
		return ParseHTMLQueries(output, block.Queries)
	}
}

func ExecutePipe(input string, command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdin = strings.NewReader(input)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func CalculateIndentation(s string) int {
	count := 0
	for _, char := range s {
		if char != ' ' {
			break
		}
		count++
	}
	return count
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
