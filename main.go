package main

import (
	"log"
	"os"

	"github.com/cvhariharan/spanner/codegen"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Expected models filename")
	}
	templatePath := os.Getenv("TEMPLATE")
	err := codegen.GenerateModel(os.Args[1], templatePath)
	if err != nil {
		log.Fatal(err)
	}

	err = codegen.GenerateRepo(os.Args[1], templatePath)
	if err != nil {
		log.Fatal(err)
	}

	err = codegen.GenerateServer(os.Args[1], templatePath)
	if err != nil {
		log.Fatal(err)
	}
}
