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

var templateList []string = []string{"knativeservice.tmpl", "handler.tmpl", "main.tmpl", "makefile.tmpl", "model.tmpl", "mongorepo.tmpl", "repo.tmpl", "dockerfile.tmpl"}
var packageMap = map[string]string{
	"handler.tmpl":        "handler",
	"authmiddleware.tmpl": "handler",
	"main.tmpl":           "main",
	"model.tmpl":          "model",
	"mongorepo.tmpl":      "repo",
	"repo.tmpl":           "repo",
}

var fileName = map[string]string{
	"handler.tmpl":        "handler.go",
	"authmiddleware.tmpl": "middleware.go",
	"main.tmpl":           "main.go",
	"model.tmpl":          "model.gen.go",
	"mongorepo.tmpl":      "mongorepo.gen.go",
	"repo.tmpl":           "repo.gen.go",
	"makefile.tmpl":       "Makefile",
	"env.tmpl":            ".env",
	"dockerfile.tmpl":     "Dockerfile",
	"knativeservice.tmpl": "service.yml",
}

func Generate(filename string, cfg config.Config) error {
	if cfg.OAuth.Enable {
		templateList = append(templateList, "authmiddleware.tmpl")
		templateList = append(templateList, "env.tmpl")
	}

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

	var prefixPathMap = map[string]string{
		"handler.tmpl":        "handler",
		"authmiddleware.tmpl": "handler",
		"model.tmpl":          filepath.Join("modules", modelName, "model"),
		"mongorepo.tmpl":      filepath.Join("modules", modelName, "repo"),
		"repo.tmpl":           filepath.Join("modules", modelName, "repo"),
	}

	generateDirectories(prefixPathMap)

	genData := struct {
		PackageName    string
		ModelName      string
		ModuleName     string
		AuthEnable     bool
		MainModelName  string
		ModelMap       map[string]map[string]string
		ClientID       string
		ClientSecret   string
		RedirectURL    string
		ConfigURL      string
		Port           string
		DockerUsername string
	}{
		"",
		modelName,
		cfg.ModulePath,
		cfg.OAuth.Enable,
		modelName,
		modelMap,
		cfg.OAuth.ClientId,
		cfg.OAuth.ClientSecret,
		cfg.OAuth.RedirectUrl,
		cfg.OAuth.ConfigUrl,
		cfg.Port,
		cfg.DockerUsername,
	}

	// log.Println(prefixPathMap["model.tmpl"])

	for _, v := range templateList {
		genData.PackageName = packageMap[v]

		tf, err := pkger.Open("/codegen/templates/" + v)
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(tf)
		templateString := string(b)

		t := template.Must(template.New(v).Funcs(
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
				"Upper": strings.ToUpper,
			}).Parse(templateString))

		out, err := os.Create(filepath.Join(prefixPathMap[v], fileName[v]))
		if err != nil {
			return err
		}

		err = t.Execute(out, genData)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateDirectories(directoriesMap map[string]string) {
	for _, v := range directoriesMap {
		if _, err := os.Stat(v); os.IsNotExist(err) {
			os.MkdirAll(v, os.ModePerm)
		}
	}
}
