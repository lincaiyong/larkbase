package larkbase

import (
	"fmt"
)

type SingleSelectField struct {
	BaseField
	value string
}

func (s *SingleSelectField) Type() FieldType {
	return FieldTypeSingleSelect
}

func (s *SingleSelectField) Value() string {
	return s.value
}

func (s *SingleSelectField) SetValue(v any) error {
	if vv, ok := v.(string); !ok {
		return fmt.Errorf("value should be []string, actual: %v", v)
	} else {
		s.value = vv
		return nil
	}
}

func (s *SingleSelectField) Build() any {
	return s.value
}

func (s *SingleSelectField) Parse(v any) IField {
	ret := &SingleSelectField{BaseField: BaseField{name: s.name}}
	ret.value = v.(string)
	return ret
}
