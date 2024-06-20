package main

import (
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := os.Create(os.Args[2])

	PackageTpl.Execute(out, node.Name.Name)

	structs := make(Structures, 16)

	for _, f := range node.Decls {
		structs.add(f)
	}

	structs.gen(out)
}
