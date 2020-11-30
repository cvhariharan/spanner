package codegen

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/cvhariharan/spanner/config"
)

func GenerateRepo(filename, templatePath string, cfg config.Config) error {
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

	// modelDef, ok := model[modelName].(map[string]interface{})
	// if !ok {
	// 	return errors.New("Problem parsing model definition")
	// }
	// modelMap := parser.ParseModel(modelDef, modelName)

	genData := struct {
		ModelName  string
		ModuleName string
	}{
		strings.Title(modelName),
		cfg.ModulePath,
	}

	repo := template.Must(template.New("repo.tmpl").Funcs(
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
		}).ParseFiles(templatePath + "/repo.tmpl"))

	mongo := template.Must(template.New("mongorepo.tmpl").Funcs(
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
		}).ParseFiles(templatePath + "/mongorepo.tmpl"))

	if _, err := os.Stat("modules"); os.IsNotExist(err) {
		os.Mkdir("modules", os.ModePerm)
	}

	if _, err := os.Stat(filepath.Join("modules", modelName)); os.IsNotExist(err) {
		os.Mkdir(filepath.Join("modules", modelName), os.ModePerm)
	}

	if _, err := os.Stat(filepath.Join("modules", modelName, "repo")); os.IsNotExist(err) {
		os.Mkdir(filepath.Join("modules", modelName, "repo"), os.ModePerm)
	}

	outrepo, err := os.Create(filepath.Join("modules", modelName, "repo", "repo.gen.go"))
	if err != nil {
		return err
	}

	outmongo, err := os.Create(filepath.Join("modules", modelName, "repo", "mongorepo.gen.go"))
	if err != nil {
		return err
	}

	err = repo.Execute(outrepo, genData)
	if err != nil {
		return err
	}

	err = mongo.Execute(outmongo, genData)
	if err != nil {
		return err
	}

	return nil
}
