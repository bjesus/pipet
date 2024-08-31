package outputs

import (
	"log"
	"strings"

	"github.com/bjesus/erol/common"
	"github.com/bjesus/erol/utils"
)

func OutputText(app *common.ErolApp) error {
	flattenedData := utils.FlattenData(app.Data, 0)
	log.Println(strings.Join(flattenedData, utils.GetSeparator(app, 0)))
	return nil
}
