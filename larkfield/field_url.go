package larkfield

type UrlField struct {
	BaseField
}

func (f *UrlField) SetValue(v string) {
	f.SetUnderlayValue(v)
}

func (f *UrlField) Is(value string) *Condition {
	return conditionIs(f.id, f.name, value)
}

func (f *UrlField) IsNot(value string) *Condition {
	return conditionIsNot(f.id, f.name, value)
}

func (f *UrlField) Contains(value string) *Condition {
	return conditionContains(f.id, f.name, value)
}

func (f *UrlField) DoesNotContains(value string) *Condition {
	return conditionDoesNotContains(f.id, f.name, value)
}

func (f *UrlField) IsEmpty() *Condition {
	return conditionIsEmpty(f.id, f.name)
}

func (f *UrlField) IsNotEmpty() *Condition {
	return conditionIsNotEmpty(f.id, f.name)
}
