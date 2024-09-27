package utils

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/bjesus/pipet/common"
)

func FlattenNestedSlices(app *common.PipetApp, data interface{}, level int) string {
	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Slice {
		return fmt.Sprint(data)
	}

	var result []string

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		flattened := FlattenNestedSlices(app, elem, level+1)
		result = append(result, flattened)
	}

	sep := GetSeparator(app, level)
	return strings.Join(result, sep)
}

func RemoveUnnecessaryNesting(data interface{}) interface{} {
	for {
		val := reflect.ValueOf(data)
		if val.Kind() == reflect.Slice && val.Len() == 1 {
			firstElem := val.Index(0).Interface()
			if reflect.ValueOf(firstElem).Kind() == reflect.Slice {
				data = firstElem
				continue
			}
		}
		break
	}
	return data
}

func GetSeparator(app *common.PipetApp, depth int) string {
	if depth < len(app.Separator) {
		sep, _ := strconv.Unquote(`"` + app.Separator[depth] + `"`)
		return sep
	}
	return ", "
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
