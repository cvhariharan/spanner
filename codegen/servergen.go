package codegen

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/cvhariharan/spanner/config"
)

func GenerateServer(filename, templatePath string, cfg config.Config) error {
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

	t := template.Must(template.New("server.tmpl").Funcs(
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
		}).ParseFiles(templatePath + "/server.tmpl"))

	// if _, err := os.Stat("server"); os.IsNotExist(err) {
	// 	os.Mkdir("server", os.ModePerm)
	// }

	out, err := os.Create("main.go")
	if err != nil {
		return err
	}
	err = t.Execute(out, genData)
	if err != nil {
		return err
	}

	return nil
}
