package larkbase

import (
	"fmt"
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

func (f *TextField) SetValue(v any) error {
	if vv, ok := v.(string); !ok {
		return fmt.Errorf("value should be string, actual: %v", v)
	} else {
		f.value = vv
		return nil
	}
}

func (f *TextField) Build() any {
	return f.value
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

func (f *TextField) Contains(v string) *larkbitable.Condition {
	return larkbitable.NewConditionBuilder().
		FieldName(f.name).
		Operator(FilterTypeContains).
		Value([]string{v}).
		Build()
}
