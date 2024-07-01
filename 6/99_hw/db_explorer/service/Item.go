package service

import "db_explorer/entity"

type ItemRepository interface {
	GetById(string, int) (entity.CR, error)
}

type Item struct {
	itemRepo  ItemRepository
	tableRepo TableRepository
}

func (i *Item) GetById(table string, id int) (entity.CR, error) {
	if !i.tableRepo.CheckTable(table) {
		return nil, UNKNOWN_TABLE_ERROR
	}
	records, err := i.itemRepo.GetById(table, id)
	if err != nil {
		return nil, err
	}
	response := make(entity.CR, 1)
	response["record"] = records
	return response, nil
}

func NewItemService(itemRepo ItemRepository, tableRepo TableRepository) *Item {
	return &Item{
		itemRepo:  itemRepo,
		tableRepo: tableRepo,
	}
}
