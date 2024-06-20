package main

import "text/template"

const AuthCode = `	if r.Header.Get("X-Auth") != "100500" {
		return "", http.StatusForbidden, errors.New("unauthorized")
	}
`

const MethodPostCode = `	if r.Method != "POST" {
		return "", http.StatusNotAcceptable, errors.New("bad method")
	}
	r.ParseForm()
`

var (
	PostTpl = template.Must(
		template.New("postTpl").Parse(`	{{.ParamName}} = r.Form.Get("{{.ParamName}}")
`),
	)

	PackageTpl = template.Must(
		template.New("PackageTpl").Parse(`package {{.}}

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)
`),
	)

	GetTpl = template.Must(
		template.New("getTpl").Parse(`	{{.ParamName}} = r.URL.Query().Get("{{.ParamName}}")
`),
	)

	RequiredTpl = template.Must(
		template.New("requiredTpl").Parse(`	if {{.ParamName}} == "" {
		return "", http.StatusBadRequest, errors.New("{{.ParamName}} must be not empty")
	}
`),
	)
	IntCheckTpl = template.Must(
		template.New("intCheckTpl").
			Parse(`	{{.ParamName}}Int, err := strconv.Atoi({{.ParamName}})
	if err != nil {
		return "", http.StatusBadRequest, errors.New("{{.ParamName}} must be int")
	}
`),
	)
	StrMinCheckTpl = template.Must(
		template.New("strMinCheckTpl").
			Parse(`	if len({{.ParamName}}) < {{.Minimum}} {
		return "", http.StatusBadRequest, errors.New("{{.ParamName}} len must be >= {{.Minimum}}")
	}
`),
	)
	StrMaxCheckTpl = template.Must(
		template.New("strMaxCheckTpl").
			Parse(`	if len({{.ParamName}}) < {{.Maximum}} {
		return "", http.StatusBadRequest, errors.New("{{.ParamName}} len must be <= {{.Maximum}}")
	}
`),
	)
	IntMinCheckTpl = template.Must(
		template.New("intMinCheckTpl").
			Parse(`	if {{.ParamName}}Int < {{.Minimum}} {
		return "", http.StatusBadRequest, errors.New("{{.ParamName}} must be >= {{.Minimum}}")
	}
`),
	)
	IntMaxCheckTpl = template.Must(
		template.New("intMaxCheckTpl").
			Parse(`	if {{.ParamName}}Int > {{.Maximum}} {
		return "", http.StatusBadRequest, errors.New("{{.ParamName}} must be <= {{.Maximum}}")
	}
`),
	)
	DefaultValueTpl = template.Must(
		template.New("defaultValueTpl").
			Parse(`	if {{.ParamName}} == "" {
		{{.ParamName}} = "{{.DefaultValue}}"
	}
`),
	)
	ResponseTpl = template.Must(
		template.New("responseTpl").
			Parse(`
	response{{.Name}}, err := srv.{{.Name}}(r.Context(), params)
	if err != nil {
		switch err.(type) {
		case ApiError:
			return "", err.(ApiError).HTTPStatus, err
		default:
			return "", http.StatusInternalServerError, err
		}
	}

	response, err := json.Marshal(response{{.Name}})
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
`),
	)
)
