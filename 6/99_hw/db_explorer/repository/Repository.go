package repository

import "database/sql"

type Repository struct {
	*Item
	*Table
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Item:  NewItem(db),
		Table: NewTable(db),
	}
}
