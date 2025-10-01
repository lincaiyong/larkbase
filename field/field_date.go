package field

import (
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"time"
)

type DateField struct {
	BaseField
}

func (f *DateField) FilterIsToday() *larkbitable.Condition {
	return filterDateIsToday(f.name)
}
func (f *DateField) FilterIsTomorrow() *larkbitable.Condition {
	return filterDateIsTomorrow(f.name)
}
func (f *DateField) FilterIsYesterday() *larkbitable.Condition {
	return filterDateIsYesterday(f.name)
}
func (f *DateField) FilterIs(time time.Time) *larkbitable.Condition {
	return filterDateIs(f.name, time)
}
func (f *DateField) FilterIsGreaterThanToday() *larkbitable.Condition {
	return filterDateIsGreaterThanToday(f.name)
}
func (f *DateField) FilterIsGreaterThanTomorrow() *larkbitable.Condition {
	return filterDateIsGreaterThanTomorrow(f.name)
}
func (f *DateField) FilterIsGreaterThanYesterday() *larkbitable.Condition {
	return filterDateIsGreaterThanYesterday(f.name)
}
func (f *DateField) FilterIsGreater(time time.Time) *larkbitable.Condition {
	return filterDateIsGreater(f.name, time)
}
func (f *DateField) FilterIsLessThanToday() *larkbitable.Condition {
	return filterDateIsLessThanToday(f.name)
}
func (f *DateField) FilterIsLessThanTomorrow() *larkbitable.Condition {
	return filterDateIsLessThanTomorrow(f.name)
}
func (f *DateField) FilterIsLessThanYesterday() *larkbitable.Condition {
	return filterDateIsLessThanYesterday(f.name)
}

func (f *DateField) FilterIsLess(time time.Time) *larkbitable.Condition {
	return filterDateIsLess(f.name, time)
}
func (f *DateField) FilterIsEmpty() *larkbitable.Condition {
	return filterIsEmpty(f.name)
}
func (f *DateField) FilterIsNotEmpty() *larkbitable.Condition {
	return filterIsNotEmpty(f.name)
}

func (f *DateField) FilterIsCurrentWeek() *larkbitable.Condition {
	return filterDateIsCurrentWeek(f.name)
}

func (f *DateField) FilterIsLastWeek() *larkbitable.Condition {
	return filterDateIsLastWeek(f.name)
}
func (f *DateField) FilterIsCurrentMonth() *larkbitable.Condition {
	return filterDateIsCurrentMonth(f.name)
}

func (f *DateField) FilterIsLastMonth() *larkbitable.Condition {
	return filterDateIsLastMonth(f.name)
}

func (f *DateField) FilterIsTheLastWeek() *larkbitable.Condition {
	return filterDateIsTheLastWeek(f.name)
}

func (f *DateField) FilterTheNextWeek() *larkbitable.Condition {
	return filterDateTheNextWeek(f.name)
}

func (f *DateField) FilterIsTheLastMonth() *larkbitable.Condition {
	return filterDateIsTheLastMonth(f.name)
}

func (f *DateField) FilterTheNextMonth() *larkbitable.Condition {
	return filterDateTheNextMonth(f.name)
}
