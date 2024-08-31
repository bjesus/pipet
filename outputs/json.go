package outputs

import (
	"encoding/json"
	"log"

	"github.com/bjesus/pipet/common"
)

func OutputJSON(app *common.PipetApp) error {
	jsonData, err := json.MarshalIndent(app.Data, "", "  ")
	if err != nil {
		return err
	}
	log.Println(string(jsonData))
	return nil
}
