package main

import (
	"fmt"
	"go/ast"
	"os"
	"reflect"
	"strings"
)

type Structure struct {
	Name      string
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
			case "required":
				isRequired = true
			case "min":
				minimum = keyVal[1]
			case "max":
				maximum = keyVal[1]
			case "paramname":
				paramName = keyVal[1]
				fmt.Println(paramName)
			case "default":
				defaultValue = keyVal[1]
			case "enum":
				enum = strings.Split(keyVal[1], "|")
			}
		}
	}

	if paramName == "" {
		paramName = strings.ToLower(fieldName)
	}

	newField := Field{
		Name:         fieldName,
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

func (s *Structure) genServeHTTP(out *os.File) {
	out.Write([]byte("func (srv *"))
	out.Write([]byte(s.Name))
	out.Write([]byte(`) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		response   string
		statusCode int
	)
	switch r.URL.Path {
`))
	for key := range s.Functions {
		out.Write([]byte("\tcase \""))
		out.Write([]byte(s.Functions[key].Url))
		out.Write([]byte("\":\n\t\tresponse, statusCode, err = srv.handle"))
		out.Write([]byte(s.Functions[key].Name))
		out.Write([]byte("(w, r)\n"))
	}
	out.Write([]byte(`	default:
		response, statusCode, err = "", 404, errors.New("unknown method")
	}
	w.WriteHeader(statusCode)
	if err != nil {
		w.Write([]byte("{\"error\": \""))
		w.Write([]byte(err.Error()))
		w.Write([]byte("\"}"))
	} else {
		w.Write([]byte("{\"error\":\"\", \"response\":"))
		w.Write([]byte(response))
		w.Write([]byte("}"))
	}
}
`))
}
