package larkbase

import (
	"fmt"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"time"
)

// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/record-filter-guide
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/record-filter-guide#29d9dc89

const FilterTypeIs = "is"
const FilterTypeIsNot = "isNot"
const FilterTypeContains = "contains"
const FilterTypeDoesNotContain = "doesNotContain"
const FilterTypeIsEmpty = "isEmpty"
const FilterTypeIsNotEmpty = "isNotEmpty"
const FilterTypeIsGreater = "isGreater"
const FilterTypeIsGreaterEqual = "isGreaterEqual"
const FilterTypeIsLess = "isLess"
const FilterTypeIsLessEqual = "isLessEqual"

func filterIs(name, value string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIs).Value([]string{value}).Build()
}

func filterIsNot(name, value string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsNot).Value([]string{value}).Build()
}

func filterContains(name, value string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeContains).Value([]string{value}).Build()
}

func filterDoesNotContains(name, value string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeDoesNotContain).Value([]string{value}).Build()
}

func filterIsEmpty(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsEmpty).Value([]string{}).Build()
}

func filterIsNotEmpty(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsNotEmpty).Value([]string{}).Build()
}

func filterIsGreater(name, value string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsGreater).Value([]string{value}).Build()
}

func filterIsGreaterEqual(name, value string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsGreaterEqual).Value([]string{value}).Build()
}

func filterIsLess(name, value string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsLess).Value([]string{value}).Build()
}

func filterIsLessEqual(name, value string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsLessEqual).
		Value([]string{value}).Build()
}

func filterDateIsToday(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIs).Value([]string{"Today"}).Build()
}

func filterDateIsTomorrow(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIs).Value([]string{"Tomorrow"}).Build()
}

func filterDateIsYesterday(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIs).Value([]string{"Yesterday"}).Build()
}

func filterDateIs(name string, time time.Time) *larkbitable.Condition {
	value := fmt.Sprintf("%d", time.UnixMilli())
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIs).Value([]string{"ExactDate", value}).Build()
}

func filterDateIsGreaterThanToday(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsGreater).Value([]string{"Today"}).Build()
}

func filterDateIsGreaterThanTomorrow(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsGreater).Value([]string{"Tomorrow"}).Build()
}

func filterDateIsGreaterThanYesterday(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsGreater).Value([]string{"Yesterday"}).Build()
}

func filterDateIsGreater(name string, time time.Time) *larkbitable.Condition {
	value := fmt.Sprintf("%d", time.UnixMilli())
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsGreater).Value([]string{"ExactDate", value}).Build()
}

func filterDateIsLessThanToday(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsLess).Value([]string{"Today"}).Build()
}

func filterDateIsLessThanTomorrow(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsLess).Value([]string{"Tomorrow"}).Build()
}

func filterDateIsLessThanYesterday(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsLess).Value([]string{"Yesterday"}).Build()
}

func filterDateIsLess(name string, time time.Time) *larkbitable.Condition {
	value := fmt.Sprintf("%d", time.UnixMilli())
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIsLess).Value([]string{"ExactDate", value}).Build()
}

func filterDateIsCurrentWeek(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIs).Value([]string{"CurrentWeek"}).Build()
}

func filterDateIsLastWeek(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIs).Value([]string{"LastWeek"}).Build()
}

func filterDateIsCurrentMonth(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIs).Value([]string{"CurrentMonth"}).Build()
}

func filterDateIsLastMonth(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIs).Value([]string{"LastMonth"}).Build()
}

func filterDateIsTheLastWeek(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIs).Value([]string{"TheLastWeek"}).Build()
}

func filterDateTheNextWeek(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIs).Value([]string{"TheNextWeek"}).Build()
}

func filterDateIsTheLastMonth(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIs).Value([]string{"TheLastMonth"}).Build()
}

func filterDateTheNextMonth(name string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().FieldName(name).Operator(FilterTypeIs).Value([]string{"TheNextMonth"}).Build()
}
