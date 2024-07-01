package repository

import (
	"database/sql"
	"fmt"

	"db_explorer/entity"
)

type Item struct {
	db      *sql.DB
	getFunc func(Scanner, []string) (entity.CR, error)
}

func NewItem(db *sql.DB, getFunc func(Scanner, []string) (entity.CR, error)) *Item {
	return &Item{
		db:      db,
		getFunc: getFunc,
	}
}

func (i *Item) GetById(table string, id int) (entity.CR, error) {
	rows, err := i.db.Query(fmt.Sprintf(`
		SELECT
		    *
		FROM
		    %s
		WHERE
		    id = ?`, table),
		id)
	if err != nil {
		return nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	rows.Next()
	return i.getFunc(rows, cols)
}
