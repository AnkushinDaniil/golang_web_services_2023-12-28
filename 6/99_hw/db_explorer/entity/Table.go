package entity

type Table []Field

func (t Table) FieldsNames() []string {
	fieldsNames := make([]string, len(t))
	for i := 0; i < len(t); i++ {
		fieldsNames[i] = t[i].FieldName
	}
	return fieldsNames
}
