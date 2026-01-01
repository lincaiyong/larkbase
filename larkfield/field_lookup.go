package larkfield

type LookupField struct {
	BaseField
}

func (f *LookupField) IsEmpty() *Condition {
	return conditionIsEmpty(f.id, f.name)
}
func (f *LookupField) IsNotEmpty() *Condition {
	return conditionIsNotEmpty(f.id, f.name)
}
