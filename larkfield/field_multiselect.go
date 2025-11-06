package larkfield

type MultiSelectField struct {
	BaseField
}

func (f *MultiSelectField) SetValue(v []string) {
	//if f.BaseField == nil {
	//	f.BaseField = NewBaseField(f, "", "", TypeMultiSelect)
	//}
	f.SetUnderlayValue(v)
}

func (f *MultiSelectField) Is(value string) *Condition {
	return conditionIs(f.id, f.name, value)
}
func (f *MultiSelectField) IsNot(value string) *Condition {
	return conditionIsNot(f.id, f.name, value)
}
func (f *MultiSelectField) Contains(value string) *Condition {
	return conditionContains(f.id, f.name, value)
}
func (f *MultiSelectField) DoesNotContains(value string) *Condition {
	return conditionDoesNotContains(f.id, f.name, value)
}
func (f *MultiSelectField) IsEmpty() *Condition {
	return conditionIsEmpty(f.id, f.name)
}
func (f *MultiSelectField) IsNotEmpty() *Condition {
	return conditionIsNotEmpty(f.id, f.name)
}
