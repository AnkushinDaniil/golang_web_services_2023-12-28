package repository

import "database/sql"

type Repository struct {
	*Item
	*Table
}

func NewRepository(db *sql.DB) (*Repository, error) {
	table, err := NewTable(db)
	if err != nil {
		return nil, err
	}
	return &Repository{
		Item:  NewItem(db),
		Table: table,
	}, nil
}
