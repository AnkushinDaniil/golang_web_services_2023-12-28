package main

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Field struct {
	Name         string
	IsInt        bool
	Minimum      string
	Maximum      string
	IsRequired   bool
	ParamName    string
	Enum         []string
	DefaultValue string
}

func (f *Field) genPost(out *os.File) error {
	return PostTpl.Execute(out, f)
}

func (f *Field) genGet(out *os.File) error {
	return GetTpl.Execute(out, f)
}

func (f *Field) genRequired(out *os.File) error {
	if f.IsRequired {
		err := RequiredTpl.Execute(out, f)
		if err != nil {
			return errors.Wrap(err, RequiredTpl.Name()+" execution error")
		}
	}
	return nil
}

func (f *Field) genBounds(out *os.File) error {
	if f.IsInt {
		err := IntCheckTpl.Execute(out, f)
		if err != nil {
			return errors.Wrap(err, IntCheckTpl.Name()+"execution error")
		}
		if f.Minimum != "" {
			err := IntMinCheckTpl.Execute(out, f)
			if err != nil {
				return errors.Wrap(err, IntMinCheckTpl.Name()+" execution error")
			}
		}
		if f.Maximum != "" {
			err := IntMaxCheckTpl.Execute(out, f)
			if err != nil {
				return errors.Wrap(err, IntMaxCheckTpl.Name()+"execution error")
			}
		}
	} else {
		if f.Minimum != "" {
			err := StrMinCheckTpl.Execute(out, f)
			if err != nil {
				return errors.Wrap(err, StrMinCheckTpl.Name()+"execution error")
			}
		}
		if f.Maximum != "" {
			err := StrMaxCheckTpl.Execute(out, f)
			if err != nil {
				return errors.Wrap(err, StrMaxCheckTpl.Name()+"maximum integer genBounds error")
			}
		}
	}
	return nil
}

func (f *Field) genDefault(out *os.File) error {
	if f.DefaultValue != "" {
		err := DefaultValueTpl.Execute(out, f)
		if err != nil {
			return errors.Wrap(err, DefaultValueTpl.Name()+"execution error")
		}
	}
	return nil
}

func (f *Field) genEnum(out *os.File) {
	if len(f.Enum) > 0 {
		out.Write([]byte("\tif "))
		out.Write([]byte(f.ParamName))
		out.Write([]byte(" != \""))
		out.Write([]byte(f.Enum[0]))
		out.Write([]byte("\" "))
		for i := 1; i < len(f.Enum); i++ {
			out.Write([]byte("&& "))
			out.Write([]byte(f.ParamName))
			out.Write([]byte(" != \""))
			out.Write([]byte(f.Enum[i]))
			out.Write([]byte("\" "))
		}
		out.Write([]byte(`{
		return "", http.StatusBadRequest, errors.New("`))
		out.Write([]byte(f.ParamName))
		out.Write([]byte(" must be one of ["))
		out.Write([]byte(strings.Join(f.Enum, ", ")))
		out.Write([]byte(`]")
	}`))
	}
}

func (f *Field) genValidation(out *os.File) error {
	err := f.genRequired(out)
	if err != nil {
		return errors.Wrap(err, "genRequired error")
	}
	err = f.genBounds(out)
	if err != nil {
		return errors.Wrap(err, "genBounds error")
	}
	err = f.genDefault(out)
	if err != nil {
		return errors.Wrap(err, "genDefault error")
	}
	f.genEnum(out)
	out.Write([]byte("\n"))
	return nil
}
