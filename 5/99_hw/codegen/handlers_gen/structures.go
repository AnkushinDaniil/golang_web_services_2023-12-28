package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"os"
	"strings"
)

type Structures map[string]*Structure

func (s *Structures) add(f ast.Decl) {
	switch f.(type) {
	case *ast.GenDecl:
		s.addStruct(f.(*ast.GenDecl))
	case *ast.FuncDecl:
		s.addFunc(f.(*ast.FuncDecl))
	default:
		fmt.Printf("SKIP %#T is not *ast.GenDecl or *ast.FuncDecl\n", f)
	}
}

func (s *Structures) addStruct(g *ast.GenDecl) {
	for _, spec := range g.Specs {
		curType, ok := spec.(*ast.TypeSpec)
		if !ok {
			fmt.Printf("SKIP %#T is not ast.TypeSpec\n", spec)
			continue
		}
		name := curType.Name.Name

		curStruct, ok := curType.Type.(*ast.StructType)
		if !ok {
			fmt.Printf("SKIP %T is not ast.StructType\n", curStruct)
			continue
		}

		fmt.Printf("Processing %v struct\n", name)

		if _, ok := (*s)[name]; !ok {
			(*s)[name] = &Structure{
				Name:      name,
				Fields:    make(map[string]*Field),
				Functions: make(map[string]*Function),
			}
		}

		for _, newField := range curStruct.Fields.List {
			(*s)[name].addField(newField)
		}
	}
}

func (s *Structures) addFunc(f *ast.FuncDecl) {
	if f.Doc == nil || f.Recv == nil {
		return
	}
	funcName := f.Name.Name
	fmt.Printf("\t\t\t%v\n", f.Type.Params.List[1].Type)
	structName := ""
	switch a := f.Recv.List[0].Type.(type) {
	case *ast.StarExpr:
		structName = a.X.(*ast.Ident).Name
	}
	fmt.Printf("Processing %v's method %v ", structName, funcName)
	var newFunc Function
	for _, comment := range f.Doc.List {
		if strings.HasPrefix(comment.Text, "// apigen:api") {
			js := strings.TrimLeft(comment.Text, "// apigen:api")
			err := json.Unmarshal([]byte(js), &newFunc)
			newFunc.Name = funcName
			newFunc.Struct = (*s)[structName]
			newFunc.inputParamsList = f.Type.Params.List
			if err == nil {
				fmt.Printf("%+v\n", newFunc)
				(*s)[structName].Functions[funcName] = &newFunc
			}
		}
	}
}

func (s *Structures) gen(out *os.File) {
	for structName := range *s {
		isGen := false
		for funcName := range (*s)[structName].Functions {
			paramsName := fmt.Sprintf(
				"%v",
				(*s)[structName].Functions[funcName].inputParamsList[1].Type,
			)
			if _, ok := (*s)[paramsName]; ok {
				(*s)[structName].Functions[funcName].ParamsStructure = (*s)[paramsName]
				fmt.Printf(
					"Parameters detected:\n\tstructure: %v,\n\tfunction: %v,\n\tparameters: %v\n",
					structName,
					funcName,
					paramsName,
				)
				if (*s)[structName].Functions[funcName].Url != "" {
					isGen = true
					fmt.Printf("\turl: %v\n", (*s)[structName].Functions[funcName].Url)
					(*s)[structName].Functions[funcName].gen(out)
				}
			}
		}
		if isGen {
			(*s)[structName].genServeHTTP(out)
		}
	}
}
