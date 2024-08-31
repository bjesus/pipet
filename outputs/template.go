package outputs

import (
	"html/template"
	"os"

	"github.com/bjesus/pipet/common"
)

func OutputTemplate(app *common.PipetApp, templateFile string) error {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	return tmpl.Execute(os.Stdout, app.Data)
}
