package repository

import "database/sql"

type Table struct {
	db *sql.DB
}

func NewTable(db *sql.DB) *Table {
	return &Table{
		db: db,
	}
}

func (t *Table) GetAll(table string) (map[string]interface{}, error) {
	res := make(map[string]interface{}, 1)
	response := make(map[string]interface{}, 1)
	tables := make([]interface{}, 0, 2)
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
		tables = append(tables, name)
	}
	rows.Close()
	res["tables"] = tables
	response["response"] = res
	return response, nil
}
