package outputs

import (
	"github.com/bjesus/pipet/common"
	"github.com/bjesus/pipet/utils"
)

func OutputText(app *common.PipetApp) string {
	cleanData := utils.RemoveUnnecessaryNesting(app.Data)
	return utils.FlattenNestedSlices(app, cleanData, 0)
}
