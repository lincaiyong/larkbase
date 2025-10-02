package field

import larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"

type CheckboxField struct {
	BaseField
}

func (f *CheckboxField) Is(value bool) *larkbitable.Condition {
	return filterIs(f.name, map[bool]string{true: "1", false: ""}[value])
}
