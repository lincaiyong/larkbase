package larkbase

import (
	"fmt"
)

type AutoNumberField struct {
	BaseField
	value *int
}

func (n *AutoNumberField) Type() FieldType {
	return FieldTypeAutoNumber
}

func (n *AutoNumberField) Value() string {
	if n.value == nil {
		return ""
	}
	return fmt.Sprintf("%d", *n.value)
}

func (n *AutoNumberField) SetValue(v any) error {
	if vv, ok := v.(int); ok {
		n.value = &vv
		return nil
	}
	return fmt.Errorf("value should be number, actual: %v", v)
}

func (n *AutoNumberField) Build() any {
	return n.value
}

func (n *AutoNumberField) Parse(v any) IField {
	ret := &AutoNumberField{BaseField: BaseField{name: n.name}}
	f := v.(int)
	ret.value = &f
	return ret
}
