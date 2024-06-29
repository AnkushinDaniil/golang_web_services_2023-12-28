package service

type TableRepository interface {
	GetAll(string) (map[string]interface{}, error)
}

type Table struct {
	repo TableRepository
}

func (t *Table) GetAll(table string) (map[string]interface{}, error) {
	return t.repo.GetAll(table)
}

func NewTableService(repo TableRepository) *Table {
	return &Table{repo: repo}
}
