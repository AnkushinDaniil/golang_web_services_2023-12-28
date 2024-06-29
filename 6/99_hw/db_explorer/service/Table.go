package service

import (
	"errors"
	"slices"

	"db_explorer/entity"
)

var UNKNOWN_TABLE_ERROR = errors.New("unknown table")

type TableRepository interface {
	GetAll() ([]string, error)
	GetList(string, int, int) ([]entity.CR, error)
}

type Table struct {
	repo TableRepository
}

func (t *Table) GetAll() (entity.CR, error) {
	tables, err := t.repo.GetAll()
	if err != nil {
		return nil, err
	}
	response := make(entity.CR, 1)
	response["tables"] = tables
	return response, nil
}

func (t *Table) GetList(table string, limit, offset int) (entity.CR, error) {
	tables, err := t.repo.GetAll()
	if err != nil {
		return nil, err
	}
	if !slices.Contains(tables, table) {
		return nil, UNKNOWN_TABLE_ERROR
	}
	records, err := t.repo.GetList(table, limit, offset)
	if err != nil {
		return nil, err
	}
	response := make(entity.CR, 1)
	response["records"] = records
	return response, nil
}

func NewTableService(repo TableRepository) *Table {
	return &Table{repo: repo}
}
