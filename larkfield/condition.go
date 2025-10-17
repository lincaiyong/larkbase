package larkfield

import (
	"encoding/json"
	"fmt"
	"github.com/lincaiyong/larkbase/larksuite/bitable"
	"time"
)

func NewCondition(fieldId, fieldName, operator string, value []string) *Condition {
	return &Condition{
		fieldId:   fieldId,
		fieldName: fieldName,
		operator:  operator,
		value:     value,
	}
}

type Condition struct {
	fieldId   string
	fieldName string
	operator  string
	value     []string
}

func (c *Condition) FieldId() string {
	return c.fieldId
}

func (c *Condition) FieldName() string {
	return c.fieldName
}

func (c *Condition) Operator() string {
	return c.operator
}

func (c *Condition) Value() []string {
	return c.value
}

func (c *Condition) ToLarkCondition() *bitable.Condition {
	if c.operator == ConditionOpIsEmpty || c.operator == ConditionOpIsNotEmpty {
		return bitable.NewConditionBuilder().FieldName(c.fieldName).Operator(c.operator).Value(make([]string, 0)).Build()
	}
	return bitable.NewConditionBuilder().FieldName(c.fieldName).Operator(c.operator).Value(c.value).Build()
}

func (c *Condition) ToLarkViewCondition() *bitable.AppTableViewPropertyFilterInfoCondition {
	if c.operator == ConditionOpIsEmpty || c.operator == ConditionOpIsNotEmpty {
		return bitable.NewAppTableViewPropertyFilterInfoConditionBuilder().FieldId(c.fieldId).Operator(c.operator).Build()
	}
	b, _ := json.Marshal(c.value)
	return bitable.NewAppTableViewPropertyFilterInfoConditionBuilder().FieldId(c.fieldId).Operator(c.operator).Value(string(b)).Build()
}

// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/record-filter-guide
// https://open.larkoffice.com/document/docs/bitable-v1/app-table-record/record-filter-guide#29d9dc89

const ConditionOpIs = "is"
const ConditionOpIsNot = "isNot"
const ConditionOpContains = "contains"
const ConditionOpDoesNotContain = "doesNotContain"
const ConditionOpIsEmpty = "isEmpty"
const ConditionOpIsNotEmpty = "isNotEmpty"
const ConditionOpIsGreater = "isGreater"
const ConditionOpIsGreaterEqual = "isGreaterEqual"
const ConditionOpIsLess = "isLess"
const ConditionOpIsLessEqual = "isLessEqual"

func conditionIs(id, name, value string) *Condition {
	return NewCondition(id, name, ConditionOpIs, []string{value})
}

func conditionIsNot(id, name, value string) *Condition {
	return NewCondition(id, name, ConditionOpIsNot, []string{value})
}

func conditionContains(id, name, value string) *Condition {
	return NewCondition(id, name, ConditionOpContains, []string{value})
}

func conditionDoesNotContains(id, name, value string) *Condition {
	return NewCondition(id, name, ConditionOpDoesNotContain, []string{value})
}

func conditionIsEmpty(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIsEmpty, []string{})
}

func conditionIsNotEmpty(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIsNotEmpty, []string{})
}

func conditionIsGreater(id, name, value string) *Condition {
	return NewCondition(id, name, ConditionOpIsGreater, []string{value})
}

func conditionIsGreaterEqual(id, name, value string) *Condition {
	return NewCondition(id, name, ConditionOpIsGreaterEqual, []string{value})
}

func conditionIsLess(id, name, value string) *Condition {
	return NewCondition(id, name, ConditionOpIsLess, []string{value})
}

func conditionIsLessEqual(id, name, value string) *Condition {
	return NewCondition(id, name, ConditionOpIsLessEqual, []string{value})
}

func conditionDateIsToday(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIs, []string{"Today"})
}

func conditionDateIsTomorrow(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIs, []string{"Tomorrow"})
}

func conditionDateIsYesterday(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIs, []string{"Yesterday"})
}

func conditionDateIs(id, name string, time time.Time) *Condition {
	value := fmt.Sprintf("%d", time.UnixMilli())
	return NewCondition(id, name, ConditionOpIs, []string{"ExactDate", value})
}

func conditionDateIsGreaterThanToday(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIsGreater, []string{"Today"})
}

func conditionDateIsGreaterThanTomorrow(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIsGreater, []string{"Tomorrow"})
}

func conditionDateIsGreaterThanYesterday(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIsGreater, []string{"Yesterday"})
}

func conditionDateIsGreater(id, name string, time time.Time) *Condition {
	value := fmt.Sprintf("%d", time.UnixMilli())
	return NewCondition(id, name, ConditionOpIsGreater, []string{"ExactDate", value})
}

func conditionDateIsLessThanToday(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIsLess, []string{"Today"})
}

func conditionDateIsLessThanTomorrow(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIsLess, []string{"Tomorrow"})
}

func conditionDateIsLessThanYesterday(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIsLess, []string{"Yesterday"})
}

func conditionDateIsLess(id, name string, time time.Time) *Condition {
	value := fmt.Sprintf("%d", time.UnixMilli())
	return NewCondition(id, name, ConditionOpIsLess, []string{"ExactDate", value})
}

func conditionDateIsCurrentWeek(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIs, []string{"CurrentWeek"})
}

func conditionDateIsLastWeek(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIs, []string{"LastWeek"})
}

func conditionDateIsCurrentMonth(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIs, []string{"CurrentMonth"})
}

func conditionDateIsLastMonth(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIs, []string{"LastMonth"})
}

func conditionDateIsTheLastWeek(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIs, []string{"TheLastWeek"})
}

func conditionDateTheNextWeek(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIs, []string{"TheNextWeek"})
}

func conditionDateIsTheLastMonth(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIs, []string{"TheLastMonth"})
}

func conditionDateTheNextMonth(id, name string) *Condition {
	return NewCondition(id, name, ConditionOpIs, []string{"TheNextMonth"})
}
