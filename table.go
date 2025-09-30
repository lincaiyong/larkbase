package larkbase

func NewTable(id, name string, fields map[string]FieldType) *Table {
	return &Table{
		id:     id,
		name:   name,
		fields: fields,
	}
}

type Table struct {
	id     string
	name   string
	fields map[string]FieldType
}
