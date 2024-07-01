package repository

import (
	"database/sql"

	"db_explorer/entity"
)

type Repository struct {
	*Item
	*Table
}

type Scanner interface {
	Scan(...any) error
}

func getFunc(scanner Scanner, cols []string) (entity.CR, error) {
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := 0; i < len(columns); i++ {
		columnPointers[i] = &columns[i]
	}
	if err := scanner.Scan(columnPointers...); err != nil {
		return nil, err
	}

	response := make(entity.CR, len(cols))
	for i := 0; i < len(cols); i++ {
		switch val := (*(columnPointers[i].(*interface{}))).(type) {
		case []uint8:
			response[cols[i]] = string(val)
		default:
			response[cols[i]] = val
		}
	}

	return response, nil
}

func NewRepository(db *sql.DB) (*Repository, error) {
	table, err := NewTable(db, getFunc)
	if err != nil {
		return nil, err
	}
	return &Repository{
		Item:  NewItem(db, getFunc),
		Table: table,
	}, nil
}
