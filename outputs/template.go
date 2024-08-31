package outputs

import (
	"html/template"
	"os"

	"github.com/bjesus/erol/common"
)

func OutputTemplate(app *common.ErolApp, templateFile string) error {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	return tmpl.Execute(os.Stdout, app.Data)
}
