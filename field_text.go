package larkbase

import (
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

type TextField struct {
	BaseField
	value string
}

func (f *TextField) Type() FieldType {
	return FieldTypeText
}

func (f *TextField) Value() string {
	return f.value
}

func (f *TextField) SetValue(v string) {
	f.value = v
}

func (f *TextField) Parse(v any) IField {
	ret := &TextField{BaseField: BaseField{name: f.name}}
	if v, ok := v.([]any); ok && len(v) == 1 {
		v2 := v[0]
		if v3, ok2 := v2.(map[string]any); ok2 {
			ret.value = v3["text"].(string)
		}
	}
	return ret
}

func (f *TextField) FilterIs(value ...string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().
		FieldName(f.Name()).
		Operator(FilterTypeIs).
		Value(value).
		Build()
}

func (f *TextField) FilterIsNot(value ...string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().
		FieldName(f.Name()).
		Operator(FilterTypeIsNot).
		Value(value).
		Build()
}

func (f *TextField) FilterContains(v string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().
		FieldName(f.name).
		Operator(FilterTypeContains).
		Value([]string{v}).
		Build()
}

func (f *TextField) FilterDoesNotContains(v string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().
		FieldName(f.name).
		Operator(FilterTypeDoesNotContain).
		Value([]string{v}).
		Build()
}

func (f *TextField) FilterIsEmpty() *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().
		FieldName(f.name).
		Operator(FilterTypeIsEmpty).
		Value([]string{}).
		Build()
}

func (f *TextField) FilterIsNotEmpty() *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().
		FieldName(f.name).
		Operator(FilterTypeIsNotEmpty).
		Value([]string{}).
		Build()
}
