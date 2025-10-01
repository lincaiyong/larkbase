package field

import larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"

type NumberField struct {
	BaseField
}

func (f *NumberField) FilterIs(value string) *larkbitable.Condition {
	return filterIs(f.name, value)
}
func (f *NumberField) FilterIsNot(value string) *larkbitable.Condition {
	return filterIsNot(f.name, value)
}
func (f *NumberField) FilterIsGreater(value string) *larkbitable.Condition {
	return filterIsGreater(f.name, value)
}
func (f *NumberField) FilterIsGreaterEqual(value string) *larkbitable.Condition {
	return filterIsGreaterEqual(f.name, value)
}
func (f *NumberField) FilterIsLess(value string) *larkbitable.Condition {
	return filterIsLess(f.name, value)
}
func (f *NumberField) FilterIsLessEqual(value string) *larkbitable.Condition {
	return filterIsLessEqual(f.name, value)
}
func (f *NumberField) FilterIsEmpty() *larkbitable.Condition {
	return filterIsEmpty(f.name)
}
func (f *NumberField) FilterIsNotEmpty() *larkbitable.Condition {
	return filterIsNotEmpty(f.name)
}
