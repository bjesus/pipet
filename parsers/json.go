package parsers

import (
	"encoding/json"
	"strings"

	"github.com/tidwall/gjson"
)

func ParseJSONQueries(jsonData []byte, queries []string) (interface{}, error) {

	result := []interface{}{}

	for i, line := range queries {

		indentation := CalculateIndentation(line)
		if len(queries) > i+1 && CalculateIndentation(queries[i+1]) > indentation {
			// if line has children
			elements := gjson.GetBytes(jsonData, line)

			// get new lines
			var lines []string

			for subi := i + 1; subi < len(queries); subi++ {
				if CalculateIndentation(queries[subi]) > indentation {
					lines = append(lines, queries[subi][2:])
				} else {
					break
				}
			}

			subresult := []interface{}{}
			elements.ForEach(func(subi gjson.Result, subdoc gjson.Result) bool {
				html := subdoc.String()

				value2, _ := ParseJSONQueries([]byte(html), lines)
				subresult = append(subresult, value2)
				return true
			})

			result = append(result, subresult)
		} else if indentation == 0 {
			parts := strings.Split(line, "|")
			query := strings.TrimSpace(parts[0])
			value := gjson.GetBytes(jsonData, query)

			html := ""

			if len(parts) > 1 {
				html = value.String()
				for _, pipe := range parts[1:] {
					pipedValue, err := ExecutePipe(html, strings.TrimSpace(pipe))
					if err != nil {
						// Handle error if needed
						break
					}
					html = strings.TrimRight(pipedValue, "\n")
				}
			} else {
				html = value.String()
			}

			if json.Valid([]byte(html)) {
				var parsedJSON interface{}
				json.Unmarshal([]byte(html), &parsedJSON)
				result = append(result, parsedJSON)
			} else {
				result = append(result, html)

			}

		}

	}
	return result, nil
}
