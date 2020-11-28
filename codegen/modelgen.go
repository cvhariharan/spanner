package codegen

import (
	"os"
	"strings"
	"text/template"

	"github.com/cvhariharan/spanner/parser"
)

func GenerateModel(filename string) error {
	modelMap, err := parser.ParseModelFromJSON(filename)
	if err != nil {
		return err
	}

	genData := struct {
		PackageName string
		ModelMap    map[string]map[string]string
	}{
		"model",
		modelMap,
	}

	t := template.Must(template.New("model.tmpl").Funcs(
		template.FuncMap{
			"Title": strings.Title,
		}).ParseFiles("codegen/templates/model.tmpl"))

	if _, err := os.Stat("model"); os.IsNotExist(err) {
		os.Mkdir("model", os.ModePerm)
	}

	out, err := os.Create("model/model.gen.go")
	if err != nil {
		return err
	}
	err = t.Execute(out, genData)
	if err != nil {
		return err
	}

	return nil
}
