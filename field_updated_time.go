package larkbase

type UpdatedTimeField struct {
	BaseField
	value string
}

func (f *UpdatedTimeField) Type() FieldType {
	return FieldTypeUpdatedTime
}

func (f *UpdatedTimeField) Value() string {
	return f.value
}

func (f *UpdatedTimeField) SetValue(v string) {
	f.value = v
}

func (f *UpdatedTimeField) Parse(v any) IField {
	ret := &UpdatedTimeField{BaseField: BaseField{name: f.name}}
	tt := unixSecondsToBeijingDateTimeStr(int64(v.(float64) / 1000))
	ret.value = tt
	return ret
}
