package repository

import (
	"database/sql"

	"db_explorer/entity"
)

type Table struct {
	db *sql.DB
}

func NewTable(db *sql.DB) *Table {
	return &Table{
		db: db,
	}
}

func (t *Table) GetAll() ([]string, error) {
	response := make([]string, 0, 1)
	rows, err := t.db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		response = append(response, name)
	}
	rows.Close()
	return response, nil
}

func (t *Table) GetList(table string, limit, offset int) ([]entity.CR, error) {
	rows, err := t.db.Query("SELECT * FROM ? LIMIT ? OFFSET ?", table, limit, offset)
	if err != nil {
		return nil, err
	}
	response := make([]entity.CR, 0, 2)
	for rows.Next() {
		var row entity.CR
		err = rows.Scan(&row)
		if err != nil {
			return nil, err
		}
		response = append(response, row)
	}
	return response, nil
}
