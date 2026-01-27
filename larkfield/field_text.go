package larkfield

import "github.com/lincaiyong/log"

type TextField struct {
	BaseField
}

func (f *TextField) SetValue(v string) {
	if len(v) > 80000 {
		log.WarnLog("set value too long: %d, truncate it", len(v))
		v = v[:80000] + "..."
	}
	f.SetUnderlayValue(v)
}

func (f *TextField) Is(value string) *Condition {
	return conditionIs(f.id, f.name, value)
}
func (f *TextField) IsNot(value string) *Condition {
	return conditionIsNot(f.id, f.name, value)
}
func (f *TextField) Contains(value string) *Condition {
	return conditionContains(f.id, f.name, value)
}
func (f *TextField) DoesNotContains(value string) *Condition {
	return conditionDoesNotContains(f.id, f.name, value)
}
func (f *TextField) IsEmpty() *Condition {
	return conditionIsEmpty(f.id, f.name)
}
func (f *TextField) IsNotEmpty() *Condition {
	return conditionIsNotEmpty(f.id, f.name)
}
