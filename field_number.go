package larkbase

import (
	"fmt"
	"strconv"
)

type NumberField struct {
	BaseField
	value *float64
}

func (f *NumberField) Type() FieldType {
	return FieldTypeNumber
}

func (f *NumberField) Value() string {
	if f.value == nil {
		return ""
	}
	return fmt.Sprintf("%g", *f.value)
}

func (f *NumberField) SetValue(v string) {
	n, _ := strconv.ParseFloat(v, 64)
	f.value = &n
}

func (f *NumberField) Parse(v any) IField {
	ret := &NumberField{BaseField: BaseField{name: f.name}}
	n := v.(float64)
	ret.value = &n
	return ret
}
