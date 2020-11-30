package codegen

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/cvhariharan/spanner/config"
	"github.com/cvhariharan/spanner/parser"
	"github.com/markbates/pkger"
)

func GenerateModel(filename string, cfg config.Config) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var model map[string]interface{}
	err = json.Unmarshal(data, &model)
	if err != nil {
		return err
	}

	keys := reflect.ValueOf(model).MapKeys()
	if len(keys) != 1 {
		return errors.New("Expected only one model definition")
	}

	modelName := keys[0].String()

	modelDef, ok := model[modelName].(map[string]interface{})
	if !ok {
		return errors.New("Problem parsing model definition")
	}
	modelMap := parser.ParseModel(modelDef, modelName)

	genData := struct {
		PackageName   string
		MainModelName string
		ModelMap      map[string]map[string]string
	}{
		"model",
		modelName,
		modelMap,
	}

	tf, err := pkger.Open("/codegen/templates/model.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(tf)
	templateString := string(b)

	t := template.Must(template.New("model.tmpl").Funcs(
		template.FuncMap{
			"Title": strings.Title,
			"TitleLower": func(s string) string {
				if len(s) < 2 {
					return strings.ToLower(s)
				}
				bts := []byte(s)
				lc := bytes.ToLower([]byte{bts[0]})
				rest := bts[1:]
				return string(bytes.Join([][]byte{lc, rest}, nil))
			},
		}).Parse(templateString))

	if _, err := os.Stat("modules"); os.IsNotExist(err) {
		os.Mkdir("modules", os.ModePerm)
	}

	if _, err := os.Stat(filepath.Join("modules", modelName)); os.IsNotExist(err) {
		os.Mkdir(filepath.Join("modules", modelName), os.ModePerm)
	}

	if _, err := os.Stat(filepath.Join("modules", modelName, "model")); os.IsNotExist(err) {
		os.Mkdir(filepath.Join("modules", modelName, "model"), os.ModePerm)
	}

	out, err := os.Create(filepath.Join("modules", modelName, "model", "model.gen.go"))
	if err != nil {
		return err
	}
	err = t.Execute(out, genData)
	if err != nil {
		return err
	}

	return nil
}
