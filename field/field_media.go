package field

import larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"

type MediaField struct {
	BaseField
}

func (f *MediaField) FilterIsEmpty() *larkbitable.Condition {
	return filterIsEmpty(f.name)
}
func (f *MediaField) FilterIsNotEmpty() *larkbitable.Condition {
	return filterIsNotEmpty(f.name)
}
