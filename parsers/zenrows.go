package parsers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bjesus/pipet/common"
	"github.com/google/shlex"
	scraperapi "github.com/zenrows/zenrows-go-sdk/service/api"
	"os"
	"strconv"
	"time"

	"log"
	"net/url"
)

var BlockTypeZenRows = common.BlockType{
	Name:       "zenrows",
	LinePrefix: "zenrows",
	Handler:    ExecuteZenRowsBlock,
}

var zenrows *scraperapi.Client

func init() {
	var options []scraperapi.Option
	if baseURL := os.Getenv("ZENROWS_BASE_URL"); baseURL != "" {
		options = append(options, scraperapi.WithBaseURL(baseURL))
	}

	if apiKey := os.Getenv("ZENROWS_API_KEY"); apiKey != "" {
		options = append(options, scraperapi.WithAPIKey(apiKey))
	}

	if maxConcurrentRequests := os.Getenv("ZENROWS_MAX_CONCURRENT_REQUESTS"); maxConcurrentRequests != "" {
		if maxConcurrentRequestsInt, err := strconv.Atoi(maxConcurrentRequests); err == nil {
			options = append(options, scraperapi.WithMaxConcurrentRequests(maxConcurrentRequestsInt))
		}
	}

	if maxRetryCount := os.Getenv("ZENROWS_MAX_RETRY_COUNT"); maxRetryCount != "" {
		if maxRetryCountInt, err := strconv.Atoi(maxRetryCount); err == nil {
			options = append(options, scraperapi.WithMaxRetryCount(maxRetryCountInt))
		}
	}

	if retryWaitTime := os.Getenv("ZENROWS_RETRY_WAIT_TIME"); retryWaitTime != "" {
		if retryWaitTimeDuration, err := time.ParseDuration(retryWaitTime); err == nil {
			options = append(options, scraperapi.WithRetryWaitTime(retryWaitTimeDuration))
		}
	}

	if retryMaxWaitTime := os.Getenv("ZENROWS_RETRY_MAX_WAIT_TIME"); retryMaxWaitTime != "" {
		if retryMaxWaitTimeDuration, err := time.ParseDuration(retryMaxWaitTime); err == nil {
			options = append(options, scraperapi.WithRetryMaxWaitTime(retryMaxWaitTimeDuration))
		}
	}

	zenrows = scraperapi.NewClient(options...)
}

func ExecuteZenRowsBlock(block common.Block) (interface{}, string, error) {
	var parts []string
	switch cmd := block.Command.(type) {
	case string:
		parts, _ = shlex.Split(cmd)
	case []string:
		parts = cmd
	default:
	}

	var requestParameters *scraperapi.RequestParameters
	if len(parts) > 2 {
		query, err := url.ParseQuery(parts[2])
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse zenrows parameters: %w", err)
		}

		if requestParameters, err = scraperapi.ParseQueryRequestParameters(query); err != nil {
			fmt.Println(err)
			return nil, "", fmt.Errorf("failed to parse zenrows parameters: %w", err)
		}
	}

	res, err := zenrows.Get(context.Background(), parts[1], requestParameters)
	if err != nil {
		return nil, "", err
	}

	if err = res.Error(); err != nil {
		return nil, "", err
	}

	if isJSON := json.Valid(bytes.TrimSpace(res.Body())); isJSON {
		log.Println("JSON detected")
		return ParseJSONQueries(res.Body(), block.Queries)
	} else {
		log.Println("HTML detected")
		return ParseHTMLQueries(res.Body(), block.Queries, block.NextPage)
	}
}
