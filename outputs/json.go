package outputs

import (
	"encoding/json"

	"github.com/bjesus/pipet/common"
)

func OutputJSON(app *common.PipetApp) string {
	jsonData, _ := json.MarshalIndent(app.Data, "", "  ")
	return string(jsonData)

}
