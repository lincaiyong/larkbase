package larkbase

import (
	"fmt"
)

type UrlField struct {
	BaseField
	value string
}

func (u *UrlField) Type() FieldType {
	return FieldTypeUrl
}

func (u *UrlField) Value() string {
	return u.value
}

func (u *UrlField) SetValue(v any) error {
	if vv, ok := v.(string); !ok {
		return fmt.Errorf("value should be string, actual: %v", v)
	} else {
		u.value = vv
		return nil
	}
}

func (u *UrlField) Build() any {
	return map[string]any{
		"text": u.value,
		"link": u.value,
	}
}

func (u *UrlField) Parse(v any) IField {
	ret := &UrlField{BaseField: BaseField{name: u.name}}
	if v, ok := v.(map[string]any); ok {
		ret.value = v["link"].(string)
	}
	return ret
}
