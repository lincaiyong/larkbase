package field

import (
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"time"
)

type DateField struct {
	BaseField
}

func (f *DateField) IsToday() *larkbitable.Condition {
	return filterDateIsToday(f.name)
}
func (f *DateField) IsTomorrow() *larkbitable.Condition {
	return filterDateIsTomorrow(f.name)
}
func (f *DateField) IsYesterday() *larkbitable.Condition {
	return filterDateIsYesterday(f.name)
}
func (f *DateField) Is(time time.Time) *larkbitable.Condition {
	return filterDateIs(f.name, time)
}
func (f *DateField) IsGreaterThanToday() *larkbitable.Condition {
	return filterDateIsGreaterThanToday(f.name)
}
func (f *DateField) IsGreaterThanTomorrow() *larkbitable.Condition {
	return filterDateIsGreaterThanTomorrow(f.name)
}
func (f *DateField) IsGreaterThanYesterday() *larkbitable.Condition {
	return filterDateIsGreaterThanYesterday(f.name)
}
func (f *DateField) IsGreater(time time.Time) *larkbitable.Condition {
	return filterDateIsGreater(f.name, time)
}
func (f *DateField) IsLessThanToday() *larkbitable.Condition {
	return filterDateIsLessThanToday(f.name)
}
func (f *DateField) IsLessThanTomorrow() *larkbitable.Condition {
	return filterDateIsLessThanTomorrow(f.name)
}
func (f *DateField) IsLessThanYesterday() *larkbitable.Condition {
	return filterDateIsLessThanYesterday(f.name)
}

func (f *DateField) IsLess(time time.Time) *larkbitable.Condition {
	return filterDateIsLess(f.name, time)
}
func (f *DateField) IsEmpty() *larkbitable.Condition {
	return filterIsEmpty(f.name)
}
func (f *DateField) IsNotEmpty() *larkbitable.Condition {
	return filterIsNotEmpty(f.name)
}

func (f *DateField) IsCurrentWeek() *larkbitable.Condition {
	return filterDateIsCurrentWeek(f.name)
}

func (f *DateField) IsLastWeek() *larkbitable.Condition {
	return filterDateIsLastWeek(f.name)
}
func (f *DateField) IsCurrentMonth() *larkbitable.Condition {
	return filterDateIsCurrentMonth(f.name)
}

func (f *DateField) IsLastMonth() *larkbitable.Condition {
	return filterDateIsLastMonth(f.name)
}

func (f *DateField) IsTheLastWeek() *larkbitable.Condition {
	return filterDateIsTheLastWeek(f.name)
}

func (f *DateField) TheNextWeek() *larkbitable.Condition {
	return filterDateTheNextWeek(f.name)
}

func (f *DateField) IsTheLastMonth() *larkbitable.Condition {
	return filterDateIsTheLastMonth(f.name)
}

func (f *DateField) TheNextMonth() *larkbitable.Condition {
	return filterDateTheNextMonth(f.name)
}
