package outputs

import (
	"log"
	"strings"

	"github.com/bjesus/pipet/common"
	"github.com/bjesus/pipet/utils"
)

func OutputText(app *common.PipetApp) error {
	flattenedData := utils.FlattenData(app.Data, 0)
	log.Println(strings.Join(flattenedData, utils.GetSeparator(app, 0)))
	return nil
}
