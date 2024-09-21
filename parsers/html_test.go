package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHTMLQueries(t *testing.T) {
	htmlData := []byte(`
		<ul><li>one</li><li>two</li><li>three</li></ul>
		<p><ol><li>foo</li><li>bar</li></ol>
		<table>
			<tr><td>penny</td><td>lane</td></tr>
			<tr><td>strawberry</td><td>fields</td></tr>
			<tr><td>everybody's got something to hide</td><td>except for me and my monkey</td></tr>
		</table>
		<table>
			<tr><td>fool</td><td>on</td><td>the</td><td>hill</td></tr>
		</table>
		<a href="/more">next</a>
	`)

	queries := []string{
		"ul li",
		"  li",
		"ol li",
		"ol li | wc -c",
		"table",
		"  tr",
		"    td",
		"      td",
	}

	result, nextPage, err := ParseHTMLQueries(htmlData, queries, "a")

	expectedResult := []interface{}{[]interface{}{"one", "two", "three"}, "foo", "12", []interface{}{[]interface{}{[]interface{}{"penny", "lane"}, []interface{}{"strawberry", "fields"}, []interface{}{"everybody's got something to hide", "except for me and my monkey"}}, []interface{}{[]interface{}{"fool", "on", "the", "hill"}}}}

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, "/more", nextPage)
}
