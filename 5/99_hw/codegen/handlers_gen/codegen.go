package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"
	"strings"
)

type Field struct {
	IsInt        bool
	Minimum      string
	Maximum      string
	IsRequired   bool
	ParamName    string
	Enum         []string
	DefaultValue string
}

type Function struct {
	Url    string `json:"url"`
	Auth   bool   `json:"auth"`
	Method string `json:"method"`
}

type Structure struct {
	Fields    map[string]*Field
	Functions map[string]*Function
}

func (s *Structure) addField(f *ast.Field) {
	var isInt bool
	fieldName := f.Names[0].Name

	fmt.Printf("	Processing %v field: ", fieldName)

	switch a := f.Type.(type) {
	case *ast.Ident:
		isInt = a.Name == "int"
	default:
		isInt = false
	}

	isRequired := false
	minimum := ""
	maximum := ""
	paramName := ""
	defaultValue := ""
	var enum []string

	if f.Tag != nil {
		tag := reflect.StructTag(f.Tag.Value[1 : len(f.Tag.Value)-1])
		validators := strings.Split(tag.Get("apivalidator"), ",")
		for _, validator := range validators {
			keyVal := strings.Split(validator, "=")
			switch keyVal[0] {
			case "min":
				minimum = keyVal[1]
			case "max":
				maximum = keyVal[1]
			case "paramname":
				paramName = keyVal[1]
			case "default":
				defaultValue = keyVal[1]
			case "enum":
				enum = strings.Split(keyVal[1], "|")
			}
		}
	}

	newField := Field{
		IsInt:        isInt,
		Minimum:      minimum,
		Maximum:      maximum,
		IsRequired:   isRequired,
		ParamName:    paramName,
		Enum:         enum,
		DefaultValue: defaultValue,
	}

	fmt.Printf("%+v\n", newField)

	s.Fields[fieldName] = &newField
}

type structures map[string]*Structure

func (s *structures) add(f ast.Decl) {
	switch f.(type) {
	case *ast.GenDecl:
		s.addStruct(f.(*ast.GenDecl))
	case *ast.FuncDecl:
		s.addFunc(f.(*ast.FuncDecl))
	default:
		fmt.Printf("SKIP %#T is not *ast.GenDecl or *ast.FuncDecl\n", f)
	}
}

func (s *structures) addStruct(g *ast.GenDecl) {
	for _, spec := range g.Specs {
		currType, ok := spec.(*ast.TypeSpec)
		if !ok {
			fmt.Printf("SKIP %#T is not ast.TypeSpec\n", spec)
			continue
		}
		name := currType.Name.Name

		currStruct, ok := currType.Type.(*ast.StructType)
		if !ok {
			fmt.Printf("SKIP %T is not ast.StructType\n", currStruct)
			continue
		}

		fmt.Printf("Processing %v struct\n", name)

		if _, ok := (*s)[name]; !ok {
			(*s)[name] = &Structure{
				Fields:    make(map[string]*Field),
				Functions: make(map[string]*Function),
			}
		}

		for _, newField := range currStruct.Fields.List {
			(*s)[name].addField(newField)
		}
	}
}

func (s *structures) addFunc(f *ast.FuncDecl) {
	if f.Doc == nil || f.Recv == nil {
		return
	}
	funcName := f.Name.Name
	structName := ""
	switch a := f.Recv.List[0].Type.(type) {
	case *ast.StarExpr:
		structName = a.X.(*ast.Ident).Name
	}
	fmt.Printf("Processing %v's method %v ", structName, funcName)
	var newFunc Function
	for _, comment := range f.Doc.List {
		js := strings.TrimLeft(comment.Text, "// apigen:api")
		err := json.Unmarshal([]byte(js), &newFunc)
		if err == nil {
			fmt.Printf("%+v\n", newFunc)
			(*s)[structName].Functions[funcName] = &newFunc
		}
	}
}

func main() {
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := os.Create(os.Args[2])

	out.WriteString("package ")
	out.WriteString(node.Name.Name)
	out.WriteString(`

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)
`)

	structs := make(structures, 16)

	for _, f := range node.Decls {
		structs.add(f)
	}
}
