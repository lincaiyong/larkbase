package larkbase

import (
	"fmt"
	"time"
)

type DateField struct {
	BaseField
	value *time.Time
}

func (d *DateField) Type() FieldType {
	return FieldTypeDate
}

func timestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).In(time.FixedZone("CST", 8*60*60))
}

func timestampToDateTimeStr(timestamp int64) string {
	t := timestampToTime(timestamp)
	return t.Format("2006-01-02 15:04:05")
}

func (d *DateField) Value() string {
	if d.value == nil {
		return ""
	}
	return timestampToDateTimeStr(d.value.Unix())
}

func (d *DateField) SetValue(v any) error {
	if vv, ok := v.(time.Time); !ok {
		return fmt.Errorf("value should be time.Time, actual: %v", v)
	} else {
		d.value = &vv
		return nil
	}
}

func (d *DateField) Build() any {
	return d.value.Unix() * 1000
}

func (d *DateField) Parse(v any) IField {
	ret := &DateField{BaseField: BaseField{name: d.name}}
	t := timestampToTime(int64(v.(float64) / 1000))
	ret.value = &t
	return ret
}
