package larkbase

import (
	"fmt"
	"strings"
)

type PersonField struct {
	BaseField
	value []string
}

func (p *PersonField) Type() FieldType {
	return FieldTypePerson
}

func (p *PersonField) Value() string {
	return strings.Join(p.value, ",")
}

func (p *PersonField) Build() any {
	return p.value
}

func (p *PersonField) SetValue(v any) error {
	if vv, ok := v.([]string); !ok {
		return fmt.Errorf("value should be []string, actual: %v", v)
	} else {
		p.value = vv
		return nil
	}
}

func (p *PersonField) Parse(v any) IField {
	ret := &PersonField{BaseField: BaseField{name: p.name}}
	if v, ok := v.([]any); ok {
		for _, v2 := range v {
			if v3, ok2 := v2.(map[string]any); ok2 {
				ret.value = append(ret.value, v3["name"].(string))
			}
		}
	}
	return ret
}
