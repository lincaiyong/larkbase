package larkbase

import (
	"encoding/json"
	"fmt"
)

type MultiSelectField struct {
	BaseField
	value []string
}

func (m *MultiSelectField) Type() FieldType {
	return FieldTypeMultiSelect
}

func (m *MultiSelectField) Value() string {
	b, _ := json.Marshal(m.value)
	return string(b)
}

func (m *MultiSelectField) SetValue(v any) error {
	if vv, ok := v.([]string); !ok {
		return fmt.Errorf("value should be []string, actual: %v", v)
	} else {
		m.value = vv
		return nil
	}
}

func (m *MultiSelectField) Build() any {
	tmp := make([]any, len(m.value))
	for i, sel := range m.value {
		tmp[i] = sel
	}
	return tmp
}

func (m *MultiSelectField) Parse(v any) IField {
	ret := &MultiSelectField{BaseField: BaseField{name: m.name}}
	if v, ok := v.([]any); ok {
		for _, v2 := range v {
			ret.value = append(ret.value, v2.(string))
		}
	}
	return ret
}
