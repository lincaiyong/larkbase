package field

import larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"

type NumberField struct {
	BaseField
}

func (f *NumberField) Is(value string) *larkbitable.Condition {
	return filterIs(f.name, value)
}
func (f *NumberField) IsNot(value string) *larkbitable.Condition {
	return filterIsNot(f.name, value)
}
func (f *NumberField) IsGreater(value string) *larkbitable.Condition {
	return filterIsGreater(f.name, value)
}
func (f *NumberField) IsGreaterEqual(value string) *larkbitable.Condition {
	return filterIsGreaterEqual(f.name, value)
}
func (f *NumberField) IsLess(value string) *larkbitable.Condition {
	return filterIsLess(f.name, value)
}
func (f *NumberField) IsLessEqual(value string) *larkbitable.Condition {
	return filterIsLessEqual(f.name, value)
}
func (f *NumberField) IsEmpty() *larkbitable.Condition {
	return filterIsEmpty(f.name)
}
func (f *NumberField) IsNotEmpty() *larkbitable.Condition {
	return filterIsNotEmpty(f.name)
}
