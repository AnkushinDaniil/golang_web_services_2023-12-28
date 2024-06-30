package handler

import (
	"errors"
	"net/http"

	"db_explorer/context"
	"db_explorer/entity"
)

const (
	DEFAULT_LIMIT  = 5
	DEFAULT_OFFSET = 0
)

type TableService interface {
	GetAll() (entity.CR, error)
	GetList(string, int, int) (entity.CR, error)
}

type TableContext interface {
	GetTableName() (string, error)
	SetResponse(int, entity.CR, error)
	SendResponse()
	Method() string
	GetInt(string) (int, error)
}

type Table struct {
	service TableService
}

func NewTableHandler(srv TableService) *Table {
	return &Table{
		service: srv,
	}
}

func (h *Table) Handle(ctx TableContext) {
	switch ctx.Method() {
	case http.MethodGet:
		h.handleGet(ctx)
	}
	ctx.SendResponse()
}

func (h *Table) handleGet(ctx TableContext) {
	table, err := ctx.GetTableName()
	var response entity.CR
	if err != nil {
		ctx.SetResponse(http.StatusBadRequest, nil, err)
		return
	}
	if table == "" {
		response, err = h.service.GetAll()
	} else {
		limit, err := ctx.GetInt("limit")
		if err != nil {
			if errors.Is(err, context.EMPTY_PARAM_ERROR) {
				limit = DEFAULT_LIMIT
			} else {
				ctx.SetResponse(http.StatusInternalServerError, nil, err)
				return
			}
		}
		offset, err := ctx.GetInt("offset")
		if err != nil {
			if errors.Is(err, context.EMPTY_PARAM_ERROR) {
				offset = DEFAULT_OFFSET
			} else {
				ctx.SetResponse(http.StatusInternalServerError, nil, err)
				return
			}
		}
		response, err = h.service.GetList(table, limit, offset)
		if err != nil {
			ctx.SetResponse(http.StatusNotFound, nil, err)
			return
		}
	}

	ctx.SetResponse(http.StatusOK, response, nil)
}
