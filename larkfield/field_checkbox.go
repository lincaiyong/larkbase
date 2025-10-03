package larkfield

import (
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

func (f *CheckboxField) Is(value bool) *Condition {
	return conditionIs(f.id, f.name, map[bool]string{true: "1", false: ""}[value])
}
