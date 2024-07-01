package service

import "errors"

type Service struct {
	*Item
	*Table
}

var UNKNOWN_TABLE_ERROR = errors.New("unknown table")

func NewService(itemRepo ItemRepository, tableRepo TableRepository) *Service {
	return &Service{
		Item:  NewItemService(itemRepo, tableRepo),
		Table: NewTableService(tableRepo),
	}
}
