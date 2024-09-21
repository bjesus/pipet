package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseJSONQueries(t *testing.T) {
	JSONData := []byte(`
		{"person": {"name": "bjesus", "friends": [{ "name": "bob", "nickname": "rocky"}, {"name":"alice", "nickname": "racoon"}]}}
	`)

	queries := []string{
		"person.name",
		"person.name | wc -c",
		"person.friends",
		"  name",
		"  nickname",
	}

	result, nextPage, err := ParseJSONQueries(JSONData, queries)

	expectedResult := []interface{}{
		"bjesus",
		float64(6),
		[]interface{}{[]interface{}{"bob", "rocky"}, []interface{}{"alice", "racoon"}},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, "", nextPage)
}
