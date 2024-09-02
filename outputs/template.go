package outputs

import (
	"bytes"
	"html/template"

	"github.com/bjesus/pipet/common"
)

func OutputTemplate(app *common.PipetApp, templateFile string) string {
	tmpl, _ := template.ParseFiles(templateFile)
	var doc bytes.Buffer
	tmpl.Execute(&doc, app.Data)
	return doc.String()
}
