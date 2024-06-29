package repository

import "database/sql"

type Item struct {
	db *sql.DB
}

func NewItem(db *sql.DB) *Item {
	return &Item{
		db: db,
	}
}
