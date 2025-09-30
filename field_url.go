package larkbase

type UrlField struct {
	BaseField
	value string
}

func (f *UrlField) Type() FieldType {
	return FieldTypeUrl
}

func (f *UrlField) Value() string {
	return f.value
}

func (f *UrlField) SetValue(v string) {
	f.value = v
}

func (f *UrlField) Parse(v any) IField {
	ret := &UrlField{BaseField: BaseField{name: f.name}}
	if v, ok := v.(map[string]any); ok {
		ret.value = v["link"].(string)
	}
	return ret
}
