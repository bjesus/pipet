package outputs

import (
	"encoding/json"
	"log"

	"github.com/bjesus/erol/common"
)

func OutputJSON(app *common.ErolApp) error {
	jsonData, err := json.MarshalIndent(app.Data, "", "  ")
	if err != nil {
		return err
	}
	log.Println(string(jsonData))
	return nil
}
