package larkbase

import (
	"fmt"
	"strings"
)

type MediaField struct {
	BaseField
	value []string
}

func (m *MediaField) Type() FieldType {
	return FieldTypeMedia
}

func (m *MediaField) Value() string {
	return strings.Join(m.value, ",")
}

func (m *MediaField) SetValue(v any) error {
	if vv, ok := v.([]string); !ok {
		return fmt.Errorf("value should be []string, actual: %v", v)
	} else {
		m.value = vv
		return nil
	}
}

func (m *MediaField) Build() any {
	tmp := make([]any, len(m.value))
	for i, sel := range m.value {
		tmp[i] = map[string]string{
			"file_token": sel,
		}
	}
	return tmp
}

func (m *MediaField) Parse(v any) IField {
	ret := &MediaField{BaseField: BaseField{name: m.name}}
	if v, ok := v.([]any); ok {
		for _, v2 := range v {
			if v3, ok2 := v2.(map[string]any); ok2 {
				ret.value = append(ret.value, v3["file_token"].(string))
			}
		}
	}
	return ret
}
