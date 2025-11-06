package larkfield

import "time"

type ModifiedTimeField struct {
	*BaseField
}

func (f *ModifiedTimeField) IsGreater(time time.Time) *Condition {
	return conditionDateIsGreater(f.id, f.name, time)
}
