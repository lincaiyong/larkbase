package larkfield

type MediaField struct {
	BaseField
}

func (f *MediaField) IsEmpty() *Condition {
	return conditionIsEmpty(f.id, f.name)
}
func (f *MediaField) IsNotEmpty() *Condition {
	return conditionIsNotEmpty(f.id, f.name)
}
