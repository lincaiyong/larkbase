package larkbase

import (
	"fmt"
)

type NumberField struct {
	BaseField
	value *float64
}

func (n *NumberField) Type() FieldType {
	return FieldTypeNumber
}

func (n *NumberField) Value() string {
	if n.value == nil {
		return ""
	}
	return fmt.Sprintf("%g", *n.value)
}

func (n *NumberField) SetValue(v any) error {
	if vv, ok := v.(float64); ok {
		n.value = &vv
		return nil
	}
	if vv, ok := v.(int); ok {
		f := float64(vv)
		n.value = &f
		return nil
	}
	return fmt.Errorf("value should be number, actual: %v", v)
}

func (n *NumberField) Build() any {
	return n.value
}

func (n *NumberField) Parse(v any) IField {
	ret := &NumberField{BaseField: BaseField{name: n.name}}
	f := v.(float64)
	ret.value = &f
	return ret
}
