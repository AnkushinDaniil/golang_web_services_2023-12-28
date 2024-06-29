package handler

type ItemService interface{}

type Item struct {
	service ItemService
}

type ItemContext interface{}

func NewItemHandler(srv ItemService) *Item {
	return &Item{
		service: srv,
	}
}
