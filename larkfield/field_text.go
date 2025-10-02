package larkfield

import larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"

type TextField struct {
	BaseField
}

func (f *TextField) SetValue(v string) {
	f.SetUnderlayValue(v)
}

func (f *TextField) Is(value string) *larkbitable.Condition {
	return filterIs(f.name, value)
}
func (f *TextField) IsNot(value string) *larkbitable.Condition {
	return filterIsNot(f.name, value)
}
func (f *TextField) Contains(value string) *larkbitable.Condition {
	return filterContains(f.name, value)
}
func (f *TextField) DoesNotContains(value string) *larkbitable.Condition {
	return filterDoesNotContains(f.name, value)
}
func (f *TextField) IsEmpty() *larkbitable.Condition {
	return filterIsEmpty(f.name)
}
func (f *TextField) IsNotEmpty() *larkbitable.Condition {
	return filterIsNotEmpty(f.name)
}
