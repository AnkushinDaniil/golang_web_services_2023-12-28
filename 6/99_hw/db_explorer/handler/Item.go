package handler

import (
	"net/http"

	"db_explorer/entity"
)

type ItemService interface {
	GetById(string, int) (entity.CR, error)
}

type Item struct {
	service ItemService
}

type ItemContext interface {
	GetTableName() (string, error)
	GetItemId() (int, error)
	SetResponse(int, entity.CR, error)
	SendResponse()
	Method() string
}

func NewItemHandler(srv ItemService) *Item {
	return &Item{
		service: srv,
	}
}

func (h *Item) Handle(ctx ItemContext) {
	switch ctx.Method() {
	case http.MethodGet:
		h.handleGet(ctx)
	}
	ctx.SendResponse()
}

func (h *Item) handleGet(ctx ItemContext) {
	table, err := ctx.GetTableName()
	if err != nil {
		ctx.SetResponse(http.StatusBadRequest, nil, err)
		return
	}
	id, err := ctx.GetItemId()
	if err != nil {
		ctx.SetResponse(http.StatusBadRequest, nil, err)
		return
	}
	response, err := h.service.GetById(table, id)
	ctx.SetResponse(http.StatusOK, response, nil)
}
