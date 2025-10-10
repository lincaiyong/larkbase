package larkfield

type UrlField struct {
	BaseField
}

func (f *UrlField) SetValue(v string) {
	f.SetUnderlayValue(v)
}
