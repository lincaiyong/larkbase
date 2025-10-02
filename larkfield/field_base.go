package larkfield

import (
	"encoding/json"
	"time"
)

type BaseField struct {
	name  string
	type_ string
	value any
	dirty bool
}

func (f *BaseField) Name() string {
	return f.name
}

func (f *BaseField) SetName(name string) {
	f.name = name
}

func (f *BaseField) Type() string {
	return f.type_
}

func (f *BaseField) SetType(type_ string) {
	f.type_ = type_
}

func (f *BaseField) UnderlayValue() any {
	return f.value
}

func (f *BaseField) SetUnderlayValue(value any) {
	if f.value != value {
		f.dirty = true
		f.value = value
	}
}

func (f *BaseField) SetUnderlayValueNoDirty(value any) {
	f.value = value
}

func (f *BaseField) Dirty() bool {
	return f.dirty
}

func (f *BaseField) SetDirty(v bool) {
	f.dirty = v
}

func stringValue(value any) string {
	if value == nil {
		return ""
	}
	if v, ok := value.(string); ok {
		return v
	}
	if v, ok := value.(time.Time); ok {
		return timeToBeijingDateTimeStr(v)
	}
	b, _ := json.Marshal(value)
	return string(b)
}

func (f *BaseField) StringValue() string {
	return stringValue(f.value)
}

func (f *BaseField) Fork() Field {
	panic("should not happen")
}

func (f *BaseField) Parse(_ any) {
	panic("should not happen")
}

func (f *BaseField) Build() any {
	panic("should not happen")
}

type HackBaseField BaseField

func (f HackBaseField) Name() string {
	return f.name
}

func (f HackBaseField) Type() string {
	return f.type_
}

func (f HackBaseField) Value() any {
	return f.value
}

func (f HackBaseField) StringValue() string {
	return stringValue(f.value)
}

func (f HackBaseField) Dirty() bool {
	return f.dirty
}
