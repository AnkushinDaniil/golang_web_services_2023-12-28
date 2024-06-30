package context

import (
	"db_explorer/entity"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var EMPTY_PARAM_ERROR = errors.New("Parameter is empty")

type Response struct {
	data       entity.CR
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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
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

func (ctx *ExplorerContext) SetResponse(statusCode int, data entity.CR, err error) {
	ctx.response = Response{
		data:       data,
		statusCode: statusCode,
		err:        err,
	}
}

func (r *Response) Bytes() ([]byte, error) {
	response := make(entity.CR, 1)
	if r.err != nil {
		response["error"] = r.err.Error()
	} else {
		response["response"] = r.data
	}
	return json.Marshal(response)
}

func (ctx *ExplorerContext) SendResponse() {
	data, err := ctx.response.Bytes()
	if err != nil {
		ctx.w.WriteHeader(http.StatusInternalServerError)
		ctx.w.Write([]byte(`{"error": "` + err.Error() + `"}`))
	}
	ctx.w.WriteHeader(ctx.response.statusCode)
	ctx.w.Write(data)
	fmt.Println(string(data))
}

func (ctx *ExplorerContext) PathLen() int {
	return len(ctx.path)
}

func (ctx *ExplorerContext) Method() string {
	return ctx.r.Method
}

func (ctx *ExplorerContext) GetStr(param string) string {
	var response string
	if ctx.r.Method != http.MethodGet {
		ctx.r.ParseForm()
		response = ctx.r.Form.Get(param)
	} else {
		response = ctx.r.URL.Query().Get(param)
	}
	return response
}

func (ctx *ExplorerContext) GetInt(param string) (int, error) {
	responseStr := ctx.GetStr(param)
	if responseStr == "" {
		return 0, EMPTY_PARAM_ERROR
	}
	return strconv.Atoi(responseStr)
}
