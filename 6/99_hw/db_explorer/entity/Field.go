package entity

type Field struct {
	FieldName  string
	FieldType  string
	Collation  *string
	Null       string
	Key        string
	Default    *string
	Extra      string
	Privileges string
	Comment    string
}
