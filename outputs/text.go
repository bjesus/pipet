package outputs

import (
	"strings"

	"github.com/bjesus/pipet/common"
	"github.com/bjesus/pipet/utils"
)

func OutputText(app *common.PipetApp) string {
	flattenedData := utils.FlattenData(app.Data, 0)
	return (strings.Join(flattenedData, utils.GetSeparator(app, 0)))

}
