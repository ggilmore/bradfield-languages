package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	_ "embed"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: generate_ast <output directory>\n")
		os.Exit(64)
	}

	outputDir := os.Args[1]
	err := defineAST(outputDir, "Expr", []Struct{
		{"Binary", []Field{
			{"Left", "Expr"},
			{"Right", "Expr"},
			{"Operator", "Token"},
		}},
		{"Grouping", []Field{
			{"Expression", "Expr"},
		}},
		{"Literal", []Field{
			{"Value", "interface{}"},
		}},
		{"Unary", []Field{
			{"Operator", "Token"},
			{"Right", "Expr"},
		}},
	})

	if err != nil {
		log.Fatal(err)
	}

}

func defineAST(dir, basename string, structs []Struct) error {
	filePath := path.Join(dir, fmt.Sprintf("%s.go", strings.ToLower(basename)))
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("when opening file: %s", err)
	}

	defer f.Close()

	err = astTemplate.Execute(f, Data{
		Basename: basename,
		Structs:  structs,
	})

	if err != nil {
		return fmt.Errorf("while running template: %s", err)
	}

	return nil
}

type Data struct {
	Basename string
	Structs  []Struct
}

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
}

//go:embed ast.go.tmpl
var rawTemplate string

var astTemplate = template.Must(template.New("").Funcs(template.FuncMap{"ToLower": strings.ToLower}).Parse(rawTemplate))
