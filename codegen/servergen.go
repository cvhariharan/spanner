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
	"github.com/markbates/pkger"
)

func GenerateServer(filename string, cfg config.Config) error {
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

	genData := struct {
		PackageName string
		ModelName   string
		ModuleName  string
	}{
		"server",
		modelName,
		cfg.ModulePath,
	}

	tfm, err := pkger.Open("/codegen/templates/main.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	bm, err := ioutil.ReadAll(tfm)
	mainTemplateString := string(bm)

	tfh, err := pkger.Open("/codegen/templates/handler.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	bh, err := ioutil.ReadAll(tfh)
	handlerTemplateString := string(bh)

	tm := template.Must(template.New("main.tmpl").Funcs(
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
		}).Parse(mainTemplateString))

	th := template.Must(template.New("handler.tmpl").Funcs(
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
		}).Parse(handlerTemplateString))

	if _, err := os.Stat("handler"); os.IsNotExist(err) {
		os.Mkdir("handler", os.ModePerm)
	}

	outm, err := os.Create("main.go")
	if err != nil {
		return err
	}
	outh, err := os.Create(filepath.Join("handler", "handler.go"))
	if err != nil {
		return err
	}
	err = tm.Execute(outm, genData)
	if err != nil {
		return err
	}
	err = th.Execute(outh, genData)
	if err != nil {
		return err
	}

	return nil
}
