package larkbase

import (
	"fmt"
)

type TextField struct {
	BaseField
	value string
}

func (t *TextField) Type() FieldType {
	return FieldTypeText
}

func (t *TextField) Value() string {
	return t.value
}

func (t *TextField) SetValue(v any) error {
	if vv, ok := v.(string); !ok {
		return fmt.Errorf("value should be string, actual: %v", v)
	} else {
		t.value = vv
		return nil
	}
}

func (t *TextField) Build() any {
	return t.value
}

func (t *TextField) Parse(v any) IField {
	ret := &TextField{BaseField: BaseField{name: t.name}}
	if v, ok := v.([]any); ok && len(v) == 1 {
		v2 := v[0]
		if v3, ok2 := v2.(map[string]any); ok2 {
			ret.value = v3["text"].(string)
		}
	}
	return ret
}
