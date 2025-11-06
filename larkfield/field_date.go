package larkfield

import (
	"time"
)

type DateField struct {
	BaseField
}

func (f *DateField) SetValue(t time.Time) {
	//if f.BaseField == nil {
	//	f.BaseField = NewBaseField(f, "", "", TypeDate)
	//}
	f.SetUnderlayValue(t)
}

func (f *DateField) IsToday() *Condition {
	return conditionDateIsToday(f.id, f.name)
}
func (f *DateField) IsTomorrow() *Condition {
	return conditionDateIsTomorrow(f.id, f.name)
}
func (f *DateField) IsYesterday() *Condition {
	return conditionDateIsYesterday(f.id, f.name)
}
func (f *DateField) Is(time time.Time) *Condition {
	return conditionDateIs(f.id, f.name, time)
}
func (f *DateField) IsGreaterThanToday() *Condition {
	return conditionDateIsGreaterThanToday(f.id, f.name)
}
func (f *DateField) IsGreaterThanTomorrow() *Condition {
	return conditionDateIsGreaterThanTomorrow(f.id, f.name)
}
func (f *DateField) IsGreaterThanYesterday() *Condition {
	return conditionDateIsGreaterThanYesterday(f.id, f.name)
}
func (f *DateField) IsGreater(time time.Time) *Condition {
	return conditionDateIsGreater(f.id, f.name, time)
}
func (f *DateField) IsLessThanToday() *Condition {
	return conditionDateIsLessThanToday(f.id, f.name)
}
func (f *DateField) IsLessThanTomorrow() *Condition {
	return conditionDateIsLessThanTomorrow(f.id, f.name)
}
func (f *DateField) IsLessThanYesterday() *Condition {
	return conditionDateIsLessThanYesterday(f.id, f.name)
}

func (f *DateField) IsLess(time time.Time) *Condition {
	return conditionDateIsLess(f.id, f.name, time)
}
func (f *DateField) IsEmpty() *Condition {
	return conditionIsEmpty(f.id, f.name)
}
func (f *DateField) IsNotEmpty() *Condition {
	return conditionIsNotEmpty(f.id, f.name)
}

func (f *DateField) IsCurrentWeek() *Condition {
	return conditionDateIsCurrentWeek(f.id, f.name)
}

func (f *DateField) IsLastWeek() *Condition {
	return conditionDateIsLastWeek(f.id, f.name)
}
func (f *DateField) IsCurrentMonth() *Condition {
	return conditionDateIsCurrentMonth(f.id, f.name)
}

func (f *DateField) IsLastMonth() *Condition {
	return conditionDateIsLastMonth(f.id, f.name)
}

func (f *DateField) IsTheLastWeek() *Condition {
	return conditionDateIsTheLastWeek(f.id, f.name)
}

func (f *DateField) TheNextWeek() *Condition {
	return conditionDateTheNextWeek(f.id, f.name)
}

func (f *DateField) IsTheLastMonth() *Condition {
	return conditionDateIsTheLastMonth(f.id, f.name)
}

func (f *DateField) TheNextMonth() *Condition {
	return conditionDateTheNextMonth(f.id, f.name)
}
