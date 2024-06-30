package repository

import (
	"database/sql"
	"fmt"

	"db_explorer/entity"
)

type Table struct {
	db     *sql.DB
	tables map[string]entity.Table
}

func NewTable(db *sql.DB) (*Table, error) {
	t := Table{
		db:     db,
		tables: make(map[string]entity.Table, 2),
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

func (t *Table) GetFields(table string) ([]entity.Field, error) {
	rows, err := t.db.Query(fmt.Sprintf(`SHOW FULL COLUMNS FROM %s`, table))
	if err != nil {
		return nil, err
	}
	fields := make([]entity.Field, 0, 4)
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
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := 0; i < len(columns); i++ {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(entity.CR, len(cols))
		for i := 0; i < len(cols); i++ {
			switch val := (*(columnPointers[i].(*interface{}))).(type) {
			case int:
				m[cols[i]] = val
			case string:
				m[cols[i]] = val
			case []uint8:
				m[cols[i]] = string(val)
			default:
				m[cols[i]] = val
			}
		}
		response = append(response, m)
	}
	return response, nil
}
