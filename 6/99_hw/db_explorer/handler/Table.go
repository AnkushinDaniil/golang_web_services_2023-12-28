package handler

import (
	"net/http"
)

type TableService interface {
	GetAll(string) (map[string]interface{}, error)
}

type TableContext interface {
	GetTableName() (string, error)
	SetResponse(int, map[string]interface{}, error)
	SendResponse()
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
	defer ctx.SendResponse()
	table, err := ctx.GetTableName()
	if err != nil {
		ctx.SetResponse(http.StatusBadRequest, nil, err)
		return
	}
	if table == "" {
		h.GetAll(ctx)
	} else {
		h.GetList(ctx)
	}
}

func (h *Table) GetAll(ctx TableContext) {
	table, err := ctx.GetTableName()
	if err != nil {
		ctx.SetResponse(http.StatusBadRequest, nil, err)
		return
	}
	response, err := h.service.GetAll(table)
	if err != nil {
		ctx.SetResponse(http.StatusInternalServerError, nil, err)
		return
	}
	ctx.SetResponse(http.StatusOK, response, err)
}

func (h *Table) GetList(ctx TableContext) {
}
