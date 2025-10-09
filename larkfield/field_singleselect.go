package larkfield

type SingleSelectField struct {
	BaseField
}

func (f *SingleLinkField) SetValue(v string) {
	f.SetUnderlayValue(v)
}
