package larkbase

import (
	"fmt"
	"time"
)

type UpdatedTimeField struct {
	BaseField
	value *time.Time
}

func (t *UpdatedTimeField) Type() FieldType {
	return FieldTypeUpdatedTime
}

func (t *UpdatedTimeField) Value() string {
	if t.value == nil {
		return ""
	}
	return timestampToDateTimeStr(t.value.Unix())
}

func (t *UpdatedTimeField) SetValue(v any) error {
	if vv, ok := v.(time.Time); !ok {
		return fmt.Errorf("value should be time.Time, actual: %v", v)
	} else {
		t.value = &vv
		return nil
	}
}

func (t *UpdatedTimeField) Build() any {
	return t.value
}

func (t *UpdatedTimeField) Parse(v any) IField {
	ret := &UpdatedTimeField{BaseField: BaseField{name: t.name}}
	tt := timestampToTime(int64(v.(float64) / 1000))
	ret.value = &tt
	return ret
}
