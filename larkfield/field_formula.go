package larkfield

type FormulaField struct {
	BaseField
}

func (f *FormulaField) IsEmpty() *Condition {
	return conditionIsEmpty(f.id, f.name)
}
func (f *FormulaField) IsNotEmpty() *Condition {
	return conditionIsNotEmpty(f.id, f.name)
}
