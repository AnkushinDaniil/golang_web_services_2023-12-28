package repository

import (
	"database/sql"
	"fmt"

	"db_explorer/entity"
)

type Table struct {
	db      *sql.DB
	tables  map[string]entity.Table
	getFunc func(Scanner, []string) (entity.CR, error)
}

func NewTable(db *sql.DB, getFunc func(Scanner, []string) (entity.CR, error)) (*Table, error) {
	t := Table{
		db:      db,
		tables:  make(map[string]entity.Table, 2),
		getFunc: getFunc,
	}
	tablesArr, err := t.GetAll()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(tablesArr); i++ {
		fields, err := t.GetFields(tablesArr[i])
		if err != nil {
			return nil, err
		}
		t.tables[tablesArr[i]] = fields
	}
	return &t, nil
}

func (t *Table) CheckTable(table string) bool {
	_, ok := t.tables[table]
	return ok
}

func (t *Table) GetAll() ([]string, error) {
	response := make([]string, 0, 1)
	rows, err := t.db.Query(`SHOW TABLES`)
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

func (t *Table) GetFields(table string) (entity.Table, error) {
	rows, err := t.db.Query(fmt.Sprintf(`SHOW FULL COLUMNS FROM %s`, table))
	if err != nil {
		return nil, err
	}
	fields := make(entity.Table, 0, 4)
	for rows.Next() {
		var field entity.Field
		err = rows.Scan(
			&field.FieldName,
			&field.FieldType,
			&field.Collation,
			&field.Null,
			&field.Key,
			&field.Default,
			&field.Extra,
			&field.Privileges,
			&field.Comment,
		)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	rows.Close()
	return fields, nil
}

func (t *Table) GetList(table string, limit, offset int) ([]entity.CR, error) {
	rows, err := t.db.Query(fmt.Sprintf(`
		SELECT
		    *
		FROM
		    %s
		LIMIT ? OFFSET ?`, table),
		limit, offset)
	if err != nil {
		return nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	response := make([]entity.CR, 0, limit)
	for rows.Next() {
		m, err := t.getFunc(rows, cols)
		if err != nil {
			return nil, err
		}
		response = append(response, m)
	}
	return response, nil
}
