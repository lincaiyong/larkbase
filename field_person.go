package larkbase

import (
	"strings"
)

type PersonField struct {
	BaseField
	value []string
}

func (f *PersonField) Type() FieldType {
	return FieldTypePerson
}

func (f *PersonField) Value() string {
	return strings.Join(f.value, ",")
}

func (f *PersonField) SetValue(v string) {
	f.value = strings.Split(v, ",")
}

func (f *PersonField) Parse(v any) IField {
	ret := &PersonField{BaseField: BaseField{name: f.name}}
	if v, ok := v.([]any); ok {
		for _, v2 := range v {
			if v3, ok2 := v2.(map[string]any); ok2 {
				ret.value = append(ret.value, v3["name"].(string))
			}
		}
	}
	return ret
}
