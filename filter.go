package larkbase

import (
	"fmt"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"github.com/lincaiyong/log"
	"time"
)

// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/record-filter-guide

/*
is：等于
isNot：不等于（不支持日期字段，了解如何查询日期字段，参考日期字段填写说明）
contains：包含（不支持日期字段）
doesNotContain：不包含（不支持日期字段）
isEmpty：为空
isNotEmpty：不为空

isGreater：大于
isGreaterEqual：大于等于（不支持日期字段）
isLess：小于
isLessEqual：小于等于（不支持日期字段）
like：LIKE 运算符。暂未支持
in：IN 运算符。暂未支持
*/

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
const FilterTypeDateIs = "is"
const FilterTypeDateIsGreater = "isGreater"
const FilterTypeDateIsLess = "isLess"

var gFieldFilterMap map[FieldType]map[string]bool

func init() {
	gFieldFilterMap = make(map[FieldType]map[string]bool)
	gFieldFilterMap[FieldTypeText] = map[string]bool{
		FilterTypeIs:             true,
		FilterTypeIsNot:          true,
		FilterTypeContains:       true,
		FilterTypeDoesNotContain: true,
		FilterTypeIsEmpty:        true,
		FilterTypeIsNotEmpty:     true,
	}
	gFieldFilterMap[FieldTypeNumber] = map[string]bool{
		FilterTypeIs:             true,
		FilterTypeIsNot:          true,
		FilterTypeIsGreater:      true,
		FilterTypeIsGreaterEqual: true,
		FilterTypeIsLess:         true,
		FilterTypeIsLessEqual:    true,
		FilterTypeIsEmpty:        true,
		FilterTypeIsNotEmpty:     true,
	}
	gFieldFilterMap[FieldTypeSingleSelect] = gFieldFilterMap[FieldTypeText]
	gFieldFilterMap[FieldTypeMultiSelect] = gFieldFilterMap[FieldTypeText]
	gFieldFilterMap[FieldTypeDate] = map[string]bool{
		FilterTypeDateIs:        true,
		FilterTypeDateIsGreater: true,
		FilterTypeDateIsLess:    true,
		FilterTypeIsEmpty:       true,
		FilterTypeIsNotEmpty:    true,
	}
	gFieldFilterMap[FieldTypeCheckbox] = map[string]bool{
		FilterTypeIs: true,
	}
	gFieldFilterMap[FieldTypePerson] = gFieldFilterMap[FieldTypeText]
	gFieldFilterMap[FieldTypeUrl] = gFieldFilterMap[FieldTypeText]
	gFieldFilterMap[FieldTypeMedia] = map[string]bool{
		FilterTypeIsEmpty:    true,
		FilterTypeIsNotEmpty: true,
	}
}

