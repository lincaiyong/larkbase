package field

import larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"

type TextField struct {
	BaseField
}

func (f *TextField) FilterIs(value string) *larkbitable.Condition {
	return filterIs(f.name, value)
}
func (f *TextField) FilterIsNot(value string) *larkbitable.Condition {
	return filterIsNot(f.name, value)
}
func (f *TextField) FilterContains(value string) *larkbitable.Condition {
	return filterContains(f.name, value)
}
func (f *TextField) FilterDoesNotContains(value string) *larkbitable.Condition {
	return filterDoesNotContains(f.name, value)
}
func (f *TextField) FilterIsEmpty() *larkbitable.Condition {
	return filterIsEmpty(f.name)
}
func (f *TextField) FilterIsNotEmpty() *larkbitable.Condition {
	return filterIsNotEmpty(f.name)
}
