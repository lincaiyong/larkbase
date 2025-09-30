package larkbase

type SingleSelectField struct {
	BaseField
	value string
}

func (f *SingleSelectField) Type() FieldType {
	return FieldTypeSingleSelect
}

func (f *SingleSelectField) Value() string {
	return f.value
}

func (f *SingleSelectField) SetValue(v string) {
	f.value = v
}

func (f *SingleSelectField) Parse(v any) IField {
	ret := &SingleSelectField{BaseField: BaseField{name: f.name}}
	ret.value = v.(string)
	return ret
}
