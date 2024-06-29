package context

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Response struct {
	data       map[string]interface{}
	statusCode int
	err        error
}

type ExplorerContext struct {
	path     []string
	w        http.ResponseWriter
	r        *http.Request
	response Response
}

func NewExplorerContext(w http.ResponseWriter, r *http.Request) *ExplorerContext {
	return &ExplorerContext{
		path: strings.Split(r.URL.Path, "/"),
		w:    w,
		r:    r,
	}
}

func (ctx *ExplorerContext) GetTableName() (string, error) {
	if len(ctx.path) < 2 {
		return "", errors.New("wrong path")
	}
	return ctx.path[1], nil
}

func (ctx *ExplorerContext) SetResponse(statusCode int, data map[string]interface{}, err error) {
	ctx.response = Response{
		data:       data,
		statusCode: statusCode,
		err:        err,
	}
}

func (r *Response) Bytes() ([]byte, error) {
	if r.err != nil {
		r.data["error"] = r.err.Error()
	}
	return json.Marshal(r.data)
}

func (ctx *ExplorerContext) SendResponse() {
	data, err := ctx.response.Bytes()
	if err != nil {
		ctx.w.WriteHeader(http.StatusInternalServerError)
		ctx.w.Write([]byte(`{"error": "` + err.Error() + `"}`))
	}
	ctx.w.WriteHeader(ctx.response.statusCode)
	ctx.w.Write([]byte(data))
}

func (ctx *ExplorerContext) PathLen() int {
	return len(ctx.path)
}
