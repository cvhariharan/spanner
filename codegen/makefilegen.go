package codegen

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/cvhariharan/spanner/config"
	"github.com/markbates/pkger"
)

func GenerateMakefile(filename, templatePath string, cfg config.Config) error {
	tf, err := pkger.Open("/codegen/templates/makefile.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(tf)
	templateString := string(b)

	t := template.Must(template.New("makefile.tmpl").Funcs(
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

	out, err := os.Create("Makefile")
	if err != nil {
		return err
	}
	err = t.Execute(out, struct{}{})
	if err != nil {
		return err
	}

	return nil
}
