package main

import (
	"go/ast"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type Function struct {
	Name            string
	Struct          *Structure
	Url             string `json:"url"`
	Auth            bool   `json:"auth"`
	Method          string `json:"method"`
	ParamsStructure *Structure
	inputParamsList []*ast.Field
}

func (f *Function) genStart(out *os.File) {
	out.Write([]byte("\nfunc (srv *"))
	out.Write([]byte(f.Struct.Name))
	out.Write([]byte(") handle"))
	out.Write([]byte(f.Name))
	out.Write([]byte("(w http.ResponseWriter, r *http.Request) (string, int, error) {\n"))
	out.Write([]byte("\tw.Header().Set(\"Content-Type\", \"application/json\")\n"))
}

func (f *Function) genAuth(out *os.File) {
	if f.Auth {
		out.Write([]byte(AuthCode))
	}
}

func (f *Function) genCheckPost(out *os.File) {
	if f.Method == http.MethodPost {
		out.Write([]byte(MethodPostCode))
	}
}

func (f *Function) genVar(out *os.File) {
	out.Write([]byte("\tvar (\n"))
	for key := range f.ParamsStructure.Fields {
		out.Write([]byte("\t\t"))
		out.Write([]byte(f.ParamsStructure.Fields[key].ParamName))
		out.Write([]byte(" string\n"))
	}
	out.Write([]byte("\t)\n\n"))
}

func (f *Function) genGet(out *os.File) error {
	if f.Method == http.MethodPost {
		for key := range f.ParamsStructure.Fields {
			err := f.ParamsStructure.Fields[key].genPost(out)
			if err != nil {
				return errors.Wrap(err, "field genGet error")
			}
		}
	} else {
		out.Write([]byte(`	if r.Method != "GET" {
	r.ParseForm()
	`))
		for key := range f.ParamsStructure.Fields {
			err := f.ParamsStructure.Fields[key].genPost(out)
			if err != nil {
				return errors.Wrap(err, "field genPost error")
			}
		}
		out.Write([]byte(`	} else {
	`))
		for key := range f.ParamsStructure.Fields {
			err := f.ParamsStructure.Fields[key].genGet(out)
			if err != nil {
				return errors.Wrap(err, "field genGet error")
			}
		}
		out.Write([]byte("\t}\n"))
	}
	out.Write([]byte("\n"))

	return nil
}

func (f *Function) genValidation(out *os.File) error {
	for key := range f.ParamsStructure.Fields {
		err := f.ParamsStructure.Fields[key].genValidation(out)
		if err != nil {
			return errors.Wrap(err, "function genValidation error")
		}
	}
	return nil
}

func (f *Function) genParams(out *os.File) {
	out.Write([]byte("\tparams := "))
	out.Write([]byte(f.ParamsStructure.Name))
	out.Write([]byte("{\n"))
	for key := range f.ParamsStructure.Fields {
		out.Write([]byte("\t\t"))
		out.Write([]byte(key))
		out.Write([]byte(":\t"))
		out.Write([]byte(f.ParamsStructure.Fields[key].ParamName))
		if f.ParamsStructure.Fields[key].IsInt {
			out.Write([]byte("Int"))
		}
		out.Write([]byte(",\n"))
	}
	out.Write([]byte("\t}\n"))
}

func (f *Function) genPC(out *os.File) error {
	err := ResponseTpl.Execute(out, f)
	if err != nil {
		return errors.Wrap(err, ResponseTpl.Name()+" execution error")
	}
	return nil
}

func (f *Function) genEnd(out *os.File) {
	out.Write([]byte("\treturn string(response), http.StatusOK, nil\n}\n"))
}

func (f *Function) gen(out *os.File) error {
	f.genStart(out)
	f.genAuth(out)
	f.genCheckPost(out)
	f.genVar(out)
	err := f.genGet(out)
	if err != nil {
		return errors.Wrap(err, "function genGet error")
	}
	err = f.genValidation(out)
	if err != nil {
		return errors.Wrap(err, "function genValidation error")
	}
	f.genParams(out)
	err = f.genPC(out)
	if err != nil {
		return errors.Wrap(err, "function genPC error")
	}
	f.genEnd(out)
	return nil
}
