package utils

import (
	"fmt"

	"github.com/bjesus/pipet/common"
)

func FlattenData(data interface{}, depth int) []string {
	var result []string

	switch v := data.(type) {
	case []interface{}:
		for _, item := range v {
			result = append(result, FlattenData(item, depth+1)...)
		}
	case map[string]interface{}:
		for _, value := range v {
			result = append(result, FlattenData(value, depth+1)...)
		}
	default:
		result = append(result, fmt.Sprintf("%v", v))
	}

	return result
}

func GetSeparator(app *common.PipetApp, depth int) string {
	if depth < len(app.Separator) {
		return app.Separator[depth]
	}
	return ", "
}
