package larkfield

type SingleSelectField struct {
	BaseField
}

func (f *SingleSelectField) SetValue(v string) {
	//if f.BaseField == nil {
	//	f.BaseField = NewBaseField(f, "", "", TypeSingleSelect)
	//}
	f.SetUnderlayValue(v)
}

func (f *SingleSelectField) Is(value string) *Condition {
	return conditionIs(f.id, f.name, value)
}
func (f *SingleSelectField) IsNot(value string) *Condition {
	return conditionIsNot(f.id, f.name, value)
}
func (f *SingleSelectField) Contains(value string) *Condition {
	return conditionContains(f.id, f.name, value)
}
func (f *SingleSelectField) DoesNotContains(value string) *Condition {
	return conditionDoesNotContains(f.id, f.name, value)
}
func (f *SingleSelectField) IsEmpty() *Condition {
	return conditionIsEmpty(f.id, f.name)
}
func (f *SingleSelectField) IsNotEmpty() *Condition {
	return conditionIsNotEmpty(f.id, f.name)
}
