package larkbase

import (
	"strings"
)

type MediaField struct {
	BaseField
	value []string
}

func (f *MediaField) Type() FieldType {
	return FieldTypeMedia
}

func (f *MediaField) Value() string {
	return strings.Join(f.value, ",")
}

func (f *MediaField) SetValue(v string) {
	f.value = strings.Split(v, ",")
}

func (f *MediaField) Parse(v any) IField {
	ret := &MediaField{BaseField: BaseField{name: f.name}}
	if v, ok := v.([]any); ok {
		for _, v2 := range v {
			if v3, ok2 := v2.(map[string]any); ok2 {
				ret.value = append(ret.value, v3["file_token"].(string))
			}
		}
	}
	return ret
}
