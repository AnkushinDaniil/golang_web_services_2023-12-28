package handler

import (
	"errors"
	"net/http"

	"db_explorer/explorerContext"
)

type Handler struct {
	*Table
	*Item
}

type Context interface {
	TableContext
	ItemContext
	PathLen() int
}

func NewHandler(itemSrv ItemService, tableSrv TableService) *Handler {
	return &Handler{
		Table: NewTableHandler(tableSrv),
		Item:  NewItemHandler(itemSrv),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := explorerContext.NewExplorerContext(w, r)
	h.Handle(ctx)
}

func (h *Handler) Handle(ctx Context) {
	switch ctx.PathLen() {
	case 2:
		h.Table.Handle(ctx)
	case 3:
		h.Item.Handle(ctx)
	default:
		ctx.SetResponse(http.StatusBadRequest, nil, errors.New("wrong URL"))
	}
}
