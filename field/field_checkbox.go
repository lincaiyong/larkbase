package field

import (
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"github.com/lincaiyong/log"
)

type CheckboxField struct {
	BaseField
}

func (f *CheckboxField) Value() bool {
	if v, ok := f.value.(bool); ok {
		return v
	} else {
		log.ErrorLog("invalid underlay value for checkbox field: expect bool, got %T", f.value)
		return false
	}
}

func (f *CheckboxField) Is(value bool) *larkbitable.Condition {
	return filterIs(f.name, map[bool]string{true: "1", false: ""}[value])
}
