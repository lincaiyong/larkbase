package larkbase

import (
	"encoding/json"
)

type MultiSelectField struct {
	BaseField
	value []string
}

func (f *MultiSelectField) Type() FieldType {
	return FieldTypeMultiSelect
}

func (f *MultiSelectField) Value() string {
	b, _ := json.Marshal(f.value)
	return string(b)
}

func (f *MultiSelectField) SetValue(v string) {
	_ = json.Unmarshal([]byte(v), &f.value)
}

func (f *MultiSelectField) Parse(v any) IField {
	ret := &MultiSelectField{BaseField: BaseField{name: f.name}}
	if v, ok := v.([]any); ok {
		for _, v2 := range v {
			ret.value = append(ret.value, v2.(string))
		}
	}
	return ret
}
