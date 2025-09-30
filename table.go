package larkbase

import "sort"

func NewTable(id, name string, fields map[string]IField) *Table {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return &Table{
		id:        id,
		name:      name,
		fields:    fields,
		fieldKeys: keys,
	}
}

type Table struct {
	id        string
	name      string
	fields    map[string]IField
	fieldKeys []string
}
