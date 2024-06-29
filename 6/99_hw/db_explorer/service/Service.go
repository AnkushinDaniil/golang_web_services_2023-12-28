package service

type Service struct {
	*Item
	*Table
}

func NewService(itemRepo ItemRepository, tableRepo TableRepository) *Service {
	return &Service{
		Item:  NewItemService(itemRepo),
		Table: NewTableService(tableRepo),
	}
}
