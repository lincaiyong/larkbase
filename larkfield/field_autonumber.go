package larkfield

import (
	"github.com/lincaiyong/log"
)

type AutoNumberField struct {
	BaseField
}

func (f *AutoNumberField) Value() int {
	if v, ok := f.value.(int); ok {
		return v
	} else {
		log.ErrorLog("invalid underlay value for autonumber field: expect int, got %T", f.value)
		return 0
	}
}
