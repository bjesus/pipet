package parsers

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseHTMLQueries(htmlData []byte, queries []string, nextPage string) ([]interface{}, string, error) {

	result := []interface{}{}

	// get new HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htmlData)))
	if err != nil {
		return nil, "", err
	}

	for i, line := range queries {

		indentation := CalculateIndentation(line)
		if len(queries) > i+1 && CalculateIndentation(queries[i+1]) > indentation {
			// if line has children
			elements := doc.Find(line)

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

			elements.Each(func(subi int, subdoc *goquery.Selection) {
				html, _ := goquery.OuterHtml(subdoc)

				if strings.HasPrefix(html, "<tr") || strings.HasPrefix(html, "<td") {
					html = "<table>" + html + "</table>"
				}
				value2, _, _ := ParseHTMLQueries([]byte(html), lines, "")

				if len(value2) == 1 {
					subresult = append(subresult, value2[0])
				} else if len(value2) > 1 {
					subresult = append(subresult, value2)
				}

			})

			// result = append(result, subresult)
			// Only append non-empty results
			if len(subresult) > 0 {
				result = append(result, subresult)
			}
		} else if indentation == 0 {
			parts := strings.Split(line, "|")
			elements := doc.Find(parts[0])
			value := elements.First()

			html := ""

			if len(parts) > 1 {
				html, _ = goquery.OuterHtml(value)
				for _, pipe := range parts[1:] {
					pipedValue, err := ExecutePipe(html, strings.TrimSpace(pipe))
					if err != nil {
						// Handle error if needed
						break
					}
					html = strings.TrimRight(pipedValue, "\n")
				}
			} else {
				html = value.Text()
			}

			result = append(result, html)
		}

	}

	nextPageURL := ""
	if nextPage != "" {
		nextPageURL = doc.Find(nextPage).First().AttrOr("href", "")
	}
	return result, nextPageURL, nil
}
