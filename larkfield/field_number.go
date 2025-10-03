package larkfield

type NumberField struct {
	BaseField
}

func (f *NumberField) SetIntValue(v int) {
	f.SetUnderlayValue(float64(v))
}

func (f *NumberField) SetValue(v float64) {
	f.SetUnderlayValue(v)
}

func (f *NumberField) Is(value string) *Condition {
	return conditionIs(f.id, f.name, value)
}

func (f *NumberField) IsNot(value string) *Condition {
	return conditionIsNot(f.id, f.name, value)
}

func (f *NumberField) IsGreater(value string) *Condition {
	return conditionIsGreater(f.id, f.name, value)
}

func (f *NumberField) IsGreaterEqual(value string) *Condition {
	return conditionIsGreaterEqual(f.id, f.name, value)
}

func (f *NumberField) IsLess(value string) *Condition {
	return conditionIsLess(f.id, f.name, value)
}

func (f *NumberField) IsLessEqual(value string) *Condition {
	return conditionIsLessEqual(f.id, f.name, value)
}

func (f *NumberField) IsEmpty() *Condition {
	return conditionIsEmpty(f.id, f.name)
}

func (f *NumberField) IsNotEmpty() *Condition {
	return conditionIsNotEmpty(f.id, f.name)
}
