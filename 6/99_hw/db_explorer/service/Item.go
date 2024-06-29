package service

type ItemRepository interface{}

type Item struct {
	repo ItemRepository
}

func NewItemService(repo ItemRepository) *Item {
	return &Item{repo: repo}
}
