package app

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bjesus/pipet/common"
	"github.com/bjesus/pipet/parsers"
	"github.com/tidwall/gjson"
)

func ParseSpecFile(e *common.PipetApp, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentBlock *common.Block

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if currentBlock != nil {
				e.Blocks = append(e.Blocks, *currentBlock)
				currentBlock = nil
			}
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}
		if currentBlock == nil {
			if strings.HasPrefix(line, "curl ") {
				currentBlock = &common.Block{Type: "curl", Command: line}
			} else if strings.HasPrefix(line, "playwright ") {
				currentBlock = &common.Block{Type: "playwright", Command: line}
			} else {
				return fmt.Errorf("invalid block start: %s", line)
			}
		} else {
			if strings.HasPrefix(line, "> ") {

				currentBlock.NextPage = strings.TrimPrefix(line, "> ")
			} else {
				currentBlock.Queries = append(currentBlock.Queries, line)
			}
		}
	}

	log.Println("Found block", currentBlock)
	if currentBlock != nil {
		e.Blocks = append(e.Blocks, *currentBlock)
	}

	return scanner.Err()
}

func ExecuteBlocks(e *common.PipetApp) error {
	for _, block := range e.Blocks {
		var data interface{}
		var err error

		for page := 0; page < e.MaxPages; page++ {
			if block.Type == "curl" {
				data, err = parsers.ExecuteCurlBlock(block)
			} else if block.Type == "playwright" {
				data, err = parsers.ExecutePlaywrightBlock(block)
			}

			if err != nil {
				return err
			}

			e.Data = append(e.Data, data)

			if block.NextPage == "" {
				break
			}

			nextURL, err := getNextPageURL(block, data)
			if err != nil {
				return err
			}

			block.Command = strings.Replace(block.Command, block.Command[strings.Index(block.Command, " ")+1:], nextURL, 1)
		}
	}

	return nil
}

func getNextPageURL(block common.Block, data interface{}) (string, error) {
	parts := strings.Split(block.NextPage, "|")
	selector := strings.TrimSpace(parts[0])

	var nextURL string

	if block.Type == "curl" {
		if strings.HasPrefix(selector, ".") {
			// JSON mode
			nextURL = gjson.Get(fmt.Sprintf("%v", data), selector).String()
		} else {
			// HTML mode
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(fmt.Sprintf("%v", data)))
			if err != nil {
				return "", err
			}
			nextURL, _ = doc.Find(selector).Attr("href")
		}
	} else if block.Type == "playwright" {
		// TODO: Implement Playwright next page logic
		return "", fmt.Errorf("playwright next page not implemented")
	}

	if len(parts) > 1 {
		pipedURL, err := parsers.ExecutePipe(nextURL, strings.TrimSpace(parts[1]))
		if err != nil {
			return "", err
		}
		nextURL = strings.TrimSpace(pipedURL)
	}

	return nextURL, nil
}
