package codegen

import (
	"bytes"
	"os"
	"strings"
	"text/template"
)

func GenerateRepo(filename, templatePath string) error {
	// f, err := os.Open(filename)
	// if err != nil {
	// 	return err
	// }

	// data, err := ioutil.ReadAll(f)
	// if err != nil {
	// 	return err
	// }

	// var model map[string]interface{}
	// err = json.Unmarshal(data, &model)
	// if err != nil {
	// 	return err
	// }

	// keys := reflect.ValueOf(model).MapKeys()
	// if len(keys) != 1 {
	// 	return errors.New("Expected only one model definition")
	// }

	// modelName := keys[0].String()

	// modelDef, ok := model[modelName].(map[string]interface{})
	// if !ok {
	// 	return errors.New("Problem parsing model definition")
	// }
	// modelMap := parser.ParseModel(modelDef, modelName)

	// genData := struct {
	// 	PackageName string
	// 	ModelMap    map[string]map[string]string
	// }{
	// 	"model",
	// 	modelMap,
	// }

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

	if _, err := os.Stat("repo"); os.IsNotExist(err) {
		os.Mkdir("repo", os.ModePerm)
	}

	outrepo, err := os.Create("repo/repo.gen.go")
	if err != nil {
		return err
	}

	outmongo, err := os.Create("repo/mongorepo.gen.go")
	if err != nil {
		return err
	}

	err = repo.Execute(outrepo, struct{}{})
	if err != nil {
		return err
	}

	err = mongo.Execute(outmongo, struct{}{})
	if err != nil {
		return err
	}

	return nil
}
