package larkbase

type AutoNumberField struct {
	BaseField
	value *string
}

func (f *AutoNumberField) Type() FieldType {
	return FieldTypeAutoNumber
}

func (f *AutoNumberField) Value() string {
	if f.value == nil {
		return ""
	}
	return *f.value
}

func (f *AutoNumberField) SetValue(v string) {
	f.value = &v
}

func (f *AutoNumberField) Parse(v any) IField {
	ret := &AutoNumberField{BaseField: BaseField{name: f.name}}
	n := v.(string)
	ret.value = &n
	return ret
}
