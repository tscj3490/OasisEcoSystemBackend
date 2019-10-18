package mail

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

var mainTemplate = "main.html"
var templateDir = "mail/templates"

/*ParseEmailFromTemplate Generates the html template from the route return the coplited template*/
func ParseEmailFromTemplate(templateFileName string, data interface{}) (string, error) {

	// Get the absolute template location
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	path := filepath.Join(exPath, templateDir)

	mainPath := filepath.Join(path, mainTemplate)
	templatePath := filepath.Join(path, templateFileName)

	t, err := template.ParseFiles(templatePath, mainPath)
	if err != nil {
		return "", err
	}

	buff := new(bytes.Buffer)
	if err := t.ExecuteTemplate(buff, "base", data); err != nil {

		return "", err
	}

	return buff.String(), nil
}
