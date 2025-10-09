package larkfield

type MultiSelectField struct {
	BaseField
}

func (f *MultiSelectField) SetValue(v []string) {
	f.SetUnderlayValue(v)
}
