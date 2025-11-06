package larkfield

import (
	"github.com/lincaiyong/log"
	"strconv"
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

func (f *AutoNumberField) Is(value int) *Condition {
	return conditionIs(f.id, f.name, strconv.Itoa(value))
}

func (f *AutoNumberField) IsNot(value int) *Condition {
	return conditionIsNot(f.id, f.name, strconv.Itoa(value))
}

func (f *AutoNumberField) IsGreater(value int) *Condition {
	return conditionIsGreater(f.id, f.name, strconv.Itoa(value))
}

func (f *AutoNumberField) IsGreaterEqual(value int) *Condition {
	return conditionIsGreaterEqual(f.id, f.name, strconv.Itoa(value))
}

func (f *AutoNumberField) IsLess(value int) *Condition {
	return conditionIsLess(f.id, f.name, strconv.Itoa(value))
}

func (f *AutoNumberField) IsLessEqual(value int) *Condition {
	return conditionIsLessEqual(f.id, f.name, strconv.Itoa(value))
}
