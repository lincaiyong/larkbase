package larkbase

func NewTable(url string, appToken string, tableId string, fields []Field) *Table {
	fieldNames := make([]string, len(fields))
	fieldMap := make(map[string]Field, len(fields))
	for i, f := range fields {
		fieldNames[i] = f.Name()
		fieldMap[f.Name()] = f
	}
	return &Table{
		tableUrl:   url,
		appToken:   appToken,
		tableId:    tableId,
		fields:     fields,
		fieldNames: fieldNames,
		fieldMap:   fieldMap,
	}
}

type Table struct {
	tableUrl   string
	appToken   string
	tableId    string
	fields     []Field
	fieldNames []string
	fieldMap   map[string]Field
}

func (t Table) TableUrl() string {
	return t.tableUrl
}

func (t Table) AppToken() string {
	return t.appToken
}

func (t Table) TableId() string {
	return t.tableId
}

func (t Table) Fields() []Field {
	return t.fields
}

func (t Table) FieldNames() []string {
	return t.fieldNames
}

func (t Table) GetField(name string) Field {
	return t.fieldMap[name]
}
