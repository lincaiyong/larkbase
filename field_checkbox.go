package larkbase

import (
	"fmt"
)

type CheckboxField struct {
	BaseField
	value bool
}

func (c *CheckboxField) Type() FieldType {
	return FieldTypeCheckbox
}

func (c *CheckboxField) Value() string {
	if c.value {
		return "true"
	} else {
		return "false"
	}
}

func (c *CheckboxField) SetValue(v any) error {
	if vv, ok := v.(bool); !ok {
		return fmt.Errorf("value should be bool, actual: %v", v)
	} else {
		c.value = vv
		return nil
	}
}

func (c *CheckboxField) Build() any {
	return c.value
}

func (c *CheckboxField) Parse(v any) IField {
	ret := &CheckboxField{BaseField: BaseField{name: c.name}}
	ret.value = v.(bool)
	return ret
}
