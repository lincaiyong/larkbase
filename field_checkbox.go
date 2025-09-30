package larkbase

type CheckboxField struct {
	BaseField
	value bool
}

func (f *CheckboxField) Type() FieldType {
	return FieldTypeCheckbox
}

func (f *CheckboxField) Value() string {
	if f.value {
		return "true"
	} else {
		return "false"
	}
}

func (f *CheckboxField) SetValue(v string) {
	if v == "true" {
		f.value = true
	} else {
		f.value = false
	}
}

func (f *CheckboxField) Parse(v any) IField {
	ret := &CheckboxField{BaseField: BaseField{name: f.name}}
	ret.value = v.(bool)
	return ret
}
