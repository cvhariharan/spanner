package codegen

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cvhariharan/spanner/config"
	"github.com/markbates/pkger"
)

func GenerateMiddlewares(filename string, cfg config.Config) error {
	// Right now there is only auth middleware, so just return if it is not enabled
	// There can be more middlewares added later. If so, remove this if-return
	if !cfg.OAuth.Enable {
		return nil
	}

	env := `
AUTH_CLIENTID={{.ClientID}}
AUTH_CLIENTSECRET={{.ClientSecret}}
AUTH_CONFIGURL={{.ConfigURL}}
AUTH_REDIRECTURL={{.RedirectURL}}
`

	envTmpl := template.Must(template.New("env").Funcs(
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
		}).Parse(env))

	envGendata := struct {
		ClientID     string
		ClientSecret string
		RedirectURL  string
		ConfigURL    string
	}{
		cfg.OAuth.ClientId,
		cfg.OAuth.ClientSecret,
		cfg.OAuth.RedirectUrl,
		cfg.OAuth.ConfigUrl,
	}

	envOut, err := os.Create(".env")
	if err != nil {
		return err
	}
	err = envTmpl.Execute(envOut, envGendata)
	if err != nil {
		return err
	}

	tfr, err := pkger.Open("/codegen/templates/authmiddleware.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	br, err := ioutil.ReadAll(tfr)
	templateString := string(br)

	t := template.Must(template.New("authmiddleware.tmpl").Funcs(
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

	out, err := os.Create(filepath.Join("handler", "authmiddleware.go"))
	if err != nil {
		return err
	}

	err = t.Execute(out, struct{}{})
	if err != nil {
		return err
	}

	return nil
}