func FilterIs(field IField, value ...string) *larkbitable.Condition {
	filterType := FilterTypeIs
	if !gFieldFilterMap[field.Type()][filterType] {
		log.FatalLog("field type %d doesn't support filter %s", field.Type(), filterType)
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(filterType).
		Value(value).
		Build()
}

func FilterIsNot(field IField, value ...string) *larkbitable.Condition {
	filterType := FilterTypeIsNot
	if !gFieldFilterMap[field.Type()][filterType] {
		log.FatalLog("field type %d doesn't support filter %s", field.Type(), filterType)
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(filterType).
		Value(value).
		Build()
}

func FilterContains(field IField, value ...string) *larkbitable.Condition {
	filterType := FilterTypeContains
	if !gFieldFilterMap[field.Type()][filterType] {
		log.FatalLog("field type %d doesn't support filter %s", field.Type(), filterType)
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(filterType).
		Value(value).
		Build()
}

func FilterDoesNotContains(field IField, value ...string) *larkbitable.Condition {
	filterType := FilterTypeDoesNotContain
	if !gFieldFilterMap[field.Type()][filterType] {
		log.FatalLog("field type %d doesn't support filter %s", field.Type(), filterType)
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(filterType).
		Value(value).
		Build()
}

func FilterIsEmpty(field IField) *larkbitable.Condition {
	filterType := FilterTypeIsEmpty
	if !gFieldFilterMap[field.Type()][filterType] {
		log.FatalLog("field type %d doesn't support filter %s", field.Type(), filterType)
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(filterType).
		Value([]string{}).
		Build()
}

func FilterIsNotEmpty(field IField) *larkbitable.Condition {
	filterType := FilterTypeIsNotEmpty
	if !gFieldFilterMap[field.Type()][filterType] {
		log.FatalLog("field type %d doesn't support filter %s", field.Type(), filterType)
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(filterType).
		Value([]string{}).
		Build()
}

func FilterIsGreater(field IField, value string) *larkbitable.Condition {
	filterType := FilterTypeIsGreater
	if !gFieldFilterMap[field.Type()][filterType] {
		log.FatalLog("field type %d doesn't support filter %s", field.Type(), filterType)
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(filterType).
		Value([]string{value}).
		Build()
}

func FilterIsGreaterEqual(field IField, value string) *larkbitable.Condition {
	filterType := FilterTypeIsGreaterEqual
	if !gFieldFilterMap[field.Type()][filterType] {
		log.FatalLog("field type %d doesn't support filter %s", field.Type(), filterType)
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(filterType).
		Value([]string{value}).
		Build()
}

func FilterIsLess(field IField, value string) *larkbitable.Condition {
	filterType := FilterTypeIsLess
	if !gFieldFilterMap[field.Type()][filterType] {
		log.FatalLog("field type %d doesn't support filter %s", field.Type(), filterType)
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(filterType).
		Value([]string{value}).
		Build()
}

func FilterIsLessEqual(field IField, value string) *larkbitable.Condition {
	filterType := FilterTypeIsLessEqual
	if !gFieldFilterMap[field.Type()][filterType] {
		log.FatalLog("field type %d doesn't support filter %s", field.Type(), filterType)
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(filterType).
		Value([]string{value}).
		Build()
}

// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/record-filter-guide#29d9dc89

const DateToday = "Today"
const DateTomorrow = "Tomorrow"
const DateYesterday = "Yesterday"

func dateTimeStrToTime(s string) time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.FixedZone("CST", 8*60*60))
	if err != nil {
		log.FatalLog("fail to parse timestamp: %v", err)
	}
	return t
}

func dateTimeStrToTimestamp(s string) int64 {
	t := dateTimeStrToTime(s)
	return t.UnixMilli()
}

func FilterDateIs(field IField, value string) *larkbitable.Condition {
	if field.Type() != FieldTypeDate {
		log.FatalLog("expect date field type, actual: %d", field.Type())
	}
	if value == DateToday || value == DateTomorrow || value == DateYesterday {
		return larkbitable.NewConditionBuilder().
			FieldName(field.Name()).
			Operator(FilterTypeIs).
			Value([]string{value}).
			Build()
	}
	value = fmt.Sprintf("%d", dateTimeStrToTimestamp(value))
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(FilterTypeIs).
		Value([]string{"ExactDate", value}).
		Build()
}

func FilterDateIsGreater(field IField, value string) *larkbitable.Condition {
	if field.Type() != FieldTypeDate && field.Type() != FieldTypeUpdatedTime {
		log.FatalLog("expect date field type, actual: %s", field.Type().String())
	}
	if value == DateToday || value == DateTomorrow || value == DateYesterday {
		return larkbitable.NewConditionBuilder().
			FieldName(field.Name()).
			Operator(FilterTypeIsGreater).
			Value([]string{value}).
			Build()
	}
	value = fmt.Sprintf("%d", dateTimeStrToTimestamp(value))
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(FilterTypeIsGreater).
		Value([]string{"ExactDate", value}).
		Build()
}

func FilterDateIsLess(field IField, value string) *larkbitable.Condition {
	if field.Type() != FieldTypeDate && field.Type() != FieldTypeUpdatedTime {
		log.FatalLog("expect date field type, actual: %s", field.Type().String())
	}
	if value == DateToday || value == DateTomorrow || value == DateYesterday {
		return larkbitable.NewConditionBuilder().
			FieldName(field.Name()).
			Operator(FilterTypeIsLess).
			Value([]string{value}).
			Build()
	}
	value = fmt.Sprintf("%d", dateTimeStrToTimestamp(value))
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(FilterTypeIsLess).
		Value([]string{"ExactDate", value}).
		Build()
}

func FilterDateIsCurrentWeek(field IField) *larkbitable.Condition {
	if field.Type() != FieldTypeDate && field.Type() != FieldTypeUpdatedTime {
		log.FatalLog("expect date field type, actual: %s", field.Type().String())
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(FilterTypeIs).
		Value([]string{"CurrentWeek"}).
		Build()
}

func FilterDateIsLastWeek(field IField) *larkbitable.Condition {
	if field.Type() != FieldTypeDate && field.Type() != FieldTypeUpdatedTime {
		log.FatalLog("expect date field type, actual: %s", field.Type().String())
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(FilterTypeIs).
		Value([]string{"LastWeek"}).
		Build()
}

func FilterDateIsCurrentMonth(field IField) *larkbitable.Condition {
	if field.Type() != FieldTypeDate && field.Type() != FieldTypeUpdatedTime {
		log.FatalLog("expect date field type, actual: %s", field.Type().String())
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(FilterTypeIs).
		Value([]string{"CurrentMonth"}).
		Build()
}

func FilterDateIsLastMonth(field IField) *larkbitable.Condition {
	if field.Type() != FieldTypeDate && field.Type() != FieldTypeUpdatedTime {
		log.FatalLog("expect date field type, actual: %s", field.Type().String())
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(FilterTypeIs).
		Value([]string{"LastMonth"}).
		Build()
}

func FilterDateIsTheLastWeek(field IField) *larkbitable.Condition {
	// 过去七天内
	if field.Type() != FieldTypeDate && field.Type() != FieldTypeUpdatedTime {
		log.FatalLog("expect date field type, actual: %s", field.Type().String())
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(FilterTypeIs).
		Value([]string{"TheLastWeek"}).
		Build()
}

func FilterDateTheNextWeek(field IField) *larkbitable.Condition {
	// 未来七天内
	if field.Type() != FieldTypeDate && field.Type() != FieldTypeUpdatedTime {
		log.FatalLog("expect date field type, actual: %s", field.Type().String())
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(FilterTypeIs).
		Value([]string{"TheNextWeek"}).
		Build()
}

func FilterDateIsTheLastMonth(field IField) *larkbitable.Condition {
	// 过去三十天内
	if field.Type() != FieldTypeDate && field.Type() != FieldTypeUpdatedTime {
		log.FatalLog("expect date field type, actual: %s", field.Type().String())
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(FilterTypeIs).
		Value([]string{"TheLastMonth"}).
		Build()
}

func FilterDateTheNextMonth(field IField) *larkbitable.Condition {
	// 未来三十天内
	if field.Type() != FieldTypeDate && field.Type() != FieldTypeUpdatedTime {
		log.FatalLog("expect date field type, actual: %s", field.Type().String())
	}
	return larkbitable.NewConditionBuilder().
		FieldName(field.Name()).
		Operator(FilterTypeIs).
		Value([]string{"TheNextMonth"}).
		Build()
}
