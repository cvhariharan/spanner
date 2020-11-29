package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/cvhariharan/spanner/config"

	"github.com/cvhariharan/spanner/codegen"
	"github.com/kkyr/fig"
	"golang.org/x/mod/modfile"
)

func main() {
	var cfg config.Config
	err := fig.Load(&cfg)

	if len(os.Args) < 2 {
		log.Fatal("Expected models filename")
	}
	templatePath := os.Getenv("TEMPLATE")

	exec.Command("rm", "go.mod", "go.sum").Run()

	goModCmd := exec.Command("go", "mod", "init")
	err = goModCmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	cfg.ModulePath = getModPath()

	err = codegen.GenerateModel(os.Args[1], templatePath, cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = codegen.GenerateRepo(os.Args[1], templatePath, cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = codegen.GenerateServer(os.Args[1], templatePath, cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func getModPath() string {
	f, err := os.Open("go.mod")
	if err != nil {
		log.Fatal(err)
	}
	modData, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	return modfile.ModulePath(modData)
}
