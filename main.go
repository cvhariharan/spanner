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
	err := codegen.GenerateModel(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}
