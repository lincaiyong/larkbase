package larkbase

type DateField struct {
	BaseField
	value string // beijing
}

func (f *DateField) Type() FieldType {
	return FieldTypeDate
}

func (f *DateField) Value() string {
	return f.value
}

func (f *DateField) SetValue(v string) {
	f.value = v
}

func (f *DateField) Parse(v any) IField {
	ret := &DateField{BaseField: BaseField{name: f.name}}
	t := unixSecondsToBeijingDateTimeStr(int64(v.(float64) / 1000))
	ret.value = t
	return ret
}
